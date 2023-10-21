package model

type RowData []interface{}

// Model 用于填充 template/text 模板中的值
type Model struct {
	// 生成 go 代码使用
	GoStructName string
	Fields       []*Meta

	// 生成 json 文件使用
	JsonName string
	Dataset  []RowData
}

type Meta struct {
	Idx int
	Key string
	Typ string
	Des string // 代码中的注释
}

type FileInfo struct {
	Path string
	Name string
}
