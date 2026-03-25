"""
CollectiveFlow Storage Abstraction Layer

This module implements the Strategy pattern for proposal storage, allowing the
Flask web app to swap between YAML file storage and SQLite database storage
without changing any route or template code.

Architecture overview:
    StorageBackend (Protocol)    -- defines the interface every backend must satisfy
        |
        +-- YAMLStorage          -- reads/writes YAML files (the original behavior)
        +-- SQLiteStorage        -- reads/writes a SQLite database (new, opt-in)

How to choose a backend at runtime:
    Set the environment variable STORAGE_BACKEND:
        STORAGE_BACKEND=yaml     (default)  -- uses YAMLStorage
        STORAGE_BACKEND=sqlite              -- uses SQLiteStorage

    The get_storage() factory function reads this variable and returns the
    appropriate backend instance.  All Flask route code calls get_storage()
    instead of touching files directly, so the rest of the app is backend-agnostic.

Why a Protocol and not an ABC?
    Protocols use structural subtyping (duck typing).  Any class that has the
    right methods satisfies the protocol, even without inheriting from it.
    This keeps the code simple and avoids import-time coupling between backends.

Teaching notes (for any agent extending this):
    1. To add a new backend (e.g., PostgreSQL), create a class with the same
       three methods: load_proposals(), get_proposal(id), save_proposal(data).
    2. Register it in get_storage() under a new STORAGE_BACKEND value.
    3. That's it.  No base class to import, no registration decorator.
"""

import json
import logging
import os
import re
import sqlite3
import uuid
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional, Protocol, runtime_checkable

import yaml

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# Protocol: the contract every storage backend must satisfy
# ---------------------------------------------------------------------------

@runtime_checkable
class StorageBackend(Protocol):
    """
    The interface for proposal storage.

    Any class that implements these three methods can be used as a storage
    backend.  The Flask routes only depend on this protocol, never on a
    concrete class.

    Methods:
        load_proposals() -> List[dict]
            Return all proposals, sorted newest-first by date.

        get_proposal(proposal_id: str) -> Optional[dict]
            Return a single proposal by ID, or None if not found.

        save_proposal(proposal_data: dict) -> str
            Persist a new proposal and return its generated ID.
    """

    def load_proposals(self) -> List[dict]:
        ...

    def get_proposal(self, proposal_id: str) -> Optional[dict]:
        ...

    def save_proposal(self, proposal_data: dict) -> str:
        ...


# ---------------------------------------------------------------------------
# YAML Storage (the original implementation, extracted from app.py)
# ---------------------------------------------------------------------------

