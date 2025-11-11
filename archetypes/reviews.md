+++
title = "{{ replace .File.ContentBaseName "-" " " | title }}"
date = {{ .Date }}
draft = true
description = ""

[params]
director = ""
year = 
rating = 5

# Uncomment for spoiler screenshots
# spoiler_screenshots = "movie-name-screenshots"
+++

Your review goes here...

## Screenshots (Optional)

{{</* spoiler-gallery file="movie-name-screenshots.toml" columns="3" label="Show Screenshots (Spoilers)" */>}}

