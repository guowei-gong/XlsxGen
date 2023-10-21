package gosrc

const templateText = `// 勇敢🐶, 不怕困难
// 该文件由代码生成器自动生成
// 创建于 %v

package %v
		{{range .}}
type {{.GoStructName}} struct {
		{{range .Fields}}
		// {{.Des}}
		{{.Key}} {{.Typ}}
        {{end}}
}
{{end}}
`
