# Zola to Hugo Migration Notes

This document summarizes the migration of the **apacible** theme from Zola to Hugo.

## Migration Date
November 10, 2025

## What Was Migrated

### ✅ Templates
- [x] Base layout (`_base.html` → `baseof.html`)
- [x] Home page (`home.html` → `index.html`)
- [x] Single page/post (`post.html` → `single.html`)
- [x] List/blog pages (`blog.html` → `list.html`)
- [x] Section pages (`prose.html` → `section.html`)
- [x] Taxonomy pages (`taxonomy.html`, `term.html`)
- [x] 404 page
- [x] Footer partial
- [x] Mermaid partial

### ✅ Shortcodes
- [x] `img` - Advanced image with dithering support
- [x] `img_compare` - Side-by-side image comparison
- [x] `gallery` - Image grid with lightbox
- [x] `collection` - Structured data display
- [x] `callout` - Callout boxes
- [x] `note`, `tip`, `important`, `warning`, `caution` - Specific callouts
- [x] `detail` - Collapsible sections
- [x] `mermaid` - Diagrams
- [x] `quote` - Styled quotes

### ✅ Styles (CSS)
- [x] All 18 CSS files copied from Zola
- [x] Base styles (tokens, theme, typography)
- [x] Layout styles (layout, navigation, footer)
- [x] Component styles (blog, prose, collection)
- [x] Feature styles (anchors, syntax, responsive-images, shortcodes)
- [x] Interactive styles (cursor, back-to-top)

### ✅ JavaScript
- [x] `main.js` - Core functionality
- [x] `lightense.min.js` - Image lightbox

### ✅ Scripts
- [x] `dither_images.sh` - Image dithering script
- [x] `build.sh` - Build script with image processing

### ✅ Configuration
- [x] `hugo.toml` - Site configuration
- [x] `theme.toml` - Theme metadata
- [x] Content structure setup

### ✅ Documentation
- [x] Theme README
- [x] Image processing guide
- [x] Site README
- [x] Migration notes (this file)

## Key Changes

### Image Processing

**Zola → Hugo Differences:**
- Zola: `resize_image()` function with multi-format
- Hugo: `.Resize`, `.Fit`, `.Fill` methods on resources
- **AVIF removed** per user request (smolweb)
- **WebP retained** with fallback to JPEG/PNG

### Template Syntax

| Feature | Zola (Tera) | Hugo (Go) |
|---------|-------------|-----------|
| Variables | `{{ page.title }}` | `{{ .Title }}` |
| Site config | `{{ config.extra.var }}` | `{{ .Site.Params.var }}` |
| Conditionals | `{% if condition %}` | `{{ if condition }}` |
| Loops | `{% for item in items %}` | `{{ range .Pages }}` |
| Blocks | `{% block name %}` | `{{ block "name" . }}` |
| Partials | `{% include "file.html" %}` | `{{ partial "file.html" . }}` |
| URLs | `{{ get_url(path="...") }}` | `{{ "..." \| relURL }}` |
| Safe HTML | `{{ content \| safe }}` | `{{ .Content }}` |

### Shortcode Syntax

| Feature | Zola | Hugo |
|---------|------|------|
| Parameters | `{% img(src="...") %}` | `{{< img src="..." >}}` |
| Body content | `{{ body }}` | `{{ .Inner }}` |
| Markdown in body | `{{ body \| markdown }}` | `{{ .Inner \| .Page.RenderString }}` |

### File Structure

```
Zola:                          Hugo:
├── templates/                 ├── layouts/
│   ├── _base.html            │   ├── _default/baseof.html
│   ├── home.html             │   ├── index.html
│   ├── post.html             │   ├── _default/single.html
│   ├── blog.html             │   ├── _default/list.html
│   ├── shortcodes/           │   ├── shortcodes/
│   └── macros/ (unused)      │   └── partials/
├── static/                    ├── static/ (same)
│   ├── css/                  │   ├── css/
│   └── js/                   │   └── js/
├── content/ (same)            ├── content/ (same)
└── config.toml               └── hugo.toml
```

