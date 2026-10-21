<div align="center">

# 🐊 Gator

**A command-line RSS feed aggregator built with Go and PostgreSQL.**

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

Gator lets multiple users register, add RSS feeds, follow or unfollow feeds, and continuously fetch new posts into a PostgreSQL database.

Built while working through the [boot.dev](https://boot.dev/) Go course.

</div>

---

## 📑 Contents

- [Requirements](#-requirements)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
- [Commands](#-commands)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Notes & Warnings](#%EF%B8%8F-notes--warnings)

---

## 📦 Requirements

| Tool | Version | Purpose |
|---|---:|---|
| [Go](https://go.dev/) | `1.22+` | Install and build the CLI |
| [PostgreSQL](https://www.postgresql.org/) | `15+` | Stores users, feeds, follows, and posts |
| [goose](https://github.com/pressly/goose) | Latest | Runs database migrations |
| [sqlc](https://sqlc.dev/) | Latest | Regenerates database code after SQL changes |

---

## 🚀 Quick Start

### 1. Install Gator

```bash
go install github.com/junaidmd0980/Gator@latest
```

If the `gator` command is not found, add Go’s binary directory to your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

To make this permanent, add the line above to `~/.bashrc` or `~/.zshrc`.

### 2. Create the database

```bash
createdb gator
```

If `createdb` is unavailable, connect to PostgreSQL and run:

```sql
CREATE DATABASE gator;
```

### 3. Apply migrations

```bash
git clone [https://github.com/junaidmd0980/Gator.git](https://github.com/junaidmd0980/Gator.git)
cd Gator

goose -dir sql/schema postgres \
  "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" \
  up
```

### 4. Create the config file

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

> [!NOTE]  
> Leave `current_user_name` empty. Gator updates it automatically when you register or log in.

---

## 🧭 Usage

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

---

## ⌨️ Commands

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
| `gator browse` | `[limit]` | Prints the latest posts from followed feeds; defaults to `2` |

### Example

```bash
gator register alice
gator addfeed "Hacker News" [https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)
gator feeds
gator following
gator agg 1m
gator browse 10
```

---

## 🛠 Tech Stack

- [Go](https://go.dev/) — CLI, HTTP requests, and XML parsing
- [PostgreSQL](https://www.postgresql.org/) — data storage
- [sqlc](https://sqlc.dev/) — generates type-safe Go code from SQL
- [goose](https://github.com/pressly/goose) — database migrations
- [lib/pq](https://github.com/lib/pq) — PostgreSQL driver
- [google/uuid](https://github.com/google/uuid) — primary key generation

---

## 🗂 Project Structure

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
├── scrapeFeeds.go        # Helper function for agg
├── sqlc.yaml             # sqlc configuration
├── sql/
│   ├── schema/           # Goose migrations
│   └── queries/          # sqlc query source
└── internal/
    ├── config/           # Reads and writes ~/.gatorconfig.json
    └── database/         # sqlc-generated database code
```

---

## 👨‍💻 Development

Run Gator directly from a local clone:

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

New migration files should go in `sql/schema/` using the `NNN_name.sql` naming convention and standard Goose markers:

```sql
-- +goose Up

-- +goose Down
```

---

## ⚠️ Notes & Warnings

> [!WARNING]  
> `gator reset` deletes all users, feeds, follows, and posts. There is no confirmation prompt and no undo.

> [!IMPORTANT]  
> `gator agg` is a long-running process. It fetches the least recently updated feed on each tick.

- Use reasonable intervals such as `1m` or `1h`. Requesting feeds too frequently can get your IP blocked by feed servers.
- `gator login` only changes the username stored in your local config file; it does not provide real authentication.
- Your database connection string is stored in plaintext in `~/.gatorconfig.json`. Do not use production credentials locally.

---

<div align="center">

**Built with Go, PostgreSQL, sqlc, and goose.**

</div>
