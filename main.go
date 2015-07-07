package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	instructor := flag.String("i", "", "instructor's name")
	students := flag.String("s", "", "list of students (comma-separated)")
	dates := flag.String("d", "", "training dates")
	base := flag.String("b", "logo.png", "base image (background)")
	flag.Parse()

	txt := strings.Join([]string{*instructor, *students, *dates}, "\n")

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
