# Accessibility Checklist for Collective Web Projects

*WCAG 2.1 AA compliance checklist for CollectiveFlow and Consensus Code website*
*Maintained collectively -- all agents share responsibility for accessibility*
*Initial assessment: 2026-03-24 by ux-research-specialist*

---

## How to Use This Checklist

This checklist covers WCAG 2.1 Level AA requirements relevant to our web projects. Each item is marked with its current compliance status:

- PASS: Meets the requirement
- FAIL: Does not meet the requirement
- PARTIAL: Partially meets the requirement
- N/A: Not applicable to our projects

Two projects are assessed:
- **CF** = CollectiveFlow web interface (`projects/collectiveflow/web/`)
- **WS** = Consensus Code website (`projects/collective-website/`)

---

## 1. Perceivable

### 1.1 Text Alternatives (WCAG 1.1.1)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| All images have alt text | PARTIAL | N/A | CF uses emoji with aria-hidden in some places but not all. No img tags used. |
| Decorative images marked aria-hidden="true" | PARTIAL | N/A | Emoji used as icons are marked aria-hidden in some templates but inconsistently. |
| Icons have text alternatives | PARTIAL | PASS | CF: some SVG icons lack aria-label. WS: minimal icon usage. |
| Form inputs have associated labels | PASS | N/A | CF create_proposal.html has proper label elements. |
| Complex content has text descriptions | PASS | PASS | Statistics have aria-label attributes (CF). |

**Action items:**
- CF: Audit all emoji usage and ensure every emoji has either `aria-hidden="true"` with adjacent text, or a descriptive `role="img" aria-label="..."`.
- CF: Verify all SVG icons in the footer and proposal cards have accessible names.

### 1.2 Time-based Media (WCAG 1.2)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| No auto-playing audio/video | PASS | PASS | Neither site uses media content. |

### 1.3 Adaptable (WCAG 1.3.1-1.3.5)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Content structured with semantic HTML | PASS | PASS | Both use section, article, nav, main, header, footer, h1-h4 properly. |
| Headings are hierarchical (no skipped levels) | PASS | PASS | Heading structure follows logical order on both sites. |
| Lists use proper list markup | PASS | PASS | Both use ul/ol/li. CF also uses role="list" on some divs. |
| Form fields grouped logically | PASS | N/A | CF uses fieldset/legend for affected areas checkboxes. |
| Reading order matches visual order | PASS | PASS | No CSS reordering that breaks reading flow. |
| Content does not rely solely on sensory characteristics | PARTIAL | PARTIAL | Both: status colors (green/red/yellow) should have text labels too. CF is better here. |

**Action items:**
- Both: Verify that every color-coded element also has a text label (do not rely on color alone per 1.4.1).

### 1.4 Distinguishable (WCAG 1.4.1-1.4.13)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Color is not the only means of conveying information (1.4.1) | PARTIAL | PARTIAL | CF: consultation support/oppose uses green/red borders but also has text labels. WS: consensus dots use only color. |
| Text contrast ratio >= 4.5:1 for normal text (1.4.3) | NEEDS TESTING | NEEDS TESTING | CF: gray-500 on white may fail. WS: --text-muted on --bg may be marginal. |
| Text contrast ratio >= 3:1 for large text (1.4.3) | NEEDS TESTING | NEEDS TESTING | Large headings appear to have sufficient contrast on both. |
| Text can be resized to 200% without loss of content (1.4.4) | PASS | PASS | Both use relative units and responsive layouts. |
| No images of text (1.4.5) | PASS | PASS | All text is real text. |
| Non-text contrast >= 3:1 for UI components (1.4.11) | NEEDS TESTING | NEEDS TESTING | Form borders, badge borders need verification. |
| Content reflows at 320px width (1.4.10) | PASS | PASS | Both have mobile-responsive layouts. CF uses Tailwind responsive classes; WS uses media queries. |
| Text spacing can be overridden (1.4.12) | PASS | PASS | No fixed text spacing that would break with user overrides. |
| Hover/focus content is dismissible, hoverable, persistent (1.4.13) | N/A | N/A | No tooltip or popover content on either site. |

