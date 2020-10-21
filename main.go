package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/axgle/mahonia"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
)

func creat_bar(code string, name string, path string, loop int) {
	pdf := gofpdf.New("L", "mm", "A4", "")

	// pdf := gofpdf.New("P", "mm", 267, 200)
	pdf.AddUTF8Font("宋体", "R", "Font/simsun.ttf")
	pdf.AddUTF8Font("黑体", "L", "Font/msyhl.ttf")
	pdf.AddUTF8Font("黑体", "R", "Font/msyhR.ttf")

	for i := 0; i < loop; i++ {

		pdf.AddPage()
		// 创建条形码
		bcode, err := code128.Encode(code)

		if err == nil {
			key := barcode.Register(bcode)
			var width float64 = 254.5
			var height float64 = 77.6
			barcode.BarcodeUnscalable(pdf, key, 21, 11.5, &width, &height, false)
			// err = pdf.OutputFileAndClose("barcode.pdf")
		}

		//创建文本
		pdf.SetFont("宋体", "R", 84)
		pdf.SetLeftMargin(11.5)
		pdf.SetRightMargin(11.5)
		pdf.SetY(111)
		pdf.CellFormat(0, 0, code, "0", 0, "C", false, 0, "")

		pdf.SetFont("黑体", "L", 60)
		pdf.SetWordSpacing(30)
		pdf.SetLeftMargin(11.5)
		pdf.SetRightMargin(11.5)
		pdf.SetY(148)
		pdf.CellFormat(0, 0, name, "0", 0, "C", false, 0, "")

		pdf.SetFont("黑体", "R", 58)
		pdf.SetLeftMargin(11.5)
		pdf.SetRightMargin(11.5)
		pdf.SetY(186)
		pdf.CellFormat(0, 0, "Made in China", "0", 0, "C", false, 0, "")
	}
	if err := pdf.OutputFileAndClose(path); err != nil {
		panic(err.Error())
	}

}

func main() {
	var path string
	fmt.Println("请输入CSV文件路径（或者拖入文件）：")
	fmt.Scanln(&path)

	file, err := os.Open(path)
	paths := filepath.Dir(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	decoder := mahonia.NewDecoder("GBK")
	csvreder := csv.NewReader(decoder.NewReader(file))

	ss, _ := csvreder.ReadAll()
	sz := len(ss)

	for i := 0; i < sz; i++ {
		if len(ss[i]) == 3 {
			int, err := strconv.Atoi(ss[i][2])
			if err != nil {
				fmt.Println("ss[i][1]", "：未创建，请填入正确的打印数量", err)
				continue
			}
			creat_bar(ss[i][1], ss[i][0], filepath.Join(paths, ss[i][0]+".pdf"), int)
		} else if len(ss[i]) == 2 {
			creat_bar(ss[i][1], ss[i][0], filepath.Join(paths, ss[i][0]+".pdf"), 1)
		} else {
			fmt.Println("您的文件有问题，请重新创建。")
		}

	}
	fmt.Println("输入任意字符退出！")
	fmt.Scanln()
}
