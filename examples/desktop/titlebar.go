package main

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/action"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/overlay"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

const (
	titleBarHeight  = 38
	titleBarPadding = 12
	trafficLightEnd = 80
	navGroupMinimum = 116
)

func small(title string) action.ActionIconProps {
	return action.ActionIconProps{Size: theme.Small, Title: title}
}

func titleBar(c *ui.Context) {
	tk := theme.TokensOf(c)
	bar := c.TitleBar()
	height := float32(titleBarHeight)
	if bar.Height > 0 {
		height = bar.Height
	}
	inset := float32(0)
	if bar.Left > 0 {
		inset = max(bar.Left+8, trafficLightEnd) - titleBarPadding
	}
	ui.Row(c).Height(height).Shrink(0).AlignItems(ui.Center).
		Padding(0, bar.Right+titleBarPadding, 0, titleBarPadding).DragWindow().Children(func() {
		nav := ui.Row(c).Shrink(0).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Padding(0, 8, 0, inset)
		if state.sidebarOpen {
			nav.Width(max(state.sidebarWidth-titleBarPadding, navGroupMinimum+inset))
		} else {
			nav.Width(navGroupMinimum + inset)
		}
		nav.Children(func() { titleBarNav(c) })
		tabStrip(c, tk)
		ui.Row(c).Shrink(0).Gap(8).AlignItems(ui.Center).Children(func() {
			action.ActionIcon(c, lucide.HardDrive(), small(tr(c).TitleBarDevice()))
			action.ActionIcon(c, lucide.Wifi(), small(tr(c).TitleBarOnline()))
		})
	})
}

func titleBarNav(c *ui.Context) {
	icon := lucide.PanelLeftClose()
	if !state.sidebarOpen {
		icon = lucide.PanelLeftOpen()
	}
	if action.ActionIcon(c, icon, small(tr(c).TitleBarToggleSidebar())).Clicked() {
		state.sidebarOpen = !state.sidebarOpen
	}
	ui.Row(c).Gap(2).AlignItems(ui.Center).Children(func() {
		if action.ActionIcon(c, lucide.ArrowLeft(), small(tr(c).TitleBarBack())).Disabled(!state.canBack()).Clicked() {
			state.back()
		}
		if action.ActionIcon(c, lucide.ArrowRight(), small(tr(c).TitleBarForward())).Disabled(!state.canForward()).Clicked() {
			state.forward()
		}
		historyMenu(c)
	})
}

func historyMenu(c *ui.Context) {
	open := &state.historyOpen
	t := state.tab()
	items := make([]overlay.MenuItem, 0, len(t.history))
	for i := len(t.history) - 1; i >= 0; i-- {
		r := routeOf(t.history[i])
		items = append(items, overlay.MenuItem{Key: strconv.Itoa(i), Label: r.label(c), Icon: r.icon(), Disabled: i == t.at})
	}
	chosen := overlay.DropdownMenu(c, open, overlay.DropdownMenuProps{Items: items}, func(b ui.Element) {
		action.StyleActionIcon(c, b, lucide.Clock(), action.ActionIconProps{Size: theme.Small, Title: tr(c).TitleBarHistory(), Active: *open})
	})
	if i, err := strconv.Atoi(chosen); err == nil {
		state.goTo(i)
	}
}

func tabStrip(c *ui.Context, tk *theme.Tokens) {
	activate, closing := -1, -1
	ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(2).AlignItems(ui.Center).Children(func() {
		for i, t := range state.tabs {
			switch tabItem(c.Key(t), tk, routeOf(t.route()), i == state.active, len(state.tabs) > 1) {
			case tabActivated:
				activate = i
			case tabClosed:
				closing = i
			}
		}
		plus := action.ActionIcon(c, lucide.Plus(), action.ActionIconProps{Size: theme.Small, Title: tr(c).TitleBarNewTab()}).
			Size(26, 22).Radius(tk.BorderRadiusSM).Margin(0, 0, 0, 4)
		if plus.Clicked() {
			state.openTab()
		}
	})
	if activate >= 0 {
		state.activate(activate)
	}
	if closing >= 0 {
		state.closeTab(closing)
	}
}
