// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// Feature-toggle map (keep this file short — details live in the ADR):
//
//	Experiments      → helper/experiments, IsExperimentEnabled (startup / pre-GA)
//	Activation flags → helper/activationflags, sys/activation-flags (operator opt-in)
//	IsFeatureFlagEnabled → engineering Strangler Fig / HCL feature_flags
//
// Do not add a fourth barrier-persisted FeatureFlagManager CRUD API.
// Full rules: docs/adr/0001-feature-flag-mechanism-boundaries.md
// Core dual-access migration: docs/adr/0002-core-service-interfaces-migration.md

const (
	// FeatureFlagHTTPHandlerDecomposition gates decomposed HTTP registrars
	// vs registerLegacyRoutes in http.handlerWithSettings.
	FeatureFlagHTTPHandlerDecomposition = "http-handler-decomposition"
)
