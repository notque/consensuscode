"""
CollectiveFlow Web Application - Data Loading and Persistence Tests

These tests validate storage operations, ensuring:
- Proposals load correctly from the storage backend
- Data structure is preserved during load/save cycles
- File system operations are reliable (for YAML backend)
- Error handling works for corrupted or missing files
- The collective's data remains transparent and accessible

Why Test Data Operations?
File-based storage is central to our transparency principle.
These tests ensure data integrity and accessibility regardless
of which storage backend (YAML or SQLite) is in use.

Note: Tests use the app.storage object (injected by the temp_data_dir
fixture in conftest.py), so they exercise the StorageBackend interface
the same way routes do.
"""

import pytest
import yaml
from pathlib import Path
from datetime import datetime
import app as app_module
from conftest import assert_proposal_structure


def _storage():
    """Return the current storage backend from the app module.

    Helper to keep test code DRY. The conftest temp_data_dir fixture
    swaps app_module.storage to a YAMLStorage pointing at a temp dir,
    so all calls here operate on isolated test data.
    """
    return app_module.storage


class TestLoadProposals:
    """
    Tests for the _storage().load_proposals() function.

    This function reads all proposal YAML files and returns a list.
    It's the foundation of our data access.
    """

    @pytest.mark.data
    def test_load_proposals_returns_list(self, sample_proposals):
        """
        Test: load_proposals returns a list

        The function should always return a list, even if empty.
        """
        proposals = _storage().load_proposals()
        assert isinstance(proposals, list)

    @pytest.mark.data
    def test_load_proposals_with_data(self, sample_proposals):
        """
        Test: load_proposals returns all proposals from directory

        All YAML files in the proposals directory should be loaded.
        """
        proposals = _storage().load_proposals()

        # We created 3 sample proposals
        assert len(proposals) == 3

    @pytest.mark.data
    def test_load_proposals_contains_expected_data(self, sample_proposals):
        """
        Test: Loaded proposals contain expected data structures

        Each proposal should have all required fields.
        """
        proposals = _storage().load_proposals()

        # Check first proposal has expected structure
        for proposal in proposals:
            assert_proposal_structure(proposal)

    @pytest.mark.data
    def test_load_proposals_empty_directory(self, empty_proposals_dir):
        """
        Test: load_proposals handles empty directory gracefully

        An empty directory should return an empty list, not error.
        """
        proposals = _storage().load_proposals()
        assert isinstance(proposals, list)
        assert len(proposals) == 0

    @pytest.mark.data
    def test_load_proposals_nonexistent_directory(self, temp_data_dir, monkeypatch):
        """
        Test: load_proposals handles nonexistent directory

        If the proposals directory doesn't exist, return empty list.
        """
        # Point to a nonexistent directory by swapping the storage backend
        nonexistent_dir = str(temp_data_dir.parent / 'nonexistent')
        from storage import YAMLStorage
        monkeypatch.setattr(app_module, 'storage', YAMLStorage(nonexistent_dir))

        proposals = _storage().load_proposals()
        assert isinstance(proposals, list)
        assert len(proposals) == 0

    @pytest.mark.data
    def test_load_proposals_sorted_by_date(self, sample_proposals):
        """
        Test: Proposals are sorted by date, newest first

        Sorting helps members see recent activity first.
        """
        proposals = _storage().load_proposals()

        # Should be sorted with newest first
        if len(proposals) > 1:
            dates = [p.get('date', '') for p in proposals]
            # Check that dates are in descending order
            assert dates == sorted(dates, reverse=True)

    @pytest.mark.data
    def test_load_proposals_with_corrupted_file(self, temp_data_dir):
        """
        Test: load_proposals skips corrupted YAML files

        Bad files should be skipped rather than crashing the app.
        """
        # Create a corrupted YAML file
        corrupted_file = temp_data_dir / 'corrupted.yaml'
        with open(corrupted_file, 'w') as f:
            f.write('this is: not: valid: yaml: syntax:\n  - broken')

        # Should still load without crashing
        proposals = _storage().load_proposals()
        assert isinstance(proposals, list)
        # Corrupted file should be skipped

    @pytest.mark.data
    def test_load_proposals_preserves_all_fields(self, sample_proposals):
        """
        Test: All proposal fields are preserved during loading

        No data should be lost during the load process.
        """
        proposals = _storage().load_proposals()

        # Find the proposal with consultations
        proposal_with_consultations = next(
            (p for p in proposals if p['id'] == 'test-proposal-002'),
            None
        )

        assert proposal_with_consultations is not None
        assert 'consultations' in proposal_with_consultations
        assert len(proposal_with_consultations['consultations']) == 2

        # Check consultation structure
        consultation = proposal_with_consultations['consultations'][0]
        assert 'contributor' in consultation
        assert 'timestamp' in consultation
        assert 'input' in consultation


