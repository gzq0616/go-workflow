package go_workflow

import (
	"testing"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

func TestInitEngine(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/workflow?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("连接数据库失败！%s", err)
	}
	if err := EngineInit(engine); err != nil {
		t.Fatalf("数据库初始化失败！:%s", err)
	}
}
