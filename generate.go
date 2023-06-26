package go_mongo

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

type FieldInfo struct {
	Name     string
	JSONName string
	Type     string
}

func GetStructsFromGoFile(fileName string) {
	structs, fileNameList := getStructsFromFile(fileName)
	for i, s := range structs {
		fields := getFieldsFromStruct(s)
		constants, err := generateConstants(fileNameList[i], fields)
		if err != nil {
			fmt.Println("Error generateConstants:", err)
			os.Exit(1)
		}
		filename := strings.ToLower(fileNameList[i]) + ".go"
		err = createDirectory(fileNameList[i])
		if err != nil {
			fmt.Println("Error creating directory:", err)
			os.Exit(1)
		}
		err = writeConstantsToFile(strings.ToLower(fileNameList[i])+"/"+filename, constants)
		if err != nil {
			fmt.Println("Error writing constants to file:", err)
			os.Exit(1)
		}
		createQuery(fileNameList[i])

	}
	createClient(fileNameList)
	createConfig()
	//createOperator()

}

func createClient(fileNameList []string) {
	err := writeConstantsToFile("client.go", generateClientCode(fileNameList))
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createConfig() {
	err := writeConstantsToFile("config.go", generateConfig())
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createQuery(structName string) {
	err := writeConstantsToFile(fmt.Sprintf("%s_query.go", strings.ToLower(structName)), generateQuery(structName))
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createDirectory(dirName string) error {
	dirName = strings.ToLower(dirName)
	err := os.Mkdir(strings.ToLower(dirName), os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

//func createOperator() {
//	err := writeConstantsToFile("operator.go", generateOperator())
//	if err != nil {
//		fmt.Println("Error writing constants to file:", err)
//		os.Exit(1)
//	}
//
//}

//	func generateOperator() []byte {
//		builder := strings.Builder{}
//		builder.WriteString("package go_mongo\n\n")
//
//		builder.WriteString("var ComparisonOperators = map[string][]string{\n")
//
//		// 添加字段类型和比较操作符的映射
//		addMapping(&builder, "int", []string{"$eq", "$ne", "$gt", "$lt", "$gte", "$lte"})
//		addMapping(&builder, "string", []string{"$eq", "$ne", "$regex"})
//		addMapping(&builder, "bool", []string{"$eq", "$ne"})
//		// 添加其他字段类型和相应的比较操作符...
//
//		builder.WriteString("}\n")
//
//		return []byte(builder.String())
//
// }
func addMapping(builder *strings.Builder, fieldType string, operators []string) {
	builder.WriteString(fmt.Sprintf("\t\"%s\": []string{", fieldType))
	for i, op := range operators {
		builder.WriteString(fmt.Sprintf("\"%s\"", op))
		if i < len(operators)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("},\n")
}

func getStructsFromFile(filename string) ([]*ast.StructType, []string) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, filename, nil, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	structs := make([]*ast.StructType, 0)
	structNameList := make([]string, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.TypeSpec:
			if s, ok := x.Type.(*ast.StructType); ok {
				structNameList = append(structNameList, x.Name.Name)
				structs = append(structs, s)
			}
		}
		return true
	})

	return structs, structNameList
}

func getFieldsFromStruct(s *ast.StructType) []FieldInfo {
	fields := make([]FieldInfo, 0)
	for _, field := range s.Fields.List {
		if field.Names == nil {
			continue
		}
		jsonName := field.Names[0].Name
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			jsonName = tag.Get("bson")
		}
		if field.Type == nil {
			continue
		}
		fields = append(fields, FieldInfo{
			Name:     field.Names[0].Name,
			JSONName: jsonName,
			Type:     fmt.Sprintf("%s", field.Type),
		})
	}
	return fields
}

func generateConstants(fileName string, fields []FieldInfo) ([]byte, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(fileName)))
	builder.WriteString("import \"go.mongodb.org/mongo-driver/bson\"\n")
	builder.WriteString("const (")
	builder.WriteString(fmt.Sprintf("%sMongo = \"%s\"\n", fileName, ToSnakeCase(fileName)))
	for _, field := range fields {
		builder.WriteString(fmt.Sprintf("%sField = \"%s\"\n", field.Name, field.JSONName))
	}
	builder.WriteString(")\n")
	builder.WriteString(generateFindFunction(fileName, fields))
	builder.WriteString(generatePredicate(fileName))

	src, err := format.Source([]byte(builder.String()))
	if err != nil {
		return []byte{}, err
	}
	return src, nil
}

