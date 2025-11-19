package doc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"portfolio/api/http/middlewares"
	"portfolio/api/http/routes"
	"portfolio/docs"
	"portfolio/domain/entities"
	"portfolio/logger"
	"strconv"
	"text/template"
	"time"
)

var redocHTML = template.Must(template.New("redoc").Parse(`

<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Portfolio API Documentation</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="/doc/theme.css">
	</head>
	<body>
		<div class="theme-toggle-wrap">
      <button id="toggle-theme" type="button" onclick="toggleTheme()"></button>
    </div>
		<div id="redoc-container"></div>
		<script type="module">
			(function(){
				var KEY = 'redoc-theme';
				var m = window.matchMedia('(prefers-color-scheme: dark)');
				function sys() { return m.matches ? 'dark' : 'light'; }
				var stored = localStorage.getItem(KEY) || 'dark';
				var effective = stored === 'dark' ? sys() : stored;

				function apply(mode, src) {
					document.documentElement.setAttribute('theme', mode);
					document.documentElement.setAttribute('data-theme-mode', src);
					document.documentElement.style.colorScheme = (mode === 'dark') ? 'dark' : 'light';
				}
				apply(effective, stored);

				m.addEventListener('change', function() {
					var mode = localStorage.getItem(KEY) || 'dark';
					if (mode === 'dark') {
						applyTheme('dark', true);
					}
				});

				function btnLabel(next){
					return next==='dark' ? '☀️' : '🌙';
				}

				window.applyTheme = function(next, reRender) {
					localStorage.setItem(KEY, next);
					var eff = next==='dark' ? sys(): next;
					apply(eff, next);
					var btn=document.getElementById('toggle-theme');
					if (btn){
						btn.textContent = btnLabel(next);
						btn.title = 'Basculer le thème ('+next+')';
					}
					if (reRender && window.__renderRedoc) {
						setTimeout(window.__renderRedoc, 75);
					}
				};
				window.toggleTheme = function() {
					var cur = localStorage.getItem(KEY) || 'dark';
					var next = (cur==='dark') ? 'light' : 'dark';
					window.applyTheme(next, true);
					setTimeout(function() {
						const url = new URL(window.location.href);
						url.searchParams.set('_ts', Date.now().toString());
						window.location.replace(url);
					}, 75);

				};
				window.__initTheme = function() { window.applyTheme(stored, false); };
			})();
		</script>
		<script type="module">
			import theme from '/doc/theme.js';

			var SPEC_URLS = ['/doc/openapi.json', '/doc/swagger.json', '/openapi.yaml'];
      function findSpecUrl() {
        return new Promise(function(resolve) {
          var i=0;
          function next() {
            if (i>=SPEC_URLS.length){ resolve(SPEC_URLS[0]); return; }
            var u=SPEC_URLS[i++], xhr=new XMLHttpRequest();
            xhr.open('HEAD', u, true);
            xhr.onreadystatechange = function() {
              if (xhr.readyState === 4) {
                if (xhr.status >= 200 && xhr.status < 400) { resolve(u); }
                else { next(); }
              }
            };
            try { xhr.send(null); } catch (e) { next(); }
          }
          next();
        });
      }

			async function render(){
        var url = await findSpecUrl();
        var el = document.getElementById('redoc-container');
        el.innerHTML = '';

        Redoc.init(
          url,
          {
            theme: theme,
            pathInMiddlePanel: true,
            hideDownloadButton: true,
            expandResponses: '200,201',
            jsonSampleExpandLevel: 'all',
            onlyRequiredInSamples: true,
						disableSearch: true
          },
          el
        );
      }
			window.__renderRedoc = render;
			if (window.__initTheme) window.__initTheme();
		</script>
		<script src="https://cdn.redoc.ly/redoc/v2.5.1/bundles/redoc.standalone.js" defer onload="window.__renderRedoc && window.__renderRedoc()"></script>
	</body>
</html>
`))

type documentationHandler struct {
	logger *logger.Logger
}

func NewDocsHandler(logger *logger.Logger) []*routes.NamedRoute {
	h := documentationHandler{
		logger: logger,
	}
	return []*routes.NamedRoute{
		{Name: "GetRedocUIHandler", Pattern: "GET /", Handler: h.RedocUIHandler},
		{Name: "GetOpenAPISpecHandler", Pattern: "GET /openapi.json", Handler: h.OpenAPISpecHandler},
		{Name: "GetSwaggerSpecHandler", Pattern: "GET /swagger.json", Handler: h.OpenAPISpecHandler},
		{Name: "GetSwaggerThemeCSSHandler", Pattern: "GET /theme.css", Handler: h.SwaggerThemeCSSHandler},
		{Name: "GetSwaggerThemeJSHandler", Pattern: "GET /theme.js", Handler: h.SwaggerThemeJSHandler},
	}
}

func (dh *documentationHandler) RedocUIHandler(w http.ResponseWriter, r *http.Request) {
	// sum := sha256.Sum256(docs.SwaggerJSON)
	// _ = sum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = redocHTML.Execute(w, nil)
}

func serveBytes(w http.ResponseWriter, r *http.Request, b []byte, contentType string, maxAge time.Duration) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	s := sha256.Sum256(b)
	etag := `W/"` + hex.EncodeToString(s[:]) + `"`
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(maxAge.Seconds())))

	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(b))
}

func (dh *documentationHandler) OpenAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middlewares.GetUserFromContext(r)
	if !ok {
		dh.logger.Error("Failed to get user ID from context")
		return
	}
	if user == nil || user.ID <= 0 {
		serveBytes(w, r, docs.SwaggerJSON, "application/json; charset=utf-8", 24*time.Hour)
		return
	}
	if user.Role == entities.RoleAdmin {
		serveBytes(w, r, docs.PrivateSwaggerJSON, "application/json; charset=utf-8", 24*time.Hour)
		return
	}
	serveBytes(w, r, docs.SwaggerJSON, "application/json; charset=utf-8", 24*time.Hour)
}

func (dh *documentationHandler) SwaggerThemeJSHandler(w http.ResponseWriter, r *http.Request) {
	serveBytes(w, r, docs.SwaggerThemeJS, "application/javascript; charset=utf-8", 24*time.Hour)
}

func (dh *documentationHandler) SwaggerThemeCSSHandler(w http.ResponseWriter, r *http.Request) {
	serveBytes(w, r, docs.SwaggerThemeCSS, "text/css; charset=utf-8", 24*time.Hour)
}
