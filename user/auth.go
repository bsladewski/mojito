package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bsladewski/mojito/data"
	"github.com/bsladewski/mojito/httperror"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// jwtAccessMetadata stores information embedded in a JWT access token.
type jwtAccessMetadata struct {
	authUUID  string
	userID    uint
	createdAt time.Time
	expiresAt time.Time
}

// jwtRefreshMetadata stores information embedded in a JWT refresh token.
type jwtRefreshMetadata struct {
	authUUID  string
	userID    uint
	createdAt time.Time
	expiresAt time.Time
}

const (
	// authFailedGeneric is a generic error returned when a request cannot be
	// authenticated.
	authFailedGeneric = "request not authorized"
)

// GenerateSecretToken creates a base64 encoded token that includes both the
// supplied user id as well as the supplied payload encrypted with the user
// secret key.
func GenerateSecretToken(ctx context.Context, u *User,
	payload string) (string, error) {

	// create cipher with user secret key
	cipherBlock, err := aes.NewCipher([]byte(u.SecretKey))
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// encrypt and base64 encode payload
	payload = base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce,
		[]byte(payload), nil))

	// marshal token contents to json
	contents, err := json.Marshal(struct {
		UserID  uint
		Payload string
	}{
		UserID:  u.ID,
		Payload: payload,
	})
	if err != nil {
		return "", err
	}

	// base64 encode json token contents
	return base64.StdEncoding.EncodeToString(contents), nil

}

// ParseSecretToken parses the supplied secret token and returns the user id
// associated with the token as well as the decrypted payload string.
func ParseSecretToken(ctx context.Context,
	token string) (u *User, payload string, err error) {

	// base64 decode token contents
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, "", err
	}

	// unmarshal json token contents
	var tokenData = struct {
		UserID  uint
		Payload string
	}{}

	if err = json.Unmarshal(tokenBytes, &tokenData); err != nil {
		return nil, "", err
	}

	// get user record
	u, err = GetUserByID(ctx, data.DB(), tokenData.UserID)
	if err != nil {
		return nil, "", err
	}

	// base64 decode encrypted payload
	encryptData, err := base64.URLEncoding.DecodeString(tokenData.Payload)
	if err != nil {
		return nil, "", err
	}

	// create cipher with user secret key
	cipherBlock, err := aes.NewCipher([]byte(u.SecretKey))
	if err != nil {
		return nil, "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, "", err
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return nil, "", err
	}

	// decrypt the payload
	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	payloadBytes, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, "", err
	}

	// return string representation of payload
	return u, string(payloadBytes), nil

}

// JWTAuthMiddleware gets middleware that handles request authentication using
// a JWT bearer token.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtAccessTokenValid(c); err != nil {
			logrus.Warn(err)
			c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
				ErrorMessage: authFailedGeneric,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTGetUser extracts a user record from the request access token.
func JWTGetUser(c *gin.Context) (*User, error) {

	metadata, err := jwtGetAccessMetadata(c)
	if err != nil {
		return nil, err
	}

	return GetUserByID(c, data.DB(), metadata.userID)

}

// JWTGetUserLogin extracts a user login record from the request access token.
func JWTGetUserLogin(c *gin.Context) (*Login, error) {

	metadata, err := jwtGetAccessMetadata(c)
	if err != nil {
		return nil, err
	}

	return GetLoginByUUID(c, data.DB(), metadata.authUUID)

}

// JWTValidateRefreshToken checks whether the supplied refresh token is valid,
// returns the associated user login record if the token is valid.
func JWTValidateRefreshToken(c *gin.Context,
	refreshToken string) (*Login, error) {

	metadata, err := jwtGetRefreshMetadata(c, refreshToken)
	if err != nil {
		return nil, err
	}

	login, err := GetLoginByUUID(c, data.DB(), metadata.authUUID)
	if err != nil {
		return nil, err
	}

	if metadata.expiresAt.Before(time.Now()) ||
		login.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return login, nil

}

