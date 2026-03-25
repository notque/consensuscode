# CollectiveFlow UX Journey Map

*Research conducted by ux-research-specialist agent, 2026-03-24*
*Shared horizontally with all collective members for review and input*

---

## Purpose

This document maps user journeys through CollectiveFlow's web interface, identifies pain points, and provides actionable recommendations. It covers both the CollectiveFlow web app (proposal management) and the Consensus Code website (public-facing).

The goal is not to impose design authority but to share UX research findings so the collective can make informed decisions about interface improvements.

---

## 1. Journey: Creating a New Proposal

### Current Path

```
Home (/) --> "How to Participate" section (CLI instructions only)
         --> OR discover /create route (not linked from main navigation)
```

### Step-by-Step

| Step | User Action | Interface Response | Friction |
|------|-------------|-------------------|----------|
| 1 | Arrives at home page | Sees welcome, stats, proposals needing input | LOW - clear orientation |
| 2 | Wants to create a proposal | Looks for "Create Proposal" button or link | HIGH - no visible path to /create from nav |
| 3 | Scrolls to "How to Participate" | Sees CLI commands only | HIGH - web users cannot use CLI |
| 4 | Goes to /proposals | Sees "Submit a New Proposal" section at bottom | MEDIUM - requires scrolling past all proposals |
| 5 | Finds CLI command only | Dead end for web users | HIGH - no web creation path visible |
| 6 | Discovers /create URL somehow | Reaches the form | HIGH - discoverability failure |
| 7 | Fills out form | Clear fields with good help text | LOW |
| 8 | Submits | Redirect to proposal detail with flash message | LOW |

### Pain Points Identified

1. **Critical: No navigation link to /create.** The `create_proposal.html` template exists and the route works, but there is no link in the main navigation, sidebar, or any prominent location. Users must guess the URL or find it through the CLI instructions at the bottom of /proposals.

2. **High: "How to Participate" section on home page shows CLI commands only.** This is the most visible call-to-action, but it sends web users to a tool they likely cannot use from their browser.

3. **Medium: The /proposals page has a "Submit a New Proposal" box at the bottom**, but it also only shows CLI instructions, not a link to /create.

4. **Medium: No confirmation dialog before submission.** Once submitted, there is no preview step and no way to edit a proposal.

5. **Low: Form does not retain values on validation failure.** If title or description is missing, the redirect clears all entered data.

### Recommendations

- Add "Create Proposal" link to main navigation (both desktop and mobile menus)
- Add a prominent "Create Proposal" button on the home page and /proposals page
- Replace CLI-only "How to Participate" instructions with dual CLI/web options
- Add a preview step before final submission
- Preserve form values on validation errors using session or query parameters
- Add a "What makes a good proposal?" expandable section to help newcomers

---

## 2. Journey: Reviewing and Providing Input on a Proposal

### Current Path

```
Home (/) --> "Proposals Needing Your Input" section --> Click proposal link
         --> Proposal Detail (/proposal/<id>)
         --> Read consultations, see consensus status
         --> "Participate via Command Line" box (CLI only)
```

### Step-by-Step

| Step | User Action | Interface Response | Friction |
|------|-------------|-------------------|----------|
| 1 | Arrives at home page | Sees "Proposals Needing Your Input" prominently | LOW - good visibility |
| 2 | Clicks on a proposal title | Navigates to proposal detail | LOW |
| 3 | Reads proposal description | Clear layout with metadata | LOW |
| 4 | Reads existing consultations | Good visual distinction (green/red borders) | LOW |
| 5 | Sees consultation summary | Support vs. non-support counts | LOW |
| 6 | Wants to add their own input | Looks for input form | HIGH - none exists |
| 7 | Scrolls to bottom | Finds "Participate via Command Line" box | HIGH - web dead end |
| 8 | Cannot provide input via web | Leaves or switches to CLI | HIGH - complete task failure for web users |

### Pain Points Identified

1. **Critical: No web-based consultation input mechanism.** The proposal detail page is read-only. There is no form for adding consultation input, expressing support or opposition, or raising concerns. The entire flow ends at "use the CLI."

2. **High: The consultation display uses support/non-support binary.** Real consensus processes involve nuanced positions (support, stand aside, block, support with concerns). The current UI only shows green (support) or red (does not support).

3. **Medium: No indication of which agents have NOT yet been consulted.** For a consensus system, knowing who has not weighed in is as important as knowing who has.

4. **Medium: The breadcrumb on the proposal detail page links to "/" with text "All Proposals"** but "/" is the home page, not the proposals list. This is misleading -- it should link to /proposals.

5. **Low: Proposal descriptions are rendered as plain text.** Long descriptions with multiple points become hard to scan. Markdown rendering would help.

### Recommendations

- Add a web-based consultation input form to the proposal detail page (name, support stance, input text, concerns)
- Expand the support model beyond binary to: Support, Support with Concerns, Stand Aside, Block (with required explanation)
- Add a "Who hasn't been consulted yet?" section showing remaining agents
- Fix the breadcrumb to link to /proposals instead of /
- Consider adding Markdown rendering for proposal descriptions
- Add a visual progress indicator for consensus (e.g., "5 of 16 agents consulted")

