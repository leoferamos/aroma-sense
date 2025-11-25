# 🚀 Release Automation

This project uses **Semantic Release** to automate versioning, tag creation and GitHub releases.

## ✨ How It Works

1. **Open a PR** → **Tests run automatically** (lint, test, build)
2. **Merge PR to main** → **Release workflow triggers**
3. **Automatic release**:
   - Analyzes conventional commits
   - Determines new version
   - Creates tag and GitHub release
   - Updates `package.json`, `package-lock.json` and `CHANGELOG.md`

## 📝 Conventional Commits

Use the format: `type(scope): description`

### Types
- `feat:` - New feature (increments MINOR: 1.0.0 → 1.1.0)
- `fix:` - Bug fix (increments PATCH: 1.0.0 → 1.0.1)
- `docs:` - Documentation
- `style:` - Formatting/style
- `refactor:` - Refactoring
- `perf:` - Performance
- `test:` - Tests
- `chore:` - Maintenance

### Examples
```bash
git commit -m "feat: add user login system"
git commit -m "fix: resolve refresh token panic"
git commit -m "docs: update API documentation"
```

## 🔄 Next Versions

- **PATCH** (0.4.1): At least 1 `fix:` commit
- **MINOR** (0.5.0): At least 1 `feat:` commit
- **MAJOR** (1.0.0): Breaking change (use `!` or `BREAKING CHANGE:`)

## 📋 Release Checklist

- [ ] PR checks pass (tests, lint, build)
- [ ] Commits follow conventional commits
- [ ] Code reviewed and approved
- [ ] Merge to `main` triggers automatic release

## 🛠️ Local Development

```bash
# Install dependencies
npm install

# Development
npm run dev

# Build
npm run build

# Lint
npm run lint
```

## 📚 Useful Links

- [Conventional Commits](https://conventionalcommits.org/)
- [Semantic Release](https://semantic-release.gitbook.io/)
- [Commit Guide](./CONVENTIONAL_COMMITS.md)