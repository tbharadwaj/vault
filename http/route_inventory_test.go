// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func testHandlerProperties(t *testing.T) *vault.HandlerProperties {
	t.Helper()
	core, _, _ := vault.TestCoreUnsealed(t)
	return &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
}

// extractRegisteredRoutes enumerates patterns registered on a ServeMux by
// probing a comprehensive candidate path list (Go 1.22+ ServeMux no longer
// exposes a reflectable patterns slice).
func extractRegisteredRoutes(mux *http.ServeMux) []string {
	candidates := routeInventoryCandidates()
	seen := make(map[string]struct{})
	var routes []string
	for _, c := range candidates {
		req := httptest.NewRequest(http.MethodGet, c, nil)
		_, pattern := mux.Handler(req)
		if pattern == "" {
			continue
		}
		if _, ok := seen[pattern]; ok {
			continue
		}
		seen[pattern] = struct{}{}
		routes = append(routes, pattern)
	}
	sort.Strings(routes)
	return routes
}

func routeInventoryCandidates() []string {
	candidates := []string{
		"/",
		"/ui",
		"/ui/",
		"/robots.txt",
		"/v1/",
		"/v1/sys/",
		"/v1/sys/config/state/",
		"/v1/sys/host-info",
		"/v1/sys/raw/",
		"/v1/sys/generate-recovery-token/attempt",
		"/v1/sys/generate-recovery-token/update",
		"/v1/sys/pprof/cmdline",
		"/v1/sys/pprof/profile",
		"/v1/sys/pprof/symbol",
		"/v1/sys/pprof/trace",
	}
	candidates = append(candidates, UnauthenticatedRoutes...)
	candidates = append(candidates, injectDataIntoTopRoutes...)
	for _, name := range []string{"goroutine", "threadcreate", "heap", "allocs", "block", "mutex"} {
		candidates = append(candidates, "/v1/sys/pprof/"+name)
	}
	return candidates
}

func TestRouteInventoryParity(t *testing.T) {
	props := testHandlerProperties(t)
	settings := getHandlerSettings(props.Core)

	// Register legacy routes
	legacyRC := NewRouteContext(props, settings)
	registerLegacyRoutes(legacyRC)
	legacyRoutes := extractRegisteredRoutes(legacyRC.Mux)

	// Register decomposed routes
	decomposedRC := NewRouteContext(props, settings)
	registerSystemRoutes(decomposedRC)
	registerAuthRoutes(decomposedRC)
	registerSecretsRoutes(decomposedRC)
	registerIdentityRoutes(decomposedRC)
	registerUIRoutes(decomposedRC)
	registerTelemetryRoutes(decomposedRC)
	decomposedRoutes := extractRegisteredRoutes(decomposedRC.Mux)

	// Diff
	if diff := cmp.Diff(legacyRoutes, decomposedRoutes); diff != "" {
		t.Fatalf("Route parity violation (-legacy +decomposed):\n%s", diff)
	}
}
