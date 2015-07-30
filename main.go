package diploma

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gosimple/slug"
	"github.com/jung-kurt/gofpdf"
)

// A Session represents a training session.
type Session struct {
	Course     string
	Period     string
	Instructor string
	Students   []string
}

// A Template conains an image file path along with a map of text overlay coordinates.
type Template struct {
	Image   string
	Overlay map[string][2]float64
}

// A DiplomaSet contains an OutputDir for PDFs, and embedded Template and Session structs.
type DiplomaSet struct {
	Session
	Template
	OutputDir string
}

// Dump config to JSON file in diplomas directory.
func (d *DiplomaSet) Dump() {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("%v", data)
}

// Load reads config from JSON file, populating a DiplomaSet.
func (d *DiplomaSet) Load(configFile string) {
	return
}

// Render DiplomaSet to PDF files.
// FIXME: UTF-8 characters don't display properly, since cp1252 encoding is used
// 		See http://godoc.org/github.com/jung-kurt/gofpdf#Fpdf.SetFont
func (d *DiplomaSet) ToPDF() {
	// Create OutputDir for PDFs
	os.MkdirAll(d.OutputDir, 0700)

	for _, s := range d.Students {
		pdf := gofpdf.New("L", "in", "Letter", "")
		pdf.AddPage()

		pdf.SetFontLocation(".")
		pdf.AddFont("DroidSans", "", "DroidSans.json")

		pdf.Image(d.Image, d.Overlay["Image"][0], d.Overlay["Image"][1], 0, 0, false, "", 0, "")

		pdf.SetFont("DroidSans", "", 50)
		pdf.Text(d.Overlay["Student"][0], d.Overlay["Student"][1], s)
		pdf.SetFont("DroidSans", "", 30)
		pdf.Text(d.Overlay["Course"][0], d.Overlay["Course"][1], d.Course)
		pdf.SetFont("DroidSans", "", 10)
		pdf.Text(d.Overlay["Period"][0], d.Overlay["Period"][1], d.Period)
		pdf.Text(d.Overlay["Instructor"][0], d.Overlay["Instructor"][1], d.Instructor)

		err := pdf.OutputFileAndClose(d.OutputDir + "/" + slug.Make(s) + ".pdf")
		if err != nil {
			log.Fatal(err)
		}
	}
}
