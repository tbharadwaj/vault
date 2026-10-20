---
title: TypeScript migration order
---

# TypeScript migration order (cognitive-load aware)

The UI still mixes JavaScript and TypeScript. During a multi-phase conversion,
**mixed JS/TS is itself a cognitive tax**. To recover ForgeScore Cognitive Load
sooner, convert **shared, high-fan-in modules first**, then leaf routes.

## Recommended order (reorder vs route-first phasing)

| Phase | Target | Why first |
|-------|--------|-----------|
| **0** | `ui/app/utils/**`, shared helpers, constants | Imported by nearly every route/service |
| **1** | `ui/app/services/**`, adapters, serializers, shared resources | Cross-cutting runtime surface |
| **2** | Shared form / model infrastructure (`forms/`, model helpers) | Used by many engines |
| **3** | High-traffic product routes (auth, KV, policy, identity) | User-facing paths once foundations are typed |
| **4** | Remaining engines and leaf routes | Lowest shared fan-in |

### Mapping from prior story IDs

If an epic sequenced identity → KV/PKI → auth/policy → remaining (e.g.
WO-120 → WO-123 → WO-125 → WO-127), **keep those as delivery batches for
product areas**, but **front-load Phase 0–2** (utils/services) before or in
parallel with the first route batch so converters are not writing TS against
untyped shared imports.

## Rules of engagement

1. Prefer `.ts` for any **new** module under `ui/app`.
2. When touching a `.js` file that sits in Phase 0–2, convert it (or extract a
   typed helper) in the same change when practical.
3. Do not expand dual-paradigm APIs (exporting both JS and TS twins) without a
   short removal plan.
4. Ember patterns and form docs still apply — see [Forms](./forms.md),
   [Models](./models.md), [Routing](./routing.md).

## Related

- Contributor [language map](../../CONTRIBUTING.md#language-map)
- Cognitive Load finding: JS→TS dual-paradigm tax
