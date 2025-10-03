package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SupabaseAuthError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	IsRetryable bool
	StatusCode  int
}

func (e *SupabaseAuthError) Error() string {
	return fmt.Sprintf("Supabase Auth Error [%s]: %s", e.Code, e.Message)
}

type SupabaseUser struct {
	ID                 string                 `json:"id"`
	Email              string                 `json:"email"`
	EmailConfirmedAt   *time.Time             `json:"email_confirmed_at"`
	Phone              *string                `json:"phone"`
	PhoneConfirmedAt   *time.Time             `json:"phone_confirmed_at"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	LastSignInAt       *time.Time             `json:"last_sign_in_at"`
	Role               string                 `json:"role"`
	UserMetadata       map[string]interface{} `json:"user_metadata"`
	AppMetadata        map[string]interface{} `json:"app_metadata"`
	IdentityData       []interface{}          `json:"identities"`
	Aud                string                 `json:"aud"`
	ConfirmationSentAt *time.Time             `json:"confirmation_sent_at"`
}

type SupabaseAuthResponse struct {
	User         SupabaseUser `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	ExpiresAt    int64        `json:"expires_at"`
	TokenType    string       `json:"token_type"`
}


type SupabaseRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SupabaseTokenValidationResponse struct {
	User SupabaseUser `json:"user"`
	Exp  int64        `json:"exp"`
	Iat  int64        `json:"iat"`
	Sub  string       `json:"sub"`
}

func createSupabaseAuthError(statusCode int, body []byte) *SupabaseAuthError {
	var errorResp map[string]interface{}

	if err := json.Unmarshal(body, &errorResp); err == nil {
		if msg, ok := errorResp["error_description"].(string); ok {
			return &SupabaseAuthError{
				Code:        "auth_error",
				Message:     msg,
				IsRetryable: statusCode >= 500, // Server errors are retryable
				StatusCode:  statusCode,
			}
		}
		if msg, ok := errorResp["message"].(string); ok {
			return &SupabaseAuthError{
				Code:        "auth_error",
				Message:     msg,
				IsRetryable: statusCode >= 500,
				StatusCode:  statusCode,
			}
		}
	}

	return &SupabaseAuthError{
		Code:        fmt.Sprintf("http_%d", statusCode),
		Message:     fmt.Sprintf("HTTP %d: %s", statusCode, string(body)),
		IsRetryable: statusCode >= 500,
		StatusCode:  statusCode,
	}
}

func makeSupabaseRequest(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func ValidateSupabaseTokenHTTP(token, supabaseURL, anonKey string) (*SupabaseUser, *SupabaseAuthError) {
	if token == "" || supabaseURL == "" || anonKey == "" {
		return nil, &SupabaseAuthError{
			Code:        "invalid_config",
			Message:     "Missing token, Supabase URL, or anon key",
			IsRetryable: false,
			StatusCode:  400,
		}
	}

	url := fmt.Sprintf("%s/auth/v1/user", supabaseURL)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"apikey":        anonKey,
	}

	resp, err := makeSupabaseRequest("GET", url, headers, nil)
	if err != nil {
		return nil, &SupabaseAuthError{
			Code:        "network_error",
			Message:     fmt.Sprintf("Failed to validate token: %v", err),
			IsRetryable: true,
			StatusCode:  0,
		}
	}
	defer resp.Body.Close()

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, createSupabaseAuthError(resp.StatusCode, responseBody.Bytes())
	}

	var user SupabaseUser
	if err := json.Unmarshal(responseBody.Bytes(), &user); err != nil {
		return nil, &SupabaseAuthError{
			Code:        "parse_error",
			Message:     fmt.Sprintf("Failed to parse user response: %v", err),
			IsRetryable: false,
			StatusCode:  resp.StatusCode,
		}
	}

	return &user, nil
}


func RefreshSupabaseTokenHTTP(refreshToken, supabaseURL, anonKey string) (*SupabaseAuthResponse, *SupabaseAuthError) {
	if refreshToken == "" || supabaseURL == "" || anonKey == "" {
		return nil, &SupabaseAuthError{
			Code:        "invalid_config",
			Message:     "Missing required parameters for token refresh",
			IsRetryable: false,
			StatusCode:  400,
		}
	}

	url := fmt.Sprintf("%s/auth/v1/token?grant_type=refresh_token", supabaseURL)

	requestData := SupabaseRefreshRequest{
		RefreshToken: refreshToken,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, &SupabaseAuthError{
			Code:        "marshal_error",
			Message:     fmt.Sprintf("Failed to marshal request: %v", err),
			IsRetryable: false,
			StatusCode:  0,
		}
	}

	headers := map[string]string{
		"apikey": anonKey,
	}

	resp, err := makeSupabaseRequest("POST", url, headers, requestBody)
	if err != nil {
		return nil, &SupabaseAuthError{
			Code:        "network_error",
			Message:     fmt.Sprintf("Failed to refresh token: %v", err),
			IsRetryable: true,
			StatusCode:  0,
		}
	}
	defer resp.Body.Close()

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, createSupabaseAuthError(resp.StatusCode, responseBody.Bytes())
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(responseBody.Bytes(), &authResp); err != nil {
		return nil, &SupabaseAuthError{
			Code:        "parse_error",
			Message:     fmt.Sprintf("Failed to parse auth response: %v", err),
			IsRetryable: false,
			StatusCode:  resp.StatusCode,
		}
	}

	return &authResp, nil
}

func GetSupabaseUserByEmailHTTP(email, supabaseURL, serviceKey string) (*SupabaseUser, *SupabaseAuthError) {
	if email == "" || supabaseURL == "" || serviceKey == "" {
		return nil, &SupabaseAuthError{
			Code:        "invalid_config",
			Message:     "Missing required parameters for user lookup",
			IsRetryable: false,
			StatusCode:  400,
		}
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users?email=%s", supabaseURL, email)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", serviceKey),
		"apikey":        serviceKey,
	}

	resp, err := makeSupabaseRequest("GET", url, headers, nil)
	if err != nil {
		return nil, &SupabaseAuthError{
			Code:        "network_error",
			Message:     fmt.Sprintf("Failed to lookup user: %v", err),
			IsRetryable: true,
			StatusCode:  0,
		}
	}
	defer resp.Body.Close()

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, createSupabaseAuthError(resp.StatusCode, responseBody.Bytes())
	}

	var users struct {
		Users []SupabaseUser `json:"users"`
	}

	if err := json.Unmarshal(responseBody.Bytes(), &users); err != nil {
		return nil, &SupabaseAuthError{
			Code:        "parse_error",
			Message:     fmt.Sprintf("Failed to parse users response: %v", err),
			IsRetryable: false,
			StatusCode:  resp.StatusCode,
		}
	}

	if len(users.Users) == 0 {
		return nil, &SupabaseAuthError{
			Code:        "user_not_found",
			Message:     "User not found in Supabase",
			IsRetryable: false,
			StatusCode:  404,
		}
	}

	return &users.Users[0], nil
}

func ConvertSupabaseIDToUUID(supabaseID string) (uuid.UUID, error) {
	return uuid.Parse(supabaseID)
}

func IsSupabaseErrorRetryable(err *SupabaseAuthError) bool {
	if err == nil {
		return false
	}
	return err.IsRetryable
}