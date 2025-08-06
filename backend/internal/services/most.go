package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type MostService struct {
	url string
}

func NewMostService(url string) *MostService {
	return &MostService{
		url: url,
	}
}

type Most interface {
	Send(ctx context.Context, dto *models.CreatePostDTO) error
	Update(ctx context.Context, dto *models.UpdatePostDTO) error
}

func (s *MostService) Send(ctx context.Context, dto *models.CreatePostDTO) error {
	apiPath := "/api/posts"

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(dto); err != nil {
		return fmt.Errorf("failed to encode notification data. error: %w", err)
	}

	resp, err := http.Post(s.url+apiPath, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("failed to send data to bot. error: %w", err)
	}

	if !strings.HasPrefix(resp.Status, "2") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("client: could not read response body: %w", err)
		}

		return fmt.Errorf("request returned an error: %s", string(body))
	}

	return nil
}

func (s *MostService) Update(ctx context.Context, dto *models.UpdatePostDTO) error {
	apiPath := fmt.Sprintf("/api/posts/%s", dto.PostId)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(dto); err != nil {
		return fmt.Errorf("failed to encode notification data. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, s.url+apiPath, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request. error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data to bot. error: %w", err)
	}
	return nil
}
