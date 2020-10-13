package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const (
	USERNAME = "username"
)

func main() {
	h := newHub()
	go h.run()

	var addr string
	if len(os.Args) > 1 {
		addr = os.Args[1]
	} else {
		addr = ":8080"
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	r.LoadHTMLGlob("templates/*.html")

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", "")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", "")
	})

	r.POST("/chat", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(USERNAME, c.PostForm(USERNAME))

		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		c.HTML(http.StatusOK, "chat.html", "")
	})

	r.GET("/ws", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(USERNAME)
		switch v := user.(type) {
		case string:
			serveWs(h, c.Writer, c.Request, v)
		default:
			serveWs(h, c.Writer, c.Request, c.Request.RemoteAddr)
		}
	})

	r.Run(addr)
}
