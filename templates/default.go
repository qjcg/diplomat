package templates

import (
	"github.com/qjcg/diploma"
)

var Default = &diploma.Template{
	Image: "logo.jpg",
	Overlay: map[string][2]float64{
		"Recipient":  [2]float64{300, 200},
		"Course":     [2]float64{300, 240},
		"Period":     [2]float64{300, 260},
		"Instructor": [2]float64{300, 280},
		"Image":      [2]float64{10, 10},
	},
}
