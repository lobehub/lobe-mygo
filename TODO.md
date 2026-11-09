# Lobe UI parity TODO

What the components in `pkg/` still lack against [Lobe UI](https://ui.lobehub.com), compared prop by prop with `lobehub/lobe-ui` `src/<Component>/type.ts` and its demos. Last audited 2026-11-08.

- **P1**: common in consumer apps; do next.
- **P2**: worth having; do when the component is touched again.
- **P3**: rare, or covered well enough another way.
- **mygo**: blocked on a capability mygo lacks today; listed under [Blocked on mygo](#blocked-on-mygo).

React-only props (`className`, `classNames`, `styles`, `ref`, `as`, controlled/`default*` pairs, portal/positioner plumbing) have no Go counterpart and are left out: state lives behind the pointers the components take.

## Cross-cutting

- [x] **P1** Built-in texts in `en` and `zh-Hans` through the `lobe` library catalog (`pkg/internal/l10n`); apps override with `lobe.<key>`. A `Texts` override is still missing where Lobe has `texts`.
- [ ] **P2** `shadow` on the input family (Input, TextArea, InputNumber, InputOTP, Select, SearchBar, AutoComplete, HotkeyInput, Segmented).
- [ ] **P2** `loading` on Button, ActionIcon, Switch and Select.
- [ ] **mygo** `glass` (backdrop blur) on Block, ActionIcon, ActionIconGroup, Alert, Segmented and the image viewer's toolbar.

## General

### Button
- [ ] **P1** `loading` (spinner, keeps width), `icon` + `iconPosition`, `danger`.
- [ ] **P1** `block` (full width) and `shape` (`round`, `circle`); callers set `Radius` by hand today.
- [ ] **P3** `ghost`, link `href`/`target`.

### ActionIcon
- [ ] **P2** `loading`, `spin` (turning glyph), `disabled` as a prop (only the element's `Disabled` today).
- [x] **P3** `blockSize` and numeric `size` (`BlockSize`, `GlyphSize`).
- [ ] **P3** `borderRadius` override, `outdent` for borderless icons.

### ActionIconGroup
- [ ] **P3** `actionIconProps` passed to every icon.

### Block
- [ ] **mygo** `glass`.

### Text
- [ ] **P2** `align`, numeric `weight`, `color`, `textTransform`, `noWrap`, `wordBreak`.
- [ ] **P2** `ellipsis.tooltip` / `tooltipWhenOverflow`: show the full text in a tooltip when it is cut.
- [ ] **P3** `shinyDuration` (fixed at 1.5 s).

### Icon
- [ ] **P3** Size presets beyond a number (`size: 'small' | 'middle' | 'large'` config objects).

### Divider
- [ ] **P3** Label placement (`orientation` start / end for the text); only centered today.

## Navigation

### Tabs
- [ ] **P3** Keys instead of indices (`activeKey`, `items[].key`); Go uses `*int`.

### Pagination
- [ ] **P1** `showSizeChanger` + `pageSizeOptions` + `onPageSizeChange`.
- [ ] **P2** `showQuickJumper` ("Go to" input, optional go button).
- [ ] **P2** `simple` mode (`< [3] / 12 >`, optional read-only).
- [ ] **P3** `itemRender`, `showTitle`, `responsive` (small below 576 px).

### Breadcrumb
- [ ] **P3** `href` links (OnClick covers it in an app).

### Steps
- Complete.

### DropdownMenu
- [ ] **P1** Submenus (`SubmenuTrigger`), checkbox, radio and switch items.
- [ ] **P2** `header` / `footer` slots, group labels.
- [ ] **P2** `placement`, `trigger: 'hover'`.
- [ ] **P3** `iconAlign` / `iconSpaceMode` (reserve the icon column when only some items have icons).

### ContextMenu
- [ ] **P1** Submenus and checkbox / radio items, as DropdownMenu.
- [ ] **P2** `header` / `footer`.

### List
- [ ] **P3** `href` items; `selectable` is the caller's (`Active` is passed in).

## Data Entry

### Input
- [ ] **P1** `allowClear` (clear button), `onPressEnter`.
- [ ] **P2** `shadow`.

### TextArea
- [ ] **P1** `showCount` (with a max length).
- [ ] **P2** `allowClear`, `onPressEnter`, `resize` handle.

### InputNumber
- [ ] **P2** `changeOnWheel`, `prefix` / `suffix`, `onPressEnter`.

### InputPassword
- [ ] **P3** `visibilityToggle: false`, toggle labels for assistive technology.

### InputOTP
- Complete apart from `shadow`.

### Select
- [ ] **P1** `showSearch` with `filterOption`, `mode: 'multiple' | 'tags'` with token separators.
- [ ] **P1** `allowClear`, `loading`, `prefix`.
- [ ] **P2** `optionRender` / `labelRender`, option groups, `listHeight` with virtualization.
- [ ] **P3** `popupMatchSelectWidth`, `selectedIndicatorVariant`, `behaviorVariant`.

### AutoComplete
- [ ] **P2** `allowClear`, `suffix`, `onSearch` (async options), `disabled`.

### SearchBar
- [ ] **P2** `spotlight` (glow while focused), `enableShortKey` toggle, `onInputChange` vs `onSearch` split.

### Checkbox
- [ ] **P1** `Checkbox.Group` with `options`.
- [ ] **P3** `backgroundColor`, `textProps`.

### Radio
- [ ] **P2** `RadioGroup` from `options` (today the caller builds each `Radio`).
- [ ] **P3** `backgroundColor`, `textProps`.

### Switch
- [ ] **P1** `loading`.
- [ ] **P2** `checkedChildren` / `unCheckedChildren` (text or icon inside the track).

### Slider
- [ ] **P1** `SliderWithInput` (slider plus a number field, `unlimitedInput`).
- [ ] **P2** `onChangeComplete`, marks / dots, range (two thumbs).

### Rate
- [ ] **P2** `readOnly`, `gap`.

### Segmented
- [ ] **P2** Per-option icons and disabled options, `vertical`, `disabled`.
- [ ] **P2** `shadow`; **mygo** `glass`.

### ToggleGroup
- Complete.

### EditableText
- [ ] **P3** `gap`, `inputProps` (placeholder, max length), submit icon / `texts`.

### HotkeyInput
- [ ] **P2** `allowReset` vs `allowClear` split, `onClear` / `onReset` / `onFocus` / `onBlur` callbacks.

### ColorPicker
- [ ] **P1** Recent colors row (`recent`).
- [ ] **P2** RGB / HSL input formats (`defaultFormat`, format switcher), `showText`, `onChangeComplete`.
- [ ] **mygo** Eyedropper (EyeDropper API in browsers, screen sampling natively).

### ColorSwatches
- [ ] **P3** `renderPicker`, `enableColorSwatches: false`, `texts`.

### DatePicker / Calendar
- [ ] **P1** RangePicker (`range`, `allowEmpty`, `order`).
- [ ] **P1** Typing a date into the field.
- [ ] **P2** Time picking, `multiple` dates with tags (`maxTagCount`, `tagRender`).
- [ ] **P2** Popup aligned to the field's start edge (centered today; Lobe is bottom-left).
- [ ] **P3** `locale` (month and weekday names).

### Form
- [ ] **P1** `Form.List`: array fields with add / remove / move and per-row keys.
- [ ] **P2** `items` schema config (groups and fields from data) and a form-level `schema` validator.
- [ ] **P2** `render` fields (`FieldRenderProps`: value, onChange, error, touched, validating) and `bare` fields.
- [ ] **P2** Submit loading while an async `OnSubmit` runs; `enableUnsavedWarning` (ask before leaving).
- [ ] **P3** `extra`, `avatar`, `hidden` on fields; `activeKey` / `onCollapse` to control groups; `itemsType: 'flat'`.

### Upload
- [ ] **P2** `beforeUpload` returning a promise: `BeforeUpload` runs on the UI thread, so a slow check (reading a file's header) freezes it.
- [ ] **P3** Picking on the web needs the input's `cancel` event (Chrome 113, Firefox 91, Safari 16.4): an older browser never reports a canceled chooser and the Upload opens no other until a reload.

### ImageSelect
- [ ] **P2** Pictures by URL (`img: string`): options take a decoded `ui.ImageSource`, so the app loads remote pictures itself.

## Data Display

### Markdown
- [ ] **P2** Code blocks grow at the reveal's pace while streaming but don't fade in per character (Highlighter's text isn't one mygo's reveal tracks).
- [ ] **P3** `enableStream: false` while `animated` (re-render without pacing): use `Animated: false`.
- [ ] **P2** Images open in the Image viewer as a gallery (`enableImageGallery`); mygo draws them itself, with only `LoadImage` to hook.
- [ ] **P2** Lists in the default variant use Lobe's "-" bullets; mygo draws its own markers.
- [ ] **P2** Table cells at least 120 DIPs wide; `details` / `summary` blocks.
- [ ] **mygo** Mermaid (planned through pele) and LaTeX: blocks show their source.
- [ ] **P3** `citations`, `allowHtml` component lists, `customRender`, remark/rehype plugins: mygo's parser takes no syntax plugins, only render hooks.

### Highlighter
- [ ] **P2** File type icon before the full-featured header's name (`icon`, Lobe's MaterialFileTypeIcon), with FileTypeIcon.
- [ ] **P2** `animated`: code fading in as it streams, for Markdown's Stream.
- [ ] **P3** `allowChangeLanguage` (a language select in the header), `theme` (other code themes), `actionIconSize`.
- [ ] **P3** Languages beyond the ~30 lexers of mygo's editor (Vue/Svelte templates, Elixir, Haskell, …) show plain.

### Snippet
- [ ] **P3** `spotlight` (the glow following the pointer).

### Avatar
- [ ] **P1** `AvatarGroup` (overlapping stack, `max` with "+N").
- [ ] **P2** `bordered` + `borderedColor`, `shadow`, `title` tooltip.
- [ ] **P3** `animation` (animated emoji), `emojiScaleWithBackground`, `sliceText`, `fallbackColors`.

### Badge
- [ ] **P3** `offset`.

### Tag
- [ ] **P3** `closeIcon`, square `shape`.

### Hotkey
- [ ] **P3** `inverseTheme` (light keys on a dark tooltip and the reverse).

### Image
- [ ] **P2** Copy image to the clipboard from the viewer (mygo `WriteClipboardData` works natively and on the web).
- [ ] **P2** `maxWidth` / `maxHeight` / `minWidth` / `minHeight`; `defaultZoom` policy (`auto` / `actual` / `fit`) and `autoZoomThreshold`.
- [ ] **P2** Slide-in switch between gallery images (fades and scales today), chrome that hides after 2 s idle, wheel past the zoom floor closing the viewer.
- [ ] **P3** `toolbarAddon`, `maxScale`, flip horizontal / vertical, rubber-band past the pan edges.

### Skeleton
- [ ] **P3** `radius`, `SkeletonText` `lineHeight` / `gap` / last-row `width`, skeleton buttons.

### Empty
- [ ] **P2** `emoji` and `image` instead of an icon.
- [ ] **P3** `iconColor`, `imageSize`, `align`.

### Accordion
- [ ] **P3** `keepMounted`.

### GroupAvatar
- [ ] **P3** Squircle / smooth / iOS corners are rounded squares: mygo has no path masks.

### QRCode
- [ ] **P3** `statusRender` (a custom overlay per status), `iconSize`.

### SortableList
- [ ] **P2** Scroll the container while a row is dragged near its edge.
- [ ] **P3** The copy under the pointer stays inside the window's top only, not all its edges.

### Collapsible, Descriptions, Statistic, Table, Tree
Not started; see below.

## Feedback

### Alert
- [ ] **P1** `action` slot (a button beside the message); `onClose` / `afterClose` callbacks (closing works through `Closed`).
- [ ] **P2** `extra` expandable details (`extraDefaultExpand`, `extraIsolate`), `banner`.
- [ ] **P3** `colorfulText`, `closeText`, `icon` override; **mygo** `glass`.

### Modal
- [ ] **P1** `confirmLoading` / `loading`, custom `footer`, `closable: false`, `maskClosable: false`.
- [ ] **P1** Imperative `confirm` / `alert` helpers (Lobe's `createModal`).
- [ ] **P2** `allowFullscreen`, `centered: false` (top offset), `height`.
- [ ] **P3** `draggable`, `keyboard: false`, `destroyOnHidden`.

### Drawer
- [ ] **P2** `extra` header slot, `mask: false`, `maskClosable: false`, `noHeader`.
- [ ] **P2** Nested drawers pushing their parent (`push`).
- [ ] **P3** `sidebar` + `sidebarWidth`, `containerMaxWidth`, `afterOpenChange`.

### Toast
- [ ] **P2** `limit` (visible count), swipe to dismiss (`swipeDirection`), a host-wide `duration`.
- [ ] **P2** Promise toasts (loading that turns into success or error) as a helper over `ShowToast` with an ID.

### Progress
- [ ] **P3** Gradient `strokeColor`.

### Spin
- [ ] **P3** Custom `indicator`.

### Result
- Complete.

### FloatingPanel
- [ ] **P2** Move the bottom-right toasts out of the panel's way while it shows (Lobe's `ToastDodge`).
- [ ] **P3** `modal` (focus kept inside), separate x / y `offset`, `afterClose`.

### FloatingSheet
- Complete.

## Overlay

### Tooltip
- [ ] **P2** `openDelay` / `closeDelay`, `disabled`, controlled `open`.
- [ ] **P3** `arrow`.

### Popover
- [ ] **P1** `trigger: 'hover'` with open / close delays.
- [ ] **P2** `arrow`, controlled `disabled`.

## Layout

### DraggablePanel
- [x] **P2** `showBorder: false` (`NoBorder`).
- [ ] **P2** `top` / `bottom` placements, `mode: 'float'`.
- [ ] **P3** `DraggablePanel.Header` (`title`, `extra`), `showHandleWideArea`, `showHandleWhenCollapsed`.
- [ ] **P2** The resize handle stays live while the panel is collapsed: dragging the 8 DIP strip at the edge changes the hidden size and covers the neighbour's edge (seen in `examples/desktop`). Disable it, or let a drag there expand the panel.

### ScrollArea
- [ ] **P3** `corner` slot, `disableContentFit`.

### Grid
- [ ] **P3** `rows`.

### Header, Layout
- [ ] **P2** `Mobile` is the caller's: no breakpoint follows the window's width.
- [ ] **P2** Layout scrolls its content in a region under the header, so the footer spans the content column, not the sidebar too, and nothing shows through the header (backdrop blur is blocked on mygo).

### SideNav, DraggableSideNav, GuideCard
- Complete.

## Blocked on mygo

- Backdrop blur (`backdrop-filter`): `glass` everywhere, the viewer toolbar and floating footers. Today they fall back to a 60 % layout-color fill.
- Drop shadow of a shape (`filter: drop-shadow`): the image viewer's control halo. The close button got a pill instead.
- Save dialog on the web: `mygo.Dialog.Save` returns unsupported in the browser, so `Image` leaves downloads to `OnDownload`.
- Screen color sampling: ColorPicker's eyedropper.

## Not started

In the order agreed for consumer apps; B-end components last.

1. Mermaid through pele (goja natively, the browser's JS on the web).
2. EmojiPicker, FluentEmoji.
3. Collapsible.
4. Carousel, Video, DownloadButton, FileTypeIcon.
5. Toc, Anchor, ScrollShadow, ProgressiveBlur, MaskShadow.
6. Descriptions, Statistic, Table / ProTable, Tree.
7. CodeEditor, CodeDiff, HtmlPreview, EditorSlashMenu.
8. Subpath entries: `chat` (ChatItem, MessageInput, ChatList), `mobile`, `awesome`, `dashboard`.
