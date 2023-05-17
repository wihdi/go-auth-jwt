package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("SECRET_KEY_BEBAS")

type User struct {
	ID int `json:id`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	
	//gin router
	r := gin.Default()
	//setuo router
	r.POST("/auth/login",loginHandler)
	userRouter:=r.Group("api/vi/users")
	//midlerware
	userRouter.Use(authMiddleware())

	//set up GET
	userRouter.GET("/:id/profile", profileHanlder)
	//start server
	r.Run(":8080")

}

func loginHandler(c *gin.Context){
var user User

if err:= c.ShouldBind(&user);err!=nil{

c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
//logic auth(compare username and paswarod)

if user.Username =="enigma"&& user.Password =="12345"{

token:=jwt.New(jwt.SigningMethodHS256)
claims := token.Claims.(jwt.MapClaims)
claims["username"]=user.Username

claims["exp"]=time.Now().Add(time.Minute * 1).Unix()
tokenString,err :=token.SignedString(jwtKey)
if err != nil {

	c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to generate token"})
	return
}
c.JSON(http.StatusOK, gin.H{"token" : tokenString})
}else {
	c.JSON(http.StatusUnauthorized,gin.H{"error": "Invalid credential "})
}

}

func profileHanlder(c *gin.Context){

// ambil username dari jwt token

claims := c.MustGet("claim").(jwt.MapClaims)

username := claims ["username"].(string)
//respons user database 
c.JSON(http.StatusOK, gin.H {"username" :username})

}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)
		c.Next()
	}
}
