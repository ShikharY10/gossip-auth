// this file contains code that is used to do role based based authorization and authenticated user.
package middlewares

import "github.com/gin-gonic/gin"

// Middleware used to authorize a used based on its role
func RoleBasedAccess(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Value("role") != nil {
			savedRole := c.Value("role").(string)
			if savedRole == role {
				c.Next()
			} else {
				c.AbortWithStatusJSON(401, "invalid role")
			}
		} else {
			c.AbortWithStatusJSON(401, "role not found")
		}
	}
}
