package main

import (
	"fmt"
	"io"
	"strings"

	svg "github.com/ajstarks/svgo/float"
)

type Action map[string]interface{}

type HaliteCanvas struct {
	*svg.SVG
}

func NewHaliteCanvas(w io.Writer) *HaliteCanvas {
	return &HaliteCanvas{
		SVG: svg.New(w),
	}
}

var emptyclose = "/>\n"

func (canvas HaliteCanvas) Planet(x float64, y float64, r float64, ownerID string, s ...string) {
	d := canvas.Decimals
	fmt.Fprintf(canvas.Writer, `<circle class="planet %s" cx="%.*f" cy="%.*f" r="%.*f" %s`, ownerID, d, x, d, y, d, r, endstyle(s, emptyclose))
}
func (canvas HaliteCanvas) Entity(x float64, y float64, r float64, class []string, s ...string) {
	d := canvas.Decimals
	fmt.Fprintf(canvas.Writer, `<circle class="planet %s" cx="%.*f" cy="%.*f" r="%.*f" %s`, strings.Join(class, " "), d, x, d, y, d, r, endstyle(s, emptyclose))
}

// endstyle modifies an SVG object, with either a series of name="value" pairs,
// or a single string containing a style
func endstyle(s []string, endtag string) string {
	if len(s) > 0 {
		nv := ""
		for i := 0; i < len(s); i++ {
			if strings.Index(s[i], "=") > 0 {
				nv += (s[i]) + " "
			} else {
				nv += style(s[i]) + " "
			}
		}
		return nv + endtag
	}
	return endtag

}
func style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}
