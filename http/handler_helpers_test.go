// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestNewRouteContext_ChrootNamespace(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			ChrootNamespace: "ns1/",
		},
	}
	rc := NewRouteContext(props, handlerSettings{})
	require.Equal(t, "ns1/", rc.ChrootNamespace)
}

func TestNewRouteContext_EmptyChrootNamespace(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	require.Equal(t, "", rc.ChrootNamespace)
}

func TestRouteContext_SysPath_NoOperatorNamespace(t *testing.T) {
	rc := &RouteContext{}
	require.Equal(t, "/v1/sys/init", rc.SysPath("init"))
}

func TestRouteContext_SysPath_SingleOperatorNamespace(t *testing.T) {
	rc := &RouteContext{OperatorNamespace: "admin/"}
	require.Equal(t, "/v1/admin/sys/init", rc.SysPath("init"))
}

func TestRouteContext_SysPath_NestedOperatorNamespace(t *testing.T) {
	rc := &RouteContext{OperatorNamespace: "parent/child/"}
	require.Equal(t, "/v1/parent/child/sys/init", rc.SysPath("init"))
}
