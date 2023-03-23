package c_v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ShikharY10/gbAUTH/cmd/handlers"
	"github.com/ShikharY10/gbAUTH/cmd/middlewares"
	"github.com/ShikharY10/gbAUTH/cmd/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthController struct {
	Handler    *handlers.Handler
	Middleware *middlewares.Middleware
}

func (ac *AuthController) Test(c *gin.Context) {
	c.JSON(200, "Working")
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
			"exp":     time.Now().Add(time.Minute * (60 * 15)).Unix(),
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
		ac.Handler.Cache.RedisClient.Set(id+".email", email, time.Minute*(60*5))
		ac.Handler.Cache.RedisClient.Set(id+".purpose", purpose, time.Minute*(60*5))

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
		err := ac.Handler.DataBase.IsEmailAvailable(email)
		if err != nil {
			c.AbortWithStatusJSON(400, err.Error())
			return
		}
		token, _, err := ac.requestOTP(email, "signup")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		} else {
			c.SetCookie("SUT-AUTHORIZATION", token, 300, "/", "", false, true)
			c.JSON(201, "Successfully Sent")
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
	otp := request["otp"].(string)

	if ac.Handler.Cache.VarifyOTP(id, otp) {
		ac.Handler.Cache.RedisClient.Set(id+".status", "varified", time.Minute*(60*5))
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
			c.AbortWithStatusJSON(400, "username already exist")
		} else {
			c.JSON(200, "username awailable")
		}
	}
}

func (ac *AuthController) SignUp(c *gin.Context) {
	// setting response headers
	c.Header("service", "Gossip API")

	tokenId := c.Value("tokenid").(string)
	result := ac.Handler.Cache.RedisClient.Get(tokenId + ".status")
	if result.Val() == "varified" {

		// collecting request body
		var requestN models.SignupRequest
		c.BindJSON(&requestN)

		if err := requestN.Examine(); err != nil {
			c.AbortWithStatusJSON(400, err.Error())
			return
		}

		err := ac.Handler.DataBase.IsUsernameAwailable(requestN.Username)
		if err != nil {
			c.AbortWithStatusJSON(400, err.Error())
			return
		}

		avatar, err := ac.Handler.Cloudinary.UploadUserAvatar(tokenId, requestN.AvatarData, requestN.AvatarExt)
		if err != nil {
			c.AbortWithStatusJSON(500, "1. something went wrong, "+err.Error())
			return
		}

		deliveryId, err := ac.Handler.DataBase.AddUserPayloadsField()
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		objectId := primitive.NewObjectID()

		var user models.User
		user.Avatar = *avatar
		user.CreatedAt = time.Now().Format(time.RFC822)
		user.DeletedAt = ""
		user.Email = c.Value("email").(string)
		user.Partners = []primitive.ObjectID{}
		user.PartnerRequests = []models.PartnerRequest{}
		user.PartnerRequested = []models.PartnerRequest{}
		user.Posts = []primitive.ObjectID{}
		user.ID = objectId
		user.DeliveryId = *deliveryId
		user.Name = requestN.Name
		user.Role = "user"
		user.UpdatedAt = time.Now().Format(time.RFC822)
		user.Username = requestN.Username

		// generating access token
		accessClaim := map[string]interface{}{
			"id":       objectId.Hex(),
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}
		accessToken, err := ac.Middleware.GenerateJWT(accessClaim, "access")
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		// generating refresh token
		refreshClaim := map[string]interface{}{
			"id":  objectId.Hex(),
			"exp": time.Now().AddDate(1, 0, 0).Unix(),
		}
		refreshToken, err := ac.Middleware.GenerateJWT(refreshClaim, "refresh")
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		ac.Handler.Cache.SetAccessTokenExpiry(objectId.Hex(), accessToken, 1*time.Hour)
		ac.Handler.Cache.SetRefreshTokenExpiry(objectId.Hex(), refreshToken, 24*time.Hour)

		c.SetCookie("refresh", refreshToken, 3600*24, "/", "", false, true)

		err = ac.Handler.DataBase.CreateNewUser(user)
		err1 := ac.Handler.DataBase.InsetUserInFrequencyTable(user.ID, user.Username)
		if err != nil || err1 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "4. something went wrong, "+err.Error())
		} else {
			user.AccessToken = accessToken
			c.SetCookie("SUT-AUTHORIZATION", "", -1, "/", "", false, true)
			c.JSON(http.StatusCreated, user)
			ac.Handler.Cache.RedisClient.Del(tokenId)
			ac.Handler.Cache.RedisClient.Del(tokenId + ".auth")
			ac.Handler.Cache.RedisClient.Del(tokenId + ".status")
			ac.Handler.Cache.RedisClient.Del(tokenId + ".purpose")
		}

	} else {
		c.AbortWithStatusJSON(401, "you are not verified")
	}
}

