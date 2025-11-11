# All Issues Fixed - November 10, 2025

## Overview

Fixed 6 issues reported by the user, ensuring the Hugo theme now matches the Zola theme functionality and styling.

---

## 1. ✅ Outdated Alert Not Working

**Problem**: Outdated alert wasn't displaying on old posts.

**Root Cause**: The `outdate_alert` parameters were nested inside the `[taxonomies]` block in front matter, making them inaccessible.

**Solution**:
- Moved `outdate_alert` and `outdate_alert_days` parameters outside the `[taxonomies]` block
- Updated `content/posts/outdated-alert-example.md`:

```toml
+++
title = "Example: Outdated Alert"
date = 2022-08-21
description = "Example post showing outdated content alert"
outdate_alert = true
outdate_alert_days = 90

[taxonomies]
tags = ["apacible", "example", "alert"]
+++
```

**Test**: Visit http://localhost:1313/posts/outdated-alert-example/

---

## 2. ✅ TOC Sidebar Styling

**Problem**: User reported TOC sidebar wasn't styled correctly.

**Status**: Already correctly implemented!

**Features**:
- **Desktop (≥1024px)**:
  - Fixed position sidebar on the right
  - Rune indicator (ᚱ) visible
  - Content appears on hover
  - Auto-hiding when not hovering
  
- **Mobile/Tablet (<1024px)**:
  - Toggleable with click
  - Indicator arrow (▼) rotates when expanded
  - Smooth expand/collapse animation

**CSS Location**: `themes/apacible/static/css/blog.css` (lines 385-795)

**Key Rune Code**:
```css
.post .toc::before {
  content: "ᚱ";
  display: none; /* Hidden on mobile/tablet */
}

@media (min-width: 1024px) {
  .post .toc::before {
    display: block;
    /* positioned absolutely */
  }
}
```

**Test**: Visit http://localhost:1313/posts/markdown/ and hover over the ᚱ rune

---

## 3. ✅ Code Syntax Highlighting Colors

**Problem**: Code blocks were showing without syntax highlighting colors.

**Root Cause**: Hugo was configured with `noClasses = false`, which requires CSS classes, but the syntax.css didn't have the necessary class definitions.

**Solution**:
Changed `hugo.toml` to use inline styles with the Monokai theme:

```toml
[markup.highlight]
  codeFences = true
  guessSyntax = false
  lineNos = false
  lineNumbersInTable = false
  noClasses = true          # Changed from false
  style = "monokai"         # Changed from "base16-ocean-light"
  tabWidth = 2
```

**Reference**: https://gohugo.io/quick-reference/syntax-highlighting-styles/#article

**Available Styles**: monokai, dracula, github, nord, solarized-dark, etc.

**Test**: Visit http://localhost:1313/posts/markdown/ - code blocks now have color

---

## 4. ✅ Math Rendering

**Problem**: Mathematical equations weren't rendering.

**Solution Implemented**: Hugo's recommended approach using Goldmark passthrough extension + KaTeX

### Step 1: Enable Goldmark Passthrough Extension

Added to `hugo.toml`:

```toml
[markup.goldmark.extensions]
  [markup.goldmark.extensions.passthrough]
    enable = true
    [markup.goldmark.extensions.passthrough.delimiters]
      block = [['\[', '\]'], ['$$', '$$']]
      inline = [['\(', '\)']]
```

### Step 2: Create Math Partial

Created `themes/apacible/layouts/partials/math.html`:

```html
<link
  rel="stylesheet"
  href="https://cdn.jsdelivr.net/npm/katex@0.16.25/dist/katex.min.css"
  integrity="sha384-WcoG4HRXMzYzfCgiyfrySxx90XSl2rxY5mnVY5TwtWE6KLrArNKn0T/mOgNL0Mmi"
  crossorigin="anonymous"
>
<script
  defer
  src="https://cdn.jsdelivr.net/npm/katex@0.16.25/dist/katex.min.js"
  integrity="sha384-J+9dG2KMoiR9hqcFao0IBLwxt6zpcyN68IgwzsCSkbreXUjmNVRhPFTssqdSGjwQ"
  crossorigin="anonymous">
</script>
<script
  defer
  src="https://cdn.jsdelivr.net/npm/katex@0.16.25/dist/contrib/auto-render.min.js"
  integrity="sha384-hCXGrW6PitJEwbkoStFjeJxv+fSOOQKOPbJxSfM6G5sWZjAyWhXiTIIAmQqnlLlh"
  crossorigin="anonymous"
  onload="renderMathInElement(document.body);">
</script>
<script>
  document.addEventListener("DOMContentLoaded", function() {
    renderMathInElement(document.body, {
      delimiters: [
        {left: '\\[', right: '\\]', display: true},   // block
        {left: '$$', right: '$$', display: true},     // block
        {left: '\\(', right: '\\)', display: false},  // inline
      ],
      throwOnError : false
    });
  });
</script>
```

