package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func placeholderPage(c *ui.Context, title, empty, hint string, icon *ui.SVG) {
	tk := theme.TokensOf(c)
	ui.Column(c).Padding(24, 32).Gap(48).Children(func() {
		ui.Text(c, title).FontSize(tk.FontSizeHeading3).FontWeight(600)
		ui.Column(c).AlignItems(ui.Center).PaddingY(64).Children(func() {
			display.Empty(c, display.EmptyProps{Icon: icon, Title: empty, Description: hint})
		})
	})
}
