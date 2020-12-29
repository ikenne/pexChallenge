package model

import (
	"fmt"
)

type RGB [3]byte

type ImageResult struct {
	URL string
	// Colors []string
	Colors [3]RGB
	ErrMsg string
}

func (ir ImageResult) String() string {
	if len(ir.ErrMsg) > 0 {
		return fmt.Sprintf("%s,%s\n", ir.URL, ir.ErrMsg)
	}

	// colors := strings.Join(ir.Colors, ",")
	// return fmt.Sprintf("%s,%s\n", ir.URL, colors)

	c1 := fmt.Sprintf("#%.2X%.2X%.2X", ir.Colors[0][0], ir.Colors[0][1], ir.Colors[0][2])
	c2 := fmt.Sprintf("#%.2X%.2X%.2X", ir.Colors[1][0], ir.Colors[1][1], ir.Colors[1][2])
	c3 := fmt.Sprintf("#%.2X%.2X%.2X", ir.Colors[2][0], ir.Colors[2][1], ir.Colors[2][2])
	return fmt.Sprintf("%s,%s,%s,%s\n", ir.URL, c1, c2, c3)
}