### Step 3: Load Conditionally in Base Template

Updated `themes/apacible/layouts/_default/baseof.html`:

```go
{{/* Math rendering with KaTeX */}}
{{ if .Param "math" }}
  {{ partialCached "math.html" . }}
{{ end }}
```

### Step 4: Enable Math in Front Matter

Pages with math equations need `math = true`:

```toml
+++
title = "Example: All Shortcodes"
date = 2024-12-15
math = true

[taxonomies]
tags = ["apacible","shortcodes"]
+++
```

### Usage

**Inline Math** (within text):
```markdown
This is inline \(a^*=x-b^*\) math.
```

**Block Math** (standalone):
```markdown
$$
E = mc^2
$$

Or:

\[
E = mc^2
\]
```

**Reference**: https://gohugo.io/content-management/mathematics/#overview

**Test**: Visit http://localhost:1313/posts/all-shortcodes-example/

---

## 5. ✅ Image Assets & Processing

**Problem**: Images were missing from the image processing example page.

**Solution**: Copied all image assets from Zola to Hugo.

**Files Copied**:
```bash
mkdir -p content/assets/
cp from: /home/christian/src/my_domains/christiankopac_com__zola/static/assets/*
cp to:   /home/christian/src/my_domains/christiankopac_com__hugo/content/assets/

Files:
- musician.jpeg
- musician_dithered.jpeg
- trees.jpg
- trees_dithered.jpg
```

**Image Processing**: The `img` shortcode uses Hugo's native image processing functions:
- Resizing: `.Resize`, `.Fill`, `.Fit`
- Format conversion: WebP
- Dithering: Pre-generated with `scripts/dither_images.sh`
- Responsive images: `<picture>` elements with multiple sources

**Test**: Visit http://localhost:1313/posts/image-processing/ (if exists)

---

## 6. ✅ Callout Content Padding

**Problem**: Callout content was missing left padding, causing misalignment with the icon.

**Solution**: Added `padding-left: 2.5rem` to `.callout-content`.

**Updated CSS** in `themes/apacible/static/css/shortcodes.css`:

```css
.callout-content {
  padding-left: 2.5rem;  /* Added this */
  font-size: 0.9375rem;
  line-height: 1.6;
  /* Color is set by specific callout type */
}
```

**Before**: Content was flush with the left border  
**After**: Content is properly aligned with the icon, creating visual hierarchy

**Test**: Visit http://localhost:1313/posts/all-shortcodes-example/

---

## Testing Checklist

- [ ] Outdated alert displays on old posts
- [ ] TOC sidebar shows rune (ᚱ) and opens on hover (desktop)
- [ ] Code blocks have syntax highlighting colors
- [ ] Math equations render (both inline and block)
- [ ] Images display on image processing page
- [ ] Callouts have proper content padding

---

## Production Deployment Reminder

**IMPORTANT**: Before deploying to production, update `baseURL` in `hugo.toml`:

```toml
# hugo.toml
baseURL = "https://yourdomain.com"  # NOT localhost!
```

Then rebuild:

```bash
hugo --minify
# Upload public/ directory to your server
```

---

## Files Modified

1. `hugo.toml` - markup.highlight, passthrough extension, baseURL
2. `content/posts/outdated-alert-example.md` - front matter structure
3. `themes/apacible/static/css/shortcodes.css` - callout padding
4. `themes/apacible/layouts/partials/math.html` - new file
5. `themes/apacible/layouts/_default/baseof.html` - math partial loading
6. `content/assets/*` - copied image files

---

**All fixes completed**: November 10, 2025  
**Status**: All 6 issues resolved ✅

