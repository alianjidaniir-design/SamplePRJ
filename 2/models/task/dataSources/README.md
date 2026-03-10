# dataSources

Task entity data source contracts and implementations.

## Contracts

- `taskDS.go`: `TaskDBDS` and `TaskCacheDS` interfaces.

## Implementations

- `memoryDS/taskDBDS.go`: in-memory DB datasource.
- `memoryDS/taskCacheDS.go`: in-memory cache datasource.
- `mysqlDS/config.go`: load MySQL settings from env.
- `mysqlDS/connection.go`: open and tune MySQL `database/sql` pool.
- `mysqlDS/schema.go`: task table validator + migration helper.
- `mysqlDS/taskDBDS.go`: MySQL DB datasource implementation.
- `mysqlDS/schema.sql`: SQL schema for manual migration.
