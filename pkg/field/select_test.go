package field

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/internal/testkit"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func selectTester(v *string, p SelectProps) *ui.Tester {
	p.Options = []SelectOption{{Value: "gpt", Label: "GPT-5"}, {Value: "claude", Label: "Claude"}, {Value: "gemini", Label: "Gemini", Disabled: true}}
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(20).Children(func() { Select(c, v, p).Width(240) })
	}, 400, 300)
}

func TestSelectShowsPlaceholderThenChoice(t *testing.T) {
	v := ""
	tt := selectTester(&v, SelectProps{Placeholder: "Pick a model"})
	if !tt.HasText("Pick a model") {
		t.Fatalf("texts %q, want the placeholder", tt.Texts())
	}
	if err := tt.Click("Pick a model"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !tt.HasText("Claude") {
		t.Fatalf("popup not open: %q", tt.Texts())
	}
	tt.Click("Claude")
	tt.Advance(time.Second)
	tt.Frame()
	if v != "claude" {
		t.Fatalf("value %q after choosing Claude", v)
	}
	if tt.HasText("GPT-5") || !tt.HasText("Claude") {
		t.Errorf("after choosing: %q, want the popup closed and Claude in the trigger", tt.Texts())
	}
}

func TestSelectPopupMatchesTriggerAndMarksChoice(t *testing.T) {
	tk := theme.Light()
	v := "gpt"
	tt := selectTester(&v, SelectProps{})
	tt.Click("GPT-5")
	tt.Advance(time.Second)
	tt.Frame()
	claude, ok := tt.Find("Claude")
	if !ok {
		t.Fatal("popup not open")
	}
	if claude.X < 20 || claude.X > 40 {
		t.Errorf("option at x %v, want the popup under the trigger at 20", claude.X)
	}
	check := false
	row := claude.Y - 32 + claude.H/2
	for x := float32(200); x < 260; x++ {
		for dy := float32(-3); dy <= 3; dy++ {
			if testkit.Near(testkit.PixelAt(tt, 400, x, row+dy), tk.Primary) {
				check = true
			}
		}
	}
	if !check {
		t.Error("chosen option has no Primary check at the end of its row")
	}
	tt.Click("Gemini")
	tt.Advance(time.Second)
	if v != "gpt" {
		t.Errorf("disabled option chosen: %q", v)
	}
}

func TestSelectPopupScrollsAndFitsWindow(t *testing.T) {
	options := make([]SelectOption, 30)
	for i := range options {
		lbl := fmt.Sprintf("Option %02d", i)
		options[i] = SelectOption{Value: lbl, Label: lbl}
	}
	inBox := func(r, box ui.Rect) bool {
		return r.H > 0 && r.Y >= box.Y && r.Y+r.H <= box.Y+box.H
	}

	t.Run("wheel_and_keys_in_256_popup", func(t *testing.T) {
		v := "Option 00"
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Padding(20).Children(func() {
				Select(c, &v, SelectProps{Options: options}).Width(240)
			})
		}, 400, 400)

		if err := tt.Click("Option 00"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(time.Second)
		tt.Frame()

		lb := role(t, tt, ui.RoleListBox, "")
		if lb.H != 256 {
			t.Fatalf("listbox height = %v, want 256", lb.H)
		}

		cx, cy := lb.X+lb.W/2, lb.Y+lb.H/2
		tt.Move(cx, cy)
		tt.Scroll(cx, cy, 0, 200)
		tt.Frame()
		opt08, ok := tt.Find("Option 08")
		if !ok || !inBox(opt08, lb) {
			t.Fatalf("Option 08 after wheel scroll = %+v (found %v), want inside listbox %+v", opt08, ok, lb)
		}

		tt.Key(0, ui.KeyEnd)
		tt.Frame()
		opt29, ok := tt.Find("Option 29")
		if !ok || !inBox(opt29, lb) {
			t.Fatalf("Option 29 after End = %+v (found %v), want inside listbox %+v", opt29, ok, lb)
		}

		tt.Key(0, ui.KeyHome)
		tt.Frame()
		opt00, ok := tt.AccessByRole(ui.RoleOption, "Option 00")
		if !ok || !inBox(opt00.Bounds, lb) {
			t.Fatalf("Option 00 after Home = %+v (found %v), want inside listbox %+v", opt00.Bounds, ok, lb)
		}
	})

	t.Run("opens_scrolled_to_selected", func(t *testing.T) {
		v := "Option 25"
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Padding(20).Children(func() {
				Select(c, &v, SelectProps{Options: options}).Width(240)
			})
		}, 400, 400)

		if err := tt.Click("Option 25"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(time.Second)
		tt.Frame()

		lb := role(t, tt, ui.RoleListBox, "")
		opt25, ok := tt.AccessByRole(ui.RoleOption, "Option 25")
		if !ok || !inBox(opt25.Bounds, lb) {
			t.Fatalf("selected Option 25 on open = %+v (found %v), want inside listbox %+v", opt25.Bounds, ok, lb)
		}
		opt24, ok := tt.Find("Option 24")
		if !ok || !inBox(opt24, lb) {
			t.Fatalf("Option 24 on open = %+v (found %v), want inside listbox %+v", opt24, ok, lb)
		}
		tt.ClickAt(opt24.X+opt24.W/2, opt24.Y+opt24.H/2)
		tt.Advance(time.Second)
		tt.Frame()
		if v != "Option 24" {
			t.Fatalf("selected = %q, want Option 24", v)
		}
	})

	t.Run("clamps_to_window_room", func(t *testing.T) {
		v := "Option 00"
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Padding(140, 20).Children(func() {
				Select(c, &v, SelectProps{Options: options}).Width(240)
			})
		}, 400, 320)

		if err := tt.Click("Option 00"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(time.Second)
		tt.Key(0, ui.KeyEnd)
		tt.Frame()

		lb := role(t, tt, ui.RoleListBox, "")
		if lb.Y+lb.H > 320-4 {
			t.Fatalf("listbox %+v overflows 320 DIP window", lb)
		}
		opt29, ok := tt.Find("Option 29")
		if !ok || !inBox(opt29, lb) || opt29.Y+opt29.H > 320-4 {
			t.Fatalf("Option 29 after End = %+v (found %v), want inside listbox %+v and window", opt29, ok, lb)
		}
	})
}
