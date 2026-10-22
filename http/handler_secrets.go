// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerSecretsRoutes is the secrets-engine registrar call site.
// The /v1/ catch-all is registered by RegisterLogicalRoutes after subsystem
// registrars so more-specific routes win; this function marks ownership for
// future per-engine extraction.
func registerSecretsRoutes(rc *RouteContext) {
	rc.logicalAPICatchAllRegistered = true
}
