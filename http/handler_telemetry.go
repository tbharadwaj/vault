// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"net/http/pprof"
)

// registerTelemetryRoutes registers metrics, pprof, and in-flight
// request endpoints with conditional authentication bypass.
func registerTelemetryRoutes(rc *RouteContext) {
	lc := rc.Props.ListenerConfig

	// Metrics
	if lc != nil && lc.Telemetry.UnauthenticatedMetricsAccess {
		rc.Mux.Handle("/v1/sys/metrics", handleMetricsUnauthenticated(rc.Core))
	} else {
		rc.Mux.Handle("/v1/sys/metrics",
			handleLogicalNoForward(rc.Core, rc.ChrootNamespace))
	}

	// Pprof
	if lc != nil && lc.Profiling.UnauthenticatedPProfAccess {
		for _, name := range []string{
			"goroutine", "threadcreate", "heap", "allocs", "block", "mutex",
		} {
			rc.Mux.Handle("/v1/sys/pprof/"+name, pprof.Handler(name))
		}
		rc.Mux.Handle("/v1/sys/pprof/", http.HandlerFunc(pprof.Index))
		rc.Mux.Handle("/v1/sys/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		rc.Mux.Handle("/v1/sys/pprof/profile", http.HandlerFunc(pprof.Profile))
		rc.Mux.Handle("/v1/sys/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		rc.Mux.Handle("/v1/sys/pprof/trace", http.HandlerFunc(pprof.Trace))
	} else {
		rc.Mux.Handle("/v1/sys/pprof/",
			handleLogicalNoForward(rc.Core, rc.ChrootNamespace))
	}

	// In-flight requests
	if lc != nil && lc.InFlightRequestLogging.UnauthenticatedInFlightAccess {
		rc.Mux.Handle("/v1/sys/in-flight-req",
			handleUnAuthenticatedInFlightRequest(rc.Core))
	} else {
		rc.Mux.Handle("/v1/sys/in-flight-req",
			handleLogicalNoForward(rc.Core, rc.ChrootNamespace))
	}
}