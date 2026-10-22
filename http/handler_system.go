// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// registerSystemRoutes registers /v1/sys/ routes excluding the catch-alls
// (see RegisterLogicalRoutes) and including seal, auth/policy top overrides,
// conditional generate-root/rekey, and raft endpoints.
func registerSystemRoutes(rc *RouteContext) {
	core := rc.Core
	lc := rc.Props.ListenerConfig

	// Non-forwarded paths
	rc.Mux.Handle(rc.SysPath("config/state/"),
		handleLogicalNoForward(core, rc.ChrootNamespace))
	rc.Mux.Handle(rc.SysPath("host-info"),
		handleLogicalNoForward(core, rc.ChrootNamespace))

	rc.Mux.Handle(rc.SysPath("init"), handleSysInit(core))

	var sealOpts []ListenerConfigOption
	if lc != nil {
		sealOpts = []ListenerConfigOption{
			WithRedactClusterName(lc.RedactClusterName),
			WithRedactVersion(lc.RedactVersion),
		}
	}
	RegisterSealRoutes(rc.Mux, core, sealOpts...)

	rc.Mux.Handle(rc.SysPath("leader"), handleSysLeader(core,
		WithRedactAddresses(lc != nil && lc.RedactAddresses)))
	rc.Mux.Handle(rc.SysPath("health"), handleSysHealth(core,
		WithRedactClusterName(lc != nil && lc.RedactClusterName),
		WithRedactVersion(lc != nil && lc.RedactVersion)))
	rc.Mux.Handle(rc.SysPath("monitor"),
		handleLogicalNoForward(core, rc.ChrootNamespace))

	// Conditional generate-root
	if rc.Settings.unauthGenerateRoot {
		rc.Mux.Handle(rc.SysPath("generate-root/attempt"),
			handleRequestForwarding(core,
				handleAuditNonLogical(core,
					handleSysGenerateRootAttempt(core, vault.GenerateStandardRootTokenStrategy))))
		rc.Mux.Handle(rc.SysPath("generate-root/update"),
			handleRequestForwarding(core,
				handleAuditNonLogical(core,
					handleSysGenerateRootUpdate(core, vault.GenerateStandardRootTokenStrategy))))
	}

	// Conditional rekey
	if rc.Settings.unauthRekey {
		rc.Mux.Handle(rc.SysPath("rekey/init"),
			handleRequestForwarding(core, handleSysRekeyInit(core, false)))
		rc.Mux.Handle(rc.SysPath("rekey/update"),
			handleRequestForwarding(core, handleSysRekeyUpdate(core, false)))
		rc.Mux.Handle(rc.SysPath("rekey/verify"),
			handleRequestForwarding(core, handleSysRekeyVerify(core, false)))
		rc.Mux.Handle(rc.SysPath("rekey-recovery-key/init"),
			handleRequestForwarding(core, handleSysRekeyInit(core, true)))
		rc.Mux.Handle(rc.SysPath("rekey-recovery-key/update"),
			handleRequestForwarding(core, handleSysRekeyUpdate(core, true)))
		rc.Mux.Handle(rc.SysPath("rekey-recovery-key/verify"),
			handleRequestForwarding(core, handleSysRekeyVerify(core, true)))
	}

	// Raft and feature flags
	rc.Mux.Handle(rc.SysPath("storage/raft/bootstrap"),
		handleSysRaftBootstrap(core))
	rc.Mux.Handle(rc.SysPath("storage/raft/join"),
		handleSysRaftJoin(core))
	rc.Mux.Handle(rc.SysPath("internal/ui/feature-flags"),
		handleSysInternalFeatureFlags(core))

	RegisterSysAuthRoutes(rc.Mux, core, rc.ChrootNamespace)
	RegisterSysPolicyRoutes(rc.Mux, core, rc.ChrootNamespace)
	registerRemainingTopRouteOverrides(rc.Mux, core, rc.ChrootNamespace)
}

// registerRemainingTopRouteOverrides registers TopRouteOverrides entries not
// already handled by RegisterSysAuthRoutes or RegisterSysPolicyRoutes.
func registerRemainingTopRouteOverrides(mux *http.ServeMux, core *vault.Core, chrootNamespace string) {
	for _, path := range TopRouteOverrides {
		if isSysAuthTopRoute(path) || isSysPolicyTopRoute(path) {
			continue
		}
		mux.Handle(path,
			handleRequestForwarding(core, handleLogicalWithInjector(core, chrootNamespace)))
	}
}
