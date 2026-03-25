# CollectiveFlow Architecture

CollectiveFlow is the collective's decision-making engine. It has two interfaces (CLI and Web) that share the same storage layer. This document shows how they are built and how they fit together.

## High-Level Structure

```
+-------------------------------------------------------------------+
|                         CollectiveFlow                             |
|                                                                   |
|   +---------------------------+  +----------------------------+   |
|   |      Go CLI Binary        |  |     Flask Web App          |   |
|   |      (cmd/collectiveflow) |  |     (web/app.py)           |   |
|   |                           |  |                            |   |
|   |  cobra commands:          |  |  Routes:                   |   |
|   |   proposal create/list/   |  |   / (index)                |   |
|   |     show/update           |  |   /proposals               |   |
|   |   consensus start/input/  |  |   /proposal/<id>           |   |
|   |     complete              |  |   /create (form + POST)    |   |
|   |   status active           |  |   /dashboard               |   |
|   |   dashboard               |  |   /collective              |   |
|   |   config show/set         |  |   /about                   |   |
|   |   web start               |  |                            |   |
|   +------------+--------------+  |  API:                      |   |
|                |                 |   GET  /api/proposals       |   |
|                |                 |   GET  /api/proposals/<id>  |   |
|                |                 |   POST /api/proposals       |   |
|                |                 |   POST /api/proposals/      |   |
|                |                 |         <id>/consultation   |   |
|                |                 |   PUT  /api/proposals/      |   |
|                |                 |         <id>/status         |   |
|                |                 |   GET  /api/collective/stats|   |
|                |                 |   GET  /api/events (SSE)    |   |
|                |                 +-------------+--------------+   |
|                |                               |                  |
|   +------------v-------------------------------v--------------+   |
|   |                  Storage Layer                             |   |
|   |                                                           |   |
|   |   ProposalStore interface (Go)   StorageBackend protocol  |   |
|   |         |                              (Python)           |   |
|   |    +----+----+                    +----+----+             |   |
|   |    |         |                    |         |             |   |
|   |    v         v                    v         v             |   |
|   | FileStore  SQLiteStore      YAMLStorage  SQLiteStorage    |   |
|   | (YAML)     (SQLite)         (YAML)       (SQLite)         |   |
|   +-----------------------------------------------------------+   |
|                         |                                         |
|   +---------------------v---------------------+                   |
|   |           Shared Data                      |                   |
|   |                                            |                   |
|   |   data/proposals/*.yaml  (file backend)    |                   |
|   |   data/proposals/*.json  (API mirror)      |                   |
|   |         -- or --                           |                   |
|   |   data/collectiveflow.db (sqlite backend)  |                   |
|   +--------------------------------------------+                   |
+-------------------------------------------------------------------+
```

## Go CLI Architecture

```
cmd/collectiveflow/main.go
    |
    v
internal/cli/app.go          -- Cobra root command, Viper config
    |
    +-- proposal.go           -- proposal create/list/show/update
    +-- consensus.go          -- consensus start/input/complete
    +-- status.go             -- status active/all
    +-- dashboard.go          -- dashboard (terminal UI summary)
    +-- config.go             -- config show/set
    +-- web.go                -- web start (launches Flask)
    |
    v
internal/proposal/
    |
    +-- proposal.go           -- Proposal struct, status types, validation
    |                            Statuses: proposed -> consultation ->
    |                                      consensus -> implemented
    |                            Also: withdrawn, blocked
    |
    +-- operations.go         -- Create, List, Get, UpdateProposal,
    |                            UpdateStatus, AddConsultationInput,
    |                            RecordDecision
    |                            (all business logic lives here)
    |
    +-- storage_adapter.go    -- Wraps storage.ProposalStore to convert
                                 between interface{} and *Proposal
    |
    v
internal/storage/
    |
    +-- interface.go          -- ProposalStore interface:
    |                            Save, Load, ListAll, Delete,
    |                            GenerateID, GetFilePath
    |
    +-- file.go               -- FileStore: one YAML file per proposal
    |                            Also writes .json mirror for API access
    |                            Backup support, sequence-based IDs
    |
    +-- sqlite.go             -- SQLiteStore: shared SQLite DB
                                 WAL mode, foreign keys, audit log
                                 Same schema as web/scripts/schema.sql
```

## Flask Web Architecture

