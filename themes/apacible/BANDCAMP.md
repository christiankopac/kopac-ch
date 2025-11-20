# Bandcamp Shortcode

Embed Bandcamp tracks and albums in your playlists with different styles.

## Getting the Embed URL

1. Go to your track/album on Bandcamp
2. Click **Share / Embed**
3. Click **Embed this track/album**
4. Copy the `src` URL from the iframe code

The URL will look like: `https://bandcamp.com/EmbeddedPlayer/album=3553509833`

## Usage

### Standard (Default)
220px width, floats left by default, text wraps around like a newspaper. 120px height, small artwork, no tracklist.

```markdown
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore.

{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=3553509833" >}}

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Text will wrap around the embed on the right side, creating a newspaper-style layout.
```

### Slim
220px width, floats left by default, text wraps around. 42px height, minimal player.

```markdown
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor.

{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=3553509833" style="slim" >}}

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore. Text flows around the compact player.
```

### Artwork Only
350x350px, just the album artwork

```markdown
{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=3553509833" style="artwork" >}}
```

### Full Album
350x522px with tracklist

```markdown
{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=3553509833" style="album" >}}
```

## Parameters

- `url` (required): The Bandcamp embed URL from Share → Embed
- `style` (optional): `"standard"` (default), `"slim"`, `"artwork"`, or `"album"`
- `float` (optional): `"left"` (default for standard/slim), `"right"`, or `"none"` (centered, no text wrapping)
- `bgcol` (optional): Background color in hex (default: `"ffffff"`)
- `linkcol` (optional): Link color in hex (default: `"333333"`)

## Example Playlist

```markdown
+++
title = "December 2025"
date = 2025-12-01
+++

## Track 1

This is a great album that I've been listening to. The production is excellent and the songwriting is top-notch.

{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=3553509833" >}}

The text automatically wraps around the player on the right side, creating a nice newspaper-style layout. You can write as much text as you want here and it will flow naturally around the embed.

## Track 2 (Slim, Centered)

If you want the embed centered without text wrapping, use `float="none"`:

{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/track=1234567890" style="slim" float="none" >}}

## Track 3 (Full Album, Float Right)

You can also float the embed to the right side:

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor.

{{< bandcamp url="https://bandcamp.com/EmbeddedPlayer/album=9876543210" style="album" float="right" >}}

Text will wrap around on the left side when using `float="right"`.
```

## Notes

- **Default behavior**: Standard and slim embeds float left by default (220px width), allowing text to wrap around them like in a newspaper
- **Disable wrapping**: Use `float="none"` to center the embed without text wrapping
- **Float options**: Use `float="left"` or `float="right"` to control which side text wraps on
- The shortcode uses theme-friendly colors by default (white background, dark gray links)
- The embed URL format is: `https://bandcamp.com/EmbeddedPlayer/album=XXXXX` or `track=XXXXX`
- All embeds are responsive and work on mobile devices (floating is disabled on screens < 640px)

