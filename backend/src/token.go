package main

//解码Token

import (
	"encoding/base64"
	"fmt"
)

func main() {

	//token:="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyb290IiwiZXhwIjo3MjAwLCJpc3MiOiJCYWNrZW5kIn0.sJqcHanOdpqvKHGm1c-YJzChf6odbV9Tnjl7qouJGzo"
	header := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	payload := "eyJhdWQiOiJyb290IiwiZXhwIjo3MjAwLCJpc3MiOiJCYWNrZW5kIn0"
	signature := "sJqcHanOdpqvKHGm1c-YJzChf6odbV9Tnjl7qouJGzo"

	//解码出，包含签名算法和JWT标识
	decodedHeader, err := base64.URLEncoding.DecodeString(header)
	if err != nil {
		fmt.Println("decodedHeader error:", err)
		return
	}
	fmt.Println("header:", string(decodedHeader))

	//解码出，包含核心数据
	decodedPayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("decodedPayload error:", err)
		return
	}
	fmt.Println("payload:", string(decodedPayload))

	//解码出，包含HS256基于header和payload还有key生成的摘要
	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("decodedSignature error:", err)
		return
	}
	fmt.Println("signature:", string(decodedSignature))

}
