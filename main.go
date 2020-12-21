package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Secret struct {
	gorm.Model
	UniqueID string
	Secret   string
}

func main() {
	db, err := gorm.Open(sqlite.Open("secrets.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Secret{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", nil)
	})
	r.POST("/create", func(c *gin.Context) {
		secret := c.PostForm("the_secret")
		ok := false
		uuidString := ""
		for i := 0; i < 100; i++ {
			uuidString = uuid.New().String()
			result := Secret{}
			if err := db.First(&result, "unique_id = ?", uuidString).Error; err == gorm.ErrRecordNotFound {
				ok = true
				break
			}
		}
		if !ok {
			c.HTML(http.StatusInternalServerError, "secret_created.html", nil)
			return
		}
		// return link to url in output template
		db.Create(&Secret{UniqueID: uuidString, Secret: secret})
		c.HTML(http.StatusOK, "secret_created.html", gin.H{
			"secret_url": "/view/" + uuidString,
		})
	})

	r.GET("/view/:uniqueid", func(c *gin.Context) {
		uniqueID := c.Param("uniqueid")
		result := Secret{}
		if err := db.First(&result, "unique_id = ?", uniqueID).Error; err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		theSecret := result.Secret
		db.Delete(&result)
		c.HTML(http.StatusOK, "view.html", gin.H{
			"the_secret": theSecret,
		})
	})
	r.Run("0.0.0.0:8000")
}
