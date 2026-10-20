// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// UnauthenticatedRoutes is the canonical allowlist of paths that bypass
// Vault's token authentication. CI validates that handlerWithSettings
// registers exactly these paths as unauthenticated handlers.
var UnauthenticatedRoutes = []string{
	"/v1/sys/init",
	"/v1/sys/seal-status",
	"/v1/sys/seal-backend-status",
	"/v1/sys/seal",
	"/v1/sys/step-down",
	"/v1/sys/unseal",
	"/v1/sys/leader",
	"/v1/sys/health",
	"/v1/sys/monitor",
	"/v1/sys/generate-root/attempt", // conditional: unauthGenerateRoot
	"/v1/sys/generate-root/update",  // conditional: unauthGenerateRoot
	"/v1/sys/rekey/init",            // conditional: unauthRekey
	"/v1/sys/rekey/update",          // conditional: unauthRekey
	"/v1/sys/rekey/verify",          // conditional: unauthRekey
	"/v1/sys/rekey-recovery-key/init",   // conditional: unauthRekey
	"/v1/sys/rekey-recovery-key/update", // conditional: unauthRekey
	"/v1/sys/rekey-recovery-key/verify", // conditional: unauthRekey
	"/v1/sys/storage/raft/bootstrap",
	"/v1/sys/storage/raft/join",
	"/v1/sys/internal/ui/feature-flags",
	"/v1/sys/metrics",       // conditional: UnauthenticatedMetricsAccess
	"/v1/sys/pprof/",        // conditional: UnauthenticatedPProfAccess
	"/v1/sys/in-flight-req", // conditional: UnauthenticatedInFlightAccess
}
