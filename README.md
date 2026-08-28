# Atlas

Atlas is the AI Deck's standalone structural-evidence provider. It maps
declarations and bounded symbol relations without calling a model or retaining
source bodies in its export.

```powershell
atlas --root C:\path\to\repo --output atlas.snapshot.json snapshot
```

The resulting `atlas.snapshot.v1` is the optional provider-neutral input for
Lexicon's code-map tier. Atlas can be used by itself; Lexicon does not own,
index, or require it.

Atlas owns the neutral snapshot engine. The initial tree-sitter parser is an
explicit `atlas/adapters/pith` compatibility adapter, so it can be replaced by
an LSP/compiler/index provider without altering Lexicon or the snapshot wire
contract. Legacy Pith is no longer an implicit Atlas dependency.

## Thread

`Thread` is the bounded structural path format used by Lens for a live
behavioral map. It replaces the internal name “BehaviorSlice”: it is grounded
only in Atlas declarations and call edges, never a model-invented explanation.

```powershell
atlas --root C:\path\to\repo --output atlas.snapshot.json snapshot
atlas --snapshot atlas.snapshot.json --symbol Apply --cap 32 thread
```

The output is `atlas.thread.v1`: a jumpable entry declaration plus a bounded
breadth-first call path, including known caller and callee identities. When the
node cap is reached, `truncated: true` makes that limitation explicit. It states only what Atlas actually knows; purposes,
side effects, invariants, and prose explanations are later evidence-backed
annotations, not fields silently fabricated by the structural engine.
