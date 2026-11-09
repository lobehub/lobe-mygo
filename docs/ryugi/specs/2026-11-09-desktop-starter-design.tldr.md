# TL;DR: desktop starter app

**What gets built:** `examples/desktop/`, a runnable skeleton of a desktop app laid out like the
LobeHub desktop client: title bar with tabs, sidebar, rounded content card, status bar. Developers
copy the folder and replace the placeholder pages.

## What you'll see

1. `go run ./examples/desktop` → a window, traffic lights inset, no system title bar.
2. Click a sidebar entry → the page changes in the active tab.
3. `+` opens a tab, `×` closes one; back / forward walk that tab's history.
4. The title-bar icon collapses the sidebar; dragging its edge resizes it.
5. The gear opens Settings; the top-right icon in the card opens a right panel.
6. Switch macOS to dark mode → the whole app turns dark.

## Decisions

| Topic | Decision |
| --- | --- |
| Kind | skeleton, placeholder content (you chose) |
| Layout | the LobeHub desktop screenshot (you chose) |
| Location | `examples/desktop/`, `package main`, same module |
| Title bar | `TitleBarHiddenInset` + `DragWindow`, own tab strip |
| Sidebar | `DraggablePanel` (hides fully, like the reference) + `List` + collapsible groups |
| State | one `app` struct, change methods tested without UI |
| Pages | Home, Search, Tasks, Projects, Docs, Notifications (bell), Settings, generic placeholder |

## Check these (my calls)

- Generic labels ("Acme"), no LobeHub branding or AI-chat parts.
- Tab strip lines up with the content card's left edge and moves with the sidebar.
- Real tabs + per-tab history, not just a drawn tab strip.
- A toggleable right panel inside the card (three-pane skeleton).
- Desktop only; not on the docs site.
- No persistence; no changes to `pkg/` (gaps go to `TODO.md`).

## Time

About 2–3 hours to build and test, plus 20 minutes of screenshot checks in light and dark.
