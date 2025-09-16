# Agent Instructions

This file contains instructions for AI agents working on this Terraform provider.

## Development Commands

### Build and Test
```bash
# Build the provider
make build

# Run tests
make test

# Run acceptance tests
make testacc

# Install locally for testing
make install

# Lint code
make lint

# Format code
make fmt

# Generate documentation
make generate
```

### Local Development Setup
```bash
# Copy dev overrides for local testing
cp examples/.terraformrc.bak ~/.terraformrc

# Build and test changes in examples directory
cd examples/data-sources/dirt_metadata
terraform init
terraform plan
terraform apply
```

## Publishing a New Release

### Prerequisites
- Write access to the GitHub repository
- GPG key configured for signing releases
- GoReleaser installed (if publishing manually)

### Release Process

1. **Update the changelog** in `CHANGELOG.md`:
   ```markdown
   ## 0.X.0 (Month DD, YYYY)

   FEATURES:
   * **new feature**: Description

   ENHANCEMENTS:  
   * **component**: Improvement description

   BUG FIXES:
   * **component**: Bug fix description
   ```

2. **Commit and push changes**:
   ```bash
   git add CHANGELOG.md
   git commit -m "docs: update changelog for vX.X.X release"
   git push
   ```

3. **Create and push a git tag**:
   ```bash
   git tag vX.X.X
   git push origin vX.X.X
   ```

4. **Automated release process**:
   - GitHub Actions will automatically build and release using GoReleaser
   - The release will be published to GitHub Releases
   - Binaries will be built for multiple platforms (Linux, macOS, Windows)

5. **Manual release** (if needed):
   ```bash
   # Set GPG fingerprint environment variable
   export GPG_FINGERPRINT="your-gpg-key-fingerprint"
   
   # Run GoReleaser
   goreleaser release --clean
   ```

### Version Numbering
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Breaking changes require a major version bump
- New features require a minor version bump  
- Bug fixes require a patch version bump

### Release Verification
After release, verify:
1. GitHub Release is created with proper assets
2. All platform binaries are present
3. Checksums and signatures are included
4. Release notes match the changelog

## Project Structure

- `internal/client/` - API client for DirtCloud server
- `internal/provider/` - Terraform provider implementation
- `examples/` - Example Terraform configurations
- `docs/` - Provider documentation
- `.goreleaser.yml` - Release configuration
- `GNUmakefile` - Build and development commands

## Testing

### Local Server
The provider requires a running DirtCloud server at `http://localhost:8080/v1`. 
The server code is located at `~/dirtcloud-server`.

### Database Schema
If you encounter database errors, ensure the metadata table has the correct schema:
```sql
CREATE TABLE metadata (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Recent Changes

### v0.2.0
- Updated metadata data source to use `path` instead of `id` as the required parameter
- Added `GetMetadataByPath` method to client for path-based metadata retrieval
