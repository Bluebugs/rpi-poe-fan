package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIEndpoints(ctx *gin.Context) {
	if !strings.HasSuffix(ctx.Request.URL.Path, "/api/") &&
		ctx.GetHeader("Content-Type") == "application/json" {
		ctx.Error(fmt.Errorf("invalid content type"))
		return
	}
	ctx.Next()
}
