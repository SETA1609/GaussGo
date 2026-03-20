# Linear Algebra Mod Content Plan (from course PDFs)

## Scope
Build a progressive content roadmap for `mods/linearAlgebra` using the extracted curriculum from:
- `pdfs/Modulo-Algebra-lineal.pdf` (course module structure)
- `pdfs/Linear+Algebra.formulas.pdf` (formula reference)

This plan targets the current GaussGo learning model (units -> concepts -> exercises) and keeps IDs/data ready for state tracking (`progress.<modId>.readConcepts`, stats, quiz pools).

## Current Baseline
- Existing unit file: `mods/linearAlgebra/data/units/vectors.json`
- Existing concepts are minimal (`vector-basics`, `dot-product`)
- Manifest already depends on `core` and provides learning features

## Curriculum Mapping (PDF -> Mod Units)

### Unit 1: 2x2 Linear Systems Foundations
Source topics:
- System types (unique, none, infinite solutions)
- Methods: graphing, equalization, substitution, elimination, Cramer
- Application problem translation workflow

Planned file:
- `mods/linearAlgebra/data/units/systems-2x2.json`

Planned concepts:
- `systems.classification`
- `systems.graphical-method`
- `systems.substitution`
- `systems.elimination`
- `systems.cramer-2x2`
- `systems.word-problems`

Exercise clusters:
- classify-solution-type (multiple_choice)
- solve-2x2-step (compute)
- choose-method (multiple_choice)
- model-from-text (compute/multiple_choice hybrid)

### Unit 2: Matrix Fundamentals and Algebra
Source topics:
- Matrix definitions, order, matrix types, equality
- Matrix operations and properties

Planned file:
- `mods/linearAlgebra/data/units/matrices-basics.json`

Planned concepts:
- `matrices.definition-order`
- `matrices.types`
- `matrices.equality`
- `matrices.add-subtract`
- `matrices.scalar-mul`
- `matrices.matrix-mul`

Exercise clusters:
- identify-order-and-type
- compute-basic-operations
- check-dimension-compatibility

### Unit 3: Matrix Methods for Linear Systems
Source topics:
- Matrix form of systems (`A`, `X`, `B`, augmented)
- Row operations, echelon form, pivots
- Gaussian and Gauss-Jordan methods
- Inverse method and determinant method
- Homogeneous systems

Planned file:
- `mods/linearAlgebra/data/units/systems-matrix-methods.json`

Planned concepts:
- `systems.matrix-form`
- `systems.row-operations`
- `systems.gaussian-elimination`
- `systems.gauss-jordan`
- `systems.inverse-method`
- `systems.homogeneous`

Exercise clusters:
- row-operation-next-step
- reduced-echelon-result
- solve-system-by-method
- consistency-diagnosis

### Unit 4: Vector Spaces, Subspaces, Basis, Dimension
Source topics:
- Vector space axioms and examples
- Subspace criteria
- Linear combinations, dependence/independence
- Basis and dimension

Planned file:
- `mods/linearAlgebra/data/units/vector-spaces.json`

Planned concepts:
- `spaces.vector-space-axioms`
- `spaces.subspace-test`
- `spaces.linear-combination`
- `spaces.independence`
- `spaces.basis`
- `spaces.dimension`

Exercise clusters:
- verify-subspace
- independence-check
- find-basis
- compute-dimension

### Unit 5: Eigenvalues, Eigenvectors, Similarity, Diagonalization
Source topics:
- Similar/equivalent/congruent matrix distinctions
- Characteristic polynomial and eigenpairs
- Diagonalization criteria and process
- Applications

Planned file:
- `mods/linearAlgebra/data/units/eigen-diagonalization.json`

Planned concepts:
- `eigen.characteristic-polynomial`
- `eigen.eigenvalues`
- `eigen.eigenvectors`
- `eigen.similarity`
- `eigen.diagonalization`

Exercise clusters:
- compute-characteristic-polynomial
- find-eigenvalues-vectors
- check-diagonalizable
- build-P-D-inverse

### Unit 6: Vectors, Lines, and Planes in R2/R3
Source topics:
- Vector representation and operations in R2/R3
- Norm, direction angles, direction cosines
- Equations of lines and planes

Planned file:
- `mods/linearAlgebra/data/units/geometry-r2-r3.json`

Planned concepts:
- `geometry.vector-operations`
- `geometry.norm-and-direction`
- `geometry.dot-product-geometry`
- `geometry.cross-product`
- `geometry.lines-r3`
- `geometry.planes-r3`

Exercise clusters:
- norm-angle-cosines
- line-equation-from-points
- plane-equation-from-point-normal
- relation-line-plane (parallel/intersect)

## Formula Coverage Checklist (from formula sheet)
- 2x2 determinant and Cramer rules
- Matrix system form `AX=B`, augmented matrix usage
- Characteristic equation `det(A-lambda*I)=0`
- Similarity/diagonalization forms (`B=P^-1 A P`, `A=P D P^-1`)
- Vector norms in R2/R3
- Direction cosines and identity

Each formula group should appear in at least one concept explanation and one exercise set.

## Data Authoring Rules
- Keep unit file schema aligned with existing loader expectations:
  - root: `id`, `title`, `description`, `concepts`
  - each concept: `id`, `title`, `explanation`, optional `examples`, optional `exercise`
- Use stable namespaced concept IDs (`systems.*`, `matrices.*`, `spaces.*`, `eigen.*`, `geometry.*`).
- Keep explanations short first, then expand in examples.
- Start with one graded exercise per concept, then add extra drills in second pass.

## Incremental Delivery Plan

### Slice A (quick unlock)
1. Add `systems-2x2.json`
2. Expand existing `vectors.json` with geometry-ready IDs
3. Ensure random quiz can draw from at least 10 read concepts

### Slice B (core matrix path)
1. Add `matrices-basics.json`
2. Add `systems-matrix-methods.json`
3. Add mixed matrix/system compute exercises

### Slice C (advanced linear algebra)
1. Add `vector-spaces.json`
2. Add `eigen-diagonalization.json`
3. Add theorem-guided examples and checks

### Slice D (geometry completion)
1. Add `geometry-r2-r3.json`
2. Add line/plane applied problems
3. Final quiz pool balancing across all units

## Suggested File Backlog
- `mods/linearAlgebra/data/units/systems-2x2.json`
- `mods/linearAlgebra/data/units/matrices-basics.json`
- `mods/linearAlgebra/data/units/systems-matrix-methods.json`
- `mods/linearAlgebra/data/units/vector-spaces.json`
- `mods/linearAlgebra/data/units/eigen-diagonalization.json`
- `mods/linearAlgebra/data/units/geometry-r2-r3.json`

## Acceptance Criteria for Content Plan Completion
- All six PDF units are represented in mod unit files.
- Every concept has at least one exercise.
- Concept IDs are stable and quiz-eligible.
- Formula groups are linked to at least one concept and one exercise.
- Content is balanced from introductory to advanced topics for progressive learning.
