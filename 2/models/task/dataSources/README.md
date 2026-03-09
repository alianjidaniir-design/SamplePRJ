# dataSources

Task entity data source contracts and implementations.

## Contracts

- `task_ds.go`: `TaskDBDS` and `TaskCacheDS` interfaces.

## Implementations

- `memory/task_db_ds.go`: in-memory DB datasource.
- `memory/task_cache_ds.go`: in-memory cache datasource.
- `mysql/config.go`: load MySQL settings from env.
- `mysql/connection.go`: open and tune MySQL `database/sql` pool.
- `mysql/schema.go`: task table validator + migration helper.
- `mysql/task_db_ds.go`: MySQL DB datasource implementation.
- `mysql/schema.sql`: SQL schema for manual migration.
