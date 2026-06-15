package apis

import (
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

func OpenAPISpecJSON(c *fiber.Ctx) error {
	data, err := openapiFS.ReadFile("openapi.yaml")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read openapi spec",
		})
	}

	var spec any
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse openapi spec",
		})
	}

	return c.JSON(spec)
}
