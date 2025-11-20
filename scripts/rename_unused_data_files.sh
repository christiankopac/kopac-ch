#!/bin/bash
# Rename unused data files with reference_ prefix

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DATA_DIR="$BASE_DIR/data"

# Files that are still used (don't rename these)
USED_FILES=(
    "data/movies/screenshots-*.toml"  # Still used by spoiler-gallery (being migrated)
)

# Files that are only used by scripts (keep for now, scripts will be updated)
# These will be renamed after scripts are updated to work with pages
SCRIPT_ONLY_FILES=(
    # "data/movies/movies.toml"  # Keep for now - scripts still use it
    # "data/music/music.toml"     # Keep for now - scripts still use it
    # "data/books/books.toml"    # Keep for now - scripts still use it
)

# Files to check if they're used in content
POTENTIALLY_UNUSED=(
    "data/bookmarks.toml"
    "data/bookmarks_simple.toml"
    "data/experiences.toml"
    "data/experiences_simple.toml"
    "data/podcasts.toml"
    "data/projects.toml"
    "data/projects_simple.toml"
    "data/publications.toml"
    "data/skills.toml"
    "data/uses.toml"
)

echo "Checking which data files are used..."
echo ""

# Check if files are referenced in content
check_file_usage() {
    local file="$1"
    local basename=$(basename "$file" .toml)
    local dirname=$(dirname "$file")
    local search_path=$(echo "$file" | sed "s|^data/||" | sed "s|\.toml$||")
    
    # Check in content files
    if grep -r "$basename\|$search_path" content/ themes/ --include="*.md" --include="*.html" 2>/dev/null | grep -v "reference_" | grep -q .; then
        return 0  # File is used
    fi
    
    return 1  # File is not used
}

# Rename files that are not used
renamed=0
skipped=0

for file in "${POTENTIALLY_UNUSED[@]}"; do
    if [ ! -f "$BASE_DIR/$file" ]; then
        continue
    fi
    
    if check_file_usage "$file"; then
        echo "✓ Keeping $file (used in content)"
        ((skipped++))
    else
        dirname=$(dirname "$file")
        basename=$(basename "$file")
        new_name="reference_$basename"
        new_path="$BASE_DIR/$dirname/$new_name"
        
        if [ -f "$new_path" ]; then
            echo "⚠ Skipping $file (reference_ version already exists)"
            ((skipped++))
        else
            mv "$BASE_DIR/$file" "$new_path"
            echo "✓ Renamed $file → $new_name"
            ((renamed++))
        fi
    fi
done

# Rename script-only files (movies, music, books)
echo ""
echo "Renaming script-only files (used by scripts, not templates)..."
for file in "${SCRIPT_ONLY_FILES[@]}"; do
    if [ ! -f "$BASE_DIR/$file" ]; then
        continue
    fi
    
    dirname=$(dirname "$file")
    basename=$(basename "$file")
    new_name="reference_$basename"
    new_path="$BASE_DIR/$dirname/$new_name"
    
    if [ -f "$new_path" ]; then
        echo "⚠ Skipping $file (reference_ version already exists)"
        ((skipped++))
    else
        mv "$BASE_DIR/$file" "$new_path"
        echo "✓ Renamed $file → $new_name"
        ((renamed++))
    fi
done

echo ""
echo "Done! Renamed $renamed file(s), skipped $skipped file(s)."
echo ""
echo "Note: Screenshot files (screenshots-*.toml) are being migrated to frontmatter."
echo "      They will be renamed after migration is complete."

