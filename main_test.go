package diploma

import (
	"os"
	"testing"

	"github.com/gosimple/slug"
)

var sessionTests = []Session{
	Session{
		Course:     "Fun with JavaScript",
		Period:     "July 12-15 2015 (22.5 hours)",
		Instructor: "Joe Instructor",
		Students: []string{
			"Joe Student",
			"Jenny Student",
			"Jean Étudiant",
			"Jasmine Estudianté",
			"Jerry de VeryLongNameThatKeepsGoing",
		},
	},
}

var templateTests = []Template{
	Template{
		Image: "logo.png",
		Overlay: map[string][2]float64{
			"Student":    [2]float64{4, 1},
			"Course":     [2]float64{4, 2},
			"Period":     [2]float64{4, 3},
			"Instructor": [2]float64{4, 4},
			"Image":      [2]float64{0.3, 0.3},
		},
	},
}

var diplomaSetTests = []*DiplomaSet{
	&DiplomaSet{
		Session:   sessionTests[0],
		Template:  templateTests[0],
		OutputDir: "./diplomas/" + slug.Make(sessionTests[0].Course),
	},
}

func TestToPDF(t *testing.T) {
	for _, d := range diplomaSetTests {
		d.ToPDF()
	}
}

func TestDump(t *testing.T) {
	for _, d := range diplomaSetTests {
		d.Dump()
	}
}

func TestMain(m *testing.M) {
	exitStatus := m.Run()
	os.RemoveAll("diplomas")
	os.Exit(exitStatus)
}
