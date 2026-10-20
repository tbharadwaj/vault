// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	gziphandler "github.com/klauspost/compress/gzhttp"
)

// registerUIRoutes registers the Vault web UI routes including
// static asset serving, robots.txt, and the root redirect.
func registerUIRoutes(rc *RouteContext) {
	if !rc.Core.UIEnabled() {
		return
	}
	if uiBuiltIn {
		rc.Mux.Handle("/ui/", http.StripPrefix("/ui/",
			gziphandler.GzipHandler(
				handleUIHeaders(rc.Core,
					handleUI(http.FileServer(
						&UIAssetWrapper{FileSystem: assetFS()}))))))
		rc.Mux.Handle("/robots.txt",
			gziphandler.GzipHandler(
				handleUIHeaders(rc.Core,
					handleUI(http.FileServer(
						&UIAssetWrapper{FileSystem: assetFS()})))))
	} else {
		rc.Mux.Handle("/ui/", handleUIHeaders(rc.Core, handleUIStub()))
	}
	rc.Mux.Handle("/ui", handleUIRedirect())
	rc.Mux.Handle("/", handleUIRedirect())
}
