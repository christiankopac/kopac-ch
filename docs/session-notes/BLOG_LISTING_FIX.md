# Blog Listing Minimalism Fix

## Issue
The user noted that the blog listing page in Zola was more minimal than the Hugo version.

## Comparison

### Zola Version (`templates/post.html`)
Looking at the Zola single post template for reference to the listing style:

```jinja
<div class="post-meta">
    <div class="post-meta-left">
        <!-- Date, word count, reading time -->
    </div>
    <div class="post-meta-right">
        <!-- Tags only -->
    </div>
</div>
```

### Zola Blog Listing (`templates/blog.html`)
```jinja
<div class="post-meta">
    <div class="post-meta-left">
        <!-- EMPTY - minimal approach -->
    </div>
    <div class="post-meta-right">
        <!-- Tags + Date together -->
    </div>
</div>
```

Key observation: **The blog listing is more minimal** with everything aligned to the right.

### Hugo Fixed Version (`layouts/_default/list.html`)
```go-html-template
<div class="post-meta">
    <div class="post-meta-left">
        <!-- EMPTY - matches Zola -->
    </div>
    <div class="post-meta-right">
        <!-- Tags + Date (in that order) -->
        {{ with .Params.tags }}
        <div class="tags">
            {{ range . }}
                <span class="tag">{{ . }}</span>
            {{ end }}
        </div>
        {{ end }}
        {{ with .Date }}
            <time datetime="{{ .Format "2006-01-02" }}">
                {{ .Format ($.Site.Params.date_format | default "Jan 2, 2006") }}
            </time>
        {{ end }}
    </div>
</div>
```

## Key Points

1. **Minimal Layout**: Left column is empty on listing pages
2. **Right Alignment**: Both tags and date are on the right side
3. **Order Matters**: Tags first, then date
4. **Single Post vs Listing**:
   - Single post: More metadata (date, word count, reading time on left; tags on right)
   - Blog listing: Less metadata (empty left; tags + date on right)

## Visual Layout

```
┌─────────────────────────────────────────────────┐
│  Post Title                                     │
├─────────────────────────┬───────────────────────┤
│  post-meta-left         │  post-meta-right      │
│  (EMPTY)                │  [tag1] [tag2] Jan 15 │
├─────────────────────────┴───────────────────────┤
│  Post summary text...                           │
└─────────────────────────────────────────────────┘
```

## CSS Considerations

The CSS likely uses flexbox/grid to create this two-column layout:

```css
.post-meta {
  display: flex;
  justify-content: space-between;
}

.post-meta-left {
  /* Empty on listing pages */
}

.post-meta-right {
  display: flex;
  gap: var(--space-sm);
  align-items: center;
}
```

## Status

✅ **Fixed** - Blog listing now matches Zola's minimal approach with empty left column and tags + date on the right.

