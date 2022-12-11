package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"dbmodel/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunp13/dbtool"
)

var (
	confPath string
	table    string
	outpath  string
)

func init() {
	flag.StringVar(&confPath, "c", "./conf/datasource_dev.yml", "datasource path")
	flag.StringVar(&table, "t", "", "table name")
	flag.StringVar(&outpath, "o", "./models", "output path")
	flag.Parse()

	if err := dbtool.Init(confPath); err != nil {
		panic(err)
	}
}

func main() {
	if table == "" {
		panic("table name is null")
	}
	spl := strings.Split(table, ".")
	if len(spl) < 2 {
		panic("table name error")
	}
	// tablePerfix := spl[0]
	tableName := spl[1]

	// 拿到表结构
	descData, err := descTable(table)
	if err != nil {
		panic(err)
	}

	// 主键
	priKey := ""
	fields := make([]string, 0)
	fieldsUpper := make([]string, 0)
	fieldsQuestion := make([]string, 0)
	fieldModify := make([]string, 0)
	for _, v := range descData {
		if fmt.Sprintf("%v", v["Key"]) == "PRI" {
			priKey = fmt.Sprintf("%v", v["Field"])
		} else {
			fields = append(fields, fmt.Sprintf("%v", v["Field"]))
			fieldsUpper = append(fieldsUpper, convField(fmt.Sprintf("%v", v["Field"])))
			fieldsQuestion = append(fieldsQuestion, "?")
			fieldModify = append(fieldModify, fmt.Sprintf("%v = ?", v["Field"]))
		}
	}

	args := &entity.Args{
		ModelName:       convField(tableName),
		TableName:       table,
		PrimaryKey:      priKey,
		PrimaryKeyUpper: convField(priKey),
		FieldsKey:       fields,
		FieldsUpper:     fieldsUpper,
		FieldsQuestion:  fieldsQuestion,
		FieldsModify:    fieldModify,
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.go", outpath, tableName))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	qTemp, err := template.New("normal.tmpl").Funcs(template.FuncMap{
		"JOIN": func(data []string, spec string) string {
			return strings.Join(data, spec)
		},
	}).ParseFiles("./models/normal.tmpl")
	if err != nil {
		panic(err)
	}
	err = qTemp.Execute(f, &args)
	if err != nil {
		panic(err)
	}

}

func descTable(tableName string) (res []map[string]interface{}, err error) {
	sql := fmt.Sprintf(`
	desc %s
	`, tableName)

	res, err = dbtool.D.QuerySQL(sql, nil)
	return
}

func convField(fname string) string {
	result := ""
	for i := 0; i < len(fname); i++ {
		if fname[i] == '_' {
			result += strings.ToUpper(string(fname[i+1]))
			i++
			continue
		} else {
			result += string(fname[i])
		}
	}
	return result
}
