package main

import (
	"github.com/gin-gonic/gin"
)

// HandlerFuncWithApp is a type for handler functions that require the application object.
type HandlerFuncWithApp func(*gin.Context, *application)

// InjectApp injects the application object into a handler function.
func InjectApp(app *application, handler HandlerFuncWithApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(ctx, app)
	}
}
