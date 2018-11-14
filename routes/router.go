package routes

import (
	"fmt"
	"net/http"

	"github.com/siskonemilia/CloudGo-IO/model"

	"github.com/gin-gonic/gin"
)

type registerJSON struct {
	Username string `form:"username" binding:"required"`
	Stuid    string `form:"stuid" binding:"required"`
	Tel      string `form:"tel" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

type test struct {
	Events string
}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func register(c *gin.Context) {

	var user registerJSON
	flagUsername, flagStuid, flagEmail, flagPhone := true, true, true, true

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": "Something wrong with the server",
		})
		fmt.Println(err.Error())
		return
	}

	flagUsername = flagUsername && model.CheckUsername(user.Username)
	flagStuid = flagStuid && model.CheckStuID(user.Stuid)
	flagEmail = flagEmail && model.CheckEmail(user.Email)
	flagPhone = flagPhone && model.CheckPhone(user.Tel)

	if flagUsername && flagStuid && flagEmail && flagPhone {
		model.AddUser(user.Username, user.Stuid, user.Email, user.Tel)
	}
	c.JSON(http.StatusOK, gin.H{
		"username": flagUsername,
		"stuid":    flagStuid,
		"tel":      flagPhone,
		"email":    flagEmail,
	})
}

func detail(c *gin.Context) {
	username := c.Query("username")
	user, err := model.FetchInfo(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "No such user.",
		})
		return
	}
	c.HTML(http.StatusOK, "Detail.tmpl", gin.H{
		"username": user.Username,
		"stuid":    user.Stuid,
		"tel":      user.Phone,
		"email":    user.Email,
	})
}

// Router method defines the routing behaviors
// And generate a router for others to use
func Router() *gin.Engine {
	// Generate a default router to configure
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	// Set the handler function for paths
	router.GET("/ginTest", func(c *gin.Context) {
		// Use JSON as the response
		c.JSON(http.StatusOK, gin.H{
			"message": "You've successfully received a message from a gin server.",
		})
	})

	router.GET("/", homePage)
	router.GET("/detail", detail)
	router.POST("/sign_req", register)
	router.Static("/public", "public")

	// NOT_IMPLEMENTED for /unknown
	router.GET("/unknown", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "This page has not been implemented.",
		})
	})

	// PAGE_NOT_FOUND page for all paths without routing
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Target page not found.",
		})
	})

	// Configuration done, router returned
	return router
}
