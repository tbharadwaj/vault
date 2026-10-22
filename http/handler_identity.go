// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerIdentityRoutes is the identity registrar call site.
// Identity mounts use the /v1/ catch-all from RegisterLogicalRoutes until
// dedicated identity routes are extracted.
func registerIdentityRoutes(rc *RouteContext) {
	if !rc.logicalAPICatchAllRegistered {
		panic("registerIdentityRoutes: registerSecretsRoutes must run first")
	}
}
