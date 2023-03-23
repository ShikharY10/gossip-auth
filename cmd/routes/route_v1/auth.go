package route_v1

import (
	"github.com/ShikharY10/gbAUTH/cmd/controllers/c_v1"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller *c_v1.AuthController) {
	router.GET("/test", controller.Test)

	// Signup Routes
	router.POST("/requestsignupotp", controller.RequestOtpForSignup)
	router.POST("varifysignupotp", controller.Middleware.SingleUseTokenVarification("signup", false), controller.VarifySignupOTP)
	router.GET("/isusernameawailable", controller.IsUsernameAwailable)
	router.POST("/signup", controller.Middleware.SingleUseTokenVarification("signup", false), controller.SignUp)

	// Login Routes
	router.POST("requestloginotp", controller.RequestOtpForLogin)
	router.POST("/login", controller.Middleware.SingleUseTokenVarification("login", true), controller.LogIn)

	// refresh token router
	router.PUT("/refreshaccesstoken/:id", controller.RefreshAccessToken)

	// Secured Routes
	authorizedRoutes := router.Group("/")
	authorizedRoutes.Use(controller.Middleware.APIV1_Authorization())

	authorizedRoutes.GET("/auth", func(c *gin.Context) {
		c.JSON(200, "Working")
	})

	// Update Routes
	authorizedRoutes.PUT("updateavatar", controller.UpdateAvatar)
	authorizedRoutes.PUT("updatedname", controller.UpdateUserName)

	// Critical Update Routes
	authorizedRoutes.PUT("/updateusername", controller.UpdateUsername)
	authorizedRoutes.POST("/varifyusernameupdate", controller.Middleware.SingleUseTokenVarification("username", true), controller.VarifyUsernameUpdateOTP)
	authorizedRoutes.PUT("/updateemail", controller.UpdateEmail)
	authorizedRoutes.POST("/varifyemailupdate", controller.Middleware.APIV3EmailUpdateVarification(), controller.VarifyEmailUpdateOTP)

	// Logout Route
	authorizedRoutes.DELETE("/logout", controller.LogOut)

}
