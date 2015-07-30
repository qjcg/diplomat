package templates

import (
	"github.com/qjcg/diploma"
)

var Default = &diploma.Template{
	Image: "logo.png",
	Overlay: map[string][2]float64{
		"Student":    [2]float64{4, 1},
		"Course":     [2]float64{4, 2},
		"Period":     [2]float64{4, 3},
		"Instructor": [2]float64{4, 4},
		"Image":      [2]float64{0.3, 0.3},
	},
}
