// This files contains all the function and code used to authorize a authenticated user.
package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	config "github.com/ShikharY10/gbAUTH/cmd/configs"
	"github.com/ShikharY10/gbAUTH/cmd/handlers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	jWT_ACCESS_TOKEN_SECRET_KEY  []byte
	jWT_REFRESH_TOKEN_SECRET_KEY []byte
	DataBase                     *handlers.DataBase
	Cache                        *handlers.Cache
}

// Initializes JWT struct
func InitializeMiddleware(env *config.ENV, database *handlers.DataBase, cache *handlers.Cache) *Middleware {
	return &Middleware{
		jWT_ACCESS_TOKEN_SECRET_KEY:  []byte(env.JWT_ACCESS_TOKEN_SECRET_KEY),
		jWT_REFRESH_TOKEN_SECRET_KEY: []byte(env.JWT_REFRESH_TOKEN_SECRET_KEY),
		DataBase:                     database,
		Cache:                        cache,
	}
}

// Creates a JWT token using SHA256 hashing algorithm.
func (j *Middleware) GenerateJWT(claim map[string]interface{}, tokenType string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range claim {
		claims[k] = v
	}

	var key []byte
	if tokenType == "access" {
		key = j.jWT_ACCESS_TOKEN_SECRET_KEY
	} else if tokenType == "refresh" {
		key = j.jWT_REFRESH_TOKEN_SECRET_KEY
	} else if tokenType == "update" {
		key = append(j.jWT_ACCESS_TOKEN_SECRET_KEY, j.jWT_REFRESH_TOKEN_SECRET_KEY...)
	}

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// varifies JWT access token and the claims the where set while creating the token
func (j *Middleware) VarifyAccessToken(token string) (claim jwt.MapClaims, err error) {
	newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("something went wrong")
		}
		return j.jWT_ACCESS_TOKEN_SECRET_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("bad token")
	}
}

// varifies JWT refresh token and the claims the where set while creating the token
func (j *Middleware) VarifyRefreshToken(token string) (claim jwt.MapClaims, err error) {
	newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("something went wrong")
		}
		return j.jWT_REFRESH_TOKEN_SECRET_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("bad token")
	}
}

func (mw *Middleware) SingleUseTokenVarification(use string, delete bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		token, err := c.Cookie("SUT-AUTHORIZATION")
		if err != nil {
			c.AbortWithStatusJSON(401, "token not found")
		} else {
			newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("something went wrong")
				}
				key := append(mw.jWT_ACCESS_TOKEN_SECRET_KEY, mw.jWT_REFRESH_TOKEN_SECRET_KEY...)
				return key, nil
			})
			if err != nil {
				c.AbortWithStatusJSON(401, "invalid SUT token")
				return
			}

			if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
				result1 := mw.Cache.RedisClient.Get(claims["tokenid"].(string) + ".email")
				email := result1.Val()

				result2 := mw.Cache.RedisClient.Get(claims["tokenid"].(string) + ".purpose")
				purpose := result2.Val()

				if email == claims["email"].(string) && purpose == claims["purpose"].(string) {
					c.Set("tokenid", claims["tokenid"].(string))
					c.Set("email", claims["email"].(string))
					if delete {
						c.SetCookie("SUT-AUTHORIZATION", "", -1, "/", "", false, true)
					}

					c.Next()
					return
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, "token integrity compromised")
					return
				}
			} else {
				c.AbortWithStatusJSON(401, "invalid SUT token")
				return
			}
		}
	}
}

func (mw *Middleware) APIV3EmailUpdateVarification() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("New Email Update Request=== " + c.Request.URL.Path + " ===")
		token, err := c.Cookie("SUT-AUTHORIZATION")
		if err == nil {
			newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("something went wrong")
				}
				key := append(mw.jWT_ACCESS_TOKEN_SECRET_KEY, mw.jWT_REFRESH_TOKEN_SECRET_KEY...)
				return key, nil
			})
			if err != nil {
				c.AbortWithStatusJSON(401, "invalid SUT token")
				return
			}

			if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
				tokenID1 := claims["tokenid1"].(string)
				tokenID2 := claims["tokenid2"].(string)

				purpose := mw.Cache.RedisClient.Get(tokenID1 + tokenID2 + ".purpose").Val()

				if purpose == claims["purpose"].(string) {
					c.Set("tokenid1", tokenID1)
					c.Set("tokenid2", tokenID2)
					fmt.Println("===Email Update Request Varified ===")
					c.Next()
					return
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, "token data compromised")
					return
				}
			}
		} else {
			c.AbortWithStatusJSON(401, "token not found")
		}
	}
}

// Middleware for authorizing user using access token
func (j *Middleware) APIV1_Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
			return
		} else {
			token := bearer[len("Bearer "):]
			if token == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
				return
			} else {
				claim, err := j.VarifyAccessToken(token)
				if err != nil {
					if err.Error() == "Token is expired" {
						c.AbortWithStatusJSON(401, err.Error())
					} else {
						c.AbortWithStatusJSON(400, err.Error())
					}
				} else {
					isTokenValid := j.Cache.IsTokenValid(claim["id"].(string), token, "access")
					if isTokenValid {
						data := map[string]interface{}{
							"id":       claim["id"].(string),
							"username": claim["username"].(string),
							"role":     claim["role"].(string),
						}
						c.Keys = data
						c.Next()
					} else {
						c.AbortWithStatus(401)
					}

				}
			}
		}
	}
}

// 4101 1404 2076
