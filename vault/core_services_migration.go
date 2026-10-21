// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// Core service-interface migration (reduces Core system gravity over time).
//
// Prefer narrow interfaces (audit, mounts, expiration, seal, …) over new
// direct *Core field coupling. During the transition both patterns exist;
// new code should use the smallest interface that fits.
//
// Full guide: docs/adr/0002-core-service-interfaces-migration.md
