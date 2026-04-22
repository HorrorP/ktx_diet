package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func TelegramInitDataAuth(botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.GetHeader("X-Telegram-Init-Data")
		if initData == "" {
			initData = c.Query("init_data")
		}
		if initData == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing telegram initData"})
			return
		}

		if !verifyTelegramInitData(initData, botToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram initData"})
			return
		}

		c.Next()
	}
}

func verifyTelegramInitData(initData, botToken string) bool {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return false
	}

	var keys []string
	for k := range values {
		if k == "hash" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, k+"="+values.Get(k))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	secretKeyHmac := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHmac.Write([]byte(botToken))
	secretKey := secretKeyHmac.Sum(nil)

	hashHmac := hmac.New(sha256.New, secretKey)
	hashHmac.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(hashHmac.Sum(nil))

	return hmac.Equal([]byte(strings.ToLower(receivedHash)), []byte(strings.ToLower(computedHash)))
}
