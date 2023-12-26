// @Title
// @Description
// @Author  Wangwengang  2023/12/25 05:42
// @Update  Wangwengang  2023/12/25 05:42
package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

func generateModelFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message) *protogen.GeneratedFile {
	name := string(message.Desc.Name())
	fullName := string(message.Desc.FullName())

	afterName, _ := strings.CutSuffix(name, "Model")

	lowerName := lowerFirstLatter(afterName)
	filename := lowerName + "_model.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-simple. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-simple v", version)
	g.P("// - protoc          ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", SimpleStorePackage.Ident("TODO"))
	g.P()
	g.P(fmt.Sprintf(`// %s Model
		type %s struct {
			store.BASE_MODEL
`, afterName, afterName))
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		generateModelFiled(g, field)

	}
	g.P(`		}`)
	g.P()
	g.P(fmt.Sprintf(`

		func (model *%[1]s) Proto() *%[2]s {
			return &%[2]s{
				Id: model.ID,
				CreatedAt: model.CreatedAt.Format(time.DateTime),
				UpdatedAt: model.UpdatedAt.Format(time.DateTime),
`, afterName, fullName))
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		g.P(fmt.Sprintf(`				%s: model.%s,`, field.GoName, field.GoName))

	}
	g.P(`			}
		}`)

	g.P(fmt.Sprintf(`
		func %[1]sProtoToModel(proto *%[2]s) *%[1]s {
			%[3]s := %[1]s{
				BASE_MODEL: store.BASE_MODEL{
					ID:        proto.Id,
				},
`, afterName, fullName, lowerFirstLatter(afterName)))
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		g.P(fmt.Sprintf(`				%s: proto.%s,`, field.GoName, field.GoName))

	}
	g.P(fmt.Sprintf(`			}
			if createdAt,err := time.Parse(time.DateTime,proto.CreatedAt);err == nil{
				user.CreatedAt = createdAt
			}
			if updatedAt,err := time.Parse(time.DateTime,proto.UpdatedAt);err == nil{
				user.CreatedAt = updatedAt
			}
			return &%[1]s
		}`, lowerFirstLatter(afterName)))

	g.P()
	g.P(fmt.Sprintf(`// Create%[1]s Func 创建
		func Create%[1]s(a %[1]s) (err error) {
			err = global.DB_.Create(&a).Error
			return err
		}
		
		// Delete%[1]s  删除
		func Delete%[1]s(a %[1]s) (err error) {
			err = global.DB_.Delete(&a).Error
			return err
		}

		// Update%[1]s 修改
		func Update%[1]s(a *%[1]s) (err error) {
			err = global.DB_.Save(a).Error
			return err
		}

		// Update%[1]s 查询
		func Get%[1]s(id int64) (result %[1]s, err error) {
			err = global.DB_.Where("id = ?", id).First(&result).Error
			return
		}

		// 分页查询
		func Get%[1]sList(info pbcommon.PageInfo) (list []%[1]s, total int64, err error) {
			limit := info.PageSize
			offset := info.PageSize * (info.Page - 1)
			db := global.DB_.Model(&%[1]s{})
			var %[1]sList []%[1]s
			// 此处增加查询条件
			//if info.Keyword != "" {
			//	db.Where("keywaord = ?", info.Keyword)
			//}
			err = db.Count(&total).Error
			if err != nil {
				return %[1]sList, total, err
			} else {
				err = db.Limit(int(limit)).Offset(int(offset)).Find(&%[1]sList).Error
			}
			return %[1]sList, total, err
		}
`, afterName, fullName))
	return g
}

func generateModelFiled(g *protogen.GeneratedFile, field *protogen.Field) {
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		g.P(fmt.Sprintf(`		%s  string `, field.GoName) + "`" + fmt.Sprintf(`json:"%s" gorm:"column:%s;comment: ;type:varchar(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.DoubleKind:
		g.P(fmt.Sprintf(`		%s  float64 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s" gorm:"column:%s;comment: ;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.Int64Kind:
		g.P(fmt.Sprintf(`		%s  int64 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s" gorm:"column:%s;comment: ;type:bigint(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.Int32Kind:
		g.P(fmt.Sprintf(`		%s  int32 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s" gorm:"column:%s;comment: ;type:smallint(6);size:6;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	default:
		g.P(fmt.Sprintf(`		%s  interface{} `, field.GoName) + "`" + fmt.Sprintf(`json:"%s" gorm:"column:%s;comment: ;type:any(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")

	}

}

func generateSimpleServerCode(gen *protogen.Plugin, file *protogen.File, service *protogen.Service) {
	serviceName := upperFirstLatter(service.GoName)

	filename := lowerFirstLatter(serviceName) + "_service.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-simple. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-simple v", version)
	g.P("// - protoc          ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", SimpleStorePackage.Ident("TODO"))
	g.P("var _ = ", contextPackage.Ident("TODO"))
	g.P()
	g.P(fmt.Sprintf(`type %[1]s struct {}
	`, serviceName))
	for _, method := range service.Methods {
		//inType := g.QualifiedGoIdent(method.Input.GoIdent)
		//method.Input.Desc.FullName()
		//outType := g.QualifiedGoIdent(method.Output.GoIdent)
		outType := method.Output.Desc.FullName()
		methodName := upperFirstLatter(method.GoName)
		g.P(fmt.Sprintf(`// %s is server rpc method as defined
		func (s *%s) %s(ctx context.Context, args *%s, reply *%s) (err error){
			// TODO: add business logics

			// TODO: setting return values
			*reply = %s{}

			return nil
		}
	`, methodName, serviceName, methodName, method.Input.Desc.FullName(), outType, outType))
	}

}
