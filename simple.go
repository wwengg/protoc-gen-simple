package main

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"regexp"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	contextPackage      = protogen.GoImportPath("context")
	rpcxServerPackage   = protogen.GoImportPath("github.com/smallnest/rpcx/server")
	rpcxClientPackage   = protogen.GoImportPath("github.com/smallnest/rpcx/client")
	rpcxProtocolPackage = protogen.GoImportPath("github.com/smallnest/rpcx/protocol")
	SimplesrpcPackage   = protogen.GoImportPath("github.com/wwengg/simple/core/srpc")
	SimpleStorePackage  = protogen.GoImportPath("github.com/wwengg/simple/core/store")
	GormPackage         = protogen.GoImportPath("gorm.io/gorm")
)

// generateFile generates a _grpc.pb.go file containing gRPC service definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + ".simple.pb.go"
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
	generateFileContent(gen, file, g)
	return g
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// generateFileContent generates the gRPC service definitions, excluding thstarts a server only registers one service.e package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", contextPackage.Ident("TODO"))
	g.P("var _ = ", rpcxServerPackage.Ident("NewServer"))
	g.P("var _ = ", rpcxClientPackage.Ident("NewClient"))
	g.P("var _ = ", rpcxProtocolPackage.Ident("NewMessage"))
	g.P("var _ = ", SimplesrpcPackage.Ident("TODO"))
	g.P("var _ = ", SimpleStorePackage.Ident("TODO"))
	g.P("var _ = ", GormPackage.Ident("ErrEmptySlice"))
	g.P()
	g.P("//================== Model ===================")
	for _, message := range file.Messages {

		name := string(message.Desc.Name())
		if strings.Contains(name, "Model") {
			generateModelCode(g, message)
		}
	}
	g.P("//================== Model End ===================")

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := upperFirstLatter(service.GoName)

	g.P("//================== interface skeleton ===================")
	g.P(fmt.Sprintf(`type %s interface {`, serviceName))
	g.P(fmt.Sprintf(`// %s can be used for interface verification.`, serviceName))
	g.P()
	for _, method := range service.Methods {
		generateAbleCode(g, method)
	}
	g.P(fmt.Sprintf(`}`))

	g.P()
	g.P("//================== server skeleton ===================")
	g.P(fmt.Sprintf(`type %[1]sImpl struct {}

		// ServeFor%[1]s starts a server only registers one service.
		// You can register more services and only start one server.
		// It blocks until the application exits.
		func ServeFor%[1]s(addr string, rpc sconfig.RPC, rpcService sconfig.RpcService) error{
			s := server.NewServer()
			// 开启rpcx监控
			s.EnableProfile = true
			// 服务注册中心
			srpc.AddRegistryPlugin(s, rpc, rpcService)
			s.RegisterName("%[1]s", new(%[1]sImpl), "")
			return s.Serve("tcp", addr)
		}
	`, serviceName))
	g.P()
	for _, method := range service.Methods {
		generateServerCode(g, service, method)
	}

	g.P()
	g.P("//================== client stub ===================")
	g.P(fmt.Sprintf(`// %[1]s is a client wrapped XClient.
		type %[1]sClient struct{
			xclient client.XClient
		}

		// New%[1]sClient wraps a XClient as %[1]sClient.
		// You can pass a shared XClient object created by NewXClientFor%[1]s.
		func New%[1]sClient(xclient client.XClient) *%[1]sClient {
			return &%[1]sClient{xclient: xclient}
		}

		// NewXClientFor%[1]s creates a XClient.
		// You can configure this client with more options such as etcd registry, serialize type, select algorithm and fail mode.
		func NewXClientFor%[1]s(addr string) (client.XClient, error) {
			d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
			if err != nil {
				return nil, err
			}
			
			opt := client.DefaultOption
			opt.SerializeType = protocol.ProtoBuffer

			xclient := client.NewXClient("%[1]s", client.Failtry, client.RoundRobin, d, opt)

			return xclient,nil
		}
	`, serviceName))
	for _, method := range service.Methods {
		generateClientCode(g, service, method)
	}

	// one client
	g.P()
	g.P("//================== oneclient stub ===================")
	g.P(fmt.Sprintf(`// %[1]sOneClient is a client wrapped oneClient.
		type %[1]sOneClient struct{
			serviceName string
			oneclient *client.OneClient
		}

		// New%[1]sOneClient wraps a OneClient as %[1]sOneClient.
		// You can pass a shared OneClient object created by NewOneClientFor%[1]s.
		func New%[1]sOneClient(oneclient *client.OneClient) *%[1]sOneClient {
			return &%[1]sOneClient{
				serviceName: "%[1]s",
				oneclient: oneclient,
			}
		}

		// ======================================================
	`, serviceName))
	for _, method := range service.Methods {
		generateOneClientCode(g, service, method)
	}
}

