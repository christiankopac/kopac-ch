#!/bin/bash
# Build script for Hugo site with optional image dithering
#
# Usage:
#   ./build.sh                  # Build with dithering
#   ./build.sh --skip-dither    # Build without dithering
#   ./build.sh --force          # Force re-dithering of all images

set -e -o pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SKIP_DITHER=false
FORCE_DITHER=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-dither)
            SKIP_DITHER=true
            shift
            ;;
        --force)
            FORCE_DITHER=true
            shift
            ;;
        -*)
            echo "Unknown option: $1"
            echo "Usage: $0 [--skip-dither] [--force]"
            exit 1
            ;;
        *)
            shift
            ;;
    esac
done

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      Hugo Site Build Process           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Image dithering
if [ "$SKIP_DITHER" = false ]; then
    echo -e "${YELLOW}[1/2]${NC} Processing images..."
    
    DITHER_OPTS="--recursive"
    if [ "$FORCE_DITHER" = true ]; then
        DITHER_OPTS="$DITHER_OPTS --overwrite"
        echo -e "      ${BLUE}↳${NC} Force mode: Re-dithering all images"
    fi
    
    # Dither images in assets directory
    if [ -d "assets" ]; then
        echo -e "      ${BLUE}↳${NC} Dithering images in assets/"
        ./scripts/dither_images.sh assets $DITHER_OPTS
    fi
    
    # Dither images in content directory
    if [ -d "content" ]; then
        echo -e "      ${BLUE}↳${NC} Dithering images in content/"
        ./scripts/dither_images.sh content $DITHER_OPTS
    fi
    
    # Dither images in static directory
    if [ -d "static" ]; then
        echo -e "      ${BLUE}↳${NC} Dithering images in static/"
        ./scripts/dither_images.sh static $DITHER_OPTS
    fi
    
    echo -e "${GREEN}✓${NC} Image processing complete"
    echo ""
else
    echo -e "${YELLOW}[1/2]${NC} Skipping image dithering"
    echo ""
fi

# Step 2: Build site
echo -e "${YELLOW}[2/2]${NC} Building Hugo site..."
echo ""

if hugo; then
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║        Build Successful! 🎉            ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  ${BLUE}→${NC} Site generated in: ${GREEN}public/${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║          Build Failed! ❌              ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}"
    echo ""
    exit 1
fi

