package main

import "C"
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Learning about middlewares in go

func myLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request has arrive")

		c.Next()

		fmt.Println("Request has finished")

	}

}

//creating a middle to check if the user is admin

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check the route params for the name
		name := c.Param("name")
		if name == "mubarak" {
			c.Set("role", "admin")
			c.Next()
		} else {
			c.JSON(200, gin.H{
				"message": "access denied",
				"name":    name,
			})
			c.Abort()
		}

	}

}

func main() {

	//initial gin using the default
	r := gin.Default()

	r.GET("/ping", myLogger(), func(c *gin.Context) {
		//using the helper to send json response
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//creating user route that accept the name as the param

	r.GET("/user/:name", func(c *gin.Context) {
		//get the name from the params

		name := c.Param("name")

		c.JSON(200, gin.H{
			"Message": "hello " + name,
		})
	})

	//trying a querry param route in go

	r.GET("/search", func(c *gin.Context) {

		search := c.DefaultQuery("topic", "everything")
		c.String(200, "search for ", search)

	})

	//Trying a post request

	//first create a type to represent the data you want to represent

	type User struct {
		// These MUST start with Capital Letters to be "Exported" (Public)
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Your route remains almost the same, but now it can "see" the fields
	r.POST("/create", func(c *gin.Context) {
		var newUser User

		// Gin can now fill the CAPITALIZED fields using the lowercase json tags
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "received",
			"data":   newUser,
		})
	})

	//creating an admin route

	r.GET("/userss/:name", IsAdmin(), func(c *gin.Context) {
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"status": "authorized",
			"name":   c.Param("name"),
			"role":   role,
		})
	})
	r.Run()

}
