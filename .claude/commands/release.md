# Release Command

Prepare and execute a release for the wt project.

## Instructions

1. Read the current version from the `VERSION` file
2. Find the last git tag matching pattern `v*` (e.g., v0.1.0, v1.0.0). If VERSION is 0.0.0 or no tags exist, analyze all commits on main
3. Get all commits since the last tag (or all commits if no tag):
   - Run `git log v<version>..HEAD --oneline` (or `git log --oneline` if no tags)
4. Analyze the commits and determine:
   - A concise summary of changes for release notes
   - Recommended bump type (major/minor/patch) based on:
     - **major**: Breaking changes, major new features, API changes
     - **minor**: New features, significant enhancements
     - **patch**: Bug fixes, small improvements, documentation
5. Present to the user:
   - Current version
   - List of commits being included
   - Proposed release notes (concise, 1-3 sentences)
   - Recommended bump type with reasoning
6. Ask the user to approve or modify the release notes and bump type
7. Upon approval, run: `make release <bump-type> "<release-notes>"`

## Important

- Keep release notes concise and user-focused (what changed, not how)
- If this is the first release (v0.0.0), recommend minor bump to get to v0.1.0
- Do not proceed with the release until the user explicitly approves