func generatePredicate(typeName string) string {
	return fmt.Sprintf(
		"type %sPredicate func(*bson.D)", typeName)
}

func generateQuery(structName string) []byte {
	clientCode := "package go_mongo\n\n"

	clientCode += fmt.Sprintf("import (\n")
	clientCode += fmt.Sprintf("\t\"context\"\n")
	clientCode += fmt.Sprintf("\t\"cc/go-mongo/%s\"\n", strings.ToLower(structName))
	clientCode += fmt.Sprintf("\t\"go.mongodb.org/mongo-driver/bson\"\n")
	clientCode += fmt.Sprintf(")\n")

	clientCode += fmt.Sprintf("type %sQuery struct {\n", structName)
	clientCode += fmt.Sprintf("\tconfig\n")
	clientCode += fmt.Sprintf("\tPredicates []%s.%sPredicate\n", strings.ToLower(structName), structName)
	clientCode += fmt.Sprintf("\tdbName string\n\n")
	clientCode += fmt.Sprintf("}\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) Where(ps ...%s.%sPredicate)*%sQuery{\n", structName, strings.ToLower(structName), structName, structName)
	//clientCode += fmt.Sprintf("\tfilter := bson.D{}\n")
	clientCode += fmt.Sprintf("\tfor _, p := range ps {\n")
	clientCode += fmt.Sprintf("\t\tuq.Predicates = append(uq.Predicates, p)\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\treturn uq\n")
	clientCode += fmt.Sprintf("}\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) All(ctx context.Context)([]*%s,error) {\n", structName, structName)
	clientCode += fmt.Sprintf("\tfilter := bson.D{}\n")
	clientCode += fmt.Sprintf("\tfor _, p := range uq.Predicates {\n")
	clientCode += fmt.Sprintf("\t\tp(&filter)\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\tcur, err := uq.Database(uq.dbName).Collection(%s.%sMongo).Find(ctx, filter)\n", strings.ToLower(structName), structName)
	clientCode += fmt.Sprintf("\tif err != nil {\n")
	clientCode += fmt.Sprintf("\t\treturn nil, err\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\tdefer cur.Close(ctx)\n")
	// todo:修改引用的包
	clientCode += fmt.Sprintf("\ttemp := make([]*%s, 0, 0)\n", structName)
	clientCode += fmt.Sprintf("\tfor cur.Next(ctx) {\n")
	clientCode += fmt.Sprintf("\t\tvar u %s\n", structName)
	clientCode += fmt.Sprintf("\t\terr = cur.Decode(&u)\n")
	clientCode += fmt.Sprintf("\t\tif err != nil {\n")
	clientCode += fmt.Sprintf("\t\t\treturn nil, err\n")
	clientCode += fmt.Sprintf("\t\t}\n")
	clientCode += fmt.Sprintf("\t\ttemp = append(temp, &u)\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\tif err = cur.Err(); err != nil {\n")
	clientCode += fmt.Sprintf("\t\treturn nil, err\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\treturn temp, nil\n")
	clientCode += fmt.Sprintf("}\n")

	return []byte(clientCode)

}

func generateFindFunction(structName string, fields []FieldInfo) string {
	var function string

	for _, field := range fields {
		function += getFindFunctionTemplate(structName, field, "")
		// todo:根据字段类型生成对应的比较符号
		if v, ok := ComparisonOperators[field.Type]; ok {
			for _, s := range v {
				function += getFindFunctionTemplate(structName, field, s)
			}
		}
	}
	return function
}

