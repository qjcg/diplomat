package diploma

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/signintech/gopdf"
)

const (
	DroidSansPath string = "/usr/share/fonts/TTF/DroidSans.ttf"
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

// Dump JSON config to an io.Writer.
func (d *DiplomaSet) Dump(w io.Writer) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s\n", data)
}

// Load reads config from JSON file, populating a DiplomaSet.
func (d *DiplomaSet) Load(configFile string) {
	return
}

// Render DiplomaSet to PDF files.
// FIXME: Make this DRY by writing a utility function
func (d *DiplomaSet) ToPDF() {
	// Create OutputDir for PDFs
	os.MkdirAll(d.OutputDir, 0700)

	for _, s := range d.Students {
		pdf := gopdf.GoPdf{}
		// Letter: 612x792
		// See https://www.gnu.org/software/gv/manual/html_node/Paper-Keywords-and-paper-size-in-points.html
		pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 792, H: 612}})
		pdf.AddPage()

		err := pdf.AddTTFFont("DroidSans", DroidSansPath)
		if err != nil {
			log.Fatal(err)
		}

		pdf.Image(d.Image, d.Overlay["Image"][0], d.Overlay["Image"][1], nil)

		// Student
		err = pdf.SetFont("DroidSans", "", 26)
		if err != nil {
			log.Fatal(err)
		}
		pdf.SetX(d.Overlay["Student"][0])
		pdf.SetY(d.Overlay["Student"][1])
		pdf.Cell(nil, s)

		// Course
		err = pdf.SetFont("DroidSans", "", 13)
		if err != nil {
			log.Fatal(err)
		}
		pdf.SetX(d.Overlay["Course"][0])
		pdf.SetY(d.Overlay["Course"][1])
		pdf.Cell(nil, d.Course)

		// Period
		err = pdf.SetFont("DroidSans", "", 11)
		if err != nil {
			log.Fatal(err)
		}
		pdf.SetX(d.Overlay["Period"][0])
		pdf.SetY(d.Overlay["Period"][1])
		pdf.Cell(nil, d.Period)

		// Instructor
		pdf.SetX(d.Overlay["Instructor"][0])
		pdf.SetY(d.Overlay["Instructor"][1])
		pdf.Cell(nil, d.Instructor)

		// FIXME: do this via a separate function for greater testability?
		pdfPath := filepath.Join(d.OutputDir, slug.Make(s)+".pdf")
		pdf.WritePdf(pdfPath)
	}
}
