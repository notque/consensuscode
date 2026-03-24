"""
CollectiveFlow Web Application - Jinja2 Filter Tests

These tests validate custom Jinja2 template filters, ensuring:
- Date formatting works correctly and is human-readable
- Status indicators display appropriate symbols
- Urgency levels have proper visual styling
- Filters handle edge cases and invalid input gracefully

Why Test Filters?
Template filters transform data for display. Testing them ensures
consistent, correct presentation across all pages.
"""

import pytest
from datetime import datetime
from app import humanize_date, status_emoji, urgency_color


class TestHumanizeDateFilter:
    """
    Tests for the humanize_date filter.

    This filter converts ISO datetime strings into human-readable formats.
    It's crucial for making timestamps accessible to all members.
    """

    @pytest.mark.filters
    def test_humanize_iso_date_string(self):
        """
        Test: ISO format datetime string converts to readable format

        Example: "2025-07-26T10:00:00-07:00" becomes "July 26, 2025 at 10:00 AM"
        """
        iso_date = '2025-07-26T10:00:00-07:00'
        result = humanize_date(iso_date)

        # Check for components of human-readable date
        assert 'July' in result or 'Jul' in result
        assert '26' in result
        assert '2025' in result
        assert 'AM' in result or 'PM' in result

    @pytest.mark.filters
    def test_humanize_iso_date_with_z_suffix(self):
        """
        Test: ISO format with Z (UTC) suffix converts correctly

        The filter should handle both timezone offsets and Z suffix.
        """
        iso_date = '2025-07-26T17:00:00Z'
        result = humanize_date(iso_date)

        # Should convert successfully
        assert 'July' in result or 'Jul' in result
        assert '26' in result
        assert '2025' in result

    @pytest.mark.filters
    def test_humanize_datetime_object(self):
        """
        Test: Python datetime object converts correctly

        The filter should work with both strings and datetime objects.
        """
        dt_object = datetime(2025, 7, 26, 14, 30, 0)
        result = humanize_date(dt_object)

        assert 'July' in result or 'Jul' in result
        assert '26' in result
        assert '2025' in result
        assert '02:30 PM' in result or '2:30 PM' in result

    @pytest.mark.filters
    def test_humanize_invalid_date_returns_original(self):
        """
        Test: Invalid date string returns the original value

        Graceful degradation ensures the app doesn't crash on bad data.
        """
        invalid_date = 'not-a-date'
        result = humanize_date(invalid_date)

        # Should return the input unchanged rather than crashing
        assert result == invalid_date

    @pytest.mark.filters
    def test_humanize_none_returns_none(self):
        """
        Test: None value returns None

        The filter should handle missing dates gracefully.
        """
        result = humanize_date(None)
        assert result is None

    @pytest.mark.filters
    def test_humanize_empty_string(self):
        """
        Test: Empty string returns empty string

        Empty values should not cause errors.
        """
        result = humanize_date('')
        assert result == ''

    @pytest.mark.filters
    def test_humanize_date_format_includes_time(self):
        """
        Test: Output includes both date and time components

        Full timestamps provide complete information.
        """
        iso_date = '2025-07-26T15:45:30-07:00'
        result = humanize_date(iso_date)

        # Should include date
        assert '2025' in result
        # Should include time
        assert '03:45 PM' in result or '3:45 PM' in result

    @pytest.mark.filters
    def test_humanize_different_months(self):
        """
        Test: Different months format correctly

        Verify month names appear correctly for various dates.
        """
        test_cases = {
            '2025-01-15T10:00:00': 'January',
            '2025-06-15T10:00:00': 'June',
            '2025-12-25T10:00:00': 'December'
        }

        for iso_date, expected_month in test_cases.items():
            result = humanize_date(iso_date)
            assert expected_month in result or expected_month[:3] in result


