# Database Template Directory

This directory is intended to contain a pre-initialized PostgreSQL data directory
for the Windows embedded database deployment.

## Purpose

When the Tallarin application starts with an embedded PostgreSQL database on Windows,
it checks if a database directory exists. If not, it attempts to copy from this 
database_template directory to create the initial database.

## Contents

For a production deployment, this directory should contain:
1. A complete PostgreSQL data directory initialized with `initdb`
2. The application schema pre-loaded
3. Proper configuration files (postgresql.conf, pg_hba.conf)

## Fallback Behavior

If this template directory is empty or missing, the embedded database code will
automatically initialize a fresh database using `initdb` at runtime.

## Note

Creating a complete database template is complex and platform-specific.
The runtime initialization approach is more flexible and maintainable.