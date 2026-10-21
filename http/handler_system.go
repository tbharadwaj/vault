// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"github.com/hashicorp/vault/vault"
)

// registerSystemRoutes registers all unauthenticated /v1/sys/ routes
// including conditional generate-root, rekey, and raft endpoints.
func registerSystemRoutes(rc *RouteContext) {
	core := rc.Core

	// Non-forwarded paths
	rc.Mux.Handle(rc.SysPath("config/state/"),
		handleLogicalNoForward(core, rc.ChrootNamespace))
	rc.Mux.Handle(rc.SysPath("host-info"),
		handleLogicalNoForward(core, rc.ChrootNamespace))

	rc.Mux.Handle(rc.SysPath("init"), handleSysInit(core))
	rc.Mux.Handle(rc.SysPath("seal-status"), handleSysSealStatus(core,
		WithRedactClusterName(rc.Props.ListenerConfig.RedactClusterName),
		WithRedactVersion(rc.Props.ListenerConfig.RedactVersion)))
	rc.Mux.Handle(rc.SysPath("seal-backend-status"),
		handleSysSealBackendStatus(core))
	rc.Mux.Handle(rc.SysPath("seal"), handleSysSeal(core))
	rc.Mux.Handle(rc.SysPath("step-down"),
		handleRequestForwarding(core, handleSysStepDown(core)))
	rc.Mux.Handle(rc.SysPath("unseal"), handleSysUnseal(core))
	rc.Mux.Handle(rc.SysPath("leader"), handleSysLeader(core,
		WithRedactAddresses(rc.Props.ListenerConfig.RedactAddresses)))
	rc.Mux.Handle(rc.SysPath("health"), handleSysHealth(core,
		WithRedactClusterName(rc.Props.ListenerConfig.RedactClusterName),
		WithRedactVersion(rc.Props.ListenerConfig.RedactVersion)))
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

	// Data-injection top routes
	for _, path := range injectDataIntoTopRoutes {
		rc.Mux.Handle(path,
			handleRequestForwarding(core,
				handleLogicalWithInjector(core, rc.ChrootNamespace)))
	}

	// Catch-all for remaining sys/ paths. Non-sys logical mounts (/v1/) are
	// owned by registerSecretsRoutes / registerIdentityRoutes.
	rc.Mux.Handle(rc.SysPath(""),
		handleRequestForwarding(core, handleLogical(core, rc.ChrootNamespace)))
}
