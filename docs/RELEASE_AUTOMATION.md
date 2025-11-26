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

## 📋 Release Notes

For more detailed releases, we use custom templates in `./releases/`:
- `TEMPLATE_MINOR.md` for minor releases
- `TEMPLATE_PATCH.md` for patch releases
- `TEMPLATE_MAJOR.md` for major releases

The `CHANGELOG.md` is generated automatically by semantic-release, but GitHub release descriptions use the templates for better narrative.

## 📝 Conventional Commits

See [Conventional Commits Guide](./CONVENTIONAL_COMMITS.md) for detailed rules and examples.

## 🔄 Next Versions

- **PATCH** (0.4.1): At least 1 `fix:` commit
- **MINOR** (0.5.0): At least 1 `feat:` commit
- **MAJOR** (1.0.0): Breaking change (use `!` or `BREAKING CHANGE:`)

## 📋 Release Checklist

- [ ] PR checks pass (tests, lint, build)
- [ ] Commits follow conventional commits
- [ ] Code reviewed and approved
- [ ] Merge to `main` triggers automatic release
## 🔧 Troubleshooting

### Release not triggered
- Check if commits follow conventional commits (`feat:`, `fix:`, etc.)
- Confirm there are changes since the last tag
- Check workflow logs in GitHub Actions

### Semantic-release error
- Verify configuration in `.releaserc`
- For manual releases, use `npx semantic-release --dry-run`

## 📚 Useful Links

- [Conventional Commits](https://conventionalcommits.org/)
- [Semantic Release](https://semantic-release.gitbook.io/)
- [Commit Guide](./CONVENTIONAL_COMMITS.md)
- [Release Templates](./releases/)