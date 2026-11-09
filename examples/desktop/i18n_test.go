package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestUIFollowsTheWindowLanguage(t *testing.T) {
	tt := newTester(t)
	tt.SetLocale("zh-Hans")
	settle(tt)
	for _, q := range []struct {
		role ui.Role
		name string
	}{
		{ui.RoleButton, "首页"},
		{ui.RoleTab, "首页"},
		{ui.RoleButton, "新建标签页"},
		{ui.RoleButton, "切换侧边栏"},
		{ui.RoleTextField, "有什么想问的？"},
		{ui.RoleTab, "我的"},
	} {
		if _, ok := tt.AccessByRole(q.role, q.name); !ok {
			t.Errorf("no %q in Chinese; texts %q", q.name, tt.Texts())
		}
	}
	state.navigate(routeSettings)
	settle(tt)
	for _, s := range []string{"通用", "个人资料", "显示名称"} {
		if !tt.HasText(s) {
			t.Errorf("settings lack %q; texts %q", s, tt.Texts())
		}
	}
}

func TestLanguageSettingSwitchesTheAppLanguage(t *testing.T) {
	var asked []string
	setAppLocale = func(tag string) { asked = append(asked, tag) }
	tt := newTester(t)
	state.navigate(routeSettings)
	settle(tt)
	state.settings.Values.Language = "zh-Hans"
	settle(tt)
	settle(tt)
	if len(asked) == 0 || asked[len(asked)-1] != "zh-Hans" {
		t.Errorf("app asked for %v, want zh-Hans", asked)
	}
	state.settings.Values.Language = ""
	settle(tt)
	settle(tt)
	if asked[len(asked)-1] != "" {
		t.Errorf("app asked for %v, want the system language back", asked)
	}
}
