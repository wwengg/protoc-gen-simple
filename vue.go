package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func generateTableFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message) *protogen.GeneratedFile {
	name := string(message.Desc.Name())
	// fullName := string(message.Desc.FullName())

	afterName, _ := strings.CutSuffix(name, "Model")

	lowerName := lowerFirstLatter(afterName)
	filename := lowerName + ".vue"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P(`<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="query.title" placeholder="Title" style="width: 200px;" class="filter-item"
        @keyup.enter.native="handleFilter" />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit"
        @click="handleCreate">
        新建
      </el-button>
    </div>
    <el-table :key="tableKey" v-loading="listLoading" :data="tableData" border fit highlight-current-row
      style="width: 100%;">
      <el-table-column label="ID" prop="id" sortable="custom" align="center" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="CreatedAt" width="150px" align="center" prop="createdAt">
      </el-table-column>
      <el-table-column label="UpdatedAt" width="150px" align="center" prop="updatedAt">
      </el-table-column>`)
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		generateTableColumnFiled(g, field)

	}
	g.P(`      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-popover v-model="row.visible" placement="top" width="160">
            <p>确定要删除此用户吗</p>
            <div style="text-align: right; margin: 0">
              <el-button size="mini" type="text" @click="row.visible = false">取消</el-button>
              <el-button type="primary" size="mini" @click="handleDelete(row)">确定</el-button>
            </div>
            <el-button slot="reference" size="mini" type="danger">删除</el-button>
          </el-popover>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total > 0" :total="total" :page.sync="page" :limit.sync="pageSize" @pagination="getTableData" />
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px"
        style="width: 450px; margin-left:50px;">`)
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		generateFormFiled(g, field)

	}
	g.P(fmt.Sprintf(`      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus === 'create' ? createData() : updateData()">
          完成
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { create%[1]s, update%[1]s, delete%[1]s, find%[1]sById, find%[1]sList } from '@/api/%[2]s'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import tableList from '@/mixins/tableList'

export default {
  name: '%[1]sTable',
  components: { Pagination },
  directives: { waves },
  mixins: [tableList],
  data() {
    return {
      listApi: find%[1]sList,
      tableKey: 0,
      temp: {
        id: undefined,
        createdAt: '',
        updatedAt: '',
`, afterName, lowerName))
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		generateTempFiled(g, field)
	}
	g.P(`      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      rules: {
        //   type: [{ required: true, message: 'type is required', trigger: 'change' }],
        //   timestamp: [{ type: 'date', required: true, message: 'timestamp is required', trigger: 'change' }],
        //   title: [{ required: true, message: 'title is required', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    handleFilter() {
      this.page = 1
      this.getTableData()
    },
    handleModifyStatus(row, status) {
      this.$message({
        message: '操作Success',
        type: 'success'
      })
      row.status = status
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        createdAt: '',
        updatedAt: '',`)
	for _, field := range message.Fields {
		if field.GoName == "Id" || field.GoName == "CreatedAt" || field.GoName == "UpdatedAt" || field.GoName == "DeletedAt" {
			continue
		}
		generateTempFiled(g, field)
	}
	g.P(fmt.Sprintf(`      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    async createData() {
      this.$refs['dataForm'].validate(async (valid) => {
        if (valid) {
          const res = await create%[1]s(this.temp)
          if (res.code === 'Success') {
            this.handleFilter();
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
          }
        }
      })
    },
    async handleUpdate(row) {
      const res = await find%[1]sById({ id: row.id })
      console.log(res)
      if (res.code === 'Success') {
        this.temp = res.data
        this.dialogStatus = 'update'
        this.dialogFormVisible = true
        this.$nextTick(() => {
          this.$refs['dataForm'].clearValidate()
        })
      }
    },
    async updateData() {
      this.$refs['dataForm'].validate(async (valid) => {
        if (valid) {
          const res = await update%[1]s(this.temp)
          if (res.code === 'Success') {
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
            this.getTableData()
          }

        }
      })
    },
    async handleDelete(row) {
      await delete%[1]s({id:row.id})
      this.getTableData()
    }
  }
}
</script>
`, afterName))
	return g
}

