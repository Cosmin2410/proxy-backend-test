package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/Cosmin2410/proxy-backend-test/src/helper"
	"github.com/Cosmin2410/proxy-backend-test/src/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DBCreate struct {
	DB *gorm.DB
}

func (h *DBCreate) ReverseProxyHandler(c *fiber.Ctx) error {
	target := "https://jsonplaceholder.typicode.com" + c.OriginalURL()

	rpURL, err := url.Parse(target)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error in parsing target URL")
	}

	reverseProxy := httptest.NewServer(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(rpURL)
		},
	})
	defer reverseProxy.Close()

	resp, err := http.Get(reverseProxy.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching data from target")
	}
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadGateway).SendString("Status code from target not OK")
	}

	contentType := resp.Header.Get("Content-Type")

	requestURL := c.OriginalURL()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading response body")
	}

	if contentType != "application/json; charset=utf-8" {
		return c.Status(fiber.StatusUnsupportedMediaType).SendString("Response is not JSON")
	}

	modifiedJSON, err := helper.AddRandomAttribute(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to modify JSON response")
	}

	logEntry := model.SaveLog{
		Request:  requestURL,
		Response: string(modifiedJSON),
	}

	if err := h.DB.Create(&logEntry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating new entry in db")
	}

	_, err = c.Write(modifiedJSON)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error writing modified JSON response")
	}

	return nil
}
