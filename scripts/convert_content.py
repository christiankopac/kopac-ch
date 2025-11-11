#!/usr/bin/env python3
"""
Convert Zola content to Hugo format.
- Removes [extra] nesting in front matter
- Converts shortcode syntax from Zola to Hugo
"""

import re
import sys
from pathlib import Path

def convert_front_matter(content):
    """Convert TOML front matter from Zola to Hugo format."""
    lines = content.split('\n')
    if not (lines[0].strip() == '+++'):
        return content
    
    # Find the closing +++
    end_idx = None
    for i in range(1, len(lines)):
        if lines[i].strip() == '+++':
            end_idx = i
            break
    
    if end_idx is None:
        return content
    
    # Process front matter
    front_matter = lines[1:end_idx]
    body = '\n'.join(lines[end_idx+1:])
    
    new_front_matter = []
    in_extra = False
    extra_indent = 0
    
    for line in front_matter:
        stripped = line.strip()
        
        # Skip [extra] line
        if stripped == '[extra]':
            in_extra = True
            # Calculate indent level for extra content
            extra_indent = len(line) - len(line.lstrip())
            continue
        
        # Check if we're leaving extra section (new section starts)
        if in_extra and line and not line[0].isspace() and '[' in line:
            in_extra = False
        
        # If we're in extra, remove the indent
        if in_extra and line.startswith(' ' * (extra_indent + 2)):
            # Remove extra indent (typically 2 spaces)
            new_front_matter.append(line[2:])
        elif not in_extra:
            new_front_matter.append(line)
        else:
            new_front_matter.append(line)
    
    return '+++\n' + '\n'.join(new_front_matter) + '\n+++\n' + body

def convert_shortcodes(content):
    """Convert Zola shortcode syntax to Hugo syntax."""
    
    # Pattern: {% shortcode(param="value", param2=value) %} ... {% end %}
    # Convert to: {{< shortcode param="value" param2="value" >}} ... {{< /shortcode >}}
    
    # First, handle opening tags with parameters
    def replace_opening(match):
        shortcode_name = match.group(1)
        params = match.group(2)
        
        # Convert parameters: remove parentheses, keep quotes as is
        if params:
            params = params.strip('()')
            # Convert boolean true/false without quotes
            params = re.sub(r'=\s*true\b', '="true"', params)
            params = re.sub(r'=\s*false\b', '="false"', params)
            # Convert numbers without quotes  
            params = re.sub(r'=\s*(\d+)\b', r'="\1"', params)
            return f'{{{{< {shortcode_name} {params} >}}}}'
        else:
            return f'{{{{< {shortcode_name} >}}}}'
    
    # Handle self-closing shortcodes: {% shortcode(...) %}
    content = re.sub(r'\{%\s*(\w+)\((.*?)\)\s*%\}', replace_opening, content)
    
    # Handle closing tags: {% end %}
    def replace_closing(match):
        # Look backwards to find the last opened shortcode
        before = content[:match.start()]
        # Find last opening tag
        last_open = None
        for m in re.finditer(r'\{\{<\s*(\w+)', before):
            last_open = m.group(1)
        if last_open:
            return f'{{{{< /{last_open} >}}}}'
        return match.group(0)
    
    # Simpler approach: convert {% end %} to generic closing
    # We'll do this in a second pass after we know all shortcodes
    lines = content.split('\n')
    converted_lines = []
    shortcode_stack = []
    
    for line in lines:
        # Check for opening shortcodes
        open_matches = list(re.finditer(r'\{\{<\s*(\w+)', line))
        for match in open_matches:
            shortcode_stack.append(match.group(1))
        
        # Check for {% end %}
        if re.search(r'\{%\s*end\s*%\}', line):
            if shortcode_stack:
                shortcode_name = shortcode_stack.pop()
                line = re.sub(r'\{%\s*end\s*%\}', f'{{{{< /{shortcode_name} >}}}}', line)
        
        converted_lines.append(line)
    
    return '\n'.join(converted_lines)

def process_file(filepath):
    """Process a single markdown file."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Convert front matter
    content = convert_front_matter(content)
    
    # Convert shortcodes
    content = convert_shortcodes(content)
    
    # Write back
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ Converted: {filepath}")

def main():
    if len(sys.argv) < 2:
        print("Usage: python convert_content.py <directory>")
        sys.exit(1)
    
    directory = Path(sys.argv[1])
    if not directory.exists():
        print(f"Error: Directory {directory} does not exist")
        sys.exit(1)
    
    # Process all .md files
    md_files = list(directory.rglob('*.md'))
    print(f"Found {len(md_files)} markdown files")
    
    for filepath in md_files:
        try:
            process_file(filepath)
        except Exception as e:
            print(f"✗ Error processing {filepath}: {e}")
    
    print(f"\nDone! Processed {len(md_files)} files")

if __name__ == '__main__':
    main()

