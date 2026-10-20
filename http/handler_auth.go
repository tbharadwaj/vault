// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

// registerAuthRoutes registers enterprise-specific authenticated routes
// such as DR operation token endpoints.
func registerAuthRoutes(rc *RouteContext) {
	if rc.Settings.unauthDROperationToken {
		entDROperationRoutes(rc.Mux, rc.Core)
	}
}
