package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/pkg/feedback"
	"github.com/lobehub/lobe-mygo/pkg/field"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/overlay"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

var sendGlyph = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 14 14"><path fill="currentColor" d="M.743 3.773c-.818-.555-.422-1.834.567-1.828l11.496.074a1 1 0 01.837 1.538l-6.189 9.689c-.532.833-1.822.47-1.842-.518L5.525 8.51a1 1 0 01.522-.9l1.263-.686a.808.808 0 00-.772-1.42l-1.263.686a1 1 0 01-1.039-.051L.743 3.773z"/></svg>`))

var (
	modelChoices = []overlay.MenuItem{{Key: "Acme Default", Label: "Acme Default"}, {Key: "Acme Fast", Label: "Acme Fast"}}
)

func composer(c *ui.Context) {
	tk := theme.TokensOf(c)
	card(c, tk)
}

func card(c *ui.Context, tk *theme.Tokens) {
	edge := tk.Fill
	if tk.Dark {
		edge = tk.FillSecondary
	}
	ui.Column(c).Radius(20).Border(1, edge).Background(tk.BgElevated).ZIndex(1).
		Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.03)).Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.04)).
		Children(func() {
			// The borderless text area keeps a transparent 1 DIP border; the margin cancels it so the text
			// sits at lobe-chat's 12 DIP editor padding.
			ui.Box(c).MinHeight(108).Margin(-1, -1, 0, -1).Children(func() {
				field.TextArea(c, &state.draft, field.TextAreaProps{
					Variant:     field.InputBorderless,
					Placeholder: tr(c).ComposerPlaceholder(),
					MinRows:     4,
					MaxRows:     12,
				}).Label(tr(c).ComposerPlaceholder()).Padding(0).FontSize(14).LineHeight(1.4).LetterSpacing(0.28)
			})
			ui.Row(c).Padding(4, 8, 8, 8).Gap(4).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Row(c).Gap(2).AlignItems(ui.Center).Children(func() {
					modeSelect(c, tk)
					plusButton(c, tk)
				})
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					modelChip(c, tk)
					sendButton(c, tk)
				})
			})
		})
}

func chipLook(t ui.Element, tk *theme.Tokens, open bool, height, radius float32) {
	t.Row().Height(height).PaddingX(8).Gap(6).AlignItems(ui.Center).Radius(radius).FontSize(12).
		TextColor(tk.TextSecondary).Cursor(ui.CursorPointer)
	if t.Hovered() || open {
		t.Background(tk.FillSecondary).TextColor(tk.Text)
	}
}

func modeSelect(c *ui.Context, tk *theme.Tokens) {
	open := &state.modeOpen
	m := tr(c)
	modes := []overlay.MenuItem{
		{Key: "agent", Label: m.ComposerModeAgent(), Icon: lucide.Infinity()},
		{Key: "task", Label: m.ComposerModeTask(), Icon: lucide.ListTodo()},
	}
	current := modes[0]
	if state.mode == "task" {
		current = modes[1]
	}
	chosen := overlay.DropdownMenu(c, open, overlay.DropdownMenuProps{Items: modes}, func(t ui.Element) {
		chipLook(t.Label(m.ComposerMode()), tk, *open, 32, tk.BorderRadius)
		t.Children(func() {
			ui.Icon(c, current.Icon).FontSize(14)
			ui.Text(c, current.Label)
			ui.Icon(c, lucide.ChevronDown()).FontSize(12)
		})
	})
	if chosen != "" {
		state.mode = chosen
	}
}

func plusButton(c *ui.Context, tk *theme.Tokens) {
	b := ui.ButtonBase(c).Label(tr(c).ComposerAttach()).Size(32, 32).Radius(16).Center().Margin(0, 0, 0, 6).
		TextColor(tk.TextTertiary).Cursor(ui.CursorPointer)
	if b.Hovered() {
		b.Background(tk.FillSecondary).TextColor(tk.TextSecondary)
	}
	b.Children(func() { ui.Icon(c, lucide.Plus()).FontSize(18) })
}

func modelChip(c *ui.Context, tk *theme.Tokens) {
	open := &state.modelOpen
	chosen := overlay.DropdownMenu(c, open, overlay.DropdownMenuProps{Items: modelChoices}, func(t ui.Element) {
		chipLook(t.Label(tr(c).ComposerModel()), tk, *open, 28, tk.BorderRadiusSM)
		effort := tk.TextTertiary
		if t.Hovered() || *open {
			effort = tk.TextSecondary
		}
		t.Children(func() {
			ui.Text(c, state.model).MaxWidth(200).SingleLine()
			ui.Text(c, tr(c).ComposerEffortHigh()).TextColor(effort)
			ui.Icon(c, lucide.ChevronDown()).FontSize(12)
		})
	})
	if chosen != "" {
		state.model = chosen
	}
}

func sendButton(c *ui.Context, tk *theme.Tokens) {
	empty := strings.TrimSpace(state.draft) == ""
	b := ui.ButtonBase(c).Label(tr(c).ComposerSend()).Size(32, 32).Radius(16).Center().Disabled(empty).Cursor(ui.CursorPointer)
	bg := tk.Primary
	if b.Hovered() && !empty {
		bg = tk.PrimaryHover
	}
	b.Background(bg).TextColor(tk.BgLayout)
	if empty {
		b.Opacity(0.5).Cursor(ui.CursorNotAllowed)
	}
	b.Children(func() { ui.Icon(c, sendGlyph).FontSize(13) })
	if b.Clicked() {
		feedback.ShowToast(c, feedback.ToastProps{Title: tr(c).ComposerSent(), Description: state.draft, Type: feedback.ToastSuccess})
		state.draft = ""
	}
}
