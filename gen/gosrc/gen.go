package gosrc

import (
	"bytes"
	"fmt"
	"guowei.com/XlsxGen/gen/model"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

func Generate(modelData []*model.Model, goPackage, excelToGoPath string) {
	// 解析并装载模板
	tpl := fmt.Sprintf(templateText, time.Now().Format("2006-01-02 15:04:05"), goPackage)

	t := template.New("go_text")
	tp, err := t.Parse(tpl)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer

	if err := tp.Execute(&buffer, modelData); err != nil {
		panic(err)
	}

	// 格式化生成的源文件
	file := excelToGoPath + string(filepath.Separator) + "dataset.go"
	os.WriteFile(file, buffer.Bytes(), 0666)

	exec.Command("gofmt", "-w", file).Run()
}