func generateModelCode(g *protogen.GeneratedFile, message *protogen.Message) {
	name := g.QualifiedGoIdent(message.GoIdent)
	afterName, _ := strings.CutSuffix(name, "Model")
	g.P(fmt.Sprintf("//================== %s Model ===================", afterName))
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
			return &%[2]s{}
		}

		
		// Create%[1]s Func 创建
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
			err = global.DB_.Save(e).Error
			return err
		}

		// Update%[1]s 查询
		func Get%[1]s(id int64) (result %[1]s, err error) {
			err = global.DB_.Where("id = ?", id).First(&result).Error
			return
		}

		// 分页查询
		func Get%[1]sList(info pbbase.PageInfo) (list []%[1]s, total int64, err error) {
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
`, afterName, name))
}

func generateModelFiled(g *protogen.GeneratedFile, field *protogen.Field) {
	lowerFirstLatter(field.GoName)
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		g.P(fmt.Sprintf(`		%s  string `, field.GoName) + "`" + fmt.Sprintf(`json:"%s"gorm:"column:%s;comment: ;type:varchar(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.DoubleKind:
		g.P(fmt.Sprintf(`		%s  float64 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s"gorm:"column:%s;comment: ;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.Int64Kind:
		g.P(fmt.Sprintf(`		%s  int64 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s"gorm:"column:%s;comment: ;type:bigint(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	case protoreflect.Int32Kind:
		g.P(fmt.Sprintf(`		%s  int32 `, field.GoName) + "`" + fmt.Sprintf(`json:"%s"gorm:"column:%s;comment: ;type:smallint(6);size:6;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")
	default:
		g.P(fmt.Sprintf(`		%s  interface{} `, field.GoName) + "`" + fmt.Sprintf(`json:"%s"gorm:"column:%s;comment: ;type:any(20);size:20;"`, field.Desc.JSONName(), ToSnakeCase(field.GoName)) + "`")

	}

}

func generateServerCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is server rpc method as defined
		func (s *%sImpl) %s(ctx context.Context, args *%s, reply *%s) (err error){
			// TODO: add business logics

			// TODO: setting return values
			*reply = %s{}

			return nil
		}
	`, methodName, serviceName, methodName, inType, outType, outType))
}

func generateAbleCode(g *protogen.GeneratedFile, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %[1]s is server rpc method as defined
		%[1]s(ctx context.Context, args *%[2]s, reply *%[3]s) (err error)
	`, methodName, inType, outType))
}

func generateClientCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is client rpc method as defined
		func (c *%sClient) %s(ctx context.Context, args *%s)(reply *%s, err error){
			reply = &%s{}
			err = c.xclient.Call(ctx,"%s",args, reply)
			return reply, err
		}
	`, methodName, serviceName, methodName, inType, outType, outType, method.GoName))
}

func generateOneClientCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is client rpc method as defined
		func (c *%sOneClient) %s(ctx context.Context, args *%s)(reply *%s, err error){
			reply = &%s{}
			err = c.oneclient.Call(ctx,c.serviceName,"%s",args, reply)
			return reply, err
		}
	`, methodName, serviceName, methodName, inType, outType, outType, method.GoName))
}

// upperFirstLatter make the fisrt charater of given string  upper class
func upperFirstLatter(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// upperFirstLatter make the fisrt charater of given string  upper class
func lowerFirstLatter(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return strings.ToLower(string(s[0]))
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

var matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	str = matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                             //全部转小写
}
