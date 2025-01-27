package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
)

func (a *API) Auth(c *gin.Context) {
	userId := c.Request.Header.Get("AuthUserId")
	if userId == "" {
		c.AbortWithError(http.StatusForbidden, errors.AuthorizationError)
		return
	}
	c.Set("userId", userId)

	/* TODO
	signature := c.Request.Header.Get("AuthSign")
	if signature == "" {
		c.AbortWithError(http.StatusForbidden, errors.AuthorizationError)
		return
	}

	if !checkSignature(publicKey, signature) {
		c.AbortWithError(http.StatusForbidden, errors.AuthorizationError)
		return
	}
	*/

	c.Next()
}

func checkSignature(pubKey, signature string) bool {
	//TODO with NaCl sign

	return true
}
