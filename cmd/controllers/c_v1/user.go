package c_v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ShikharY10/gbAUTH/cmd/handlers"
	"github.com/ShikharY10/gbAUTH/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Handler    *handlers.Handler
	Middleware *middlewares.Middleware
}

func (ac *AuthController) requestOTP(email string, purpose string) (string, string, error) {
	// generating random OTP and UID
	id, otp := ac.Handler.Cache.RegisterOTP()

	// sending OTP to OTP_SERVICE for sending it to the user
	var otpData map[string]string = map[string]string{
		"otp":     otp,
		"email":   email,
		"purpose": purpose,
	}
	_, err := json.Marshal(otpData)
	if err == nil {
		// produced to the queue that is listened by OTP_SERVICE
		// v3.Handler.QueueHandler.Produce("OTPd3hdzl8", b)
		fmt.Println("OTP: ", otp) // temp

		// generating authorization token
		claim := map[string]interface{}{
			"exp":     time.Now().Add(time.Minute * (60 * 5)).Unix(),
			"tokenid": id,
			"email":   email,
			"purpose": purpose,
		}
		token, err := ac.Middleware.GenerateJWT(claim, "update")
		fmt.Println(len(token), " | Token Generated: ", token)
		if err != nil {
			return "", "", err
		}

		// storing id and number for future authorization.
		ac.Handler.Cache.RedisClient.Set(id+"_id", email, time.Minute*(60*5))
		ac.Handler.Cache.RedisClient.Set(id+"_purpose", purpose, time.Minute*(60*5))

		return token, id, nil
	} else {
		return "", "", err
	}
}

func (ac *AuthController) RequestOtpForSignup(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)
	email := request["email"].(string)
	if email != "" {
		token, _, err := ac.requestOTP(email, "signup")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		} else {
			response := map[string]string{
				"token": token,
			}
			c.JSON(http.StatusCreated, response)
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, "email not found")
	}
}

func (ac *AuthController) VarifySignupOTP(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	id := c.Value("tokenid").(string)
	fmt.Println("id from VarifyOTP: ", id)
	otp := request["otp"].(string)

	if ac.Handler.Cache.VarifyOTP(id, otp) {
		ac.Handler.Cache.RedisClient.Set(id+"_status", "varified", time.Minute*(60*5))
		response := map[string]string{
			"status": "successful",
		}
		c.JSON(http.StatusAccepted, response)
	} else {
		response := map[string]string{
			"status": "unsucessful",
		}
		c.JSON(http.StatusNotAcceptable, response)
	}
}

func (ac *AuthController) IsUsernameAwailable(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	username := c.Query("username")
	if username == "" {
		c.AbortWithStatus(400)
	} else {
		err := ac.Handler.DataBase.IsUsernameAwailable(username)
		if err != nil {
			ac.Handler.Logger.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(200, "")
		}
	}
}

func (ac *AuthController) requestEmailUpdateOTP(oldEmail string, newEmail string) (string, error) {
	id1, otp1 := ac.Handler.Cache.RegisterOTP()
	id2, otp2 := ac.Handler.Cache.RegisterOTP()

	var otpData map[string]string = map[string]string{
		"purpose":     "email",
		"oldemail":    oldEmail,
		"oldemailotp": otp1,
		"newemail":    newEmail,
		"newEmailotp": otp2,
	}
	_, err := json.Marshal(otpData)
	if err != nil {
		return "", err
	} else {
		fmt.Println("Old Email: ", oldEmail, "OTP: ", otp1)
		fmt.Println("New Email: ", newEmail, "OTP: ", otp2)

		claim := map[string]interface{}{
			"exp":      time.Now().Add(time.Minute * (60 * 5)).Unix(),
			"tokenid1": id1,
			"tokenid2": id2,
			"purpose":  "email",
		}
		token, err := ac.Middleware.GenerateJWT(claim, "update")
		if err != nil {
			return "", err
		} else {
			ac.Handler.Cache.RedisClient.Set(id2+"_updateemail", newEmail, time.Minute*(60*5))
			ac.Handler.Cache.RedisClient.Set(id1+id2+"_purpose", "email", time.Minute*(60*5))
			return token, nil
		}
	}
}
