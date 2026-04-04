package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"planning-trip-be/internal/response"
)

const (
	contextUserIDKey    = "auth_user_id"
	contextUserEmailKey = "auth_user_email"
)

func AuthRequired(secret string) gin.HandlerFunc {
	trimmedSecret := strings.TrimSpace(secret)

	return func(c *gin.Context) {
		if trimmedSecret == "" {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusInternalServerError,
				Message: "auth is not configured",
			})
			c.Abort()
			return
		}

		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			writeUnauthorized(c)
			return
		}

		userID, email, err := verifyAccessToken(token, trimmedSecret)
		if err != nil {
			writeUnauthorized(c)
			return
		}

		c.Set(contextUserIDKey, userID)
		c.Set(contextUserEmailKey, email)
		c.Next()
	}
}

func AuthUserID(c *gin.Context) (string, bool) {
	value, ok := c.Get(contextUserIDKey)
	if !ok {
		return "", false
	}
	userID, ok := value.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}
	return userID, true
}

func extractBearerToken(header string) string {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func verifyAccessToken(token string, secret string) (string, string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 4 {
		return "", "", http.ErrNoCookie
	}

	userID := strings.TrimSpace(parts[0])
	email := strings.TrimSpace(parts[1])
	expiresAtRaw := strings.TrimSpace(parts[2])
	signature := strings.TrimSpace(parts[3])

	if userID == "" || email == "" || expiresAtRaw == "" || signature == "" {
		return "", "", http.ErrNoCookie
	}

	payload := userID + "|" + email + "|" + expiresAtRaw
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(payload)); err != nil {
		return "", "", err
	}

	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return "", "", http.ErrNoCookie
	}

	expiresAt, err := strconv.ParseInt(expiresAtRaw, 10, 64)
	if err != nil || time.Now().UTC().Unix() > expiresAt {
		return "", "", http.ErrNoCookie
	}

	return userID, email, nil
}

func writeUnauthorized(c *gin.Context) {
	response.WriteError(c.Writer, response.APIError{
		Status:  http.StatusUnauthorized,
		Message: "unauthorized",
	})
	c.Abort()
}
