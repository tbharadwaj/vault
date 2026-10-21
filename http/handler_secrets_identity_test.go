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

func TestRegisterSecretsAndIdentityRoutes_ShareV1CatchAll(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})

	registerSecretsRoutes(rc)
	registerIdentityRoutes(rc) // idempotent; must not panic on second call

	req := httptest.NewRequest(http.MethodGet, "/v1/secret/foo", nil)
	_, pattern := rc.Mux.Handler(req)
	require.Equal(t, "/v1/", pattern)
	require.True(t, rc.v1LogicalCatchAllSet)
}

func TestRegisterSystemRoutes_DoesNotOwnV1CatchAll(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerSystemRoutes(rc)

	req := httptest.NewRequest(http.MethodGet, "/v1/secret/foo", nil)
	_, pattern := rc.Mux.Handler(req)
	// sys catch-all is /v1/sys/; non-sys logical paths need secrets/identity registrars
	require.NotEqual(t, "/v1/", pattern)
	require.False(t, rc.v1LogicalCatchAllSet)
}
