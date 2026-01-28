# .gitignore Guide for Monorepo

## Overview
This monorepo uses a hierarchical `.gitignore` structure to keep the repository clean and avoid committing unnecessary files.

---

## Structure

```
mono-starter/
├── .gitignore                  # Root - Handles common patterns for entire monorepo
├── apps/
│   ├── api/
│   │   └── .gitignore         # API-specific ignores (Go)
│   └── frontend/
│       └── .gitignore         # Frontend-specific ignores (React Router)
└── packages/
    └── (inherits from root)
```

---

## Root `.gitignore`

### Purpose
The root `.gitignore` handles patterns common to the entire monorepo:
- Node.js dependencies (`node_modules/`)
- Environment files (`.env*`)
- Build outputs (`dist/`, `build/`)
- OS files (`.DS_Store`, `Thumbs.db`)
- IDE files (`.vscode/`, `.idea/`)
- Logs (`*.log`)
- Testing artifacts (`coverage/`)
- Go binaries (`*.exe`, `*.so`)

### Why This Approach?
- **DRY Principle**: Avoid duplicating ignore patterns across multiple files
- **Consistency**: Same rules apply everywhere unless explicitly overridden
- **Maintainability**: Update common patterns in one place

---

## App-Specific `.gitignore` Files

### `apps/api/.gitignore` (Go API)

**API-specific patterns only**:
- Go binary outputs (`api`, `api.exe`, `main`)
- Air live reload temp directory (`tmp/`)
- API-specific data directories
- Swagger/OpenAPI generated files

**Why minimal?**
Most Go patterns (`.DS_Store`, `*.log`, `vendor/`, etc.) are already handled by root `.gitignore`.

### `apps/frontend/.gitignore` (React Router)

**Frontend-specific patterns only**:
- React Router build cache (`.react-router/`)
- Frontend build output (`build/`)
- Type generation artifacts

**Why minimal?**
Most JavaScript patterns (`node_modules/`, `.env`, `*.log`, etc.) are already handled by root `.gitignore`.

---

## What Gets Ignored

### ✅ Always Ignored
- `node_modules/` - NPM dependencies
- `.env*` - Environment variables (secrets)
- `build/`, `dist/` - Build outputs
- `.DS_Store`, `Thumbs.db` - OS files
- `*.log` - Log files
- `coverage/` - Test coverage reports
- `.turbo/` - Turbo cache
- `*.test`, `*.out` - Go test artifacts
- `bin/` - Binary outputs
- `.vscode/`, `.idea/` - IDE settings (with exceptions)

### ⚠️ Conditionally Ignored
- `.claude/` - Currently tracked (commented out in root gitignore)
- `.vscode/settings.json` - Tracked (useful for team)
- `go.work` - Ignored (local Go workspace)

### ✅ Always Tracked
- Source code (`*.ts`, `*.tsx`, `*.go`)
- Configuration files (`package.json`, `tsconfig.json`, `go.mod`)
- Documentation (`*.md`)
- Public assets
- Lock files (`pnpm-lock.yaml`, `go.sum`)

---

## Best Practices

### 1. **Hierarchy Principle**
```
Root .gitignore (broad patterns)
    ↓
App .gitignore (specific overrides/additions)
    ↓
Individual files
```

### 2. **When to Add to Root**
Add patterns to root `.gitignore` when:
- Pattern applies to multiple apps/packages
- It's a common development artifact
- It's OS or IDE related

**Example**: All apps generate logs, so `*.log` goes in root.

### 3. **When to Add to App-Specific**
Add patterns to app `.gitignore` when:
- Pattern is unique to that app
- It's a tool-specific artifact (Air for Go, React Router for frontend)
- It's app-specific data/cache

**Example**: `.react-router/` is frontend-only, so it goes in `apps/frontend/.gitignore`.

### 4. **Don't Duplicate**
❌ **Bad**:
```
# Root .gitignore
.DS_Store
node_modules/

# apps/frontend/.gitignore
.DS_Store          # Already in root!
node_modules/      # Already in root!
.react-router/     # Frontend-specific, good!
```

✅ **Good**:
```
# Root .gitignore
.DS_Store
node_modules/

# apps/frontend/.gitignore
.react-router/     # Only frontend-specific patterns
```

---

## Common Scenarios

### Scenario 1: New Tool Generates Artifacts
**Question**: Should I add the pattern to root or app-specific?

**Decision Tree**:
```
Will other apps use this tool?
├─ Yes → Add to root .gitignore
└─ No → Add to app-specific .gitignore
```