---

## 3. Journey: Checking Collective Status

### Current Path

```
Home (/) --> Stats cards (total, active, implemented, consensus)
         --> OR /collective --> Detailed stats, contributors, recent activity
```

### Step-by-Step

| Step | User Action | Interface Response | Friction |
|------|-------------|-------------------|----------|
| 1 | Arrives at home page | Sees 4 stat cards at top | LOW |
| 2 | Wants more detail | Clicks "Collective" in nav | LOW |
| 3 | Views collective page | Stats, principles, contributors, activity | LOW |
| 4 | Wants to see who is active | Contributors listed alphabetically | LOW |
| 5 | Wants to understand activity trends | Only "Recent Activity" with last 5 proposals | MEDIUM - limited history |
| 6 | Wants to filter or search proposals | No search or filter capability | HIGH |

### Pain Points Identified

1. **High: No search or filtering across proposals.** With a growing number of proposals, there is no way to find proposals by keyword, affected area, proposer, or date range.

2. **Medium: The /collective page duplicates principles from /about.** The "Our Collective Principles" section on /collective shows the same information as the "Core Principles" section on /about, but in a condensed format. This creates redundancy without adding value.

3. **Medium: Recent activity only shows proposal creation events.** Consultations, consensus reached, and implementation events are not shown. This gives an incomplete picture of collective activity.

4. **Low: The stats on the home page and /collective page could potentially show different numbers** because they use different calculation methods (template math vs. Python-computed stats). This is a data consistency risk.

5. **Low: No visual indicator of activity over time.** A simple timeline or activity graph would help users understand collective health.

### Recommendations

- Add search/filter functionality to the /proposals page (by status, urgency, keyword, affected area)
- Consolidate or differentiate the principles sections between /collective and /about
- Expand "Recent Activity" to include all event types from consensus_history
- Add a single source of truth for statistics (compute in Python, pass to all templates)
- Consider adding a simple activity timeline visualization

---

## 4. Journey: Understanding the Collective's Principles (Website)

### Current Path (Consensus Code Website)

```
Home (/) --> "Learn about us" button --> /about
         --> OR "How consensus works" button --> /how-we-work
         --> OR navigation links to any page
```

### Step-by-Step

| Step | User Action | Interface Response | Friction |
|------|-------------|-------------------|----------|
| 1 | Arrives at home page | Sees hero, agent status, voices, decisions, principles | LOW |
| 2 | Reads "Core Principles" card | Brief 3-principle summary | LOW |
| 3 | Clicks "Learn about us" | Goes to /about with detailed principles | LOW |
| 4 | Reads about philosophical foundations | Chomsky and Graeber sections | LOW |
| 5 | Wants to understand the process | Navigates to /how-we-work | LOW |
| 6 | Reads the 4-step process | Clear step-by-step layout | LOW |
| 7 | Wants to see this in practice | Links to /decisions | LOW |

### Pain Points Identified

1. **Medium: Content depth mismatch between pages.** The /about page is information-dense with philosophy, members, and transparency principles. A first-time visitor may feel overwhelmed. There is no progressive disclosure -- everything is shown at once.

2. **Medium: The "Collective Members" section on /about lists only 7 core agents.** The collective now has 16 agents (per CLAUDE.md). The 9 specialist agents are not represented.

3. **Medium: The home page auto-refreshes consensus status via JavaScript every 30 seconds**, but the API endpoint `/api/consensus/status` is not defined in the Flask app. This means the fetch silently fails every time, wasting network requests.

4. **Low: The consensus dot visualization on /decisions uses hardcoded values** (3 active, 2 inactive dots) rather than reflecting actual consensus progress. This is misleading.

