/*
	网易云音乐请求参数加密,来自https://github.com/yitimo/api-163-go/blob/master/madoka/encrypt.go，感谢
*/

package api

import (
	"MyCloudMusic_Server_Go/mylog"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
)

const MODULUS = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
const NONCE = "0CoJUm6Qyw8W8jud"
const PUBKEY = "010001"
const KEYS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"
const IV = "0102030405060708"

//
func encryptoParams(param map[string]interface{}) (string, string, error) {
	//创建key
	bytesData, err := json.Marshal(param)
	if err != nil {
		mylog.Error(err.Error())
	}

	secKey := createSecretKey(16)
	aes1, err1 := aesEncrypt(bytesData, NONCE)
	//第一次加密，使用固定nonce
	if err1 != nil {
		return "", "", err1
	}

	aes2, err2 := aesEncrypt([]byte(aes1), secKey)
	//第二次加密，使用创建的nonce
	if err2 != nil {
		return "", "", err2
	}

	//得到加密好的param以及key
	return aes2, rsaEncrypt(secKey, PUBKEY, MODULUS), nil
}

//创建指定长度的key
func createSecretKey(size int) string {
	//也就是从a-9以及+/中随机拿出指定数量的字符拼成一个key
	rs := ""
	for i := 0; i < size; i++ {
		pos := rand.Intn(len(KEYS))
		rs += KEYS[pos : pos+1]
	}
	return rs
}

//通过CBC模式的AES加密用sKey加密sSrc
func aesEncrypt(sSrc []byte, sKey string) (string, error) {
	iv := []byte(IV)
	block, err := aes.NewCipher([]byte(sKey))
	if err != nil {
		return "", nil
	}

	padding := block.BlockSize() - len(sSrc)%block.BlockSize()
	src := append(sSrc, bytes.Repeat([]byte{byte(padding)}, padding)...)
	model := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(src))
	model.CryptBlocks(cipherText, src)
	//最后使用base64编码输出
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

//将key也加密
func rsaEncrypt(key string, pubKey string, modulus string) string {
	//倒序key
	rKey := ""
	for i := len(key) - 1; i >= 0; i-- {
		rKey += key[i : i+1]
	}
	//将key转ascii编码，然后转成16进制字符串
	hexRKey := ""
	for _, char := range []rune(rKey) {
		hexRKey += fmt.Sprintf("%x", int(char))
	}
	//将16进制的三个参数转为10进制的bigint
	bigRKey, _ := big.NewInt(0).SetString(hexRKey, 16)
	bigPubKey, _ := big.NewInt(0).SetString(pubKey, 16)
	bigModulus, _ := big.NewInt(0).SetString(modulus, 16)
	//执行幂乘取模运算得到最终的bigint结果
	bigRs := bigRKey.Exp(bigRKey, bigPubKey, bigModulus)
	//将结果转为16进制字符串
	hexRs := fmt.Sprintf("%x", bigRs)
	//可能在不满256位的情况，要在前面补0补满256位
	return addPadding(hexRs, modulus)
}

//补0步骤
func addPadding(encText string, modulus string) string {
	ml := len(modulus)
	for i := 0; ml > 0 && modulus[i:i+1] == "0"; i++ {
		ml--
	}
	num := ml - len(encText)
	prefix := ""
	for i := 0; i < num; i++ {
		prefix += "0"
	}

	return prefix + encText
}
