package svcHealthcheck

import (
	"net/http"

	"github.com/FDSAP-Git-Org/hephaestus/apilogs"
	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func HealthCheck(c fiber.Ctx) error {
	response := v1.JSONResponse(c, respcode.SUC_CODE_200, "This service is running.", http.StatusOK)
	apilogs.ApplicationLogger(c.Path(), "System", "healthCheck", "Check Health", respcode.SUC_CODE_200_MSG, nil, response)
	return response
}
