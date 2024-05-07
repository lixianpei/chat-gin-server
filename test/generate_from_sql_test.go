package test

import (
	"GoChatServer/helper"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
	"testing"
)

// 通过读取sql文件生成相关数据表结构体
func TestSqlGen(t *testing.T) {
	dbName := "chat"
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../dal/query/" + dbName + "_query",
		ModelPkgPath: "../model/" + dbName + "_model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		WithUnitTest: false,
	})
	g.WithJSONTagNameStrategy(func(c string) string {
		return helper.UnderscoreCamelCase(c) //指定生成的model对应数据表字段json格式化为驼峰
	})
	gormDb, _ := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{
			fmt.Sprintf("./sql/%s.sql", dbName),
		},
	}))
	g.UseDB(gormDb) // reuse your gorm db
	g.ApplyBasic(
		g.GenerateAllTable()...,
	)
	g.Execute()
}

func TestS(t *testing.T) {
	str := "abc_Def_ggg"
	fmt.Println(helper.UnderscoreCamelCase(str))
}
