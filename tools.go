package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
)

var (
	db *gorm.DB
)

func ConnectDB() *gorm.DB {
	var (
		err error
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		"yz_select_xuyunlong", "Xylks3&19$h1dKtR^$+", "rm-bp15401h920m1vr5s.mysql.rds.aliyuncs.com", 3306, "qy_wx")

	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}

	return db
}

func GetDB() *gorm.DB {
	return db
}

func GenerateTableStruct(db *gorm.DB) {
	//根据配置实例化 gen
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./dao",   //curd代码的输出路径
		ModelPkgPath:      "./model", //model代码的输出路径
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})
	//使用数据库的实例
	g.UseDB(db)
	//模型结构体的命名规则
	g.WithModelNameStrategy(func(tableName string) (modelName string) {
		if strings.HasPrefix(tableName, "tbl") {
			return firstUpper(tableName[3:])
		}
		if strings.HasPrefix(tableName, "table") {
			return firstUpper(tableName[5:])
		}
		return firstUpper(tableName)
	})
	//模型文件的命名规则
	g.WithFileNameStrategy(func(tableName string) (fileName string) {
		if strings.HasPrefix(tableName, "tbl") {
			return firstLower(tableName[3:])
		}
		if strings.HasPrefix(tableName, "table") {
			return firstLower(tableName[5:])
		}
		return tableName
	})
	//数据类型映射
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"int": func(columnType gorm.ColumnType) (dataType string) {
			//if strings.Contains(detailType, "unsigned") {
			//    return "uint64"
			//}
			return "int64"
		},
		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			//if strings.Contains(detailType, "unsigned") {
			//    return "uint64"
			//}
			return "int64"
		},
	}
	//使用上面的类型映射
	g.WithDataTypeMap(dataMap)
	//生成数据库内所有表的结构体
	//g.GenerateAllTable()
	//生成某张表的结构体
	//g.GenerateModelAs("tblUser", "User")
	g.ApplyBasic(g.GenerateModel("chat_data"))
	//执行
	g.Execute()
}

// 字符串第一位改成小写
func firstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// 字符串第一位改成大写
func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
