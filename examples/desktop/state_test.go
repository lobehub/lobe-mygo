package main

import "testing"

func TestNavigatePushesHistoryAndDropsForwardEntries(t *testing.T) {
	a := newApp()
	a.navigate(routeTasks)
	a.navigate(routeProjects)
	a.back()
	a.navigate(routeDocs)
	if got := a.current(); got != routeDocs {
		t.Fatalf("current %q, want %q", got, routeDocs)
	}
	if a.canForward() {
		t.Error("forward still possible after navigating from the middle of history")
	}
	a.back()
	if got := a.current(); got != routeTasks {
		t.Errorf("back went to %q, want %q", got, routeTasks)
	}
}

func TestNavigateToCurrentRouteKeepsHistory(t *testing.T) {
	a := newApp()
	a.navigate(routeHome)
	if a.canBack() {
		t.Error("navigating to the current route added a history entry")
	}
}

func TestBackAndForwardStopAtTheEnds(t *testing.T) {
	a := newApp()
	a.back()
	if got := a.current(); got != routeHome {
		t.Fatalf("back at start moved to %q", got)
	}
	a.navigate(routeTasks)
	a.back()
	if !a.canForward() || a.canBack() {
		t.Fatalf("canBack %v canForward %v, want false true", a.canBack(), a.canForward())
	}
	a.forward()
	a.forward()
	if got := a.current(); got != routeTasks {
		t.Errorf("forward went to %q, want %q", got, routeTasks)
	}
}

func TestEachTabKeepsItsOwnHistory(t *testing.T) {
	a := newApp()
	a.navigate(routeTasks)
	a.openTab()
	if len(a.tabs) != 2 || a.active != 1 || a.current() != routeHome {
		t.Fatalf("tabs %d active %d current %q, want a new Home tab in front", len(a.tabs), a.active, a.current())
	}
	if a.canBack() {
		t.Error("new tab inherited history")
	}
	a.activate(0)
	if got := a.current(); got != routeTasks {
		t.Errorf("first tab shows %q, want %q", got, routeTasks)
	}
}

func TestCloseTabKeepsTheActiveTabInFront(t *testing.T) {
	a := newApp()
	a.navigate(routeTasks)
	a.openTab()
	a.navigate(routeDocs)
	a.openTab()
	a.navigate(routeProjects)
	a.activate(1)
	a.closeTab(0)
	if a.active != 0 || a.current() != routeDocs {
		t.Errorf("after closing a tab before it: active %d current %q, want 0 %q", a.active, a.current(), routeDocs)
	}
	a.closeTab(0)
	if len(a.tabs) != 1 || a.current() != routeProjects {
		t.Errorf("after closing the active tab: %d tabs current %q, want 1 %q", len(a.tabs), a.current(), routeProjects)
	}
}

func TestLastTabCannotClose(t *testing.T) {
	a := newApp()
	a.closeTab(0)
	if len(a.tabs) != 1 {
		t.Errorf("%d tabs, want the last one kept", len(a.tabs))
	}
}

func TestClosingTheLastActiveTabActivatesTheOneBefore(t *testing.T) {
	a := newApp()
	a.openTab()
	a.navigate(routeDocs)
	a.closeTab(1)
	if a.active != 0 || a.current() != routeHome {
		t.Errorf("active %d current %q, want 0 %q", a.active, a.current(), routeHome)
	}
}

func TestGoToJumpsWithinHistoryWithoutAddingAnEntry(t *testing.T) {
	a := newApp()
	a.navigate(routeTasks)
	a.navigate(routeDocs)
	a.goTo(0)
	if a.current() != routeHome || !a.canForward() || len(a.tab().history) != 3 {
		t.Errorf("current %q history %v, want Home with Tasks and Docs ahead", a.current(), a.tab().history)
	}
	a.goTo(9)
	if a.current() != routeHome {
		t.Errorf("out-of-range jump moved to %q", a.current())
	}
}
