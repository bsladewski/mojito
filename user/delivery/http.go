package delivery

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/bsladewski/mojito/auth"
	"github.com/bsladewski/mojito/email"
	"github.com/bsladewski/mojito/httperror"
	"github.com/bsladewski/mojito/server"
	"github.com/bsladewski/mojito/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// init registers the user API with the application router.
func init() {

	// bind public endpoints
	server.Router().POST(signupEndpoint, signup)
	server.Router().POST(verifyEndpoint, verify)
	server.Router().POST(recoverEndpoint, recover)
	server.Router().POST(recoverResetEndpoint, recoverReset)

	// bind private endpoints
	server.Router().POST(resetEndpoint, auth.JWTAuthMiddleware(), reset)
}

const (
	// signupEndpoint the API endpoint used to create new user accounts.
	signupEndpoint = "/signup"
	// verifyEndpoint the API endpoint used to verify a new user's email
	// address.
	verifyEndpoint = "/verify"
	// recoverEndpoint the API endpoint used to send account recovery emails.
	recoverEndpoint = "/recover"
	// recoverResetEndpoint the API endpoint for resetting an account password
	// as part of the account recovery process.
	recoverResetEndpoint = "/recover/reset"
	// resetEndpoint the API endpoint used to reset the logged in user's
	// password.
	resetEndpoint = "/reset"
	// invalidToken is an error returned if if a user validation token is
	// supplied that cannot be parsed or contains invalid data.
	invalidToken = "invalid token"
	// resetFailedGeneric is a generic error message returned when resetting
	// the user account password fails.
	resetFailedGeneric = "failed to reset password"
)

// signup creates a new user account.
func signup(c *gin.Context) {

	var req signupRequest

	// read request parameters
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

	// check if a verified user account with the same email address already
	// exists
	u, err := user.GetUserByEmail(c, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	} else if u != nil && u.Verified {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "email address is already registered",
		})
		return
	}

	// if no unverified user account exists, create a new user account
	if u == nil {

		// generate user secret key
		secretKey := md5.Sum(uuid.NewV4().Bytes())

		u = &user.User{
			Email:     req.Email,
			SecretKey: fmt.Sprintf("%x", secretKey),
		}

		// create the user account record
		if err := user.SaveUser(c, u); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
				ErrorMessage: httperror.InternalServerError,
			})
			return
		}

	}

	// set user password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%d:%s", u.ID, req.Password)), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	u.Password = string(hash)

	if err := user.SaveUser(c, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// generate the verification token
	token, err := user.GenerateSecretToken(c, u, u.Email)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// send the verification email
	if err := email.SendEmailTemplate(
		email.DefaultFromAddress(),
		email.DefaultReplyToAddress(),
		[]string{u.Email},
		nil,
		nil,
		email.TemplateTitleSignup,
		signupEmailData{
			ClientHost:        server.ClientHost(),
			VerificationToken: token,
		},
	); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: "failed to send verification email, please try again later",
		})
	}

}

// verify checks the supplied verification token to determine if the user has
// access to the account email address.
func verify(c *gin.Context) {

	var req verifyRequest

	// read request parameters
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// decode the verification token
	u, payload, err := user.ParseSecretToken(c, req.Token)
	if err != nil {
		logrus.Warn(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidToken,
		})
		return
	}

	// validate the token payload
	if payload != u.Email {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidToken,
		})
		return
	}

	// mark the user record as verified
	u.Verified = true

	// save user record
	if err := user.SaveUser(c, u); err != nil {
		logrus.WithError(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidToken,
		})
		return
	}

	// response with 200 - OK if verification was successful
	c.Status(http.StatusOK)

}

// recover sends an email to the user with a link to reset the user account
// password.
func recover(c *gin.Context) {

	var req recoverRequest

	// read request parameters
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// retrieve user account by email address
	u, err := user.GetUserByEmail(c, req.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "email address not found",
		})
		return
	} else if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// generate the verification token
	token, err := user.GenerateSecretToken(c, u, u.Email)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// send the verification email
	if err := email.SendEmailTemplate(
		email.DefaultFromAddress(),
		email.DefaultReplyToAddress(),
		[]string{u.Email},
		nil,
		nil,
		email.TemplateTitleRecover,
		recoverEmailData{
			ClientHost:        server.ClientHost(),
			VerificationToken: token,
		},
	); err != nil {
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: "failed to send verification email, please try again later",
		})
	}

}

// recoverReset is used to change a user account password as part of the account
// recovery process.
func recoverReset(c *gin.Context) {

	var req recoverResetRequest

	// read request parameters
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// decode the verification token
	u, payload, err := user.ParseSecretToken(c, req.Token)
	if err != nil {
		logrus.Warn(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidToken,
		})
		return
	}

	// validate the token payload
	if payload != u.Email {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: invalidToken,
		})
		return
	}

	// set user password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%d:%s", u.ID, req.Password)), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	u.Password = string(hash)

	if err := user.SaveUser(c, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// response with 200 - OK if password reset was successful
	c.Status(http.StatusOK)

}

// reset is used to change the logged in user's account password.
func reset(c *gin.Context) {

	// get user from JWT
	u, err := auth.JWTGetUser(c)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, httperror.ErrorResponse{
			ErrorMessage: resetFailedGeneric,
		})
		return
	}

	var req resetRequest

	// read request parameters
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "invalid request body",
		})
		return
	}

	// validate request parameters
	if req.CurrentPassword == "" {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "current password is required",
		})
		return
	}

	if req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "new password is required",
		})
		return
	}

	// verify current password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(fmt.Sprintf("%d:%s", u.ID, req.CurrentPassword)),
	); err != nil {
		logrus.Debug(err)
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "current password is incorrect",
		})
		return
	}

	// check that current password is not the same as the new password
	if req.CurrentPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, httperror.ErrorResponse{
			ErrorMessage: "new and current passwords are the same",
		})
		return
	}

	// set user password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%d:%s", u.ID, req.NewPassword)), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	u.Password = string(hash)

	if err := user.SaveUser(c, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// response with 200 - OK if password reset was successful
	c.Status(http.StatusOK)

}
