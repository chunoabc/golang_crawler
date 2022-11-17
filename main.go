package main

import (
	"fmt"
	"homework/colly"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexData struct {
	Title   string
	Content string
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/html/*")
	//設定靜態資源的讀取
	server.Static("/assets", "./template/assets")
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/logout", logout)
	server.GET("/index", indexPage)
	server.GET("/index/getData", getData)
	server.Run(":8888")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func indexPage(c *gin.Context) {
	_, err := c.Cookie("site_cookie")
	if err != nil {
		fmt.Println("沒有cookie 回到登入頁面")
		c.Redirect(http.StatusFound, "/login")
		// c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	c.HTML(http.StatusOK, "index.html", nil)
}
func LoginAuth(c *gin.Context) {
	var (
		username string
		password string
	)
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "請輸入使用者名稱"})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "請輸入密碼"})
		return
	}
	if err := Auth(username, password); err == nil {
		c.SetCookie("site_cookie", "cookievalue", 600, "/", "/", false, true)
		c.JSON(http.StatusOK, gin.H{})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "帳號或密碼錯誤"})
		return
	}
}
func logout(c *gin.Context) {
	c.SetCookie("site_cookie", "cookievalue", -1, "/", "/", false, true)
	c.JSON(http.StatusOK, gin.H{})
	// c.Redirect(http.StatusMovedPermanently, "/login")
}

func getData(c *gin.Context) {
	name := c.Query("name")
	if name != "" {
		data := colly.CollyInit(name)
		// c.JSON(http.StatusOK, gin.H{"data": data})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data))
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "錯誤"})
	}
}
