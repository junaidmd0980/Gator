# 🐊 Gator

Gator is a command-line RSS feed aggregator written in Go. It lets multiple users register, add RSS feeds, follow or unfollow feeds, and continuously fetch new posts from those feeds into a PostgreSQL database.

Built while working through the [boot.dev](https://boot.dev/) Go course.

## Contents

- [Requirements](#requirements)
- [Setup](#setup)
- [Usage](#usage)
- [Commands](#commands)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Development](#development)
- [Notes & Warnings](#notes--warnings)

## Requirements

- **Go 1.22+** — to install and build the CLI
- **PostgreSQL 15+** — stores users, feeds, follows, and posts
- **[goose](https://github.com/pressly/goose)** — runs database migrations
- **[sqlc](https://sqlc.dev/)** — only needed if you modify SQL queries

## Setup

### 1. Install Gator

```bash
go install github.com/junaidmd0980/Gator@latest
```

If `gator` is not found, add Go’s binary directory to your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add that line to `~/.bashrc` or `~/.zshrc` to make it permanent.

### 2. Create the database

```bash
createdb gator
```

If `createdb` is unavailable, connect to PostgreSQL and run:

```sql
CREATE DATABASE gator;
```

### 3. Apply migrations

Clone the repository and run the goose migrations:

```bash
git clone [https://github.com/junaidmd0980/Gator.git](https://github.com/junaidmd0980/Gator.git)
cd Gator
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

### 4. Create the config file

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Leave `current_user_name` empty. Gator updates it automatically when you register or log in.

## Usage

Create a user and add a feed:

```bash
gator register alice
gator addfeed "Boot.dev Blog" [https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)
```

Start the aggregator in one terminal:

```bash
gator agg 1m
```

In another terminal, read the latest posts:

```bash
gator browse 5
```

## Commands

| Command | Arguments | Description |
|---|---|---|
| `gator register` | `<name>` | Creates a user and logs in as them |
| `gator login` | `<name>` | Switches to an existing user |
| `gator users` | — | Lists all users and marks the current one |
| `gator reset` | — | Deletes every user, feed, follow, and post |
| `gator addfeed` | `<name> <url>` | Adds a feed and follows it; requires a logged-in user |
| `gator feeds` | — | Lists every feed and the user who added it |
| `gator follow` | `<url>` | Follows an existing feed; requires a logged-in user |
| `gator following` | — | Lists feeds followed by the current user |
| `gator unfollow` | `<url>` | Unfollows a feed; requires a logged-in user |
| `gator agg` | `<time_between_reqs>` | Fetches feeds on a loop, waiting the given duration between requests |
| `gator browse` | `[limit]` | Prints the latest posts from feeds you follow; defaults to 2 |

Example:

```bash
gator register alice
gator addfeed "Hacker News" [https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)
gator feeds
gator following
gator agg 1m
gator browse 10
```

## Tech Stack

- [Go](https://go.dev/) — CLI, HTTP requests, and XML parsing
- [PostgreSQL](https://www.postgresql.org/) — data storage
- [sqlc](https://sqlc.dev/) — generates type-safe Go code from SQL
- [goose](https://github.com/pressly/goose) — database migrations
- [lib/pq](https://github.com/lib/pq) — PostgreSQL driver
- [google/uuid](https://github.com/google/uuid) — primary key generation

## Project Structure

```text
Gator/
├── main.go               # Entry point, command registration, auth middleware
├── commands.go           # Command registry
├── middleware.go         # Logged-in user middleware
├── handler_user.go       # register, login, users
├── handler_feed.go       # addfeed, feeds, follow, unfollow
├── handler_following.go  # following
├── handler_agg.go        # agg
├── handler_browse.go     # browse
├── handler_reset.go      # reset
├── rss_feed.go           # RSS fetching and parsing
├── scrapeFeeds.go        # helper function for agg
├── sqlc.yaml             # sqlc configuration
├── sql/
│   ├── schema/           # goose migrations
│   └── queries/          # sqlc query source
└── internal/
    ├── config/           # Reads and writes ~/.gatorconfig.json
    └── database/         # sqlc-generated database code
```

## Development

Run Gator directly from a clone:

```bash
git clone [https://github.com/junaidmd0980/Gator.git](https://github.com/junaidmd0980/Gator.git)
cd Gator
go run . users
```

Build a local binary:

```bash
go build -o gator .
./gator users
```

After changing files in `sql/queries/` or `sql/schema/`, regenerate the database layer:

```bash
sqlc generate
```

New migration files should go in `sql/schema/` using the `NNN_name.sql` naming convention and standard goose markers:

```sql
-- +goose Up

-- +goose Down
```

## Notes & Warnings

- `gator reset` deletes all users, feeds, follows, and posts. There is no confirmation prompt and no undo.
- `gator agg` is a long-running process. It fetches the least recently updated feed on each tick.
- Use reasonable intervals such as `1m` or `1h`. Requesting feeds too frequently can get your IP blocked by feed servers.
- `gator login` only changes the username stored in your local config file; it does not provide real authentication.
- Your database connection string is stored in plaintext in `~/.gatorconfig.json`. Do not use production credentials locally.
