package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zaihui/mongoent"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type FieldInfo struct {
	Name     string
	JSONName string
	Type     string
}

func init() {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate ent model factory",
		Run:   GetStructsFromGoFile,
	}
	RootCmd.AddCommand(cmd)
}

func GetStructsFromGoFile(cmd *cobra.Command, _ []string) {
	modelFilePath, outputPath, modPath, err := ExtraFlags(cmd)
	if err != nil {
		fmt.Println("Error createModel:", err)
		os.Exit(1)
	}
	finalOutputPath := outputPath + "/" + mongoent.MongoSchema
	err = createDirectory(finalOutputPath)
	if err != nil {
		fmt.Println("Error createDirectory:", err)
		os.Exit(1)
	}

	structs, fileNameList := getStructsFromFile(modelFilePath)
	err = createModel(modelFilePath, finalOutputPath)
	if err != nil {
		fmt.Println("Error createModel:", err)
		os.Exit(1)
	}
	for i, s := range structs {
		fields := getFieldsFromStruct(s)
		constants, err := generateConstants(fileNameList[i], fields)
		if err != nil {
			fmt.Println("Error generateConstants:", err)
			os.Exit(1)
		}
		filename := strings.ToLower(fileNameList[i]) + ".go"
		err = createDirectory(finalOutputPath + "/" + fileNameList[i])
		if err != nil {
			fmt.Println("Error creating directory:", err)
			os.Exit(1)
		}
		err = writeConstantsToFile(finalOutputPath+"/"+fileNameList[i]+"/"+filename, constants)
		if err != nil {
			fmt.Println("Error writing constants to file:", err)
			os.Exit(1)
		}
		createQuery(finalOutputPath, modPath, fileNameList[i])

	}
	createClient(finalOutputPath, modPath, fileNameList)
	createConfig(finalOutputPath)

}

