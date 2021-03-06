package diploma

import (
	"log"
	"os"
	"testing"

	"github.com/gosimple/slug"
)

var sessionTests = []Session{
	Session{
		Course:     "Fun with JavaScript",
		Period:     "July 12-15 2015 (22.5 hours)",
		Instructor: "Joe Instructor",
		Recipients: []string{
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
		Image: "logo.jpg",
		Overlay: map[string][2]float64{
			"Recipient":  [2]float64{300, 200},
			"Course":     [2]float64{300, 240},
			"Period":     [2]float64{300, 260},
			"Instructor": [2]float64{300, 280},
			"Image":      [2]float64{10, 10},
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
		f, err := os.Create(d.OutputDir + "/diplomas.json")
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		d.Dump(f)
	}
}

func TestMain(m *testing.M) {
	exitStatus := m.Run()
	//os.RemoveAll("diplomas")
	os.Exit(exitStatus)
}
