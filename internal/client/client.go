// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"time"
)

// ErrorResponse represents an error response from the server
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// parseErrorResponse parses error response body and returns a formatted error
func parseErrorResponse(resp *http.Response) error {
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("HTTP %d: %s (failed to parse error response)", resp.StatusCode, resp.Status)
	}
	
	if errResp.Message != "" {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Message)
	}
	
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
}

// Client represents the DirtCloud API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new DirtCloud API client.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8080/v1"
	}

	token := os.Getenv("DIRT_TOKEN")

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}
}

// Project represents a DirtCloud project.
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProjectRequest represents the request body for creating a project.
type CreateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateProjectRequest represents the request body for updating a project.
type UpdateProjectRequest struct {
	Name string `json:"name"`
}

// Instance represents a DirtCloud instance.
type Instance struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	CPU       int       `json:"cpu"`
	MemoryMB  int       `json:"memory_mb"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInstanceRequest represents the request body for creating an instance.
type CreateInstanceRequest struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	CPU       int    `json:"cpu"`
	MemoryMB  int    `json:"memory_mb"`
	Image     string `json:"image"`
	Status    string `json:"status,omitempty"`
}

// UpdateInstanceRequest represents the request body for updating an instance.
type UpdateInstanceRequest struct {
	Name     *string `json:"name,omitempty"`
	CPU      *int    `json:"cpu,omitempty"`
	MemoryMB *int    `json:"memory_mb,omitempty"`
	Image    *string `json:"image,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// Metadata represents DirtCloud metadata.
type Metadata struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// doRequest performs an HTTP request with proper authentication.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	fullURL := c.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}

// Projects API

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/projects", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &project, nil
}

// GetProject retrieves a project by ID.
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	resp, err := c.doRequest(ctx, "GET", "/projects/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &project, nil
}

// ListProjects retrieves all projects, optionally filtered by name.
func (c *Client) ListProjects(ctx context.Context, nameFilter string) ([]Project, error) {
	endpoint := "/projects"
	if nameFilter != "" {
		endpoint += "?name=" + url.QueryEscape(nameFilter)
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return projects, nil
}

// UpdateProject updates a project.
func (c *Client) UpdateProject(ctx context.Context, id string, req UpdateProjectRequest) (*Project, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "PATCH", "/projects/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &project, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/projects/"+id, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("project not found")
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Instances API

// CreateInstance creates a new instance.
func (c *Client) CreateInstance(ctx context.Context, req CreateInstanceRequest) (*Instance, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/instances", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &instance, nil
}

// GetInstance retrieves an instance by ID.
func (c *Client) GetInstance(ctx context.Context, id string) (*Instance, error) {
	resp, err := c.doRequest(ctx, "GET", "/instances/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("instance not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &instance, nil
}

// ListInstances retrieves all instances, optionally filtered.
func (c *Client) ListInstances(ctx context.Context, projectID, nameFilter, statusFilter string) ([]Instance, error) {
	endpoint := "/instances"
	params := url.Values{}

	if projectID != "" {
		params.Add("project_id", projectID)
	}
	if nameFilter != "" {
		params.Add("name", nameFilter)
	}
	if statusFilter != "" {
		params.Add("status", statusFilter)
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instances []Instance
	if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return instances, nil
}

// UpdateInstance updates an instance.
func (c *Client) UpdateInstance(ctx context.Context, id string, req UpdateInstanceRequest) (*Instance, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "PATCH", "/instances/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("instance not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &instance, nil
}

// DeleteInstance deletes an instance.
func (c *Client) DeleteInstance(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/instances/"+id, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("instance not found")
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Metadata API

// CreateMetadataRequest represents the request body for creating metadata.
type CreateMetadataRequest struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// UpdateMetadataRequest represents the request body for updating metadata.
type UpdateMetadataRequest struct {
	Path  *string `json:"path,omitempty"`
	Value *string `json:"value,omitempty"`
}

// CreateMetadata creates new metadata.
func (c *Client) CreateMetadata(ctx context.Context, req CreateMetadataRequest) (*Metadata, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/metadata", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &metadata, nil
}

// GetMetadata retrieves metadata by ID.
func (c *Client) GetMetadata(ctx context.Context, id string) (*Metadata, error) {
	resp, err := c.doRequest(ctx, "GET", "/metadata/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("metadata not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &metadata, nil
}

// GetMetadataByPath retrieves metadata by path.
func (c *Client) GetMetadataByPath(ctx context.Context, path string) (*Metadata, error) {
	// Use ListMetadata to find the exact path match
	metadataList, err := c.ListMetadata(ctx, path)
	if err != nil {
		return nil, err
	}

	// Find exact match
	for _, metadata := range metadataList {
		if metadata.Path == path {
			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("metadata not found")
}

// ListMetadata lists all metadata, optionally filtered by prefix.
func (c *Client) ListMetadata(ctx context.Context, prefix string) ([]Metadata, error) {
	endpoint := "/metadata"
	if prefix != "" {
		endpoint += "?prefix=" + url.QueryEscape(prefix)
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata []Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return metadata, nil
}

// UpdateMetadata updates metadata by ID.
func (c *Client) UpdateMetadata(ctx context.Context, id string, req UpdateMetadataRequest) (*Metadata, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.doRequest(ctx, "PATCH", "/metadata/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("metadata not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &metadata, nil
}

// DeleteMetadata deletes metadata by ID.
func (c *Client) DeleteMetadata(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/metadata/"+id, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("metadata not found")
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
