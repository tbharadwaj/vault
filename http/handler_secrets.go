// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerSecretsRoutes owns the /v1/ logical catch-all used by secrets
// engines (and other router-mounted backends until they are extracted).
// Mount-specific dispatch remains in Vault's radix router behind this Handle.
func registerSecretsRoutes(rc *RouteContext) {
	if rc.logicalAPICatchAllRegistered {
		return
	}
	rc.Mux.Handle("/v1/",
		handleRequestForwarding(rc.Core, handleLogical(rc.Core, rc.ChrootNamespace)))
	rc.logicalAPICatchAllRegistered = true
}
