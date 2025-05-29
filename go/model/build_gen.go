package model

import (
	"goproject/model/orm/dal"
	//"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
)

// TODO dsn处理host
func dsn() string {
	return "host=db user=postgres password=123456 dbname=postgres" +
		" port=5432 sslmode=disable TimeZone=Asia/Shanghai search_path=game"
}

func BuildGen() *gorm.DB {
	dsn := dsn()
	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("[%s] postgresql 连接失败", dsn)
	}

	// gen 配置
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./model/orm/dal",
		ModelPkgPath: "./model/orm/model",
		Mode:         gen.WithDefaultQuery | gen.WithoutContext,
	})
	g.UseDB(db)
	// 生成所有表的DAO接口 配置文件修改
	// TODO 控制这个变量决定是否生成模型
	isCreate := false
	if isCreate {
		g.ApplyBasic(g.GenerateAllTable()...)
	}
	g.Execute()
	//log.Println("<-----------gen success----------->")
	//logrus.Info("gen init success")
	dal.SetDefault(db)
	return db
}
