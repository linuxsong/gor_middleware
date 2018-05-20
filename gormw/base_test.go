package gormw

import (
	"encoding/hex"
	"testing"
)

var passby map[string]int = make(map[string]int)

func init() {
	passby["counter"] = 0
}

func incrCounter(gor *Gor, msg *GorMessage, kwargs ...interface{}) *GorMessage {
	passbyReadonly, _ := kwargs[0].(map[string]int)
	increase, _ := kwargs[1].(int)
	passby["counter"] += passbyReadonly["counter"] + increase
	return msg
}

func TestMessageLogic(t *testing.T) {
	gor := CreateGor()
	gor.On("message", incrCounter, "", &passby, 1)
	gor.On("request", incrCounter, "", &passby, 2)
	gor.On("response", incrCounter, "2", &passby, 3)
	if len(gor.retainQueue) != 2 {
		t.Errorf("gor retain queue length %d != 2", len(gor.retainQueue))
	}
	if len(gor.tempQueue) != 1 {
		t.Errorf("gor temp queue length %d != 1", len(gor.tempQueue))
	}
	req, err := gor.ParseMessage(hex.EncodeToString([]byte("1 2 3\nGET / HTTP/1.1\r\n\r\n")))
	if err != nil {
		t.Error(err.Error())
	}
	resp, err := gor.ParseMessage(hex.EncodeToString([]byte("2 2 3\nHTTP/1.1 200 OK\r\n\r\n")))
	if err != nil {
		t.Error(err.Error())
	}
	resp2, err := gor.ParseMessage(hex.EncodeToString([]byte("2 3 3\nHTTP/1.1 200 OK\r\n\r\n")))
	if err != nil {
		t.Error(err.Error())
	}
	gor.Emit(req)
	gor.Emit(resp)
	gor.Emit(resp2)
	if passby["counter"] != 8 {
		t.Errorf("passby counter %d != 8", passby["counter"])
	}
}
