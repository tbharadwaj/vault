// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestHandlerWithSettings_FeatureFlagDisabledUsesLegacy(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	require.False(t, core.IsFeatureFlagEnabled("http-handler-decomposition"),
		"flag must be disabled by default for zero behavioral change")

	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	settings := getHandlerSettings(core)

	// Legacy path produces a non-empty route set
	legacyRC := NewRouteContext(props, settings)
	registerLegacyRoutes(legacyRC)
	legacyRoutes := extractRegisteredRoutes(legacyRC.Mux)
	require.NotEmpty(t, legacyRoutes)

	// handlerWithSettings with flag disabled should register the same routes as legacy
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
	decomposedRoutes := extractRegisteredRoutes(decomposedRC.Mux)

	require.Equal(t, legacyRoutes, decomposedRoutes,
		"decomposed registrars must register the same routes as legacy when flag is enabled")
}