package routes

import (
	"akh/blog/controllers/back"

	"github.com/gin-gonic/gin"
)

func addUserRoute(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.POST("/", back.Cteate)
	user.POST("/login", back.Login)
}

func backRoute(rg *gin.RouterGroup) {
	blog := rg.Group("/back")
	blog.POST("/blog", back.BlogCreate)
}
