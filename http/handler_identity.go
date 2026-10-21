// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerIdentityRoutes owns the identity HTTP surface.
// Identity mounts share the /v1/ logical catch-all with secrets engines; the
// registration is idempotent so both registrars remain first-class call sites
// without concentrating ownership back into registerSystemRoutes.
func registerIdentityRoutes(rc *RouteContext) {
	rc.registerV1LogicalCatchAll()
}
