// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathUseDecomposedRoutes = "feature-flags/use-decomposed-routes"
)

func (b *SystemBackend) httpFeatureFlagPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathUseDecomposedRoutes + "$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "feature-flags",
				OperationVerb:   "configure",
				OperationSuffix: "use-decomposed-routes",
			},
			Fields: map[string]*framework.FieldSchema{
				"enabled": {
					Type:        framework.TypeBool,
					Description: "When true, HTTP routing uses decomposed registrars; when false, the legacy monolithic path.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleUseDecomposedRoutesRead,
					Summary:  "Read the use-decomposed-routes feature flag.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"enabled": {Type: framework.TypeBool, Required: true},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.handleUseDecomposedRoutesUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
					Summary:                     "Set the use-decomposed-routes feature flag (requires sudo).",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{Description: "OK"}},
						http.StatusOK:        {{Description: "OK"}},
					},
				},
			},
			HelpSynopsis:    "Toggle HTTP handler decomposition (Strangler Fig) at runtime.",
			HelpDescription: "PUT requires sudo. When enabled, Vault rebuilds the HTTP mux on the next request to use decomposed route registrars.",
		},
	}
}

func (b *SystemBackend) handleUseDecomposedRoutesRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"enabled": b.Core.GetUseDecomposedRoutes(),
		},
	}, nil
}

func (b *SystemBackend) handleUseDecomposedRoutesUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	enabledRaw, ok := d.GetOk("enabled")
	if !ok {
		return logical.ErrorResponse("missing required field: enabled"), nil
	}
	enabled, ok := enabledRaw.(bool)
	if !ok {
		return logical.ErrorResponse("field enabled must be a boolean"), nil
	}
	if err := b.Core.SetUseDecomposedRoutes(ctx, enabled); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"enabled": enabled,
		},
	}, nil
}
