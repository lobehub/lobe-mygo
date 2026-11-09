package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/action"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

type tabEvent int

const (
	tabNone tabEvent = iota
	tabActivated
	tabClosed
)

func tabItem(c *ui.Context, tk *theme.Tokens, r route, active, closable bool) tabEvent {
	ev := tabNone
	title := r.label(c)
	row := ui.Row(c).Height(26).Basis(200).Grow(0).Shrink(1).MinWidth(40).MaxWidth(200).Clip().
		Radius(tk.BorderRadius).AlignItems(ui.Center)
	hovered := row.Hovered()
	switch {
	case active && tk.Dark:
		row.MinWidth(150).Background(tk.FillSecondary).InsetShadow(0, 0, 0, 1, tk.BorderSecondary)
	case active:
		tk.Elevate(row.MinWidth(150).Background(tk.BgElevated), theme.BoxShadowTertiary)
	case hovered:
		row.Background(tk.FillQuaternary)
	}
	showClose := closable && hovered
	row.Children(func() {
		b := ui.ButtonBase(c).Role(ui.RoleTab).AccessSelected(active).Label(title).
			Grow(1).MinWidth(0).AlignSelf(ui.Stretch).Row().Justify(ui.Start).AlignItems(ui.Center).Gap(6).
			Padding(0, 6, 0, 8).Cursor(ui.CursorPointer)
		b.Children(func() {
			ui.Icon(c, r.icon()).FontSize(14).TextColor(tk.TextSecondary).Shrink(0)
			text := ui.Text(c, title).FontSize(tk.FontSizeSM).TextColor(tk.Text).SingleLine()
			if showClose {
				text.Margin(0, 17, 0, 0)
			}
		})
		if b.Clicked() {
			ev = tabActivated
		}
		if !closable {
			return
		}
		x := action.ActionIcon(c, lucide.X(), action.ActionIconProps{Size: theme.Small, Title: tr(c).TitleBarCloseTab(title)}).
			Absolute().Right(3).Top(3).Size(20, 20).Radius(5)
		if !showClose {
			x.Opacity(0)
		}
		if x.Clicked() {
			ev = tabClosed
		}
	})
	return ev
}
