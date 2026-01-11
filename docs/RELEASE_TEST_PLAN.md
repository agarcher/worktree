# Release Flow Test Plan

## Phase 1: Local Testing (No Side Effects)

### 1.1 Test Formula Generation
```bash
VERSION=0.1.0 \
SHA_DARWIN_AMD64=abc123 \
SHA_DARWIN_ARM64=def456 \
SHA_LINUX_AMD64=ghi789 \
SHA_LINUX_ARM64=jkl012 \
./scripts/generate-formula.sh
```
Verify: Valid Ruby, correct URLs (`agarcher/wt`), placeholders replaced.

### 1.2 Test Release Script Logic (Without Executing)
```bash
# Verify VERSION file
cat VERSION  # Should be 0.0.0

# Dry-run the version bump logic
bash -c 'VERSION="0.0.0"; IFS="." read -r M m p <<< "$VERSION"; echo "minor: $M.$((m+1)).0"'
# Should output: minor: 0.1.0
```

### 1.3 Verify Build Works
```bash
make build && ./build/wt version
```

### 1.4 Test `/release` Command (Decline at Confirmation)
Run `/release` in Claude Code. It should:
- Read VERSION (0.0.0) and detect no prior tags
- Analyze all commits on the branch
- Propose release notes and bump type (likely `minor` for initial release)
- Ask for approval before executing

**Decline the approval** to verify the workflow without actually releasing.
Check: Notes are concise, bump recommendation makes sense.

## Phase 2: Verify Prerequisites

### 2.1 Tap Repo Setup
```bash
# Check tap repo exists with Formula directory
gh repo view agarcher/homebrew-tap --json name
git ls-remote git@github.com:agarcher/homebrew-tap.git  # Tests SSH access
```

### 2.2 Deploy Key (Manual Check)
- `agarcher/homebrew-tap` → Settings → Deploy keys → Should have write-enabled key
- `agarcher/wt` → Settings → Secrets → Should have `TAP_DEPLOY_KEY`

## Phase 3: Real Release

Once local tests pass:
```bash
make release minor "Initial release"
```

Watch GitHub Actions: https://github.com/agarcher/wt/actions

## Cleanup (If Something Fails)

```bash
# Delete tag locally and remotely
git tag -d v0.1.0
git push origin --delete v0.1.0

# Delete GitHub release (if created)
gh release delete v0.1.0 --yes

# Reset VERSION
echo "0.0.0" > VERSION
git add VERSION && git commit -m "Reset VERSION" && git push
```