**Example**:
- Prettier cache → Root (all apps might use it)
- Go Air tmp → API-specific (only Go API uses it)

### Scenario 2: Secret Files
**Always use root `.gitignore`** for secrets:
- `.env*`
- `*.pem`, `*.key`, `*.cert`
- `secrets/` directories

**Why?** Security must be consistent everywhere.

### Scenario 3: Adding New App
When adding a new app (e.g., `apps/mobile`):

1. Create `apps/mobile/.gitignore`
2. Add ONLY mobile-specific patterns
3. Rely on root for common patterns

```
# apps/mobile/.gitignore
# Mobile-specific patterns
*.apk
*.ipa
android/build/
ios/build/
```

---

## Verification Commands

### Check what's ignored
```bash
git status --ignored
```

### Test if file would be ignored
```bash
git check-ignore -v path/to/file
```

### List tracked files
```bash
git ls-files
```

### Remove tracked file that should be ignored
```bash
# If you accidentally committed a file
git rm --cached path/to/file
git commit -m "Remove accidentally committed file"
```

---

## Cleaning Up

### Remove All `.DS_Store` Files
```bash
# Find and remove (already done)
find . -name ".DS_Store" -type f -delete

# Prevent creation on network drives (macOS)
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool true
```

### Remove All Untracked Files
```bash
# See what would be removed (dry run)
git clean -xdn

# Actually remove (careful!)
git clean -xdf
```

---

## IDE Configuration

### VSCode
The root `.gitignore` allows specific VSCode files:
```
.vscode/*
!.vscode/settings.json      # Team settings
!.vscode/tasks.json         # Build tasks
!.vscode/launch.json        # Debug configs
!.vscode/extensions.json    # Recommended extensions
```

### GoLand / IntelliJ
All `.idea/` files are ignored (personal preference varies).

### Vim
All swap files ignored (`*.swp`, `*.swo`, `*~`).

---

## Security Considerations

### Never Commit
- API keys, secrets, tokens
- Private keys (`.pem`, `.key`)
- Database credentials
- `.env` files with sensitive data

### Always Review
Before committing, check:
```bash
git diff --cached
```

### Use Environment Variables
Store secrets in environment variables, not in code:
```typescript
// ❌ Bad
const API_KEY = "sk-1234567890abcdef";

// ✅ Good
const API_KEY = process.env.API_KEY;
```

---

## Troubleshooting

### File Not Being Ignored
**Problem**: Added pattern to `.gitignore` but file still tracked.

**Solution**: Remove from tracking:
```bash
git rm --cached path/to/file
git commit -m "Stop tracking file"
```

### Pattern Not Working
**Problem**: Pattern doesn't seem to work.

**Debug**:
```bash
# Check which gitignore rule matches
git check-ignore -v path/to/file

# Test pattern matching
git ls-files --ignored --exclude-standard --others
```

### Too Much Ignored
**Problem**: Important file being ignored.

**Solution**: Use negation in `.gitignore`:
```
# Ignore all .env files
.env*

# But track example
!.env.example
```

---

## Maintenance

### Regular Review
Periodically check if patterns are still relevant:
```bash
# See what's ignored
git status --ignored

# Review patterns
cat .gitignore
```

### Update When Adding Dependencies
New tools often generate new artifacts:
- Check tool documentation
- Add necessary patterns
- Test with `git status`

### Team Communication
When updating `.gitignore`:
1. Communicate changes to team
2. Team members should run `git clean -xdf` (backup first!)
3. Reinstall dependencies if needed

---

## Quick Reference

| Pattern | Meaning | Example |
|---------|---------|---------|
| `file.txt` | Exact filename | Ignores all `file.txt` |
| `*.log` | Extension | All log files |
| `dir/` | Directory | Entire directory |
| `**/foo` | Any nested | `a/foo`, `a/b/foo` |
| `!file.txt` | Negate | Don't ignore this |
| `#` | Comment | For documentation |

---

## Resources

- [Git Documentation](https://git-scm.com/docs/gitignore)
- [GitHub's .gitignore Templates](https://github.com/github/gitignore)
- [Check Online](https://www.toptal.com/developers/gitignore)

---

## Summary

✅ **Root .gitignore**: Common patterns for entire monorepo
✅ **App .gitignore**: App-specific patterns only
✅ **No duplication**: DRY principle
✅ **Security first**: Always ignore secrets
✅ **Team-friendly**: Keep team settings, ignore personal files

This structure keeps the repository clean, secure, and maintainable across all apps and packages in the monorepo.
