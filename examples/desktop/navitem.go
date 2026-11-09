package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

type navItemProps struct {
	label  string
	icon   *ui.SVG
	avatar func()
	active bool
}

func navItem(c *ui.Context, tk *theme.Tokens, p navItemProps) bool {
	row := ui.ButtonBase(c).Label(p.label).Row().Height(36).PaddingX(4).Gap(8).AlignItems(ui.Center).
		Justify(ui.Start).Radius(tk.BorderRadius).Cursor(ui.CursorPointer)
	hovered := row.Hovered()
	text, glyph := tk.TextSecondary, tk.TextDescription
	switch {
	case p.active && hovered:
		row.Background(tk.FillSecondary)
	case p.active, hovered:
		row.Background(tk.FillTertiary)
	}
	if p.active {
		text, glyph = tk.Text, tk.Text
	}
	row.Children(func() {
		ui.Row(c).Size(28, 28).Shrink(0).Center().TextColor(glyph).Children(func() {
			switch {
			case p.avatar != nil:
				p.avatar()
			case p.icon != nil:
				ui.Icon(c, p.icon).FontSize(18)
			}
		})
		ui.Text(c, p.label).TextColor(text).SingleLine()
	})
	return row.Clicked()
}