func getFindFunctionTemplate(structName string, field FieldInfo, op string) string {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func %s%s(v %s) %sPredicate {\n", field.Name, OpSplit(op), field.Type, structName))
	builder.WriteString("\treturn func(d *bson.D) {\n")
	builder.WriteString(fmt.Sprintf("\t\t*d = append(*d, bson.E{\n"))
	builder.WriteString(fmt.Sprintf("\t\t\tKey:   %sField,\n", field.Name))
	if op == "" {
		op = Eq
	}
	builder.WriteString(fmt.Sprintf("\t\t\tValue: bson.M{\"%s\": v},\n", op))
	builder.WriteString("\t\t})\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return builder.String()

}

func generateClientCode(structNameList []string) []byte {
	packageCode := fmt.Sprintf("package %s\n\n", strings.ToLower("go_mongo"))
	//clientCode += fmt.Sprintf("import \"go.mongodb.org/mongo-driver/bson\"\n\n")
	importCode := "import (\n"
	clientCode := fmt.Sprintf("type Client struct {\n")
	clientCode += "\tconfig\n"
	initFuncStr := "func (c *Client) init(){\n"
	for _, field := range structNameList {
		importCode += fmt.Sprintf("\t\"cc/go-mongo/%s\"\n", strings.ToLower(field))
		clientCode += fmt.Sprintf("\t%s *%sClient\n", field, field)
		initFuncStr += fmt.Sprintf("\tc.%s = New%sClient(c.config)\n", field, field)
	}
	importCode += ")\n"
	initFuncStr += "}\n"
	clientCode += "}\n\n"
	clientCode = packageCode + importCode + "\n" + clientCode + initFuncStr
	clientCode += fmt.Sprintf("func NewClient(opts ...Option) *Client {\n")
	clientCode += fmt.Sprintf("\tcfg := config{}\n")
	clientCode += fmt.Sprintf("\tcfg.options(opts...)\n")
	clientCode += fmt.Sprintf("\tclient := &Client{config: cfg}\n")
	clientCode += fmt.Sprintf("\tclient.init()\n")
	clientCode += fmt.Sprintf("\treturn client\n")
	clientCode += fmt.Sprintf("}\n")

	for _, field := range structNameList {
		// struct
		clientCode += fmt.Sprintf("type %sClient struct {\n", field)
		clientCode += "\tconfig\n"
		clientCode += "\tdbName string\n"
		clientCode += "}\n"

		clientCode += fmt.Sprintf("func (c *%sClient)SetDBName(dbName string)*%sClient{\n", field, field)
		clientCode += fmt.Sprintf("\tc.dbName=dbName\n")
		clientCode += fmt.Sprintf("\treturn c\n")
		clientCode += fmt.Sprintf("}\n")

		//new field client
		clientCode += fmt.Sprintf("func New%sClient(c config) *%sClient {\n", field, field)
		clientCode += fmt.Sprintf("\treturn &%sClient{ config: c }\n", field)
		clientCode += "}\n"

		clientCode += fmt.Sprintf("func(c *%sClient) Query() *%sQuery {\n", field, field)
		clientCode += fmt.Sprintf("\treturn &%sQuery{ \n", field)
		clientCode += fmt.Sprintf("\t\tconfig: c.config,\n")
		clientCode += fmt.Sprintf("\t\tPredicates: []%s.%sPredicate{},\n", strings.ToLower(field), field)

		clientCode += fmt.Sprintf("\t\tdbName: c.dbName,\n")
		clientCode += fmt.Sprintf("\t}\n")
		clientCode += "}\n"
	}

	return []byte(clientCode)
}

func generateConfig() []byte {
	clientCode := fmt.Sprintf("package %s\n\n", strings.ToLower("go_mongo"))
	clientCode += fmt.Sprintf("import \"go.mongodb.org/mongo-driver/mongo\"\n\n")
	clientCode += fmt.Sprintf("type config struct {\n")
	clientCode += fmt.Sprintf("\tmongo.Client\n")
	clientCode += "}\n"
	clientCode += "type Option func(*config)\n"
	clientCode += "func (c *config) options(opts ...Option) {"
	clientCode += fmt.Sprintf("\tfor _, opt := range opts {\n")
	clientCode += fmt.Sprintf("\t\topt(c)\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("}\n")
	clientCode += fmt.Sprintf("func Driver(driver mongo.Client) Option {\n")
	clientCode += fmt.Sprintf("\treturn func(c *config) {\n")
	clientCode += fmt.Sprintf("\t\tc.Client = driver\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("}\n")

	return []byte(clientCode)

}

func writeConstantsToFile(filename string, constants []byte) error {
	return os.WriteFile(filename, constants, 0644)
}
