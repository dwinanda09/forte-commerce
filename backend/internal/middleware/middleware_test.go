package middleware

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test context
func newMiddlewareContext(method, path string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// Helper function to create a valid JWT token
func makeTestToken(secret, userID, username, role string) string {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

// Dummy handler for testing
func dummyHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}

// TestRequestID tests the RequestID middleware
func TestRequestID(t *testing.T) {
	tests := []struct {
		name          string
		headerValue   string
		expectHeader  bool
		validateID    func(string) bool
		description   string
	}{
		{
			name:         "Header provided - passes through",
			headerValue:  "test-request-id-123",
			expectHeader: true,
			validateID: func(id string) bool {
				return id == "test-request-id-123"
			},
			description: "X-Request-ID header provided should be passed through",
		},
		{
			name:         "Header not provided - generates UUID",
			headerValue:  "",
			expectHeader: true,
			validateID: func(id string) bool {
				// UUID validation: should be 36 chars (8-4-4-4-12)
				return len(id) == 36
			},
			description: "Missing X-Request-ID header should generate UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newMiddlewareContext("GET", "/test")
			if tt.headerValue != "" {
				c.Request().Header.Set("X-Request-ID", tt.headerValue)
			}

			middleware := RequestID(dummyHandler)
			err := middleware(c)

			require.NoError(t, err)
			assert.Equal(t, 200, rec.Code)

			// Check context value
			requestID := c.Get("request_id")
			require.NotNil(t, requestID)
			requestIDStr, ok := requestID.(string)
			require.True(t, ok)
			assert.True(t, tt.validateID(requestIDStr), "%s", tt.description)

			// Check response header
			if tt.expectHeader {
				assert.Equal(t, requestIDStr, rec.Header().Get("X-Request-ID"))
			}
		})
	}
}

// TestJWTAuth tests the JWTAuth middleware
func TestJWTAuth(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name            string
		authHeader      string
		expectedStatus  int
		expectedError   string
		expectedUserID  string
		expectedRole    string
		expectedUsername string
		description     string
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			expectedStatus: 401,
			expectedError:  "Missing authorization header",
			description:    "Missing Authorization header should return 401",
		},
		{
			name:           "Invalid format - not Bearer",
			authHeader:     "Basic abc123",
			expectedStatus: 401,
			expectedError:  "Invalid authorization format",
			description:    "Authorization header not starting with Bearer should return 401",
		},
		{
			name:           "Invalid format - Bearer too short",
			authHeader:     "Bearer ",
			expectedStatus: 401,
			expectedError:  "Invalid token",
			description:    "Invalid JWT token string should return 401",
		},
		{
			name:           "Invalid token - malformed",
			authHeader:     "Bearer invalid.token.string",
			expectedStatus: 401,
			expectedError:  "Invalid token",
			description:    "Malformed JWT token should return 401",
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer " + makeTestToken(secret, "user-123", "testuser", "seller"),
			expectedStatus: 200,
			expectedUserID: "user-123",
			expectedUsername: "testuser",
			expectedRole:   "seller",
			description:    "Valid JWT token should call next handler and set claims in context",
		},
		{
			name:           "Valid token with buyer role",
			authHeader:     "Bearer " + makeTestToken(secret, "user-456", "buyeruser", "buyer"),
			expectedStatus: 200,
			expectedUserID: "user-456",
			expectedUsername: "buyeruser",
			expectedRole:   "buyer",
			description:    "Valid JWT token with buyer role should set correct claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newMiddlewareContext("GET", "/test")
			if tt.authHeader != "" {
				c.Request().Header.Set("Authorization", tt.authHeader)
			}

			middleware := JWTAuth(secret)(dummyHandler)
			err := middleware(c)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code, tt.description)

			if tt.expectedStatus != 200 {
				// Check error response
				body := rec.Body.String()
				assert.Contains(t, body, tt.expectedError)
			} else {
				// Check context values for successful auth
				userID, ok := c.Get("user_id").(string)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedUserID, userID)

				username, ok := c.Get("username").(string)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedUsername, username)

				role, ok := c.Get("role").(string)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedRole, role)
			}
		})
	}
}

// TestSellerOnly tests the SellerOnly middleware
func TestSellerOnly(t *testing.T) {
	tests := []struct {
		name           string
		roleValue      interface{}
		expectedStatus int
		expectedError  string
		description    string
	}{
		{
			name:           "Role is seller",
			roleValue:      "seller",
			expectedStatus: 200,
			description:    "Context with role=seller should call next handler",
		},
		{
			name:           "Role is buyer",
			roleValue:      "buyer",
			expectedStatus: 403,
			expectedError:  "Seller access required",
			description:    "Context with role=buyer should return 403",
		},
		{
			name:           "Role is admin",
			roleValue:      "admin",
			expectedStatus: 403,
			expectedError:  "Seller access required",
			description:    "Context with role=admin should return 403",
		},
		{
			name:           "No role in context",
			roleValue:      nil,
			expectedStatus: 403,
			expectedError:  "Seller access required",
			description:    "Context without role should return 403",
		},
		{
			name:           "Role is wrong type",
			roleValue:      123,
			expectedStatus: 403,
			expectedError:  "Seller access required",
			description:    "Context with non-string role should return 403",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newMiddlewareContext("GET", "/test")
			if tt.roleValue != nil {
				c.Set("role", tt.roleValue)
			}

			middleware := SellerOnly(dummyHandler)
			err := middleware(c)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code, tt.description)

			if tt.expectedStatus != 200 {
				body := rec.Body.String()
				assert.Contains(t, body, tt.expectedError)
			}
		})
	}
}

// TestJWTAuthWithSellerOnly tests JWT auth combined with SellerOnly
func TestJWTAuthWithSellerOnly(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		description    string
	}{
		{
			name:           "Valid seller token",
			token:          makeTestToken(secret, "seller-123", "seller-user", "seller"),
			expectedStatus: 200,
			description:    "Valid seller token should pass through both middlewares",
		},
		{
			name:           "Valid buyer token",
			token:          makeTestToken(secret, "buyer-123", "buyer-user", "buyer"),
			expectedStatus: 403,
			description:    "Valid buyer token should fail at SellerOnly middleware",
		},
		{
			name:           "Invalid token",
			token:          "invalid",
			expectedStatus: 401,
			description:    "Invalid token should fail at JWT middleware",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newMiddlewareContext("GET", "/test")
			c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			// Chain middlewares
			handler := JWTAuth(secret)(SellerOnly(dummyHandler))
			err := handler(c)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code, tt.description)
		})
	}
}

// TestRequestIDAndJWTAuthChaining tests multiple middlewares together
func TestRequestIDAndJWTAuthChaining(t *testing.T) {
	secret := "test-secret-key"
	token := makeTestToken(secret, "user-123", "testuser", "seller")

	c, rec := newMiddlewareContext("GET", "/test")
	c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("X-Request-ID", "custom-request-id")

	// Chain middlewares: RequestID -> JWTAuth
	handler := RequestID(JWTAuth(secret)(dummyHandler))
	err := handler(c)

	require.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	// Verify RequestID was set
	requestID := c.Get("request_id")
	require.NotNil(t, requestID)
	assert.Equal(t, "custom-request-id", requestID.(string))

	// Verify JWT claims were set
	userID, ok := c.Get("user_id").(string)
	assert.True(t, ok)
	assert.Equal(t, "user-123", userID)

	role, ok := c.Get("role").(string)
	assert.True(t, ok)
	assert.Equal(t, "seller", role)

	// Verify response header
	assert.Equal(t, "custom-request-id", rec.Header().Get("X-Request-ID"))
}