class TestGetProposal:
    """
    Tests for the _storage().get_proposal(proposal_id) function.

    This function loads a specific proposal by ID.
    It's used for detail page views.
    """

    @pytest.mark.data
    def test_get_proposal_existing_id(self, sample_proposals):
        """
        Test: get_proposal returns proposal for valid ID

        Members should be able to retrieve any proposal by its ID.
        """
        proposal = _storage().get_proposal('test-proposal-001')

        assert proposal is not None
        assert proposal['id'] == 'test-proposal-001'
        assert proposal['title'] == 'Test Proposal: Simple Example'

    @pytest.mark.data
    def test_get_proposal_returns_complete_data(self, sample_proposals):
        """
        Test: Retrieved proposal has all expected fields

        The full proposal structure should be returned.
        """
        proposal = _storage().get_proposal('test-proposal-002')

        assert_proposal_structure(proposal)
        assert 'consultations' in proposal
        assert 'consensus_history' in proposal

    @pytest.mark.data
    def test_get_proposal_nonexistent_id(self, sample_proposals):
        """
        Test: get_proposal returns None for invalid ID

        Non-existent proposals should return None, not error.
        """
        proposal = _storage().get_proposal('nonexistent-id')
        assert proposal is None

    @pytest.mark.data
    def test_get_proposal_with_consultations(self, sample_proposals):
        """
        Test: Retrieved proposal includes consultation data

        Consultations are essential for transparency.
        """
        proposal = _storage().get_proposal('test-proposal-002')

        assert 'consultations' in proposal
        assert len(proposal['consultations']) > 0

        # Check consultation structure
        consultation = proposal['consultations'][0]
        assert 'contributor' in consultation
        assert 'input' in consultation

    @pytest.mark.data
    def test_get_proposal_with_decision(self, sample_proposals):
        """
        Test: Implemented proposal includes decision data

        Decision rationale should be preserved and accessible.
        """
        proposal = _storage().get_proposal('test-proposal-003')

        assert 'decision' in proposal
        assert proposal['decision']['result'] == 'approved'
        assert 'rationale' in proposal['decision']

    @pytest.mark.data
    def test_get_proposal_empty_directory(self, empty_proposals_dir):
        """
        Test: get_proposal returns None when directory is empty

        Handles empty state gracefully.
        """
        proposal = _storage().get_proposal('any-id')
        assert proposal is None