class YAMLStorage:
    """
    File-based storage using one YAML file per proposal.

    This is the original storage mechanism from CollectiveFlow.  Each proposal
    lives in its own .yaml file under the proposals directory.  Human-readable,
    git-friendly, and zero-dependency beyond PyYAML.

    Directory layout:
        {data_dir}/proposals/
            proposal-2025-07-26-001.yaml
            proposal-2025-07-26-001.json   (API mirror, written on save)
            ...
    """

    def __init__(self, data_dir: str):
        """
        Args:
            data_dir: Path to the CollectiveFlow data directory (contains
                      a 'proposals/' subdirectory with YAML files).
        """
        self.proposals_dir = Path(data_dir) / "proposals"

    def load_proposals(self) -> List[dict]:
        """
        Load all proposals from YAML files, sorted newest-first by date.

        Silently skips files that fail to parse so one corrupt file doesn't
        take down the whole listing.
        """
        proposals: List[dict] = []

        if not self.proposals_dir.exists():
            return proposals

        for yaml_file in self.proposals_dir.glob("*.yaml"):
            try:
                with open(yaml_file, "r") as f:
                    proposal = yaml.safe_load(f)
                    if proposal:
                        proposals.append(proposal)
            except Exception as e:
                logger.warning("Error loading %s: %s", yaml_file, e)

        def sort_key(p: dict) -> str:
            date = p.get("date", "")
            if isinstance(date, datetime):
                return date.isoformat()
            return str(date)

        proposals.sort(key=sort_key, reverse=True)
        return proposals

    def get_proposal(self, proposal_id: str) -> Optional[dict]:
        """
        Load a single proposal by ID.

        Security: rejects IDs containing path-traversal characters.
        Only alphanumeric, hyphens, and underscores are allowed.
        """
        if not re.match(r"^[a-zA-Z0-9_-]+$", proposal_id):
            return None

        yaml_path = self.proposals_dir / f"{proposal_id}.yaml"

        # Defense in depth: verify resolved path stays within proposals dir
        try:
            yaml_path.resolve().relative_to(self.proposals_dir.resolve())
        except ValueError:
            return None

        if yaml_path.exists():
            with open(yaml_path, "r") as f:
                return yaml.safe_load(f)

        return None

    def save_proposal(self, proposal_data: dict) -> str:
        """
        Save a new proposal to disk as both YAML and JSON.

        Generates an ID if one isn't already set, stamps metadata
        (date, status, consensus_history), and writes two files.

        Returns the proposal ID.
        """
        self.proposals_dir.mkdir(parents=True, exist_ok=True)

        if "id" not in proposal_data:
            proposal_data["id"] = (
                f"proposal-{datetime.now().strftime('%Y-%m-%d')}"
                f"-{str(uuid.uuid4())[:8]}"
            )

        proposal_data["date"] = datetime.now().isoformat()
        proposal_data["status"] = "proposed"
        proposal_data["consensus_status"] = "New proposal submitted"
        proposal_data["consensus_history"] = [
            {
                "timestamp": proposal_data["date"],
                "event": "proposal_created",
                "actor": proposal_data.get("proposer", "web-user"),
                "details": (
                    f"Created with urgency: "
                    f"{proposal_data.get('urgency', 'medium')}"
                ),
            }
        ]
        proposal_data["consultations"] = []

        yaml_path = self.proposals_dir / f"{proposal_data['id']}.yaml"
        with open(yaml_path, "w") as f:
            yaml.safe_dump(proposal_data, f, default_flow_style=False, sort_keys=False)

        json_path = self.proposals_dir / f"{proposal_data['id']}.json"
        with open(json_path, "w") as f:
            json.dump(proposal_data, f, indent=2)

        return proposal_data["id"]


# ---------------------------------------------------------------------------
# SQLite Storage (new backend, opt-in via STORAGE_BACKEND=sqlite)
# ---------------------------------------------------------------------------

