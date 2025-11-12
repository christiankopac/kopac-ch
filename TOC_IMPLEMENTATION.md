# Substack-Style Table of Contents Implementation

## Overview

The Table of Contents (TOC) has been redesigned to feature a Substack-style interface with horizontal line indicators. This provides a more elegant and intuitive navigation experience for long-form content.

## Key Features

### 1. **Horizontal Line Indicators**
- Initially displays as a vertical stack of horizontal lines
- Each line represents a heading in your document
- Lines are indented based on heading level (H1, H2, H3)
- Width varies by level:
  - Level 1 (H1): 28px
  - Level 2 (H2): 24px, 4px left margin
  - Level 3+ (H3): 20px, 8px left margin

### 2. **Card Expansion on Click**
- Clicking the line indicators reveals a card with full heading text
- Card features:
  - Smooth scale animation (0.95 to 1.0)
  - Semi-transparent background matching site theme
  - Subtle shadow for depth
  - Border radius of 8px
  - Auto-adjusts for dark mode
  - Closes when clicking outside or after selecting a heading
  - Toggle behavior: click again to close

### 3. **Active Section Highlighting**
- As you scroll through the document, the corresponding line highlights
- Active line changes to accent color with increased opacity
- Smooth transitions between active states

### 4. **Smooth Scroll Animation**
- Clicking on any line or heading in the expanded card smoothly scrolls to that section
- 80px offset to account for any fixed headers
- Native smooth scrolling behavior

### 5. **Responsive Design**
- Only visible on screens wider than 1401px
- Mobile and tablet devices continue to use the original collapsible TOC
- Automatically re-initializes on window resize

## Technical Implementation

### CSS Changes (`blog.css`)

1. **Line Container** (`.toc-lines`)
   - Flexbox column layout with 6px gap
   - Positioned absolutely, centered vertically
   - 8px padding for larger click target
   - Cursor pointer to indicate clickability

2. **Individual Lines** (`.toc-line`)
   - 2px height with rounded borders
   - Level-specific widths and indentation
   - Hover feedback (slight opacity increase)
   - Active states with transitions

3. **Expanded Card** (`.toc-content`)
   - Shows when `.toc-expanded` class is present
   - Positioned 220px to the left of the indicators
   - Scale transform for smooth expansion
   - Theme-aware background and shadows
   - Max height respects viewport with scrolling

### JavaScript Changes (`main.js`)

1. **Line Generation**
   - `initializeTOCLines()` dynamically creates line elements
   - Reads TOC structure and heading levels
   - Attaches click handler to toggle TOC visibility

2. **Toggle Functionality**
   - `toggleTOC()` adds/removes `.toc-expanded` class
   - `closeTOC()` removes `.toc-expanded` class
   - Click outside TOC to close automatically
   - Click inside card content doesn't close it

3. **Scroll Tracking**
   - `updateActiveTOC()` monitors scroll position
   - Updates both line and link active states
   - Throttled with requestAnimationFrame for performance

4. **Navigation**
   - `scrollToHeading()` provides smooth scroll behavior
   - Works for all heading links in the card
   - Automatically closes TOC after navigation
   - Consistent 80px offset for all navigation

## Files Modified

- `themes/apacible/static/css/blog.css`
- `themes/apacible/static/js/main.js`
- `public/css/blog.css` (deployed version)
- `public/js/main.js` (deployed version)

## Usage

To enable the TOC on a post, add `toc = true` to the front matter:

```toml
+++
title = "Your Post Title"
date = 2025-01-15
toc = true
+++
```

The TOC will automatically generate line indicators for all headings (H1-H6) in your post.

## Browser Compatibility

- Works in all modern browsers with CSS Grid/Flexbox support
- Smooth scrolling uses native `scroll-behavior` when available
- Graceful degradation for older browsers

## How to Use

1. **Opening the TOC**: Click on the horizontal lines on the right side
2. **Navigating**: Click any heading in the expanded card to jump to that section
3. **Closing**: 
   - Click the lines again to toggle closed
   - Click anywhere outside the TOC
   - Click a heading (automatically closes after navigation)
4. **Tracking**: Watch the active line highlight as you scroll through the document

## Testing

A test post has been created at `content/posts/toc-test-post.md` with multiple heading levels to demonstrate the feature.

To test:
1. Start Hugo server: `hugo server -D`
2. Navigate to the test post
3. View on a large screen (>1401px width)
4. Observe the line indicators on the right side
5. **Click** the lines to expand the TOC card
6. Click headings to navigate (card closes automatically)
7. Scroll to see active highlighting
8. Click outside to close the card

## Future Enhancements

Potential improvements:
- Configurable line colors via CSS variables
- Animation speed customization
- Option to keep card expanded/pinned
- Keyboard navigation for line indicators
- Progress indicator showing reading position


