package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

func HandleInternalError(c *gin.Context, err error) {
	handle(http.StatusInternalServerError, c, err)
}

// HandleBadRequest writes http.StatusBadRequest and log error
func HandleBadRequest(c *gin.Context, err error) {
	handle(http.StatusBadRequest, c, err)
}

func HandleNotFound(c *gin.Context, err error) {
	handle(http.StatusNotFound, c, err)
}

func HandleForbidden(c *gin.Context, err error) {
	handle(http.StatusForbidden, c, err)
}

func HandleUnauthorized(c *gin.Context, err error) {
	handle(http.StatusUnauthorized, c, err)
}

func HandleTooManyRequests(c *gin.Context, err error) {
	handle(http.StatusTooManyRequests, c, err)
}

func HandleConflict(c *gin.Context, err error) {
	handle(http.StatusConflict, c, err)
}

func HandleError(c *gin.Context, err error) {
	var statusCode int
	switch t := err.(type) {
	case errors.APIStatus:
		statusCode = int(t.Status().Code)
	default:
		statusCode = http.StatusInternalServerError
	}
	handle(statusCode, c, err)
}

func handle(statusCode int, c *gin.Context, err error) {
	_, fn, line, _ := runtime.Caller(2)
	klog.Errorf("%s:%d %v", fn, line, err)
	c.JSON(statusCode, err)
}
