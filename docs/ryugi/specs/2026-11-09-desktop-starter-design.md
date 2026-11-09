# Desktop starter app: design

## Intent

lobe-mygo has 70+ components and a catalog that shows them one at a time, but nothing shows how they
fit into a whole desktop app. A developer starting a MyGo app today has to work out the window chrome,
the sidebar, the content area and the state layout alone.

The starter is a **skeleton**: a runnable desktop app with the standard layout of a modern desktop
client (the LobeHub desktop app is the reference, see "Reference" below). Developers copy the
directory, rename it, delete the placeholder pages and fill in their own. It also shows lobe-mygo
working as a complete app.

Success: `go run ./examples/desktop` opens a window that looks like the reference layout in light and
dark mode. A developer can find where to add a page, a sidebar entry or a status-bar item within a
minute of opening the code.

## Reference

`docs/ryugi/specs/2026-11-09-desktop-starter-reference.png`: the user's screenshot of the LobeHub
desktop app (macOS, light mode). The layout regions to copy:

1. **Title bar** (about 40 DIPs, window background, drags the window): macOS traffic lights inset, a
   sidebar toggle, back / forward / history icons, then a browser-style tab strip (icon + title per
   tab, active tab raised on a white card, `+` to open a tab), and a couple of small status icons at
   the right end. The toggle and history icons sit in the sidebar's column; the tab strip starts at
   the content card's left edge and follows it as the sidebar resizes. With the sidebar hidden, the
   tab strip starts right after the history icons (my call).
2. **Sidebar** (about 300 DIPs, window background, no card): a workspace switcher row (avatar, name,
   chevron) with a bell icon carrying a red dot; a nav list (Search, Home, Tasks, Projects, Docs);
   collapsible groups with a small gray header and a caret ("Private", "Workspace"), each holding
   avatar rows; at the bottom, help and settings icons.
3. **Content card**: a rounded white card (`BgContainer`, large radius, hairline border) inset from
   the window edge on top, right and bottom, holding the page. The Home page has a "Mine / Team" tab
   switch at the top left, panel-toggle icons at the top right, a greeting heading, a large input
   card, and a "Recent topics" list (avatar, title, one-line preview, date at the right).
4. **Status bar** (about 24 DIPs, window background, small monospace-ish gray text) across the whole
   bottom: items at the left ("Local Database", "Replica", a path) and at the right (stats).

Nothing in the starter copies LobeHub's branding, mascot or AI-specific content. Labels are generic
("Acme", "Home", "Inbox"...).

## Goals

- One directory, `examples/desktop/`, `package main`, in the repo's own module, importing only
  `github.com/egoist/mygo`, `github.com/egoist/mygo/ui` and lobe-mygo's public `pkg/...` packages, so
  it copies out cleanly (the developer changes the module path, nothing else).
- The four regions above, built mostly from existing components:
  - sidebar: `layout.DraggablePanel` at `Left` with `NoToggle` (resize by drag; `expanded` driven
    by the title-bar toggle; collapses fully, as the reference does). Not `DraggableSideNav`: it
    stops at a 64-DIP rail (my call),
    `display.List` for nav rows, `display.Avatar`, `overlay.DropdownMenu` for the workspace switcher,
    `feedback.Badge` for the bell dot, `action.ActionIcon` for the icon buttons;
  - content: `nav.Tabs` (Point variant) for "Mine / Team", `field.TextArea` in the input card,
    `display.List` or plain rows for the recent list, `display.Empty` on placeholder pages,
    `form` for the settings page.
- Window: `mygo.TitleBarHiddenInset`, title bar row marked `DragWindow()`, left padding from
  `c.TitleBar().Left` so the traffic lights never overlap anything; on Windows / Linux the right
  padding from `c.TitleBar().Right`.
