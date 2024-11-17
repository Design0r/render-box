package api

import (
	"fmt"
	"text/template"

	"github.com/labstack/echo/v4"

	"render-box/client/assets"
)

func HandleIndex(c echo.Context) error {
	tmpl, err := template.ParseFS(assets.TemplatesFS, "templates/index.html")
	if err != nil {
		return fmt.Errorf("Coult not load template: %v", err)
	}

	err = tmpl.Execute(c.Response().Writer, nil)
	if err != nil {
		return fmt.Errorf("Could not render template: %v", err)
	}

	return nil
}
