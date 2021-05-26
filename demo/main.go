package main

import (
	"junkoho/gin_jwt/jwt"
	"log"
)

func main() {

	//测试
	userName := "admin"
	token, err := jwt.GenToken(userName)
	if err != nil {
		log.Println("生成token失败" + err.Error())
	}

	log.Println("生成的token：" + token)

	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		log.Println("解析token失败" + err.Error())
	}

	log.Println("myClaims")
	log.Printf("%#v", myClaims)
	log.Println("解析后的token：" + myClaims.Username)
}
