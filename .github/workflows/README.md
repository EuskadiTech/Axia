# GitHub Release Workflow

This directory contains the GitHub Actions workflow for automated releases of the Tallarin project.

## Workflow: `release.yml`

### Trigger
The workflow is triggered when a new GitHub release is created (not just tagged).

### What it does

#### 1. Build Release Binaries
- Builds cross-platform binaries for:
  - Linux AMD64 (`r3_linux_amd64.tar.gz`)
  - Linux ARM64 (`r3_linux_arm64.tar.gz`) 
  - Windows AMD64 (`r3_windows_amd64.zip`)
- Each binary is built with the release version injected via `-ldflags`
- Archives include configuration templates and documentation
- Automatically uploads all artifacts to the GitHub release

#### 2. Build and Push Docker Image
- Creates a multi-architecture Docker image (linux/amd64, linux/arm64)
- Pushes to GitHub Container Registry (`ghcr.io`)
- Tags with the release version and `latest` (if on default branch)
- Includes runtime dependencies: PostgreSQL client, ImageMagick, Ghostscript

### Docker Image Features
- Multi-stage build for smaller final image
- Non-root user for security
- Health check endpoint
- Exposes ports 80 and 443
- Runs with `-run` flag by default

### How to use

1. **Create a new release on GitHub:**
   - Go to your repository's Releases page
   - Click "Create a new release"
   - Choose or create a tag (e.g., `v1.0.0`)
   - Add release notes
   - Click "Publish release"

2. **The workflow will automatically:**
   - Build binaries for all target platforms
   - Create a Docker image
   - Upload all artifacts to your release

3. **Users can then:**
   - Download platform-specific binaries from the release page
   - Pull the Docker image: `docker pull ghcr.io/your-org/tallarin:latest`

### Configuration

The workflow uses these build parameters:
- Go version: 1.24.4 (matches go.mod)
- CGO disabled for static binaries
- Build flags: `-s -w` for smaller binaries
- Version injection: `-X main.appVersion=<release_tag>`

### Permissions Required

The workflow requires:
- `contents: read` - to checkout code
- `packages: write` - to push Docker images to GHCR

These are automatically provided by GitHub Actions for release events.