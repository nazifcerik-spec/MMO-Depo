# Git (Version Control)

## 1. Amaç
Kod, içerik ve doküman değişimlerinin izlenebilir, auditable, kolay geri alınabilir olması. Branşlaşma stratejisinin herkesin aynı sayfada olduğu bir sözleşme olması.

## 2. Repository Yapısı
- **Tek monorepo (`mmo`)** başlangıçta:
  - `client/` (Flutter)
  - `server/` (Go)
  - `content/` (data pack: yaml/json)
  - `docs/` (mevcut)
  - `deploy/` (dockerfile, helm, terraform Phase 02)
  - `scripts/` (dev/ci/release)
  - `.github/` veya `.gitlab-ci/` (workflow tanımları)
- Ayırma ihtimali (Phase 02): client ve server ayrı repolara bölmek — sadece takım büyüdüğünde.

## 3. Branch Model — Trunk-Based
- **`main`**: her zaman deploy-edilebilir.
- **Feature branches**: `feat/<scope>-<short>` (örn: `feat/afk-reveal-batch`).
- **Fix branches**: `fix/<scope>-<short>`.
- **Chore/Docs branches**: `chore/`, `docs/`, `refactor/`, `perf/`, `test/`, `ci/`.
- Ömür: kısa (< 3 gün ideal); büyük feature'lar feature flag ile trunk'a incremental merge.
- **Release branch yok** başlangıçta; tag'ler `main` üstünden.
- Hotfix: `main`'den branch → merge back (release cut'a paralel değil, tek yol).

## 4. Commit Konvansiyonu — Conventional Commits
```
<type>(<scope>): <short imperative summary>

[Optional body — why, not what]

[Optional footer — BREAKING CHANGE, Refs: #123]
```
Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `build`, `revert`.

Örnekler:
- `feat(afk): add deterministic reveal batch worker`
- `fix(marketplace): prevent double-charge on cancel race`
- `docs(gdd): finalize Phase 00 architecture link`
- `refactor(character)!: split repository interface — BREAKING`

Kurallar:
- Imperative mood (İngilizce), 72 karakter altı özet.
- Body wrap 100 karakter.
- Referans zorunlu değil ama tercih (issue/PR).
- Attribution (Claude Code kullanıldığında):
  ```
  Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
  ```

## 5. Pull Request Kuralları
- **Küçük PR** (< 400 satır diff ideal); büyükse split.
- Description formatı:
  - **Summary** (1-3 madde)
  - **Test plan** (checkbox)
  - **Screenshots** (UI change ise)
  - **Migration/impact notes** (varsa)
- Reviewer: min 1 approval.
- CI green zorunlu (lint + test + build).
- Squash merge default (commit message = PR title).
- Merge sonrası branch delete (otomatik).

## 6. Tagging & Releases
- Semver: `vMAJOR.MINOR.PATCH`.
- Release note otomatik (Conventional Commits'ten changelog gen).
- Artifact:
  - Docker image tag = git sha + semver.
  - Flutter build artifact (web + android + ios) — CI upload.
- Pre-release: `v1.0.0-rc.1`, `v1.0.0-beta.3`.

## 7. Protected Branches
- `main` protected:
  - Force push yasak.
  - Direct push yasak; sadece PR merge.
  - Required status checks: lint, test, build.
  - Required review: 1 approval.
- Tag protection: `v*` tag'leri korumalı; sadece release worker/otomasyon.

## 8. Signing
- Commit signing (`gpg` veya `ssh-key sign`) — opsiyonel v1, zorunlu v2.
- Tag signing zorunlu release'lerde.

## 9. .gitignore Öncelikleri
- `.env*` (dosyaları — `.env.example` hariç).
- `**/build/`, `**/dist/`, `**/.dart_tool/`, `**/coverage/`.
- `*.log`, `*.tmp`.
- Local IDE: `.idea/`, `.vscode/settings.json` (paylaşımlı `.vscode/extensions.json` hariç).
- Container volume dumps.

## 10. Large Files
- Git LFS: içerik pack'inde büyük asset (art, atlas, audio) için.
- LFS quota izleme; alternatif: object storage + manifest (Phase 02).

## 11. Hooks
- `pre-commit`: format (`gofmt`, `dart format`), lint fast path.
- `commit-msg`: Conventional Commit validator.
- CI'da hook'ları duplicate ederiz — pre-commit atlanırsa CI yakalar.

## 12. History Hygiene
- Rebase before merge (linear history hedefi + squash).
- Force-push yalnızca kendi feature branch'ine, PR öncesi.
- `main` üstünde asla `git reset --hard`; rollback = revert PR.
- Interactive rebase yalnızca yerel; paylaşılan history değişmez.

## 13. Content Repo Discipline
- `content/` yaml/json değişimleri: mutlaka `docs/decisions/` veya PR body'de gerekçe.
- Content patch'ler bir "content release" olarak version'lu (schema_version).

## 14. Anti-Patterns
- 1000+ satır PR.
- "wip", "fix stuff", "asdf" commit message.
- Uzun ömürlü feature branch (haftalar).
- Force-push shared branch.
- `.env` commit.

## 15. İlişkiler
- [CI/CD](cicd.md), [Config](config.md), [Environment](environment.md), [Docker](docker.md)

## 16. Bekleyen Kararlar
- Host: GitHub vs GitLab vs self-hosted (öneri: GitHub).
- LFS host: GitHub LFS quota mı dedicated (S3-benzeri) mı.
- Monorepo tool (Nx, Turborepo) gerekli mi — v1'de gerekmiyor; v2 hızlanma için değerlendir.