- Light and dark follow the system through `theme.Apply`, with no colors outside the tokens.
- Working skeleton behavior (this is what makes it a template, not a picture):
  - **Navigation**: clicking a sidebar entry shows that page in the active tab.
  - **Tabs**: `+` opens a new tab on Home; a tab closes with its `×` (shown on hover and on the
    active tab); the last tab cannot close; clicking a tab activates it.
  - **History**: back / forward walk the active tab's own history; the icons are disabled at the
    ends.
  - **Sidebar toggle**: the title-bar icon hides and shows the sidebar completely (with
    DraggablePanel's own slide); dragging the edge resizes it.
  - **Groups**: the sidebar group headers collapse and expand their rows.
  - **Right panel**: an icon at the content card's top right toggles a right-hand
    `DraggablePanel` (`Right`, `NoToggle`) inside the card (placeholder content), so the template
    shows a three-pane layout.
  - **Bell**: opens the Notifications placeholder page.
  - **Settings**: the gear in the sidebar footer opens a Settings page (a short `form` with a text
    field, a select and a switch; values kept in memory only).
- Pages: Home (as in the reference); Search, Tasks, Projects, Docs, Notifications as placeholders
  (title + `Empty`); Settings. The nav list is Search, Home, Tasks, Projects, Docs, as in the
  reference (my call). Sidebar group rows (e.g. "Design", "Notes" under Private; three people under Workspace)
  open a generic placeholder page titled with the row's name.
- Tests with `ui.NewTester` covering the behaviors above (navigate, open / close tab, back / forward,
  sidebar collapse, group collapse, right panel toggle), and that the title bar keeps clear of a
  non-zero `TitleBar().Left` (set it with `ui.Tester.SetTitleBar`).
- README: a short "Examples" note with `go run ./examples/desktop`.

## Non-goals

- Not on the docs site / WebAssembly for now (title bar and window chrome are desktop-only).
- No persistence (window size, open tabs, settings are lost on quit).
- No real data, routing library, i18n, menus in the native menu bar, or keyboard shortcuts beyond
  what components already give.
- No new or changed components in `pkg/`. If a component lacks something the layout needs, work
  around it locally in the example and note it in `TODO.md`.

## Anti-goals

- No AI-chat specifics (model pickers, quota cards, agent lists). It is a generic skeleton.
- No LobeHub logo, mascot or brand names.
- Not a component showcase: each component appears only where the layout needs it.

## Code shape

Files, each under 300 lines, one region per file so a reader finds things by name:

| File | Holds |
| --- | --- |
| `main.go` | window options, `main`, the root view composing the regions |
| `state.go` | the app state struct: tabs (each with a history stack and index), active tab, sidebar width / expanded, right panel width / expanded, group open flags, settings values; and the methods that change it (`navigate`, `back`, `forward`, `openTab`, `closeTab`) |
| `routes.go` | the page table: route key → title, icon, view function. Adding a page = one entry here + one view function |
| `titlebar.go` | title bar: toggle, history buttons, tab strip, right icons |
| `sidebar.go` | sidebar: workspace switcher, nav list, groups, footer |
| `content.go` | the content card, its header row, the right panel |
| `statusbar.go` | the status bar |
| `page_home.go` | the Home page |
| `page_settings.go` | the Settings page |
| `page_placeholder.go` | the generic placeholder page |
| `app_test.go` | the tests |

State lives in one package-level `*app` (as `catalog/home` keeps its showcase state), so views stay
plain functions of `(c *ui.Context)`; logic that changes state is in methods on `app`, testable
without a view.

## Constraints

- Go 1.27, mygo at the version in `go.mod`.
- Follow `catalog/home`'s style: plain view functions, `tk := theme.TokensOf(c)`, `c.Key(...)` for
  repeated rows.
- No comments except where behavior is surprising (repo rule).
- Lint: `gofmt` and `go vet` on `./examples/desktop` only; `go test ./examples/desktop`.

## Left to the builder

- Exact sizes, paddings and font sizes, matched to the screenshot by eye, using tokens.
- Window default size (around 1280 × 820) and minimum size.
- The tab strip is drawn locally (lobe-mygo's `nav.Tabs` is not a closable browser-style strip).
- Placeholder copy and sample names.
- Whether the recent list on Home uses `display.List` or local rows, whichever matches the reference
  better.

## Verification

- `go test ./examples/desktop` passes.
- `go run ./examples/desktop` on macOS, checked by screenshot in light and dark: the four regions
  present and aligned as in the reference, the traffic lights clear of the toggle, each behavior in
  Goals clicked through once.
