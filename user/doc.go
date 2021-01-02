// Package user provides functionality for managing user accounts, permissions,
// and authentication.
//
// Environment:
//     MOJITO_ACCESS_KEY:
//         string - the key used to sign JWT access tokens
//     MOJITO_REFRESH_KEY:
//         string - the key used to sign JWT refresh tokens
//     MOJITO_ACCESS_EXPIRATION_HOURS:
//         int - the number of hours before an access token is expired
//     MOJITO_REFRESH_EXPIRATION_HOURS:
//         int - the number of hours before a refresh token is expired
package user
