![r3_logo_git](https://github.com/EuskadiTech/tallarin/assets/91060542/a759e7ec-e1a0-4a4e-a426-509abc764352)
<h1 align="center">Axia</h1>
<p align="center"><strong>Free and open low code</strong><br />Build and host powerful applications with full control and ownership</p>


<p align="center">
	<a href="https://github.com/EuskadiTech/tallarin/releases" target="_blank">
		<img src="https://img.shields.io/github/v/release/EuskadiTech/tallarin" alt="Latest GitHub release" />
	</a>
	<a href="https://tech.eus/t4/downloads.php" target="_blank">
		<img src="https://img.shields.io/badge/linux-x64-yellow" alt="Latest Linux x64" />
	</a>
	<a href="https://tech.eus/t4/downloads.php" target="_blank">
		<img src="https://img.shields.io/badge/linux-arm64-yellow" alt="Latest Linux arm64" />
	</a>
	<a href="https://tech.eus/t4/downloads.php" target="_blank">
		<img src="https://img.shields.io/badge/windows-x64-00a8e8" alt="Latest Windows x64" />
	</a>
	<a href="https://img.shields.io/github/go-mod/go-version/EuskadiTech/tallarin" target="_blank">
		<img src="https://img.shields.io/github/go-mod/go-version/EuskadiTech/tallarin" alt="GitHub go.mod Go version" />
	</a>
	<a href="https://github.com/EuskadiTech/tallarin/stargazers" target="_blank">
		<img src="https://img.shields.io/github/stars/EuskadiTech/tallarin" alt="GitHub repo stars" />
	</a>
	<a href="https://github.com/EuskadiTech/tallarin/commits/main" target="_blank">
		<img src="https://img.shields.io/github/commit-activity/t/EuskadiTech/tallarin" alt="GitHub commit activity" />
	</a>
	<a href="https://github.com/EuskadiTech/tallarin/blob/main/LICENSE" target="_blank">
		<img src="https://img.shields.io/github/license/EuskadiTech/tallarin" alt="License" />
	</a>
</p>
<p align="center">
	<a href="https://tech.eus/t4/downloads.php" target="_blank">Downloads</a>
	-
	<a href="https://tech.eus/t4/docs.php" target="_blank">Documentation</a>
	-
	<a href="https://tech.eus/t4/apps.php" target="_blank">Applications</a>
</p>

<p align="center">Free yourself from walled gardens and cloud-only SaaS offerings. Axia enables powerful low code applications, selfhosted in the cloud or on-premise. Create and then use, share or even sell your Axia applications.</p>

![DEMO - Orgas](https://github.com/user-attachments/assets/5506d0c1-4bf3-4011-bc3a-2650cb5ec0b9)
![DEMO - Gantt](https://github.com/user-attachments/assets/1e413540-f9e8-4c2f-bd91-f46f51137d8b)

## :star: Features
* **Fast results**: Quickly replace spreadsheet based 'solutions' with proper multi-user applications.
* **It can count**: Summarize records, do date calculations, apply business rules and much more.
* **Make things visible**: Show tasks on Gantt charts, generate diagrams or display information-dense lists.
* **Workflows included**: Adjust forms based on the current state of a record, export to PDF or send notifications.
* **Compliance tools**: With roles and access policies, Axia can give and restrict access globally or for specific records.
* **End-to-end encryption**: Built-in support for E2EE - easy to use with integrated key management features.
* **Integration options**: Axia can serve as and call REST endpoints, create or import CSV files and offer ICS for accessing calendars.
* **Ready for mobile**: Works well on all devices, with specific mobile settings and PWA features for great-feeling apps.
* **Fulltext search**: Users can quickly find desired content by using search phrases and language specific lookups.
* **Many inputs available**: From simple date ranges, to drawing inputs for signatures, to bar- & QR code inputs that can scan codes via camera - Axia offers a growing list of input types for various needs.
* **Blazingly fast**: Axia takes advantage of multi-core processors and communicates with clients over bi-directional data channels.
* **Security features**: Apply password policies, block brute-force attempts and enable MFA for your users.
* **Fully transparent**: Directly read and even change data in the Axia database - everything is human-readable.
* **Selfhosted**: Run Axia as you wish, locally or in the cloud - with full control on where your data is located.
* **Enterprise-ready**: Adjust Axia to your corporate identity, manage users & access via LDAP and grow with your organization by extending applications and clustering Axia.

![DEMO - PW-Safe](https://github.com/user-attachments/assets/e9161bf2-027e-409d-a9eb-ed97dfe76f7e)
![DEMO - IT-Assets](https://github.com/user-attachments/assets/c5273f72-24cb-40cc-a947-c6a42b78f7bb)
![DEMO - TimeTracker](https://github.com/user-attachments/assets/e6b6e0e9-558a-4bad-ad52-45700e7d229e)

## :rocket: Quickstart
### Linux
1. Extract the Axia package ([x64](https://tech.eus/t4/downloads.php)/[arm64](https://tech.eus/t4/downloads.php)) to any location (like `/opt/rei3`) and make the binary `r3` executable (`chmod u+x r3`).
1. Copy the file `config_template.json` to `config.json` and fill in details to an empty, UTF8 encoded Postgres database. The DB user needs full permissions to this database.
1. Install optional dependencies - ImageMagick & Ghostscript for image and PDF thumbnails (`sudo apt install imagemagick ghostscript`), PostgreSQL client utilities for integrated backups (`sudo apt install postgresql-client`).
1. Register (`sudo ./r3 -install`) and start Axia with your service manager (`sudo systemctl start rei3`).
### Windows
1. Setup the standalone version directly on any Windows Server with the portable build
1. Optionally, install [Ghostscript](https://www.ghostscript.com/) on the same Windows Server for PDF thumbnails.

Once running, Axia is available at https://localhost (default port 6129) with both username and password being `admin`. For the full documentation, visit [rei3.de](https://tech.eus/t4/docs.php).

If you plan to run Axia behind a proxy, please make sure to disable client timeouts for websockets. More details [here](https://rei3.de/en/docs/admin#proxies).

There are also Docker Compose files ([x64](https://rei3.de/docker_x64)/[arm64](https://rei3.de/docker_arm64)) and a [portable version](https://rei3.de/latest/x64_portable) for Windows available to quickly setup a test or development system.

## :bulb: Where to get help
You can visit our [community forum](https://community.rei3.de) for anything related to Axia. The full documentation is available on [rei3.de](https://rei3.de/en/docs), including documentation for [admins](https://rei3.de/en/docs/admin) and [application authors](https://rei3.de/en/docs/builder) as well as [Youtube videos](https://www.youtube.com/channel/UCKb1YPyUV-O4GxcCdHc4Csw).

## :clap: Thank you
Axia would not be possible without the help of our contributors and people using Axia and providing feedback for continuous improvement. So thank you to everybody involved with the Axia project!

[![Stargazers repo roster for @EuskadiTech/tallarin](https://reporoster.com/stars/dark/EuskadiTech/tallarin)](https://github.com/EuskadiTech/tallarin/stargazers)

Axia is built on-top of amazing open source software and technologies. Naming them all would take pages, but here are some core libraries and software that helped shape Axia:
* [Golang](https://golang.org/) to enable state-of-the-art web services and robust code even on multi-threaded systems.
* [PostgreSQL](https://www.postgresql.org/) for powerful features and the most reliable database management system we´ve ever had the pleasure to work with.
* [Vue.js](https://vuejs.org/) to provide stable and efficient frontend components and to make working with user interfaces fun.

## :+1: How to contribute
Contributions are always welcome - feel free to fork and submit pull requests.

Axia follows a four-digit versioning syntax, such as `3.2.0.4246` (MAJOR.MINOR.PATCH.BUILD). The major release will stay at `3` indefinitely, while we introduce new features and database changes with each minor release. Patch releases primarily focus on fixes, but may include small features as long as the database is not changed.

The branch `main` will contain the currently released minor version of Axia; patches for this version can directly be submitted for the main branch. Each new minor release will use a separate branch, which will merge with `main` once the latest minor version is released.

## :pick: Third party tools and resources
We want to give a shout-out to a number of projects around Axia. Often created for specific requirements, these projects have been prepared and made public by awesome people to help others do more with Axia.
1. [R3 Toolshop](https://github.com/Umb-Astardo/R3-Toolshop): A toolset for Axia operations - including data importers, bulk user creation and relation duplication.
1. [Axia-Tickets-MCP-Server](https://github.com/lgndluke/Axia-Tickets-MCP-Server): A FastMCP server for LLM integration for [Axia Tickets](https://rei3.de/en/applications/tickets).
1. [Google Material Icons for Axia](https://github.com/fmvalsera/r3_material_icons_app): A Axia application that can be built on, providing the Google Material icon pack for use in your apps.

We are humbled by the effort put into these projects and want to say thank you.

## :nut_and_bolt: Build Axia yourself
If you want to build Axia itself, you can fork this repo or download the source code to build your own executable. The master branch contains the current minor release, while new minor releases are managed in new branches.

1. Install the latest version of [Golang](https://golang.org/dl/).
1. Go into the source code directory (where `r3.go` is located) and execute: `go build -ldflags "-X main.appVersion={YOUR_APP_VERSION}"`.
   * Replace `{YOUR_APP_VERSION}` with the version of the extracted source code. Example: `go build -ldflags "-X main.appVersion=2.5.1.2980"`
   * You can change the build version anytime. If you want to upgrade the major/minor version numbers however, you need to deal with upgrading the Axia database (see `db/upgrade/upgrade.go`).
   * By setting the environment parameter `GOOS`, you can cross-compile for other systems (`GOOS=windows`, `GOOS=linux`, ...).
   * Static resource files (HTML, JS, CSS, etc.) are embedded into the binary during compilation - so changes to these files are only reflected after you recompile. Alternatively, you can use the `-wwwpath` command line argument to load Axia with an external `www` directory, in which you can make changes directly.
1. Use your new, compiled binary of Axia to replace an already installed one.
1. You can now start your own Axia version. Make sure to clear all browser caches after creating/updating your own version.

## :page_with_curl: Copyright, license & trademark
Axia (C) 2025 EuskadiTech (TM)

The Axia source code is released under the [MIT license](https://opensource.org/license/mit).

Axia is based on r3_team/r3. The official license is MIT, which allows modifications and commercial use. If any trademark issue occurs, please open a Issue or send a Email to soporte (at) tech (dot) eus
