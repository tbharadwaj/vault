// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerSecretsRoutes owns the secrets-engine HTTP surface.
// Mount-specific paths are dispatched by the radix router behind the shared
// /v1/ logical catch-all (see RouteContext.registerV1LogicalCatchAll).
func registerSecretsRoutes(rc *RouteContext) {
	rc.registerV1LogicalCatchAll()
}