## Breaking Changes

### For Content Authors

1. **Shortcode syntax changed**:
   ```markdown
   # Zola
   {{ img(src="photo.jpg", alt="Photo", dithered=true) }}
   
   # Hugo
   {{< img src="photo.jpg" alt="Photo" dithered="true" >}}
   ```

2. **Data file paths**:
   ```markdown
   # Zola
   {{ collection(file="projects.toml") }}
   # Looks in: content/collections/projects.toml
   
   # Hugo
   {{< collection file="projects.toml" >}}
   # Looks in: data/projects.toml
   ```

3. **Front matter params**:
   ```yaml
   # Zola
   [extra]
   math = true
   
   # Hugo
   math: true
   # (directly in front matter, not nested)
   ```

### For Theme Developers

1. **No macros**: Zola's macros system doesn't exist in Hugo. Use partials instead.

2. **Data loading**: 
   - Zola: `load_data(path="content/file.toml")`
   - Hugo: `index .Site.Data "file"`

3. **Image processing**: Completely different API (see above)

4. **Taxonomy access**:
   - Zola: `get_taxonomy(kind="tags")`
   - Hugo: Automatic via `.Data.Terms`

## What Was NOT Migrated

- ❌ Zola-specific features not applicable to Hugo
- ❌ Build-time computed values (Hugo handles differently)
- ❌ AVIF image format (dropped per requirements)

## Known Limitations

1. **Gallery data**: Must be in `data/` directory in Hugo, not `content/`
2. **Collection data**: Same as gallery
3. **Page resources**: Slightly different behavior between Zola and Hugo
4. **TOC structure**: Hugo's TOC HTML is different from Zola's

## Testing Checklist

Before going live, test:

- [ ] Home page renders correctly
- [ ] Blog listing shows posts
- [ ] Single post pages work
- [ ] Taxonomy pages (tags, categories)
- [ ] Dark mode toggle
- [ ] Image processing (dithered/original toggle)
- [ ] Gallery lightbox
- [ ] Collection displays
- [ ] All shortcodes render
- [ ] Mermaid diagrams (if used)
- [ ] Math rendering (if used)
- [ ] RSS feed generation
- [ ] 404 page
- [ ] Mobile responsive
- [ ] No JavaScript fallback

## Migration Steps for Existing Content

1. **Update shortcode syntax** in all `.md` files:
   ```bash
   # Find Zola shortcodes
   grep -r "{{.*(" content/
   
   # Manual replacement needed (syntax too different for simple sed)
   ```

2. **Move data files**:
   ```bash
   # Move from content/ to data/
   mv content/collections/*.toml data/
   ```

3. **Update front matter**:
   ```bash
   # Remove [extra] nesting if present
   # Update any Zola-specific front matter
   ```

4. **Run image dithering**:
   ```bash
   ./scripts/dither_images.sh assets --recursive
   ./scripts/dither_images.sh content --recursive
   ```

5. **Test build**:
   ```bash
   ./build.sh
   hugo server -D
   ```

## Build Commands

```bash
# Development
hugo server -D

# Production build (with image processing)
./build.sh

# Production build (skip image processing)
./build.sh --skip-dither

# Force re-dither all images
./build.sh --force
```

## Resources

- [Hugo Documentation](https://gohugo.io/documentation/)
- [Hugo Image Processing](https://gohugo.io/content-management/image-processing/)
- [Hugo Shortcodes](https://gohugo.io/content-management/shortcodes/)
- [Hugo Templates](https://gohugo.io/templates/)

## Support

For issues with the theme migration:
1. Check this document
2. Review theme README
3. Check Hugo documentation
4. Open an issue on the theme repo

## License

Theme code: MIT  
Original design: Christian Kopac

