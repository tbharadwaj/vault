# ADR 0001: Feature-flag mechanism boundaries

## Status

Accepted (consolidation guidance for contributors and epic work)

## Context

ForgeScore / Cognitive Load analysis flagged **three competing flag-like
mechanisms** (plus specialized variants). A feature-flag infrastructure epic
risked adding a fourth (barrier-persisted `FeatureFlagManager` CRUD under
`sys/feature-flags/`) on top of:

| Mechanism | Location | Lifetime | Typical use |
|-----------|----------|----------|-------------|
| **Experiments** | `helper/experiments`, `Core.IsExperimentEnabled` | Process start | Pre-GA / alpha APIs |
| **Activation flags** | `helper/activationflags`, `sys/activation-flags` | Persisted, operator-driven | Opt-in after ship |
| **Config `feature_flags` / `IsFeatureFlagEnabled` | HCL + `Core.IsFlagEnabled` / `IsFeatureFlagEnabled` | Config (enterprise); engineering gates | Strangler Fig / licensed flags |

Specialized (out of scope for general gating):

- `sys/internal/ui/feature-flags` env vars (`http/sys_feature_flags.go`)
- Internal `vault.FeatureFlags` cubbyhole replication state (`persistFeatureFlags`)

Stage reflections noted intent for **config-stored, reloadable** flags, while
some task breakdowns sketched a **barrier-persisted CRUD** subsystem. Those
differ; a single source of truth must be chosen for new work.

## Decision

1. **Do not introduce** a new barrier-persisted FeatureFlagManager CRUD API that
   overlaps activation flags. That path is **rejected** for new development.
2. **Use activation flags** for operator-facing, persisted product opt-in.
3. **Use experiments** only for startup-gated pre-GA API surface.
4. **Use `IsFeatureFlagEnabled` / HCL `feature_flags`** for engineering
   Strangler Fig toggles (example: `http-handler-decomposition`) and enterprise
   config feature flags. Wire new gates through `IsFeatureFlagEnabled`, not a
   parallel manager type.
5. Document the map in `vault/feature_flags.go` and `CONTRIBUTING.md` so the
   decision stays discoverable.

## Consequences

- Cognitive load drops: contributors have one decision table instead of three
  overlapping mental models plus a fourth proposed store.
- HTTP handler decomposition (WO-102) correctly calls
  `IsFeatureFlagEnabled("http-handler-decomposition")`.
- Future epics that need runtime operator toggles extend **activation flags**,
  not a new subsystem.
- Experiments remain intentionally narrow (startup / pre-GA).

## When to use which (cheat sheet)

```
Need pre-GA API behind a startup switch?     → Experiments
Need operator opt-in after the code ships?   → Activation flags
Need engineering cutover / Strangler Fig?    → IsFeatureFlagEnabled
Need UI/cloud env hint for the browser?      → FeatureFlag_EnvVariables only
Need cluster-internal cubbyhole state?       → vault.FeatureFlags storage helpers only
```
