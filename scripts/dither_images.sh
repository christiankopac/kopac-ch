#!/bin/bash
# HARD Lo-fi Image Dithering using ImageMagick
# Applies aggressive 4x4 ordered dithering with high contrast and posterization
# Creates a chunky, high-contrast black & white lo-fi aesthetic
#
# Usage:
#   ./scripts/dither_images.sh [directory] [--recursive] [--overwrite]
#
# Examples:
#   ./scripts/dither_images.sh assets
#   ./scripts/dither_images.sh content/posts --recursive
#   ./scripts/dither_images.sh assets --overwrite

set -e -o pipefail

# Check if ImageMagick is installed and determine which command to use
if command -v magick &> /dev/null; then
    # ImageMagick v7+ uses 'magick' command
    CONVERT_CMD="magick"
elif command -v convert &> /dev/null; then
    # ImageMagick v6 uses 'convert' command
    CONVERT_CMD="convert"
else
    echo "❌ Error: ImageMagick is not installed"
    echo "Install with:"
    echo "  - Ubuntu/Debian: sudo apt install imagemagick"
    echo "  - macOS: brew install imagemagick"
    echo "  - Windows: https://imagemagick.org/script/download.php"
    exit 1
fi

# Default values
DIRECTORY="."
RECURSIVE=false
OVERWRITE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --recursive|-r)
            RECURSIVE=true
            shift
            ;;
        --overwrite|-o)
            OVERWRITE=true
            shift
            ;;
        -*)
            echo "Unknown option: $1"
            echo "Usage: $0 [directory] [--recursive] [--overwrite]"
            exit 1
            ;;
        *)
            DIRECTORY="$1"
            shift
            ;;
    esac
done

# Validate directory
if [ ! -d "$DIRECTORY" ]; then
    echo "❌ Directory not found: $DIRECTORY"
    exit 1
fi

echo ""
echo "🔍 Searching for images in: $DIRECTORY"
[ "$RECURSIVE" = true ] && echo "   (including subdirectories)"
echo ""

# Find images
if [ "$RECURSIVE" = true ]; then
    FIND_DEPTH=""
else
    FIND_DEPTH="-maxdepth 1"
fi

# Process each image
count=0
processed=0

while IFS= read -r -d '' file; do
    # Skip if already dithered
    if [[ "$file" == *"_dithered"* ]]; then
        continue
    fi
    
    ((count++))
    
    # Get filename info
    dir=$(dirname "$file")
    filename=$(basename "$file")
    extension="${filename##*.}"
    name="${filename%.*}"
    output="$dir/${name}_dithered.$extension"
    
    # Skip if output exists and not overwriting
    if [ -f "$output" ] && [ "$OVERWRITE" = false ]; then
        echo "⏭️  Skipping $filename (dithered version exists)"
        continue
    fi
    
    echo "🎨 Processing $filename..."
    
    # Apply HARD lo-fi dithering with ImageMagick
    # -colorspace Gray: Convert to grayscale
    # -contrast-stretch: Increase contrast dramatically
    # -posterize 2: Reduce to 2 levels before dithering (harder edges)
    # -ordered-dither o4x4: Use 4x4 ordered dithering (more visible/chunky pattern)
    # -normalize: Maximize contrast
    $CONVERT_CMD "$file" \
        -colorspace Gray \
        -contrast-stretch 5%x5% \
        -normalize \
        -posterize 2 \
        -ordered-dither o4x4 \
        "$output"
    
    echo "✅ Created ${name}_dithered.$extension"
    ((processed++))
done < <(find "$DIRECTORY" $FIND_DEPTH -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.webp" \) -print0) || true

echo ""
if [ $count -eq 0 ]; then
    echo "❌ No images found"
else
    echo "✨ Done! Processed $processed of $count image(s)"
    echo "   Dithering: HARD lo-fi 4x4 ordered (high contrast, posterized)"
fi

