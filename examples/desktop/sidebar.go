package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/action"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/layout"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

var navRoutes = []string{routeSearch, routeHome, routeTasks, routeProjects, routeDocs}

var workspaces = []string{"Acme", "Personal"}

type sidebarGroup struct {
	name    string
	label   text
	members []string
	people  bool
}

var sidebarGroups = []sidebarGroup{
	{name: "private", label: locales.Messages.GroupPrivate, members: []string{"Design", "Notes"}},
	{name: "workspace", label: locales.Messages.GroupWorkspace, members: []string{"Anna", "Ben", "Chloe"}, people: true},
}

func sidebar(c *ui.Context) {
	layout.DraggablePanel(c, &state.sidebarWidth, &state.sidebarOpen, layout.DraggablePanelProps{
		Placement: layout.Left,
		MinSize:   240,
		MaxSize:   400,
		NoToggle:  true,
		NoBorder:  true,
	}, func() {
		tk := theme.TokensOf(c)
		ui.Column(c).Fill().TextColor(tk.TextSecondary).Children(func() {
			sidebarHeader(c, tk)
			ui.Scroll(c).Grow(1).Basis(0).Children(func() {
				ui.Column(c).PaddingX(4).Gap(8).Children(func() {
					ui.Column(c).Gap(1).Children(func() {
						for _, key := range navRoutes {
							r := routeOf(key)
							if navItem(c.Key(key), tk, navItemProps{label: r.label(c), icon: r.icon(), active: state.current() == key}) {
								state.navigate(key)
							}
						}
					})
					for _, g := range sidebarGroups {
						groupSection(c.Key(g.name), tk, g)
					}
				})
			})
			sidebarFooter(c)
		})
	})
}

func sidebarHeader(c *ui.Context, tk *theme.Tokens) {
	ui.Row(c).Shrink(0).Padding(8, 6).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
		workspaceSwitcher(c, tk)
		ui.Row(c).Gap(2).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Children(func() {
				bell := action.ActionIcon(c, lucide.Bell(), action.ActionIconProps{BlockSize: 28, GlyphSize: 16, Title: tr(c).RouteNotifications()})
				if bell.Clicked() {
					state.navigate(routeNotifications)
				}
				ui.Box(c).Absolute().Top(2).Right(2).Size(8, 8).Radius(4).Background(tk.Error).IgnoreInput(true)
			})
		})
	})
}

var caret = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M7 3.5v17a1 1 0 0 0 1.5.86l14-8.5a1 1 0 0 0 0-1.72l-14-8.5A1 1 0 0 0 7 3.5z" fill="currentColor"/></svg>`))

func groupSection(c *ui.Context, tk *theme.Tokens, g sidebarGroup) {
	folded := state.folded[g.name]
	label := g.label(tr(c))
	ui.Column(c).Gap(1).Children(func() {
		head := ui.ButtonBase(c).Label(label).Row().Gap(8).AlignItems(ui.Center).Justify(ui.Start).
			Padding(4, 4, 4, 8).Radius(tk.BorderRadius).Cursor(ui.CursorPointer).
			FontSize(tk.FontSizeSM).FontWeight(500).TextColor(tk.TextDescription)
		if head.Hovered() {
			head.Background(tk.FillTertiary)
		}
		head.Children(func() {
			ui.Text(c, label)
			turn := float32(90)
			if folded {
				turn = 0
			}
			ui.Row(c).Size(18, 18).Margin(0, 0, 0, -6).Center().Children(func() {
				ui.Icon(c, caret).FontSize(7).Rotate(turn)
			})
		})
		if head.Clicked() {
			state.folded[g.name] = !folded
		}
		if folded {
			return
		}
		for _, m := range g.members {
			p := navItemProps{label: m, icon: lucide.Hash(), active: state.current() == pageRoute(m)}
			if g.people {
				p.icon = nil
				p.avatar = func() { display.Avatar(c, m, display.AvatarProps{Size: 22, Shape: display.AvatarSquare}) }
			}
			if navItem(c.Key(m), tk, p) {
				state.navigate(pageRoute(m))
			}
		}
	})
}

func sidebarFooter(c *ui.Context) {
	ui.Row(c).Shrink(0).Padding(8).Gap(2).AlignItems(ui.Center).Children(func() {
		if action.ActionIcon(c, lucide.CircleHelp(), action.ActionIconProps{GlyphSize: 16, Title: tr(c).SidebarHelp()}).Clicked() {
			state.navigate(routeDocs)
		}
	})
}
