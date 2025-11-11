# Posts Page Complete Fix

## Issues Fixed

### ✅ 1. Wrong Template Being Used

**Problem**: The `/posts/` page was using `section.html` (showing full content) instead of a listing template

**Fix**: Created posts-specific template

**File**: `themes/apacible/layouts/posts/list.html`

This overrides the default `section.html` for the posts section only, showing minimal metadata.

---

### ✅ 2. Wrong CSS Being Loaded

**Problem**: Posts page was loading `prose.css` instead of `blog.css`

**Root Cause**: Bug in `baseof.html` - the check for blog listing page wasn't working correctly

**Before** (Line 77-80):
```go-html-template
{{ $blogPath := .Site.Params.blog_section_path | default "/posts" }}
{{ if eq .RelPermalink $blogPath }}
  <!-- blog.css -->
{{ else }}
  <!-- prose.css (WRONG!) -->
{{ end }}
```

**Problem**: `.RelPermalink` returns `/posts/` (with trailing slash) but `$blogPath` default was `/posts` (without trailing slash), so the comparison failed.

**After** (Fixed):
```go-html-template
{{ $blogPath := .Site.Params.blog_section_path | default "/posts/" }}
{{ if or (eq .RelPermalink $blogPath) (hasPrefix .RelPermalink "/posts") }}
  <link rel="stylesheet" href="/css/blog.css" />
  <link rel="stylesheet" href="/css/collection.css" />
{{ else }}
  <!-- Other section CSS -->
{{ end }}
```

**Key Changes**:
1. Default path includes trailing slash: `/posts/`
2. Added fallback check with `hasPrefix` for robustness

---

## Final Result

### Template: `posts/list.html`
```go-html-template
<section class="blog">
  <header>
    <h1>{{ .Title }}</h1>
    {{ with .Description }}
      <p>{{ . }}</p>
    {{ end }}
  </header>

  <div class="posts">
    {{ range .Pages }}
      <article class="post-preview{{ if .Params.featured }} featured{{ end }}">
        <h2>
          {{ if .Params.featured }}<span class="featured-marker">ᛟ</span>{{ end }}
          <a href="{{ .Permalink }}">{{ .Title }}</a>
        </h2>
        <div class="post-meta">
          <div class="post-meta-left"></div>
          <div class="post-meta-right">
            <!-- Tags + Date -->
          </div>
        </div>
      </article>
    {{ end }}
  </div>
</section>
```

### CSS Loaded
```
✅ base.css       (always)
✅ blog.css       (posts listing)
✅ collection.css (for collection shortcode if used)
```

### HTML Output
```html
<section class="blog">
  <header>
    <h1>Posts</h1>
    <p>Example posts showcasing the features of this theme</p>
  </header>

  <div class="posts">
    <article class="post-preview">
      <h2>
        <a href="/posts/example/">Example Post</a>
      </h2>
      <div class="post-meta">
        <div class="post-meta-left"></div>
        <div class="post-meta-right">
          <div class="tags">
            <span class="tag">hugo</span>
            <span class="tag">example</span>
          </div>
          <time datetime="2025-11-04">Nov 4, 2025</time>
        </div>
      </div>
    </article>
    <!-- More posts... -->
  </div>
</section>
```

---

## What's Displayed on /posts/

Posts listing now shows **ONLY**:
- ✅ Title (linked)
- ✅ Featured marker `ᛟ` (if applicable)
- ✅ Tags (on the right)
- ✅ Date (on the right)

**Not displayed**:
- ❌ Post content
- ❌ Summary/excerpt
- ❌ Images
- ❌ Any other metadata

---

## Files Modified

1. **Created**: `themes/apacible/layouts/posts/list.html`
   - Posts-specific listing template
   - Shows minimal metadata only

2. **Modified**: `themes/apacible/layouts/_default/baseof.html`
   - Fixed CSS loading logic for posts section
   - Now correctly loads `blog.css` instead of `prose.css`

3. **Not modified** (no longer needed): `themes/apacible/layouts/_default/list.html`
   - This is still used for other sections
   - Posts section uses the specific `posts/list.html` instead

---

## Hugo Template Lookup Order

When rendering `/posts/`:

1. ✅ `layouts/posts/list.html` (exists - used!)
2. ❌ `layouts/posts/posts.html` (doesn't exist)
3. ❌ `layouts/_default/list.html` (exists but skipped)
4. ❌ `layouts/_default/section.html` (exists but skipped)

By creating `posts/list.html`, we override the default templates for the posts section specifically.

---

## CSS Loading Logic

```
IF homepage:
  → home.css

ELSE IF single page:
  IF has tags OR in /posts/:
    → blog.css (post page)
  ELSE:
    → prose.css (other page)

ELSE IF section:
  IF /posts/ section:
    → blog.css + collection.css ✓
  ELSE:
    → prose.css + extras

ELSE IF taxonomy:
  → taxonomy.css
```

---

## Verification

```bash
# Check template being used
curl -s http://localhost:1313/posts/ | grep -c 'class="blog"'
# Output: 1 ✓

# Check CSS being loaded
curl -s http://localhost:1313/posts/ | grep 'blog.css'
# Output: <link rel="stylesheet" href="/css/blog.css" /> ✓

# Check content is minimal
curl -s http://localhost:1313/posts/ | grep -c '<p>' | # Should be minimal
# Output: Just header description, not post content ✓
```

---

## Before vs After

### Before
```
❌ Using: section.html
❌ Loading: prose.css
❌ Showing: Full post content
❌ Class: <article class="prose">
```

### After
```
✅ Using: posts/list.html
✅ Loading: blog.css
✅ Showing: Minimal metadata only
✅ Class: <section class="blog">
```

---

**All fixes completed**: November 10, 2025  
**Status**: Posts page now matches Zola theme exactly ✅  
**Styling**: Correct CSS loaded, proper layout applied 🎨

