package diploma_test

import (
	"log"
	"os"
	"testing"

	"github.com/gosimple/slug"

	"github.com/qjcg/diplomat/pkg/diploma"
)

var sessionTests = []diploma.Session{
	{
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

var templateTests = []diploma.Template{
	{
		Image: "logo.jpg",
		Overlay: map[string][2]float64{
			"Recipient":  {300, 200},
			"Course":     {300, 240},
			"Period":     {300, 260},
			"Instructor": {300, 280},
			"Image":      {10, 10},
		},
	},
}

var diplomaSetTests = []*diploma.DiplomaSet{
	{
		Session:   sessionTests[0],
		Template:  templateTests[0],
		OutputDir: "./diplomas/" + slug.Make(sessionTests[0].Course),
	},
}

func TestToPDF(t *testing.T) {
	fontData, err := os.ReadFile("../../cmd/diplomat/fonts/DroidSans.ttf")
	if err != nil {
		t.Fatalf("Error loading font: %v", err)
	}

	for _, d := range diplomaSetTests {
		d.ToPDF("DroidSans", fontData)
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
