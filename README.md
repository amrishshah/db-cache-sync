# DB Cache Sync

**DB Cache Sync** reads MySQL binlogs asynchronously using [go-mysql/canal](https://github.com/go-mysql-org/go-mysql) and updates Redis cache based on database changes (Insert/Update/Delete).

---

## Features

- Read MySQL binlog events (insert, update, delete)
- Update Redis cache asynchronously
- Auto-reconnect on MySQL/Redis failures
- Lightweight

---

## Requirements

- Go 1.20+
- MySQL (Binlog enabled, Row format)
- Redis server

---

## Installation

```bash
git clone https://your-repository-link/db-cache-sync.git
cd db-cache-sync
go mod tidy
go run main.go
```

---

## Configuration

Edit `main.go` to update your database and Redis connection settings:

```go
cfg.Addr = "host:port"
cfg.User = "your_mysql_user"
cfg.Password = "your_mysql_password"
cfg.Flavor = "mysql"
cfg.ServerID = 1001 // random unique number

// Redis config
redisClient = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
```

---

## MySQL Binlog Settings (important)

Make sure MySQL server is configured properly:

```sql
SHOW VARIABLES LIKE 'binlog_format'; -- Should be ROW
SHOW VARIABLES LIKE 'binlog_row_image'; -- Should be FULL
```

If not, set in your `my.cnf`:

```ini
[mysqld]
log-bin=mysql-bin
server-id=1
binlog-format=ROW
binlog-row-image=FULL
```

Restart MySQL after changes.

---

## Usage

Simply run:

```bash
go run main.go
```

The service will listen to binlogs and update Redis automatically.

Example Redis Key format:

- Insert: `cache:{table_name}:{primary_key}` → JSON value
- Update: `cache:{table_name}:{primary_key}` → updated JSON
- Delete: Redis key deleted

---

## Error Handling

If MySQL dump fails with:

```
Couldn't execute 'SET SQL_QUOTE_SHOW_CREATE=1...': Variable 'sql_mode' can't be set...
```

**Solution:**  
Set `DumpExecPath` to empty and skip initial snapshot in code:

```go
canalCfg.Dump.ExecutionPath = ""
canalCfg.Dump.TableDB = ""
```

You can also disable `Dump` mode entirely if binlog position is already known.

---

## Libraries Used

- [go-mysql/canal](https://github.com/go-mysql-org/go-mysql)
- [go-redis/redis](https://github.com/go-redis/redis/v8)

---

## License

This project is open-sourced under the MIT license.

