package field

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/internal/testkit"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func searchTester(q *string, p SearchBarProps) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(20).Children(func() {
			SearchBar(c, q, p).Label("search")
			ui.Button(c, "Elsewhere")
		})
	}, 400, 200)
}

func TestSearchBarShortKeyTagAndFocus(t *testing.T) {
	q := ""
	tt := searchTester(&q, SearchBarProps{ShortKey: "mod+k", Platform: 2})
	if !tt.HasText("Ctrl") || !tt.HasText("K") {
		t.Fatalf("texts %q, want the Ctrl K tag", tt.Texts())
	}
	tt.Key(ui.Ctrl, ui.KeyK)
	tt.Frame()
	if !tt.Focused("search") {
		t.Fatal("Ctrl+K did not focus the search")
	}
	if tt.HasText("Ctrl") {
		t.Error("tag still shows while the search has the focus")
	}
	tt.Type("lobe")
	tt.Frame()
	if q != "lobe" {
		t.Errorf("query %q", q)
	}
}

func TestSearchBarEnterSearchesAndClearEmpties(t *testing.T) {
	q := ""
	var searched []string
	tt := searchTester(&q, SearchBarProps{OnSearch: func(s string) { searched = append(searched, s) }})
	if err := tt.Click("search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("agents")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(searched) != 1 || searched[0] != "agents" {
		t.Fatalf("searched %q, want [agents]", searched)
	}
	info, ok := tt.AccessByRole(ui.RoleButton, "Clear")
	if !ok {
		t.Fatal("no clear button with a query")
	}
	tt.ClickAt(info.Bounds.X+info.Bounds.W/2, info.Bounds.Y+info.Bounds.H/2)
	tt.Frame()
	if q != "" {
		t.Errorf("query %q after clearing", q)
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Clear"); ok {
		t.Error("clear button shows with an empty query")
	}
}

func TestSmallSearchBarHotkeyHasVerticalSpacing(t *testing.T) {
	tk := theme.Light()
	chipBg := testkit.Over(tk.BgContainer, tk.FillTertiary)

	for _, tc := range []struct {
		name     string
		platform display.HotkeyPlatform
		mods     ui.Modifiers
		wantText string
	}{
		{"apple", display.HotkeyApple, ui.Super, "K"},
		{"other", display.HotkeyOther, ui.Ctrl, "Ctrl"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := ""
			tt := searchTester(&q, SearchBarProps{ShortKey: "mod+k", Size: theme.Small, Platform: tc.platform})
			if !tt.HasText(tc.wantText) || !tt.HasText("K") {
				t.Fatalf("texts %q, want %q and K", tt.Texts(), tc.wantText)
			}
			k, ok := tt.Find("K")
			if !ok {
				t.Fatal("missing hotkey text K")
			}
			cy := k.Y + k.H/2
			// 10 DIPs above/below the vertical center is inside the 24 DIP bar (±12) but outside a 16 DIP chip (±8).
			for _, dy := range []float32{-10, 10} {
				got := testkit.PixelAt(tt, 400, k.X, cy+dy)
				if testkit.Near(got, chipBg) || !testkit.Near(got, tk.BgContainer) {
					t.Fatalf("at dy=%v: pixel=%+v, want BgContainer %+v (outside chipBg %+v)", dy, got, tk.BgContainer, chipBg)
				}
			}

			// Verify shortcut focus hides the badge and typing + Clear keeps working.
			tt.Key(tc.mods, ui.KeyK)
			tt.Frame()
			if !tt.Focused("search") {
				t.Fatalf("shortcut did not focus small search bar")
			}
			if tt.HasText("K") {
				t.Errorf("hotkey badge still visible while small search bar is focused: %q", tt.Texts())
			}
			tt.Type("aion")
			tt.Frame()
			if q != "aion" {
				t.Fatalf("query=%q, want aion", q)
			}
			clear := role(t, tt, ui.RoleButton, "Clear")
			tt.ClickAt(clear.X+clear.W/2, clear.Y+clear.H/2)
			tt.Frame()
			if q != "" {
				t.Fatalf("query=%q after Clear, want empty", q)
			}
		})
	}
}
