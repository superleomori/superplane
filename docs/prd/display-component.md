# PRD: Display Component

**Status:** Draft  
**Author:** Leo (AI assistant)  
**Date:** 2026-05-16  

---

## Problem

When building and debugging workflows on the SuperPlane canvas, there is no way to see the value of intermediate data inline on the canvas itself. Developers currently have to:

- Click into a node's run history and navigate the payload JSON manually
- Add a temporary HTTP Request pointing to a request bin (noisy, requires external setup)
- Re-run the whole canvas and hope the right node tells them what they need

There is no `console.log` equivalent — a lightweight way to say "show me this value right here, right now, during a run."

---

## Goal

Introduce a **Display** action component that:

1. Renders an expression-evaluated value **directly on the canvas node face**
2. Lets the user control the **color** of the display badge using an expression
3. Passes the upstream event through **unchanged** so it doesn't interrupt the workflow chain
4. Requires zero external integrations or side effects

---

## Non-Goals

- This is not a logging system — it does not persist a log stream or support log levels
- It does not replace the run payload inspector in the sidebar
- It is not a alerting or notification component
- No new storage is introduced beyond what execution records already track

---

## User Stories

**As a workflow builder**, I want to place a Display node after any action node and see the value of a specific payload field on the canvas face, so I can confirm the data shape without clicking through run history.

**As a workflow builder**, I want to color the display badge red or green based on an expression (e.g. success/failure), so I can spot failures at a glance across a complex canvas.

**As a workflow builder**, I want the Display node to be completely transparent to the rest of the workflow — events pass through unchanged — so I can add/remove it without rewiring anything.

---

## Component Spec

### Identity

| Field | Value |
|---|---|
| Type | Action (participates in runs) |
| Name | `core.display` |
| Label | Display |
| Category | Core |
| Icon | `monitor` (or similar inspect/eye icon) |

### Configuration Fields

| Field | Key | Type | Required | Description |
|---|---|---|---|---|
| Value | `value` | text (expression) | Yes | The value to display on the node. Supports `{{ }}` expression syntax. |
| Color | `color` | text (expression) | No | Badge color. Supports `{{ }}` expressions. Resolves to a color name (see Color Values). Default: `gray`. |

#### Color Values

| Value | Renders As |
|---|---|
| `green` | Green badge (success / healthy) |
| `yellow` | Yellow badge (warning / pending) |
| `red` | Red badge (failure / critical) |
| `blue` | Blue badge (info / in-progress) |
| `gray` | Gray badge (neutral / default) |

Any unrecognized value falls back to `gray`.

### Execution Behavior

1. The Display node receives the upstream event
2. It evaluates `value` and `color` against the run payload using the standard expression engine
3. The resolved `{value, color}` is stored on the execution record as `display_result`
4. The upstream event is re-emitted **unchanged** on the default output channel
5. The node never fails due to an expression error — if evaluation fails, `value` falls back to `[expression error]` and `color` falls back to `gray`

### Output Event

The Display component emits the upstream event unmodified. The output channel is named `default` (same pattern as No-Op and Filter pass-through).

### Canvas Node Rendering

The node face shows a colored badge with the resolved value from the **latest run**:

```
┌──────────────────────────────────┐
│  👁  Display                      │
│                                  │
│  ┌────────────────────────────┐  │
│  │  ✅  deployed: abc123      │  │  ← green
│  └────────────────────────────┘  │
│                                  │
│  ○ default                       │
└──────────────────────────────────┘
```

- Badge is only shown when there is a resolved value from the latest run
- Badge updates live as new runs complete
- If the node has never run, the badge area is empty (no placeholder text)
- Badge text is truncated at ~60 chars with an ellipsis if too long; full value is visible in the sidebar

### Expression Examples

**Show a field value:**
```
Value:  {{$['HTTP Request'].data.body.status}}
Color:  gray
```

**Green/red based on boolean:**
```
Value:  {{$['Deploy'].data.success ? "✅ deployed" : "❌ failed"}}
Color:  {{$['Deploy'].data.success ? "green" : "red"}}
```

**Show previous node's result:**
```
Value:  {{previous().data.result}}
Color:  {{previous().data.result == "ok" ? "green" : "red"}}
```

**Nil-safe with fallback:**
```
Value:  {{$['Filter'].data?.body?.error ?? "no error"}}
Color:  {{$['Filter'].data?.body?.error != nil ? "red" : "green"}}
```

---

## Canvas YAML

```yaml
- name: Debug Status
  type: core.display
  configuration:
    value: "{{$['Deploy'].data.status}}"
    color: "{{$['Deploy'].data.status == \"success\" ? \"green\" : \"red\"}}"
  connections:
    - from: Deploy
      channel: default
```

---

## Implementation Notes

### Backend

- New `ComponentType = core.display` in the component registry
- Execution handler: evaluate `value` and `color` expressions → store as `display_result` on the execution record → forward upstream event on `default` channel
- Expression errors should be caught and stored as `{value: "[expression error: <msg>]", color: "gray"}` — never fail the run
- `display_result` should be included in the execution record API response

### Frontend

- Canvas node renderer: if `display_result` is present on the latest execution, render the colored badge
- Color → CSS class mapping: `green` → success token, `red` → danger token, `yellow` → warning token, `blue` → info token, `gray` → neutral token
- Sidebar: add a "Display" section showing the resolved value (untruncated) and color for each run in the history

### No new storage required

`display_result` is stored alongside the existing execution payload data — no schema migration needed beyond adding a new optional field to the execution record.

---

## Open Questions

1. Should the Display node support **multiple values** (i.e. a list of `{label, value, color}` rows)? Kept simple for V1 — single value only.
2. Should there be a **title** field separate from `value`? e.g. "Status: ✅ ok" vs just "✅ ok". Deferred to V1 feedback.
3. Should the badge persist after the node is disconnected from a run (e.g. if you pause the node)? Proposal: yes, show last known value with a faded style.

---

## Success Metrics

- Workflow builders can confirm payload values without leaving the canvas or opening the sidebar
- No existing workflow behavior is altered by adding a Display node
- Color-coded badges are readable at typical canvas zoom levels
