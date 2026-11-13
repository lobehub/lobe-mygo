package field

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/internal/testkit"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func TestDatePickerOpensPicksAndClears(t *testing.T) {
	var v time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(20).Width(260).Children(func() {
			DatePicker(c, &v, DatePickerProps{AllowClear: true}).Label("Due")
		})
	}, 400, 500)
	if !tt.HasText("Select date") {
		t.Fatalf("texts %q, want the placeholder", tt.Texts())
	}
	b := role(t, tt, ui.RoleButton, "Due")
	tt.ClickAt(b.X+b.W/2, b.Y+b.H/2)
	tt.Advance(time.Second)
	tt.Frame()
	now := tt.Now()
	name := time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	d := role(t, tt, ui.RoleButton, name)
	tt.ClickAt(d.X+d.W/2, d.Y+d.H/2)
	tt.Advance(time.Second)
	tt.Frame()
	if v.Format("2006-01-02") != name {
		t.Fatalf("value %v, want %s", v, name)
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Today"); ok {
		t.Error("the calendar stays open after choosing a day")
	}
	if !tt.HasText(name) {
		t.Errorf("trigger shows %q, want %s", tt.Texts(), name)
	}
	c := role(t, tt, ui.RoleButton, "Clear")
	tt.ClickAt(c.X+c.W/2, c.Y+c.H/2)
	tt.Frame()
	if !v.IsZero() || !tt.HasText("Select date") {
		t.Errorf("after Clear: value %v texts %q", v, tt.Texts())
	}
}

func TestDatePickerClearButtonCenteredAndStates(t *testing.T) {
	tk := theme.Light()
	picked := time.Date(2026, 11, 12, 0, 0, 0, 0, time.Local)
	for _, tc := range []struct {
		name string
		size theme.Size
	}{
		{"small", theme.Small},
		{"middle", theme.Middle},
		{"large", theme.Large},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := picked
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Padding(20).Width(260).Children(func() {
					DatePicker(c, &v, DatePickerProps{AllowClear: true, Size: tc.size}).Label("Due")
				})
			}, 400, 200)
			due := role(t, tt, ui.RoleButton, "Due")
			clear := role(t, tt, ui.RoleButton, "Clear")
			fieldCenterY := due.Y + due.H/2
			if d := math.Abs(float64((clear.Y + clear.H/2) - fieldCenterY)); d > 1 {
				t.Fatalf("clear button not vertically centered: clear=%v field=%v delta=%.2f", clear, due, (clear.Y+clear.H/2)-fieldCenterY)
			}
			// Also verify the rendered icon glyph inside the Clear button is vertically centered.
			img := tt.Image()
			scale := float32(img.Bounds().Dx()) / 400
			minY, maxY := img.Bounds().Max.Y, -1
			for py := int(clear.Y * scale); py < int((clear.Y+clear.H)*scale); py++ {
				for px := int(clear.X * scale); px < int((clear.X+clear.W)*scale); px++ {
					p := img.RGBAAt(px, py)
					if !testkit.Near(ui.RGB(p.R, p.G, p.B), tk.BgContainer) {
						if py < minY {
							minY = py
						}
						if py > maxY {
							maxY = py
						}
					}
				}
			}
			if maxY < minY {
				t.Fatalf("no clear icon ink rendered inside %v", clear)
			}
			glyphCenterY := (float32(minY) + float32(maxY) + 1) / (2 * scale)
			if d := math.Abs(float64(glyphCenterY - fieldCenterY)); d > 1 {
				t.Fatalf("clear icon glyph not vertically centered: glyphCenterY=%.2f fieldCenterY=%.2f delta=%.2f", glyphCenterY, fieldCenterY, glyphCenterY-fieldCenterY)
			}
		})
	}

	for _, tc := range []struct {
		name     string
		value    time.Time
		disabled bool
	}{
		{"empty", time.Time{}, false},
		{"disabled", picked, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.value
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Padding(20).Width(260).Children(func() {
					DatePicker(c, &v, DatePickerProps{AllowClear: true, Disabled: tc.disabled}).Label("Due")
				})
			}, 400, 200)
			if _, ok := tt.AccessByRole(ui.RoleButton, "Clear"); ok {
				t.Fatalf("%s DatePicker must not render a Clear button", tc.name)
			}
		})
	}
}
