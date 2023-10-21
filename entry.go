package codegen

import (
	"github.com/xuri/excelize/v2"
	"guowei.com/XlsxGen/gen/gosrc"
	"guowei.com/XlsxGen/gen/json"
	"guowei.com/XlsxGen/gen/model"
	"os"
	"path/filepath"
	"strings"
)

var (
	fileExt         = ".xlsx"
	excelPath       = ""
	goPackage       = "codegen"
	excelToGoPath   = ""
	excelToJsonPath = ""
)

// SetExcelPath 设置 Excel 文件所在路径
func SetExcelPath(path string) {
	p, err := filepath.Abs(path)
	if err != nil {
		panic("set excel file path: " + path)
	}
	excelPath = p
}

// SetExcelToGoPath 设置根据 Excel 导出的 Go 文件地址
func SetExcelToGoPath(path string) {
	p, err := filepath.Abs(path)
	if err != nil {
		panic("set excel file path: " + path)
	}
	excelToGoPath = p
}

// SetExcelToJsonPath 设置根据 Excel 导出的 Json 文件地址
func SetExcelToJsonPath(path string) {
	p, err := filepath.Abs(path)
	if err != nil {
		panic("set excel file path: " + path)
	}
	excelToJsonPath = p
}

// ExcelExport 根据 Excel 生成 Json 和 Go 文件
func ExcelExport() {
	files, err := getFileList(excelPath, fileExt)
	if err != nil {
		panic(err.Error())
	}

	var structs []*model.Model
	for _, file := range files {
		structs = append(structs, parseXlsx(file)...)
	}

	if len(excelToJsonPath) != 0 {
		json.Generate(structs, excelToJsonPath)
	}

	if len(excelToGoPath) != 0 {
		gosrc.Generate(structs, goPackage, excelToGoPath)
	}

}

func parseXlsx(file *model.FileInfo) []*model.Model {
	xlsx, err := excelize.OpenFile(file.Path)
	if err != nil {
		panic(err.Error())
	}

	sheets := xlsx.GetSheetList()
	result := make([]*model.Model, len(sheets))

	for i, sheet := range sheets {
		rows, err := xlsx.GetRows(sheet)
		if err != nil {
			panic(err.Error())
		}

		// 考虑健壮性, 不应该使用 rows[1] 和 len(rows)-2 这种方式
		// 定义的类型数量
		typeQuantity := len(rows[1])
		metaSet := make([]*model.Meta, 0, typeQuantity)
		dataset := make([]model.RowData, 0, len(rows)-2)

		for line, row := range rows {

			switch line {
			case 0: // 字段注释
				for idx, doc := range row {
					metaSet = append(metaSet, &model.Meta{Des: doc, Idx: idx})
				}
			case 1: // 字段类型
				for idx, typ := range row {
					if len(typ) == 0 || typ == "#" {
						continue
					}
					metaSet[idx].Typ = typ
				}

			case 2: // 字段名

				for idx, colName := range row {
					if colName == "#" || len(colName) == 0 {
						continue
					}
					metaSet[idx].Key = colName
				}

			default: // 从第 6 行开始, 表数据
				if line < 5 {
					continue
				}
				data := make(model.RowData, typeQuantity)
				for k := 0; k < typeQuantity; k++ {
					if i < len(row) {
						data[k] = row[k]
					}
				}
				dataset = append(dataset, data)

			}
		}

		item := &model.Model{
			Fields:       metaSet,
			Dataset:      dataset,
			JsonName:     file.Name + "_" + sheet,
			GoStructName: file.Name,
		}
		result[i] = item
	}

	return result
}

// getFileList 获取文件目录
func getFileList(path, ext string) (ret []*model.FileInfo, err error) {
	files, err := os.ReadDir(path)

	for _, file := range files {
		if strings.Contains(file.Name(), ".~") {
			continue
		}

		if filepath.Ext(file.Name()) == ext {

			realPath := path + "/" + file.Name()

			// TODO: 递归子文件夹
			if !file.IsDir() {
				ret = append(ret, &model.FileInfo{
					Path: realPath,
					// TODO: ".xlsx", ".xls", ".xlsm"
					Name: strings.TrimSuffix(file.Name(), ".xlsx"),
				})
			}

		}
	}

	return
}
