package service

import (
	"fmt"
	"testing"
)

func TestRsaService_EncodeDecode(t *testing.T) {
	data, err := RsaEncrypt([]byte("Test 123131414"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	decoded, err := RsaDecrypt(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(decoded))
}
