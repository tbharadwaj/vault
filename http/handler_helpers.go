// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/vault"
)

// RouteContext is the shared input for every subsystem route registrar.
//
//	RecoveryMode (Props) — only recovery routes; see registerRecoveryRoutes
//	ChrootNamespace      — listener chroot for handleLogical* ("" = none)
//	OperatorNamespace    — canonical "ns/" prefix for SysPath ("" = root)
//	logicalAPICatchAllRegistered — set by registerSecretsRoutes for /v1/
type RouteContext struct {
	Core              *vault.Core
	Mux               *http.ServeMux
	OperatorNamespace string
	ChrootNamespace   string
	Settings          handlerSettings
	Props             *vault.HandlerProperties

	logicalAPICatchAllRegistered bool
}

// NewRouteContext builds a RouteContext from listener/core handler properties.
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

// SysPath returns "/v1/{operatorNamespace}sys/{suffix}".
func (rc *RouteContext) SysPath(suffix string) string {
	return "/v1/" + rc.OperatorNamespace + "sys/" + suffix
}
