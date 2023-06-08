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
	//builder.WriteString(generateQueryStruct(fileName))
	builder.WriteString("const (")
	builder.WriteString(fmt.Sprintf("%sMongo = \"%s\"\n", fileName, ToSnakeCase(fileName)))
	for _, field := range fields {
		builder.WriteString(fmt.Sprintf("%sField = \"%s\"\n", field.Name, field.JSONName))
	}
	builder.WriteString(")\n")
	builder.WriteString(generateFindFunction(fileName, fields))

	src, err := format.Source([]byte(builder.String()))
	if err != nil {
		return []byte{}, err
	}
	return src, nil
}

func generateQuery(structName string) []byte {
	clientCode := "package go_mongo\n\n"
	clientCode += fmt.Sprintf("import \"go.mongodb.org/mongo-driver/bson\"\n\n")
	clientCode += fmt.Sprintf("type %sQuery struct {\n", structName)
	clientCode += fmt.Sprintf("\tconfig\n")
	clientCode += fmt.Sprintf("\tConditions bson.M\n")
	clientCode += fmt.Sprintf("}\n")

	// Where adds a new predicate for the UserQuery builder.
	// func (uq *UserInfoQuery) Where(ps ...bson.M) *UserInfoQuery {
	//	for _, p := range ps {
	//		for s, v := range p {
	//			uq.Conditions[s] = v
	//		}
	//	}
	//	return uq
	//}
	clientCode += fmt.Sprintf("func (uq *%sQuery) Where(ps ...bson.M)*%sQuery{\n", structName, structName)
	clientCode += fmt.Sprintf("\tfor _, p := range ps {\n")
	clientCode += fmt.Sprintf("\t\tfor s, v := range p {\n")
	clientCode += fmt.Sprintf("\t\t\tuq.Conditions[s] = v\n")
	clientCode += fmt.Sprintf("\t\t}\n")
	clientCode += fmt.Sprintf("\t}\n")
	clientCode += fmt.Sprintf("\treturn uq\n")
	clientCode += fmt.Sprintf("}\n")

	return []byte(clientCode)

}

func generateFindFunction(structName string, fields []FieldInfo) string {
	var function string
	for _, field := range fields {
		function += fmt.Sprintf("func Find%sBy", structName)
		function += fmt.Sprintf("%s(%s %s) bson.M {\n", field.Name, ConvertToCamelCase(field.JSONName), field.Type)
		function += fmt.Sprintf("\treturn bson.M{%sField: %s}", field.Name, ConvertToCamelCase(field.JSONName))
		function += "}\n\n"
	}
	return function
}

func generateClientCode(structNameList []string) []byte {
	clientCode := fmt.Sprintf("package %s\n\n", strings.ToLower("go_mongo"))

	clientCode += fmt.Sprintf("type Client struct {\n")
	clientCode += "\tconfig\n"
	for _, field := range structNameList {
		clientCode += fmt.Sprintf("\t%s *%sClient\n", field, field)
	}
	clientCode += "}\n\n"
	for _, field := range structNameList {
		// struct
		clientCode += fmt.Sprintf("type %sClient struct {\n", field)
		clientCode += "\tconfig\n"
		clientCode += "}\n"

		//new client
		clientCode += fmt.Sprintf("func New%s(c config) *%sClient {\n", field, field)
		clientCode += fmt.Sprintf("\treturn &%sClient{ config: c }\n", field)
		clientCode += "}\n"

		// query
		//func (c *UserClient) Query() *UserQuery {
		//	return &UserQuery{
		//	config: c.config,
		//}
		//}
		clientCode += fmt.Sprintf("func(c *%sClient) Query() *%sQuery {\n", field, field)
		clientCode += fmt.Sprintf("\treturn &%sQuery{ config: c.config }\n", field)
		clientCode += "}\n"
	}

	return []byte(clientCode)
}

func generateConfig() []byte {
	clientCode := fmt.Sprintf("package %s\n\n", strings.ToLower("go_mongo"))
	clientCode += fmt.Sprintf("type config struct {\n")
	clientCode += "}\n"
	return []byte(clientCode)

}

func writeConstantsToFile(filename string, constants []byte) error {
	return os.WriteFile(filename, constants, 0644)
}
