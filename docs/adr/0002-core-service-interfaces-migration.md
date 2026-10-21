# ADR 0002: Core service interface migration (transitional dual access)

## Status

Accepted (documentation / migration guide for WO-118-style extraction)

## Context

Extracting services from `vault.Core` (audit, mounts, expiration, seal, …)
reduces the god-object surface but, during the Strangler Fig transition,
**both** patterns coexist:

- Legacy: direct `*Core` field and method access
- New: narrow interfaces (`AuditService`, `MountManager`, …) with accessors

Without an explicit migration guide, that dual pattern increases cognitive load
(ForgeScore finding: interfaces added while Core aliases remain).

## Decision

1. **New code** depends on the smallest interface that satisfies the need.
2. **Legacy aliases** may remain for compatibility; do not teach them as the
   preferred API in new modules.
3. Keep this ADR as the migration guide; link it from `vault/feature_flags.go`
   rather than duplicating a doc-only `.go` file in package vault.
4. When a field is aliased, comment it as transitional and point at the
   interface accessor.

## Consequences

- Short-term: two patterns exist (unavoidable during strangulation).
- Documented preference prevents new direct Core coupling from spreading.
- Cognitive load stabilizes because the “which pattern?” question has an answer.

## Related

- HTTP route Strangler Fig: `http/handler.go` + `http/handler_legacy.go`
