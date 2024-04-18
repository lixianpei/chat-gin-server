package helper

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

// 设置雪花算法起始的时间，上线后不可轻易修改
const startTime = "2024-01-01"

// 生成客户ID的雪花节点
var clientIdSnowflakeNode *snowflake.Node

// InitAllSnowflakeNode 初始化所有的雪花节点
func InitAllSnowflakeNode() {
	//TODO 根据配置读取当前机器的ID
	//machineId := pkg.Config.Read("app").GetInt64("server.machine_id")
	machineId := int64(1)

	//根据雪花算法初始化客户ID的节点
	initClientIdBySnowflake(machineId)
}

func initSnowFlake(startTime string, machineId int64, dataCenter int64) (snowFlakeNode *snowflake.Node, err error) {
	var ts time.Time
	ts, err = time.Parse(time.DateOnly, startTime)
	if err != nil {
		return
	}
	if machineId < 0 || machineId > 9 {
		err = errors.New("InitSnowFlake fail：机器编号ID只能在[0-9]")
		return
	}
	if dataCenter < 0 || dataCenter > 99 {
		err = errors.New("InitSnowFlake fail：数据中心编号ID只能在[0-9]")
		return
	}
	//Node number must be between 0 and 1023：按机器+数据中心组合拼装node值，实现同一机器上拥有不同的node节点
	node := machineId*100 + dataCenter

	snowflake.Epoch = ts.UnixNano() / 1e6
	snowFlakeNode, err = snowflake.NewNode(node)
	if err != nil {
		return
	}
	return
}

// 根据雪花算法初始化客户ID的节点
func initClientIdBySnowflake(machineId int64) {
	//初始化客户ID生成的节点
	node, err := initSnowFlake(startTime, machineId, 1)
	if err != nil {
		fmt.Println(errors.New("InitOrderNoSnowflake fail：" + err.Error()))
	}
	clientIdSnowflakeNode = node
}

// GenerateClientId 根据雪花算法生成客户ID
func GenerateClientId() string {
	return clientIdSnowflakeNode.Generate().String()
}
