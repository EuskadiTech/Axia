#!/usr/bin/env python3
"""
Distribution Package Builder for Axia
Packages binary releases for different operating systems with appropriate dependencies.
"""

import argparse
import json
import os
import shutil
import sys
import tarfile
import tempfile
import urllib.request
import zipfile
from pathlib import Path

# Windows PostgreSQL dependency URL
WINDOWS_POSTGRES_URL = "https://github.com/Axia4/windows-postgres-packages/releases/latest/download/pgsql16.zip"


def download_file(url, dest_path):
    """Download a file from URL to destination path."""
    print(f"Downloading {url}...")
    try:
        with urllib.request.urlopen(url) as response, open(dest_path, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
        print(f"Downloaded to {dest_path}")
        return True
    except Exception as e:
        print(f"Error downloading {url}: {e}", file=sys.stderr)
        return False


def extract_zip(zip_path, extract_to):
    """Extract a zip file to a destination directory."""
    print(f"Extracting {zip_path} to {extract_to}...")
    try:
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_to)
        print(f"Extracted successfully")
        return True
    except Exception as e:
        print(f"Error extracting {zip_path}: {e}", file=sys.stderr)
        return False


def create_config_template_for_windows(source_portable_path):
    """Create config_template.json with contents from config_portable.json for Windows."""
    print("Creating config_template.json from config_portable.json for Windows...")
    try:
        with open(source_portable_path, 'r') as f:
            portable_config = json.load(f)
        return portable_config
    except Exception as e:
        print(f"Error reading config_portable.json: {e}", file=sys.stderr)
        return None


def package_linux(binary_path, output_name):
    """Package Linux distribution as tar.gz."""
    print(f"Packaging Linux distribution: {output_name}")
    
    # Create temporary directory for packaging
    with tempfile.TemporaryDirectory() as temp_dir:
        package_dir = os.path.join(temp_dir, "package")
        os.makedirs(package_dir)
        
        # Copy binary
        binary_dest = os.path.join(package_dir, "axia4")
        shutil.copy2(binary_path, binary_dest)
        os.chmod(binary_dest, 0o755)  # Make executable
        print(f"Added binary: {binary_dest}")
        
        # Copy config template
        config_src = "config_template.json"
        if os.path.exists(config_src):
            shutil.copy2(config_src, os.path.join(package_dir, "config_template.json"))
            print(f"Added config_template.json")
        else:
            print(f"Warning: {config_src} not found", file=sys.stderr)
        
        # Copy LICENSE
        license_src = "LICENSE"
        if os.path.exists(license_src):
            shutil.copy2(license_src, os.path.join(package_dir, "LICENSE"))
            print(f"Added LICENSE")
        else:
            print(f"Warning: {license_src} not found", file=sys.stderr)
        
        # Create tar.gz archive
        output_path = f"{output_name}.tar.gz"
        print(f"Creating archive: {output_path}")
        with tarfile.open(output_path, "w:gz") as tar:
            for item in os.listdir(package_dir):
                item_path = os.path.join(package_dir, item)
                tar.add(item_path, arcname=item)
        
        print(f"Successfully created {output_path}")
        return True


def package_windows(binary_path, output_name):
    """Package Windows distribution as zip with PostgreSQL dependencies."""
    print(f"Packaging Windows distribution: {output_name}")
    
    # Create temporary directory for packaging
    with tempfile.TemporaryDirectory() as temp_dir:
        package_dir = os.path.join(temp_dir, "package")
        os.makedirs(package_dir)
        
        # Copy binary
        binary_dest = os.path.join(package_dir, "axia4.exe")
        shutil.copy2(binary_path, binary_dest)
        print(f"Added binary: {binary_dest}")
        
        # Create config_template.json from config_portable.json for Windows
        config_portable_src = "config_portable.json"
        if os.path.exists(config_portable_src):
            portable_config = create_config_template_for_windows(config_portable_src)
            if portable_config:
                config_dest = os.path.join(package_dir, "config_template.json")
                with open(config_dest, 'w') as f:
                    json.dump(portable_config, f, indent="\t")
                print(f"Added config_template.json (from config_portable.json)")
            else:
                print("Error creating config_template.json", file=sys.stderr)
                return False
        else:
            print(f"Warning: {config_portable_src} not found", file=sys.stderr)
            return False
        
        # Copy LICENSE
        license_src = "LICENSE"
        if os.path.exists(license_src):
            shutil.copy2(license_src, os.path.join(package_dir, "LICENSE"))
            print(f"Added LICENSE")
        else:
            print(f"Warning: {license_src} not found", file=sys.stderr)
        
        # Download and extract PostgreSQL dependencies
        print("Downloading Windows PostgreSQL dependencies...")
        postgres_zip_path = os.path.join(temp_dir, "pgsql16.zip")
        if not download_file(WINDOWS_POSTGRES_URL, postgres_zip_path):
            print("Error downloading PostgreSQL dependencies", file=sys.stderr)
            return False
        
        # Extract PostgreSQL to package directory
        postgres_extract_dir = os.path.join(temp_dir, "postgres_extract")
        os.makedirs(postgres_extract_dir)
        if not extract_zip(postgres_zip_path, postgres_extract_dir):
            print("Error extracting PostgreSQL dependencies", file=sys.stderr)
            return False
        
        # Copy pgsql16 folder to package directory
        pgsql_src = os.path.join(postgres_extract_dir, "pgsql16")
        if os.path.exists(pgsql_src):
            pgsql_dest = os.path.join(package_dir, "pgsql16")
            shutil.copytree(pgsql_src, pgsql_dest)
            print(f"Added PostgreSQL dependencies: pgsql16/")
        else:
            print(f"Warning: pgsql16 folder not found in extracted archive", file=sys.stderr)
            return False
        
        # Create zip archive
        output_path = f"{output_name}.zip"
        print(f"Creating archive: {output_path}")
        with zipfile.ZipFile(output_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
            for root, dirs, files in os.walk(package_dir):
                for file in files:
                    file_path = os.path.join(root, file)
                    arcname = os.path.relpath(file_path, package_dir)
                    zipf.write(file_path, arcname)
        
        print(f"Successfully created {output_path}")
        return True


def main():
    parser = argparse.ArgumentParser(
        description="Package Axia releases for distribution"
    )
    parser.add_argument(
        "binary_path",
        help="Path to the binary to package"
    )
    parser.add_argument(
        "os",
        choices=["linux", "windows"],
        help="Target operating system"
    )
    parser.add_argument(
        "--output",
        default=None,
        help="Output package name (without extension)"
    )
    
    args = parser.parse_args()
    
    # Verify binary exists
    if not os.path.exists(args.binary_path):
        print(f"Error: Binary not found: {args.binary_path}", file=sys.stderr)
        return 1
    
    # Determine output name
    if args.output:
        output_name = args.output
    else:
        # Use binary name without extension
        binary_name = os.path.basename(args.binary_path)
        output_name = os.path.splitext(binary_name)[0]
    
    # Package based on OS
    success = False
    if args.os == "linux":
        success = package_linux(args.binary_path, output_name)
    elif args.os == "windows":
        success = package_windows(args.binary_path, output_name)
    
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
