# Corpus Level Review

This file tracks level questions found during manual review of generated L1
support overlays. Corrections with high confidence should be encoded in
`jlpt-overrides.jsonl`; uncertain items stay here until reviewed.

## 2026-04-30 manual review

High-confidence corrections applied:

- N2 -> N5: `お休み`, `お帰り`, `お邪魔します`, `お願いします`
- N2 -> N4: `お代わり`, `お先に`, `お参り`, `お大事に`, `お待たせしました`, `お手伝いさん`
- N4 -> N5: `うん`, `そう`, `ケーキ`, `サラダ`, `サンドイッチ`, `テニス`, `ピアノ`, `プレゼント`, `ベル`

Still needs review:

- Current N5 but possibly N4: `勤める`, `交番`, `並ぶ`, `並べる`, `上げる`, `問題`, `喫茶店`
- Current N4 but possibly N5/N3 depending list: `パソコン`, `スーツ`, `スーツケース`, `アルバイト`, `ガソリンスタンド`, `アクセサリー`
- Current N2 loanwords may be corpus-list artifacts rather than pedagogical N2: `ガム`, `クーラー`, `コンセント`, `サラリーマン`, `スタート`, `ストップ`, `スリッパ`

Policy:

- Do not trust imported `jlpt_level` blindly when creating learner-facing overlays.
- Move only high-confidence corrections into `jlpt-overrides.jsonl`.
- Keep uncertain items listed here until checked against the project's chosen
  JLPT reference source or manually accepted by review.
