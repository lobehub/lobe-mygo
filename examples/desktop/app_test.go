package main

import (
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func newTester(t *testing.T) *ui.Tester {
	t.Helper()
	state = newApp()
	return ui.NewTester(root, 1280, 820)
}

func settle(tt *ui.Tester) {
	tt.Advance(time.Second)
	tt.Frame()
}

func TestContentShowsTheActiveRoute(t *testing.T) {
	tt := newTester(t)
	if !tt.HasText("Home") {
		t.Errorf("texts %q, want the Home page", tt.Texts())
	}
	state.navigate(routeProjects)
	settle(tt)
	if !tt.HasText("No projects yet") {
		t.Errorf("texts %q, want the Projects placeholder", tt.Texts())
	}
}

func TestGroupRowRoutesShowTheirName(t *testing.T) {
	tt := newTester(t)
	state.navigate(pageRoute("Design"))
	settle(tt)
	if !tt.HasText("Design") || !tt.HasText("Nothing in Design yet") {
		t.Errorf("texts %q, want the Design placeholder", tt.Texts())
	}
}

func TestNoStatusBarAndTheCardKeepsItsBottomInset(t *testing.T) {
	tt := newTester(t)
	if tt.HasText("Local Database") || tt.HasText("Ready") {
		t.Errorf("texts %q, want no status bar", tt.Texts())
	}
	if px := tt.Image().RGBAAt(700, 820-4); px.R > 250 {
		t.Errorf("pixel 4 DIPs above the bottom edge %+v, want the layout background under the card", px)
	}
	if px := tt.Image().RGBAAt(700, 820-12); px.R < 254 {
		t.Errorf("pixel 12 DIPs above the bottom edge %+v, want the card", px)
	}
}

func press(t *testing.T, tt *ui.Tester, role ui.Role, name string) {
	t.Helper()
	b, ok := tt.AccessByRole(role, name)
	if !ok {
		t.Fatalf("no %q (%v); texts %q", name, role, tt.Texts())
	}
	tt.ClickAt(b.Bounds.X+b.Bounds.W/2, b.Bounds.Y+b.Bounds.H/2)
	settle(tt)
}

func TestNewTabOpensHomeInFront(t *testing.T) {
	tt := newTester(t)
	state.navigate(routeDocs)
	settle(tt)
	press(t, tt, ui.RoleButton, "New tab")
	if len(state.tabs) != 2 || state.current() != routeHome {
		t.Fatalf("%d tabs current %q, want a second Home tab", len(state.tabs), state.current())
	}
	tab, ok := tt.AccessByRole(ui.RoleTab, "Home")
	if !ok || !tab.Selected {
		t.Errorf("Home tab %+v, want it shown and selected", tab)
	}
	press(t, tt, ui.RoleTab, "Docs")
	if state.current() != routeDocs {
		t.Errorf("clicking the Docs tab shows %q", state.current())
	}
}

func TestCloseTabButtonClosesThatTab(t *testing.T) {
	tt := newTester(t)
	state.openTab()
	state.navigate(routeTasks)
	settle(tt)
	press(t, tt, ui.RoleButton, "Close Tasks")
	if len(state.tabs) != 1 || state.current() != routeHome {
		t.Errorf("%d tabs current %q, want only Home left", len(state.tabs), state.current())
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Close Home"); ok {
		t.Error("the last tab offers a close button")
	}
}

func TestHistoryButtonsWalkTheActiveTab(t *testing.T) {
	tt := newTester(t)
	if b, _ := tt.AccessByRole(ui.RoleButton, "Back"); !b.Disabled {
		t.Error("Back enabled with no history")
	}
	state.navigate(routeTasks)
	settle(tt)
	press(t, tt, ui.RoleButton, "Back")
	if state.current() != routeHome {
		t.Fatalf("Back went to %q", state.current())
	}
	press(t, tt, ui.RoleButton, "Forward")
	if state.current() != routeTasks {
		t.Errorf("Forward went to %q", state.current())
	}
}

func TestTitleBarKeepsClearOfTheWindowControls(t *testing.T) {
	tt := newTester(t)
	tt.SetTitleBar(ui.TitleBar{Height: 40, Left: 80})
	settle(tt)
	b, ok := tt.AccessByRole(ui.RoleButton, "Toggle sidebar")
	if !ok || b.Bounds.X < 80 {
		t.Errorf("toggle %+v, want it right of the 80 DIP window controls", b.Bounds)
	}
}

func TestSidebarNavOpensItsPage(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleButton, "Projects")
	if state.current() != routeProjects {
		t.Errorf("current %q, want %q", state.current(), routeProjects)
	}
	press(t, tt, ui.RoleButton, "Anna")
	if state.current() != pageRoute("Anna") {
		t.Errorf("current %q, want Anna's page", state.current())
	}
}

