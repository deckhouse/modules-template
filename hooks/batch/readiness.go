package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/deckhouse/module-sdk/pkg"
)

func ReadinessFunc(_ context.Context, input *pkg.HookInput) error {
	input.Logger.Info("start user logic for readiness probe")

	c := input.DC.GetHTTPClient()

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1/readyz", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	input.Logger.Debug("readiness probe done successfully", slog.Any("body", string(respBody)))

	return nil
}
