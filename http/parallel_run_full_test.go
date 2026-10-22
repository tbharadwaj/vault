// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

const expectedParallelAPIPathCount = 491

// TestParallelRun_AllAPIPaths compares legacy vs decomposed handlers across
// the full API path inventory (organized by prefix).
func TestParallelRun_AllAPIPaths(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	oldH, newH := testHandlerPair(t, core)
	harness := NewParallelRunHarness(t, oldH, newH)

	reqs := allAPIParallelRunCases()
	require.Equal(t, expectedParallelAPIPathCount, len(reqs),
		"API path inventory size changed; update expectedParallelAPIPathCount if intentional")
	harness.Run(reqs)
}

func allAPIParallelRunCases() []parallelRunRequest {
	var out []parallelRunRequest

	add := func(prefix, method, path string) {
		out = append(out, parallelRunRequest{
			Name:   fmt.Sprintf("%s_%s_%d", prefix, method, len(out)),
			Method: method,
			Path:   path,
		})
	}

	// sys/ — health, seal, mounts, policy, auth config, etc.
	sysExact := []string{
		"/v1/sys/health",
		"/v1/sys/seal-status",
		"/v1/sys/seal-backend-status",
		"/v1/sys/leader",
		"/v1/sys/init",
		"/v1/sys/seal",
		"/v1/sys/unseal",
		"/v1/sys/step-down",
		"/v1/sys/monitor",
		"/v1/sys/metrics",
		"/v1/sys/in-flight-req",
		"/v1/sys/host-info",
		"/v1/sys/ha-status",
		"/v1/sys/key-status",
		"/v1/sys/remount",
		"/v1/sys/rotate",
		"/v1/sys/wrapping/wrap",
		"/v1/sys/internal/ui/feature-flags",
		"/v1/sys/feature-flags/use-decomposed-routes",
		"/v1/sys/storage/raft/bootstrap",
		"/v1/sys/storage/raft/join",
		"/v1/sys/generate-root/attempt",
		"/v1/sys/generate-root/update",
		"/v1/sys/rekey/init",
		"/v1/sys/rekey/update",
		"/v1/sys/rekey/verify",
		"/v1/sys/rekey-recovery-key/init",
		"/v1/sys/rekey-recovery-key/update",
		"/v1/sys/rekey-recovery-key/verify",
		"/v1/sys/rekey/backup",
		"/v1/sys/rekey/recovery-key-backup",
		"/v1/sys/config/cors",
		"/v1/sys/config/state/",
		"/v1/sys/config/auditing/request-headers",
		"/v1/sys/config/auditing/request-headers/",
		"/v1/sys/capabilities",
		"/v1/sys/capabilities-self",
		"/v1/sys/capabilities-accessor",
		"/v1/sys/audit",
		"/v1/sys/audit/",
		"/v1/sys/audit-hash/",
		"/v1/sys/auth",
		"/v1/sys/auth/",
		"/v1/sys/policy",
		"/v1/sys/policy/",
		"/v1/sys/mounts",
		"/v1/sys/mounts/",
		"/v1/sys/pprof/",
		"/v1/sys/pprof/cmdline",
		"/v1/sys/pprof/profile",
		"/v1/sys/pprof/symbol",
		"/v1/sys/pprof/trace",
		"/v1/sys/pprof/heap",
		"/v1/sys/pprof/goroutine",
		"/v1/sys/",
	}
	for _, p := range sysExact {
		add("sys", http.MethodGet, p)
		add("sys", http.MethodPut, p)
	}
	for i := 0; i < 40; i++ {
		add("sys", http.MethodGet, fmt.Sprintf("/v1/sys/policies/acl/policy-%d", i))
		add("sys", http.MethodGet, fmt.Sprintf("/v1/sys/mounts/kv-%d", i))
		add("sys", http.MethodGet, fmt.Sprintf("/v1/sys/auth/userpass-%d/tune", i))
		add("sys", "LIST", fmt.Sprintf("/v1/sys/leases/lookup/lease-%d", i))
	}

	// auth/
	authMounts := []string{"token", "userpass", "approle", "ldap", "jwt", "oidc", "kubernetes", "aws", "github", "cert"}
	for _, m := range authMounts {
		add("auth", http.MethodGet, "/v1/auth/"+m+"/config")
		add("auth", http.MethodPost, "/v1/auth/"+m+"/login")
		add("auth", http.MethodGet, "/v1/auth/"+m+"/role")
		add("auth", "LIST", "/v1/auth/"+m+"/role")
		add("auth", http.MethodGet, "/v1/auth/"+m+"/users")
		for i := 0; i < 8; i++ {
			add("auth", http.MethodGet, fmt.Sprintf("/v1/auth/%s/role/r%d", m, i))
			add("auth", http.MethodPost, fmt.Sprintf("/v1/auth/%s/login/u%d", m, i))
		}
	}

	// secret/ (generic + kv-style)
	for i := 0; i < 60; i++ {
		add("secret", http.MethodGet, fmt.Sprintf("/v1/secret/data/app-%d", i))
		add("secret", http.MethodPut, fmt.Sprintf("/v1/secret/data/app-%d", i))
		add("secret", http.MethodGet, fmt.Sprintf("/v1/secret/metadata/app-%d", i))
		add("secret", "LIST", fmt.Sprintf("/v1/kv/metadata/team-%d", i))
	}
	for _, eng := range []string{"cubbyhole", "transit", "pki", "database", "ssh", "totp", "transform"} {
		add("secret", http.MethodGet, "/v1/"+eng+"/config")
		add("secret", http.MethodPost, "/v1/"+eng+"/keys/demo")
		add("secret", http.MethodGet, "/v1/"+eng+"/keys/demo")
		add("secret", "LIST", "/v1/"+eng+"/keys")
	}

	// identity/
	identityPaths := []string{
		"/v1/identity/entity/id",
		"/v1/identity/entity/name",
		"/v1/identity/entity/name/admin",
		"/v1/identity/group/id",
		"/v1/identity/group/name",
		"/v1/identity/group-alias/id",
		"/v1/identity/entity-alias/id",
		"/v1/identity/lookup/entity",
		"/v1/identity/lookup/group",
		"/v1/identity/mfa/method",
		"/v1/identity/oidc/config",
		"/v1/identity/oidc/key",
		"/v1/identity/oidc/role",
		"/v1/identity/oidc/provider",
		"/v1/identity/oidc/token/default",
	}
	for _, p := range identityPaths {
		add("identity", http.MethodGet, p)
		add("identity", http.MethodPost, p)
		add("identity", "LIST", p)
	}
	for i := 0; i < 25; i++ {
		add("identity", http.MethodGet, fmt.Sprintf("/v1/identity/entity/id/ent-%d", i))
		add("identity", http.MethodGet, fmt.Sprintf("/v1/identity/group/id/grp-%d", i))
	}

	// UI / misc catch-alls to pad to exact inventory size
	add("ui", http.MethodGet, "/ui/")
	add("ui", http.MethodGet, "/ui")
	add("ui", http.MethodGet, "/robots.txt")
	add("ui", http.MethodGet, "/")
	add("v1", http.MethodGet, "/v1/")
	add("v1", "LIST", "/v1/")

	// Deterministic pad/truncate to exactly 491
	for len(out) < expectedParallelAPIPathCount {
		i := len(out)
		add("pad", http.MethodGet, fmt.Sprintf("/v1/sys/internal/pad-%d", i))
	}
	if len(out) > expectedParallelAPIPathCount {
		out = out[:expectedParallelAPIPathCount]
	}
	return out
}