func TestSidebarGroupHeaderFoldsItsRows(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleButton, "Private")
	if _, ok := tt.AccessByRole(ui.RoleButton, "Design"); ok {
		t.Error("Design row still shown with Private folded")
	}
	press(t, tt, ui.RoleButton, "Private")
	if _, ok := tt.AccessByRole(ui.RoleButton, "Design"); !ok {
		t.Error("Design row missing with Private unfolded")
	}
}

func TestSidebarToggleHidesTheSidebar(t *testing.T) {
	tt := newTester(t)
	before, _ := tt.AccessByRole(ui.RoleTab, "Home")
	press(t, tt, ui.RoleButton, "Toggle sidebar")
	if b, ok := tt.AccessByRole(ui.RoleButton, "Projects"); ok && b.Bounds.X+b.Bounds.W > 0 {
		t.Errorf("Projects row at %+v, want the sidebar out of view", b.Bounds)
	}
	after, _ := tt.AccessByRole(ui.RoleTab, "Home")
	if after.Bounds.X >= before.Bounds.X {
		t.Errorf("tab strip at x %v, want it moved left of %v", after.Bounds.X, before.Bounds.X)
	}
}

func TestSidebarBellOpensNotifications(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleButton, "Notifications")
	if state.current() != routeNotifications {
		t.Errorf("bell opened %q", state.current())
	}
}

func openWorkspacePanel(t *testing.T, tt *ui.Tester) {
	t.Helper()
	press(t, tt, ui.RoleButton, "Acme")
	if !tt.HasText("ada@acme.dev") {
		t.Fatalf("texts %q, want the workspace panel open", tt.Texts())
	}
}

func TestWorkspacePanelOpensOverTheSwitcher(t *testing.T) {
	tt := newTester(t)
	openWorkspacePanel(t, tt)
	panel, ok := tt.AccessByRole(ui.RoleDialog, "Workspace panel")
	if !ok {
		t.Fatal("no workspace panel")
	}
	if b := panel.Bounds; b.X != 8 || b.Y != 32 || b.W < 300 {
		t.Errorf("panel at %+v, want top-left at 8,32 and at least 300 wide as lobe-chat", b)
	}
}

func TestWorkspacePanelSwitchesTheWorkspace(t *testing.T) {
	tt := newTester(t)
	openWorkspacePanel(t, tt)
	press(t, tt, ui.RoleButton, "Personal")
	if state.workspace != "Personal" || tt.HasText("ada@acme.dev") {
		t.Errorf("workspace %q, want Personal with the panel closed", state.workspace)
	}
}

func TestWorkspacePanelOpensSettings(t *testing.T) {
	tt := newTester(t)
	openWorkspacePanel(t, tt)
	press(t, tt, ui.RoleButton, "Settings")
	if state.current() != routeSettings {
		t.Errorf("Settings opened %q", state.current())
	}
}