// jwtAccessTokenValid checks whether the request access token is valid.
func jwtAccessTokenValid(c *gin.Context) error {

	metadata, err := jwtGetAccessMetadata(c)
	if err != nil {
		return err
	}

	u, err := GetUserByID(c, data.DB(), metadata.userID)
	if err != nil {
		return err
	}

	if metadata.expiresAt.Before(time.Now()) ||
		metadata.createdAt.Before(u.LoggedOutAt) {
		return errors.New("access token expired")
	}

	return nil

}

// jwtGetAccessMetadata extracts metdata from the request access token.
func jwtGetAccessMetadata(c *gin.Context) (*jwtAccessMetadata, error) {

	// parse JWT
	token, err := jwt.Parse(getAccessToken(c),
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v",
					token.Header["alg"])
			}
			return []byte(accessKey), nil
		})
	if err != nil {
		return nil, err
	}

	// define generic error to return return if parsing details fails
	genericErr := errors.New("failed to read JWT metadata")

	// extract claims from JWT
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, genericErr
	}

	authUUID, ok := claims["auth_uuid"].(string)
	if !ok {
		return nil, genericErr
	}

	userID, err := jwtParseIntFromClaims(claims, "user_id")
	if err != nil {
		return nil, genericErr
	}

	createdAtUnix, err := jwtParseIntFromClaims(claims, "created_at")
	if err != nil {
		return nil, genericErr
	}

	expiresAtUnix, err := jwtParseIntFromClaims(claims, "expires_at")
	if err != nil {
		return nil, genericErr
	}

	return &jwtAccessMetadata{
		authUUID:  authUUID,
		userID:    uint(userID),
		createdAt: time.Unix(int64(createdAtUnix), 0),
		expiresAt: time.Unix(int64(expiresAtUnix), 0),
	}, nil
}

// getAccessToken retrieves the bearer auth token from the supplied request.
func getAccessToken(c *gin.Context) string {

	tokenParts := strings.Split(c.Request.Header.Get("Authorization"), " ")

	if len(tokenParts) == 2 {
		return tokenParts[1]
	}

	return ""

}

// jwtGetRefreshMetadata extracts metdata from the supplied refresh token.
func jwtGetRefreshMetadata(c *gin.Context,
	refreshToken string) (*jwtRefreshMetadata, error) {

	// parse JWT
	token, err := jwt.Parse(refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v",
					token.Header["alg"])
			}
			return []byte(refreshKey), nil
		})
	if err != nil {
		return nil, err
	}

	// define generic error to return return if parsing details fails
	genericErr := errors.New("failed to read JWT metadata")

	// extract claims from JWT
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, genericErr
	}

	authUUID, ok := claims["auth_uuid"].(string)
	if !ok {
		return nil, genericErr
	}

	userID, err := jwtParseIntFromClaims(claims, "user_id")
	if err != nil {
		return nil, genericErr
	}

	createdAtUnix, err := jwtParseIntFromClaims(claims, "created_at")
	if err != nil {
		return nil, genericErr
	}

	expiresAtUnix, err := jwtParseIntFromClaims(claims, "expires_at")
	if err != nil {
		return nil, genericErr
	}

	return &jwtRefreshMetadata{
		authUUID:  authUUID,
		userID:    uint(userID),
		createdAt: time.Unix(int64(createdAtUnix), 0),
		expiresAt: time.Unix(int64(expiresAtUnix), 0),
	}, nil
}

// jwtParseIntFromClaims extracts an integer from the supplied JWT map claims.
func jwtParseIntFromClaims(claims jwt.MapClaims, key string) (int, error) {

	var value int
	var err error

	switch claims[key].(type) {
	case string:
		value, err = strconv.Atoi(claims[key].(string))
	case float64:
		value = int(claims[key].(float64))
	default:
		return 0, fmt.Errorf("valid type for claim '%s'", key)
	}

	if err != nil {
		return 0, fmt.Errorf("invalid format for claim '%s'", key)
	}

	return value, nil

}
