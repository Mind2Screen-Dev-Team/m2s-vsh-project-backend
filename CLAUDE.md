# CLAUDE.md — M2S-VSH Project Backend

Backend service untuk pilot M2S-VSH Lite v0.1.0. Bagian dari workflow ter-orchestrasi control repo `m2s-vsh-platform` (Go). Governance, task contract, dan pipeline hidup di control repo — lihat README untuk referensi.

Stack: **Go 1.26.4** (module `github.com/Mind2Screen-Dev-Team/m2s-vsh-project-backend`). HTTP service sederhana, belum ada framework/dependency eksternal.

## Perintah

```sh
go test ./...      # unit test (termasuk race: go test -race ./...)
go vet ./...       # static analysis
go build ./...     # build
go run ./cmd/server # jalankan server lokal
```

## Struktur

| Path | Isi |
|---|---|
| `cmd/server/` | Entry point `main.go` (HTTP server) |
| `internal/handler/` | HTTP handler + unit test colocated (`status.go`, `status_test.go`) |
| `.github/` | CI `path-enforcement.yml`, CODEOWNERS, PR template |

Pattern: handler colocated dengan test-nya di `internal/handler/`. Ikuti pattern yang sudah ada, jangan bikin layout baru tanpa alasan.

## Branch & alur kerja

- Branch: `main` (production, default), `staging` (pre-production), `develop` (integrasi).
- Seluruh perubahan masuk via PR. Task dikerjakan di worktree branch `agent/<TASK-ID>-slug`, PR target **`develop`**. Merge ke `main` = manusia.
- Setiap task didefinisikan di control repo (`control/tasks/specifications/<TASK-ID>.yaml`) — contract adalah source of truth, bukan asumsi lokal.
- Detail workflow: AGENT.md repo ini + control repo.

## Konvensi

- Commits: **Conventional Commits** (feat, fix, docs, chore, refactor, test).
- Go standard: `gofmt`, `go vet`, test wajib lulus sebelum selesai.
- Bahasa: README & dokumentasi berbahasa Indonesia.
