+++
# Minimal frontmatter - all metadata comes from consumed.toml
# 
# IMPORTANT: The filename (slug) must match the link in consumed.toml
# Example: if consumed.toml has link = "/reviews/lurker", this file should be lurker.md
#
# Metadata (title, year, director, rating, tmdb) is automatically pulled from consumed.toml
# You only need to write the review content here.
date = {{ .Date }}
draft = false
+++

Your review goes here...

## Screenshots (Optional)

{{</* spoiler-gallery file="screenshots-{{ .File.ContentBaseName }}" columns="3" label="Show Screenshots (Spoilers)" */>}}
