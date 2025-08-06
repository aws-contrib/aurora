# aurora

Manage your **Aurora DSQL** schema as code.

This command-line tool was built because [Atlas](https://atlasgo.io/) currently [does not support](https://github.com/ariga/atlas/issues/3539) [AWS Aurora DSQL](https://docs.aws.amazon.com/aurora-dsql/latest/userguide/what-is-aurora-dsql.html).

## Getting Started

### 1. Create a configuration file

- Create an [`aurora.hcl`](./cmd/aurora.hcl) file.
- Follow the same structure you would use with Atlas.

### 2. Migrations

- The CLI tool **does not generate migration files**. Use **Atlas** to create migrations.
- **Migrations must be idempotent** because:
  - Aurora DSQL does **not** allow executing a single transaction with multiple DDL and DML statements.
  - Each migration must be safe to run multiple times.
  - Idempotency prevents:
    - Errors from reapplying statements.
    - Duplicate schema changes.
    - Unintended side effects.
- Recommended practices for idempotent migrations:
  - Use `IF NOT EXISTS` for `CREATE` statements.
  - Use `IF EXISTS` for `DROP` statements.
  - Avoid operations that cannot be safely re-run.
- **All SQL statements in migration files must end with `;`** because:
  - The `aurora` CLI splits migration files by `;`.
  - Aurora DSQL does **not** support executing multiple DDL or DML statements in a single query, so each statement must be executed individually.

#### Example: Idempotent migration

```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY,
  name TEXT
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT;
```

### 3. Index creation for local compatibility

- To ensure SQL scripts are compatible with both **PostgreSQL** and **Aurora DSQL**:
  - **Create indexes using `CONCURRENTLY` in PostgreSQL** so they can run locally outside of a transaction.
  - During execution in Aurora DSQL, the CLI tool will **replace `CONCURRENTLY` with `ASYNC`** for compatibility.
- This ensures that:
  - Local development works with PostgreSQL.
  - The same scripts are valid and executable in Aurora DSQL.

#### Example: Index creation

```sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users (email);
```

When executed in Aurora DSQL, this will be translated to:

```sql
CREATE INDEX ASYNC IF NOT EXISTS idx_users_email ON users (email);
```

### 4. Apply migrations

- Once your idempotent migrations are ready, run the following command to apply them:

```bash
aurora migrate --env aws apply
```

- This command:
  - Connects to your Aurora DSQL environment.
  - Ensures index creation syntax is correct.
  - Applies the migrations safely in the correct order. Waits all indexes to be created.

## Installation in Docker

You can install `aurora` in your container by using a multi-stage Docker build.
The first stage pulls the `aurora` binary from the published image, and the second stage copies it into a minimal container.

### Example `Dockerfile`

```dockerfile
FROM ghcr.io/aws-contrib/aurora:edge AS aurora
FROM scratch

WORKDIR /app

COPY --from=aurora /bin/aurora /app/aurora
```