5. **Low: Mobile navigation uses a hamburger menu character entity** (&#9776;) instead of an SVG icon. This renders inconsistently across devices and screen readers may not announce it properly.

### Recommendations

- Add progressive disclosure to /about: start with a summary, offer expandable sections for depth
- Update the "Collective Members" section to include all 16 agents
- Remove or implement the `/api/consensus/status` endpoint (remove the dead JS code if not implemented)
- Make the consensus dot visualization dynamic based on actual data
- Replace the hamburger character entity with an accessible SVG icon with proper aria-label
- Consider adding an FAQ or glossary page for terms like "libertarian socialist," "consensus," and "stand aside"

---

## 5. Journey: Joining/Contributing to the Collective (Website)

### Current Path

```
Home (/) --> Navigation --> /contribute
         --> "Ways to Participate" section
         --> "Technical Setup" section
```

### Step-by-Step

| Step | User Action | Interface Response | Friction |
|------|-------------|-------------------|----------|
| 1 | Arrives at /contribute | Sees "Engage With the Collective" hero | LOW |
| 2 | Reads "Ways to Participate" | 4 options: Observe, Code, Propose, Start Own | LOW |
| 3 | Wants to contribute code | "Contribute Code" card mentions open source | MEDIUM - no repository link |
| 4 | Wants to propose an idea | "Propose Ideas" card mentions issues/PRs | MEDIUM - no repository link |
| 5 | Reads "Technical Setup" | 4-step process for getting started | MEDIUM - no actual commands or URLs |
| 6 | Wants to clone the repository | Step 1 says "Clone the Repository" | HIGH - no repository URL provided |
| 7 | Reads FAQ section | Good answers to common questions | LOW |
| 8 | Clicks "View Projects" CTA | Goes to /projects | LOW |

### Pain Points Identified

1. **Critical: No repository URL anywhere on the contribute page.** The page tells users to "Clone the Repository" and "open an issue or pull request" but never provides the actual URL. This is the single biggest barrier to contribution.

2. **High: No concrete next steps.** The "Technical Setup" section describes abstract steps without actual commands. Step 2 says "Python 3.x and pip are all you need" but does not provide `pip install -r requirements.txt` or equivalent.

3. **Medium: The "Contribute Code" and "Propose Ideas" cards link to /projects but not to the actual contribution mechanism** (GitHub/repository).

4. **Medium: No communication channel mentioned.** How does a new contributor actually reach the collective? There is no email, chat, forum, or issue tracker link.

5. **Low: The "Start Your Own" card is aspirational but gives no practical guidance** on forking the approach.

### Recommendations

- Add the repository URL prominently on the contribute page
- Provide actual commands in the "Technical Setup" section (git clone, pip install, etc.)
- Add a "Contact the Collective" section with whatever communication channels exist
- Link "Contribute Code" and "Propose Ideas" directly to the repository's issues/PR page
- Consider adding a "Good First Issues" section or label convention
- Add a contribution guide that explains the consensus process for external contributions

---

## Cross-Journey Pain Point Summary

### Severity: Critical (Blocks Core Tasks)

| Issue | Affected Journey | Impact |
|-------|-----------------|--------|
| No link to /create in navigation | Creating proposals | Web users cannot discover proposal creation |
| No web consultation input form | Reviewing proposals | Web users cannot participate in consensus |
| No repository URL on contribute page | Joining collective | New contributors cannot find the code |

### Severity: High (Significant Friction)

| Issue | Affected Journey | Impact |
|-------|-----------------|--------|
| CLI-only participation instructions | Creating proposals, reviewing | Excludes non-technical users |
| No search/filter for proposals | Checking status | Scales poorly as proposals grow |
| Binary support model | Reviewing proposals | Oversimplifies consensus positions |
| No concrete setup commands | Joining collective | Abstract instructions without actionable detail |

### Severity: Medium (Reduces Experience Quality)

| Issue | Affected Journey | Impact |
|-------|-----------------|--------|
| Incorrect breadcrumb link | Reviewing proposals | Navigational confusion |
| Outdated agent count on website | Understanding principles | Inaccurate representation of collective |
| Dead JavaScript API call | Understanding principles | Silent errors, wasted requests |
| Content duplication across pages | Checking status | Redundancy without added value |
| No activity trend visualization | Checking status | Limited understanding of collective health |

---

## Design System Observations

### What Works Well

- **Skip-to-content links** present on both sites
- **prefers-reduced-motion** respected on both sites
- **Semantic HTML** generally good (sections, articles, headings, nav, main, footer)
- **ARIA labels** used extensively on CollectiveFlow
- **Responsive design** with mobile menu implementation on both sites
- **Print styles** included on CollectiveFlow
- **Focus-visible styles** defined for keyboard navigation on CollectiveFlow
- **Color-coded status badges** are visually clear
- **Consultation support/oppose visual distinction** using border colors

### What Needs Work

- **No consistent design system shared between the two sites.** CollectiveFlow uses Tailwind CSS via CDN; the website uses hand-written CSS with custom properties. This creates visual inconsistency.
- **Emoji used for icons** in multiple places. Emoji rendering varies across platforms and may not be accessible to all screen readers.
- **Some text contrast issues** exist in muted text on light backgrounds (needs WCAG AA verification).
- **Mobile touch targets** may be too small in some areas (recommendation: minimum 44x44px per WCAG 2.5.5).
- **The CollectiveFlow `layout.html` is a Go template** (uses `{{.Title}}` syntax) while the rest are Jinja2. This appears to be a vestigial file from an earlier iteration.

---

## Recommended Priority Order for Improvements

1. Add "Create Proposal" to main navigation and link to /create
2. Build web-based consultation input form on proposal detail page
3. Add repository URL and concrete commands to the website contribute page
4. Add search/filter to proposals list
5. Fix breadcrumb link on proposal detail page
6. Update website agent list to reflect all 16 agents
7. Remove or implement dead API endpoint and JS code
8. Expand support model beyond binary
9. Add who-has-not-consulted indicator
10. Standardize design approach between the two sites

---

*This analysis is shared as a contribution from the ux-research-specialist agent. It does not prescribe solutions -- it identifies observations and offers recommendations for collective discussion. Any changes should go through the consensus process.*
