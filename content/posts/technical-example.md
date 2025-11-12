+++
title = "Example: Math and Code Snippets"
date = 2024-11-04
description = "A comprehensive guide with various code examples"


tags = ["apacible", "example", "latex", "code"]
+++

This page demonstrates various technical examples with code blocks and detailed explanations.

## Web Development

### HTML Structure

A basic HTML5 template:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <main>
        <h1>Hello World</h1>
        <p>Welcome to my site.</p>
    </main>
    <script src="script.js"></script>
</body>
</html>
```

### CSS Styling

Modern CSS with custom properties:

```css
:root {
    --primary-color: #3498db;
    --secondary-color: #2ecc71;
    --font-stack: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

@media (prefers-color-scheme: dark) {
    :root {
        --primary-color: #5dade2;
        --secondary-color: #58d68d;
    }
}
```

## Backend Development

### Python API Example

Using FastAPI for modern Python web services:

```python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI()

class Item(BaseModel):
    id: int
    name: str
    description: Optional[str] = None
    price: float

items_db = []

@app.get("/items", response_model=List[Item])
async def get_items():
    return items_db

@app.post("/items", response_model=Item)
async def create_item(item: Item):
    items_db.append(item)
    return item

@app.get("/items/{item_id}")
async def get_item(item_id: int):
    for item in items_db:
        if item.id == item_id:
            return item
    raise HTTPException(status_code=404, detail="Item not found")
```

### Database Queries

SQL examples for common operations:

```sql
-- Create a users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO users (username, email) VALUES
    ('alice', 'alice@example.com'),
    ('bob', 'bob@example.com');

-- Query with joins
SELECT 
    u.username,
    u.email,
    COUNT(p.id) as post_count
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
GROUP BY u.id, u.username, u.email
HAVING COUNT(p.id) > 5
ORDER BY post_count DESC;
```

## Shell Scripting

### Bash Script Example

A useful deployment script:

```bash
#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment...${NC}"

# Run tests
if npm test; then
    echo -e "${GREEN}Tests passed!${NC}"
else
    echo -e "${RED}Tests failed!${NC}"
    exit 1
fi

# Build project
npm run build

# Deploy
rsync -avz --delete dist/ user@server:/var/www/html/

echo -e "${GREEN}Deployment complete!${NC}"
```

## Configuration Files

### TOML Configuration

Example Zola configuration:

```toml
base_url = "https://example.com"
title = "My Site"
description = "A demo site"

compile_sass = true
minify_html = false

[markdown]
highlight_code = true
highlight_theme = "base16-ocean-dark"
render_emoji = false

[extra]
sections = [
  { name = "posts", path = "/posts", is_external = false },
  { name = "about", path = "/about", is_external = false },
]
```

### JSON Configuration

Package configuration:

```json
{
  "name": "my-project",
  "version": "1.0.0",
  "description": "A sample project",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "dev": "nodemon index.js",
    "test": "jest",
    "build": "webpack --mode production"
  },
  "dependencies": {
    "express": "^4.18.0",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "nodemon": "^2.0.0",
    "jest": "^29.0.0"
  }
}
```

## Algorithm Examples

### Sorting Algorithm

Quick sort implementation:

```python
def quicksort(arr):
    """
    Quick sort algorithm implementation.
    Time complexity: O(n log n) average case
    Space complexity: O(log n)
    """
    if len(arr) <= 1:
        return arr
    
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    
    return quicksort(left) + middle + quicksort(right)

# Example usage
numbers = [3, 6, 8, 10, 1, 2, 1]
sorted_numbers = quicksort(numbers)
print(f"Sorted: {sorted_numbers}")
```

### Binary Search

Efficient searching in sorted arrays:

```javascript
function binarySearch(arr, target) {
  let left = 0;
  let right = arr.length - 1;
  
  while (left <= right) {
    const mid = Math.floor((left + right) / 2);
    
    if (arr[mid] === target) {
      return mid; // Found!
    } else if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }
  
  return -1; // Not found
}

// Example
const numbers = [1, 3, 5, 7, 9, 11, 13];
const index = binarySearch(numbers, 7);
console.log(`Found at index: ${index}`); // Output: 3
```

## Conclusion

These examples demonstrate the versatility of code blocks on this site, with proper syntax highlighting and formatting for multiple programming languages.

The site supports many more languages including Go, Ruby, PHP, Java, C++, and more. All code is displayed with appropriate syntax highlighting based on the language identifier in the code fence.

