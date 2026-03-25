# YAML to SQLite Migration Log

**Date**: 2026-03-24
**Performed by**: database-design-specialist agent
**Branch**: `collective/data-migration`

## Summary

Migrated all CollectiveFlow proposal data from YAML files to a SQLite database
using the migration script at `scripts/migrate_to_sqlite.py` and the schema at
`scripts/schema.sql`.

## Migration Results

### Import

| Metric | Count |
|--------|-------|
| YAML files processed | 13 |
| Proposals imported | 13 |
| Consultations imported | 24 |
| Consensus events imported | 52 |
| Decisions imported | 4 |
| Audit log entries created | 13 |
| Errors | 0 |

### Verification

The `--verify` flag compared every YAML file against the database and found
**zero mismatches** across all checked fields: titles, statuses, consultation
counts, consensus event counts, and decision results.

### Web App Test

The Flask web app was tested with `STORAGE_BACKEND=sqlite` pointing at the
migrated database. The `/api/proposals` endpoint returned HTTP 200 with all 13
proposals.

### Data Issues Found

None. All 13 YAML files parsed cleanly and mapped to the schema without data
loss or type conversion problems.

## Database Details

| Property | Value |
|----------|-------|
| File | `data/collectiveflow.db` |
| Size | 172 KB |
| Schema version | 1 |
| Journal mode | WAL |
| Foreign keys | Enforced |

### Tables

| Table | Rows |
|-------|------|
| proposals | 13 |
| consultations | 24 |
| consensus_events | 52 |
| decisions | 4 |
| audit_log | 13 |
| schema_version | 1 |

### Views

- `active_proposals` -- proposals not in `implemented` or `withdrawn` status,
  with consultation and support counts
- `proposal_details` -- proposals joined with their decision (if any)
- `agent_activity` -- per-contributor consultation summary

### Agent Participation (from `agent_activity` view)

| Contributor | Consultations | Supports | Blocks |
|-------------|--------------|----------|--------|
| go-systems-developer | 5 | 5 | 0 |
| product-steward | 5 | 5 | 0 |
| devops-coordinator | 4 | 4 | 0 |
| flask-web-developer | 4 | 4 | 0 |
| consensus-coordinator | 2 | 2 | 0 |
| david-graeber-agent | 2 | 2 | 0 |
| noam-chomsky-agent | 2 | 2 | 0 |

## How to Switch Between YAML and SQLite

### Using YAML (current default)

No configuration needed. The web app reads from `data/proposals/*.yaml` by
default:

```bash
cd projects/collectiveflow/web
python3 app.py
```

### Using SQLite

Set two environment variables before starting the web app:

```bash
export STORAGE_BACKEND=sqlite
export SQLITE_DB_PATH=../data/collectiveflow.db
python3 app.py
```

Or inline:

```bash
STORAGE_BACKEND=sqlite SQLITE_DB_PATH=../data/collectiveflow.db python3 app.py
```

### Re-running the Migration

If YAML files change and you need to refresh the database:

```bash
# Remove old database
rm data/collectiveflow.db

# Re-import
python3 scripts/migrate_to_sqlite.py --import \
  --data-dir data/proposals \
  --db data/collectiveflow.db

# Verify
python3 scripts/migrate_to_sqlite.py --verify \
  --data-dir data/proposals \
  --db data/collectiveflow.db
```

### Exporting SQLite Back to YAML

The migration is reversible:

```bash
python3 scripts/migrate_to_sqlite.py --export \
  --data-dir data/proposals-export \
  --db data/collectiveflow.db
```

### Dry Run (Validate Without Writing)

```bash
python3 scripts/migrate_to_sqlite.py --import \
  --data-dir data/proposals \
  --db data/collectiveflow.db \
  --dry-run
```

## Notes

- The YAML files remain the source of truth. The SQLite database is a derived
  artifact that can be regenerated at any time.
- Both storage backends are equal peers -- neither has authority over the other,
  consistent with the collective's horizontal principles.
- The database file (`data/collectiveflow.db`) is committed so any agent can use
  SQLite immediately without running the migration.
