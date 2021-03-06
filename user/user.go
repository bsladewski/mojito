package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"mojito/data"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

// accessKey is used to sign JWT access tokens.
var accessKey string

// refreshKey is used to sign JWT refresh tokens.
var refreshKey string

// authExpirationHours determines the number of hours before we consider an
// access token to be expired.
var accessExpirationHours time.Duration

// refreshExpirationHours determines the number of hours before we consider a
// refresh token to be expired.
var refreshExpirationHours time.Duration

// CreateAuth generates JWT access and refresh tokens for the supplied user.
func CreateAuth(ctx context.Context, u *User) (accessToken,
	refreshToken string, err error) {

	// generate UUID to track issued credentials in peristent storage
	authUUID := uuid.NewV4().String()

	// create the access token
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth_uuid":  authUUID,
		"user_id":    u.ID,
		"created_at": time.Now().Unix(),
		"expires_at": time.Now().Add(accessExpirationHours).Unix(),
	})

	accessToken, err = accessJWT.SignedString([]byte(accessKey))
	if err != nil {
		return "", "", err
	}

	// create the refresh token
	refreshExpiration := time.Now().Add(refreshExpirationHours)

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth_uuid":  authUUID,
		"user_id":    u.ID,
		"created_at": time.Now().Unix(),
		"expires_at": refreshExpiration.Unix(),
	})

	refreshToken, err = refreshJWT.SignedString([]byte(refreshKey))
	if err != nil {
		return "", "", err
	}

	// add the user auth record
	if err := SaveLogin(ctx, data.DB(), &Login{
		UserID:    u.ID,
		UUID:      authUUID,
		ExpiresAt: refreshExpiration,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

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
