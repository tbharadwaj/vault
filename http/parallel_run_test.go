// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// parallelRunRequest is a single request exercised against both handlers.
type parallelRunRequest struct {
	Name    string
	Method  string
	Path    string
	Body    []byte
	Headers http.Header
}

// ParallelRunHarness compares responses from a legacy (old) and decomposed (new)
// http.Handler for the same list of requests.
type ParallelRunHarness struct {
	t      *testing.T
	oldH   http.Handler
	newH   http.Handler
	ignore map[string]struct{}
}

// NewParallelRunHarness builds a harness that ignores Date and common
// request-specific headers when comparing responses.
func NewParallelRunHarness(t *testing.T, oldHandler, newHandler http.Handler) *ParallelRunHarness {
	t.Helper()
	return &ParallelRunHarness{
		t:    t,
		oldH: oldHandler,
		newH: newHandler,
		ignore: map[string]struct{}{
			"Date":             {},
			"Content-Length":   {},
			"X-Request-Id":     {},
			"X-Vault-Request":  {},
			"X-Correlation-Id": {},
		},
	}
}

// Run sends each request to both handlers and fails on any divergence.
func (h *ParallelRunHarness) Run(reqs []parallelRunRequest) {
	h.t.Helper()
	for _, req := range reqs {
		h.t.Run(req.Name, func(t *testing.T) {
			oldResp := h.serve(h.oldH, req)
			newResp := h.serve(h.newH, req)
			if diff := h.diffResponses(req, oldResp, newResp); diff != "" {
				t.Fatalf("parallel-run divergence:\n%s", diff)
			}
		})
	}
}

type capturedResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

func (h *ParallelRunHarness) serve(handler http.Handler, pr parallelRunRequest) capturedResponse {
	h.t.Helper()
	var body io.Reader
	if pr.Body != nil {
		body = bytes.NewReader(pr.Body)
	}
	req := httptest.NewRequest(pr.Method, pr.Path, body)
	for k, vals := range pr.Headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return capturedResponse{
		Status: rec.Code,
		Header: rec.Header().Clone(),
		Body:   append([]byte(nil), rec.Body.Bytes()...),
	}
}

func (h *ParallelRunHarness) diffResponses(pr parallelRunRequest, oldR, newR capturedResponse) string {
	var b strings.Builder
	prefix := fmt.Sprintf("%s %s", pr.Method, pr.Path)

	if oldR.Status != newR.Status {
		fmt.Fprintf(&b, "%s\n  field=status old=%d new=%d\n", prefix, oldR.Status, newR.Status)
	}

	oldHeaders := h.filteredHeaders(oldR.Header)
	newHeaders := h.filteredHeaders(newR.Header)
	keys := make(map[string]struct{})
	for k := range oldHeaders {
		keys[k] = struct{}{}
	}
	for k := range newHeaders {
		keys[k] = struct{}{}
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	for _, k := range sorted {
		ov := strings.Join(oldHeaders[k], ", ")
		nv := strings.Join(newHeaders[k], ", ")
		if ov != nv {
			fmt.Fprintf(&b, "%s\n  field=header[%s] old=%q new=%q\n", prefix, k, ov, nv)
		}
	}

	if !bytes.Equal(oldR.Body, newR.Body) {
		off := firstBodyDiffOffset(oldR.Body, newR.Body)
		var oldB, newB byte
		if off < len(oldR.Body) {
			oldB = oldR.Body[off]
		}
		if off < len(newR.Body) {
			newB = newR.Body[off]
		}
		fmt.Fprintf(&b, "%s\n  field=body offset=%d old_len=%d new_len=%d old_byte=0x%02x new_byte=0x%02x\n",
			prefix, off, len(oldR.Body), len(newR.Body), oldB, newB)
	}

	return b.String()
}

func (h *ParallelRunHarness) filteredHeaders(hdr http.Header) http.Header {
	out := make(http.Header)
	for k, vals := range hdr {
		if _, skip := h.ignore[http.CanonicalHeaderKey(k)]; skip {
			continue
		}
		if _, skip := h.ignore[k]; skip {
			continue
		}
		out[k] = append([]string(nil), vals...)
	}
	return out
}

func firstBodyDiffOffset(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func testHandlerPair(t *testing.T, core *vault.Core) (oldH, newH http.Handler) {
	t.Helper()
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	base := getHandlerSettings(core)
	legacy := base
	legacy.useDecomposedRoutes = false
	decomp := base
	decomp.useDecomposedRoutes = true
	return handlerWithSettings(props, legacy), handlerWithSettings(props, decomp)
}

func TestParallelRun_SealEndpoints(t *testing.T) {
	// Shared core: only non-mutating / pre-mutation-failure cases so the first
	// handler does not change state observed by the second.
	core, _, token := vault.TestCoreUnsealed(t)
	oldH, newH := testHandlerPair(t, core)
	harness := NewParallelRunHarness(t, oldH, newH)

	auth := http.Header{"X-Vault-Token": {token}}
	reqs := []parallelRunRequest{
		{Name: "seal-status-get", Method: http.MethodGet, Path: "/v1/sys/seal-status"},
		{Name: "seal-status-get-authed", Method: http.MethodGet, Path: "/v1/sys/seal-status", Headers: auth},
		{Name: "seal-status-post-rejected", Method: http.MethodPost, Path: "/v1/sys/seal-status", Body: []byte(`{}`)},
		{Name: "seal-status-delete-rejected", Method: http.MethodDelete, Path: "/v1/sys/seal-status"},
		{Name: "seal-backend-status-get", Method: http.MethodGet, Path: "/v1/sys/seal-backend-status"},
		{Name: "seal-put-unauthenticated", Method: http.MethodPut, Path: "/v1/sys/seal", Body: []byte(`{}`)},
		{Name: "seal-get-method-not-allowed", Method: http.MethodGet, Path: "/v1/sys/seal"},
		{Name: "unseal-put-invalid-json", Method: http.MethodPut, Path: "/v1/sys/unseal", Body: []byte(`{`)},
		{Name: "unseal-put-missing-key", Method: http.MethodPut, Path: "/v1/sys/unseal", Body: []byte(`{"reset":false}`)},
		{Name: "unseal-put-bad-key", Method: http.MethodPut, Path: "/v1/sys/unseal", Body: []byte(`{"key":"not-a-real-key"}`)},
		{Name: "unseal-get-method-not-allowed", Method: http.MethodGet, Path: "/v1/sys/unseal"},
		{Name: "step-down-put-unauthenticated", Method: http.MethodPut, Path: "/v1/sys/step-down", Body: []byte(`{}`)},
		{Name: "step-down-get-method-not-allowed", Method: http.MethodGet, Path: "/v1/sys/step-down"},
	}
	require.GreaterOrEqual(t, len(reqs), 10)
	harness.Run(reqs)
}
