package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
)

type Message struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type Log struct {
	ErrorString string `json:"error_string"`
	Operation   string `json:"operation"`
	Stack       string `json:"stack"`
}

type ApiResp struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

func Initialize() {
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			// the maximum allowable packet length is 128 bytes，use \n to split, strip delimiter
			AddLast(frame.DelimiterCodec(128*1024*1024, "\n", true)).
			// convert to string
			AddLast(format.TextCodec()).
			// LoggerHandler, print connected/disconnected event and received messages
			AddLast(MessageHandler{})
	}

	// create bootstrap & listening & accepting
	netty.NewBootstrap(netty.WithChildInitializer(childInitializer)).
		Listen(":8888").Sync()
}

type MessageHandler struct {
}

func (MessageHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	msgStr, ok := message.(string)
	if !ok {
		return
	}
	// base 64 解码
	msgRsaBytes, err := base64.StdEncoding.DecodeString(msgStr)
	if err != nil {
		fmt.Println("Base64 decode error:", err)
		return
	}
	// rsa解码
	rsaDecoded, err := RsaDecrypt(msgRsaBytes)
	if err != nil {
		fmt.Println("Rsa decode error:", err)
		return
	}
	msgInter := &Message{}
	err = json.Unmarshal(rsaDecoded, msgInter)
	if err != nil {
		fmt.Println("Json decode error:", err)
		return
	}
	fmt.Println(msgInter)
	if msgInter.Type == "api_key" {
		resp := &ApiResp{
			Key:   "key",
			Count: 33,
		}
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Json encode err:", err)
			return
		}
		resultByte, err := EncodeForRsaBase64String([]byte(jsonStr))
		if err != nil {
			fmt.Println("Encode for resp byte error:", err)
			return
		}
		fmt.Println(len(([]byte(resultByte))))
		ctx.Write(resultByte)
	}
	ctx.HandleRead(message)
}
