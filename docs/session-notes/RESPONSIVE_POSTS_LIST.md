# Responsive Posts Listing

## Implementation

Added responsive CSS to hide tags on smaller screens while keeping title and date visible.

## CSS Changes

**File**: `themes/apacible/static/css/blog.css`

```css
/* Responsive: Hide tags on small screens */
@media (max-width: 768px) {
  .blog .post-preview .tags {
    display: none;
  }
}
```

## Responsive Behavior

### Large Screens (> 768px)
```
┌─────────────────────────────────────────────────┐
│  ᛟ Post Title ──────────────────────────────→  │
│                        [tag1] [tag2] Nov 10     │
└─────────────────────────────────────────────────┘
```

Shows:
- ✅ Title (with featured marker if applicable)
- ✅ Tags
- ✅ Date

### Small Screens (≤ 768px)
```
┌─────────────────────────────────┐
│  ᛟ Post Title ─────────────→   │
│                      Nov 10     │
└─────────────────────────────────┘
```

Shows:
- ✅ Title (with featured marker if applicable)
- ❌ Tags (hidden)
- ✅ Date

## Breakpoint

- **768px**: Common tablet/mobile breakpoint
- Devices ≤ 768px width will hide tags
- Devices > 768px width will show tags

## Testing

### Desktop View (Wide Screen)
```bash
# Open browser, visit http://localhost:1313/posts/
# Browser width > 768px
# Should see: Title + Tags + Date
```

### Mobile View (Narrow Screen)
```bash
# Open browser DevTools (F12)
# Toggle device toolbar (Ctrl+Shift+M)
# Select mobile device (iPhone, etc.)
# Should see: Title + Date only (no tags)
```

### Manual Resize Test
```bash
# Visit http://localhost:1313/posts/
# Slowly resize browser window narrower
# At 768px width, tags should disappear
# Date should remain visible
```

## Why This Approach?

1. **Progressive Enhancement**: Full information on large screens, essential info on small screens
2. **Readability**: Fewer elements on small screens = less clutter
3. **Performance**: CSS-only solution, no JavaScript needed
4. **Accessibility**: Content still accessible (tags visible in actual posts)

## CSS Specificity

The selector `.blog .post-preview .tags` ensures:
- Only affects posts listing tags
- Doesn't affect tags in actual post content
- Doesn't affect tags on single post pages

## Alternative Breakpoints

If you want to adjust the breakpoint:

```css
/* Tablets and below */
@media (max-width: 1024px) {
  .blog .post-preview .tags {
    display: none;
  }
}

/* Mobile only */
@media (max-width: 640px) {
  .blog .post-preview .tags {
    display: none;
  }
}
```

## Future Enhancements (Optional)

If you want to add more responsive behavior:

```css
/* Stack metadata on very small screens */
@media (max-width: 480px) {
  .blog .post-preview .post-meta {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .blog .post-preview .post-meta-right {
    margin-top: var(--space-sm);
  }
}

/* Smaller font on mobile */
@media (max-width: 768px) {
  .blog .post-preview h2 {
    font-size: var(--font-size-xl);
  }
  
  .blog .post-preview time {
    font-size: var(--font-size-sm);
  }
}
```

---

**Feature added**: November 10, 2025  
**Breakpoint**: 768px  
**Status**: Responsive posts listing ✅

