// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/vault"
)

// RouteContext holds shared state needed by all subsystem route registrars.
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
func (rc *RouteContext) SysPath(suffix string) string {
	return "/v1/" + rc.OperatorNamespace + "sys/" + suffix
}
