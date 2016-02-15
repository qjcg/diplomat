package templates

import (
	"github.com/qjcg/diploma"
)

var Default = &diploma.Template{
	Image: "logo.png",
	Overlay: map[string][2]float64{
		"Recipient":  [2]float64{4, 200},
		"Course":     [2]float64{4, 240},
		"Period":     [2]float64{4, 260},
		"Instructor": [2]float64{4, 280},
		"Image":      [2]float64{10, 10},
	},
}