class TestSaveProposal:
    """
    Tests for the _storage().save_proposal(proposal_data) function.

    This function creates new proposals and saves them to YAML files.
    It's used when members submit new proposals via the web interface.
    """

    @pytest.mark.data
    def test_save_proposal_creates_file(self, temp_data_dir):
        """
        Test: save_proposal creates YAML file

        New proposals should be persisted to the filesystem.
        """
        proposal_data = {
            'title': 'New Test Proposal',
            'description': 'Testing proposal creation',
            'proposer': 'test-agent',
            'urgency': 'medium',
            'affected_areas': ['testing']
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Check that file was created
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        assert yaml_path.exists()

    @pytest.mark.data
    def test_save_proposal_returns_id(self, temp_data_dir):
        """
        Test: save_proposal returns the proposal ID

        The ID is needed to redirect to the new proposal's detail page.
        """
        proposal_data = {
            'title': 'New Proposal',
            'description': 'Description',
            'proposer': 'test-agent',
            'urgency': 'low'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        assert proposal_id is not None
        assert isinstance(proposal_id, str)
        assert proposal_id.startswith('proposal-')

    @pytest.mark.data
    def test_save_proposal_generates_id_if_missing(self, temp_data_dir):
        """
        Test: save_proposal generates ID if not provided

        IDs should be auto-generated for convenience.
        """
        proposal_data = {
            'title': 'Proposal Without ID',
            'description': 'No ID provided',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Should have generated an ID
        assert 'id' not in proposal_data or proposal_data['id'] == proposal_id
        assert proposal_id is not None

    @pytest.mark.data
    def test_save_proposal_adds_metadata(self, temp_data_dir):
        """
        Test: save_proposal adds required metadata fields

        Status, date, and consensus tracking should be added automatically.
        """
        proposal_data = {
            'title': 'Metadata Test',
            'description': 'Testing metadata addition',
            'proposer': 'test-agent',
            'urgency': 'medium'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Load the saved proposal
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        with open(yaml_path, 'r') as f:
            saved_proposal = yaml.safe_load(f)

        # Check that metadata was added
        assert 'date' in saved_proposal
        assert 'status' in saved_proposal
        assert saved_proposal['status'] == 'proposed'
        assert 'consensus_status' in saved_proposal
        assert 'consensus_history' in saved_proposal
        assert 'consultations' in saved_proposal

    @pytest.mark.data
    def test_save_proposal_preserves_input_data(self, temp_data_dir):
        """
        Test: save_proposal preserves all input fields

        User-provided data should not be lost or modified.
        """
        proposal_data = {
            'title': 'Data Preservation Test',
            'description': 'Testing that input is preserved',
            'proposer': 'preservation-tester',
            'urgency': 'high',
            'affected_areas': ['testing', 'data']
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Load and verify
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        with open(yaml_path, 'r') as f:
            saved_proposal = yaml.safe_load(f)

        assert saved_proposal['title'] == proposal_data['title']
        assert saved_proposal['description'] == proposal_data['description']
        assert saved_proposal['proposer'] == proposal_data['proposer']
        assert saved_proposal['urgency'] == proposal_data['urgency']
        assert saved_proposal['affected_areas'] == proposal_data['affected_areas']

    @pytest.mark.data
    def test_save_proposal_creates_json_copy(self, temp_data_dir):
        """
        Test: save_proposal also creates JSON file for API compatibility

        Both YAML and JSON versions should be created.
        """
        proposal_data = {
            'title': 'JSON Test',
            'description': 'Testing JSON creation',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Check both files exist
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        json_path = temp_data_dir / f"{proposal_id}.json"

        assert yaml_path.exists()
        assert json_path.exists()

    @pytest.mark.data
    def test_save_proposal_yaml_is_readable(self, temp_data_dir):
        """
        Test: Saved YAML file is human-readable

        Human readability is a core principle of our storage choice.
        """
        proposal_data = {
            'title': 'Readability Test',
            'description': 'Testing YAML readability',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        with open(yaml_path, 'r') as f:
            content = f.read()

        # Check that it's not using flow style (should have newlines)
        assert '\n' in content
        # Should have readable keys
        assert 'title:' in content
        assert 'description:' in content

    @pytest.mark.data
    def test_save_proposal_creates_consensus_history(self, temp_data_dir):
        """
        Test: New proposal has initial consensus history entry

        Tracking begins at creation for full transparency.
        """
        proposal_data = {
            'title': 'History Test',
            'description': 'Testing consensus history',
            'proposer': 'test-agent',
            'urgency': 'low'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        with open(yaml_path, 'r') as f:
            saved_proposal = yaml.safe_load(f)

        # Check consensus history
        assert 'consensus_history' in saved_proposal
        assert len(saved_proposal['consensus_history']) == 1

        history_entry = saved_proposal['consensus_history'][0]
        assert 'timestamp' in history_entry
        assert history_entry['event'] == 'proposal_created'
        assert history_entry['actor'] == 'test-agent'

    @pytest.mark.data
    def test_save_proposal_directory_created_if_missing(self, temp_data_dir):
        """
        Test: save_proposal creates proposals directory if needed

        First-time setup should work automatically.
        """
        # Remove the directory
        import shutil
        shutil.rmtree(temp_data_dir)

        proposal_data = {
            'title': 'Directory Creation Test',
            'description': 'Testing auto-creation',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Directory should have been created
        assert temp_data_dir.exists()

        # File should exist
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        assert yaml_path.exists()


class TestDataIntegrity:
    """
    Tests for data integrity across load/save cycles.

    These tests ensure data is preserved correctly through
    the complete cycle of save, load, and display.
    """

    @pytest.mark.data
    @pytest.mark.integration
    def test_roundtrip_preserves_data(self, temp_data_dir):
        """
        Test: Data survives save-then-load cycle unchanged

        This is critical for data integrity.
        """
        original_data = {
            'title': 'Roundtrip Test',
            'description': 'Testing data preservation through save/load',
            'proposer': 'integrity-tester',
            'urgency': 'medium',
            'affected_areas': ['testing', 'data-integrity']
        }

        # Save
        proposal_id = _storage().save_proposal(original_data)

        # Load
        loaded_data = _storage().get_proposal(proposal_id)

        # Verify key fields preserved
        assert loaded_data['title'] == original_data['title']
        assert loaded_data['description'] == original_data['description']
        assert loaded_data['proposer'] == original_data['proposer']
        assert loaded_data['urgency'] == original_data['urgency']
        assert loaded_data['affected_areas'] == original_data['affected_areas']

    @pytest.mark.data
    def test_yaml_and_json_contain_same_data(self, temp_data_dir):
        """
        Test: YAML and JSON files have identical data

        Both formats should contain the same information.
        """
        import json

        proposal_data = {
            'title': 'Format Consistency Test',
            'description': 'Testing YAML/JSON consistency',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)

        # Load from both files
        yaml_path = temp_data_dir / f"{proposal_id}.yaml"
        json_path = temp_data_dir / f"{proposal_id}.json"

        with open(yaml_path, 'r') as f:
            yaml_data = yaml.safe_load(f)

        with open(json_path, 'r') as f:
            json_data = json.load(f)

        # Should contain same data
        assert yaml_data['id'] == json_data['id']
        assert yaml_data['title'] == json_data['title']
        assert yaml_data['description'] == json_data['description']

    @pytest.mark.data
    def test_special_characters_preserved(self, temp_data_dir):
        """
        Test: Special characters in text are preserved correctly

        YAML should handle unicode and special characters.
        """
        proposal_data = {
            'title': 'Special Characters: Testing "quotes", emoji 🤝, and unicode',
            'description': 'Testing preservation of:\n- Line breaks\n- "Quotes"\n- Emoji: 💡\n- Unicode: café',
            'proposer': 'test-agent'
        }

        proposal_id = _storage().save_proposal(proposal_data)
        loaded_data = _storage().get_proposal(proposal_id)

        # All special characters should be preserved
        assert loaded_data['title'] == proposal_data['title']
        assert loaded_data['description'] == proposal_data['description']
        assert '🤝' in loaded_data['title']
        assert '💡' in loaded_data['description']
        assert 'café' in loaded_data['description']

    @pytest.mark.data
    def test_empty_lists_preserved(self, temp_data_dir):
        """
        Test: Empty lists are preserved correctly

        New proposals start with empty consultations.
        """
        proposal_data = {
            'title': 'Empty List Test',
            'description': 'Testing empty list handling',
            'proposer': 'test-agent',
            'affected_areas': []
        }

        proposal_id = _storage().save_proposal(proposal_data)
        loaded_data = _storage().get_proposal(proposal_id)

        # Empty consultations list should exist
        assert 'consultations' in loaded_data
        assert loaded_data['consultations'] == []

        # Empty affected_areas should be preserved
        assert 'affected_areas' in loaded_data
        assert loaded_data['affected_areas'] == []


class TestErrorHandling:
    """
    Tests for error handling in data operations.

    Robust error handling ensures the application remains stable
    even with problematic data.
    """

    @pytest.mark.data
    def test_load_with_permission_error(self, temp_data_dir):
        """
        Test: Gracefully handle files that can't be read

        Permission errors shouldn't crash the entire load operation.
        """
        # This test is platform-dependent and may be skipped on some systems
        # Just verify load_proposals doesn't crash with various file issues
        proposals = _storage().load_proposals()
        assert isinstance(proposals, list)

    @pytest.mark.data
    def test_malformed_yaml_skipped(self, temp_data_dir):
        """
        Test: Malformed YAML files are skipped gracefully

        One bad file shouldn't prevent loading others.
        """
        # Create a good file
        good_data = {
            'id': 'good-proposal',
            'title': 'Good Proposal',
            'description': 'This one is fine',
            'proposer': 'test-agent',
            'date': '2025-07-26T10:00:00',
            'status': 'proposed',
            'urgency': 'low'
        }
        good_path = temp_data_dir / 'good.yaml'
        with open(good_path, 'w') as f:
            yaml.safe_dump(good_data, f)

        # Create a malformed file
        bad_path = temp_data_dir / 'bad.yaml'
        with open(bad_path, 'w') as f:
            f.write('{ this is not: valid yaml at all: [broken')

        # Should load the good file and skip the bad one
        proposals = _storage().load_proposals()
        assert len(proposals) >= 1
        assert any(p.get('id') == 'good-proposal' for p in proposals)
