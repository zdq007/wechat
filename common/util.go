package common

import (
	"math/rand"
	"time"
)

//生成随机字符串
func GetRandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytearr := []byte(str)
	length := len(bytearr)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytearr[r.Intn(length)])
	}
	return string(result)
}

//生成随机数字串
func GetRandomStringNum(size int) string {
	str := "0123456789"
	bytearr := []byte(str)
	length := len(bytearr)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytearr[r.Intn(length)])
	}
	return string(result)
}
