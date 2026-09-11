package supabase

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	SupabaseURL    string
	ServiceRoleKey string
	httpClient     *http.Client
}

func NewClient(supabaseURL, serviceRoleKey string) *Client {
	return &Client{
		SupabaseURL:    strings.TrimRight(supabaseURL, "/"),
		ServiceRoleKey: serviceRoleKey,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}
}

// UploadFile uploads raw binary data/reader to a specified bucket and path in Supabase Storage.
// It returns the full public URL of the uploaded object upon success.
func (c *Client) UploadFile(bucketName, filePathInBucket string, body io.Reader, contentType string) (string, error) {
	if c.SupabaseURL == "" || c.ServiceRoleKey == "" {
		return "", fmt.Errorf("supabase client is not configured properly (missing SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY)")
	}

	cleanPath := strings.TrimPrefix(filePathInBucket, "/")
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.SupabaseURL, bucketName, cleanPath)

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("apiKey", c.ServiceRoleKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute upload request to supabase: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("supabase storage upload error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return c.GetPublicURL(bucketName, cleanPath), nil
}

// UploadMultipartFile is a convenience helper for uploading a multipart.FileHeader (e.g. from Fiber FormFile).
// It generates a unique filename to prevent collisions and returns the public URL.
func (c *Client) UploadMultipartFile(bucketName string, fileHeader *multipart.FileHeader, folder string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open multipart file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(fileHeader.Filename)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	filePath := uniqueName
	if folder != "" {
		filePath = fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), uniqueName)
	}

	return c.UploadFile(bucketName, filePath, src, contentType)
}

// DeleteFile removes an object from a Supabase Storage bucket.
func (c *Client) DeleteFile(bucketName, filePathInBucket string) error {
	if c.SupabaseURL == "" || c.ServiceRoleKey == "" {
		return fmt.Errorf("supabase client is not configured properly")
	}

	cleanPath := strings.TrimPrefix(filePathInBucket, "/")
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.SupabaseURL, bucketName, cleanPath)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("apiKey", c.ServiceRoleKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase storage delete error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetPublicURL formats and returns the public CDN URL for an object in a public bucket.
func (c *Client) GetPublicURL(bucketName, filePathInBucket string) string {
	cleanPath := strings.TrimPrefix(filePathInBucket, "/")
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.SupabaseURL, bucketName, cleanPath)
}
