# AIUB Notice Fetcher

AIUB Notice Fetcher is a CLI tool and background service for automatically fetching,
and notifying users about new notices from the AIUB website.

## Preview

![Windows Toast Notification Example](https://github.com/user-attachments/assets/411999be-0da6-4a23-9fd3-15bd97c7b44f)

```text
$ aiub-notice.exe --help
AIUB Notice Notifier is a command-line tool that fetches and
displays notices from AIUB's official website.

Usage:
  aiub-notice [command]

Available Commands:
  appid       Manage AppID registration for Windows notifications
  autostart   Manage autostart settings for AIUB Notice Fetcher service
  close       Close the AIUB Notice Fetcher service
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  last        Display the last fetched notice
  list        List all fetched notices
  log         View the log of notices
  start       Start the AIUB Notice Fetcher service
  status      Check the status of the AIUB Notice Fetcher service

Flags:
  -h, --help      help for aiub-notice
  -v, --version   version for aiub-notice

Use "aiub-notice [command] --help" for more information about a command.
```

## Features

- Periodically checks for new notices from the AIUB website
- Caches fetched notices locally for offline access
- Sends desktop notifications for new notices
- Tracks seen notices to avoid duplicate notifications
- CLI commands to view the last notice, manage autostart, and more
- Supports autostart on Windows

## Requirements

- Windows 10/11
- Go 1.25 or later
- GNU make
- Git

## Installation

### Install dependencies

```sh
scoop install make go git
```

**Note:** If you don't have [Scoop](https://scoop.sh/) installed,
please follow the instructions on their website to install it.
Or you can install Go and Make manually.

### Steps

1. Clone this repository:

   ```sh
   git clone https://github.com/AtifChy/aiub-notice.git
   cd aiub-notice
   ```

2. Install the project:

   ```sh
   make install-all
   ```

3. Restart (or Sign out) your computer to ensure that autostart works correctly.

4. Profit!

## Usage

### Start the Service

```sh
aiub-notice start
```

- Use `--interval` or `-i` to set the custom check interval (default: 30m).

### Show Last Notice

```sh
aiub-notice last
```

**Note:** This command will show the last fetched notice,
or an error if no notices have been fetched yet.

### Register

To register the program and ensure that toast notifications display
the correct icon and name, run the following command once:

```sh
aiub-notice appid --register
```

**Note:** Registration is recommended before using other features.

### Manage Autostart (Windows)

```sh
aiub-notice autostart --enable   # Enable autostart
aiub-notice autostart --disable  # Disable autostart
aiub-notice autostart --status   # Show autostart status
```

## Project Structure

- `cmd/` — Entrypoints for CLI applications and subcommands
  - `aiub-notice/` — Main CLI application
  - `aiub-notice-launcher/` — Launcher utility
- `internal/appid/` — AppID registration for Windows notifications
- `internal/autostart/` — Windows autostart management
- `internal/common/` — Shared constants, paths, and helpers
- `internal/list/` — Notice List TUI
- `internal/notice/` — Notice fetching, parsing, caching, and seen notice tracking
- `internal/service/` — Main service logic: periodic checks, notifications
- `internal/toast/` — Windows Toast notification logic and icon handling

## Contributing

Pull requests and issues are welcome!

## License

[MIT](LICENSE)
