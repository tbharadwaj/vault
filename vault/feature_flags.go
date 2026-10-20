// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// Package-level guidance: Vault feature toggles (read this first)
//
// Vault historically grew several independent "flag-like" mechanisms. Using the
// wrong one (or inventing a fourth) is a major source of contributor cognitive
// load. Canonical decision record:
//
//	docs/adr/0001-feature-flag-mechanism-boundaries.md
//
// Quick map — choose ONE mechanism for new work:
//
//  1. Experiments (helper/experiments, Core.IsExperimentEnabled)
//     Startup-only, pre-GA / alpha API surface. Enabled via server config or
//     -experiment CLI flags. Not reloadable at runtime; not for gradual
//     product opt-in after GA.
//
//  2. Activation flags (helper/activationflags, sys/activation-flags)
//     One-way (or carefully controlled) operator opt-in for shipping features
//     that must not turn on until the operator is ready. Persisted in storage.
//     Use for product readiness gates (e.g. identity deduplication).
//
//  3. Config feature_flags / IsFlagEnabled / IsFeatureFlagEnabled
//     HCL listener/server feature_flags (enterprise) and the engineering
//     Strangler Fig gate API. Use IsFeatureFlagEnabled for in-process
//     decomposition toggles such as "http-handler-decomposition".
//     Do NOT build a parallel barrier-persisted CRUD FeatureFlagManager
//     subsystem — that duplicates activation flags and is rejected by the ADR.
//
// Specialized (do not extend for general gating):
//
//   - http/sys_feature_flags.go FeatureFlag_EnvVariables: UI/client env hints
//     exposed at sys/internal/ui/feature-flags.
//   - vault.FeatureFlags (persistFeatureFlags/readFeatureFlags): internal
//     cluster replication state (e.g. namespace cubbyholes), not a product toggle.
//
// When unsure: prefer activation flags for operator-facing product opt-in, and
// IsFeatureFlagEnabled for engineering cutovers. See the ADR for full rules.

// Known engineering feature-flag names consumed via IsFeatureFlagEnabled.
const (
	// FeatureFlagHTTPHandlerDecomposition gates the Strangler Fig path in
	// http.handlerWithSettings (decomposed route registrars vs registerLegacyRoutes).
	FeatureFlagHTTPHandlerDecomposition = "http-handler-decomposition"
)
