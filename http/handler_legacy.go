// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"net/http/pprof"

	"github.com/hashicorp/vault/vault"
	gziphandler "github.com/klauspost/compress/gzhttp"
)

// registerLegacyRoutes is the transitional fallback used while
// FeatureFlagHTTPHandlerDecomposition is off. Prefer the decomposed
// registrars for new route work; do not extend this function except to
// keep parity until the flag is the default.
func registerLegacyRoutes(rc *RouteContext) {
	core := rc.Core
	mux := rc.Mux
	chrootNamespace := rc.ChrootNamespace
	operatorNamespace := rc.OperatorNamespace
	props := rc.Props
	settings := rc.Settings
	lc := props.ListenerConfig

	// Handle non-forwarded paths
	mux.Handle("/v1/"+operatorNamespace+"sys/config/state/", handleLogicalNoForward(core, chrootNamespace))
	mux.Handle("/v1/"+operatorNamespace+"sys/host-info", handleLogicalNoForward(core, chrootNamespace))

	mux.Handle("/v1/"+operatorNamespace+"sys/init", handleSysInit(core))

	var sealOpts []ListenerConfigOption
	if lc != nil {
		sealOpts = []ListenerConfigOption{
			WithRedactClusterName(lc.RedactClusterName),
			WithRedactVersion(lc.RedactVersion),
		}
	}
	RegisterSealRoutes(mux, core, sealOpts...)

	mux.Handle("/v1/"+operatorNamespace+"sys/leader", handleSysLeader(core,
		WithRedactAddresses(lc != nil && lc.RedactAddresses)))
	mux.Handle("/v1/"+operatorNamespace+"sys/health", handleSysHealth(core,
		WithRedactClusterName(lc != nil && lc.RedactClusterName),
		WithRedactVersion(lc != nil && lc.RedactVersion)))
	mux.Handle("/v1/"+operatorNamespace+"sys/monitor", handleLogicalNoForward(core, chrootNamespace))

	// Register generate-root endpoints as unauthenticated handlers only if unauthGenerateRoot is true.
	// When false, these endpoints will be handled by the sys backend as authenticated endpoints.
	if settings.unauthGenerateRoot {
		mux.Handle("/v1/"+operatorNamespace+"sys/generate-root/attempt", handleRequestForwarding(core,
			handleAuditNonLogical(core, handleSysGenerateRootAttempt(core, vault.GenerateStandardRootTokenStrategy))))
		mux.Handle("/v1/"+operatorNamespace+"sys/generate-root/update", handleRequestForwarding(core,
			handleAuditNonLogical(core, handleSysGenerateRootUpdate(core, vault.GenerateStandardRootTokenStrategy))))
	}

	// Register rekey endpoints as unauthenticated handlers only if unauthRekey is true.
	// When false (the default), these endpoints will be handled by the sys backend as authenticated endpoints.
	if settings.unauthRekey {
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey/init", handleRequestForwarding(core, handleSysRekeyInit(core, false)))
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, false)))
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey/verify", handleRequestForwarding(core, handleSysRekeyVerify(core, false)))
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey-recovery-key/init", handleRequestForwarding(core, handleSysRekeyInit(core, true)))
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey-recovery-key/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, true)))
		mux.Handle("/v1/"+operatorNamespace+"sys/rekey-recovery-key/verify", handleRequestForwarding(core, handleSysRekeyVerify(core, true)))
	}

	mux.Handle("/v1/"+operatorNamespace+"sys/storage/raft/bootstrap", handleSysRaftBootstrap(core))
	mux.Handle("/v1/"+operatorNamespace+"sys/storage/raft/join", handleSysRaftJoin(core))
	mux.Handle("/v1/"+operatorNamespace+"sys/internal/ui/feature-flags", handleSysInternalFeatureFlags(core))

	RegisterSysAuthRoutes(mux, core, chrootNamespace)
	RegisterSysPolicyRoutes(mux, core, chrootNamespace)
	registerRemainingTopRouteOverrides(mux, core, chrootNamespace)
	RegisterLogicalRoutes(mux, core, chrootNamespace)

	if core.UIEnabled() {
		if uiBuiltIn {
			mux.Handle("/ui/", http.StripPrefix("/ui/", gziphandler.GzipHandler(handleUIHeaders(core, handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()}))))))
			mux.Handle("/robots.txt", gziphandler.GzipHandler(handleUIHeaders(core, handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()})))))
		} else {
			mux.Handle("/ui/", handleUIHeaders(core, handleUIStub()))
		}
		mux.Handle("/ui", handleUIRedirect())
		mux.Handle("/", handleUIRedirect())

	}

	// Register metrics path without authentication if enabled
	if props.ListenerConfig != nil && props.ListenerConfig.Telemetry.UnauthenticatedMetricsAccess {
		mux.Handle("/v1/sys/metrics", handleMetricsUnauthenticated(core))
	} else {
		mux.Handle("/v1/sys/metrics", handleLogicalNoForward(core, chrootNamespace))
	}

	if props.ListenerConfig != nil && props.ListenerConfig.Profiling.UnauthenticatedPProfAccess {
		for _, name := range []string{"goroutine", "threadcreate", "heap", "allocs", "block", "mutex"} {
			mux.Handle("/v1/sys/pprof/"+name, pprof.Handler(name))
		}
		mux.Handle("/v1/sys/pprof/", http.HandlerFunc(pprof.Index))
		mux.Handle("/v1/sys/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/v1/sys/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/v1/sys/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/v1/sys/pprof/trace", http.HandlerFunc(pprof.Trace))
	} else {
		mux.Handle("/v1/sys/pprof/", handleLogicalNoForward(core, chrootNamespace))
	}

	if props.ListenerConfig != nil && props.ListenerConfig.InFlightRequestLogging.UnauthenticatedInFlightAccess {
		mux.Handle("/v1/sys/in-flight-req", handleUnAuthenticatedInFlightRequest(core))
	} else {
		mux.Handle("/v1/sys/in-flight-req", handleLogicalNoForward(core, chrootNamespace))
	}
	if settings.unauthDROperationToken {
		entDROperationRoutes(mux, core)
	}
}
