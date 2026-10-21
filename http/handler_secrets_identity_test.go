// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestRegisterSecretsRoutes_OwnsV1CatchAll(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	rc := NewRouteContext(&vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}, handlerSettings{})

	registerSecretsRoutes(rc)
	registerIdentityRoutes(rc)

	req := httptest.NewRequest(http.MethodGet, "/v1/secret/foo", nil)
	_, pattern := rc.Mux.Handler(req)
	require.Equal(t, "/v1/", pattern)
	require.True(t, rc.logicalAPICatchAllRegistered)
}

func TestRegisterIdentityRoutes_RequiresSecretsFirst(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	rc := NewRouteContext(&vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}, handlerSettings{})

	require.Panics(t, func() { registerIdentityRoutes(rc) })
	registerSecretsRoutes(rc)
	require.NotPanics(t, func() { registerIdentityRoutes(rc) })
}

func TestRegisterSystemRoutes_DoesNotOwnV1CatchAll(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	rc := NewRouteContext(&vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}, handlerSettings{})
	registerSystemRoutes(rc)

	req := httptest.NewRequest(http.MethodGet, "/v1/secret/foo", nil)
	_, pattern := rc.Mux.Handler(req)
	require.NotEqual(t, "/v1/", pattern)
	require.False(t, rc.logicalAPICatchAllRegistered)
}
