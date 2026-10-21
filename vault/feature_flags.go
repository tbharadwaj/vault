// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// Feature toggles — pick one (details: docs/adr/0001-feature-flag-mechanism-boundaries.md):
//
//	Experiments           → IsExperimentEnabled (startup / pre-GA)
//	Activation flags      → sys/activation-flags (operator opt-in)
//	IsFeatureFlagEnabled  → engineering Strangler Fig / HCL feature_flags
//
// Core interface migration: docs/adr/0002-core-service-interfaces-migration.md

const (
	// FeatureFlagHTTPHandlerDecomposition selects decomposed HTTP registrars
	// over registerLegacyRoutes in http.handlerWithSettings.
	FeatureFlagHTTPHandlerDecomposition = "http-handler-decomposition"
)
