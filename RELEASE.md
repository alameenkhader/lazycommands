# Release Guide

LazyCommands follows [Semantic Versioning](https://semver.org/): `MAJOR.MINOR.PATCH`

- **PATCH**: Bug fixes (0.1.0 → 0.1.1)
- **MINOR**: New features (0.1.0 → 0.2.0)
- **MAJOR**: Breaking changes (0.9.0 → 1.0.0)

## Release Steps

1. **Update version** in `internal/version/version.go`
2. **Update CHANGELOG.md** with changes
3. **Commit and tag:**
   ```bash
   git commit -am "Bump version to vX.Y.Z"
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin main vX.Y.Z
   ```
4. **Build binaries:**
   ```bash
   make build-all
   ```
5. **Publish the GitHub Release** (the curl install script pulls from `releases/latest`):
   ```bash
   gh release create vX.Y.Z bin/lazycommands-* --title "vX.Y.Z" --notes-file CHANGELOG.md
   ```
6. **Update the Homebrew formula** in [alameenkhader/homebrew-tap](https://github.com/alameenkhader/homebrew-tap):
   ```bash
   curl -sL https://github.com/alameenkhader/lazycommands/archive/refs/tags/vX.Y.Z.tar.gz | shasum -a 256
   ```
   Then edit `Formula/lazycommands.rb`: bump the version in `url` and replace `sha256` with the new hash. Commit and push — `brew upgrade` picks it up immediately.

## Check Version

```bash
lazycommands --version
```
