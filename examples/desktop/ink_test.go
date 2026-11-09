package main

import "github.com/egoist/mygo/ui"

func inkRows(tt *ui.Tester, r ui.Rect) (top, bottom int) {
	img := tt.Image()
	top, bottom = -1, -1
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		dark := false
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); int(p.R)+int(p.G)+int(p.B) < 3*150 {
				dark = true
				break
			}
		}
		if dark {
			if top < 0 {
				top = y
			}
			bottom = y
		} else if top >= 0 {
			return
		}
	}
	return
}

func center(a ui.AccessInfo) float32 { return a.Bounds.Y + a.Bounds.H/2 }

func inkCenter(tt *ui.Tester, r ui.Rect) float32 {
	top, bottom := inkRows(tt, r)
	return float32(top+bottom) / 2
}
