package serve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asano69/archiver/internal/archive"
	"github.com/asano69/archiver/internal/static"
	"github.com/asano69/archiver/internal/version"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// registerRoutes wires up every HTTP route served by archiver. It is passed
// to app.OnServe().BindFunc in serve.go, keeping all route/handler
// definitions in this file while serve.go stays focused on server setup
// and startup.
func registerRoutes(e *core.ServeEvent) error {
	// Public routes: no auth required. Keep this list limited to
	// endpoints that return no user data (version info, health checks,
	// the static SPA shell below).
	e.Router.GET("/api/version", func(re *core.RequestEvent) error {
		return re.JSON(http.StatusOK, map[string]string{"version": version.Version})
	})

	e.Router.GET("/health", func(re *core.RequestEvent) error {
		return re.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Custom API routes that return or mutate user data go under this
	// group so RequireSuperuserAuth only has to be declared once here,
	// instead of on every individual route.
	admin := e.Router.Group("/api/admin")
	admin.Bind(apis.RequireSuperuserAuth())
	// e.g. admin.POST("/jobs/rescan", rescanHandler)

	// Kept outside the "/api/admin" group on purpose: it authenticates
	// against the "users" collection (see singleFileAuthRecord), not
	// _superusers, so RequireSuperuserAuth doesn't apply here.
	e.Router.POST("/api/singlefile", singleFileUploadHandler)

	// Serves the whole Vite build output (index.html, hashed JS/CSS
	// under assets/, and public/ files like favicon.svg copied to the
	// root) from a single route. indexFallback=true makes any unmatched
	// path (e.g. /manifests/abc, /settings) fall back to index.html, so
	// Solid Router can handle it client-side even on a hard refresh.
	// This shell is left unauthenticated on purpose: it's an empty
	// HTML/JS bundle with no data in it. Every route that actually
	// returns collection data is guarded below with
	// RequireSuperuserAuth, so an unauthenticated visitor only ever
	// sees the login screen the SPA renders client-side.
	e.Router.GET("/{path...}", apis.Static(static.FS, true))

	return e.Next()
}

// singleFileAuthRecord resolves the Authorization header sent by the
// SingleFile browser extension's "REST Form API" destination, which is
// always "Bearer <token>". This is checked by hand instead of
// apis.RequireAuth() because PocketBase itself expects the bare token
// without the "Bearer " prefix.
func singleFileAuthRecord(re *core.RequestEvent) (*core.Record, error) {
	token := strings.TrimPrefix(re.Request.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}
	return re.App.FindAuthRecordByToken(token, core.TokenTypeAuth)
}

// singleFileUploadHandler accepts a page snapshot uploaded by the
// SingleFile extension and saves it directly as a "done" archive,
// without running monolith. Field names ("file", "url") must match
// what's configured in the extension's "REST Form API" settings.
func singleFileUploadHandler(re *core.RequestEvent) error {
	if _, err := singleFileAuthRecord(re); err != nil {
		return apis.NewUnauthorizedError("invalid or expired token", err)
	}

	rawURL := re.Request.FormValue("url")
	if rawURL == "" {
		return apis.NewBadRequestError("missing url field", nil)
	}

	_, fileHeader, err := re.Request.FormFile("file")
	if err != nil {
		return apis.NewBadRequestError("missing file field", err)
	}

	if _, err := archive.FromUpload(re.App, rawURL, fileHeader); err != nil {
		return apis.NewBadRequestError("failed to save archive", err)
	}

	// Must be 200/201 with a JSON body, or SingleFile treats it as a
	// failed upload (see the extension's rest-form-api client).
	return re.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}
