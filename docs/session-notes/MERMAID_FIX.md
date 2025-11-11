# Mermaid Charts Fix

## Problem
Mermaid charts were not displaying on pages that had `mermaid = true` in their front matter.

## Root Cause
The `mermaid = true` parameter was incorrectly placed inside the `[taxonomies]` section instead of at the top level of the front matter.

### ❌ Incorrect (Before)
```toml
+++
title = "Example: Mermaid Charts"
date = 2024-06-03
description = "Example post showing Mermaid diagrams"
[taxonomies]
tags = ["apacible", "example", "chart", "mermaid"]
mermaid = true  # ❌ Wrong location - inside taxonomies
+++
```

Hugo's template checks `.Params.mermaid` which requires the parameter to be at the root level, not nested under taxonomies.

### ✅ Correct (After)
```toml
+++
title = "Example: Mermaid Charts"
date = 2024-06-03
description = "Example post showing Mermaid diagrams"
mermaid = true  # ✅ Correct location - root level

[taxonomies]
tags = ["apacible", "example", "chart", "mermaid"]
+++
```

## Files Fixed

1. **content/posts/mermaid-charts-example.md**
   - Moved `mermaid = true` from `[taxonomies]` to root level

2. **content/posts/all-shortcodes-example.md**
   - Moved `mermaid = true` from `[taxonomies]` to root level
   - Also moved `math = true` to root level

## How It Works

### Template Logic (baseof.html)
```go
{{ if .Params.mermaid }}
  {{ partial "mermaid.html" . }}
{{ end }}
```

This checks for `.Params.mermaid` which maps to the root-level parameter in front matter.

### Mermaid Partial (partials/mermaid.html)
- Loads Mermaid.js from CDN (ESM module)
- Configures theme based on light/dark mode
- Initializes on page load
- Re-renders charts when theme changes
- Stores original code to preserve it during re-renders

### Theme Integration
The mermaid partial includes theme-aware configuration:
- **Light theme**: Uses earthy sage green palette
- **Dark theme**: Uses burgundy/plum palette
- Automatically re-renders diagrams when user toggles theme

## Verification

### Test Page
```bash
http://localhost:1313/posts/mermaid-charts-example/
```

### Expected Behavior
✅ Sequence diagrams render
✅ Flowcharts render
✅ Class diagrams render
✅ Entity relationship diagrams render
✅ Gantt charts render
✅ Charts update when toggling light/dark theme
✅ Charts match site color palette

### HTML Output Check
```bash
# Should find 1 mermaid import
grep -c "import mermaid" public/posts/mermaid-charts-example/index.html
# Output: 1
```

### Browser Console
No errors should appear when viewing pages with mermaid charts.

## Summary

| Issue | Status | Cause | Fix |
|-------|--------|-------|-----|
| Mermaid not loading | ✅ Fixed | Param under `[taxonomies]` | Moved to root level |
| Charts not rendering | ✅ Fixed | `.Params.mermaid` not found | Front matter restructured |

**Status**: Mermaid charts now render correctly ✅

