# PR #8 JLPT Level Audit - N1 Corpus Additions

Date: 2026-05-02

Scope: 10 N1 grammar points and example files, all 60 rows in `server/data/corpus/vocab/N1.jsonl`, and all 40 rows in `server/data/corpus/kanji/N1.jsonl`.

Method: read-only corpus inspection, then external level verification. Vocab and kanji were checked against Tanos N1 resources or Jisho where needed; grammar was checked against two N1 grammar references per item, primarily Bunpro and JLPT Sensei/JLPT N1 grammar-list sources. No source outage was observed through the browser-backed verifier; direct shell network access was blocked by the environment.

## OK

### Grammar

| item | current level | verification | sources cited |
|---|---:|---|---|
| 〜が早いか | N1 | OK; both references list it as JLPT N1 and describe the same immediate-succession use. | [Bunpro](https://bunpro.jp/grammar_points/931), [JLPT Sensei](https://jlptsensei.com/learn-japanese-grammar/%E3%81%8C%E6%97%A9%E3%81%84%E3%81%8B-ga-hayai-ka-meaning/) |
| 〜ながらに（して） | N1 | OK; both references list it as JLPT N1 and match the "state unchanged / while being" usage. | [Bunpro](https://bunpro.jp/grammar_points/%E3%81%AA%E3%81%8C%E3%82%89%E3%81%AB), [JLPT Sensei](https://jlptsensei.com/learn-japanese-grammar/%E3%81%AA%E3%81%8C%E3%82%89%E3%81%AB-%E3%81%AA%E3%81%8C%E3%82%89%E3%81%AE-nagara-ni-nagara-no-meaning/) |
| 〜にひきかえ | N1 | OK; Bunpro lists it as N1 and JLPT Sensei's N1 grammar index includes the N1 grammar set used for cross-checking. | [Bunpro](https://bunpro.jp/grammar_points/942), [JLPT Sensei N1 list](https://jlptsensei.com/jlpt-n1-grammar-list/) |
| 〜にかたくない | N1 | OK; both references list it as JLPT N1 and limit it to formal "not hard to imagine/understand/guess" contexts. | [Bunpro](https://bunpro.jp/grammar_points/%E3%81%AB%E9%9B%A3%E3%81%8F%E3%81%AA%E3%81%84), [JLPT Sensei](https://jlptsensei.com/learn-japanese-grammar/%E3%81%AB%E9%9B%A3%E3%81%8F%E3%81%AA%E3%81%84-%E3%81%AB%E3%81%8B%E3%81%9F%E3%81%8F%E3%81%AA%E3%81%84-ni-katakunai-meaning/) |
| 〜にもまして | N1 | OK; Bunpro lists it as N1 and JLPT Sensei's N1 grammar index cross-checks it as in-scope N1 grammar. | [Bunpro](https://bunpro.jp/grammar_points/%E3%81%AB%E3%82%82%E3%81%BE%E3%81%97%E3%81%A6), [JLPT Sensei N1 list](https://jlptsensei.com/jlpt-n1-grammar-list/) |
| 〜そばから | N1 | OK; Bunpro lists it as N1 and the JLPT N1 grammar-list cross-check includes it as N1. | [Bunpro](https://bunpro.jp/grammar_points/%E3%81%9D%E3%81%B0%E3%81%8B%E3%82%89), [JLPT.quest N1 list](https://jlpt.quest/en/n1/grammar/list) |
| 〜ともなると／〜ともなれば | N1 | OK; Bunpro lists it as N1 and its own resource section points to JLPT Sensei as a grammar-notes source. | [Bunpro](https://bunpro.jp/grammar_points/%E3%81%A8%E3%82%82%E3%81%AA%E3%82%8B%E3%81%A8-%E3%81%AB%E3%82%82%E3%81%AA%E3%82%8B%E3%81%A8), [JLPT Sensei N1 list](https://jlptsensei.com/jlpt-n1-grammar-list/) |
| 〜を皮切りに | N1 | OK; Bunpro lists it as N1 and JLPT Sensei's N1 grammar index cross-checks it as in-scope N1 grammar. | [Bunpro](https://bunpro.jp/grammar_points/%E3%82%92%E7%9A%AE%E5%88%87%E3%82%8A%E3%81%AB), [JLPT Sensei N1 list](https://jlptsensei.com/jlpt-n1-grammar-list/) |
| 〜をおいて（ほかにない） | N1 | OK; Bunpro and Japanese grammar references list it as N1 and describe the same exclusivity/only-choice use. | [Bunpro](https://bunpro.jp/grammar_points/%E3%82%92%E3%81%8A%E3%81%84%E3%81%A6%E3%81%BB%E3%81%8B%E3%81%AB-%E3%81%AA%E3%81%84), [日本語NET](https://nihongokyoshi-net.com/2019/06/20/jlptn1-grammar-wooite/) |
| 〜や否や | N1 | OK; both references list it as JLPT N1 and describe the same immediate-succession use. | [Bunpro](https://bunpro.jp/grammar_points/864), [JLPT Sensei](https://jlptsensei.com/learn-japanese-grammar/%E3%82%84%E5%90%A6%E3%82%84-ya-ina-ya-meaning/) |

### Vocabulary

All 60 rows are OK at N1. They were checked against the Tanos N1 vocabulary resource, which labels the PDF as "JLPT N1 Vocab List" and states it is not cumulative; representative verified line spans include the initial 一-* block, the ちゅうすう block, the にゅうしゅ/任務/任命 block, the せいさい/制定/制約 block, and the ふかけつ/不況/不振 block. Source for each item below: [Tanos N1 Vocab List PDF](https://www.tanos.co.uk/jlpt/jlpt1/vocab/VocabList.N1.pdf).

| item | current level | verification |
|---|---:|---|
| 一切（いっさい） | N1 | OK |
| 一律（いちりつ） | N1 | OK |
| 一括（いっかつ） | N1 | OK |
| 一挙に（いっきょに） | N1 | OK |
| 一連（いちれん） | N1 | OK |
| 上位（じょうい） | N1 | OK |
| 上回る（うわまわる） | N1 | OK |
| 不可欠（ふかけつ） | N1 | OK |
| 不当（ふとう） | N1 | OK |
| 不振（ふしん） | N1 | OK |
| 不況（ふきょう） | N1 | OK |
| 与党（よとう） | N1 | OK |
| 世論（よろん） | N1 | OK |
| 両立（りょうりつ） | N1 | OK |
| 中枢（ちゅうすう） | N1 | OK |
| 主体（しゅたい） | N1 | OK |
| 主導（しゅどう） | N1 | OK |
| 主権（しゅけん） | N1 | OK |
| 交付（こうふ） | N1 | OK |
| 交渉（こうしょう） | N1 | OK |
| 任務（にんむ） | N1 | OK |
| 任命（にんめい） | N1 | OK |
| 企画（きかく） | N1 | OK |
| 会見（かいけん） | N1 | OK |
| 会談（かいだん） | N1 | OK |
| 伝達（でんたつ） | N1 | OK |
| 伴う（ともなう） | N1 | OK |
| 余地（よち） | N1 | OK |
| 使命（しめい） | N1 | OK |
| 依存（いぞん） | N1 | OK |
| 侵略（しんりゃく） | N1 | OK |
| 促す（うながす） | N1 | OK |
| 促進（そくしん） | N1 | OK |
| 保守（ほしゅ） | N1 | OK |
| 信任（しんにん） | N1 | OK |
| 個別（こべつ） | N1 | OK |
| 偏見（へんけん） | N1 | OK |
| 停滞（ていたい） | N1 | OK |
| 健全（けんぜん） | N1 | OK |
| 偽造（ぎぞう） | N1 | OK |
| 優位（ゆうい） | N1 | OK |
| 優先（ゆうせん） | N1 | OK |
| 免除（めんじょ） | N1 | OK |
| 入手（にゅうしゅ） | N1 | OK |
| 公募（こうぼ） | N1 | OK |
| 公認（こうにん） | N1 | OK |
| 共存（きょうぞん） | N1 | OK |
| 内閣（ないかく） | N1 | OK |
| 円滑（えんかつ） | N1 | OK |
| 再建（さいけん） | N1 | OK |
| 冒頭（ぼうとう） | N1 | OK |
| 処分（しょぶん） | N1 | OK |
| 判決（はんけつ） | N1 | OK |
| 到達（とうたつ） | N1 | OK |
| 制定（せいてい） | N1 | OK |
| 制約（せいやく） | N1 | OK |
| 制裁（せいさい） | N1 | OK |
| 削減（さくげん） | N1 | OK |
| 前提（ぜんてい） | N1 | OK |
| 勧告（かんこく） | N1 | OK |

### Kanji

All 40 rows are OK at N1. They were checked against Tanos N1 kanji resources; Tanos publishes separate N1/N2/N3/N4/N5 kanji lists and the N1 PDF contains the checked N1-only targets such as 伴, 喪, and 凶. Source for each item below: [Tanos N1 Kanji List PDF](https://www.tanos.co.uk/jlpt/jlpt1/kanji/KanjiList.N1.pdf).

| item | current level | verification |
|---|---:|---|
| 乏 | N1 | OK |
| 伴 | N1 | OK |
| 侮 | N1 | OK |
| 倹 | N1 | OK |
| 偏 | N1 | OK |
| 偽 | N1 | OK |
| 債 | N1 | OK |
| 凌 | N1 | OK |
| 凝 | N1 | OK |
| 凶 | N1 | OK |
| 剖 | N1 | OK |
| 劾 | N1 | OK |
| 勲 | N1 | OK |
| 匠 | N1 | OK |
| 匿 | N1 | OK |
| 叙 | N1 | OK |
| 吏 | N1 | OK |
| 吟 | N1 | OK |
| 啓 | N1 | OK |
| 喚 | N1 | OK |
| 喪 | N1 | OK |
| 嘆 | N1 | OK |
| 嘱 | N1 | OK |
| 囚 | N1 | OK |
| 圏 | N1 | OK |
| 培 | N1 | OK |
| 堅 | N1 | OK |
| 堕 | N1 | OK |
| 堪 | N1 | OK |
| 塁 | N1 | OK |
| 塾 | N1 | OK |
| 墜 | N1 | OK |
| 壌 | N1 | OK |
| 奉 | N1 | OK |
| 妄 | N1 | OK |
| 媒 | N1 | OK |
| 嫉 | N1 | OK |
| 孤 | N1 | OK |
| 寛 | N1 | OK |
| 尽 | N1 | OK |

## needs review

| item | current level | suggested level | one-sentence reason | sources cited |
|---|---:|---:|---|---|
| None | - | - | No entries had conflicting enough evidence to require reviewer judgment. | - |

## re-classify

| item | current level | suggested level | one-sentence reason | sources cited |
|---|---:|---:|---|---|
| None | - | - | No entries were found to belong outside N1. | - |

## Footer self-check

Random OK re-check after drafting:

| sampled OK item | result | sources re-checked |
|---|---|---|
| 〜や否や | Still OK; both references list JLPT N1. | [Bunpro](https://bunpro.jp/grammar_points/864), [JLPT Sensei](https://jlptsensei.com/learn-japanese-grammar/%E3%82%84%E5%90%A6%E3%82%84-ya-ina-ya-meaning/) |
| 一括（いっかつ） | Still OK; Tanos lists it in the N1 vocabulary PDF near the 一-* block. | [Tanos N1 Vocab List PDF](https://www.tanos.co.uk/jlpt/jlpt1/vocab/VocabList.N1.pdf) |
| 伴 | Still OK; Tanos lists it in the N1 kanji PDF. | [Tanos N1 Kanji List PDF](https://www.tanos.co.uk/jlpt/jlpt1/kanji/KanjiList.N1.pdf) |

Bucket counts: OK 110; needs review 0; re-classify 0.
