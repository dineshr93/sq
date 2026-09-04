# Synthesise SoftwareAgent stubs to load cross-document references

SPDX 3.x documents legitimately reference elements defined in other
documents (e.g. `FROM` a source package in a sibling SBOM). tools-golang
v3.0.0-rc4 expands such a document with JSON-LD frame semantics and then
fails the whole parse with `unable to find viable external IRI` when the
reference target is typed `AnyElement` and its id is not in the document —
rc4 registers external-IRI constructors for concrete types but none for
`AnyElement`. Dropping the references loses real data sq prints.

We decided to keep the retry-injection loader: parse once; if the failure
is the external-`AnyElement` kind, collect the offending ids, append
minimal stub nodes `{"@id": id, "type": "SoftwareAgent"}` to `@graph`, and
re-parse. The stub makes rc4 treat the endpoint as an internal reference,
the document loads, and every id that was synthesised is counted as
`ExternalElement` in the meta sidecar so the user sees the document has
out-of-document references.

## Considered Options

- Fix upstream in tools-golang and wait: correct long-term — file it —
  but sq targets an rc channel we do not control; a loader fallback is
  deletable once upstream ships the `AnyElement` constructor.
- Custom JSON-LD frame marking externals as named individuals: avoids
  fake types but requires owning a fork of rc4's frame/context pipeline.
- Drop unresolvable relationship endpoints: silently loses the
  cross-document links the fixtures exist to exercise.

## Consequences

- Stubs carry a type the document never asserted (`SoftwareAgent` is a
  placeholder, not a claim). They are reported only as `ExternalElement`
  counts, never rendered as agents — do not let them leak into agent or
  package views.
- The stub must use the `@id` keyword, not the compact `id` alias; the
  offline context does not guarantee `id` expansion.
- rc4 validation would warn about stubs lacking `creationInfo`; the
  warning is suppressed exactly when injection occurred
  (`len(injected) > 0`), keeping warnings meaningful for real documents.
- If upstream adds the constructor, this path becomes dead code: the
  first parse succeeds, `injected` stays empty, behaviour is unchanged.
