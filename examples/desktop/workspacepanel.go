package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/motion"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

var workspaceStats = []struct {
	title text
	count string
	wide  bool
}{{locales.Messages.PanelProjects, "4", false}, {locales.Messages.PanelTasks, "86", false}, {locales.Messages.PanelDocs, "15.7K", true}}

func workspaceSwitcher(c *ui.Context, tk *theme.Tokens) {
	open := &state.switcherOpen
	t := ui.ButtonBase(c).Label(state.workspace).Row().Gap(8).AlignItems(ui.Center).Padding(2, 8, 2, 2).
		Radius(tk.BorderRadius).Cursor(ui.CursorPointer).MinWidth(32).Shrink(1)
	if t.Hovered() || *open {
		t.Background(tk.FillTertiary)
	}
	t.Children(func() {
		display.Avatar(c, state.workspace, display.AvatarProps{Size: 28, Shape: display.AvatarSquare})
		ui.Row(c).Gap(4).AlignItems(ui.Center).MinWidth(0).Children(func() {
			ui.Text(c, state.workspace).FontWeight(500).TextColor(tk.Text).SingleLine()
			ui.Icon(c, lucide.ChevronDown()).FontSize(14).TextColor(tk.TextDescription).Shrink(0)
		})
	})
	if t.Clicked() {
		*open = !*open
	}
	ui.PopoverBase(c, t, open, func(panel ui.Element) {
		panel.AttachTo(t, ui.AnchorTopStart, ui.AnchorTopStart).Margin(-14, 0, 0, 2).AttachedOrigin().
			Role(ui.RoleDialog).Label(tr(c).PanelLabel()).MinWidth(300).Radius(tk.BorderRadius).
			Background(tk.BgElevated).TextColor(tk.Text).Transition(motion.Popup())
		tk.Elevate(panel, theme.BoxShadowSecondary)
		panel.Children(func() {
			ui.Column(c).Gap(2).Padding(0, 0, 4, 0).Children(func() { workspacePanel(c, tk, open) })
		})
	})
}

func workspacePanel(c *ui.Context, tk *theme.Tokens, open *bool) {
	ui.Row(c).Padding(12).Gap(10).AlignItems(ui.Center).Children(func() {
		display.Avatar(c, state.prefs.Name, display.AvatarProps{Size: 36, Shape: display.AvatarSquare, Background: tk.Fill})
		ui.Column(c).Children(func() {
			ui.Text(c, state.prefs.Name).FontWeight(600).LineHeight(1.4)
			ui.Text(c, "ada@acme.dev").FontSize(12).LineHeight(1.4).TextColor(tk.TextSecondary)
		})
	})
	ui.Row(c).Gap(4).PaddingX(8).Margin(0, 0, 8, 0).Children(func() {
		for _, s := range workspaceStats {
			card := ui.Column(c.Key(s.count)).Grow(1).Basis(0).Gap(2).Padding(6, 8).Radius(tk.BorderRadius).
				Background(tk.FillTertiary)
			if s.wide {
				card.Grow(2)
			}
			card.Children(func() {
				ui.Text(c, s.count).FontSize(16).FontWeight(600).LineHeight(1.2)
				ui.Text(c, s.title(tr(c))).FontSize(12).LineHeight(1.2).TextColor(tk.TextDescription)
			})
		}
	})
	panelDivider(c, tk)
	for _, w := range workspaces {
		check := func() {}
		if w == state.workspace {
			check = func() { ui.Icon(c, lucide.Check()).FontSize(16).TextColor(tk.TextSecondary) }
		}
		avatar := func() { display.Avatar(c, w, display.AvatarProps{Size: 16, Shape: display.AvatarSquare}) }
		if panelRow(c.Key(w), tk, w, avatar, check) {
			state.workspace = w
			*open = false
		}
	}
	panelDivider(c, tk)
	icon := func(svg *ui.SVG) func() { return func() { ui.Icon(c, svg).FontSize(16) } }
	hotkey := func() { display.Hotkey(c, "mod+,", display.HotkeyProps{}) }
	if panelRow(c.Key("settings"), tk, tr(c).RouteSettings(), icon(lucide.Settings2()), hotkey) {
		state.navigate(routeSettings)
		*open = false
	}
	if panelRow(c.Key("docs"), tk, tr(c).RouteDocs(), icon(lucide.BookOpen()), nil) {
		state.navigate(routeDocs)
		*open = false
	}
}

func panelDivider(c *ui.Context, tk *theme.Tokens) {
	ui.Box(c).Height(1).Margin(4, 0).Background(tk.BorderSecondary)
}

func panelRow(c *ui.Context, tk *theme.Tokens, label string, lead, extra func()) bool {
	row := ui.ButtonBase(c).Label(label).Row().MinHeight(32).Padding(6, 12).Margin(0, 4).Gap(12).
		AlignItems(ui.Center).Justify(ui.Start).Radius(tk.BorderRadius).FontSize(14).TextColor(tk.Text).
		Cursor(ui.CursorPointer)
	if row.Hovered() {
		row.Background(tk.FillTertiary)
	}
	row.Children(func() {
		ui.Row(c).Size(16, 16).Shrink(0).Center().Children(lead)
		ui.Text(c, label).Grow(1).SingleLine()
		if extra != nil {
			ui.Row(c).Shrink(0).FontSize(12).TextColor(tk.TextDescription).Children(extra)
		}
	})
	return row.Clicked()
}
