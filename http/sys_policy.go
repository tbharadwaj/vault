// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// sysPolicyTopRoutes are the TopRouteOverrides entries owned by RegisterSysPolicyRoutes.
var sysPolicyTopRoutes = []string{
	"/v1/sys/policy",
	"/v1/sys/policy/",
}

// RegisterSysPolicyRoutes registers policy top-route overrides with
// handleLogicalWithInjector. Additional policy-specific routes can be added
// here as they are extracted from the catch-all.
func RegisterSysPolicyRoutes(mux *http.ServeMux, core *vault.Core, chrootNamespace string) {
	for _, path := range sysPolicyTopRoutes {
		mux.Handle(path, handleRequestForwarding(core, handleLogicalWithInjector(core, chrootNamespace)))
	}
}

func isSysPolicyTopRoute(path string) bool {
	for _, p := range sysPolicyTopRoutes {
		if p == path {
			return true
		}
	}
	return false
}