func createClient(outputPath string, modPath string, fileNameList []string) {
	err := writeConstantsToFile(outputPath+"/"+"client.go", generateClientCode(fileNameList, modPath))
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createConfig(outputPath string) {
	err := writeConstantsToFile(outputPath+"/"+"config.go", generateConfig())
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createQuery(outputPath, modPath string, structName string) {
	err := writeConstantsToFile(outputPath+"/"+fmt.Sprintf("%s_query.go", strings.ToLower(structName)), generateQuery(structName, modPath))
	if err != nil {
		fmt.Println("Error writing constants to file:", err)
		os.Exit(1)
	}
}

func createDirectory(dirName string) error {
	dirName = strings.ToLower(dirName)

	// 根据传入的路径获取所有中间文件夹
	directories := getAllDirectories(dirName)

	// 遍历中间文件夹并逐个创建
	for _, dir := range directories {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getAllDirectories(path string) []string {
	// 根据路径获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	// 分隔路径中的目录部分
	dirs := strings.Split(absPath, string(filepath.Separator))

	// 构建所有中间文件夹的路径
	var directories []string
	currentPath := ""
	for _, dir := range dirs {
		currentPath = filepath.Join(currentPath, "/"+dir)
		directories = append(directories, currentPath)
	}

	return directories
}

func createModel(filename string, outputPath string) error {
	modelData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read model file:", err)
		return err
	}
	// 按分割符号 "type " 将内容拆分为每个结构体的定义
	structureDefinitions := strings.Split(string(modelData), "type ")
	// 遍历每个结构体定义并生成对应的文件
	for _, structDef := range structureDefinitions {
		structDef = strings.TrimSpace(structDef)
		if structDef == "" {
			continue
		}

		lines := strings.Split(structDef, "\n")
		if len(lines) == 0 {
			continue
		}

		structName := strings.SplitN(lines[0], " ", 2)[0]
		if structName == "package" {
			continue
		}

		c := fmt.Sprintf("package %s\n\n"+"type "+structDef, mongoent.MongoSchema)
		outputFile := outputPath + "/" + strings.ToLower(structName) + ".go"
		err = writeConstantsToFile(outputFile, []byte(c))
		if err != nil {
			fmt.Println("Failed to create output file:", err)
			return err
		}

		fmt.Println(outputFile, "generated successfully.")
	}
	return nil
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
	builder.WriteString(fmt.Sprintf("%sMongo = \"%s\"\n", fileName, mongoent.ToSnakeCase(fileName)))
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

func generateQuery(structName string, modPath string) []byte {
	clientCode := fmt.Sprintf("package %s\n\n", mongoent.MongoSchema)

	clientCode += fmt.Sprintf("import (\n")
	clientCode += fmt.Sprintf("\t\"context\"\n")

	clientCode += fmt.Sprintf("\t\"%s/%s\"\n", modPath, strings.ToLower(structName))
	clientCode += fmt.Sprintf("\t\"go.mongodb.org/mongo-driver/bson\"\n")
	clientCode += fmt.Sprintf("\t\"go.mongodb.org/mongo-driver/mongo\"\n")
	clientCode += fmt.Sprintf("\t\"go.mongodb.org/mongo-driver/mongo/options\"\n")

	clientCode += fmt.Sprintf(")\n\n")

	clientCode += fmt.Sprintf("type %sQuery struct {\n", structName)
	clientCode += fmt.Sprintf("\tconfig\n")
	clientCode += fmt.Sprintf("\tPredicates []%s.%sPredicate\n", strings.ToLower(structName), structName)
	clientCode += fmt.Sprintf("\tlimit  *int64\n")
	clientCode += fmt.Sprintf("\toffset *int64\n")
	clientCode += fmt.Sprintf("\tdbName string\n")
	clientCode += fmt.Sprintf("\toptions bson.D\n\n")
	clientCode += fmt.Sprintf("}\n\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) Limit(limit int64) *%sQuery{\n", structName, structName)
	clientCode += fmt.Sprintf("\tuq.limit = &limit\n")
	clientCode += fmt.Sprintf("\treturn uq\n")
	clientCode += fmt.Sprintf("}\n\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) Offset(offset int64) *%sQuery{\n", structName, structName)
	clientCode += fmt.Sprintf("\tuq.offset = &offset\n")
	clientCode += fmt.Sprintf("\treturn uq\n")
	clientCode += fmt.Sprintf("}\n\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) Order(o ...OrderFunc) *%sQuery {\n"+
		"\tfor _, fn := range o {\n"+
		"\t\tfn(&uq.options)\n"+
		"\t}\n"+
		"\treturn uq\n}\n\n", structName, structName)

	clientCode += fmt.Sprintf("func (uq *%sQuery) Where(ps ...%s.%sPredicate)*%sQuery{\n", structName, strings.ToLower(structName), structName, structName)
	//clientCode += fmt.Sprintf("\tfilter := bson.D{}\n")
	clientCode += fmt.Sprintf("\tfor _, p := range ps {\n")
	clientCode += fmt.Sprintf("\t\tuq.Predicates = append(uq.Predicates, p)\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\treturn uq\n")
	clientCode += fmt.Sprintf("}\n\n")

	clientCode += fmt.Sprintf("func (uq *%sQuery) All(ctx context.Context)([]*%s,error) {\n", structName, structName)
	clientCode += fmt.Sprintf("\tfilter := bson.D{}\n")
	clientCode += fmt.Sprintf("\tfor _, p := range uq.Predicates {\n")
	clientCode += fmt.Sprintf("\t\tp(&filter)\n")
	clientCode += fmt.Sprintf("\t}\n\n")

	clientCode += fmt.Sprintf("\to := options.Find()\n")
	clientCode += fmt.Sprintf("\tif uq.limit != nil && *uq.limit != 0 {\n")
	clientCode += fmt.Sprintf("\t\to = o.SetLimit(*uq.limit)\n")
	clientCode += fmt.Sprintf("\t}\n")

	clientCode += fmt.Sprintf("\tif uq.offset != nil && *uq.offset != 0 {\n")
	clientCode += fmt.Sprintf("\t\to = o.SetSkip(*uq.offset)\n")
	clientCode += fmt.Sprintf("\t}\n")

	clientCode += fmt.Sprintf("\to.SetSort(uq.options)\n")

	clientCode += fmt.Sprintf("\tcur, err := uq.Database(uq.dbName).Collection(%s.%sMongo).Find(ctx, filter,o)\n", strings.ToLower(structName), structName)
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
	clientCode += fmt.Sprintf("func (uq *%sQuery) First(ctx context.Context) (*%s, error) {\n"+
		"\tdocument, err := uq.Limit(1).All(ctx)\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tif len(document) == 0 {\n"+
		"\t\treturn nil, mongo.ErrNilDocument\n"+
		"\t}\n"+
		"\treturn document[0], err\n}", structName, structName)

	return []byte(clientCode)

}

func generateFindFunction(structName string, fields []FieldInfo) string {
	var function string

	for _, field := range fields {
		function += getFindFunctionTemplate(structName, field, "")
		if v, ok := mongoent.ComparisonOperators[field.Type]; ok {
			for _, s := range v {
				function += getFindFunctionTemplate(structName, field, s)
			}
		}
	}
	return function
}

func getFindFunctionTemplate(structName string, field FieldInfo, op string) string {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func %s%s(v %s) %sPredicate {\n", field.Name, mongoent.OpSplit(op), field.Type, structName))
	builder.WriteString("\treturn func(d *bson.D) {\n")
	builder.WriteString(fmt.Sprintf("\t\t*d = append(*d, bson.E{\n"))
	builder.WriteString(fmt.Sprintf("\t\t\tKey:   %sField,\n", field.Name))
	if op == "" {
		op = mongoent.Eq
	}
	builder.WriteString(fmt.Sprintf("\t\t\tValue: bson.M{\"%s\": v},\n", op))
	builder.WriteString("\t\t})\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return builder.String()

}

func generateClientCode(structNameList []string, modPath string) []byte {
	packageCode := fmt.Sprintf("package %s\n\n", mongoent.MongoSchema)
	//clientCode += fmt.Sprintf("import \"go.mongodb.org/mongo-driver/bson\"\n\n")
	importCode := "import (\n"
	importCode += "\t\"go.mongodb.org/mongo-driver/bson\"\n"
	clientCode := fmt.Sprintf("type Client struct {\n")
	clientCode += "\tconfig\n"
	initFuncStr := "func (c *Client) init(){\n"
	for _, field := range structNameList {
		importCode += fmt.Sprintf("\t\"%s/%s\"\n", modPath, strings.ToLower(field))
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
		clientCode += fmt.Sprintf("\t\toptions:    bson.D{},\n")
		clientCode += fmt.Sprintf("\t}\n")
		clientCode += "}\n\n"
	}
	clientCode += fmt.Sprintf("type OrderFunc func(*bson.D)\n\n")
	clientCode += fmt.Sprintf("func Desc(field string) OrderFunc {\n" +
		"\treturn func(sort *bson.D) {\n" +
		"\t\t*sort = append(*sort, bson.E{Key: field, Value: -1})\n\n" +
		"\t}\n}\n")
	clientCode += fmt.Sprintf("func Asc(field string) OrderFunc {\n" +
		"\treturn func(sort *bson.D) {\n" +
		"\t\t*sort = append(*sort, bson.E{Key: field, Value: 1})\n" +
		"\t}\n}\n")

	return []byte(clientCode)
}

func generateConfig() []byte {
	clientCode := fmt.Sprintf("package %s\n\n", mongoent.MongoSchema)
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
	filename = strings.ToLower(filename)
	return os.WriteFile(filename, constants, 0644)
}
