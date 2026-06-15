package apis

import (
	"embed"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi.yaml
var openapiFS embed.FS

const swaggerUITemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>W API — Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/api/openapi.yaml",
      dom_id: "#swagger-ui",
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset,
      ],
      layout: "BaseLayout",
    });
  </script>
</body>
</html>`

func OpenAPISpec(c *fiber.Ctx) error {
	data, err := openapiFS.ReadFile("openapi.yaml")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read openapi spec",
		})
	}
	c.Set("Content-Type", "text/yaml; charset=utf-8")
	return c.Send(data)
}

func SwaggerUI(c *fiber.Ctx) error {
	tmpl := template.Must(template.New("swagger").Parse(swaggerUITemplate))
	c.Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(c, nil)
}