func (ac *AuthController) RefreshAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(400)
		return
	}

	isTokenValid := ac.Handler.Cache.IsTokenValid(id, refreshToken, "refresh")
	if isTokenValid {
		refreshClaims, err := ac.Middleware.VarifyRefreshToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(401, "logged out")
			return
		}
		_id, err := primitive.ObjectIDFromHex(refreshClaims["id"].(string))
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		opts := options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 1},
			{Key: "username", Value: 1},
			{Key: "role", Value: 1},
		})

		user, err := ac.Handler.DataBase.GetUserData(bson.M{"_id": _id}, opts)
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		newAccessTokenClaim := map[string]interface{}{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}
		accessToken, err := ac.Middleware.GenerateJWT(newAccessTokenClaim, "access")
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			return
		}

		ac.Handler.Cache.SetAccessTokenExpiry(user.ID.Hex(), accessToken, time.Hour*1)

		c.JSON(200, map[string]string{
			"accessToken": accessToken,
		})
	} else {
		c.AbortWithStatusJSON(401, "logged out")
		return
	}
}

func (ac AuthController) LogOut(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	username := c.Value("username").(string)
	id := c.Value("id").(string)

	ac.Handler.Cache.DeleteTokenExpiry(id)
	c.SetCookie("refresh", "", -1, "/", "", false, true)
	c.JSON(200, "Successfully Logout")

	result := ac.Handler.DataBase.UpdateLogoutStatus(username, true)
	if result == nil {
		c.JSON(http.StatusCreated, "sucessfully logged out")
	} else {
		c.AbortWithStatus(http.StatusPreconditionFailed)
	}
}

func (ac *AuthController) RequestOtpForLogin(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// {
	// 	  "type": "email"
	//    "email": "yshikharfzd10@gmail.com"
	// }

	// {
	//    "type": "username",
	// 	  "username": "shikhary10"
	// }

	// collecting request body
	var request models.RequestLoginRequest
	c.BindJSON(&request)
	if err := request.Examine(); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	var email string
	var err error
	if request.Type == "username" {
		if request.Username != "" {
			email, err = ac.Handler.DataBase.GetUserEmail(request.Username)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, "wrong username")
			}
		}

	} else if request.Type == "email" {
		if request.Email != "" {
			email = request.Email
		}
	}
	if email == "" {
		c.AbortWithStatusJSON(500, "email not found")
		return
	}
	token, _, err := ac.requestOTP(email, "login")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
	} else {
		c.SetCookie("SUT-AUTHORIZATION", token, 300, "/", "", false, true)
		c.JSON(201, "Successfully Sent")
	}
}

func (ac *AuthController) LogIn(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	otp := request["otp"].(string)
	if !ac.Handler.Cache.VarifyOTP(c.Value("tokenid").(string), otp) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "wrong OTP")
	} else {
		opts := options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 1},
			{Key: "name", Value: 1},
			{Key: "username", Value: 1},
			{Key: "email", Value: 1},
			{Key: "avatar", Value: 1},
			{Key: "deliveryId", Value: 1},
			{Key: "posts", Value: 1},
			{Key: "partners", Value: 1},
			{Key: "partnerrequests", Value: 1},
			{Key: "partnerrequested", Value: 1},
		})
		user, err := ac.Handler.DataBase.GetUserData(
			bson.M{"email": c.Value("email").(string)},
			opts,
		)
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
		} else {
			c.JSON(200, user)
		}
	}
}

// user token based update
func (ac *AuthController) UpdateUserName(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	id := c.Value("id").(string)

	// name
	fullName := request["fullname"].(string)
	if fullName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "name not found")
	} else {
		err := ac.Handler.DataBase.UpdateUserDetail(id, "name", fullName)
		// err := ac.Handler.DataBase.UpdateUserName(id, fullName)
		if err == nil {
			response := map[string]string{
				"name": fullName,
			}
			c.JSON(201, response)
		} else {
			c.AbortWithStatusJSON(500, err.Error())
		}
	}
}

