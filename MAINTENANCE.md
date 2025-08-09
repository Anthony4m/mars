# Maintenance policy

- Track `main` branch; tag stable releases.
- Support window: critical bug fixes and security only.
- Hotfix flow: branch from tag → add tests → bump patch → tag `v1.0.x`.

## Cutting a hotfix
```
git checkout -b hotfix/v1.0.x
# commit fix + tests
npm run ? # (no js tooling; skip)
go test ./...
git commit -m "fix: ... (v1.0.x)"
git tag v1.0.x
```

