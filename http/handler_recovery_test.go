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
	"go.uber.org/atomic"
)

func TestRegisterRecoveryRoutes_OnlyThreeRoutes(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		RecoveryMode:   true,
		RecoveryToken:  atomic.NewString("test-recovery-token"),
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerRecoveryRoutes(rc)

	expected := []string{
		"/v1/sys/raw/",
		"/v1/sys/generate-recovery-token/attempt",
		"/v1/sys/generate-recovery-token/update",
	}
	routes := extractRegisteredRoutes(rc.Mux)
	require.ElementsMatch(t, expected, routes)

	for _, path := range expected {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := rc.Mux.Handler(req)
		require.NotEmpty(t, pattern, "expected recovery route %s to be registered", path)
	}

	// Non-recovery paths should not match a dedicated recovery registration
	req := httptest.NewRequest(http.MethodGet, "/v1/sys/health", nil)
	h, pattern := rc.Mux.Handler(req)
	require.Equal(t, "", pattern)
	require.NotNil(t, h) // default NotFound handler
}
