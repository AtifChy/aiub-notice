# AIUB Notice Fetcher

AIUB Notice Fetcher is a cross-platform CLI tool and background service for automatically fetching, caching,
and notifying users about new notices from the AIUB website.

## Features

- Periodically checks for new notices from the AIUB website
- Caches fetched notices locally for offline access
- Sends desktop notifications for new notices
- Tracks seen notices to avoid duplicate notifications
- CLI commands to view the last notice, manage autostart, and more
- Supports autostart on Windows

## Installation

1. Clone this repository:

   ```sh
   git clone https://github.com/AtifChy/aiub-notice.git
   cd aiub-notice
   ```

2. Build the project:

   ```sh
   go build -o aiub-notice.exe
   ```

## Usage

### Start the Service

```sh
./aiub-notice start
```

- Use `--interval` or `-i` to set the check interval (default: 30m).

### Show Last Notice

```sh
./aiub-notice last
```

**Note:** This command will show the last fetched notice, or an error if no notices have been fetched yet.

### Manage Autostart (Windows)

```sh
./aiub-notice autostart --enable   # Enable autostart
./aiub-notice autostart --disable  # Disable autostart
./aiub-notice autostart --status   # Show autostart status
```

## Project Structure

- `cmd/` — CLI commands
- `internal/service/` — Main service logic (periodic checks, notifications)
- `internal/notice/` — Notice fetching, parsing, and caching
- `internal/toast/` — Notification logic
- `internal/common/` — Shared constants and helpers
- `internal/autostart/` — Autostart management (Windows)

## Contributing

Pull requests and issues are welcome!

## License

[MIT](./LICENSE)

