package test

import (
	"GoChatServer/helper"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化配置
	helper.InitConfig("../config")

	fmt.Println("TestMain...")
	os.Exit(m.Run())
}

func TestConfig(t *testing.T) {
}
