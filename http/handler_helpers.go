// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/vault"
)

// RouteContext holds shared state needed by all subsystem route registrars.
//
// Internals that previously lived undocumented inside handlerWithSettings:
//
//   - RecoveryMode (Props.RecoveryMode): when true, Vault exposes only the
//     recovery raw/generate-recovery-token surface via registerRecoveryRoutes.
//     Normal sys/UI/telemetry registrars are not used.
//
//   - ChrootNamespace: optional listener chroot that prefixes logical path
//     resolution inside handleLogical* helpers. Empty means no chroot.
//
//   - OperatorNamespace: canonicalized operator namespace path (trailing slash)
//     prepended into /v1/{ns}sys/... via SysPath. Empty string means root.
//
// Registrars must use SysPath (or these fields) instead of recomputing
// namespace prefixes ad hoc.
type RouteContext struct {
	Core              *vault.Core
	Mux               *http.ServeMux
	OperatorNamespace string
	ChrootNamespace   string
	Settings          handlerSettings
	Props             *vault.HandlerProperties
}

// NewRouteContext extracts the common namespace-prefix and chroot
// injection values from HandlerProperties.
func NewRouteContext(props *vault.HandlerProperties, settings handlerSettings) *RouteContext {
	core := props.Core
	var chrootNamespace string
	if props.ListenerConfig != nil {
		chrootNamespace = props.ListenerConfig.ChrootNamespace
	}
	operatorNamespace := ""
	if p := core.OperatorNamespacePath(); p != "" {
		operatorNamespace = namespace.Canonicalize(p)
	}
	return &RouteContext{
		Core:              core,
		Mux:               http.NewServeMux(),
		OperatorNamespace: operatorNamespace,
		ChrootNamespace:   chrootNamespace,
		Settings:          settings,
		Props:             props,
	}
}

// SysPath returns the versioned, namespace-prefixed path for a sys endpoint.
// Example: SysPath("init") → "/v1/sys/init" or "/v1/admin/sys/init".
func (rc *RouteContext) SysPath(suffix string) string {
	return "/v1/" + rc.OperatorNamespace + "sys/" + suffix
}
