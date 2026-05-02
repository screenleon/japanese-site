# Kanji N4/N3/N2 Batch 1 Audit - 2026-05-02

## Counts

| Level | File | Entries |
|---|---|---:|
| N4 | `server/data/corpus/kanji/N4.jsonl` | 40 |
| N3 | `server/data/corpus/kanji/N3.jsonl` | 40 |
| N2 | `server/data/corpus/kanji/N2.jsonl` | 40 |

## Sources Used

- Tanos JLPT N4/N3/N2 kanji PDFs, local mirrors under `server/data/tanos_raw/jlpt{4,3,2}/KanjiList.N{4,3,2}.pdf`, cross-checked with public pages:
  - https://www.tanos.co.uk/jlpt/jlpt4/kanji/KanjiList.N4.pdf
  - https://www.tanos.co.uk/jlpt/jlpt3/kanji/KanjiList.N3.pdf
  - https://www.tanos.co.uk/jlpt/jlpt2/kanji/KanjiList.N2.pdf
- Tanos Anki character decks under `server/data/tanos_raw/jlpt{4,3,2}/n{4,3,2}-kanji-char-eng.anki`, used as a machine-readable cross-check of the same source family.
- Jisho kanji search pages (`https://jisho.org/search/<kanji>%20%23kanji`) and the bundled KANJIDIC2 English dataset (`server/data/external/kanjidic2-en-3.6.2%2B20260420131912.json.tgz`) for meaning sanity checks.
- Existing overlays `server/data/corpus/kanji/N5.jsonl` and `server/data/corpus/kanji/N1.jsonl` for schema and deduplication.

## Self-Check Sample

Random sample command: `shuf -n 3 server/data/corpus/kanji/N{4,3,2}.jsonl`.

| Level | Kanji | Level evidence | Meaning check | Result |
|---|---|---|---|---|
| N4 | 題 | Tanos N4 PDF/deck lists 題 as topic/subject; Jisho/KANJIDIC2 has topic/subject. | `問題。文章などの名。` matches topic/title/problem sense. | Pass |
| N4 | 開 | Tanos N4 PDF/deck lists 開 as open/unfold/unseal; Jisho/KANJIDIC2 has open. | `ひらくこと。あけること。` matches open sense. | Pass |
| N4 | 動 | Tanos N4 PDF/deck lists 動 as move/motion/change; Jisho/KANJIDIC2 has move/motion. | `うごくこと。うごかすこと。` matches move sense. | Pass |
| N3 | 法 | Tanos N3 PDF/deck lists 法 as method/law/rule; Jisho/KANJIDIC2 agrees. | `法律。決められた方法。` matches law/method sense. | Pass |
| N3 | 和 | Tanos N3 PDF/deck lists 和 as harmony/Japanese style/peace; Jisho/KANJIDIC2 agrees. | `やわらぐこと。日本風。` matches harmony/Japanese-style sense. | Pass |
| N3 | 性 | Tanos N3 PDF/deck lists 性 as nature/gender; Jisho/KANJIDIC2 has nature/sex/gender. | `ものの性質。生まれつきのあり方。` matches nature/quality sense. | Pass |
| N2 | 版 | Tanos N2 PDF/deck lists 版 as printing block/edition; Jisho/KANJIDIC2 agrees. | `印刷のもと。出版された形。` matches plate/edition sense. | Pass |
| N2 | 門 | Tanos N2 PDF/deck lists 門 as gates; Jisho/KANJIDIC2 has gate/class/branch. | `出入り口。専門の分野。` matches gate/specialty field sense. | Pass |
| N2 | 細 | Tanos N2 PDF/deck lists 細 as slender/narrow/fine; Jisho/KANJIDIC2 agrees. | `ほそいこと。こまかいこと。` matches thin/fine sense. | Pass |

## Judgment Calls

- `発` is listed in the Tanos N4 source, but the existing N5 overlay already contains very basic kanji and may later expand; kept because the user asked to avoid only existing N5/N1 rows and `発` is foundational for N4 vocabulary.
- `地` is a high-frequency basic kanji and appears in the Tanos N4 source; kept even though some modern learner lists sometimes place it earlier.
- `党` was selected early for N2 because it is high-frequency in public/political vocabulary and present in the Tanos N2 source, even though the Japanese gloss uses `政党`, an N2/N3-adjacent compound.
- `校` was treated as a tempting N4 core candidate because it is foundational in beginner learning, but it was not present in the local Tanos N4 PDF/deck source. Replaced with Tanos-confirmed N4 `館`.
- `悪` was initially considered for N3 because it is high-frequency, but the local Tanos source places it in N4. Replaced with Tanos-confirmed N3 `性`.

## Substitutions

- Excluded existing N5/N1 collisions from candidate consideration, including `会`, `作`, `使`, `安`, `少`, `口`, `古`, `多`, `手`, `店`, and N1 entries such as `乏`, `伴`, `偽`, `債`.
- Post-write validation substituted `校` with `館` because `校` was not present in the local Tanos N4 PDF/deck source used for this batch.
- Post-write validation substituted N3 `悪` with `性` because `悪` appears in the local Tanos N4 source, not the N3 source.
