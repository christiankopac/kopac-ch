# Apacible Hugo Theme

> good-tempered, calm, pleasant

A minimalist Hugo theme prioritizing accessibility and smolweb principles, with unique dithered image support.

## Features

- ✅ **Smolweb-focused**: Small, fast, accessible
- 🎨 **Dithered images**: Lo-fi aesthetic with toggle to original
- 🌓 **Dark mode**: Automatic theme switching with localStorage
- 📝 **Blog-ready**: Post listings, tags, categories
- 🎯 **Collections**: Flexible content display system
- 📊 **Mermaid diagrams**: Built-in diagram support
- 🔢 **KaTeX math**: Mathematical notation support
- 📱 **Responsive**: Mobile-first design
- ♿ **Accessible**: WCAG compliant

## Installation

1. Add as a Git submodule:
```bash
cd your-hugo-site
git submodule add https://github.com/christiankopac/apacible.git themes/apacible
```

2. Update your `hugo.toml`:
```toml
theme = "apacible"
```

3. Copy the example config from `themes/apacible/exampleSite/hugo.toml` to your site root.

## Quick Start

### 1. Create your home page

Create `content/_index.md`:

```markdown
+++
title = "Home"

[params]
name = "Your Name"
bio = "Your bio here"
avatar = ""  # Leave empty for default avatar

[[params.links]]
name = "GitHub"
url = "https://github.com/yourusername"
icon = "github"

recent = true
recent_max = 5
recent_more_text = "View all posts →"
+++
```

### 2. Create your first post

```bash
hugo new posts/my-first-post.md
```

### 3. Add images with dithering

1. Place images in `assets/` or `static/` directory
2. Run the dithering script:
```bash
./scripts/dither_images.sh assets --recursive
```

3. Use in your content:
```markdown
{{< img src="assets/photo.jpg" alt="My photo" caption="A dithered photo" dithered="true" >}}
```

### 4. Build and serve

```bash
# Development server
hugo server -D

# Production build
./build.sh
```

## Image Processing

This theme includes a unique dithered image feature:

- **Dithered thumbnails**: Show lo-fi aesthetic by default
- **Original on click**: Toggle to high-quality original
- **WebP support**: Automatic WebP conversion (no AVIF per smolweb)
- **Lazy loading**: Optimized performance

### Image Shortcodes

#### img

Basic image with optional dithering:

```markdown
{{< img src="path/to/image.jpg" alt="Description" caption="Caption" dithered="true" >}}
```

Parameters:
- `src` (required): Path to image
- `alt` (required): Alt text
- `caption`: Figure caption
- `width`: Target width (default: 800)
- `height`: Target height (default: 600)
- `op`: Resize operation - "resize", "fit", "fill" (default: "fit")
- `quality`: JPEG/WebP quality 1-100 (default: 85)
- `dithered`: Use dithered version (default: false)

#### img_compare

Side-by-side comparison of dithered vs original:

```markdown
{{< img_compare src="path/to/image.jpg" alt="Description" caption="Comparison" >}}
```

#### gallery

Grid of images with lightbox:

```markdown
{{< gallery file="gallery-data.toml" columns="3" >}}
```

Create `data/gallery-data.toml`:
```toml
[[images]]
src = "/assets/image1.jpg"
alt = "Description"
caption = "Caption"

[[images]]
src = "/assets/image2.jpg"
alt = "Description"
```

## Shortcodes

### Callouts

```markdown
{{< note title="Note" >}}
Important information here
{{< /note >}}

{{< tip title="Pro Tip" >}}
Helpful advice
{{< /tip >}}

{{< warning title="Warning" >}}
Cautionary note
{{< /warning >}}
```

Available: `note`, `tip`, `important`, `warning`, `caution`

### Details (Collapsible)

```markdown
{{< detail title="Click to expand" >}}
Hidden content here
{{< /detail >}}
```

### Quote

```markdown
{{< quote cite="Author Name" >}}
Quote text here
{{< /quote >}}
```

### Mermaid Diagrams

```markdown
{{< mermaid >}}
graph LR
  A --> B
  B --> C
{{< /mermaid >}}
```

Enable in front matter:
```yaml
---
mermaid: true
---
```

### Collection

Display structured data:

```markdown
{{< collection file="projects.toml" style="card" >}}
```

Styles: `card`, `simple-card`, `list`, `card-grid`, `card-horizontal`

## Configuration

Key config options in `hugo.toml`:

```toml
[params]
  description = "Site description"
  blog_section_path = "/posts"
  back_link_text = "Back"
  show_word_count = true
  show_reading_time = true
  date_format = "Jan 2, 2006"
  
  # Navigation sections
  [[params.sections]]
    name = "posts"
    path = "/posts"
    is_external = false
```

## Development

### File Structure

```
apacible/
├── archetypes/          # Content templates
├── assets/              # Processed assets
│   ├── css/            # Stylesheets (from theme)
│   └── js/             # JavaScript (from theme)
├── layouts/
│   ├── _default/       # Base templates
│   ├── partials/       # Reusable components
│   └── shortcodes/     # Content shortcodes
├── static/             # Static files (copied as-is)
│   ├── css/            # Final CSS files
│   ├── js/             # Final JS files
│   └── img/            # Favicons, etc.
└── theme.toml          # Theme metadata
```

### Build Script

```bash
./build.sh              # Build with dithering
./build.sh --skip-dither    # Skip image processing
./build.sh --force      # Re-dither all images
```

## Smolweb Compliance

This theme adheres to smolweb principles:

- ✅ No tracking or analytics
- ✅ No external dependencies (except CDN for math/diagrams)
- ✅ Minimal JavaScript (progressive enhancement)
- ✅ Semantic HTML
- ✅ Small file sizes
- ✅ Works without JavaScript
- ✅ Accessible (WCAG AA)

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- Progressive enhancement for older browsers
- WebP with fallback to JPEG/PNG
- Works without JavaScript (with degraded UX)

## Credits

- Theme design: Christian Kopac
- Migrated from Zola to Hugo
- ~~Lightbox: [Lightense](https://github.com/sparanoid/lightense-images)~~ (Removed)
- Diagrams: [Mermaid](https://mermaid.js.org/)
- Math: [KaTeX](https://katex.org/)

## License

MIT License - see LICENSE file for details