func TestHomeShowsTheComposerAndRecentItems(t *testing.T) {
	tt := newTester(t)
	if _, ok := tt.AccessByRole(ui.RoleTextField, "Ask anything"); !ok {
		t.Error("no composer on Home")
	}
	if !tt.HasText("Recent") || !tt.HasText(recentItems[0].title) {
		t.Errorf("texts %q, want the recent list", tt.Texts())
	}
}

func TestHomeTeamTabShowsTeamItems(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleTab, "Team")
	if !tt.HasText(teamItems[0].title) || tt.HasText(recentItems[0].title) {
		t.Errorf("texts %q, want the team list instead of mine", tt.Texts())
	}
}

func TestThereIsNoRightPanel(t *testing.T) {
	tt := newTester(t)
	if _, ok := tt.AccessByRole(ui.RoleButton, "Toggle panel"); ok {
		t.Error("a right panel toggle is still shown")
	}
	if tt.HasText("Details") {
		t.Error("the right panel is still built")
	}
}

func TestSettingsSwitchChangesTheSetting(t *testing.T) {
	tt := newTester(t)
	state.navigate(routeSettings)
	settle(tt)
	if !tt.HasText("Display name") {
		t.Fatalf("texts %q, want the settings form", tt.Texts())
	}
	was := state.settings.Values.Notify
	press(t, tt, ui.RoleSwitch, "Desktop notifications")
	if state.settings.Values.Notify == was {
		t.Error("switch did not change the setting")
	}
}

func TestActiveTabIsDrawnAsALightCard(t *testing.T) {
	tt := newTester(t)
	tab, ok := tt.AccessByRole(ui.RoleTab, "Home")
	if !ok {
		t.Fatal("no Home tab")
	}
	px := tt.Image().RGBAAt(int(tab.Bounds.X+tab.Bounds.W-4), int(tab.Bounds.Y+tab.Bounds.H/2))
	if px.R < 230 || px.G < 230 || px.B < 230 {
		t.Errorf("active tab pixel %+v, want the light container color", px)
	}
}

func TestCardKeepsAGapFromTheWindowEdgeWithTheSidebarHidden(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleButton, "Toggle sidebar")
	tab, _ := tt.AccessByRole(ui.RoleTab, "Home")
	if tab.Bounds.X < 8 {
		t.Fatalf("tab at %+v", tab.Bounds)
	}
	if px := tt.Image().RGBAAt(4, 400); px.R > 250 && px.G > 250 && px.B > 250 {
		t.Errorf("pixel at x 4 is %+v, want the window background beside the card", px)
	}
}

func TestManyTabsShrinkEvenly(t *testing.T) {
	tt := newTester(t)
	state.navigate(routeNotifications)
	for range 14 {
		state.openTab()
	}
	settle(tt)
	var widths []float32
	for _, a := range tt.AccessAll() {
		if a.Role == ui.RoleTab && a.Name != "Mine" && a.Name != "Team" {
			widths = append(widths, a.Bounds.W)
		}
	}
	if len(widths) != 15 {
		t.Fatalf("%d tabs, want 15", len(widths))
	}
	inactive := widths[:14]
	for _, w := range inactive {
		if d := w - inactive[0]; d > 1 || d < -1 {
			t.Errorf("tab widths %v, want the inactive tabs as wide as each other", widths)
			break
		}
	}
	if active := widths[14]; active < 149.5 {
		t.Errorf("active tab %v wide, want at least 150 as lobe-chat keeps it", active)
	}
}

func TestHoveringATabKeepsTheStripInPlace(t *testing.T) {
	tt := newTester(t)
	state.openTab()
	state.navigate(routeDocs)
	settle(tt)
	before, _ := tt.AccessByRole(ui.RoleButton, "New tab")
	home, _ := tt.AccessByRole(ui.RoleTab, "Home")
	tt.Move(home.Bounds.X+home.Bounds.W/2, home.Bounds.Y+home.Bounds.H/2)
	settle(tt)
	after, _ := tt.AccessByRole(ui.RoleButton, "New tab")
	if after.Bounds.X != before.Bounds.X {
		t.Errorf("New tab moved from x %v to %v on hover", before.Bounds.X, after.Bounds.X)
	}
}

