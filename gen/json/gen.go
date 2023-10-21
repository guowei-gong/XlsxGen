package json

import (
	"fmt"
	"guowei.com/XlsxGen/gen/model"
	"os"
	"path/filepath"
	"strings"
)

const (
	filetypeString                   = "string"
	filetypeTime                     = "time"
	filetypeArrayInt                 = "[]int"
	filetypeTwoDimensionalArrayInt32 = "[][]int32"
	filetypeArrayFloat               = "[]float64"
)

func Generate(structs []*model.Model, excelToJsonPath string) {
	// 将生成的数据分成多个 .json 文件, 粒度为 sheet, 可读性更好
	// 也可以考虑将数据合并为一个 .json 文件, 加载到内存时, 内核态切换成本更低, 加载的更快
	for _, s := range structs {
		jsonFile := fmt.Sprintf("%s.json", s.JsonName)

		if err := output(jsonFile, toJson(s.Dataset, s.Fields), excelToJsonPath); err != nil {
			panic(err.Error())
		}

	}
}

func output(filename, str, excelToJsonPath string) error {
	f, err := os.OpenFile(excelToJsonPath+string(filepath.Separator)+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		return err
	}

	return nil
}

func toJson(dataRows []model.RowData, metaSet []*model.Meta) string {
	ret := "["
	for _, row := range dataRows {
		ret += "\n\t{"
		for idx, meta := range metaSet {

			if len(meta.Key) == 0 {
				continue
			}

			ret += fmt.Sprintf("\n\t\t\"%s\":", meta.Key)
			switch meta.Typ {
			case filetypeString:
				if row[idx] == nil {
					ret += "\"\""
				} else {
					ret += fmt.Sprintf("\"%s\"", row[idx])
				}
			case filetypeTime:
			case filetypeArrayInt, filetypeArrayFloat:
				if row[idx] == nil || row[idx] == "" {
					ret += "[]"
				} else {
					s := strings.Split(row[idx].(string), ",")

					var result []interface{}

					for _, v := range s {
						result = append(result, v)
					}

					ret += fmt.Sprintf("%s", formatArray(result))
				}
			case filetypeTwoDimensionalArrayInt32:
				if row[idx] == nil || row[idx] == "" {
					ret += "[]"
				} else {

					str := row[idx].(string) // 二维数组配置有固定格式, e.g: 1,1; 2,2; 3,3 -> [[1,1], [2,2], [3,3]]

					parts := strings.Split(str, ";")

					var result [][]interface{}

					for _, part := range parts {
						s := strings.Split(part, ",")

						var item []interface{}

						for _, v := range s {
							item = append(item, v)
						}

						result = append(result, item)
					}

					ret += fmt.Sprintf("%s", formatTwoDimensionalArray(result))
				}
			default:
				if row[idx] == nil || row[idx] == "" {
					ret += "0"
				} else {
					ret += fmt.Sprintf("%s", row[idx])
				}
			}
			ret += ","
		}
		ret = ret[:len(ret)-1]
		ret += "\n\t},"
	}
	ret = ret[:len(ret)-1]
	ret += "\n]"
	return ret
}

func formatArray(result []interface{}) string {
	strArr := make([]string, len(result))
	for i, val := range result {
		strArr[i] = fmt.Sprintf("%v", val)
	}
	return "[" + strings.Join(strArr, ",") + "]"
}

func formatTwoDimensionalArray(result [][]interface{}) string {
	var output string

	output += "["
	for i, row := range result {
		output += "["
		for j, val := range row {
			if j > 0 {
				output += ","
			}
			output += fmt.Sprintf("%v", val)
		}
		output += "]"

		if i < len(result)-1 {
			output += ","
		}
	}
	output += "]"

	return output
}
