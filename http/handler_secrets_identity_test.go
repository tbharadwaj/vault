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

func TestRegisterLogicalRoutes_OwnsCatchAlls(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	rc := NewRouteContext(&vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}, handlerSettings{})

	registerSecretsRoutes(rc)
	registerIdentityRoutes(rc)
	RegisterLogicalRoutes(rc.Mux, rc.Core, rc.ChrootNamespace)

	req := httptest.NewRequest(http.MethodGet, "/v1/secret/foo", nil)
	_, pattern := rc.Mux.Handler(req)
	require.Equal(t, "/v1/", pattern)

	req = httptest.NewRequest(http.MethodGet, "/v1/sys/foo", nil)
	_, pattern = rc.Mux.Handler(req)
	require.Equal(t, "/v1/sys/", pattern)
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
}

func TestRegisterSealRoutes_RegistersFivePaths(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	mux := http.NewServeMux()
	RegisterSealRoutes(mux, core)

	for _, path := range []string{
		"/v1/sys/seal",
		"/v1/sys/unseal",
		"/v1/sys/seal-status",
		"/v1/sys/seal-backend-status",
		"/v1/sys/step-down",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := mux.Handler(req)
		require.NotEmpty(t, pattern, "expected seal route %s", path)
	}
}
