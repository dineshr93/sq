# Support SPDX 3.0.1 via prerelease tools-golang, pinned

sq parses SPDX 2.x JSON only; SPDX 3.0 is a JSON-LD rewrite that `encoding/json` against the 2.x structs cannot read, and no stable Go reader exists: `github.com/spdx/tools-golang` v0.5.7 is 2.x-only, and 3.0 support exists solely on the `v0.6.0-rc4` release candidate, whose model targets SPDX 3.0.1 exactly. We decided to take `v0.6.0-rc4` pinned, support 3.0.1 only, and gate on the `@context` IRI (`https://spdx.org/rdf/3.0.1/spdx-context.jsonld`), reporting a clean version error for other 3.x contexts rather than a parse failure.

## Considered Options

- Hand-rolled JSON-LD reader: rejected — sq would permanently own `@context` expansion, compaction, and IRI resolution for a viewer whose job is rendering.
- Track the `spdx3` branch at `main`: rejected — an unpinned moving target from a branch that may be rewritten before v0.6.0.
- Wait for stable v0.6.0: rejected — real 3.0.1 files exist now (spdx-examples), and the pin is contained by the adapter (ADR-0002), so the upgrade lands as one commit.

## Consequences

- The module's Go floor rises to 1.23.5 (rc4's requirement): `go.mod`, `Dockerfile`, and CI follow.
- Prerelease API churn risk is real but bounded to the adapter package; `go.mod` pins the exact pseudo-version and CI runs against it.
- 3.0.0 and 3.1 files are rejected at detection with an explicit "supports SPDX 3.0.1" message — the parser's own failure mode (`context is not known`) is strictly worse UX and would leak library internals.
