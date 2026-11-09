package main

import (
	"time"

	"github.com/lobehub/lobe-mygo/pkg/form"
)

const (
	routeHome          = "home"
	routeSearch        = "search"
	routeTasks         = "tasks"
	routeProjects      = "projects"
	routeDocs          = "docs"
	routeNotifications = "notifications"
	routeSettings      = "settings"
)

type tab struct {
	history []string
	at      int
}

func (t *tab) route() string { return t.history[t.at] }

type settings struct {
	Name     string
	Language string
	Notify   bool
}

var defaultSettings = settings{Name: "Ada", Notify: true}

type app struct {
	tabs         []*tab
	active       int
	sidebarWidth float32
	sidebarOpen  bool
	workspace    string
	switcherOpen bool
	historyOpen  bool
	mode         string
	modeOpen     bool
	model        string
	modelOpen    bool
	folded       map[string]bool
	homeTab      int
	draft        string
	settings     *form.Form[settings]
	prefs        settings
	savedAt      time.Time
}

func newApp() *app {
	return &app{tabs: []*tab{newTab()}, sidebarWidth: 280, sidebarOpen: true, workspace: "Acme", mode: "agent", model: "Acme Default", folded: map[string]bool{},
		settings: form.New(defaultSettings),
		prefs:    defaultSettings,
	}
}

func newTab() *tab { return &tab{history: []string{routeHome}} }

func (a *app) tab() *tab { return a.tabs[a.active] }

func (a *app) current() string { return a.tab().route() }

func (a *app) navigate(route string) {
	t := a.tab()
	if t.route() == route {
		return
	}
	t.history = append(t.history[:t.at+1], route)
	t.at++
}

func (a *app) canBack() bool { return a.tab().at > 0 }

func (a *app) canForward() bool { t := a.tab(); return t.at < len(t.history)-1 }

func (a *app) back() {
	if a.canBack() {
		a.tab().at--
	}
}

func (a *app) forward() {
	if a.canForward() {
		a.tab().at++
	}
}

func (a *app) goTo(i int) {
	if t := a.tab(); i >= 0 && i < len(t.history) {
		t.at = i
	}
}

func (a *app) openTab() {
	a.tabs = append(a.tabs, newTab())
	a.active = len(a.tabs) - 1
}

func (a *app) activate(i int) {
	if i >= 0 && i < len(a.tabs) {
		a.active = i
	}
}

func (a *app) closeTab(i int) {
	if len(a.tabs) <= 1 || i < 0 || i >= len(a.tabs) {
		return
	}
	a.tabs = append(a.tabs[:i], a.tabs[i+1:]...)
	if a.active > i || a.active == len(a.tabs) {
		a.active--
	}
}