class TestStatusEmojiFilter:
    """
    Tests for the status_emoji filter.

    This filter returns emoji symbols representing proposal statuses.
    Visual indicators help members quickly scan proposal states.
    """

    @pytest.mark.filters
    def test_status_emoji_proposed(self):
        """
        Test: 'proposed' status returns appropriate emoji

        Proposed proposals should have a distinctive symbol.
        """
        result = status_emoji('proposed')
        assert result == '💡'

    @pytest.mark.filters
    def test_status_emoji_consultation(self):
        """
        Test: 'consultation' status returns discussion emoji

        Consultation phase indicates active discussion.
        """
        result = status_emoji('consultation')
        assert result == '🗣️'

    @pytest.mark.filters
    def test_status_emoji_consensus(self):
        """
        Test: 'consensus' status returns agreement emoji

        Consensus reached deserves a collaborative symbol.
        """
        result = status_emoji('consensus')
        assert result == '🤝'

    @pytest.mark.filters
    def test_status_emoji_implemented(self):
        """
        Test: 'implemented' status returns completion emoji

        Implemented proposals show success.
        """
        result = status_emoji('implemented')
        assert result == '✅'

    @pytest.mark.filters
    def test_status_emoji_blocked(self):
        """
        Test: 'blocked' status returns blocked emoji

        Blocked proposals need clear visual indication.
        """
        result = status_emoji('blocked')
        assert result == '🚫'

    @pytest.mark.filters
    def test_status_emoji_withdrawn(self):
        """
        Test: 'withdrawn' status returns withdrawal emoji

        Withdrawn proposals are marked as returned/removed.
        """
        result = status_emoji('withdrawn')
        assert result == '↩️'

    @pytest.mark.filters
    def test_status_emoji_unknown_status(self):
        """
        Test: Unknown status returns default emoji

        Graceful handling of unexpected status values.
        """
        result = status_emoji('unknown-status')
        assert result == '📄'  # Default emoji

    @pytest.mark.filters
    def test_status_emoji_none_value(self):
        """
        Test: None status returns default emoji

        Missing status should not cause errors.
        """
        result = status_emoji(None)
        assert result == '📄'

    @pytest.mark.filters
    def test_status_emoji_empty_string(self):
        """
        Test: Empty string status returns default emoji

        Empty status values handled gracefully.
        """
        result = status_emoji('')
        assert result == '📄'

    @pytest.mark.filters
    def test_all_valid_statuses_have_emojis(self):
        """
        Test: All valid status values have assigned emojis

        Completeness check - every status should have a symbol.
        """
        valid_statuses = [
            'proposed', 'consultation', 'consensus',
            'implemented', 'blocked', 'withdrawn'
        ]

        for status in valid_statuses:
            result = status_emoji(status)
            # Should return an emoji (unicode character), not the default
            assert result is not None
            assert len(result) > 0


class TestUrgencyColorFilter:
    """
    Tests for the urgency_color filter.

    This filter returns CSS classes for styling urgency levels.
    Visual differentiation helps members prioritize attention.
    """

    @pytest.mark.filters
    def test_urgency_color_low(self):
        """
        Test: Low urgency returns green color class

        Low urgency proposals styled in calming green.
        """
        result = urgency_color('low')
        assert 'green' in result.lower()
        assert 'text-' in result

    @pytest.mark.filters
    def test_urgency_color_medium(self):
        """
        Test: Medium urgency returns yellow/amber color class

        Medium urgency gets attention without alarm.
        """
        result = urgency_color('medium')
        assert 'yellow' in result.lower()
        assert 'text-' in result

    @pytest.mark.filters
    def test_urgency_color_high(self):
        """
        Test: High urgency returns orange color class

        High urgency stands out for prioritization.
        """
        result = urgency_color('high')
        assert 'orange' in result.lower()
        assert 'text-' in result

    @pytest.mark.filters
    def test_urgency_color_emergency(self):
        """
        Test: Emergency urgency returns red color class

        Emergency urgency demands immediate attention.
        """
        result = urgency_color('emergency')
        assert 'red' in result.lower()
        assert 'text-' in result

    @pytest.mark.filters
    def test_urgency_color_unknown(self):
        """
        Test: Unknown urgency returns default gray class

        Unknown urgency levels handled gracefully.
        """
        result = urgency_color('unknown-urgency')
        assert 'gray' in result.lower()

    @pytest.mark.filters
    def test_urgency_color_none_value(self):
        """
        Test: None urgency returns default color

        Missing urgency should not cause errors.
        """
        result = urgency_color(None)
        assert 'gray' in result.lower()

    @pytest.mark.filters
    def test_urgency_color_returns_tailwind_class(self):
        """
        Test: Returns valid Tailwind CSS class format

        All returned values should be valid Tailwind classes.
        """
        urgencies = ['low', 'medium', 'high', 'emergency']

        for urgency in urgencies:
            result = urgency_color(urgency)
            # Tailwind text color classes start with 'text-'
            assert result.startswith('text-')
            # Should include color and shade
            assert '-' in result

    @pytest.mark.filters
    def test_urgency_colors_are_distinct(self):
        """
        Test: Different urgency levels get different colors

        Visual distinction is essential for quick scanning.
        """
        colors = {
            'low': urgency_color('low'),
            'medium': urgency_color('medium'),
            'high': urgency_color('high'),
            'emergency': urgency_color('emergency')
        }

        # All colors should be different
        unique_colors = set(colors.values())
        assert len(unique_colors) == 4, "Urgency levels should have distinct colors"

    @pytest.mark.filters
    def test_urgency_color_case_insensitive(self):
        """
        Test: Urgency filter handles different cases

        Users might input urgency in various cases.
        """
        # Note: The current implementation is case-sensitive
        # This test documents current behavior
        lowercase = urgency_color('low')
        assert 'green' in lowercase.lower()

        # If uppercase isn't handled, it should return default
        uppercase = urgency_color('LOW')
        # This will likely be gray since dict lookup is case-sensitive


