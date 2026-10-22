// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestHandlerWithSettings_UseDecomposedRoutesDefaultLegacy(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	require.False(t, core.GetUseDecomposedRoutes(),
		"use-decomposed-routes must default false for zero behavioral change")

	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	settings := getHandlerSettings(core)
	require.False(t, settings.useDecomposedRoutes)

	h := handlerWithSettings(props, settings)
	require.NotNil(t, h)
}

func TestHandlerWithSettings_DecomposedRegistrarsParityWithLegacy(t *testing.T) {
	props := testHandlerProperties(t)
	settings := getHandlerSettings(props.Core)

	legacyRC := NewRouteContext(props, settings)
	registerLegacyRoutes(legacyRC)
	legacyRoutes := extractRegisteredRoutes(legacyRC.Mux)

	decomposedRC := NewRouteContext(props, settings)
	registerSystemRoutes(decomposedRC)
	registerAuthRoutes(decomposedRC)
	registerSecretsRoutes(decomposedRC)
	registerIdentityRoutes(decomposedRC)
	registerUIRoutes(decomposedRC)
	registerTelemetryRoutes(decomposedRC)
	RegisterLogicalRoutes(decomposedRC.Mux, decomposedRC.Core, decomposedRC.ChrootNamespace)
	decomposedRoutes := extractRegisteredRoutes(decomposedRC.Mux)

	require.Equal(t, legacyRoutes, decomposedRoutes,
		"decomposed registrars must register the same routes as legacy when flag is enabled")
}

func TestGetHandlerSettings_ReflectsUseDecomposedRoutesToggle(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	require.False(t, getHandlerSettings(core).useDecomposedRoutes)

	require.NoError(t, core.SetUseDecomposedRoutes(context.Background(), true))
	require.True(t, getHandlerSettings(core).useDecomposedRoutes)

	require.NoError(t, core.SetUseDecomposedRoutes(context.Background(), false))
	require.False(t, getHandlerSettings(core).useDecomposedRoutes)
}

// TestHandler_UseDecomposedRoutesRuntimeToggle verifies PUT
// /v1/sys/feature-flags/use-decomposed-routes (sudo) updates Core and that
// the outer handler() rebuilds the mux on the subsequent request.
func TestHandler_UseDecomposedRoutesRuntimeToggle(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	h := handler(props)

	require.False(t, core.GetUseDecomposedRoutes())

	// Baseline request through the live wrapper (legacy mux).
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/sys/seal-status", nil)
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Toggle via the sys feature-flag path (root token has sudo).
	putReq := httptest.NewRequest(http.MethodPut, "/v1/sys/feature-flags/use-decomposed-routes",
		strings.NewReader(`{"enabled":true}`))
	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("X-Vault-Token", token)
	putRec := httptest.NewRecorder()
	h.ServeHTTP(putRec, putReq)
	require.True(t, putRec.Code == http.StatusOK || putRec.Code == http.StatusNoContent,
		"unexpected status %d body=%s", putRec.Code, putRec.Body.String())

	require.True(t, core.GetUseDecomposedRoutes(),
		"PUT must update Core so handler() observes the settings change")

	// Next request must succeed on the rebuilt (decomposed) mux.
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/v1/sys/seal-status", nil)
	h.ServeHTTP(rec2, req2)
	require.Equal(t, http.StatusOK, rec2.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/sys/feature-flags/use-decomposed-routes", nil)
	getReq.Header.Set("X-Vault-Token", token)
	getRec := httptest.NewRecorder()
	h.ServeHTTP(getRec, getReq)
	require.Equal(t, http.StatusOK, getRec.Code)
	var payload struct {
		Data struct {
			Enabled bool `json:"enabled"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &payload))
	require.True(t, payload.Data.Enabled)
}
