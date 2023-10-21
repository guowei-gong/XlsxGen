package codegen

import "testing"

func TestExcelToJson(t *testing.T) {
	SetExcelPath("./")
	SetExcelToJsonPath("./")
	ExcelExport()
}

func TestExcelToGo(t *testing.T) {
	SetExcelPath("./")
	SetExcelToGoPath("./")
	ExcelExport()
}
