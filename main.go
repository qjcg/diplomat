package diploma

import (
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

// Create an initial PNG image, with text added.
// TODO: CreatePDF should take the byteslice returned by this function as input (?).
func CreatePNG(d *Diploma) ([]byte, error) {
	return []byte{}, nil
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
