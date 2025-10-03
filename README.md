![r3_logo_git](www/images/logo.png)
<h1 align="center">Axia</h1>
<p align="center"><strong>Free and open low code platform</strong><br />Build and host powerful applications with full control</p>

<p align="center">
	<a href="https://github.com/EuskadiTech/tallarin/releases"><img src="https://img.shields.io/github/v/release/EuskadiTech/tallarin" alt="Release" /></a>
	<a href="https://github.com/EuskadiTech/tallarin/stargazers"><img src="https://img.shields.io/github/stars/EuskadiTech/tallarin" alt="Stars" /></a>
	<a href="https://github.com/EuskadiTech/tallarin/blob/main/LICENSE"><img src="https://img.shields.io/github/license/EuskadiTech/tallarin" alt="License" /></a>
</p>
<p align="center">
	<a href="https://tech.eus/t4/downloads.php">Downloads</a> · 
	<a href="https://tech.eus/t4/docs.php">Documentation</a> · 
	<a href="https://tech.eus/t4/apps.php">Applications</a>
</p>

**Axia** is a self-hosted low code platform for building multi-user applications. Break free from cloud-only SaaS and walled gardens - create, use, share, or sell your applications with complete ownership.

![DEMO - Organizations](https://github.com/user-attachments/assets/5506d0c1-4bf3-4011-bc3a-2650cb5ec0b9)
![DEMO - Gantt Chart](https://github.com/user-attachments/assets/1e413540-f9e8-4c2f-bd91-f46f51137d8b)

## Features

**Core Capabilities**
- **Rapid Development**: Replace spreadsheets with proper multi-user applications quickly
- **Business Logic**: Calculations, date operations, business rules, and data summarization
- **Rich Visualizations**: Gantt charts, diagrams, and information-dense lists
- **Workflow Automation**: State-based forms, PDF exports, and notifications
- **Access Control**: Roles and policies for global or record-specific permissions

**Technical Features**
- **End-to-End Encryption**: Built-in E2EE with integrated key management
- **Integration Ready**: REST APIs, CSV import/export, ICS calendar feeds
- **Mobile-First**: Responsive design with PWA support for app-like experience
- **Full-text Search**: Language-specific search with phrase matching
- **Diverse Inputs**: Date ranges, signatures, QR/barcode scanning, and more
- **High Performance**: Multi-core optimization with bi-directional data channels

**Enterprise Ready**
- **Security**: Password policies, brute-force protection, MFA support
- **Flexibility**: Direct database access, human-readable data storage
- **Self-Hosted**: Deploy locally or in the cloud with full data control
- **Scalable**: LDAP integration, custom branding, clustering support

![DEMO - Password Manager](https://github.com/user-attachments/assets/e9161bf2-027e-409d-a9eb-ed97dfe76f7e)
![DEMO - IT Assets](https://github.com/user-attachments/assets/c5273f72-24cb-40cc-a947-c6a42b78f7bb)
![DEMO - Time Tracker](https://github.com/user-attachments/assets/e6b6e0e9-558a-4bad-ad52-45700e7d229e)

## Quickstart

### Installation

**Linux** ([x64](https://tech.eus/t4/downloads.php) / [arm64](https://tech.eus/t4/downloads.php))
```bash
# 1. Extract and make executable
chmod u+x r3

# 2. Configure database
cp config_template.json config.json
# Edit config.json with PostgreSQL database details (UTF8, full permissions)

# 3. Install optional dependencies
sudo apt install imagemagick ghostscript postgresql-client

# 4. Register and start
sudo ./r3 -install
sudo systemctl start rei3
```

**Windows**
- Use the standalone/portable build for Windows Server
- Optionally install [Ghostscript](https://www.ghostscript.com/) for PDF thumbnails

**Docker**
- Docker Compose files available: [x64](https://rei3.de/docker_x64) / [arm64](https://rei3.de/docker_arm64)

### First Login
Access Axia at `https://localhost:6129` (default credentials: `admin` / `admin`)

**Important**: If running behind a proxy, disable client timeouts for websockets. See [proxy documentation](https://rei3.de/en/docs/admin#proxies).

Full documentation: [tech.eus/t4/docs.php](https://tech.eus/t4/docs.php)

## Resources

- **Documentation**: [tech.eus/t4/docs.php](https://tech.eus/t4/docs.php) - [Admin Guide](https://rei3.de/en/docs/admin) - [Builder Guide](https://rei3.de/en/docs/builder)
- **Community Forum**: [community.rei3.de](https://community.rei3.de)
- **Video Tutorials**: [YouTube Channel](https://www.youtube.com/channel/UCKb1YPyUV-O4GxcCdHc4Csw)

## Contributing

Contributions welcome! Fork the repository and submit pull requests.

**Versioning**: Axia uses `MAJOR.MINOR.PATCH.BUILD` format (e.g., `3.2.0.4246`). The `main` branch contains the current release. New minor versions use separate branches that merge to `main` upon release.

## Community Projects

- [R3 Toolshop](https://github.com/Umb-Astardo/R3-Toolshop) - Operational tools for data import, bulk user creation, and more
- [Axia-Tickets-MCP-Server](https://github.com/lgndluke/Axia-Tickets-MCP-Server) - LLM integration for Axia Tickets
- [Google Material Icons](https://github.com/fmvalsera/r3_material_icons_app) - Material icon pack application

## Building from Source

```bash
# 1. Install Go (latest version)
# 2. Build
go build -ldflags "-X main.appVersion={VERSION}"

# Example:
go build -ldflags "-X main.appVersion=2.5.1.2980"

# Cross-compile (optional)
GOOS=windows go build -ldflags "-X main.appVersion=2.5.1.2980"
GOOS=linux go build -ldflags "-X main.appVersion=2.5.1.2980"
```

**Development Tips**:
- Static files (HTML, JS, CSS) are embedded in the binary
- Use `-wwwpath` argument to load external `www` directory for live changes
- Clear browser cache after updates

## Built With

- [Golang](https://golang.org/) - Backend and web services
- [PostgreSQL](https://www.postgresql.org/) - Database management
- [Vue.js](https://vuejs.org/) - Frontend framework

## License

Copyright (c) 2025 EuskadiTech (TM)

Released under the [MIT License](https://opensource.org/license/mit).

Axia is based on r3_team/r3. For trademark issues, open an issue or email: soporte (at) tech (dot) eus

---

[![Stargazers](https://reporoster.com/stars/dark/EuskadiTech/tallarin)](https://github.com/EuskadiTech/tallarin/stargazers)
