// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestRegisterSecretsAndIdentityRoutes_ExistAndCallable(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	props := &vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}
	rc := NewRouteContext(props, handlerSettings{})

	require.NotPanics(t, func() { registerSecretsRoutes(rc) })
	require.NotPanics(t, func() { registerIdentityRoutes(rc) })
	require.Empty(t, extractRegisteredRoutes(rc.Mux), "secrets/identity registrars are no-ops")
}

func TestHandlerWithSettings_CallsSecretsAndIdentityRegistrars(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	dir := filepath.Dir(thisFile)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(dir, "handler.go"), nil, 0)
	require.NoError(t, err)

	var callsSecrets, callsIdentity bool
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		switch ident.Name {
		case "registerSecretsRoutes":
			callsSecrets = true
		case "registerIdentityRoutes":
			callsIdentity = true
		}
		return true
	})
	require.True(t, callsSecrets, "handlerWithSettings must call registerSecretsRoutes")
	require.True(t, callsIdentity, "handlerWithSettings must call registerIdentityRoutes")
}
