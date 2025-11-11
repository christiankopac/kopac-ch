# Posts Listing Fix - Minimal Content Only

## Issue
The posts listing page (`/posts/`) was displaying too much content, including post summaries. It should show only minimal metadata.

## What Should Be Displayed

According to the user's requirements, the posts listing should show **AT MOST**:
- ✅ Title (linked to post)
- ✅ Featured marker (if applicable)
- ✅ Tags
- ✅ Categories (if set)
- ✅ Date
- ✅ Reading time (if enabled)

## What Should NOT Be Displayed

- ❌ Post summary/excerpt
- ❌ Full post content
- ❌ Images
- ❌ Any other content

## Fix Applied

**Before** (`themes/apacible/layouts/_default/list.html`):
```go-html-template
<article class="post-preview">
  <h2>
    <a href="...">{{ .Title }}</a>
  </h2>
  <div class="post-meta">
    <!-- tags + date -->
  </div>
  {{ with .Summary }}
    <p>{{ . }}</p>  ← REMOVED THIS
  {{ end }}
</article>
```

**After** (Minimal):
```go-html-template
<article class="post-preview{{ if .Params.featured }} featured{{ end }}">
  <h2>
    {{ if .Params.featured }}<span class="featured-marker">ᛟ</span>{{ end }}
    <a href="{{ .Permalink }}">{{ .Title }}</a>
  </h2>
  <div class="post-meta">
    <div class="post-meta-left">
      <!-- Empty - minimal approach -->
    </div>
    <div class="post-meta-right">
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
</article>
```

## Comparison: Zola vs Hugo

### Zola Template (`templates/blog.html`)
```jinja
<article class="post-preview">
  <h2>
    {% if page.extra.featured %}<span class="featured-marker">ᛟ</span>{% endif %}
    <a href="{{ page.permalink }}">{{ page.title }}</a>
  </h2>
  <div class="post-meta">
    <div class="post-meta-left">
      <!-- Empty -->
    </div>
    <div class="post-meta-right">
      {% if page.taxonomies.tags %}
      <div class="tags">...</div>
      {% endif %}
      {% if page.date %}
      <time>{{ page.date }}</time>
      {% endif %}
    </div>
  </div>
  {% if page.summary %}  ← Zola also shows summary!
    <p>{{ page.summary }}</p>
  {% endif %}
</article>
```

**Note**: The Zola version also includes summary, but the user wants it removed for a truly minimal listing.

## Result

The posts listing now shows only:

```
┌─────────────────────────────────────────────────┐
│  ᛟ Post Title (featured marker if applicable)  │
│                           [tag1] [tag2] Jan 15  │
└─────────────────────────────────────────────────┘
```

**No content**, **no summary**, **no description** - just the essential metadata.

## Visual Layout

```
Posts
==================

ᛟ Featured Post Title ────────────────────→
                        [hugo] [zola] Nov 10, 2025

Regular Post Title ────────────────────────→
                     [example] [docs] Nov 9, 2025

Another Post ──────────────────────────────→
                           [tutorial] Nov 8, 2025
```

## Verification

```bash
# Visit posts listing
curl http://localhost:1313/posts/ | grep -A 10 'post-preview'

# Should show:
# - <h2> with title link
# - <div class="post-meta"> with tags and date
# - NO <p> tags with content
```

## CSS Styling

The minimal layout relies on these CSS classes:
- `.post-preview` - Article container
- `.featured` - Featured post modifier
- `.featured-marker` - ᛟ marker for featured posts
- `.post-meta` - Metadata container
- `.post-meta-left` - Left column (empty)
- `.post-meta-right` - Right column (tags + date)
- `.tags` - Tags container
- `.tag` - Individual tag
- `time` - Date element

## Future Enhancements (Optional)

If you want to add reading time or categories to the minimal listing:

```go-html-template
<div class="post-meta-right">
  {{ with .Params.categories }}
  <div class="categories">
    {{ range . }}
      <span class="category">{{ . }}</span>
    {{ end }}
  </div>
  {{ end }}
  
  {{ with .Params.tags }}
  <div class="tags">
    {{ range . }}
      <span class="tag">{{ . }}</span>
    {{ end }}
  </div>
  {{ end }}
  
  {{ if and .Site.Params.show_reading_time .ReadingTime }}
    <span class="reading-time">{{ .ReadingTime }} min read</span>
  {{ end }}
  
  {{ with .Date }}
    <time datetime="{{ .Format "2006-01-02" }}">
      {{ .Format ($.Site.Params.date_format | default "Jan 2, 2006") }}
    </time>
  {{ end }}
</div>
```

---

**Fix completed**: November 10, 2025  
**File modified**: `themes/apacible/layouts/_default/list.html`  
**Status**: Posts listing is now truly minimal ✅