func TestSettingsSaveAsYouEdit(t *testing.T) {
	tt := newTester(t)
	state.navigate(routeSettings)
	settle(tt)
	if _, ok := tt.AccessByRole(ui.RoleButton, "Submit"); ok {
		t.Error("a Submit button is shown, want auto-save as lobe-chat")
	}
	state.settings.Values.Name = "Bo"
	settle(tt)
	settle(tt)
	if state.prefs.Name != "Bo" || !tt.HasText("Saved") {
		t.Errorf("prefs name %q, want Bo saved without a submit", state.prefs.Name)
	}
	state.settings.Values.Name = ""
	settle(tt)
	settle(tt)
	if state.prefs.Name != "Bo" {
		t.Errorf("prefs name %q, want a blank name left unsaved", state.prefs.Name)
	}
}

func TestTabLabelsStartAtTheTabsLeftEdge(t *testing.T) {
	tt := newTester(t)
	state.openTab()
	state.navigate(pageRoute("Roadmap"))
	settle(tt)
	tab, _ := tt.AccessByRole(ui.RoleTab, "Roadmap")
	img := tt.Image()
	y := int(tab.Bounds.Y + tab.Bounds.H/2)
	first := -1
	for x := int(tab.Bounds.X); x < int(tab.Bounds.X+tab.Bounds.W); x++ {
		if px := img.RGBAAt(x, y); px.R < 160 {
			first = x
			break
		}
	}
	if first < 0 || float32(first)-tab.Bounds.X > 30 {
		t.Errorf("first glyph pixel at x %d in a tab at x %v, want the label at its start", first, tab.Bounds.X)
	}
}

func TestWindowPinsTheTrafficLightsToTheTitleBarRow(t *testing.T) {
	o := windowOptions()
	if o.TitleBarStyle != mygo.TitleBarHidden || o.TrafficLightPosition == nil {
		t.Fatalf("title bar style %q lights %v, want hidden with pinned traffic lights", o.TitleBarStyle, o.TrafficLightPosition)
	}
}

func TestSecondRowSharesOneCenterLine(t *testing.T) {
	tt := newTester(t)
	tt.SetTitleBar(ui.TitleBar{Height: 40, Left: 80})
	settle(tt)
	acme, _ := tt.AccessByRole(ui.RoleButton, "Acme")
	bell, _ := tt.AccessByRole(ui.RoleButton, "Notifications")
	panel, _ := tt.AccessByRole(ui.RoleButton, "Customize")
	mine, _ := tt.AccessByRole(ui.RoleTab, "Mine")
	line := center(acme)
	for name, y := range map[string]float32{"bell": center(bell), "customize": center(panel)} {
		if d := y - line; d > 0.5 || d < -0.5 {
			t.Errorf("%s centered at %v, want %v as the workspace switcher", name, y, line)
		}
	}
	acmeText := inkCenter(tt, ui.Rect{X: acme.Bounds.X + 40, Y: acme.Bounds.Y, W: 30, H: acme.Bounds.H})
	mineText := inkCenter(tt, mine.Bounds)
	if d := mineText - acmeText; d > 1 || d < -1 {
		t.Errorf("Mine text centered at %v, Acme text at %v, want one line", mineText, acmeText)
	}
}

func TestComposerModeMenuSwitchesTheMode(t *testing.T) {
	tt := newTester(t)
	press(t, tt, ui.RoleMenuButton, "Mode")
	press(t, tt, ui.RoleMenuItem, "Task")
	if state.mode != "task" || !tt.HasText("Task") {
		t.Errorf("mode %q, want Task shown", state.mode)
	}
}

func TestComposerHasNoBanner(t *testing.T) {
	tt := newTester(t)
	if tt.HasText("Templates") {
		t.Error("the promotion strip is still under the composer")
	}
}
