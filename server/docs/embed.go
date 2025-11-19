package docs

import _ "embed"

//go:embed theme.css
var SwaggerThemeCSS []byte

//go:embed theme.js
var SwaggerThemeJS []byte

//go:embed swagger.json
var SwaggerJSON []byte

//go:embed private_swagger.json
var PrivateSwaggerJSON []byte
