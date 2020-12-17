package delivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bsladewski/mojito/auth"
	"github.com/bsladewski/mojito/httperror"
	"github.com/bsladewski/mojito/server"
	"github.com/bsladewski/mojito/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// init registers the auth API with the application router.
func init() {

	// bind public endpoints
	server.Router().POST(loginEndpoint, login)
	server.Router().POST(refreshEndpoint, refresh)

	// bind private endpoints
	server.Router().POST(logoutEndpoint, auth.JWTAuthMiddleware(), logout)
}

const (
	// loginEndpoint the API endpoint that handles user login.
	loginEndpoint = "/login"
	// refreshEndpoint the API endpoint that handles refreshing access tokens.
	refreshEndpoint = "/refresh"
	// logoutEndpoint the API endpoint that handles user logout.
	logoutEndpoint = "/logout"
	// invalidUserCredentials is an error message returned when the user's email
	// or password is incorrect.
	invalidUserCredentials = "invalid email or password"
	// logoutFailedGeneric is a generic error returned when user logout fails.
	logoutFailedGeneric = "failed to log out user"
	// invalidRefreshToken is an error message returned if the user supplies an
	// invalid refresh token or a refresh token that is inconsistent with
	// persistent data.
	invalidRefreshToken = "invalid refresh token"
)

// login checks user credentials and generates access and refresh tokens for
// authenticating user requests.
func login(c *gin.Context) {

	var req loginRequest

	// read user credentials from request body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// validate request parameters
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "email is required",
		})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "password is required",
		})
		return
	}

	// retrieve user account by email address
	u, err := user.GetUserByEmail(c, req.Email)
	if err == gorm.ErrRecordNotFound {
		logrus.Warn(err)
		c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
			ErrorMessage: invalidUserCredentials,
		})
		return
	} else if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// compare supplied password with user password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(fmt.Sprintf("%d:%s", u.ID, req.Password)),
	); err != nil {
		logrus.Debug(err)
		c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
			ErrorMessage: invalidUserCredentials,
		})
		return
	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := auth.CreateAuth(c, u)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// delete expired user login records to keep persistent storage clean
	go func() {
		if err := auth.DeleteExpiredLogin(c, u.ID); err != nil {
			logrus.Error(err)
		}
	}()

	// repond with auth tokens
	c.JSON(http.StatusOK, loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

// refresh checks the supplied refresh token and generates new access and
// refresh tokens if valid.
func refresh(c *gin.Context) {

	var req refreshRequest

	// read user credentials from request body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// validate request parameters
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "refresh token is required",
		})
		return
	}

	// validate the supplied refresh token
	login, err := auth.JWTValidateRefreshToken(c, req.RefreshToken)
	if err != nil {
		logrus.Warn(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidRefreshToken,
		})
		return
	}

	// retrieve user record
	u, err := user.GetUserByID(c, login.UserID)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidRefreshToken,
		})
		return
	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := auth.CreateAuth(c, u)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// delete original refresh token
	if err := auth.DeleteLogin(c, login); err != nil {
		logrus.Error(err)
	}

	// repond with auth tokens
	c.JSON(http.StatusOK, refreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

// logout invalidates the logged in user's access and refresh tokens.
func logout(c *gin.Context) {

	// get user from JWT
	u, err := auth.JWTGetUser(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
			ErrorMessage: logoutFailedGeneric,
		})
		return
	}

	// get user auth record from JWT
	login, err := auth.JWTGetUserLogin(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
			ErrorMessage: logoutFailedGeneric,
		})
		return
	}

	// delete user auth record, this will invalidate the refresh token
	if err := auth.DeleteLogin(c, login); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: logoutFailedGeneric,
		})
		return
	}

	// set logged out at time, this will invalidate all access tokens issued
	// before this time
	u.LoggedOutAt = time.Now()

	// update the user record
	if err := user.SaveUser(c, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: logoutFailedGeneric,
		})
		return
	}

	// respond with 200 - OK if logout was successful
	c.Status(http.StatusOK)

}
