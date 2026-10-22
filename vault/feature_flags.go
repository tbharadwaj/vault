// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

// Feature toggles — pick one (details: docs/adr/0001-feature-flag-mechanism-boundaries.md):
//
//	Experiments           → IsExperimentEnabled (startup / pre-GA)
//	Activation flags      → sys/activation-flags (operator opt-in)
//	IsFeatureFlagEnabled  → engineering Strangler Fig / HCL feature_flags
//	use-decomposed-routes → GetUseDecomposedRoutes (HTTP registrar Strangler Fig)
//
// Core interface migration: docs/adr/0002-core-service-interfaces-migration.md

const (
	// FeatureFlagHTTPHandlerDecomposition selects decomposed HTTP registrars
	// over registerLegacyRoutes in http.handlerWithSettings.
	FeatureFlagHTTPHandlerDecomposition = "http-handler-decomposition"

	// UseDecomposedRoutesBarrierPath is the barrier key for the runtime
	// HTTP handler decomposition toggle (sys/feature-flags/use-decomposed-routes).
	UseDecomposedRoutesBarrierPath = "sys/feature-flags/use-decomposed-routes"
)

// GetUseDecomposedRoutes reports whether HTTP route registration should use
// the decomposed registrars. Default is false (legacy monolithic path).
func (c *Core) GetUseDecomposedRoutes() bool {
	if c == nil || c.useDecomposedRoutes == nil {
		return false
	}
	return c.useDecomposedRoutes.Load()
}

// SetUseDecomposedRoutes persists the flag to the barrier and updates the
// in-memory cache. The HTTP handler() wrapper rebuilds its mux when
// getHandlerSettings observes the change on the next request.
func (c *Core) SetUseDecomposedRoutes(ctx context.Context, enabled bool) error {
	if c == nil || c.barrier == nil {
		return fmt.Errorf("core barrier unavailable")
	}
	val := []byte("false")
	if enabled {
		val = []byte("true")
	}
	if err := c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   UseDecomposedRoutesBarrierPath,
		Value: val,
	}); err != nil {
		return fmt.Errorf("failed to persist use-decomposed-routes: %w", err)
	}
	c.useDecomposedRoutes.Store(enabled)
	return nil
}

func (c *Core) loadUseDecomposedRoutesFromBarrier(ctx context.Context) error {
	if c == nil || c.barrier == nil || c.useDecomposedRoutes == nil {
		return nil
	}
	entry, err := c.barrier.Get(ctx, UseDecomposedRoutesBarrierPath)
	if err != nil {
		return err
	}
	if entry == nil || len(entry.Value) == 0 {
		c.useDecomposedRoutes.Store(false)
		return nil
	}
	switch string(entry.Value) {
	case "1", "true", "TRUE", "True":
		c.useDecomposedRoutes.Store(true)
	default:
		c.useDecomposedRoutes.Store(false)
	}
	return nil
}
