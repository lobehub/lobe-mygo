package main

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type metric struct {
	role       ui.Role
	name       string
	x, y, w, h float32
}

const any = -1

func checkMetrics(t *testing.T, tt *ui.Tester, ms []metric) {
	t.Helper()
	for _, m := range ms {
		a, ok := tt.AccessByRole(m.role, m.name)
		if !ok {
			t.Errorf("no %q", m.name)
			continue
		}
		b := a.Bounds
		for _, f := range []struct {
			label     string
			got, want float32
		}{{"x", b.X, m.x}, {"y", b.Y, m.y}, {"w", b.W, m.w}, {"h", b.H, m.h}} {
			if f.want != any && (f.got-f.want > 0.5 || f.want-f.got > 0.5) {
				t.Errorf("%s %s = %v, want %v (lobe-chat)", m.name, f.label, f.got, f.want)
			}
		}
	}
}

func desktopTester(t *testing.T) *ui.Tester {
	tt := newTester(t)
	tt.SetTitleBar(ui.TitleBar{Height: 38, Left: 72})
	state.openTab()
	state.navigate(routeDocs)
	state.activate(0)
	settle(tt)
	return tt
}

func TestTitleBarMatchesLobeChat(t *testing.T) {
	checkMetrics(t, desktopTester(t), []metric{
		{ui.RoleButton, "Toggle sidebar", 80, 7, 24, 24},
		{ui.RoleButton, "Back", any, 7, 24, 24},
		{ui.RoleMenuButton, "History", 280 - 8 - 24, 7, 24, 24},
		{ui.RoleTab, "Home", 280, 6, any, 26},
		{ui.RoleButton, "New tab", any, 8, 26, 22},
	})
	o := windowOptions()
	if p := o.TrafficLightPosition; p == nil || p.X != 12 || p.Y != 12 {
		t.Errorf("traffic lights at %v, want {12 12} for a 38 DIP bar", p)
	}
}

func TestTabsMatchLobeChat(t *testing.T) {
	tt := desktopTester(t)
	home, _ := tt.AccessByRole(ui.RoleTab, "Home")
	docs, _ := tt.AccessByRole(ui.RoleTab, "Docs")
	if gap := docs.Bounds.X - (home.Bounds.X + home.Bounds.W); gap < 1.5 || gap > 2.5 {
		t.Errorf("gap between tabs %v, want 2", gap)
	}
	if docs.Bounds.W > 200.5 || home.Bounds.W < 149.5 {
		t.Errorf("tab widths active %v inactive %v, want active ≥150 and all ≤200", home.Bounds.W, docs.Bounds.W)
	}
}

func TestSidebarMatchesLobeChat(t *testing.T) {
	checkMetrics(t, desktopTester(t), []metric{
		{ui.RoleButton, "Acme", 6, 38 + 8, any, 32},
		{ui.RoleButton, "Notifications", any, 38 + 10, 28, 28},
		{ui.RoleButton, "Search", 4, any, 272, 36},
		{ui.RoleButton, "Help", 8, any, 28.8, 28.8},
	})
	tt := desktopTester(t)
	a, _ := tt.AccessByRole(ui.RoleButton, "Tasks")
	b, _ := tt.AccessByRole(ui.RoleButton, "Projects")
	if gap := b.Bounds.Y - (a.Bounds.Y + a.Bounds.H); gap < 0.5 || gap > 1.5 {
		t.Errorf("gap between nav rows %v, want 1", gap)
	}
}

func TestContentCardMatchesLobeChat(t *testing.T) {
	tt := desktopTester(t)
	checkMetrics(t, tt, []metric{
		{ui.RoleButton, "Customize", any, 38 + 1 + 11, 24, 24},
		{ui.RoleTab, "Mine", any, any, any, 30},
		{ui.RoleButton, "Send", any, any, 32, 32},
	})
	if px := tt.Image().RGBAAt(1280-4, 400); px.R < 240 || px.R > 250 {
		t.Errorf("pixel in the 8 DIP gap right of the card %+v, want the layout background", px)
	}
}

func TestWindowAsksForTheSidebarMaterial(t *testing.T) {
	if v := windowOptions().Vibrancy; v != mygo.VibrancySidebar {
		t.Errorf("vibrancy %q, want the sidebar material as lobe-chat", v)
	}
}

func TestComposerMatchesLobeChat(t *testing.T) {
	checkMetrics(t, desktopTester(t), []metric{
		{ui.RoleMenuButton, "Mode", any, any, any, 32},
		{ui.RoleButton, "Attach", any, any, 32, 32},
		{ui.RoleMenuButton, "Model", any, any, any, 28},
		{ui.RoleButton, "Send", any, any, 32, 32},
	})
}

func TestSettingsPageMatchesLobeChat(t *testing.T) {
	tt := desktopTester(t)
	state.navigate(routeSettings)
	settle(tt)
	name, _ := tt.AccessByRole(ui.RoleTextField, "Display name")
	if name.Bounds.W < 239.5 {
		t.Errorf("name field %v wide, want at least 240 (max(34%%, 240))", name.Bounds.W)
	}
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func TestSettingsStartsRightUnderTheHeaderRow(t *testing.T) {
	tt := desktopTester(t)
	state.navigate(routeSettings)
	settle(tt)
	title, ok := tt.Find("General")
	if !ok || abs(title.Y+title.H/2-(38+1+23)) > 1.5 {
		t.Errorf("page title %+v, want it in the header row centered at %v as lobe-chat's NavHeader", title, 38+1+23)
	}
	profile, _ := tt.Find("Profile")
	if top := profile.Y - (38 + 1 + 46); top > 24+16+4 {
		t.Errorf("first group title %v below the header row, want the form to start after the 24 DIP container padding", top)
	}
}

func TestTabCloseButtonCentersOnTheTab(t *testing.T) {
	tt := desktopTester(t)
	home, _ := tt.AccessByRole(ui.RoleTab, "Home")
	x, ok := tt.AccessByRole(ui.RoleButton, "Close Home")
	if !ok {
		t.Fatal("no close button")
	}
	tab := home.Bounds.Y + home.Bounds.H/2
	if d := x.Bounds.Y + x.Bounds.H/2 - tab; abs(d) > 0.5 {
		t.Errorf("close button centered %v off the tab's center", d)
	}
	if right := home.Bounds.X + home.Bounds.W; abs(x.Bounds.X+x.Bounds.W-(right-3)) > 0.5 {
		t.Errorf("close button ends at %v, want 3 DIPs inside the tab's end", x.Bounds.X+x.Bounds.W)
	}
}

func TestComposerTextLinesUpWithTheActionBar(t *testing.T) {
	tt := desktopTester(t)
	mode, _ := tt.AccessByRole(ui.RoleMenuButton, "Mode")
	field, ok := tt.AccessByRole(ui.RoleTextField, "Ask anything")
	if !ok {
		t.Fatal("no composer field")
	}
	card := mode.Bounds.X - 1 - 8
	if d := field.Bounds.X - card; abs(d-(1+12)) > 0.5 {
		t.Errorf("text starts %v into the card, want 13 (border and lobe-chat's 12 DIP editor padding)", d)
	}
}
