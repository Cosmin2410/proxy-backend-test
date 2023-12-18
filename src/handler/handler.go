package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Cosmin2410/proxy-backend-test/src/helper"
	"github.com/Cosmin2410/proxy-backend-test/src/model"
	"github.com/andybalholm/brotli"
	"gorm.io/gorm"
)

type DBCreate struct {
	DB *gorm.DB
}

func (h *DBCreate) ModifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			return errors.New("content not json")
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code no OK, received status: %v", resp.StatusCode)
		}

		reader := brotli.NewReader(resp.Body)

		body, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("error reading the body: %v", err)
		}
		resp.Body.Close()

		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Transfer-Encoding")

		modifiedJSON, err := helper.AddRandomAttribute(body)
		if err != nil {
			return errors.New("failed to modify JSON response")
		}

		logEntry := model.SaveLog{
			Request:  resp.Request.URL.String(),
			Response: string(modifiedJSON),
		}

		result := h.DB.Create(&logEntry)
		if result.Error != nil {
			return errors.New("failed to log in db request/response")
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(modifiedJSON))

		return nil
	}
}
