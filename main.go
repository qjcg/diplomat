package main

import (
	"flag"
	"fmt"
	//"image"
	//"image/png"
	//"io/ioutil"
	"strings"

	//"code.google.com/p/freetype-go/freetype"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	instructor := flag.String("i", "", "instructor's name")
	students := flag.String("s", "", "list of students (comma-separated)")
	dates := flag.String("d", "", "training dates")
	base := flag.String("b", "logo.png", "base image (background)")
	flag.Parse()

	txt := strings.Join([]string{*instructor, *students, *dates}, "\n")

	// Build image
	//bytes, err := ioutil.ReadFile(*base)
	//if err != nil {
	//	return err
	//}
	//img, err := png.Decode(bytes)
	//if err != nil {
	//	return err
	//}

	//ctx := freetype.NewContext()
	//ctx.SetSrc(img)
	// FIXME
	//ctx.SetDst()

	// Build PDF
	pdf := gofpdf.New("L", "in", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Times", "", 10)
	pdf.Image(*base, 0, 0, 0, 0, false, "", 0, "")
	pdf.Text(2, 2, txt)
	err := pdf.OutputFileAndClose("test.pdf")
	if err != nil {
		fmt.Println("ERROR!")
	}
}
