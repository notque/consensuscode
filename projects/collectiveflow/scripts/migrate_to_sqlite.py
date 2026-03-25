#!/usr/bin/env python3
"""
Migrate CollectiveFlow data from YAML files to SQLite.

This script reads all YAML proposal files from the data/proposals/ directory
and imports them into a SQLite database using the schema defined in schema.sql.

The migration is reversible: use --export to dump SQLite back to YAML files.

Usage:
    # Import YAML files into SQLite
    python3 migrate_to_sqlite.py --import --data-dir ./data/proposals --db ./data/collectiveflow.db

    # Export SQLite back to YAML files (reversibility)
    python3 migrate_to_sqlite.py --export --data-dir ./data/proposals-export --db ./data/collectiveflow.db

    # Dry run (validate without writing)
    python3 migrate_to_sqlite.py --import --data-dir ./data/proposals --db ./data/collectiveflow.db --dry-run

    # Verify migration integrity
    python3 migrate_to_sqlite.py --verify --data-dir ./data/proposals --db ./data/collectiveflow.db

Dependencies: PyYAML (pip install pyyaml)
"""

import argparse
import json
import os
import sqlite3
import sys
from datetime import datetime
from pathlib import Path

try:
    import yaml
except ImportError:
    print("Error: PyYAML is required. Install with: pip install pyyaml", file=sys.stderr)
    sys.exit(1)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def normalize_timestamp(ts):
    """
    Normalize a timestamp to ISO 8601 format.
    The YAML files contain Go-formatted timestamps like:
        2025-07-26T10:03:44.302413-07:00
    SQLite expects TEXT in a sortable format. We keep the original string
    since ISO 8601 with timezone is already sortable within a timezone.
    """
    if isinstance(ts, datetime):
        return ts.isoformat()
    if isinstance(ts, str):
        return ts
    return str(ts)


