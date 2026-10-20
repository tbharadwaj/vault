// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestRegisterAuthRoutes_DROperationTokenGate(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}

	t.Run("disabled", func(t *testing.T) {
		rc := NewRouteContext(props, handlerSettings{unauthDROperationToken: false})
		before := len(extractRegisteredRoutes(rc.Mux))
		registerAuthRoutes(rc)
		after := len(extractRegisteredRoutes(rc.Mux))
		require.Equal(t, before, after, "OSS entDROperationRoutes is a no-op; no routes should be added when disabled")
	})

	t.Run("enabled", func(t *testing.T) {
		rc := NewRouteContext(props, handlerSettings{unauthDROperationToken: true})
		before := len(extractRegisteredRoutes(rc.Mux))
		registerAuthRoutes(rc)
		after := len(extractRegisteredRoutes(rc.Mux))
		// On OSS, entDROperationRoutes is a no-op even when the setting is true.
		require.Equal(t, before, after)
	})
}
