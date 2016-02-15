package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gosimple/slug"

	"github.com/qjcg/diploma"
	"github.com/qjcg/diploma/templates"
)

func main() {
	course := flag.String("c", "", "course name")
	period := flag.String("p", "", "training period (dates)")
	students := flag.String("s", "Joe Learnery", "list of students (comma-separated)")
	instructor := flag.String("i", "Rory Q. Teachalot", "instructor's name")

	// TODO: add template flag handling logic
	//template := flag.String("t", "default", "template")
	flag.Parse()

	// TODO: use cli template option if provided
	template := templates.Default

	session := &diploma.Session{
		Course:     *course,
		Period:     *period,
		Instructor: *instructor,
		Students:   strings.Split(*students, ","),
	}

	d := &diploma.DiplomaSet{
		Session:   *session,
		Template:  *template,
		OutputDir: "diplomas/" + slug.Make(*course),
	}
	d.ToPDF()

	// create JSON config file
	conf, err := os.Create(d.OutputDir + "/diplomas.json")
	defer conf.Close()
	if err != nil {
		log.Fatal(err)
	}
	d.Dump(conf)
}
