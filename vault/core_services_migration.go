// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// Core service-interface migration guide (WO-118 / Strangler Fig on Core)
//
// Historical access pattern
// -------------------------
// Call sites reach into *Core fields and methods directly (e.g. c.router,
// c.expiration, c.seal, audit backends). That concentrates domain knowledge
// onto the Core "god object" and raises cognitive load for every change.
//
// Target access pattern
// ---------------------
// Domain behavior is expressed behind narrow interfaces owned by Core, for
// example (names illustrative of the decomposition epic):
//
//	AuditService, MountManager, ExpirationManager, SealManager
//
// New code should depend on the interface it needs, not on *Core wholesale.
// During the transition, Core may expose both:
//
//   - legacy direct fields / methods (existing callers), and
//   - accessors that return the new interfaces (preferred for new callers).
//
// Backward-compatible aliases (WO-124 style) may keep old names compiling
// while implementations move. Prefer the interface accessor in new code;
// do not add new direct field reads when an interface exists.
//
// What to do when you touch Core
// ------------------------------
//  1. Prefer injecting or calling the smallest interface that satisfies the need.
//  2. If only a *Core field exists today, keep the change local; do not spread
//     new direct Core coupling into unrelated packages.
//  3. When extracting a service, add a short comment on the Core field pointing
//     at the interface and mark the direct field as transitional if aliased.
//  4. Full narrative: docs/adr/0002-core-service-interfaces-migration.md
//
// This file intentionally contains no executable code — it exists so `go doc`
// and code search surface the migration rules next to package vault.
