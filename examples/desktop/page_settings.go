package main

import (
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/lobehub/lobe-mygo/examples/desktop/locales"
	"github.com/lobehub/lobe-mygo/pkg/choice"
	"github.com/lobehub/lobe-mygo/pkg/field"
	"github.com/lobehub/lobe-mygo/pkg/form"
	"github.com/lobehub/lobe-mygo/pkg/lucide"
	"github.com/lobehub/lobe-mygo/pkg/theme"
)

func languages(m locales.Messages) []field.SelectOption {
	options := []field.SelectOption{{Value: "", Label: m.SettingsLanguageSystem()}}
	for _, tag := range []string{"en", "zh-Hans"} {
		options = append(options, field.SelectOption{Value: tag, Label: i18n.Autonym(tag)})
	}
	return options
}

func settingsPage(c *ui.Context) {
	tk := theme.TokensOf(c)
	m := tr(c)
	f := state.settings
	ui.Column(c).AlignItems(ui.Center).Padding(24, 24, 128, 24).Children(func() {
		ui.Column(c).Width(1024).MaxWidthPercent(100).Children(func() {
			form.View(c, f, form.Props[settings]{
				Variant:             form.Filled,
				ItemMinWidth:        240,
				ItemMinWidthPercent: 34,
				Gap:                 36,
				ChangeDebounce:      300 * time.Millisecond,
				OnValuesChange:      saveSettings,
			}, func() {
				form.Group(c, f, form.GroupProps{Title: m.SettingsProfile(), Fixed: true, Extra: autoSaveHint(c, tk)}, func() {
					form.Field(c, f, form.FieldProps{Name: "name", Label: m.SettingsDisplayName(), Desc: m.SettingsDisplayNameDesc(), Required: true, Value: f.Values.Name}, func() {
						fill(c, func() {
							field.Input(c, &f.Values.Name, field.InputProps{Invalid: f.Error("name") != ""}).Label(m.SettingsDisplayName())
						})
					})
					form.Field(c, f, form.FieldProps{Name: "language", Label: m.SettingsLanguage(), Desc: m.SettingsLanguageDesc(), Value: f.Values.Language}, func() {
						fill(c, func() { field.Select(c, &f.Values.Language, field.SelectProps{Options: languages(m)}) })
					})
				})
				form.Group(c, f, form.GroupProps{Title: m.SettingsNotifications(), Fixed: true, Extra: autoSaveHint(c, tk)}, func() {
					form.Field(c, f, form.FieldProps{Name: "notify", Label: m.SettingsDesktopNotifications(), Desc: m.SettingsDesktopNotificationsDesc(), Value: f.Values.Notify}, func() {
						choice.Switch(c, &f.Values.Notify, choice.SwitchProps{}).Label(m.SettingsDesktopNotifications())
					})
				})
			})
		})
	})
}

// field.Input and field.Select return their inner element, so the column is what stretches the
// field to the width of the form's control column.
func fill(c *ui.Context, control func()) {
	ui.Column(c).AlignSelf(ui.Stretch).AlignItems(ui.Stretch).Children(control)
}

func saveSettings(s settings) {
	if s.Language != state.prefs.Language {
		setAppLocale(s.Language)
		state.prefs.Language = s.Language
	}
	if s.Name == "" {
		return
	}
	state.prefs = s
	state.savedAt = time.Now()
}

func settingsTitle(c *ui.Context) {
	ui.Text(c, tr(c).SettingsGeneral()).FontWeight(500).PaddingX(4)
}

func autoSaveHint(c *ui.Context, tk *theme.Tokens) func() {
	return func() {
		if state.savedAt.IsZero() {
			return
		}
		ui.Row(c).Gap(4).AlignItems(ui.Center).FontSize(12).TextColor(tk.TextDescription).Children(func() {
			ui.Icon(c, lucide.CircleCheck()).FontSize(12)
			ui.Text(c, tr(c).SettingsSaved())
		})
	}
}
