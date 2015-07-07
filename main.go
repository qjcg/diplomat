package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jung-kurt/gofpdf"
)

type Diploma struct {
	BackgroundImage string
	Course          string
	Dates           string
	Instructor      string
	Student         string
}

// DiplomaText is a string to be displayed on a Diploma, specifying coordinates
// and font information.
type DiplomaText struct {
	X        float64
	Y        float64
	FontName string
	FontSize int
	Text     string
}

// Read config from JSON file.
func ReadConfig(fileStr string) error {
	return nil
}

// Write config to JSON file in diplomas directory.
func WriteConfig(fileStr string) error {
	return nil
}

// Create a Diploma PDF and write it to disk.
// FIXME: UTF-8 characters don't display properly, since cp1252 encoding is used
// 		See http://godoc.org/github.com/jung-kurt/gofpdf#Fpdf.SetFont
func CreatePDF(d *Diploma) error {
	pdf := gofpdf.New("L", "in", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Times", "", 30)
	pdf.Image(d.BackgroundImage, 0.3, 0.3, 0, 0, false, "", 0, "")
	pdf.Text(4, 1, d.Course)
	pdf.Text(4, 2, d.Dates)
	pdf.Text(4, 3, d.Instructor)
	pdf.SetFont("Times", "", 50)
	pdf.Text(4, 4, d.Student)

	err := pdf.OutputFileAndClose(slug.Make(d.Student) + ".pdf")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	base := flag.String("b", "logo.png", "base image (background)")
	course := flag.String("c", "Web Programming for Bureaucrats", "course name")
	dates := flag.String("d", "July 7-10, 2015 (22.5 hours)", "training dates")
	students := flag.String("s", "Joe Student, Jenny Student, Français Gérèçêêëà", "list of students (comma-separated)")
	instructor := flag.String("i", "John Gosset", "instructor's name")
	flag.Parse()

	// ensure base image uses an absolute path
	absBase, err := filepath.Abs(*base)
	if err != nil {
		log.Fatal(err)
	}

	// create PDF output dir and cd there
	pdfDir := "diplomas/" + slug.Make(*course)
	os.MkdirAll(pdfDir, 0755)
	os.Chdir(pdfDir)

	studentsSlice := strings.Split(*students, ",")
	for _, s := range studentsSlice {
		d := &Diploma{absBase, *course, *dates, *instructor, strings.TrimSpace(s)}
		err := CreatePDF(d)
		if err != nil {
			log.Fatal(err)
		}
	}
}
