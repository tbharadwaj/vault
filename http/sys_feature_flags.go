// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

type FeatureFlagsResponse struct {
	FeatureFlags []string `json:"feature_flags"`
}

// FeatureFlag_EnvVariables lists client/UI environment hints exposed at
// sys/internal/ui/feature-flags. This is NOT Vault's server-side gating
// mechanism — do not add engineering or product toggles here.
// See docs/adr/0001-feature-flag-mechanism-boundaries.md.
var FeatureFlag_EnvVariables = [...]string{
	"VAULT_CLOUD_ADMIN_NAMESPACE",
}

func featureFlagIsSet(name string) bool {
	switch os.Getenv(name) {
	case "", "0":
		return false
	default:
		return true
	}
}

func handleSysInternalFeatureFlags(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			break
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}

		response := &FeatureFlagsResponse{}

		for _, f := range FeatureFlag_EnvVariables {
			if featureFlagIsSet(f) {
				response.FeatureFlags = append(response.FeatureFlags, f)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		defer logical.IncrementResponseStatusCodeMetric(http.StatusOK)
		w.WriteHeader(http.StatusOK)

		// Generate the response
		enc := json.NewEncoder(w)
		enc.Encode(response)
	})
}
