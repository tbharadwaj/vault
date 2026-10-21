# Cognitive Load remediation (ForgeScore)

This note tracks the documentation/consolidation work that addresses the
Cognitive Load dimension drop (flag sprawl, undocumented handler internals,
JS→TS dual paradigm).

## Strategies applied

### Strategy 1 — Inline documentation for new patterns

| Target | What was added |
|--------|----------------|
| WO-083 / flag map | `vault/feature_flags.go` package guide |
| WO-102 gate | Decision comment in `http/handler.go` + constant `FeatureFlagHTTPHandlerDecomposition` |
| Handler internals | Expanded `RouteContext` docs in `http/handler_helpers.go` |
| WO-118 Core dual access | [ADR 0002](adr/0002-core-service-interfaces-migration.md) |

### Strategy 2 — Flag mechanism consolidation (highest ROI)

| Deliverable | Path |
|-------------|------|
| Decision record | [ADR 0001](adr/0001-feature-flag-mechanism-boundaries.md) |
| Reject duplicate CRUD manager | Documented as rejected; comments on `IsFeatureFlagEnabled`, activationflags |
| Package pointers | `helper/experiments`, `helper/activationflags`, `vault.FeatureFlags`, UI env flags |
| Contributor guide | `CONTRIBUTING.md` → Feature-flag & toggle mechanisms |

### Strategy 3 — JS→TS conversion order

| Deliverable | Path |
|-------------|------|
| Migration order (utils/services first) | `ui/docs/typescript-migration.md` |
| Language map | `CONTRIBUTING.md` → Language map |
| Doc index link | `ui/docs/index.md` |

## Expected Cognitive Load impact

- Removes “which flag system?” ambiguity (primary driver of the 6-point drop).
- Documents recovery / chroot / operator-namespace previously buried in
  `handlerWithSettings`.
- Front-loads shared TS conversion so dual-paradigm tax shrinks on the hottest
  import graph first (story reorder guidance; conversion itself remains epic work).
