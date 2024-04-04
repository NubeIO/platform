package controller

import (
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/platform/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (inst *Controller) HandleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorized := auth.AuthorizeInternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized := auth.AuthorizeExternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized, _, err := auth.AuthorizeRoles(c.Request, roles...); authorized {
			c.Next()
			return
		} else if err != nil {
			c.JSON(http.StatusUnauthorized, model.Message{Message: err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, model.Message{Message: "token is invalid"})
		c.Abort()
	}
}

func (inst *Controller) HandleUserAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := auth.GetAuthorization(c.Request)
		if len(authorization) > 0 {
			if authorization[0] == "Internal" {
				c.JSON(http.StatusUnauthorized, model.Message{Message: "internal token is restricted"})
				c.Abort()
				return
			} else if authorization[0] == "External" {
				c.JSON(http.StatusUnauthorized, model.Message{Message: "external token is restricted"})
				c.Abort()
				return
			}
		}
		if authorized, _, err := auth.AuthorizeRoles(c.Request, roles...); authorized {
			c.Next()
			return
		} else if err != nil {
			c.JSON(http.StatusUnauthorized, model.Message{Message: err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, model.Message{Message: "token is invalid"})
		c.Abort()
	}
}
