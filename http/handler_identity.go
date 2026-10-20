// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerIdentityRoutes is a named registrar for identity routes.
// Currently delegates to the catch-all handleLogical; will be expanded
// as identity endpoints are decomposed.
func registerIdentityRoutes(rc *RouteContext) {
	// Identity routes are handled by the radix-tree router
	// via the catch-all /v1/ handler registered in registerSystemRoutes.
	// This function exists as a named extension point for the Strangler Fig
	// decomposition and is a no-op until specific routes are extracted.
}
