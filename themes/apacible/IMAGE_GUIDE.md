# Image Processing Guide for Apacible Hugo Theme

This guide explains how to use the image processing system with dithered images in the Apacible Hugo theme.

## Overview

The theme supports:
- ✅ **WebP format** (modern, efficient)
- ✅ **Dithered images** (lo-fi aesthetic)
- ✅ **Toggle to original** (via caption link)
- ✅ **Lazy loading** (performance)
- ❌ **No AVIF** (dropped per smolweb requirements)

## Quick Start

### 1. Place Images

Put images in one of these locations:
- `assets/` - Processed by Hugo (recommended)
- `static/` - Served as-is (for pre-processed images)
- `content/posts/` - Page-relative images

### 2. Generate Dithered Versions

```bash
# Dither images in assets
./scripts/dither_images.sh assets --recursive

# Dither images in content
./scripts/dither_images.sh content --recursive

# Force re-dither existing files
./scripts/dither_images.sh assets --recursive --overwrite
```

This creates `image_dithered.jpg` for each `image.jpg`.

### 3. Use in Content

#### Basic Image
```markdown
{{< img src="assets/photo.jpg" alt="Description" >}}
```

#### Dithered Image with Toggle
```markdown
{{< img src="assets/photo.jpg" alt="Description" caption="My photo (click to toggle)" dithered="true" >}}
```

## Shortcode Reference

### img

**Parameters:**

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `src` | Yes | - | Path to image |
| `alt` | Yes | - | Alt text for accessibility |
| `caption` | No | - | Figure caption (enables toggle for dithered) |
| `width` | No | 800 | Target width in pixels |
| `height` | No | 600 | Target height in pixels |
| `op` | No | "fit" | Resize operation: "resize", "fit", "fill" |
| `quality` | No | 85 | JPEG/WebP quality (1-100) |
| `dithered` | No | false | Use dithered version |
| `show_original` | No | true | Show toggle link (for dithered) |

**Examples:**

```markdown
<!-- Standard image -->
{{< img src="assets/photo.jpg" alt="Beach sunset" >}}

<!-- Dithered with caption toggle -->
{{< img src="assets/photo.jpg" alt="Beach sunset" caption="Sunset at the beach" dithered="true" >}}

<!-- Custom size -->
{{< img src="assets/photo.jpg" alt="Diagram" width="1200" height="800" op="fit" >}}

<!-- High quality -->
{{< img src="assets/artwork.jpg" alt="Digital art" quality="95" >}}
```

### img_compare

Display dithered and original side-by-side.

**Parameters:**

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `src` | Yes | - | Path to image |
| `alt` | Yes | - | Alt text |
| `caption` | No | - | Overall caption |
| `width` | No | 800 | Target width |
| `height` | No | 600 | Target height |
| `op` | No | "fit" | Resize operation |
| `quality` | No | 75 | Quality (1-100) |

**Example:**

```markdown
{{< img_compare src="assets/photo.jpg" alt="Photo comparison" caption="Dithering Effect" >}}
```

### gallery

Display grid of images with lightbox.

**Parameters:**

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `file` | Yes | - | Data file name (in `data/`) |
| `columns` | No | 3 | Number of columns |
| `thumb_width` | No | 300 | Thumbnail width |
| `thumb_height` | No | 200 | Thumbnail height |
| `full_width` | No | 2400 | Full image width |
| `full_height` | No | 1600 | Full image height |

**Example:**

```markdown
{{< gallery file="my-gallery.toml" columns="3" >}}
```

**Data file** (`data/my-gallery.toml`):

```toml
[[images]]
src = "/assets/image1.jpg"
alt = "Description"
caption = "Optional caption"

[[images]]
src = "/assets/image2.jpg"
alt = "Description"
caption = "Another photo"
```

## Image Locations

### Assets Directory (Recommended)

```
assets/
  photos/
    beach.jpg
    beach_dithered.jpg  ← auto-generated
```

Usage:
```markdown
{{< img src="photos/beach.jpg" alt="Beach" dithered="true" >}}
```

Hugo processes these at build time.

### Static Directory

```
static/
  images/
    icon.png
    icon_dithered.png
```

Usage:
```markdown
{{< img src="images/icon.png" alt="Icon" dithered="true" >}}
```

Files served as-is, no Hugo processing.

### Page Resources

```
content/posts/my-post/
  index.md
  image.jpg
  image_dithered.jpg
```

Usage in `my-post/index.md`:
```markdown
{{< img src="image.jpg" alt="Local image" dithered="true" >}}
```

## Resize Operations

### fit (default)
Fits within dimensions, preserves aspect ratio:
```markdown
{{< img src="photo.jpg" alt="Photo" width="800" height="600" op="fit" >}}
```

### resize
Resizes to exact dimensions (may distort):
```markdown
{{< img src="photo.jpg" alt="Photo" width="800" height="600" op="resize" >}}
```

### fill
Crops to exact dimensions:
```markdown
{{< img src="photo.jpg" alt="Photo" width="800" height="600" op="fill" >}}
```

## Dithering Effect

The dithering script applies:
- **Grayscale conversion**
- **High contrast** (5% stretch)
- **Posterization** (2 levels)
- **Ordered 4x4 dithering** (visible pattern)
- **Normalization** (maximize contrast)

Result: Chunky, high-contrast black & white lo-fi aesthetic.

## Build Integration

The `build.sh` script automatically dithers images:

```bash
# Full build with dithering
./build.sh

# Skip dithering (faster for testing)
./build.sh --skip-dither

# Force re-dither all images
./build.sh --force
```

## Workflow

1. **Add images** to your content/assets:
   ```bash
   cp photo.jpg assets/photos/
   ```

2. **Generate dithered versions**:
   ```bash
   ./scripts/dither_images.sh assets --recursive
   ```

3. **Use in markdown**:
   ```markdown
   {{< img src="photos/photo.jpg" alt="Photo" caption="My photo" dithered="true" >}}
   ```

4. **Build site**:
   ```bash
   ./build.sh
   ```

## Troubleshooting

### Image not found

- Check path is correct (relative to `assets/`, `static/`, or page)
- Ensure dithered version exists (run dither script)
- For `assets/`, path should not include `/assets/` prefix

### Dithered version not showing

- Run: `./scripts/dither_images.sh your-directory --overwrite`
- Verify `image_dithered.jpg` exists next to `image.jpg`
- Check filename matches exactly (case-sensitive)

### ImageMagick not found

Install ImageMagick:
```bash
# Ubuntu/Debian
sudo apt install imagemagick

# macOS
brew install imagemagick

# Verify
magick --version  # or: convert --version
```

## Performance Tips

1. **Dither before build**: Run dither script before `hugo` command
2. **Appropriate sizes**: Don't use 4K images for thumbnails
3. **WebP compression**: Already optimized in shortcodes
4. **Lazy loading**: Enabled by default in shortcodes

## Examples

See the theme demo for live examples:
- Standard images
- Dithered toggles
- Image comparisons
- Gallery grids

## Format Support

| Format | Support | Notes |
|--------|---------|-------|
| WebP | ✅ Yes | Auto-generated, with fallback |
| AVIF | ❌ No | Dropped per smolweb requirements |
| JPEG | ✅ Yes | Fallback format |
| PNG | ✅ Yes | For graphics with transparency |

## Browser Support

- Modern browsers: Full support (WebP + lazy loading)
- Older browsers: Fallback to JPEG/PNG
- No JavaScript: Images still display (no toggle)

