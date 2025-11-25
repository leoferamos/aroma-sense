# Conventional Commits Guide

This project uses [Conventional Commits](https://conventionalcommits.org/) to automate versioning and releases.

## Commit Formats

### Main Types
- `feat:` - New feature (increments MINOR)
- `fix:` - Bug fix (increments PATCH)
- `docs:` - Documentation changes
- `style:` - Style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or fixing tests
- `chore:` - Tool, config, or maintenance changes

### Examples
```
feat: add user authentication system
fix: resolve login panic on refresh token
docs: update API documentation
chore: update dependencies
```

### Breaking Changes
For changes that break compatibility, add `!` after the type:
```
feat!: change API response format
```

Or describe in the commit body:
```
feat: change API response format

BREAKING CHANGE: The response format has changed from XML to JSON
```

## How It Works

1. **Conventional commits** → **Automatic versioning** → **GitHub release**
2. Semantic-release analyzes commits since the last release
3. Determines the version type (major, minor, patch)
4. Creates tag, release and updates CHANGELOG.md automatically

## Next Versions

- `fix:` commits = **PATCH** (0.4.1 → 0.4.2)
- `feat:` commits = **MINOR** (0.4.1 → 0.5.0)
- Breaking changes = **MAJOR** (0.4.1 → 1.0.0)