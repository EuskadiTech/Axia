# Release Process Guide

This document explains how to use the new GitHub release workflow for the Tallarin project.

## Quick Start

1. **Create a GitHub Release:**
   - Go to your repository's Releases page
   - Click "Create a new release"
   - Create or choose a tag (e.g., `v1.0.0`, `v2.1.3`)
   - Add release notes describing the changes
   - Click "Publish release"

2. **Automatic Build Process:**
   - The workflow will automatically trigger
   - It will build binaries for Linux (amd64, arm64) and Windows (amd64)
   - It will create and push a Docker image
   - All artifacts will be uploaded to your release

## What Gets Built

### Binary Releases
- **Linux AMD64**: `axia4_linux_amd64.tar.gz` - Standard x64 Linux systems
- **Linux ARM64**: `axia4_linux_arm64.tar.gz` - ARM64 systems (Raspberry Pi 4, Apple Silicon, etc.)
- **Windows AMD64**: `axia4_windows_amd64.zip` - Windows 64-bit systems

Each archive contains:
- The compiled binary (`axia4` or `axia4.exe`)
- Configuration template (`config_template.json`)
  - **Note for Windows**: Uses portable configuration by default
- LICENSE file
- **Windows only**: PostgreSQL 16 dependencies in `pgsql16/` folder

### Docker Image
- **Registry**: `ghcr.io/your-org/tallarin`
- **Architectures**: linux/amd64, linux/arm64
- **Tags**: `latest` (for releases from main branch), version tag (e.g., `v1.0.0`)

## Docker Usage

### Pull and Run
```bash
# Pull the latest image
docker pull ghcr.io/your-org/tallarin:latest

# Run with default settings
docker run -p 80:80 -p 443:443 ghcr.io/your-org/tallarin:latest
```

### Using Docker Compose (Development)
```bash
# Start the full development environment
docker-compose up -d

# For development with overrides
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

## Development Workflow

### Local Testing
```bash
# Test cross-compilation
GOOS=linux GOARCH=amd64 go build -o axia4_linux_amd64
GOOS=windows GOARCH=amd64 go build -o axia4_windows_amd64.exe

# Test packaging with the distribution script
python3 package_release.py axia4_linux_amd64 linux --output axia4_linux_amd64
python3 package_release.py axia4_windows_amd64.exe windows --output axia4_windows_amd64

# Test Docker build
docker build -t tallarin:test .
```

### Manual Packaging
If you need to manually create distribution packages, use the `package_release.py` script:

```bash
# Package for Linux (creates tar.gz)
python3 package_release.py <binary_path> linux [--output <name>]

# Package for Windows (creates zip with PostgreSQL dependencies)
python3 package_release.py <binary_path> windows [--output <name>]
```

The script will:
- **Linux**: Create a tar.gz with the binary (renamed to `axia4`), `config_template.json`, and `LICENSE`
- **Windows**: Create a zip with the binary (renamed to `axia4.exe`), `config_template.json` (using portable config), `LICENSE`, and PostgreSQL 16 dependencies in `pgsql16/` folder

### Version Management
The workflow automatically injects the release tag into the binary:
- Release tag `v1.2.3` becomes version `v1.2.3` in the application
- This is done via `-ldflags "-X main.appVersion=${VERSION}"`

## Troubleshooting

### Build Failures
- Check the Actions tab in your GitHub repository
- Verify the Go version matches what's in go.mod
- Ensure all dependencies are properly vendored

### Docker Issues
- Verify Dockerfile builds locally first
- Check that all COPY paths exist
- Ensure container registry permissions are correct

### Release Issues
- Make sure you "publish" the release, not just create a draft
- Verify the workflow has necessary permissions
- Check that the tag follows semantic versioning

## File Structure

```
.github/
├── workflows/
│   ├── release.yml           # Main release workflow
│   └── README.md            # Workflow documentation
├── package_release.py        # Distribution packaging script
├── Dockerfile               # Production Docker image
├── docker-compose.yml       # Development environment
├── docker-compose.dev.yml   # Development overrides
└── .dockerignore           # Docker build optimization
```

## Security Notes

- Docker images run as non-root user `axia4` (UID 1000)
- Static binaries with CGO disabled for security and portability
- Builds use official Go and Alpine base images
- Container registry uses GitHub's GHCR with repository permissions

## Performance Optimizations

### Docker Performance
Docker instances automatically skip fsync operations during file copy/move operations for better performance:
- Container storage drivers (overlay, aufs) can make fsync extremely slow
- The application auto-detects container environments and disables fsync accordingly
- This significantly improves performance during app installs and builder changes
- Can be manually controlled via `REI3_SKIP_FSYNC` environment variable (set to `true` or `false`)
- Data integrity is maintained as the OS will eventually flush data to disk