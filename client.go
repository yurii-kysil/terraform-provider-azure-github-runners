package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Client struct {
	httpClient   *http.Client
	token        string
	baseURL      string
	organization string
	appAuth      *AppAuth
}

func NewClient(token, baseURL, organization string, insecure bool) (*Client, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	if insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		httpClient:   httpClient,
		token:        token,
		baseURL:      strings.TrimSuffix(baseURL, "/"),
		organization: organization,
		appAuth:      nil,
	}, nil
}

// NewClientWithAppAuth creates a new client with GitHub App authentication
func NewClientWithAppAuth(appAuth *AppAuth, baseURL, organization string, insecure bool) (*Client, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	if insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	var appPemFile string

	if appAuth.PEMFile != "" {
		// The Go encoding/pem package only decodes PEM formatted blocks
		// that contain new lines. Some platforms, like Terraform Cloud,
		// do not support new lines within Environment Variables.
		// Any occurrence of \n in the `pem_file` argument's value
		// (explicit value, or default value taken from
		// GITHUB_APP_PEM_FILE Environment Variable) is replaced with an
		// actual new line character before decoding.
		appPemFile = strings.Replace(appAuth.PEMFile, `\n`, "\n", -1)
	} else {
		return nil, fmt.Errorf("app_auth.pem_file must be set and contain a non-empty value")
	}

	// Generate JWT token
	jwtToken, err := generateJWT(appAuth.ID, appPemFile)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %v", err)
	}

	// Get installation access token
	installationToken, err := getInstallationTokenFromGitHub(jwtToken, appAuth.InstallationID, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %v", err)
	}

	return &Client{
		httpClient:   httpClient,
		token:        installationToken,
		baseURL:      strings.TrimSuffix(baseURL, "/"),
		organization: organization,
		appAuth:      appAuth,
	}, nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "terraform-provider-azure-github-runners")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s failed with status %d: %s", path, resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST %s failed with status %d: %s", path, resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PUT %s failed with status %d: %s", path, resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Patch(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, "PATCH", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PATCH %s failed with status %d: %s", path, resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, path string, body interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "terraform-provider-azure-github-runners")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Treat 404 as success for DELETE operations (resource already deleted)
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DELETE %s failed with status %d: %s", path, resp.StatusCode, string(body))
	}

	return nil
}

// InstallationTokenResponse represents the response from GitHub's installation token API
type InstallationTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// getInstallationTokenFromGitHub retrieves an installation access token from GitHub
func getInstallationTokenFromGitHub(jwtToken string, installationID int, baseURL string) (string, error) {
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", strings.TrimSuffix(baseURL, "/"), installationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "terraform-provider-azure-github-runners")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get installation token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp InstallationTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return tokenResp.Token, nil
}

// generateJWT creates a JWT token for GitHub App authentication
func generateJWT(appID int, privateKeyPEM string) (string, error) {
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Create the JWT claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 10).Unix(), // JWT expires in 10 minutes
		"iss": appID,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}
