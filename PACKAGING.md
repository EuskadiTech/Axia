# Distribution Packaging Script

This script automates the packaging of Axia/REI3 binaries for distribution.

## Features

- Packages binaries for Linux (tar.gz) and Windows (zip)
- Automatically includes required files (binary, config, license)
- For Windows: Downloads and includes PostgreSQL 16 dependencies
- For Windows: Uses portable configuration by default

## Requirements

- Python 3.6 or higher
- Internet connection (for downloading Windows PostgreSQL dependencies)

## Usage

### Basic Usage

```bash
# Package a Linux binary
python3 package_release.py /path/to/axia4_linux_amd64 linux

# Package a Windows binary
python3 package_release.py /path/to/axia4_windows_amd64.exe windows
```

### Custom Output Name

```bash
# Specify custom output name (without extension)
python3 package_release.py /path/to/axia4_linux_amd64 linux --output my_custom_name
python3 package_release.py /path/to/axia4_windows_amd64.exe windows --output my_custom_name
```

## Output

### Linux Package (`*.tar.gz`)

Contains:
- `axia4` - The executable binary
- `config_template.json` - Configuration template
- `LICENSE` - License file

### Windows Package (`*.zip`)

Contains:
- `axia4.exe` - The executable binary
- `config_template.json` - Configuration template (with portable settings)
- `LICENSE` - License file
- `pgsql16/` - PostgreSQL 16 binaries and dependencies

## Example Workflow

```bash
# 1. Build binaries for different platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=1.0.0" -o axia4_linux_amd64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=1.0.0" -o axia4_linux_arm64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=1.0.0" -o axia4_windows_amd64.exe

# 2. Package for distribution
python3 package_release.py axia4_linux_amd64 linux --output axia4_linux_amd64
python3 package_release.py axia4_linux_arm64 linux --output axia4_linux_arm64
python3 package_release.py axia4_windows_amd64.exe windows --output axia4_windows_amd64

# 3. Results
# - axia4_linux_amd64.tar.gz
# - axia4_linux_arm64.tar.gz
# - axia4_windows_amd64.zip
```

## Notes

- The script must be run from the repository root (where `config_template.json`, `config_portable.json`, and `LICENSE` are located)
- Windows packaging requires internet access to download PostgreSQL dependencies from GitHub
- The PostgreSQL package is downloaded from: https://github.com/Axia4/windows-postgres-packages/releases/latest/download/pgsql16.zip
