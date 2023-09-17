package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//router.LoadHTMLFiles("templates/index.tmpl")
	//router.LoadHTMLGlob("templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "gin-template-title",
			"path":  dir,
		})

	})
	MultiFile(router)
	router.Run(":8083")
}

func MultiFile(r *gin.Engine) {
	r.LoadHTMLGlob("./templates/**/*")
	r.GET("/user/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.list.html", nil)
	})
	r.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.list.html", nil)
	})
}
