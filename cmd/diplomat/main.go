package main

import (
	"flag"
	//"log"
	//"path/filepath"
	"strings"

	"github.com/gosimple/slug"

	"github.com/qjcg/diploma"
	"github.com/qjcg/diploma/templates"
)

func main() {
	//base := flag.String("b", "logo.png", "base image (background)")
	course := flag.String("c", "", "course name")
	period := flag.String("p", "", "training period (dates)")
	students := flag.String("s", "Joe Learnery", "list of students (comma-separated)")
	instructor := flag.String("i", "Rory Q. Teachalot", "instructor's name")
	flag.Parse()

	// ensure base image uses an absolute path
	//absBase, err := filepath.Abs(*base)
	//if err != nil {
	//	log.Fatal(err)
	//}

	session := &diploma.Session{
		Course:     *course,
		Period:     *period,
		Instructor: *instructor,
		Students:   strings.Split(*students, ","),
	}

	template := templates.Default

	d := &diploma.DiplomaSet{
		Session:   *session,
		Template:  *template,
		OutputDir: "diplomas/" + slug.Make(*course),
	}
	d.ToPDF()
}
