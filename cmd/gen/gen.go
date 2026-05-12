package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// generate code
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	//db, _ := gorm.Open(mysql.Open("golang:123456@(10.121.120.114:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"))
	db, _ := gorm.Open(mysql.Open("golang:123456@(10.211.55.3:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db)

	// generate a struct and specify struct's name
	g.ApplyBasic(
		g.GenerateModel("user"),
		g.GenerateModel("session"),
		g.GenerateModel("community"),
		g.GenerateModel("post"),
	)

	// execute the action of code generation
	g.Execute()
}
