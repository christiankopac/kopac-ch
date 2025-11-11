# Quick Start Guide - Apacible Hugo Theme

Your Zola theme has been successfully migrated to Hugo! 🎉

## What's Ready

✅ **Complete Hugo theme** in `themes/apacible/`  
✅ **All templates migrated** (base, home, post, list, section, taxonomy)  
✅ **All shortcodes** (img, gallery, collection, callouts, etc.)  
✅ **All CSS & JavaScript** files  
✅ **Image dithering support** (WebP only, no AVIF)  
✅ **Build scripts** with automated image processing  
✅ **Documentation** (README, guides, migration notes)

## Immediate Next Steps

### 1. Test the Setup

```bash
cd /home/christian/src/my_domains/christiankopac_com__hugo

# Start development server
hugo server -D
```

Open http://localhost:1313 in your browser.

### 2. Add Your Content

The home page starter is at `content/_index.md`. Edit it to add your info:

```yaml
+++
title = "Home"

[params]
name = "Your Name"
bio = "Your bio"
# ... etc
+++
```

### 3. Process Images

If you have images to dither:

```bash
# Copy images to assets/
cp -r /path/to/images/* assets/

# Generate dithered versions
./scripts/dither_images.sh assets --recursive
```

### 4. Create Content

```bash
# Create a new post
hugo new posts/my-first-post.md

# Edit it
vim content/posts/my-first-post.md
```

### 5. Build for Production

```bash
# Full build with image processing
./build.sh

# Output will be in public/
ls public/
```

## Important Differences from Zola

### Shortcode Syntax Changed

**Before (Zola):**
```markdown
{{ img(src="photo.jpg", alt="Photo", dithered=true) }}
```

**Now (Hugo):**
```markdown
{{< img src="photo.jpg" alt="Photo" dithered="true" >}}
```

### Data Files Moved

Collection and gallery data files should now be in `data/` directory:
- Before: `content/collections/projects.toml`
- Now: `data/projects.toml`

### Front Matter

Remove `[extra]` nesting in front matter:

**Before (Zola):**
```yaml
+++
title = "Post"
[extra]
math = true
mermaid = true
+++
```

**Now (Hugo):**
```yaml
+++
title = "Post"
math = true
mermaid = true
+++
```

## Key Features

### Image Shortcodes

```markdown
<!-- Standard image -->
{{< img src="assets/photo.jpg" alt="Beach" >}}

<!-- Dithered with toggle -->
{{< img src="assets/photo.jpg" alt="Beach" caption="Beach photo" dithered="true" >}}

<!-- Comparison -->
{{< img_compare src="assets/photo.jpg" alt="Before/After" >}}

<!-- Gallery -->
{{< gallery file="my-gallery.toml" columns="3" >}}
```

### Callouts

```markdown
{{< note title="Note" >}}
Important information
{{< /note >}}

{{< tip title="Pro Tip" >}}
Helpful advice
{{< /tip >}}
```

### Collections

```markdown
{{< collection file="projects.toml" style="card" >}}
```

### Diagrams

```markdown
{{< mermaid >}}
graph LR
  A --> B
{{< /mermaid >}}
```

(Set `mermaid: true` in front matter)

## Configuration

Edit `hugo.toml` to customize:

```toml
[params]
  description = "Your site description"
  blog_section_path = "/posts"
  show_word_count = true
  show_reading_time = true
  
  [[params.sections]]
    name = "posts"
    path = "/posts"
```

## File Structure

```
your-hugo-site/
├── content/          # Your markdown files
│   ├── _index.md    # Home page
│   ├── posts/       # Blog posts
│   ├── about/       # About section
│   └── ...
├── data/            # Data files for collections
├── assets/          # Images (Hugo processes these)
├── static/          # Static files (copied as-is)
├── themes/
│   └── apacible/    # The theme (migrated from Zola)
├── hugo.toml        # Site configuration
├── build.sh         # Build script
└── scripts/
    └── dither_images.sh
```

## Common Tasks

### Add a new blog post
```bash
hugo new posts/my-post.md
```

### Add images with dithering
```bash
# 1. Add image
cp photo.jpg assets/images/

# 2. Generate dithered version
./scripts/dither_images.sh assets --recursive

# 3. Use in markdown
# {{< img src="images/photo.jpg" alt="..." dithered="true" >}}
```

### Build and deploy
```bash
# Build
./build.sh

# Deploy (copy public/ to your server)
rsync -avz public/ user@server:/var/www/html/
```

## Getting Help

- **Theme README**: `themes/apacible/README.md`
- **Image Guide**: `themes/apacible/IMAGE_GUIDE.md`
- **Migration Notes**: `MIGRATION_NOTES.md`
- **Hugo Docs**: https://gohugo.io/documentation/

## What's Different from Zola

1. **Build command**: `hugo` instead of `zola build`
2. **Dev server**: `hugo server` instead of `zola serve`
3. **Shortcodes**: Different syntax (see above)
4. **Image processing**: Different API (handled by theme)
5. **Data files**: In `data/` not `content/`

## Troubleshooting

### "Theme not found"
- Ensure `theme = "apacible"` is in `hugo.toml`
- Check theme exists at `themes/apacible/`

### Images not showing
- For `assets/`, use path without `/assets/` prefix
- For `static/`, use full path from root
- Ensure dithered versions exist (run dither script)

### Shortcodes not working
- Check syntax: `{{< shortcode >}}` not `{{ shortcode() }}`
- Parameters: `name="value"` not `name=value`

### Build fails
- Check Hugo version: `hugo version` (need v0.100.0+)
- Check ImageMagick: `magick --version` or `convert --version`

## Next Steps

1. ✅ Test the basic setup
2. ✅ Customize `content/_index.md`
3. ✅ Update `hugo.toml` with your info
4. ✅ Migrate your existing content
5. ✅ Process your images
6. ✅ Build and deploy

Enjoy your Hugo site! 🚀