// user token based update
func (ac *AuthController) UpdateAvatar(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	id := c.Value("id").(string)

	// name
	imageData := request["imagedata"].(string)
	imageExt := request["imageext"].(string)

	if imageData != "" && imageExt != "" {

		avatar, err := ac.Handler.Cloudinary.UploadUserAvatar(id+"temp", imageData, imageExt)
		if err != nil {
			c.AbortWithStatus(http.StatusPreconditionFailed)
		} else {
			if err != nil {
				c.AbortWithStatus(http.StatusPreconditionFailed)
			} else {
				err := ac.Handler.DataBase.UpdateUserDetail(id, "avatar", *avatar)
				// err := ac.Handler.DataBase.UpdateUserAvatar(uuid, *avatar)
				if err == nil {
					c.JSON(201, "Successfully updated")
				} else {
					c.AbortWithStatus(http.StatusPreconditionFailed)
				}
			}
		}
	} else {
		c.AbortWithStatus(http.StatusPreconditionFailed)
	}
}

// special token based update
func (ac *AuthController) UpdateUsername(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	if request == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
	} else {
		id := c.Value("id").(string)
		username := request["username"].(string)
		email, err := ac.Handler.DataBase.GetUserEmail(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong: "+err.Error())
		} else {
			token, id, err := ac.requestOTP(email, "username")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong: "+err.Error())
			} else {
				ac.Handler.Cache.RedisClient.Set(id+".updateusername", username, time.Duration(time.Minute*(60*5)))
				response := map[string]string{
					"token": token,
				}
				c.JSON(http.StatusCreated, response)
			}
		}
	}
}

func (ac *AuthController) VarifyUsernameUpdateOTP(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	if request == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
	} else {
		tokenID := c.Value("tokenid").(string)
		id := c.Value("id").(string)
		otp := request["otp"].(string)

		if !ac.Handler.Cache.VarifyOTP(tokenID, otp) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "wrong OTP")
		} else {
			newUsername := ac.Handler.Cache.RedisClient.Get(tokenID + ".updateusername").Val()
			err := ac.Handler.DataBase.UpdateUserDetail(id, "username", newUsername)
			if err == nil {
				response := map[string]string{
					"username": newUsername,
				}
				c.JSON(http.StatusOK, response)
			} else {
				c.AbortWithStatusJSON(500, err.Error())
			}
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
			ac.Handler.Cache.RedisClient.Set(id2+".updateemail", newEmail, time.Minute*(60*5))
			ac.Handler.Cache.RedisClient.Set(id1+id2+".purpose", "email", time.Minute*(60*5))
			return token, nil
		}
	}
}

// special token based update
func (ac *AuthController) UpdateEmail(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	if request == nil {
		c.AbortWithStatusJSON(400, "body not found")
	} else {
		id := c.Value("id").(string)
		newEmail := request["newemail"].(string)
		oldEmail, err := ac.Handler.DataBase.GetUserEmail(id)

		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
		} else {
			token, err := ac.requestEmailUpdateOTP(oldEmail, newEmail)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
			} else {
				c.SetCookie("SUT-AUTHORIZATION", token, 300, "/", "", false, true)
				c.JSON(200, "Successfullu Sent")
			}
		}
	}
}

func (ac *AuthController) VarifyEmailUpdateOTP(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request map[string]any
	c.BindJSON(&request)

	if request == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
	} else {
		// oldemailotp
		// newemailotp

		oldEmailOTP := request["oldemailotp"].(string)
		newEmailOTP := request["newemailotp"].(string)

		tokenID1 := c.Value("tokenid1").(string)
		tokenID2 := c.Value("tokenid2").(string)

		if ac.Handler.Cache.VarifyOTP(tokenID1, oldEmailOTP) && ac.Handler.Cache.VarifyOTP(tokenID2, newEmailOTP) {
			id := c.Value("id").(string)
			newEmail := ac.Handler.Cache.RedisClient.Get(tokenID2 + "_updateemail").Val()
			err := ac.Handler.DataBase.UpdateUserDetail(id, "email", newEmail)
			if err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			} else {
				response := map[string]string{
					"email": newEmail,
				}
				c.JSON(http.StatusOK, response)
			}
		} else {
			c.AbortWithStatusJSON(401, "Invalid OTP")
		}
	}
}
