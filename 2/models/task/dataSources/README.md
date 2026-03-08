# datasourse

Task entity datasourse implementations should be placed here.

## MySQL datasource

- `mysql/config.go`: load MySQL settings from env.
- `mysql/connection.go`: open and tune MySQL `database/sql` pool.
- `mysql/schema.go`: task table validator + migration helper.
- `mysql/schema.sql`: SQL schema for manual migration.
