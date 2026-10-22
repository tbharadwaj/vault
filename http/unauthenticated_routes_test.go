// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// registrationSourceFiles are the files that register HTTP routes and may
// introduce new unauthenticated handlers that must appear in UnauthenticatedRoutes.
var registrationSourceFiles = []string{
	"handler.go",
	"handler_legacy.go",
	"handler_system.go",
	"handler_telemetry.go",
	"handler_recovery.go",
	"handler_auth.go",
	"handler_ui.go",
}

// TestUnauthenticatedRouteAllowlist parses route-registration sources and
// verifies every unauthenticated path registration appears in UnauthenticatedRoutes.
func TestUnauthenticatedRouteAllowlist(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	dir := filepath.Dir(thisFile)

	allowlist := make(map[string]struct{}, len(UnauthenticatedRoutes))
	for _, p := range UnauthenticatedRoutes {
		allowlist[p] = struct{}{}
	}

	fset := token.NewFileSet()
	registered := map[string]struct{}{}

	for _, name := range registrationSourceFiles {
		path := filepath.Join(dir, name)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			// File may not exist yet during incremental extraction; skip missing.
			if strings.Contains(err.Error(), "no such file") {
				continue
			}
			t.Fatalf("parse %s: %v", path, err)
		}

		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Handle" || len(call.Args) < 1 {
				return true
			}
			pathStr, ok := extractHandlePath(call.Args[0])
			if !ok || pathStr == "" {
				return true
			}
			normalized := normalizeRegisteredPath(pathStr)
			if !isUnauthenticatedCandidate(normalized) {
				return true
			}
			registered[normalized] = struct{}{}
			if _, found := allowlist[normalized]; !found {
				t.Errorf("unauthenticated path %q registered in %s but missing from UnauthenticatedRoutes allowlist", normalized, name)
			}
			return true
		})
	}

	if len(registered) == 0 {
		t.Fatal("no unauthenticated Handle registrations found; allowlist validation cannot run")
	}
}

// extractHandlePath resolves the first argument to mux.Handle / rc.Mux.Handle into
// a path string with an empty operator namespace.
func extractHandlePath(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return strings.Trim(e.Value, `"`), true
		}
	case *ast.CallExpr:
		// rc.SysPath("suffix") → /v1/sys/suffix
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "SysPath" && len(e.Args) == 1 {
			if lit, ok := e.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				suffix := strings.Trim(lit.Value, `"`)
				return "/v1/sys/" + suffix, true
			}
		}
	case *ast.BinaryExpr:
		// "/v1/" + operatorNamespace + "sys/..." → treat operatorNamespace as empty
		left, lok := extractHandlePath(e.X)
		right, rok := extractHandlePath(e.Y)
		if lok && rok {
			return left + right, true
		}
		if ident, ok := e.X.(*ast.Ident); ok && (ident.Name == "operatorNamespace" || ident.Name == "OperatorNamespace") {
			if rok {
				return right, true
			}
		}
		if ident, ok := e.Y.(*ast.Ident); ok && (ident.Name == "operatorNamespace" || ident.Name == "OperatorNamespace") {
			if lok {
				return left, true
			}
		}
		// "/v1/" + operatorNamespace + "sys/init" is nested BinaryExprs
		if lok && !rok {
			if bin, ok := e.Y.(*ast.BinaryExpr); ok {
				if mid, midOK := bin.X.(*ast.Ident); midOK && (mid.Name == "operatorNamespace" || mid.Name == "OperatorNamespace") {
					if r, rOK := extractHandlePath(bin.Y); rOK {
						return left + r, true
					}
				}
			}
		}
	case *ast.Ident:
		// Skip non-literal identifiers (e.g. path variable in range loops)
		return "", false
	}
	return "", false
}

func normalizeRegisteredPath(p string) string {
	// Collapse accidental double slashes from concat edge cases
	for strings.Contains(p, "//") {
		p = strings.ReplaceAll(p, "//", "/")
	}
	return p
}

// allowlistExemptPaths are /v1/sys/ Handle registrations that are authenticated
// catch-alls, inject-data routes, or recovery-only endpoints and therefore do
// not belong on UnauthenticatedRoutes.
var allowlistExemptPaths = map[string]struct{}{
	"/v1/":        {},
	"/v1/sys/":    {},
	"/v1/sys/config/state/": {},
	"/v1/sys/host-info":     {},
	"/v1/sys/raw/":          {},
	"/v1/sys/generate-recovery-token/attempt": {},
	"/v1/sys/generate-recovery-token/update":  {},
}

// isUnauthenticatedCandidate reports whether a registered path must appear on
// the UnauthenticatedRoutes allowlist.
func isUnauthenticatedCandidate(path string) bool {
	if !strings.HasPrefix(path, "/v1/sys/") {
		return false
	}
	if _, ok := allowlistExemptPaths[path]; ok {
		return false
	}
	for _, p := range TopRouteOverrides {
		if path == p {
			return false
		}
	}
	// Individual pprof profiles are covered by the /v1/sys/pprof/ allowlist entry.
	if strings.HasPrefix(path, "/v1/sys/pprof/") && path != "/v1/sys/pprof/" {
		return false
	}
	return true
}
