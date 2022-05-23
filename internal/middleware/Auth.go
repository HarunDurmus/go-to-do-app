package custommiddleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthMiddleware struct {
	SecretKey string
	Aud       string
	Iss       string
}

type MyClaims struct {
	jwt.StandardClaims
	Client     string `json:"client"`
	Authorized bool   `json:"authorized"`
}

func (a *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Auth")
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid Signing Method")
			}
			checkAudience := token.Claims.(*MyClaims).VerifyAudience(a.Aud, true)
			if !checkAudience {
				return nil, fmt.Errorf("invalid aud")
			}
			checkIss := token.Claims.(*MyClaims).VerifyIssuer(a.Iss, true)
			if !checkIss {
				fmt.Errorf("invalid iss")
				return nil, c.JSON(http.StatusUnauthorized, http.NoBody)

			}
			checkExp := token.Claims.(*MyClaims).VerifyExpiresAt(time.Now().Add(time.Minute*1).Unix(), true)
			if !checkExp {
				fmt.Errorf("this token expired")
				return nil, c.JSON(http.StatusUnauthorized, http.NoBody)
			}
			return []byte(a.SecretKey), nil
		})
		if err != nil {
			return err
		}
		myClaims := token.Claims.(*MyClaims)
		if !myClaims.Authorized {
			fmt.Errorf("not authorize")
			return c.JSON(http.StatusUnauthorized, http.NoBody)
		}
		if myClaims.Client != "matching-api" {
			fmt.Errorf("invalid client")
			return c.JSON(http.StatusUnauthorized, http.NoBody)
		}
		if token.Valid {
			next(c)
		}
		return nil
	}
}
