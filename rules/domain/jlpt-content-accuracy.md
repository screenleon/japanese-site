# Domain: JLPT Content Accuracy

Domain rules for keeping content tagged at a JLPT level usable by a learner at
that level. These rules guard against drift, especially for LLM-generated rows.

## Rules

### Rule: JLPT-001

- Owner layer: Domain
- Domain: jlpt-content-accuracy
- Stability: core
- Status: active
- Scope: any row with a `jlpt_level` column (vocab, kanji, sentence, grammar_point, question)
- Statement: A row tagged at JLPT level N MUST NOT use vocabulary or grammar above level N − 1, unless every above-level token is annotated with a furigana / gloss / explanation in the same row. Above-level tokens without scaffolding are a content bug, not a stylistic preference.
- Rationale: A learner studying N3 hits a sentence loaded with N1 grammar and abandons. The whole point of level tagging is the contract that the level is honored.
- Verification: content-validator agent compares each token against JMdict's JLPT tags and the grammar-point table; fails any row with above-level tokens that lack scaffolding.
- Supersedes: N/A
- Superseded by: N/A

### Rule: JLPT-002

- Owner layer: Domain
- Domain: jlpt-content-accuracy
- Stability: core
- Status: active
- Scope: question generation prompts (server-provider and local connector)
- Statement: The connector envelope payload for question generation MUST include `jlpt_level` and `target_grammar_point`. Generators MUST refuse to produce a question if `jlpt_level` is absent. The validator rejects any generated question whose `target_grammar_point` does not actually appear in the question stem.
- Rationale: Without this, generators silently produce off-target items (e.g., asked for ては-form, returned a て-form question). The grammar-point pointer is also what powers the corrective feedback's "see this grammar point" link.
- Verification: Schema check on the envelope; validator unit test asserting rejection of off-target items.
- Supersedes: N/A
- Superseded by: N/A

### Rule: JLPT-003

- Owner layer: Domain
- Domain: jlpt-content-accuracy
- Stability: behavior
- Status: active
- Scope: kanji rows
- Statement: When a kanji is rendered in any served sentence at level N3 or below, it MUST be accompanied by furigana in the response payload. Above N3, furigana is opt-in based on the user's furigana preference.
- Rationale: Beginner / intermediate learners cannot read raw kanji; missing furigana on the lower levels is a regression that breaks the learning loop entirely.
- Verification: Integration test that fetches an N4 sentence and asserts every kanji range has an accompanying reading in the response.
- Supersedes: N/A
- Superseded by: N/A