def read_yaml_proposal(filepath):
    """Read a single YAML proposal file and return the parsed dict."""
    with open(filepath, "r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    return data


def get_schema_sql():
    """Read the schema.sql file from the same directory as this script."""
    schema_path = Path(__file__).parent / "schema.sql"
    if not schema_path.exists():
        print(f"Error: schema.sql not found at {schema_path}", file=sys.stderr)
        sys.exit(1)
    return schema_path.read_text(encoding="utf-8")


# ---------------------------------------------------------------------------
# Import: YAML -> SQLite
# ---------------------------------------------------------------------------

def import_proposal(conn, data, actor="migration"):
    """
    Import a single proposal dict (parsed from YAML) into SQLite.

    This function inserts into all relevant tables:
    - proposals
    - consultations
    - consensus_events
    - decisions
    - audit_log
    """
    proposal_id = data.get("id", "")
    if not proposal_id:
        raise ValueError("Proposal has no 'id' field")

    # -- proposals table --
    affected_areas = data.get("affected_areas") or []
    if isinstance(affected_areas, str):
        affected_areas = [affected_areas]

    conn.execute(
        """
        INSERT OR REPLACE INTO proposals
            (id, title, description, proposer, created_at, status, urgency,
             consensus_status, affected_areas)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
        (
            proposal_id,
            data.get("title", ""),
            data.get("description", ""),
            data.get("proposer", ""),
            normalize_timestamp(data.get("date", "")),
            data.get("status", "proposed"),
            data.get("urgency", "medium"),
            data.get("consensus_status", ""),
            json.dumps(affected_areas),
        ),
    )

    # -- audit_log for proposal insert --
    conn.execute(
        """
        INSERT INTO audit_log (table_name, row_id, operation, new_values, actor)
        VALUES (?, ?, 'INSERT', ?, ?)
        """,
        (
            "proposals",
            proposal_id,
            json.dumps({"id": proposal_id, "title": data.get("title", "")}),
            actor,
        ),
    )

    # -- consensus_events table --
    for event in data.get("consensus_history") or []:
        conn.execute(
            """
            INSERT INTO consensus_events
                (proposal_id, event_type, actor, details, created_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            (
                proposal_id,
                event.get("event", ""),
                event.get("actor", ""),
                event.get("details", ""),
                normalize_timestamp(event.get("timestamp", "")),
            ),
        )

    # -- consultations table --
    for consultation in data.get("consultations") or []:
        concerns = consultation.get("concerns") or []
        if isinstance(concerns, str):
            concerns = [concerns]

        conn.execute(
            """
            INSERT INTO consultations
                (proposal_id, contributor, created_at, input, support, concerns)
            VALUES (?, ?, ?, ?, ?, ?)
            """,
            (
                proposal_id,
                consultation.get("contributor", ""),
                normalize_timestamp(consultation.get("timestamp", "")),
                consultation.get("input", ""),
                1 if consultation.get("support", False) else 0,
                json.dumps(concerns),
            ),
        )

    # -- decisions table --
    decision = data.get("decision")
    if decision and isinstance(decision, dict):
        result = decision.get("result", "")
        if result:
            conn.execute(
                """
                INSERT OR REPLACE INTO decisions
                    (proposal_id, result, rationale, created_at)
                VALUES (?, ?, ?, ?)
                """,
                (
                    proposal_id,
                    result,
                    decision.get("rationale", ""),
                    normalize_timestamp(decision.get("timestamp", "")),
                ),
            )


def run_import(data_dir, db_path, dry_run=False):
    """Import all YAML proposals from data_dir into the SQLite database at db_path."""
    data_path = Path(data_dir)
    if not data_path.exists():
        print(f"Error: Data directory not found: {data_dir}", file=sys.stderr)
        return False

    yaml_files = sorted(data_path.glob("*.yaml"))
    if not yaml_files:
        print(f"No YAML files found in {data_dir}")
        return True

    print(f"Found {len(yaml_files)} YAML proposal files in {data_dir}")

    # Create database and apply schema
    if dry_run:
        # Use in-memory database for dry run
        conn = sqlite3.connect(":memory:")
    else:
        # Ensure parent directory exists
        db_file = Path(db_path)
        db_file.parent.mkdir(parents=True, exist_ok=True)
        conn = sqlite3.connect(db_path)

    conn.executescript(get_schema_sql())

    imported = 0
    errors = 0

    for yaml_file in yaml_files:
        try:
            data = read_yaml_proposal(yaml_file)
            if not data:
                print(f"  SKIP (empty): {yaml_file.name}")
                continue

            import_proposal(conn, data, actor="yaml-migration")
            imported += 1
            print(f"  OK: {yaml_file.name} -> {data.get('id', '?')}")

        except Exception as e:
            errors += 1
            print(f"  ERROR: {yaml_file.name}: {e}", file=sys.stderr)

    if not dry_run:
        conn.commit()

    # Print summary
    print(f"\nMigration {'(DRY RUN) ' if dry_run else ''}complete:")
    print(f"  Imported: {imported}")
    print(f"  Errors:   {errors}")
    print(f"  Total:    {len(yaml_files)}")

    if not dry_run and imported > 0:
        # Verify counts
        row = conn.execute("SELECT COUNT(*) FROM proposals").fetchone()
        print(f"\n  Database proposals:   {row[0]}")
        row = conn.execute("SELECT COUNT(*) FROM consultations").fetchone()
        print(f"  Database consultations: {row[0]}")
        row = conn.execute("SELECT COUNT(*) FROM consensus_events").fetchone()
        print(f"  Database events:      {row[0]}")
        row = conn.execute("SELECT COUNT(*) FROM decisions").fetchone()
        print(f"  Database decisions:   {row[0]}")
        print(f"\n  Database file: {db_path}")

    conn.close()
    return errors == 0


# ---------------------------------------------------------------------------
# Export: SQLite -> YAML (reversibility)
# ---------------------------------------------------------------------------

def export_proposal(conn, proposal_id):
    """
    Export a single proposal from SQLite back to a dict matching the YAML format.
    """
    # Fetch proposal
    row = conn.execute(
        "SELECT * FROM proposals WHERE id = ?", (proposal_id,)
    ).fetchone()
    if not row:
        return None

    cols = [desc[0] for desc in conn.execute("SELECT * FROM proposals LIMIT 0").description]
    proposal = dict(zip(cols, row))

    # Build the YAML-compatible dict in the original field order
    result = {
        "id": proposal["id"],
        "title": proposal["title"],
        "description": proposal["description"],
        "proposer": proposal["proposer"],
        "date": proposal["created_at"],
        "status": proposal["status"],
        "urgency": proposal["urgency"],
        "affected_areas": json.loads(proposal["affected_areas"]),
        "consensus_status": proposal["consensus_status"],
    }

    # Consensus history
    events = conn.execute(
        """
        SELECT event_type, actor, details, created_at
        FROM consensus_events
        WHERE proposal_id = ?
        ORDER BY created_at ASC, id ASC
        """,
        (proposal_id,),
    ).fetchall()

    if events:
        result["consensus_history"] = [
            {
                "timestamp": e[3],
                "event": e[0],
                "actor": e[1],
                "details": e[2],
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
        (proposal_id,),
    ).fetchall()

    if consults:
        consultations_list = []
        for c in consults:
            entry = {
                "contributor": c[0],
                "timestamp": c[1],
                "input": c[2],
                "support": bool(c[3]),
            }
            concerns = json.loads(c[4])
            if concerns:
                entry["concerns"] = concerns
            consultations_list.append(entry)
        result["consultations"] = consultations_list

    # Decision
    decision = conn.execute(
        """
        SELECT result, rationale, created_at
        FROM decisions
        WHERE proposal_id = ?
        """,
        (proposal_id,),
    ).fetchone()

    if decision:
        result["decision"] = {
            "result": decision[0],
            "timestamp": decision[2],
            "rationale": decision[1],
        }

    return result


def run_export(data_dir, db_path):
    """Export all proposals from SQLite back to YAML files."""
    if not Path(db_path).exists():
        print(f"Error: Database not found: {db_path}", file=sys.stderr)
        return False

    output_path = Path(data_dir)
    output_path.mkdir(parents=True, exist_ok=True)

    conn = sqlite3.connect(db_path)
    conn.execute("PRAGMA foreign_keys = ON")

    proposal_ids = [
        row[0]
        for row in conn.execute(
            "SELECT id FROM proposals ORDER BY created_at ASC"
        ).fetchall()
    ]

    if not proposal_ids:
        print("No proposals found in database.")
        conn.close()
        return True

    print(f"Exporting {len(proposal_ids)} proposals to {data_dir}")

    exported = 0
    errors = 0

    for pid in proposal_ids:
        try:
            data = export_proposal(conn, pid)
            if not data:
                continue

            out_file = output_path / f"{pid}.yaml"
            with open(out_file, "w", encoding="utf-8") as f:
                yaml.dump(
                    data,
                    f,
                    default_flow_style=False,
                    allow_unicode=True,
                    sort_keys=False,
                    width=120,
                )
            exported += 1
            print(f"  OK: {pid} -> {out_file.name}")

        except Exception as e:
            errors += 1
            print(f"  ERROR: {pid}: {e}", file=sys.stderr)

    conn.close()

    print(f"\nExport complete:")
    print(f"  Exported: {exported}")
    print(f"  Errors:   {errors}")
    print(f"  Output:   {data_dir}")

    return errors == 0


# ---------------------------------------------------------------------------
# Verify: Compare YAML files against SQLite data
# ---------------------------------------------------------------------------

def run_verify(data_dir, db_path):
    """
    Verify that the SQLite database matches the YAML source files.
    Compares proposal IDs, consultation counts, and decision results.
    """
    data_path = Path(data_dir)
    if not data_path.exists():
        print(f"Error: Data directory not found: {data_dir}", file=sys.stderr)
        return False
    if not Path(db_path).exists():
        print(f"Error: Database not found: {db_path}", file=sys.stderr)
        return False

    yaml_files = sorted(data_path.glob("*.yaml"))
    conn = sqlite3.connect(db_path)
    conn.execute("PRAGMA foreign_keys = ON")

    mismatches = 0
    checked = 0

    for yaml_file in yaml_files:
        data = read_yaml_proposal(yaml_file)
        if not data:
            continue

        pid = data.get("id", "")
        checked += 1

        # Check proposal exists in DB
        row = conn.execute(
            "SELECT id, title, status, urgency FROM proposals WHERE id = ?",
            (pid,),
        ).fetchone()
        if not row:
            print(f"  MISSING in DB: {pid}")
            mismatches += 1
            continue

        # Check title
        if row[1] != data.get("title", ""):
            print(f"  MISMATCH title: {pid}")
            mismatches += 1

        # Check status
        if row[2] != data.get("status", ""):
            print(f"  MISMATCH status: {pid} (yaml={data.get('status')}, db={row[2]})")
            mismatches += 1

        # Check consultation count
        yaml_count = len(data.get("consultations") or [])
        db_count = conn.execute(
            "SELECT COUNT(*) FROM consultations WHERE proposal_id = ?",
            (pid,),
        ).fetchone()[0]
        if yaml_count != db_count:
            print(f"  MISMATCH consultations: {pid} (yaml={yaml_count}, db={db_count})")
            mismatches += 1

        # Check consensus event count
        yaml_events = len(data.get("consensus_history") or [])
        db_events = conn.execute(
            "SELECT COUNT(*) FROM consensus_events WHERE proposal_id = ?",
            (pid,),
        ).fetchone()[0]
        if yaml_events != db_events:
            print(f"  MISMATCH events: {pid} (yaml={yaml_events}, db={db_events})")
            mismatches += 1

        # Check decision
        yaml_decision = data.get("decision")
        db_decision = conn.execute(
            "SELECT result FROM decisions WHERE proposal_id = ?",
            (pid,),
        ).fetchone()
        if yaml_decision and not db_decision:
            print(f"  MISSING decision in DB: {pid}")
            mismatches += 1
        elif not yaml_decision and db_decision:
            print(f"  EXTRA decision in DB: {pid}")
            mismatches += 1
        elif yaml_decision and db_decision:
            if yaml_decision.get("result") != db_decision[0]:
                print(f"  MISMATCH decision: {pid}")
                mismatches += 1

    # Check for extra proposals in DB
    db_ids = {
        row[0]
        for row in conn.execute("SELECT id FROM proposals").fetchall()
    }
    yaml_ids = set()
    for yf in yaml_files:
        d = read_yaml_proposal(yf)
        if d:
            yaml_ids.add(d.get("id", ""))
    extra = db_ids - yaml_ids
    if extra:
        print(f"  EXTRA proposals in DB not in YAML: {extra}")
        mismatches += len(extra)

    conn.close()

    print(f"\nVerification complete:")
    print(f"  Checked:    {checked}")
    print(f"  Mismatches: {mismatches}")

    if mismatches == 0:
        print("  Result:     ALL OK")
    else:
        print("  Result:     MISMATCHES FOUND")

    return mismatches == 0


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------

def main():
    parser = argparse.ArgumentParser(
        description="Migrate CollectiveFlow data between YAML and SQLite.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Import YAML to SQLite
  python3 migrate_to_sqlite.py --import --data-dir ./data/proposals --db ./data/collectiveflow.db

  # Export SQLite back to YAML
  python3 migrate_to_sqlite.py --export --data-dir ./data/proposals-export --db ./data/collectiveflow.db

  # Dry run (validate without writing)
  python3 migrate_to_sqlite.py --import --data-dir ./data/proposals --db ./data/collectiveflow.db --dry-run

  # Verify migration integrity
  python3 migrate_to_sqlite.py --verify --data-dir ./data/proposals --db ./data/collectiveflow.db
        """,
    )

    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--import", dest="do_import", action="store_true",
        help="Import YAML files into SQLite database",
    )
    group.add_argument(
        "--export", dest="do_export", action="store_true",
        help="Export SQLite database back to YAML files",
    )
    group.add_argument(
        "--verify", dest="do_verify", action="store_true",
        help="Verify SQLite data matches YAML source files",
    )

    parser.add_argument(
        "--data-dir", required=True,
        help="Path to YAML proposals directory (source for import, destination for export)",
    )
    parser.add_argument(
        "--db", required=True,
        help="Path to SQLite database file",
    )
    parser.add_argument(
        "--dry-run", action="store_true",
        help="Validate import without writing to disk (import only)",
    )

    args = parser.parse_args()

    if args.do_import:
        ok = run_import(args.data_dir, args.db, dry_run=args.dry_run)
    elif args.do_export:
        ok = run_export(args.data_dir, args.db)
    elif args.do_verify:
        ok = run_verify(args.data_dir, args.db)
    else:
        parser.print_help()
        ok = False

    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
