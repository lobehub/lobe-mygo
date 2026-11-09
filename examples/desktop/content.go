package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/action"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func contentCard(c *ui.Context) {
	tk := theme.TokensOf(c)
	key := state.current()
	r := routeOf(key)
	inset := float32(0)
	if !state.sidebarOpen {
		inset = 8
	}
	edge := tk.Border
	if tk.Dark {
		edge = tk.BorderSecondary
	}
	ui.Column(c).Grow(1).Basis(0).MinWidth(0).Margin(0, 8, 8, inset).
		Radius(tk.BorderRadiusLG).Background(tk.BgContainer).Border(1, edge).Clip().
		Children(func() {
			ui.Row(c).Height(46).Shrink(0).Padding(8).Gap(4).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Row(c).Gap(2).AlignItems(ui.Center).Children(func() {
					if r.header != nil {
						r.header(c)
					}
				})
				ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
					action.ActionIcon(c, lucide.SlidersHorizontal(), small(tr(c).HeaderCustomize()))
				})
			})
			ui.Scroll(c.Key("page-" + key)).Grow(1).Basis(0).Children(func() { r.view(c) })
		})
}
