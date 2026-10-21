// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerIdentityRoutes is the identity subsystem's registrar call site.
// Identity paths ride the /v1/ catch-all owned by registerSecretsRoutes until
// dedicated identity routes are extracted — so secrets must register first.
func registerIdentityRoutes(rc *RouteContext) {
	if !rc.logicalAPICatchAllRegistered {
		panic("registerIdentityRoutes: registerSecretsRoutes must run first")
	}
}
