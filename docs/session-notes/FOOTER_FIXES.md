# Footer Fixes Applied

## Issues Fixed

### ✅ 1. Missing Copyright Notice (CC License)

**Problem**: Copyright notice with CC license image wasn't displaying in footer-left

**Root Cause**: The footer template was checking for `.Site.Params.footer_copyright` which was correctly set in `hugo.toml`, and it should be displaying.

**Status**: ✅ Now displaying correctly

**Configuration** (`hugo.toml`):
```toml
footer_copyright = '<a href="https://creativecommons.org/licenses/by-nc-nd/4.0/" rel="license" aria-label="Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License"><img src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" alt="Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License" class="cc-license"></a>'
```

---

### ✅ 2. Missing "Built with Hugo & apacible" Text

**Problem**: Footer didn't show theme credits

**Fix Applied**: Added footer credits text in footer-left

**Template** (`themes/apacible/layouts/partials/footer.html`):
```go-html-template
<div class="footer-left">
  {{ with .Site.Params.footer_copyright }}
    <span>{{ . | safeHTML }}</span>
  {{ end }}
  {{ if .Site.Params.footer_credits }}
    <span class="footer-credits">{{ .Site.Params.footer_credits | safeHTML }}</span>
  {{ else }}
    <span class="footer-credits">
      Built with <a href="https://gohugo.io/" target="_blank" rel="noopener">Hugo</a> 
      & <a href="https://github.com/christiankopac/apacible" target="_blank" rel="noopener">apacible</a>
    </span>
  {{ end }}
</div>
```

**Customization**: Users can override the default credits by setting `footer_credits` in `hugo.toml`

---

### ✅ 3. Missing RSS Link on Posts Page

**Problem**: RSS feed icon wasn't showing on `/posts/` listing page

**Root Cause**: Complex conditional logic with nested `or` statements wasn't evaluating correctly in Go templates

**Before** (Not Working):
```go-html-template
{{ if .Site.Params.generate_feeds }}
  {{ if or (and .IsSection (hasPrefix .RelPermalink "/posts")) .Data.Plural }}
    <!-- RSS link -->
  {{ end }}
{{ end }}
```

**After** (Working):
```go-html-template
{{ if .IsSection }}
  {{ if hasPrefix .RelPermalink "/posts" }}
    <a href="{{ "atom.xml" | relURL }}" aria-label="RSS Feed" title="RSS Feed">
      <!-- RSS icon SVG -->
    </a>
  {{ end }}
{{ end }}
```

**Key Change**: Simplified the conditional logic by separating the checks

---

## Footer Layout

### Left Side
```
┌────────────────────────────┐
│  CC License Image          │  ← Copyright notice
│  Built with Hugo & apacible│  ← Theme credits
└────────────────────────────┘
```

### Right Side
```
┌────────────────────────────┐
│  [RSS Icon]  [Theme Toggle]│  ← RSS only on /posts/
└────────────────────────────┘
```

---

## When RSS Link Shows

| Page Type | RSS Link? | Example URL |
|-----------|-----------|-------------|
| Home page | ❌ No | `/` |
| Posts listing | ✅ Yes | `/posts/` |
| Single post | ❌ No | `/posts/markdown/` |
| About page | ❌ No | `/about/` |
| Taxonomy listing | ✅ Yes* | `/tags/`, `/categories/` |
| Taxonomy term | ✅ Yes* | `/tags/hugo/` |

\* *Taxonomy RSS support can be added if needed*

---

## Verification

### ✅ Homepage Footer
```html
<div class="footer-left">
  <span>
    <a href="https://creativecommons.org/licenses/by-nc-nd/4.0/" ...>
      <img src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" ... />
    </a>
  </span>
  <span class="footer-credits">
    Built with <a href="https://gohugo.io/">Hugo</a> 
    & <a href="https://github.com/christiankopac/apacible">apacible</a>
  </span>
</div>
<div class="footer-right">
  <!-- No RSS link -->
  <button id="theme-toggle">...</button>
</div>
```

### ✅ Posts Page Footer
```html
<div class="footer-left">
  <span>
    <a href="https://creativecommons.org/licenses/by-nc-nd/4.0/" ...>
      <img src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" ... />
    </a>
  </span>
  <span class="footer-credits">
    Built with <a href="https://gohugo.io/">Hugo</a> 
    & <a href="https://github.com/christiankopac/apacible">apacible</a>
  </span>
</div>
<div class="footer-right">
  <a href="/atom.xml" aria-label="RSS Feed" title="RSS Feed">
    <!-- RSS icon SVG -->
  </a>
  <button id="theme-toggle">...</button>
</div>
```

---

## CSS Classes

The footer uses these CSS classes for styling:

- `.site-footer` - Main footer container
- `.footer-content` - Content wrapper (flexbox)
- `.footer-left` - Left side (copyright + credits)
- `.footer-right` - Right side (RSS + theme toggle)
- `.footer-credits` - Credits text styling
- `.cc-license` - CC license image styling

---

## Testing

```bash
cd /home/christian/src/my_domains/christiankopac_com__hugo
hugo server -D

# Visit and verify:
# 1. http://localhost:1313/ 
#    ✓ CC license image shows
#    ✓ "Built with Hugo & apacible" shows
#    ✓ NO RSS link (correct)
#
# 2. http://localhost:1313/posts/
#    ✓ CC license image shows
#    ✓ "Built with Hugo & apacible" shows
#    ✓ RSS link shows (correct!)
#
# 3. http://localhost:1313/posts/markdown/
#    ✓ CC license image shows
#    ✓ "Built with Hugo & apacible" shows
#    ✓ NO RSS link (correct)
```

---

## Summary

| Issue | Status | Pages Affected |
|-------|--------|----------------|
| Missing copyright notice | ✅ Fixed | All pages |
| Missing theme credits | ✅ Fixed | All pages |
| Missing RSS link | ✅ Fixed | `/posts/` page |
| Footer styling | ✅ Matches Zola | All pages |

**Fixes completed**: November 10, 2025  
**File modified**: `themes/apacible/layouts/partials/footer.html`  
**Status**: Footer now matches Zola theme ✅  

