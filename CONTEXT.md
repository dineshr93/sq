# sq

A read-only CLI that loads SBOM files and renders them as human-readable tables for security and legal review.

## Language

**SBOM**:
A Software Bill of Materials: one document declaring what open source components a piece of software contains and how they relate.
_Avoid_: manifest, inventory

**Document**:
One SBOM file as loaded by sq. It carries a spec version, but that version is invisible in every view.
_Avoid_: report, report file

**Element**:
SPDX 3.0's universal unit — packages, files, relationships, licenses, and agents are all elements. Internal term only; never surfaces in sq output.
_Avoid_: using "element" for packages/files in user-facing text

**Package**:
A distinct software component named by the document — the unit of `sq pkgs`.
_Avoid_: library, dependency, component

**File**:
A single named artifact a package contains — the unit of `sq files`.
_Avoid_: path, artifact

**Relationship**:
A typed, directed statement between two things a document describes — the unit of `sq rels` and `sq dig`.
_Avoid_: link, edge

**Root package**:
The package a document is primarily about: what 2.x `documentDescribes` and 3.0 `rootElement` both resolve to.
_Avoid_: main package, described package

**Declared license**:
The license a package's authors state it carries.

**Concluded license**:
The license the document's author assigns to a package after reviewing its actual licensing.

**Profile**:
A named vocabulary subset of SPDX 3.0. Core and Software are in scope; Dataset, AI, and Build are not.

**Unsupported element**:
An element from an out-of-scope profile (AI, Dataset, Build). Counted in `sq meta`, never rendered.
_Avoid_: skipped element (ambiguous with parse-dropped elements)

**Normalized vocabulary**:
sq's single package/file/relationship view regardless of spec version: 3.0 relationship types and directions are translated back into their 2.x equivalents where equivalents exist; genuinely new 3.0 relationship types appear verbatim under their SPDX names.
_Avoid_: IRI, spdxId, type URIs in output
