# Tasks

## TODO

- [x] Add tasks here

- [x] Rename module to github.com/go-go-golems/sanitize and restructure packages
- [x] Move sanitize/ to pkg/yaml/ and rename package to yamlsanitize
- [x] Precompile regexes (fixMissingSpaceAfterColon, fixExtraColonInValue)
- [x] Add functional options for Sanitize (rules to enable, indent width, max iterations)
- [x] Embed static assets with go:embed
- [x] Add comprehensive unit tests for lint rules and fix rules
- [x] Add CLI mode (stdin->stdout) in cmd/sanitize/
- [x] Move HTTP server to cmd/sanitize-server/
- [x] Write README.md with usage examples and supported rules
- [ ] Clean up error handling (return errors instead of swallowing)
