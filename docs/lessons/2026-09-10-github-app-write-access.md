# GitHub App needs write access to push

- **Date:** 2026-09-10
- **Area:** tooling

## What happened

The first push of the README cleanup failed:

```
remote: Claude doesn't have GitHub access to JL037/firehose-relay-go ...
fatal: unable to access 'https://github.com/JL037/firehose-relay-go/': 403
```

The GitHub MCP path could *read* the repo but not write — branch creation
returned `403 Resource not accessible by integration`.

## Cause

The Claude GitHub App was installed with **read-only** access to the repo.

## Fix / takeaway

Grant the app write access (app installation settings, or reconnect GitHub in
claude.ai → Settings → Connectors), then the same `git push` succeeds. When a
push 403s, check the app's repo permission scope first — it isn't a network or
auth-token problem, and retrying the push verbatim won't help.
