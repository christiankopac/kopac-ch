# Footer License Configuration

The apacible theme supports configurable footer license display. You can choose from various Creative Commons licenses, traditional copyright, or no license display.

## Configuration in `hugo.toml`

```toml
[params]
  # Footer License Configuration
  # Options: "cc-by", "cc-by-sa", "cc-by-nc", "cc-by-nc-sa", "cc-by-nd", "cc-by-nc-nd", "copyright", "none"
  # Default: "cc-by-nc-nd"
  footer_license = "cc-by-nc-nd"
  
  # Copyright year and holder (used for "copyright" type)
  copyright_holder = "Your Name"
  copyright_year = "2025"
```

## Available License Options

### Creative Commons Licenses

All CC licenses display the official badge image from Creative Commons.

| Option | Full Name | Description | Image URL |
|--------|-----------|-------------|-----------|
| `cc-by` | Attribution | Allows remixing, even commercially | https://i.creativecommons.org/l/by/4.0/88x31.png |
| `cc-by-sa` | Attribution-ShareAlike | Remix allowed, must share alike | https://i.creativecommons.org/l/by-sa/4.0/88x31.png |
| `cc-by-nc` | Attribution-NonCommercial | Non-commercial use only | https://i.creativecommons.org/l/by-nc/4.0/88x31.png |
| `cc-by-nc-sa` | Attribution-NonCommercial-ShareAlike | Non-commercial, share alike | https://i.creativecommons.org/l/by-nc-sa/4.0/88x31.png |
| `cc-by-nd` | Attribution-NoDerivatives | No derivatives allowed | https://i.creativecommons.org/l/by-nd/4.0/88x31.png |
| `cc-by-nc-nd` | Attribution-NonCommercial-NoDerivatives | Most restrictive CC license | https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png |

### Traditional Copyright

| Option | Display |
|--------|---------|
| `copyright` | Displays: "© 2025 Your Name" (using configured values) |

### No License

| Option | Display |
|--------|---------|
| `none` | No license information displayed |

## Footer Layout

The footer displays in a single row with three sections:

```
┌──────────────────────────────────────────────────────────┐
│ [License]    Built with Hugo and apacible    [RSS] [🌙] │
│   Left              Center                      Right    │
└──────────────────────────────────────────────────────────┘
```

- **Left**: License badge/copyright text
- **Center**: Site credits (Hugo + theme)
- **Right**: RSS feed (on /posts) + Theme toggle

## Examples

### CC BY-NC-ND (Default)

```toml
[params]
  footer_license = "cc-by-nc-nd"
```

Displays the CC BY-NC-ND 4.0 badge linking to the license.

### Traditional Copyright

```toml
[params]
  footer_license = "copyright"
  copyright_holder = "Jane Doe"
  copyright_year = "2025"
```

Displays: "© 2025 Jane Doe"

### CC BY (Most Permissive)

```toml
[params]
  footer_license = "cc-by"
```

Displays the CC BY 4.0 badge allowing maximum freedom for reuse.

### No License Display

```toml
[params]
  footer_license = "none"
```

Only displays the site credits and icons (no license information).

## Notes

- All CC license badges are 88x31px PNG images hosted by Creative Commons
- License badges link to the full license text on creativecommons.org
- The `copyright_year` defaults to the current year if not specified
- The `copyright_holder` defaults to "Your Name" if not specified
- Mobile displays maintain the same single-row layout with compact spacing

