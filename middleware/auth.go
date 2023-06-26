package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		kuki, err := ctx.Cookie("session_token")
		if err != nil {
			contentType := ctx.Request.Header.Get("Content-Type")
			if contentType != "application/json" {
				ctx.Redirect(http.StatusSeeOther, "/user/login")
				return
			} else {
				ctx.JSON(http.StatusUnauthorized, "cookie not found")
				return
			}
		}

		claims := &model.Claims{}

		tkn, err := jwt.ParseWithClaims(kuki, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, "token gagal")
				return
			}
			ctx.JSON(http.StatusBadRequest, "token error")
			return
		}

		if !tkn.Valid {
			ctx.JSON(http.StatusUnauthorized, "token tidak valid")
			return
		}

		ctx.Set("email", claims.Email)

		ctx.Next()
	})
}
