package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func reverseProxyHandler(c *fiber.Ctx) error {
	target := "https://jsonplaceholder.typicode.com" + c.OriginalURL()

	rpURL, err := url.Parse(target)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error in parsing target URL")
	}

	frontendProxy := httptest.NewServer(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(rpURL)
		},
	})
	defer frontendProxy.Close()

	resp, err := http.Get(frontendProxy.URL)
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

	modifiedJSON, err := addRandomAttribute(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to modify JSON response")
	}

	logEntry := ProxyLog{
		Request:  requestURL,
		Response: string(modifiedJSON),
	}
	db.Create(&logEntry)

	_, err = c.Write(modifiedJSON)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error writing modified JSON response")
	}

	return nil
}
