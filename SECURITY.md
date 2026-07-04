# Security Policy

## Reporting Vulnerabilities

If you discover a security vulnerability in ChartForge, please report it responsibly:

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Email the maintainer or use GitHub's private vulnerability reporting
3. Include a description of the vulnerability and steps to reproduce
4. Allow reasonable time for a fix before public disclosure

## Security Considerations

### Input Validation
- ChartForge validates all input data before rendering
- CSV and JSON parsers handle malformed input gracefully
- No shell execution or system commands are performed on input data
- Path traversal is prevented for output file writing

### File Operations
- Output files are written only to the specified path (no directory traversal)
- The `-o` flag only accepts file paths, not commands or URLs
- Generated SVG files are static — no JavaScript execution

### Dependencies
- Minimal dependency footprint (only `github.com/spf13/cobra`)
- Dependencies are pinned to specific versions in `go.sum`
- Regular dependency updates via automated CI

### Data Privacy
- ChartForge processes data locally — no data is sent to external services
- No network requests are made during normal operation
- All processing happens in-memory

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | ✅        |

## Best Practices

When using ChartForge:
- Validate data sources before processing
- Be cautious with untrusted CSV/JSON files (though parsers are designed to be safe)
- Use `go install` from the official repository for secure installation
