// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// sysAuthTopRoutes are the TopRouteOverrides entries owned by RegisterSysAuthRoutes.
var sysAuthTopRoutes = []string{
	"/v1/sys/auth",
	"/v1/sys/auth/",
	"/v1/sys/audit",
	"/v1/sys/audit/",
	"/v1/sys/audit-hash/",
	"/v1/sys/config/cors",
	"/v1/sys/config/auditing/request-headers/",
	"/v1/sys/config/auditing/request-headers",
}

// RegisterSysAuthRoutes registers auth/audit/cors top-route overrides with
// handleLogicalWithInjector. chrootNamespace is required for listener chroot.
func RegisterSysAuthRoutes(mux *http.ServeMux, core *vault.Core, chrootNamespace string) {
	for _, path := range sysAuthTopRoutes {
		mux.Handle(path, handleRequestForwarding(core, handleLogicalWithInjector(core, chrootNamespace)))
	}
}

func isSysAuthTopRoute(path string) bool {
	for _, p := range sysAuthTopRoutes {
		if p == path {
			return true
		}
	}
	return false
}
