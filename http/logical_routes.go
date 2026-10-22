// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// RegisterLogicalRoutes registers the catch-all /v1/sys/ and /v1/ routes that
// forward into handleLogical. The sys catch-all respects the operator namespace
// when configured on core.
func RegisterLogicalRoutes(mux *http.ServeMux, core *vault.Core, chrootNamespace string) {
	mux.Handle(operatorSysBase(core),
		handleRequestForwarding(core, handleLogical(core, chrootNamespace)))
	mux.Handle("/v1/",
		handleRequestForwarding(core, handleLogical(core, chrootNamespace)))
}
