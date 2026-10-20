// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerSecretsRoutes is a named registrar for secrets-engine routes.
// Currently delegates to the catch-all handleLogical; will be expanded
// as individual secrets engines are decomposed.
func registerSecretsRoutes(rc *RouteContext) {
	// Secrets engine routes are handled by the radix-tree router
	// via the catch-all /v1/ handler registered in registerSystemRoutes.
	// This function exists as a named extension point for the Strangler Fig
	// decomposition and is a no-op until specific routes are extracted.
}
