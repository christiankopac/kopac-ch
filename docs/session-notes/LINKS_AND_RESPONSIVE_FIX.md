# Content Links & Responsive Behavior Fix

## ✅ Issue 1: Content Page Links Not Working

**Problem**: Links to posts were pointing to `https://example.com/posts/...` instead of local URLs

**Root Cause**: `baseURL` was set to production URL in `hugo.toml`

**Fix Applied**:

```toml
# Before
baseURL = "https://example.com"

# After (for local development)
baseURL = "http://localhost:1313/"
```

**Note**: Remember to change this back to your production URL before deploying!

```toml
# For production deployment
baseURL = "https://yourdomain.com"
```

---

## ✅ Issue 2: Responsive Behavior

**Status**: Already working correctly!

The `blog.css` file already has responsive behavior built in (lines 162-174):

```css
.post-preview .post-meta .tags {
  display: none; /* Hidden on mobile */
}

@media (min-width: 1024px) {
  .post-preview .post-meta .tags {
    display: flex;
    flex-wrap: nowrap;
    gap: 0.5rem;
    align-items: center;
    white-space: nowrap;
  }
}
```

### Breakpoint: 1024px

- **< 1024px** (Mobile/Tablet): Title + Date only
- **≥ 1024px** (Desktop): Title + Tags + Date

This is more conservative than 768px, ensuring better mobile experience.

---

## About Runes

The Zola blog listing template (`templates/blog.html`) does **NOT** use runes in the posts listing.

Runes are only used in:
- Single post pages (for date, word count, reading time metadata)
- Tags in post headers
- TOC indicators
- Post footers

The minimal blog listing style deliberately omits decorative elements like runes to maintain simplicity.

If you want to add runes to the blog listing, we can add them. Current Zola style:

```
Posts
──────────────────────────────────────

Post Title ────────────────────→
                  [tag1] [tag2] Nov 10

Another Post ──────────────────→
             [example] [hugo] Nov 9
```

With runes (if desired):

```
Posts
──────────────────────────────────────

ᛟ Post Title ──────────────────→
              ᛞ [tag1] [tag2] ᛉ Nov 10

Another Post ──────────────────→
         ᛞ [example] [hugo] ᛉ Nov 9
```

---

## Verification

### Test Links
```bash
# Visit posts listing
curl -s http://localhost:1313/posts/ | grep 'href=' | grep posts

# Should show:
# href="http://localhost:1313/posts/detail/" ✓
# href="http://localhost:1313/posts/markdown/" ✓
```

### Test Responsive Behavior
```bash
# Open browser DevTools (F12)
# Toggle device toolbar (Ctrl+Shift+M)
# Resize browser width:
#
# At 1024px: Tags should appear/disappear
# < 1024px: Title + Date only
# ≥ 1024px: Title + Tags + Date
```

---

## Production Deployment Checklist

Before deploying to production:

```toml
# hugo.toml
baseURL = "https://youractualdom ain.com"  # ← Change this!
```

Then rebuild:

```bash
hugo --minify
# Upload public/ directory to your server
```

---

**Fixes completed**: November 10, 2025  
**Status**: Links working, responsive behavior confirmed ✅