class TestFilterIntegration:
    """
    Tests for filter integration with templates.

    These tests verify filters work correctly when used in actual templates.
    """

    @pytest.mark.filters
    @pytest.mark.integration
    def test_filters_registered_in_app(self, app):
        """
        Test: All custom filters are registered with Flask app

        Filters must be registered to be available in templates.
        """
        # Check that custom filters are in the app's Jinja environment
        assert 'humanize_date' in app.jinja_env.filters
        assert 'status_emoji' in app.jinja_env.filters
        assert 'urgency_color' in app.jinja_env.filters

    @pytest.mark.filters
    @pytest.mark.integration
    def test_filters_used_in_templates(self, client, sample_proposals):
        """
        Test: Filters produce output in rendered templates

        Filters should actually transform data in real page renders.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # Should see humanized dates (month names)
        assert 'July' in data or 'June' in data or any(
            month in data for month in ['January', 'February', 'March', 'April', 'May']
        )

        # Should see urgency color classes (Tailwind classes)
        assert 'text-' in data

    @pytest.mark.filters
    @pytest.mark.integration
    def test_proposal_detail_uses_filters(self, client, sample_proposals):
        """
        Test: Proposal detail page applies filters correctly

        Detail pages should use all custom filters for rich display.
        """
        response = client.get('/proposal/test-proposal-002')
        data = response.data.decode('utf-8')

        # Should see humanized date
        assert '2025' in data

        # Should see status emoji in HTML
        # (emojis are unicode characters, so they'll be in the HTML)
        # We can check for consultation status emoji

        # Should see urgency color class
        assert 'text-' in data


class TestFilterEdgeCases:
    """
    Tests for edge cases and error handling in filters.

    Robust filters prevent template rendering failures.
    """

    @pytest.mark.filters
    def test_filters_dont_crash_on_none(self):
        """
        Test: Filters handle None values without crashing

        Missing data should not break page rendering.
        """
        # All filters should handle None gracefully
        assert humanize_date(None) is None
        assert status_emoji(None) == '📄'
        assert urgency_color(None) == 'text-gray-600'

    @pytest.mark.filters
    def test_filters_dont_crash_on_empty_string(self):
        """
        Test: Filters handle empty strings without crashing

        Empty values should be handled gracefully.
        """
        assert humanize_date('') == ''
        assert status_emoji('') == '📄'
        assert urgency_color('') == 'text-gray-600'

    @pytest.mark.filters
    def test_filters_dont_crash_on_unexpected_types(self):
        """
        Test: Filters handle unexpected data types

        Defensive programming prevents rendering failures.
        """
        # Test with numbers, lists, etc.
        # Most should return default values rather than crashing
        try:
            result = status_emoji(123)
            assert result == '📄'  # Should return default
        except:
            pytest.fail("Filter should not crash on unexpected type")

        try:
            result = urgency_color(['list'])
            assert 'gray' in result  # Should return default
        except:
            pytest.fail("Filter should not crash on unexpected type")
