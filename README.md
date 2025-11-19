# Christian Kopac's Website (Hugo)

Personal website built with Hugo and the **apacible** theme.

## Quick Start

### Development

```bash
# Start development server
hugo server -D

# Build for production
./build.sh
```

### Image Processing

Process images with dithering effect:

```bash
# Process all images in assets
./scripts/dither_images.sh assets --recursive

# Force re-processing
./scripts/dither_images.sh assets --recursive --overwrite
```

### Gemini Output

The site now emits Gemini-compatible `.gmi` files alongside the HTML build:

```bash
# Build both HTML and Gemini outputs
hugo

# Inspect generated Gemini files
find public -name '*.gmi'
```

You can publish the contents of `public/` as-is to a Gemini capsule (the build mirrors the same permalink structure, with `index.gmi` files for each page).

### Atom Feed

The primary Atom feed is published at `https://christiankopac.com/atom.xml`. It aggregates updates from the `posts` and `picks` sections. Each section also exposes its own feed (for example, `https://christiankopac.com/posts/atom.xml`).

## Theme

This site uses the [apacible](themes/apacible/) Hugo theme, which focuses on:
- Accessibility
- Smolweb principles
- Dithered image aesthetics
- Minimal JavaScript
- Fast loading times

See [themes/apacible/README.md](themes/apacible/README.md) for theme documentation.

## Content Structure

- `content/posts/` - Blog posts
- `content/collections/` - Collection pages
- `content/consumed/` - Books, movies, etc.
- `content/about/` - About page
- `data/` - Structured data for collections and galleries

## Building

```bash
# Full build with image processing
./build.sh

# Build without dithering (faster for testing)
./build.sh --skip-dither

# Force re-dither all images
./build.sh --force
```

## Requirements

- Hugo v0.100.0 or later
- ImageMagick (for image dithering)

## License

Content: CC BY-NC-ND 4.0  
Code: MIT

