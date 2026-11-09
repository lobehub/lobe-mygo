package main

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/display"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

type recentItem struct {
	title, preview, date, owner string
}

var recentItems = []recentItem{
	{"Quarterly planning notes", "Goals for Q1 are drafted; the hiring plan still needs numbers from finance.", "Nov 7", "Anna"},
	{"Onboarding checklist", "Accounts, laptop setup and the first-week reading list for new members.", "Nov 7", "Ben"},
	{"Release 0.2 retro", "What went well: shorter review cycles. What to change: earlier QA sign-off.", "Nov 5", "Chloe"},
	{"Design review: settings page", "Group related switches, move danger actions to the bottom.", "Nov 4", "Anna"},
	{"Customer interview #12", "Wants keyboard shortcuts for every list action, and an offline mode.", "Nov 3", "Ben"},
}

var teamItems = []recentItem{
	{"Team handbook", "How we plan, review and ship. Edited by three people this week.", "Nov 8", "Chloe"},
	{"Roadmap 2027", "Draft themes: performance, collaboration, mobile companion app.", "Nov 6", "Anna"},
	{"Incident log", "No open incidents. Last one resolved on Oct 30.", "Nov 2", "Ben"},
}

func dayPeriod(hour int) string {
	switch {
	case hour < 5:
		return "evening"
	case hour < 12:
		return "morning"
	case hour < 18:
		return "afternoon"
	}
	return "evening"
}

var homeViews = []text{locales.Messages.HomeMine, locales.Messages.HomeTeam}

func homeTabs(c *ui.Context) {
	tk := theme.TokensOf(c)
	ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		for i, view := range homeViews {
			label := view(tr(c))
			active := state.homeTab == i
			b := ui.ButtonBase(c.Key(i)).Role(ui.RoleTab).AccessSelected(active).Label(label).
				Padding(4, 8).Cursor(ui.CursorPointer)
			color, weight := tk.TextTertiary, 400
			switch {
			case active:
				color, weight = tk.Text, 600
			case b.Hovered():
				color = tk.TextSecondary
			}
			b.Children(func() {
				ui.Text(c, label).FontSize(14).LineHeight(22.0 / 14).FontWeight(weight).TextColor(color)
				if active {
					ui.Box(c).Absolute().Bottom(-2).Left(0).Right(0).Children(func() {
						ui.Box(c).Size(4, 4).Radius(2).Background(tk.Text).AlignSelf(ui.Center)
					}).Column().AlignItems(ui.Center)
				}
			})
			if b.Clicked() {
				state.homeTab = i
			}
		}
	})
}

func homePage(c *ui.Context) {
	m := tr(c)
	title, items := m.HomeRecent(), recentItems
	if state.homeTab == 1 {
		title, items = m.HomeShared(), teamItems
	}
	ui.Column(c).AlignItems(ui.Center).Padding(50, 24, 24, 24).Children(func() {
		ui.Column(c).Width(760).MaxWidthPercent(100).Children(func() {
			ui.Text(c, m.HomeGreeting(dayPeriod(c.Now().Hour()), state.prefs.Name)).
				FontSize(22).LineHeight(1.4).LetterSpacing(-0.22).FontWeight(600).Margin(0, 0, 24, 0)
			composer(c)
			ui.Column(c).Margin(24, 0, 0, 0).Gap(32).Children(func() {
				itemList(c, title, items)
			})
		})
	})
}

func itemList(c *ui.Context, title string, items []recentItem) {
	tk := theme.TokensOf(c)
	ui.Column(c).Gap(12).Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).FontSize(12).LineHeight(18.0 / 12).FontWeight(600).LetterSpacing(0.48).TextColor(tk.TextSecondary)
			ui.Row(c).MinWidth(20).Height(18).PaddingX(5).Radius(5).Background(tk.FillTertiary).Center().Children(func() {
				ui.Text(c, strconv.Itoa(len(items))).FontSize(12).LineHeight(1).FontWeight(500).
					FontFeatures("tnum").TextColor(tk.TextTertiary)
			})
		})
		ui.Column(c).Gap(4).Children(func() {
			for i, it := range items {
				itemRow(c.Key(i), tk, it)
			}
		})
	})
}

func itemRow(c *ui.Context, tk *theme.Tokens, it recentItem) {
	row := ui.ButtonBase(c).Label(it.title).Row().Gap(12).Padding(9, 10).Margin(0, -10).Radius(tk.BorderRadiusLG).
		AlignItems(ui.Start).Justify(ui.Start).Cursor(ui.CursorPointer)
	if row.Hovered() {
		row.Background(tk.FillQuaternary)
	}
	row.Children(func() {
		ui.Box(c).PaddingY(3).Shrink(0).Children(func() {
			display.Avatar(c, it.owner, display.AvatarProps{Size: 22, Shape: display.AvatarSquare}).Margin(1, 0, 0, 0)
		})
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(3).Children(func() {
			ui.Text(c, it.title).FontSize(15).LineHeight(22.0 / 15).FontWeight(500).TextColor(tk.Text).SingleLine()
			ui.Text(c, it.preview).FontSize(13).LineHeight(20.0 / 13).TextColor(tk.TextTertiary).SingleLine()
		})
		ui.Text(c, it.date).FontSize(12).LineHeight(18.0 / 12).TextColor(tk.TextTertiary).MinWidth(56).
			TextAlign(ui.End).Shrink(0)
	})
	if row.Clicked() {
		state.navigate(pageRoute(it.title))
	}
}
