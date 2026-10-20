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

func TestRegisterTelemetryRoutes_AuthenticatedMode(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerTelemetryRoutes(rc)

	for _, path := range []string{"/v1/sys/metrics", "/v1/sys/pprof/", "/v1/sys/in-flight-req"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := rc.Mux.Handler(req)
		require.NotEmpty(t, pattern, "expected telemetry path %s to be registered", path)
	}

	// Individual pprof profiles are not registered in authenticated mode
	req := httptest.NewRequest(http.MethodGet, "/v1/sys/pprof/heap", nil)
	_, pattern := rc.Mux.Handler(req)
	require.Equal(t, "/v1/sys/pprof/", pattern)
}

func TestRegisterTelemetryRoutes_UnauthenticatedMode(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			Telemetry: configutil.ListenerTelemetry{
				UnauthenticatedMetricsAccess: true,
			},
			Profiling: configutil.ListenerProfiling{
				UnauthenticatedPProfAccess: true,
			},
			InFlightRequestLogging: configutil.ListenerInFlightRequestLogging{
				UnauthenticatedInFlightAccess: true,
			},
		},
	}
	rc := NewRouteContext(props, handlerSettings{})
	registerTelemetryRoutes(rc)

	for _, path := range []string{
		"/v1/sys/metrics",
		"/v1/sys/pprof/",
		"/v1/sys/pprof/heap",
		"/v1/sys/pprof/cmdline",
		"/v1/sys/in-flight-req",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		_, pattern := rc.Mux.Handler(req)
		require.NotEmpty(t, pattern, "expected unauthenticated telemetry path %s to be registered", path)
	}
}
