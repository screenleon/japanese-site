# JS-042 N3 mental_model rollout — codex self-review (2026-05-08)

## Summary

40 N3 entries were annotated. Confidence split: 31 high / 9 medium / 0 low. Human review should focus first on the 9 medium rows, mostly entries whose existing `explanation_ja` is concise and leaves less room to confirm nuance boundaries.

## Per-entry review

| slug | confidence | concept-fit-note | human-review-priority |
|---|---|---|---|
| bakari | high | anchored to "同じものや動作が多すぎる" and "Vたばかり" | no |
| concession-temo | high | anchored to "Aが起きても、Bは変わらない" | no |
| conditional-ba | high | anchored to "論理的な条件" and command/request caution | no |
| conditional-nara | high | anchored to "相手の話や前に出た条件" | no |
| conditional-tara | high | anchored to "仮定、時間の順序、過去の発見" | no |
| contrast-noni | high | anchored to "期待と違ってB" and speaker emotion | no |
| dake-denaku | high | anchored to "一つだけではなく、ほかにも" | no |
| dokoroka | medium | anchored to "大きく超える、または反対の状態" | yes |
| hazuda | high | anchored to "理由や情報から考えて、当然そう" | no |
| hazuganai | high | anchored to "その可能性はないと強く言う" | no |
| hodo | medium | anchored to "程度の高さや比較の基準" | yes |
| kadouka | high | anchored to "そうであるかそうでないか分からない内容" | no |
| kagiri | medium | anchored to "状態や条件が続く範囲内" | yes |
| kamoshirenai | high | anchored to "可能性があるが、はっきり分からない" | no |
| kawari-ni | medium | anchored to "代替、または交換条件" | yes |
| koto-ni-natte-iru | high | anchored to "規則、予定、取り決め" | no |
| kotoni-naru | high | anchored to "自分だけの意志ではなく、予定や決まり" | no |
| kotoni-suru | high | anchored to "自分の意志で決めたこと" | no |
| mitai | high | anchored to "見たことや聞いたことから判断した推量、または比喩" | no |
| monono | medium | anchored to "事実として認める内容" plus contrast | yes |
| ni-chigainai | medium | anchored to "根拠に基づいて強く確信" | yes |
| okage-de | high | anchored to "よい結果をもたらした原因や助け" | no |
| rashii | high | anchored to "聞いた情報に基づく推量" and "そのものらしい性質" | no |
| sae-ba | high | anchored to "その一つの条件が満たされれば十分" | no |
| sei-de | high | anchored to "望ましくない結果の原因" | no |
| souda-appearance | high | anchored to "見た感じから判断" | no |
| souda-hearsay | high | anchored to "人から聞いた情報やニュース" | no |
| tabi-ni | high | anchored to "起こるごとに、別の出来事も繰り返し" | no |
| tameni-purpose | high | anchored to "目的を表します" and volitional action | no |
| tara-dou | high | anchored to "提案や助言" and softer than command | no |
| te-bakari-iru | high | anchored to "同じ行為を何度も続けたり、それだけをして" | no |
| teiku-tekuru | medium | anchored to "今から先" versus "過去から今まで" movement/change | yes |
| teshimau | high | anchored to "全部終わる" plus regret/unintentional feeling | no |
| to-sureba | medium | anchored to "仮定したうえで、その結果や判断" | yes |
| tokoro | high | anchored to "動作の段階" and attached form | no |
| wake-niwa-ikanai | medium | anchored to "社会的、道徳的、または状況上の理由" | yes |
| youda | high | anchored to "状況から判断した推量や、何かに似ている" | no |
| youni-goal | high | anchored to "できる状態を目標にする" | no |
| youni-naru | high | anchored to "状態に変わる" | no |
| youni-suru | high | anchored to "努力したり、気をつけたり" | no |

## Sample-5 OK re-check

- concession-temo: re-checked against entry's explanation_ja — held
- conditional-nara: re-checked against entry's explanation_ja — held
- kawari-ni: re-checked against entry's explanation_ja — held
- tokoro: re-checked against entry's explanation_ja — held
- youni-suru: re-checked against entry's explanation_ja — held
- Conclusion: held / 0 revisions made before declaring done

## Native review pass — 2026-05-08

Second-pass review by Claude main thread acting as a native Japanese reader, focused on the 9 medium-confidence rows above. Reassessment + revision outcomes:

| slug | codex confidence | native verdict | action |
|---|---|---|---|
| dokoroka | medium | OK but emotionally flat — missed the surprise/emphatic-contrast charge that defines どころか | revised: re-anchored to 「強い驚きとともに」 and the "足場にせず一気に飛び越える" image |
| hodo | medium | minor inaccuracy: 「程度の**高さ**」 — ほど is a neutral degree yardstick, not necessarily "high" | revised: changed to 「程度をはかる**ものさし**」; clarified proportional usage as "一方が動けばもう一方も同じだけ動く比例関係" |
| kawari-ni | medium | latter half 「一方を受けて反対側の条件を出す」 reads stiff/translation-tang | revised: reframed with native collocations 「代償として」「引き受けと引き換えの関係」 |
| teiku-tekuru | medium | restated explanation_ja; the load-bearing N3 hook (speaker's vantage point) was buried after「または」 | revised: promoted 「話し手の立つ位置を基準に」 to the leading definition; all usages now derive from that vantage |
| ni-chigainai | medium | concept-fit good; only collocation polish: 「状況を**合わせる**」 unnatural | revised: 「合わせる」→「踏まえる」; 「可能性が」→「可能性は」 (topic continuity with の他) |
| kagiri | medium | concept-fit excellent; "範囲を区切って" + "話し手が保証できる範囲" matches Japanese pedagogical framing | unchanged — should have been high |
| monono | medium | concept-fit excellent; 「Aを否定せずに置きながら、結果や評価は別方向だと示します」 captures the bookish concession nuance precisely | unchanged — should have been high |
| to-sureba | medium | concept-fit excellent; "事実かどうかよりも、その条件を受けた論理の流れに焦点があります" sharply distinguishes from たら/ば/なら | unchanged — should have been high |
| wake-niwa-ikanai | medium | concept-fit excellent; 「能力の問題ではなく」 + 「したいが事情が許さない」 is the load-bearing distinction | unchanged — should have been high |

### Net outcome of native pass
- 5 entries revised (dokoroka, hodo, kawari-ni, teiku-tekuru, ni-chigainai); dual-write byte-identical invariant re-verified across all 40.
- 4 entries (kagiri, monono, to-sureba, wake-niwa-ikanai) retained unchanged but reassessed as high-confidence — codex's medium tagging on these was conservative.
- Final confidence read after native pass: 35 high / 5 revised-to-native / 0 low remaining.
