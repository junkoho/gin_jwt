package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"junkoho/gin_jwt/jwt"
	"log"
	"net/http"
)

type UserInfo struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中
		// 这里的具体实现方式要依据你的实际业务情况决定
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中token为空",
			})
			c.Abort()
			return
		}
		log.Println("token")
		log.Println(token)

		//tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()
	// GET：请求方式；/auth：请求的路径
	// 当客户端以GET方法请求/auth路径时，会执行后面的匿名函数
	r.POST("/auth", func(c *gin.Context) {

		// 用户发送用户名和密码过来
		var user UserInfo

		//gin获取参数的五种方法
		//1获取querystring参数
		//user.Username = c.Query("username")
		//user.Password = c.Query("password")

		//2获取form参数
		//user.Username = c.PostForm("username")
		//user.Password = c.PostForm("password")

		//3获取json参数
		body, _ := ioutil.ReadAll(c.Request.Body)
		log.Printf(string(body))
		if err := json.Unmarshal(body, &user); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "参数解析报错:" + err.Error(),
			})
			return
		}

		//4获取path参数
		//user.Username = c.Param("username")
		//user.Password = c.Param("password")

		//5参数绑定
		//err := c.ShouldBind(&user)
		//if err != nil {
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": 2001,
		//		"msg":  "无效的参数",
		//	})
		//	return
		//}

		// 校验用户名和密码是否正确
		if user.Username == "junko" && user.Password == "junko123" {
			// 生成Token
			tokenString, _ := jwt.GenToken(user.Username)
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"msg":  "success",
				"data": gin.H{"token": tokenString},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 2002,
			"msg":  "鉴权失败",
		})
		return

	})

	r.GET("/home", JWTAuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"username": username},
		})
	})

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	http.ListenAndServe(":8081", r)
	r.Run()

}
