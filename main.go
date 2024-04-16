package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"dbmodel/entity"
	"dbmodel/models"
	"dbmodel/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunp13/dbtool"
)

var (
	confPath       string
	table          string
	foreignKey     string
	withController bool
	withEntity     bool
	withModel      bool
)

func init() {
	flag.StringVar(&confPath, "c", "./conf/datasource_dev.yml", "datasource path")
	flag.StringVar(&table, "t", "", "table name")
	flag.StringVar(&foreignKey, "f", "", "foreign key")
	flag.BoolVar(&withController, "controller", false, "generate controllers")
	flag.BoolVar(&withEntity, "entity", false, "generate entity")
	flag.BoolVar(&withModel, "model", false, "generate model")
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
	// 库名 = spl[0]
	schema := spl[0]
	// 表名 = spl[1]
	table := spl[1]

	// 拿到表结构
	descData, err := models.DescModel.Desc(table)
	if err != nil {
		panic(err)
	}
	// 主键
	priKey := ""
	fields := make([]string, 0)
	fieldsHump := make([]string, 0)
	fieldsHumpUpper := make([]string, 0)
	fieldsQuestion := make([]string, 0)
	fieldModify := make([]string, 0)
	for _, v := range descData {
		if fmt.Sprintf("%v", v["Key"]) == "PRI" {
			priKey = fmt.Sprintf("%v", v["Field"])
		} else {
			field := fmt.Sprintf("%v", v["Field"])
			if field == "add_time" || field == "modify_time" || field == "is_deleted" {
				continue
			}
			fields = append(fields, fmt.Sprintf("%v", v["Field"]))
			fieldsHump = append(fieldsHump, utils.ConvField(fmt.Sprintf("%v", v["Field"])))
			fieldsHumpUpper = append(fieldsHumpUpper, utils.ConvFieldFirstUpper(fmt.Sprintf("%v", v["Field"])))
			fieldsQuestion = append(fieldsQuestion, "?")
			fieldModify = append(fieldModify, fmt.Sprintf("%v = ?", v["Field"]))
		}
	}

	args := &entity.Args{
		// 下划线 库名
		SchemaName: schema,
		// 下划线 表名
		TableName: table,
		// 驼峰单独table
		TableNameHump: utils.ConvField(table),
		// 首字母大写+驼峰
		TableNameHumpUpper: utils.ConvFieldFirstUpper(table),

		// 下划线 主键
		PrimaryKey: priKey,
		// 驼峰 主键
		PrimaryKeyHump: utils.ConvField(priKey),
		// 驼峰+大写
		PrimaryKeyHumpUpper: utils.ConvFieldFirstUpper(priKey),
		// 下划线 外键
		ForeignKey: foreignKey,
		// 驼峰 外键
		ForeignKeyHump: utils.ConvField(foreignKey),

		// 下划线 所有字段 (不含时间和isdelete)
		FieldsKey: fields,
		// 驼峰 所有字段 (不含时间和isdelete)
		FieldsKeyHump: fieldsHump,
		// 驼峰+大写
		FieldsKeyHumpUpper: fieldsHumpUpper,

		// 下划线 换成问好的字段
		FieldsQuestion: fieldsQuestion,
		// 下划线 更新用的所有键
		FieldsModify: fieldModify,
	}

	// 如果需要创建controller
	if withController {
		if err := generateController(args); err != nil {
			panic(err)
		}
	}

	// 如果需要创建entity
	if withEntity {
		if err := generateEntity(args); err != nil {
			panic(err)
		}
	}

	// 如果是创建模型
	if withModel {
		if err := generateModel(args); err != nil {
			panic(err)
		}
	}

}

func generateModel(args *entity.Args) error {
	f, err := os.Create(fmt.Sprintf("./models/%s.go", args.TableName))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	qTemp, err := template.New("normal.tmpl").Funcs(template.FuncMap{
		"JOIN": func(data []string, spec string) string {
			return strings.Join(data, spec)
		},
		"EQ": func(target, dest string) bool {
			return target == dest
		},
		"NEQ": func(target, dest string) bool {
			return target != dest
		},
	}).ParseFiles("./models/normal.tmpl")
	if err != nil {
		panic(err)
	}
	err = qTemp.Execute(f, &args)
	if err != nil {
		panic(err)
	}
	return nil
}

// 生成 controller
func generateController(args *entity.Args) error {

	f, err := os.Create(fmt.Sprintf("./controllers/%s.go", args.TableName))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	qTemp, err := template.New("normal.tmpl").Funcs(template.FuncMap{
		"JOIN": func(data []string, spec string) string {
			return strings.Join(data, spec)
		},
		"EQ": func(target, dest string) bool {
			return target == dest
		},
		"NEQ": func(target, dest string) bool {
			return target != dest
		},
	}).ParseFiles("./controllers/normal.tmpl")

	if err != nil {
		return err
	}
	err = qTemp.Execute(f, &args)
	if err != nil {
		return err
	}
	return nil
}

func generateEntity(args *entity.Args) error {
	f, err := os.Create(fmt.Sprintf("./entity/%s.go", args.TableName))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	qTemp, err := template.New("normal.tmpl").Funcs(template.FuncMap{
		"ENTJSON": func(data []string, spec string) string {
			nData := make([]string, 0)

			for _, v := range data {
				s := fmt.Sprintf("%s string `json:\"%s\"`", utils.ConvFieldFirstUpper(v), v)
				nData = append(nData, s)
			}

			return strings.Join(nData, spec)
		},
	}).ParseFiles("./entity/normal.tmpl")

	if err != nil {
		return err
	}
	err = qTemp.Execute(f, &args)
	if err != nil {
		return err
	}
	return nil

}
