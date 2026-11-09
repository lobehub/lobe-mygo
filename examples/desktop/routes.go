package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
)

type text func(m locales.Messages) string

type route struct {
	title  text
	icon   func() *ui.SVG
	view   func(c *ui.Context)
	header func(c *ui.Context)
	empty  [2]text
}

func tr(c *ui.Context) locales.Messages { return locales.In(c.Locale()) }

func (r route) label(c *ui.Context) string { return r.title(tr(c)) }

const pagePrefix = "page:"

func pageRoute(name string) string { return pagePrefix + name }

var routes = map[string]route{
	routeHome: {title: locales.Messages.RouteHome, icon: lucide.House, view: homePage, header: homeTabs},
	routeSearch: {title: locales.Messages.RouteSearch, icon: lucide.Search,
		empty: [2]text{locales.Messages.EmptySearchTitle, locales.Messages.EmptySearchHint}},
	routeTasks: {title: locales.Messages.RouteTasks, icon: lucide.ListTodo,
		empty: [2]text{locales.Messages.EmptyTasksTitle, locales.Messages.EmptyTasksHint}},
	routeProjects: {title: locales.Messages.RouteProjects, icon: lucide.FolderKanban,
		empty: [2]text{locales.Messages.EmptyProjectsTitle, locales.Messages.EmptyProjectsHint}},
	routeDocs: {title: locales.Messages.RouteDocs, icon: lucide.LibraryBig,
		empty: [2]text{locales.Messages.EmptyDocsTitle, locales.Messages.EmptyDocsHint}},
	routeNotifications: {title: locales.Messages.RouteNotifications, icon: lucide.Bell,
		empty: [2]text{locales.Messages.EmptyNotificationsTitle, locales.Messages.EmptyNotificationsHint}},
	routeSettings: {title: locales.Messages.RouteSettings, icon: lucide.Settings, view: settingsPage, header: settingsTitle},
}

func routeOf(key string) route {
	if name, ok := strings.CutPrefix(key, pagePrefix); ok {
		return route{
			title: func(locales.Messages) string { return name },
			icon:  lucide.Hash,
			view: func(c *ui.Context) {
				m := tr(c)
				placeholderPage(c, name, m.EmptyPageTitle(name), m.EmptyPageHint(name), lucide.Hash())
			},
		}
	}
	r, ok := routes[key]
	if !ok {
		return routeOf(routeHome)
	}
	if r.view == nil {
		r.view = func(c *ui.Context) {
			m := tr(c)
			placeholderPage(c, r.title(m), r.empty[0](m), r.empty[1](m), r.icon())
		}
	}
	return r
}