```
web/app.py                   -- Flask application, all routes
    |
    +-- storage.py            -- StorageBackend protocol + implementations
    |       |
    |       +-- YAMLStorage   -- Reads/writes YAML files directly
    |       +-- SQLiteStorage -- Reads/writes SQLite database
    |       +-- get_storage() -- Factory: picks backend from env var
    |
    +-- EventBus class        -- In-memory pub/sub for SSE
    |       |
    |       +-- subscribe()   -- Returns a queue for a new client
    |       +-- unsubscribe() -- Removes disconnected client
    |       +-- publish()     -- Broadcasts event to all queues
    |
    +-- CSRF protection       -- Session-based token, checked on POST/PUT/DELETE
    |
    +-- Security headers      -- X-Content-Type-Options, CSP, X-Frame-Options
    |
    +-- Template filters      -- humanize_date, status_emoji, urgency_color
    |
    +-- templates/
    |       +-- base.html / layout.html
    |       +-- index.html         (home with grouped proposals)
    |       +-- proposals.html     (filterable list)
    |       +-- proposal.html      (detail + consultation form)
    |       +-- create_proposal.html
    |       +-- dashboard.html     (stats, velocity, participation)
    |       +-- collective.html    (collective state overview)
    |       +-- about.html
    |
    +-- static/css/style.css  -- Custom styles (Tailwind via CDN)
```

## Storage Schema (SQLite)

Both Go and Python use the same schema. The Go binary embeds it inline; Python reads from `scripts/schema.sql`.

```
+-------------------+       +--------------------+
|    proposals      |       |   consultations    |
|-------------------|       |--------------------|
| id (PK, TEXT)     |<------| proposal_id (FK)   |
| title             |       | contributor        |
| description       |       | created_at         |
| proposer          |       | input              |
| created_at        |       | support (0/1)      |
| status            |       | concerns (JSON[])  |
| urgency           |       +--------------------+
| consensus_status  |
| affected_areas    |       +--------------------+
| (JSON[])          |       | consensus_events   |
+-------------------+       |--------------------|
        |                   | proposal_id (FK)   |
        |                   | event_type         |
        +------------------>| actor              |
        |                   | details            |
        |                   | created_at         |
        |                   +--------------------+
        |
        |                   +--------------------+
        |                   |    decisions       |
        +------------------>|--------------------|
                            | proposal_id (FK,   |
                            |   UNIQUE)          |
                            | result             |
                            | rationale          |
                            | created_at         |
                            +--------------------+

+-------------------+       +--------------------+
|    audit_log      |       |  schema_version    |
|-------------------|       |--------------------|
| table_name        |       | version (PK)       |
| row_id            |       | applied_at         |
| operation         |       | description        |
| old_values        |       +--------------------+
| new_values        |
| actor             |
| created_at        |
+-------------------+
```

Valid status values: `proposed`, `consultation`, `consensus`, `implemented`, `withdrawn`, `blocked`

Valid urgency values: `low`, `medium`, `high`, `emergency`

Valid decision results: `approved`, `rejected`, `deferred`, `no_consensus`

## Status State Machine

```
                +------------+
                |  proposed  |
                +-----+------+
                      |
            +---------+---------+
            |                   |
            v                   v
    +---------------+    +-----------+
    | consultation  |    | withdrawn |  (terminal)
    +-------+-------+    +-----------+
            |
    +-------+-------+--------+
    |               |        |
    v               v        v
+----------+  +---------+  +-----------+
| consensus|  | blocked |  | withdrawn |
+----+-----+  +----+----+  +-----------+
     |             |
     +------+------+
     |      |
     v      v
+-----------+    +---------------+
|implemented|    | consultation  |  (retry)
+-----------+    +---------------+
  (terminal)
```

## Choosing a Storage Backend

Both CLI and Web support two backends. Set via environment variable or CLI flag:

| Backend | CLI Flag | Env Var | What It Does |
|---------|----------|---------|-------------|
| YAML (default) | `--storage file` | `STORAGE_BACKEND=yaml` | One `.yaml` file per proposal in `data/proposals/` |
| SQLite | `--storage sqlite` | `STORAGE_BACKEND=sqlite` | All data in `data/collectiveflow.db` |

The SQLite backend is the better choice when both CLI and Web need to read/write the same data concurrently, because SQLite handles locking. YAML is better for simplicity, git-friendliness, and human inspection.

## Go Embedded Web Server

The Go binary includes its own web server (`internal/web/server.go`) with embedded templates and static files. This is a separate, lighter web interface from the Flask app. It reads from the same storage backends and provides:

- Homepage with stats
- Proposals list with filters
- Individual proposal view
- Collective status page
- About page

This lets you run `collectiveflow web start` without installing Python.
