package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/qjcg/diploma"
)

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
		d := &diploma.Diploma{absBase, *course, *dates, *instructor, strings.TrimSpace(s)}
		err := diploma.CreatePDF(d)
		if err != nil {
			log.Fatal(err)
		}
	}
}
