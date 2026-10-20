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

func TestRegisterSystemRoutes_ExpectedSysPaths(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	settings := handlerSettings{
		unauthGenerateRoot: true,
		unauthRekey:        true,
	}
	rc := NewRouteContext(props, settings)
	registerSystemRoutes(rc)

	paths := []string{
		"/v1/sys/init",
		"/v1/sys/seal-status",
		"/v1/sys/seal-backend-status",
		"/v1/sys/seal",
		"/v1/sys/step-down",
		"/v1/sys/unseal",
		"/v1/sys/leader",
		"/v1/sys/health",
		"/v1/sys/monitor",
		"/v1/sys/generate-root/attempt",
		"/v1/sys/generate-root/update",
		"/v1/sys/rekey/init",
		"/v1/sys/rekey/update",
		"/v1/sys/rekey/verify",
		"/v1/sys/rekey-recovery-key/init",
		"/v1/sys/rekey-recovery-key/update",
		"/v1/sys/rekey-recovery-key/verify",
		"/v1/sys/storage/raft/bootstrap",
		"/v1/sys/storage/raft/join",
		"/v1/sys/internal/ui/feature-flags",
		"/v1/sys/",
		"/v1/",
	}
	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := rc.Mux.Handler(req)
		require.NotEmpty(t, pattern, "expected sys path %s to be registered", path)
	}
}

func TestRegisterSystemRoutes_ConditionalGenerateRootDisabled(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{unauthGenerateRoot: false})
	registerSystemRoutes(rc)

	routes := extractRegisteredRoutes(rc.Mux)
	for _, r := range routes {
		require.NotEqual(t, "/v1/sys/generate-root/attempt", r)
		require.NotEqual(t, "/v1/sys/generate-root/update", r)
	}
}
