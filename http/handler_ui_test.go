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

func TestRegisterUIRoutes_Enabled(t *testing.T) {
	core := vault.TestCoreUI(t, true)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerUIRoutes(rc)

	for _, path := range []string{"/ui/", "/ui", "/"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := rc.Mux.Handler(req)
		require.NotEmpty(t, pattern, "expected UI path %s to be registered when UI enabled", path)
	}
}

func TestRegisterUIRoutes_Disabled(t *testing.T) {
	core := vault.TestCoreUI(t, false)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerUIRoutes(rc)

	require.Empty(t, extractRegisteredRoutes(rc.Mux), "no UI routes should be registered when UI disabled")
}