class SQLiteStorage:
    """
    SQLite-backed storage using the schema from scripts/schema.sql.

    This backend stores proposals, consultations, consensus events, decisions,
    and a full audit log in a single SQLite file.  It is designed to be a
    drop-in replacement for YAMLStorage: the dict shapes returned by
    load_proposals() and get_proposal() match the YAML format exactly, so
    templates and API routes work unchanged.

    The database file location is controlled by the SQLITE_DB_PATH environment
    variable (default: {data_dir}/collectiveflow.db).

    Key design decisions:
        - WAL mode for concurrent readers (the web server) without blocking.
        - Foreign keys enforced so orphaned consultations can't exist.
        - Every write is recorded in the audit_log table for traceability.
        - JSON arrays for affected_areas and concerns (sparse, rarely queried).
        - The schema is applied automatically on first connection if the
          database file doesn't exist yet.
    """

    def __init__(self, data_dir: str, db_path: Optional[str] = None):
        """
        Args:
            data_dir: Path to the CollectiveFlow data directory.  Used to
                      locate schema.sql and as default db location.
            db_path:  Explicit path to the SQLite database file.  If None,
                      defaults to {data_dir}/collectiveflow.db or the
                      SQLITE_DB_PATH environment variable.
        """
        self.data_dir = Path(data_dir)

        if db_path:
            self.db_path = Path(db_path)
        else:
            self.db_path = Path(
                os.environ.get(
                    "SQLITE_DB_PATH",
                    str(self.data_dir / "collectiveflow.db"),
                )
            )

        self._ensure_schema()

    def _get_connection(self) -> sqlite3.Connection:
        """
        Open a new connection with recommended pragmas.

        Each call returns a fresh connection.  For a web app served by a WSGI
        server (gunicorn, etc.), each request gets its own connection, which
        is the simplest correct approach for SQLite.
        """
        conn = sqlite3.connect(str(self.db_path))
        conn.execute("PRAGMA journal_mode = WAL")
        conn.execute("PRAGMA foreign_keys = ON")
        conn.row_factory = sqlite3.Row  # rows behave like dicts
        return conn

    def _ensure_schema(self) -> None:
        """
        Apply the schema if the database is empty or doesn't exist yet.

        Looks for schema.sql in the scripts/ directory relative to the
        collectiveflow project root.  If the schema file isn't found, the
        storage will still work with an existing database, but won't be able
        to create a new one from scratch.
        """
        if self.db_path.exists():
            # Database already exists; assume schema is applied.
            return

        # Find schema.sql -- it lives in scripts/ next to the web/ directory
        schema_candidates = [
            self.data_dir.parent / "scripts" / "schema.sql",
            Path(__file__).parent.parent / "scripts" / "schema.sql",
        ]

        schema_sql = None
        for candidate in schema_candidates:
            if candidate.exists():
                schema_sql = candidate.read_text(encoding="utf-8")
                break

        if schema_sql is None:
            logger.warning(
                "schema.sql not found; cannot auto-create database. "
                "Looked in: %s",
                [str(c) for c in schema_candidates],
            )
            return

        self.db_path.parent.mkdir(parents=True, exist_ok=True)
        conn = sqlite3.connect(str(self.db_path))
        conn.executescript(schema_sql)
        conn.close()
        logger.info("Created new SQLite database at %s", self.db_path)

    # ------- Read operations -------

    def load_proposals(self) -> List[dict]:
        """
        Load all proposals with their consultations, events, and decisions.

        Returns a list of dicts in the same shape as the YAML files, sorted
        newest-first.  Each dict contains nested 'consultations',
        'consensus_history', and optionally 'decision' keys.
        """
        conn = self._get_connection()
        try:
            rows = conn.execute(
                "SELECT * FROM proposals ORDER BY created_at DESC"
            ).fetchall()

            proposals = [self._row_to_proposal(conn, row) for row in rows]
            return proposals
        finally:
            conn.close()

    def get_proposal(self, proposal_id: str) -> Optional[dict]:
        """
        Load a single proposal by ID.

        Security: rejects IDs containing path-traversal characters, same
        as YAMLStorage, even though SQLite isn't vulnerable to path traversal.
        This keeps the interface contract consistent.
        """
        if not re.match(r"^[a-zA-Z0-9_-]+$", proposal_id):
            return None

        conn = self._get_connection()
        try:
            row = conn.execute(
                "SELECT * FROM proposals WHERE id = ?", (proposal_id,)
            ).fetchone()

            if row is None:
                return None

            return self._row_to_proposal(conn, row)
        finally:
            conn.close()

    def _row_to_proposal(self, conn: sqlite3.Connection, row: sqlite3.Row) -> dict:
        """
        Convert a proposals table row (plus related data) into the dict
        format that templates and API routes expect.

        This is the bridge between the normalized SQLite schema and the
        flat YAML structure.  The returned dict looks identical to what
        yaml.safe_load() returns from a proposal file.
        """
        proposal: Dict[str, Any] = {
            "id": row["id"],
            "title": row["title"],
            "description": row["description"],
            "proposer": row["proposer"],
            "date": row["created_at"],
            "status": row["status"],
            "urgency": row["urgency"],
            "affected_areas": json.loads(row["affected_areas"]),
            "consensus_status": row["consensus_status"],
        }

        # Consensus history (events)
        events = conn.execute(
            """
            SELECT event_type, actor, details, created_at
            FROM consensus_events
            WHERE proposal_id = ?
            ORDER BY created_at ASC, id ASC
            """,
            (row["id"],),
        ).fetchall()

        proposal["consensus_history"] = [
            {
                "timestamp": e["created_at"],
                "event": e["event_type"],
                "actor": e["actor"],
                "details": e["details"],
            }
            for e in events
        ]

        # Consultations
        consults = conn.execute(
            """
            SELECT contributor, created_at, input, support, concerns
            FROM consultations
            WHERE proposal_id = ?
            ORDER BY created_at ASC, id ASC
            """,
            (row["id"],),
        ).fetchall()

        consultations_list = []
        for c in consults:
            entry: Dict[str, Any] = {
                "contributor": c["contributor"],
                "timestamp": c["created_at"],
                "input": c["input"],
                "support": bool(c["support"]),
            }
            concerns = json.loads(c["concerns"])
            if concerns:
                entry["concerns"] = concerns
            consultations_list.append(entry)
        proposal["consultations"] = consultations_list

        # Decision (at most one per proposal)
        decision = conn.execute(
            """
            SELECT result, rationale, created_at
            FROM decisions
            WHERE proposal_id = ?
            """,
            (row["id"],),
        ).fetchone()

        if decision:
            proposal["decision"] = {
                "result": decision["result"],
                "timestamp": decision["created_at"],
                "rationale": decision["rationale"],
            }

        return proposal

    # ------- Write operations -------

    def save_proposal(self, proposal_data: dict) -> str:
        """
        Save a new proposal to the SQLite database.

        Inserts rows into the proposals, consensus_events, and audit_log
        tables.  Returns the generated proposal ID.

        The proposal_data dict should contain at minimum 'title' and
        'description'.  Other fields (id, date, status, etc.) are generated
        automatically, matching the behavior of YAMLStorage.save_proposal().
        """
        if "id" not in proposal_data:
            proposal_data["id"] = (
                f"proposal-{datetime.now().strftime('%Y-%m-%d')}"
                f"-{str(uuid.uuid4())[:8]}"
            )

        now = datetime.now().isoformat()
        proposal_data["date"] = now
        proposal_data["status"] = "proposed"
        proposal_data["consensus_status"] = "New proposal submitted"

        proposer = proposal_data.get("proposer", "web-user")
        urgency = proposal_data.get("urgency", "medium")
        affected_areas = proposal_data.get("affected_areas", [])
        if isinstance(affected_areas, str):
            affected_areas = [affected_areas]

        conn = self._get_connection()
        try:
            # Insert proposal
            conn.execute(
                """
                INSERT INTO proposals
                    (id, title, description, proposer, created_at, status,
                     urgency, consensus_status, affected_areas)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    proposal_data["id"],
                    proposal_data.get("title", ""),
                    proposal_data.get("description", ""),
                    proposer,
                    now,
                    "proposed",
                    urgency,
                    "New proposal submitted",
                    json.dumps(affected_areas),
                ),
            )

            # Insert the creation event into consensus_events
            conn.execute(
                """
                INSERT INTO consensus_events
                    (proposal_id, event_type, actor, details, created_at)
                VALUES (?, ?, ?, ?, ?)
                """,
                (
                    proposal_data["id"],
                    "proposal_created",
                    proposer,
                    f"Created with urgency: {urgency}",
                    now,
                ),
            )

            # Audit log entry
            conn.execute(
                """
                INSERT INTO audit_log
                    (table_name, row_id, operation, new_values, actor)
                VALUES (?, ?, 'INSERT', ?, ?)
                """,
                (
                    "proposals",
                    proposal_data["id"],
                    json.dumps({
                        "id": proposal_data["id"],
                        "title": proposal_data.get("title", ""),
                    }),
                    proposer,
                ),
            )

            conn.commit()
        except Exception:
            conn.rollback()
            raise
        finally:
            conn.close()

        return proposal_data["id"]


# ---------------------------------------------------------------------------
# Factory: pick the right backend based on environment configuration
# ---------------------------------------------------------------------------

def get_storage(data_dir: Optional[str] = None) -> StorageBackend:
    """
    Factory function that returns the configured storage backend.

    Reads the STORAGE_BACKEND environment variable:
        "yaml"   (default) -> YAMLStorage
        "sqlite"           -> SQLiteStorage

    Args:
        data_dir: Path to the CollectiveFlow data directory.  If None,
                  reads COLLECTIVEFLOW_DATA env var (default: "../data").

    Returns:
        An object satisfying the StorageBackend protocol.

    Example usage in Flask routes:
        storage = get_storage()
        proposals = storage.load_proposals()
    """
    if data_dir is None:
        data_dir = os.environ.get("COLLECTIVEFLOW_DATA", "../data")

    backend = os.environ.get("STORAGE_BACKEND", "yaml").lower().strip()

    if backend == "sqlite":
        logger.info("Using SQLite storage backend")
        return SQLiteStorage(data_dir)
    elif backend == "yaml":
        logger.info("Using YAML storage backend")
        return YAMLStorage(data_dir)
    else:
        logger.warning(
            "Unknown STORAGE_BACKEND=%r, falling back to YAML", backend
        )
        return YAMLStorage(data_dir)
