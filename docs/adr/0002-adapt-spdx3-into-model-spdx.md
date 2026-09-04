# Adapt SPDX 3.0 into the existing model.SPDX, not a neutral IR

Loading 3.0 requires mapping its graph (Elements, typed relationships, `@id` references) into something the five renderers, `dig`, and `ip` can print. We decided the target is today's `model.SPDX` flat struct, with a hand-written forward mapping (meta→CreationInfo fields, Software packages/files, relationship back-translation, license-expression strings) and a small sidecar struct for the 3.0-only extras (profile list, skipped-element count). All five commands run untouched on 3.0 input.

## Considered Options

- Neutral intermediate representation: rejected — cleanest long-term, but rewrites the working 2.x load path and every renderer's input type for a feature that needs three extra fields. Revisit only if 3.0-only data outgrows the sidecar.
- Native v3 renderers per command: rejected — contradicts the normalized-vocabulary rule in `CONTEXT.md` (users must not see version-specific output surfaces).

## Consequences

- The mapping is bounded by what sq prints, not by the 2.3 spec — deliberately incomplete, and every dropped field is a future gap to fix in one place.
- The relationship back-translation table and the license walk are pure functions with pinned golden fixtures (tools-golang `testdata`, spdx-examples `example11`), so spec drift on the rc channel is caught by CI.
