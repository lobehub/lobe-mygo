package hotkey

import (
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func Sizes(c *ui.Context) {
	ui.Column(c).Gap(16).Children(func() {
		for _, size := range []theme.Size{theme.Small, theme.Middle} {
			ui.Row(c.Key(size)).Gap(12).AlignItems(ui.Center).Children(func() {
				display.Hotkey(c, "mod+k", display.HotkeyProps{Size: size})
				display.Hotkey(c, "mod+k", display.HotkeyProps{Size: size, Compact: true})
			})
		}
	})
}
