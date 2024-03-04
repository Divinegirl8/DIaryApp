package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Pin      string `json:"pin"`
}

type Entry struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	ID    string `json:"ID"`
	Time  string `json:"time"`
}

var userDatabase []User
var entryDataBase []Entry

func register(c *gin.Context) {

	var requestBody User

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := requestBody.Username
	pin := requestBody.Pin

	for _, users := range userDatabase {
		if username == users.Username {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username already exist"})
			return
		}
	}

	userDatabase = append(userDatabase, requestBody)

	c.JSON(http.StatusOK, gin.H{
		"message":  "registration successful",
		"username": username,
		"pin":      pin,
	})
}

var isLogin bool = false

func login(c *gin.Context) {
	var requestBody User

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	found := false

	for _, user := range userDatabase {
		if requestBody.Username == user.Username {
			found = true
			if requestBody.Pin != user.Pin {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "login credentials are not valid"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "login successful"})
			isLogin = true
			return
		}
	}

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user does not exist"})
	}
}

func logout(c *gin.Context) {
	var requestBody User

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	found := false

	for _, user := range userDatabase {
		if requestBody.Username == user.Username {
			found = true
			if requestBody.Pin != user.Pin {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "login credentials are not valid"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "logout successfully"})
			isLogin = false
			return
		}
	}
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user does not exist"})
	}
}

func addEntry(c *gin.Context) {
	var requestBody Entry

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if isLogin == false {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must login to perform this action"})
		return
	}

	title := requestBody.Title
	body := requestBody.Body
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	requestBody.Time = formattedTime

	entryDataBase = append(entryDataBase, requestBody)

	c.JSON(http.StatusOK, gin.H{
		"title": title,
		"body":  body,
		"time":  formattedTime,
		"ID":    len(entryDataBase),
	})

}

func findEntry(c *gin.Context) {
	entryID := c.Query("Entid")

	if isLogin == false {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must login to perform this action"})
		return
	}

	found := false
	var foundEntry Entry

	for _, entry := range entryDataBase {
		if entry.ID == entryID {
			found = true
			foundEntry = entry
			break
		}
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"title": foundEntry.Title,
			"body": foundEntry.Body,
			"time": foundEntry.Time})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Entry not found"})
	}
}

func deleteEntry(c *gin.Context) {
	entryID := c.Query("Entid")

	if isLogin == false {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must login to perform this action"})
		return
	}

	found := false

	for index, entry := range entryDataBase {
		if entry.ID == entryID {
			found = true
			entryDataBase = append(entryDataBase[:index], entryDataBase[index+1:]...)
			break
		}

	}

	if found {
		c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Entry not found"})
	}
}

func main() {

	fmt.Print("hello")
	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.POST("/entry", addEntry)
	router.GET("/findEntry", findEntry)
	router.DELETE("/deleteEntry", deleteEntry)

	err := router.Run("localhost:9080")
	if err != nil {
		return
	}
}