func generateTableColumnFiled(g *protogen.GeneratedFile, field *protogen.Field) {
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		g.P(fmt.Sprintf(`      <el-table-column label="%s" width="150px" align="center" prop="%s">
      </el-table-column>`, field.GoName, field.Desc.JSONName()))
	case protoreflect.DoubleKind:
		g.P(fmt.Sprintf(`      <el-table-column label="%s" width="150px" align="center" prop="%s">
      </el-table-column>`, field.GoName, field.Desc.JSONName()))
	case protoreflect.Int64Kind:
		g.P(fmt.Sprintf(`      <el-table-column label="%s" width="150px" align="center" prop="%s">
      </el-table-column>`, field.GoName, field.Desc.JSONName()))
	case protoreflect.Int32Kind:
		g.P(fmt.Sprintf(`      <el-table-column label="%s" width="150px" align="center" prop="%s">
      </el-table-column>`, field.GoName, field.Desc.JSONName()))
	default:
		g.P(fmt.Sprintf(`      <el-table-column label="%s" width="150px" align="center" prop="%s">
      </el-table-column>`, field.GoName, field.Desc.JSONName()))
	}

}

func generateFormFiled(g *protogen.GeneratedFile, field *protogen.Field) {
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		g.P(fmt.Sprintf(`        <el-form-item label="%s" prop="%s">
          <el-input v-model="temp.%s" />
        </el-form-item>`, field.GoName, field.Desc.JSONName(), field.Desc.JSONName()))
	case protoreflect.DoubleKind:
		g.P(fmt.Sprintf(`        <el-form-item label="%s" prop="%s">
          <el-input v-model="temp.%s" />
        </el-form-item>`, field.GoName, field.Desc.JSONName(), field.Desc.JSONName()))
	case protoreflect.Int64Kind:
		g.P(fmt.Sprintf(`        <el-form-item label="%s" prop="%s">
          <el-input v-model="temp.%s" />
        </el-form-item>`, field.GoName, field.Desc.JSONName(), field.Desc.JSONName()))
	case protoreflect.Int32Kind:
		g.P(fmt.Sprintf(`        <el-form-item label="%s" prop="%s">
          <el-input v-model="temp.%s" />
        </el-form-item>`, field.GoName, field.Desc.JSONName(), field.Desc.JSONName()))
	default:
		g.P(fmt.Sprintf(`        <el-form-item label="%s" prop="%s">
          <el-input v-model="temp.%s" />
        </el-form-item>`, field.GoName, field.Desc.JSONName(), field.Desc.JSONName()))
	}

}

func generateTempFiled(g *protogen.GeneratedFile, field *protogen.Field) {
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		g.P(fmt.Sprintf(`        %s: '',`, field.Desc.JSONName()))
	case protoreflect.DoubleKind:
		g.P(fmt.Sprintf(`        %s: 0,`, field.Desc.JSONName()))
	case protoreflect.Int64Kind:
		g.P(fmt.Sprintf(`        %s: 0,`, field.Desc.JSONName()))
	case protoreflect.Int32Kind:
		g.P(fmt.Sprintf(`        %s: 0,`, field.Desc.JSONName()))
	default:
		g.P(fmt.Sprintf(`        %s: '',`, field.Desc.JSONName()))
	}

}

func generateApiCode(gen *protogen.Plugin, file *protogen.File, service *protogen.Service) {
	serviceName := upperFirstLatter(service.GoName)

	filename := lowerFirstLatter(serviceName) + ".js"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P(`import request from '@/utils/request'
import protoRoot from '@/proto/proto.js'
`)
	for _, method := range service.Methods {
		//inType := g.QualifiedGoIdent(method.Input.GoIdent)
		//method.Input.Desc.FullName()
		//outType := g.QualifiedGoIdent(method.Output.GoIdent)
		// outType := method.Output.Desc.FullName()
		// methodName := upperFirstLatter(method.GoName)
		g.P(fmt.Sprintf(`export function %[1]s(data) {
  var buffer = protoRoot.%[3]s.encode(data).finish().slice().buffer
  return request({
    url: '/v2/%[2]s/%[1]s',
    method: 'post',
    buffer,
    pb: '%[4]s'
  })
}
`, lowerFirstLatter(method.GoName), lowerFirstLatter(serviceName), method.Input.Desc.FullName(), method.Output.Desc.FullName()))
	}

}