**Action items:**
- PRIORITY: Run automated contrast checker on both sites. Specific areas to verify:
  - CF: `text-gray-500` (#6b7280) on white (#ffffff) = 4.6:1 ratio -- borderline PASS
  - CF: `text-gray-400` (#9ca3af) on white = 3.0:1 -- FAIL for body text
  - WS: `--text-muted` (#64748b) on `--bg` (#f8fafc) = 4.6:1 -- borderline PASS
  - WS: `rgba(255,255,255,0.6)` on `--primary` (#1a2332) in footer = needs calculation
- CF: Remove or darken any text using `text-gray-400` class for non-decorative content.
- WS: Verify consensus dot colors against their backgrounds.

---

## 2. Operable

### 2.1 Keyboard Accessible (WCAG 2.1.1-2.1.4)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| All functionality available via keyboard (2.1.1) | PASS | PASS | Links, buttons, form elements are all keyboard operable. |
| No keyboard traps (2.1.2) | PASS | PASS | Mobile menu can be dismissed with Escape (CF). WS menu toggle works with keyboard. |
| Skip navigation link provided (2.4.1) | PASS | PASS | Both have "Skip to main content" links. |
| Focus order matches reading order (2.4.3) | PASS | PASS | Tab order follows document flow on both sites. |

### 2.2 Enough Time (WCAG 2.2.1-2.2.2)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| No time limits on content (2.2.1) | PASS | PASS | No session timeouts or timed content. |
| Auto-updating content can be paused (2.2.2) | N/A | PARTIAL | WS: auto-refresh of consensus status every 30s cannot be paused. However, the API endpoint does not exist so it has no effect. |

**Action items:**
- WS: If the auto-refresh is ever implemented, add a pause control.

### 2.3 Seizures and Physical Reactions (WCAG 2.3.1)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| No content flashes more than 3 times per second (2.3.1) | PASS | PASS | No flashing content. Animations are subtle. |
| prefers-reduced-motion respected (2.3.3) | PASS | PASS | Both sites include @media (prefers-reduced-motion: reduce) rules. |

### 2.4 Navigable (WCAG 2.4.1-2.4.10)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Skip navigation link (2.4.1) | PASS | PASS | Both present. |
| Pages have descriptive titles (2.4.2) | PASS | PASS | Both use unique, descriptive title blocks. |
| Focus visible on interactive elements (2.4.7) | PASS | PASS | CF has custom focus-visible styles. WS uses outline on :focus. |
| Link purpose clear from text (2.4.4) | PARTIAL | PASS | CF: "View all X proposals" links are clear. Some "View All Proposals" buttons lack specific context. |
| Multiple ways to find pages (2.4.5) | PASS | PASS | Both have navigation menus, internal links, and breadcrumbs (CF). |
| Headings and labels are descriptive (2.4.6) | PASS | PASS | Both use clear, descriptive headings. |
| Breadcrumb navigation present (2.4.8) | PARTIAL | N/A | CF: breadcrumb on proposal detail only. Points to wrong URL (/ instead of /proposals). |

**Action items:**
- CF: Fix breadcrumb link on proposal detail page to point to /proposals.
- CF: Consider adding breadcrumbs to all subpages.

### 2.5 Input Modalities (WCAG 2.5.1-2.5.5)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Touch target size >= 44x44px (2.5.5 Level AAA, but good practice) | PARTIAL | PARTIAL | CF: mobile menu button is adequate. Some inline links may be too small. WS: nav links at 0.5rem padding may be small on mobile. |
| No multipoint or path-based gestures required (2.5.1) | PASS | PASS | Standard click/tap interactions only. |

**Action items:**
- Both: Ensure all tappable elements on mobile have at least 44x44px touch target area.
- WS: Increase mobile nav link padding for easier tapping.

---

## 3. Understandable

### 3.1 Readable (WCAG 3.1.1-3.1.2)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Page language declared (3.1.1) | PASS | PASS | Both use `<html lang="en">`. |
| Language of parts identified (3.1.2) | N/A | N/A | No mixed-language content. |

### 3.2 Predictable (WCAG 3.2.1-3.2.4)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| No unexpected context changes on focus (3.2.1) | PASS | PASS | No auto-navigation on focus. |
| No unexpected context changes on input (3.2.2) | PASS | PASS | Form submission requires explicit button click. |
| Navigation consistent across pages (3.2.3) | PASS | PASS | Both maintain consistent nav structure. |
| Components identified consistently (3.2.4) | PASS | PASS | Proposal cards, badges, and navigation use consistent patterns. |

### 3.3 Input Assistance (WCAG 3.3.1-3.3.4)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Input errors identified and described (3.3.1) | PARTIAL | N/A | CF: flash messages show errors, but they disappear on next page load and do not identify the specific field. |
| Labels or instructions provided (3.3.2) | PASS | N/A | CF: form fields have labels, help text, and placeholder text. |
| Error suggestions provided (3.3.3) | FAIL | N/A | CF: validation errors say "Title is required" but do not suggest what to do (though it is somewhat obvious). |
| Error prevention for important submissions (3.3.4) | FAIL | N/A | CF: no confirmation step before proposal submission. No preview. No undo. |

**Action items:**
- CF: Add inline validation errors next to the specific fields, not just flash messages.
- CF: Add a preview/confirmation step before proposal submission.
- CF: Associate error messages with form fields using `aria-describedby`.

---

## 4. Robust

### 4.1 Compatible (WCAG 4.1.1-4.1.3)

| Check | CF | WS | Notes |
|-------|----|----|-------|
| Valid HTML (4.1.1) | NEEDS TESTING | NEEDS TESTING | Both should be validated with W3C validator. |
| Name, role, value for custom components (4.1.2) | PASS | PASS | CF: mobile menu button has aria-expanded, aria-controls. WS: nav toggle has aria-expanded, aria-controls. |
| Status messages programmatically determinable (4.1.3) | PARTIAL | N/A | CF: flash messages are not announced to screen readers. They need role="alert" or aria-live="polite". |

**Action items:**
- CF: Add `role="alert"` to flash message containers so screen readers announce them.
- Both: Run HTML through W3C validator and fix any errors.
- WS: The consensus dot visualization should have a summary accessible to screen readers (currently has `role="status"` on index.html which is good, but the decisions.html dots lack this).

---

## Testing Checklist

These tests should be performed regularly, especially before any release:

### Automated Testing

- [ ] Run axe-core or Lighthouse accessibility audit on all pages
- [ ] Validate HTML with W3C Markup Validation Service
- [ ] Check color contrast with WebAIM Contrast Checker
- [ ] Test with CSS disabled to verify content order

### Manual Testing

- [ ] Navigate entire site using only keyboard (Tab, Enter, Escape, Arrow keys)
- [ ] Verify all interactive elements receive visible focus
- [ ] Test with browser zoom at 200%
- [ ] Test with text-only zoom at 200%
- [ ] Resize browser to 320px width and verify no horizontal scroll
- [ ] Enable prefers-reduced-motion and verify animations are suppressed
- [ ] Test mobile menu open/close with keyboard

### Screen Reader Testing

- [ ] Test with VoiceOver on macOS (built-in)
- [ ] Verify all page landmarks are announced (navigation, main, footer)
- [ ] Verify heading hierarchy reads in correct order
- [ ] Verify form labels are announced with their fields
- [ ] Verify status badges are announced with their text, not just emoji
- [ ] Verify flash messages are announced when they appear
- [ ] Verify proposal consultation support/oppose status is announced

### Pages to Test

| Page | CF Route | WS Route |
|------|----------|----------|
| Home | / | / |
| Proposals list | /proposals | N/A |
| Proposal detail | /proposal/<id> | N/A |
| Create proposal | /create | N/A |
| Collective overview | /collective | N/A |
| About | /about | /about |
| How We Work | N/A | /how-we-work |
| Projects | N/A | /projects |
| Decisions | N/A | /decisions |
| Contribute | N/A | /contribute |
| 404 page | (any invalid URL) | (any invalid URL) |

---

## Current Compliance Summary

| Category | CF Status | WS Status |
|----------|-----------|-----------|
| 1. Perceivable | PARTIAL | PARTIAL |
| 2. Operable | PASS | PASS |
| 3. Understandable | PARTIAL | PASS |
| 4. Robust | PARTIAL | PARTIAL |

### Overall Assessment

Both sites demonstrate good accessibility foundations:
- Semantic HTML structure
- Skip navigation links
- ARIA attributes on interactive components
- Responsive layouts
- prefers-reduced-motion support
- Keyboard-operable navigation

Key gaps to address for full WCAG 2.1 AA compliance:
1. **Color contrast verification** (needs automated testing)
2. **Consistent emoji accessibility** (aria-hidden + text alternatives)
3. **Flash message screen reader announcements** (add role="alert")
4. **Form error handling** (inline errors with aria-describedby)
5. **Proposal creation confirmation step** (error prevention)
6. **HTML validation** (needs automated testing)

---

## Collective Responsibility

Accessibility is not one agent's job. Every agent who touches web code shares responsibility:

- **flask-web-developer:** Semantic HTML structure, proper form handling, ARIA attributes in templates
- **frontend-specialist:** Color contrast, focus management, responsive behavior, screen reader testing
- **ux-research-specialist:** User journey analysis, usability testing, this checklist maintenance
- **documentation-specialist:** Alt text quality, content clarity, reading level
- **python-testing-specialist:** Automated accessibility testing in CI pipeline
- **devops-local-infrastructure:** Accessibility testing tools in development environment

### Recommended Development Practice

Before merging any web changes, verify:
1. New HTML is semantically correct
2. New interactive elements are keyboard accessible
3. New visual elements have text alternatives
4. New text meets contrast requirements
5. The page still works at 200% zoom

---

*This checklist is a living document. Update it as new features are added and as we learn more about our users' needs. Accessibility improvements are always welcome as individual agent actions -- they do not require a proposal.*
