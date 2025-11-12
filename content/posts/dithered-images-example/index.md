+++
title = "Example: Dithered Images"
date = 2024-11-10
description = "Demonstration of dithered image processing with toggle functionality"
featured = true


tags = ["apacible", "images", "dithering", "example"]
categories = ["Examples"]
+++

This page demonstrates the dithered image functionality. Images are shown in their dithered form by default, which reduces file size while maintaining a unique aesthetic. Click "**→ show original**" link below each image to toggle between dithered and original versions.

## Single Dithered Image

{{< img src="musician.jpeg" alt="Musician playing" caption="Musician in performance" dithered="true" width="800" height="600" >}}

The dithering effect reduces the image to a limited color palette, creating a distinctive retro aesthetic while significantly reducing file size—perfect for the smolweb philosophy.

## Another Example

{{< img src="trees.jpg" alt="Forest landscape" caption="Trees in natural light" dithered="true" width="800" height="600" >}}

## Side-by-Side Comparison

For direct comparison, here's the same image shown in both forms:

{{< img_compare 
    dithered="trees_dithered.jpg" 
    original="trees.jpg" 
    alt="Trees comparison" 
    caption="Left: Dithered | Right: Original"
>}}

## How It Works

### Processing

Images are dithered using ImageMagick's Floyd-Steinberg algorithm:

```bash
magick input.jpg \
    -colorspace gray \
    -ordered-dither o8x8,8 \
    output_dithered.jpg
```

### Toggle Mechanism

The shortcode generates HTML with both versions:

```html
<figure class="responsive-image has-toggle">
  <div class="image-container">
    <!-- Dithered version (shown by default) -->
    <img class="dithered-img active" src="dithered.jpg" alt="...">
    
    <!-- Original version (hidden) -->
    <picture class="original-img" style="display: none;">
      <source srcset="original.webp" type="image/webp">
      <img src="original.jpg" alt="...">
    </picture>
  </div>
  
  <figcaption>
    Caption text
    <span class="toggle-original">→ show original image</span>
  </figcaption>
</figure>
```

JavaScript toggles visibility when the link is clicked:

```javascript
function toggleOriginal(element) {
  const container = element.closest('.image-container') || 
                    element.closest('figure').querySelector('.image-container');
  const dithered = container.querySelector('.dithered-img');
  const original = container.querySelector('.original-img');
  const showText = element.querySelector('.show-text');
  const hideText = element.querySelector('.hide-text');
  
  // Toggle visibility
  dithered.classList.toggle('active');
  original.style.display = dithered.classList.contains('active') ? 'none' : 'block';
  
  // Toggle link text
  showText.style.display = dithered.classList.contains('active') ? 'inline' : 'none';
  hideText.style.display = dithered.classList.contains('active') ? 'none' : 'inline';
}
```

## Benefits

### Performance
- **Smaller Files**: Dithered images are typically 40-60% smaller
- **Faster Loading**: Reduced bandwidth usage
- **Progressive Enhancement**: Original available on demand

### Aesthetics
- **Unique Style**: Distinctive visual identity
- **Retro Appeal**: Nostalgic 8-bit aesthetic
- **High Contrast**: Often more readable

### Philosophy
- **Smolweb Aligned**: Minimal data transfer
- **User Choice**: Original available when needed
- **Accessibility**: High contrast improves readability

## Usage

To use dithered images in your posts:

```markdown
{{</* img 
    src="image.jpg" 
    alt="Description" 
    caption="Your caption here" 
    dithered="true" 
    width="800" 
    height="600" 
*/>}}
```

**Note**: You must have both the original image and its dithered version (with `_dithered` suffix) in your assets folder. Use the `scripts/dither_images.sh` script to generate dithered versions automatically.

## Technical Details

- **Algorithm**: Floyd-Steinberg ordered dithering (8x8 pattern)
- **Color Space**: Grayscale conversion before dithering
- **Format**: Supports JPEG, PNG, WebP
- **Responsive**: Hugo's image processing creates optimized sizes
- **Fallback**: If dithered version missing, shows original only

---

The dithered image approach balances modern web performance with a unique aesthetic, perfect for sites embracing the smolweb philosophy while maintaining visual interest.

