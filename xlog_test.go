// go test
// go test -v -run TestLog2
package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	L, _ := NewLogger("test.log", true, "dbg")
	L.Error("this is err")
	L.Warn("this is warn")
	L.Info("this is info")
	L.Debug("this is debug")
	// L.InfoIf(true, "this is InfoIf")

	// ChangeLevel("err")
	L.Error("this is err again")
	L.Warn("this is warn again")
	L.Info("this is info again")
	L.Debug("this is debug again")
	L.Info("this is InfoIf again")

	logger := L.With("key1", "val1", "key2", "val2")
	logger.Errorf("logger err")
	logger2 := L.With("key3", "val3")
	logger2.Errorf("logger2 err")
	L.Error("err after logger")
}

/*
// 测试文件转存
func TestLog2(t *testing.T) {
	MylogInit("/tmp/mylog_test.log", false, "dbg")

	s := ""
	for i := 0; i < 256; i++ {
		s += "a"
	}
	for i := 0; i < 50000; i++ {
		Error(s)
	}
}
*/
