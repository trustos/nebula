# hopssh Nebula Fork

This is a fork of [slackhq/nebula](https://github.com/slackhq/nebula) maintained
by [hopssh](https://hopssh.com) for performance and feature enhancements.

Upstream: `https://github.com/slackhq/nebula`
Branch: `hopssh` (based on upstream tag `v1.10.3`)
Consumed by: `github.com/trustos/hopssh` via `go.mod` replace directive

## Changes (on `hopssh` branch)

| Commit | Description | Files |
|--------|-------------|-------|
| 1. Graceful shutdown | Fix `os.Exit(2)` on service close — check `io.ErrClosedPipe` and `io.EOF` | `interface.go` |
| 2. TUN buffer reuse | Reuse read buffer + write mutex for concurrent safety on macOS | `overlay/tun_darwin.go` |
| 3. UDP socket buffers | `SO_RCVBUF`/`SO_SNDBUF` support + `SupportsMultipleReaders()` on macOS | `udp/udp_darwin.go` |
| 4. Decoupled routines | Separate TUN/UDP routine counts — macOS gets multi-reader UDP with single TUN | `interface.go` |
| 5. Packet coalescing | Batch UDP sends with length-prefix framing + Linux `StdConn` panic fix | `udp/coalesce.go`, `interface.go`, `udp/udp_linux.go` |
| 6. PMTUD support | `SetMTU` on Device interface, `SendTestRequest`, `TestReply` callback | `overlay/device.go`, `overlay/tun_*.go`, `control.go`, `outside.go`, `interface.go` |

## How hopssh Consumes This Fork

In hopssh's `go.mod`:

```go
require github.com/slackhq/nebula v1.10.3

replace github.com/slackhq/nebula => github.com/trustos/nebula v1.10.3-hopssh.1
```

All internal imports remain `github.com/slackhq/nebula` — no code changes needed
in either the fork or hopssh.

## Adding a New Change

```bash
git clone git@github.com:trustos/nebula.git
cd nebula
git checkout hopssh

# Make your changes
vim overlay/tun_darwin.go

# Commit
git add -A
git commit -m "Description of change"

# Tag and push
git tag v1.10.3-hopssh.2    # increment the suffix
git push origin hopssh --tags

# In hopssh: update go.mod
# replace github.com/slackhq/nebula => github.com/trustos/nebula v1.10.3-hopssh.2
go get github.com/trustos/nebula@v1.10.3-hopssh.2
go mod tidy
```

## Upgrading to a New Upstream Nebula Release

When upstream releases v1.11.0:

```bash
cd nebula
git fetch upstream
git checkout hopssh

# Rebase our commits on top of the new release
git rebase upstream/v1.11.0

# Resolve any conflicts — our commits are clean and isolated
# Each commit modifies different files, so conflicts are rare

# Tag the new version
git tag v1.11.0-hopssh.1
git push origin hopssh --tags --force-with-lease

# In hopssh: update go.mod
# replace github.com/slackhq/nebula => github.com/trustos/nebula v1.11.0-hopssh.1
go get github.com/trustos/nebula@v1.11.0-hopssh.1
go mod tidy
```

## Contributing Back Upstream

Some of our changes may be useful to upstream Nebula:

- **Graceful shutdown** — upstream PR #1375 already exists
- **Packet coalescing** — could be a standalone PR
- **PMTUD / SetMTU** — could be proposed as a feature

If upstream merges a change, remove the corresponding commit from our branch
during the next rebase (it will auto-resolve since the code matches).
