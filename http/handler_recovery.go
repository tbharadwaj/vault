// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"github.com/hashicorp/vault/vault"
)

// registerRecoveryRoutes registers the minimal set of routes available
// when Vault is running in recovery mode.
func registerRecoveryRoutes(rc *RouteContext) {
	raw := vault.NewRawBackend(rc.Core)
	strategy := vault.GenerateRecoveryTokenStrategy(rc.Props.RecoveryToken)
	rc.Mux.Handle("/v1/sys/raw/", handleLogicalRecovery(raw, rc.Props.RecoveryToken))
	rc.Mux.Handle("/v1/sys/generate-recovery-token/attempt",
		handleSysGenerateRootAttempt(rc.Core, strategy))
	rc.Mux.Handle("/v1/sys/generate-recovery-token/update",
		handleSysGenerateRootUpdate(rc.Core, strategy))
}
