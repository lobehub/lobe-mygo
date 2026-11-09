package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/feedback"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

var state = newApp()

var setAppLocale = func(tag string) { mygo.App.SetLocale(tag) }

func root(c *ui.Context) {
	theme.Apply(c)
	tk := theme.TokensOf(c)
	page := ui.Column(c).Fill()
	if !c.Vibrancy() {
		page.Background(tk.BgLayout)
	}
	page.Children(func() {
		titleBar(c)
		ui.Row(c).Grow(1).Basis(0).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
			sidebar(c)
			contentCard(c)
		})
		feedback.ToastHost(c, feedback.ToastHostProps{})
	})
}

func windowOptions() mygo.WindowOptions {
	return mygo.WindowOptions{
		Title:                locales.AppName(),
		Width:                1280,
		Height:               820,
		MinWidth:             960,
		MinHeight:            600,
		TitleBarStyle:        mygo.TitleBarHidden,
		TrafficLightPosition: &mygo.Point{X: 12, Y: 12},
		Vibrancy:             mygo.VibrancySidebar,
		BackgroundColor:      "light-dark(#f8f8f8, #000000)",
		Content:              ui.View(root),
	}
}

func main() {
	mygo.App.WhenReady(func() { mygo.NewWindow(windowOptions()) })
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
