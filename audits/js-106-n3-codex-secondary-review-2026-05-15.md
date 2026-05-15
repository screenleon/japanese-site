# JS-106 N3 Inline Ruby — Codex Secondary Review (PR #59)

**Date**: 2026-05-15
**Reviewer**: codex-executor (secondary pass)
**Primary author**: Claude main thread
**Scope**: 38 N3 grammar entries' `explanation_ja_blocks` ruby tokens
**Limitation**: codex is not a native Japanese speaker; readings cross-checked against general training corpus knowledge.

## Per-entry verdicts

| Slug | Ruby count | Verdict | Notes |
|---|---:|---|---|
| bakari | 3 | OK | |
| concession-temo | 14 | OK | |
| conditional-ba | 18 | tweak | Tokens 4-5 and 6-7 split grammar metalanguage compounds that should be single ruby tokens. |
| conditional-nara | 10 | OK | |
| conditional-tara | 10 | OK | |
| contrast-noni | 19 | OK | |
| dake-denaku | 2 | OK | |
| dokoroka | 10 | OK | |
| hazuda | 16 | OK | |
| hazuganai | 6 | OK | |
| hodo | 16 | OK | |
| kadouka | 6 | OK | |
| kamoshirenai | 3 | OK | |
| kawari-ni | 11 | OK | |
| koto-ni-natte-iru | 14 | OK | |
| kotoni-naru | 9 | OK | |
| kotoni-suru | 4 | OK | |
| mitai | 4 | OK | |
| monono | 10 | OK | |
| ni-chigainai | 9 | OK | |
| okage-de | 8 | OK | |
| rashii | 8 | OK | |
| sae-ba | 8 | OK | |
| sei-de | 11 | OK | |
| souda-appearance | 5 | OK | |
| souda-hearsay | 8 | OK | |
| tabi-ni | 10 | OK | |
| tameni-purpose | 9 | OK | |
| tara-dou | 7 | OK | |
| te-bakari-iru | 8 | OK | |
| teiku-tekuru | 4 | OK | |
| teshimau | 8 | OK | |
| to-sureba | 7 | OK | |
| tokoro | 5 | OK | |
| wake-niwa-ikanai | 12 | OK | |
| youda | 9 | OK | |
| youni-goal | 9 | OK | |
| youni-naru | 11 | OK | |

## Detail (for tweak / block rows only)

### conditional-ba

- Tokens 4-5: form `K=一段, r=いちだん` followed by `K=動詞, r=どうし` — issue: `一段動詞` is a grammar metalanguage compound and the migration convention says to ruby `一段動詞` as one compound token — suggest: `K=一段動詞, r=いちだんどうし`.
- Tokens 6-7: form `K=五段, r=ごだん` followed by `K=動詞, r=どうし` — issue: `五段動詞` is a grammar metalanguage compound and the migration convention says to ruby `五段動詞` as one compound token — suggest: `K=五段動詞, r=ごだんどうし`.

## Summary

- 38 entries reviewed.
- 37 OK, 1 tweak, 0 block.
- 1 tweak entry detailed above.
- No wrong readings found.
- No same-compound reading collisions found across the reviewed files.
- No confident missing-ruby issues found for obvious N3+ content-noun compounds in surrounding text.

## Process note

This is the codex secondary pass per `feedback_native_reviewer_role` two-pass model. Claude main thread is primary author. Codex's role here is independent double-check, NOT primary judgment.
