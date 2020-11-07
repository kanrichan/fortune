package core

import "C"

//export GO_Create
func GO_Create(version *C.char) *C.char {
	result := "{\"name\":\"fortune-运势\", \"pver\":\"1.0.5\", \"sver\":3, \"author\":\"Kanri\", \"desc\":\"项目地址 https://github.com/Yiwen-Chan/fortune\"}"
	return CString(result)
}

//export GO_Event
func GO_Event(botQQ *C.char, msgType C.int, subType C.int, sourceId *C.char, activeQQ *C.char, passiveQQ *C.char, msg *C.char, msgNum *C.char, msgId *C.char, rawMsg *C.char, timeStamp *C.char, retText *C.char) C.int {
	if int(msgType) == 10000 {
		Create()
	} else {
		Event(botQQ, msgType, subType, sourceId, activeQQ, passiveQQ, msg, msgNum, msgId, rawMsg, timeStamp, retText)
	}
	return 0
}

//export GO_DestroyPlugin
func GO_DestroyPlugin() C.int {
	return 0
}

//export GO_SetUp
func GO_SetUp() C.int {
	return 0
}

func main() {
	//
}
