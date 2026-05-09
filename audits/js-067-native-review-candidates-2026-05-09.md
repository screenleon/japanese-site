# JS-067 Native-Reviewer Candidate Sweep — 2026-05-09

## Scope
- Live corpus scanned: 200 entries (N1×40 + N2×40 + N3×40 + N4×40 + N5×40). The task brief says 199 entries; the repository currently contains 200 JSON files, so this audit follows the live source of truth.
- 全 1421 條 key_terms, plus 76 title_ja furigana pairs (1497 total furigana pairs).
- 全 63 個 title_ja 含漢字者.
- Read-only corpus sweep over server/data/corpus/grammar/N1..N5/*.json; output-only write to audits/js-067-native-review-candidates-2026-05-09.md.

## Methodology
- R1 single-kanji-token: one bare-kanji furigana token, requiring context verification. Tokens that already include okurigana, e.g. `伴って`, are not treated as R1.
- R2 stem-suspect-ending: bare one-kanji token with a stem-like reading ending and no corresponding okurigana in the token.
- R2-confirmed-truncation: source text contains the same bare one-kanji token followed by likely okurigana; this is HIGH because it resembles the F1 stem-cut bugs.
- R3 onyomi-kunyomi-ambiguity: common bare one-kanji tokens with multiple frequent readings.
- R4 proper-noun: location/institution/person-name watchlist hits.
- R5 compound-tokenizer-split: bare one-kanji token whose kanji also appears inside a multi-kanji compound in the same entry.
- R6 rendaku-or-special-reading: fixed irregular, jukujikun, rendaku, or grammar-title watchlist hits.
- R7 number-counter-context: number/counter/date/time context requiring native verification.
- R8 verb-form-mismatch: bare one-kanji token appears against a fuller verb/adjective/conjugated lexeme in source.
- R-title-uncovered: title_ja contains kanji not covered by annotations.furigana.title_ja or key_terms.
- Confidence: HIGH = likely content bug; MEDIUM = native check needed; LOW = small risk or overinclusive watchlist.
- `feedback_native_perspective.md` was requested but is not present in this checkout; this file uses the JS-042 two-pass audit shape and leaves Section 2 for human native review.

### F1 Fixed Calibration Cases
| entry | field | current | calibration-note |
| --- | --- | --- | --- |
| N1/wo-kawakiri-ni.json | title_ja | 皮切り/かわきり | fixed full lexeme; not flagged as R2 |
| N2/ni-sakidatte.json | title_ja | 先立って/さきだって | fixed full lexeme; not flagged as R2 |
| N2/nitomonatte.json | title_ja | 伴って/ともなって | fixed full lexeme; not flagged as R2 |
| N1/taru-mono.json | key_terms | 心構え/こころがまえ | fixed full lexeme; not flagged as R2 |

### PR-Gate Critic Baseline Check
| candidate | status | note |
| --- | --- | --- |
| N1/atte-no 支/ささ | H001 | Caught as R1 + R2-confirmed + R8. |
| N1/ya-inaya title 否/いな (brief says to-ina-ya) | M022 | Caught as R1 + R6; requires independent native check. |
| N1/atte-no 他/た | STALE IN HEAD | No 他/た pair exists in N1/atte-no.json at current HEAD. The same risk shape is caught at M021 for N1/wo-oite title_ja 他/た. |

## Section 1 — Codex Pre-Pass Findings (Auto-Detected Risks)

### 1.1 HIGH — Confirmed-truncation candidates (R2-confirmed)
| id | level | entry | field | key_term | observed | likely-correct | rationale | codex-suggested-correction | confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| H001 | HIGH | N1/atte-no.json | key_terms | 支/ささ | 〜あっての 意味核：ある存在や支えがあるからこそ後件が成り立つことを表す。感謝・前提条件・依存関係を強く示す。 接 | 支えがあるからこそ | R1+R2-confirmed-truncation+R8 | 採用 = 支え/ささえ，已修正 commit 65c336f | HIGH-likely |
| H002 | HIGH | N1/bakoso.json | key_terms | 厳/きび | こそ」を使います。 - 信頼していればこそ、厳しく言う - 大切であればこそ、慎重に扱う - 親友であればこそ、遠慮なく言える | 厳しく | R1+R2-confirmed-truncation+R8 | 採用 = 厳しく/きびしく，已修正 commit 65c336f | HIGH-likely |
| H003 | HIGH | N1/ga-hayai-ka.json | title_ja | 早/はや | 〜が早いか 意味核：ある動作が終わるのとほぼ同時に、次の動作がただちに起こることを表す。反 | 早いか | R1+R2-confirmed-truncation+R8 | 採用 = 早いか/はやいか，已修正 commit 65c336f | HIGH-likely |
| H004 | HIGH | N1/gotoki.json | key_terms | 使/つか | のようだ」を硬く、文語的に言う表現です。 使い分け： - ごとき＋名詞：夢のごとき時間 - ごとく＋動詞・形容詞：風のごとく走る | 使い | R1+R2-confirmed-truncation+R8 | 採用 = 使い/つかい，已修正 commit 65c336f | HIGH-likely |
| H005 | HIGH | N1/gotoki.json | key_terms | 分/わ | うだ」を硬く、文語的に言う表現です。 使い分け： - ごとき＋名詞：夢のごとき時間 - ごとく＋動詞・形容詞：風のごとく走る - | 分け | R1+R2-confirmed-truncation+R8 | 採用 = 分け/わけ，已修正 commit 65c336f | HIGH-likely |
| H006 | HIGH | N1/gotoki.json | key_terms | 走/はし | き時間 - ごとく＋動詞・形容詞：風のごとく走る - ごとし：文末で「〜のようだ」 古風で硬い響きがあるため、日常会話では「よう | 走る | R1+R2-confirmed-truncation+R8 | 採用 = 走る/はしる，已修正 commit 65c336f | HIGH-likely |
| H007 | HIGH | N1/kagiri-da.json | title_ja | 限/かぎ | 〜限りだ 意味核：感情や評価の程度が極めて高いことを表す。話し手の主観的な強い気持ちを、 | 限りだ | R1+R2-confirmed-truncation+R8 | 採用 = 限りだ/かぎりだ，已修正 commit 65c336f | HIGH-likely |
| H008 | HIGH | N1/te-yamanai.json | key_terms | 惜/お | やまない 「〜てやまない」は、願う・愛する・惜しむ・期待するなどの感情や評価が非常に強く、ずっと続いていることを表す硬い表現である | 惜しむ | R1+R2-confirmed-truncation+R8+R5 | 採用 = 惜しむ/おしむ，已修正 commit 65c336f | HIGH-likely |
| H009 | HIGH | N1/te-yamanai.json | key_terms | 愛/あい | 〜てやまない 「〜てやまない」は、願う・愛する・惜しむ・期待するなどの感情や評価が非常に強く、ずっと続いていることを表す硬い表 | 愛する | R1+R2-confirmed-truncation+R8+R5 | 採用 = 愛する/あいする，已修正 commit 65c336f | HIGH-likely |
| H010 | HIGH | N1/te-yamanai.json | key_terms | 願/ねが | 〜てやまない 「〜てやまない」は、願う・愛する・惜しむ・期待するなどの感情や評価が非常に強く、ずっと続いていることを表す | 願う | R1+R2-confirmed-truncation+R8 | 採用 = 願う/ねがう，已修正 commit 65c336f | HIGH-likely |
| H011 | HIGH | N1/towaie.json | key_terms | 寒/さむ | な形容詞にも使えます。 - 春とはいえ、まだ寒い - 安いとはいえ、品質は悪くない - 専門家とはいえ、すべてを知っているわけでは | 寒い | R1+R2-confirmed-truncation+R8 | 採用 = 寒い/さむい，已修正 commit 65c336f | HIGH-likely |
| H012 | HIGH | N1/zujimai.json | key_terms | 聞/き | である。「手紙を書かずじまいだった」「理由を聞けずじまいに終わった」のように、未完了への残念さや心残りを伴いやすい。接続は動詞ない | 聞けずじまいに | R1+R2-confirmed-truncation+R8 | 採用 = 聞けずじまいに/きけずじまいに，已修正 commit 65c336f | HIGH-likely |
| H013 | HIGH | N1/zuniwa-irarenai.json | key_terms | 笑/わら | て自然にそうしてしまう感じがあります。 - 笑わずにはいられない - 心配せずにはいられない 「〜ないわけにはいかない」は社会的 | 笑わずにはいられない | R1+R2-confirmed-truncation+R8 | 採用 = 笑わずにはいられない/わらわずにはいられない，已修正 commit 65c336f | HIGH-likely |
| H014 | HIGH | N2/ageku.json | key_terms | 至/いた | ：長い迷い・努力・経緯の末に、よくない結果に至ることを表す。過程の負担や徒労感が含まれる。 接続：動詞た形＋あげく（に）。名詞＋ | 至ることを | R1+R2-confirmed-truncation+R8 | 採用 = 至ることを/いたることを，已修正 commit 65c336f | HIGH-likely |
| H015 | HIGH | N2/dokoro-dewa-nai.json | key_terms | 厳/きび | 〜どころではない 意味核：状況が厳しく、前件をする余裕がない、または前件の程度では表せないほどであることを表す。 接 | 厳しく | R1+R2-confirmed-truncation+R8 | 採用 = 厳しく/きびしく，已修正 commit 65c336f | HIGH-likely |
| H016 | HIGH | N2/nai-koto-wa-nai.json | key_terms | 控/ひか | 全に否定はしないが積極的にも認めない、という控えめな肯定を表す。例として「読めないことはない」「行けないことはないが時間がかかる」 | 控えめな | R1+R2-confirmed-truncation+R8 | 採用 = 控えめな/ひかえめな，已修正 commit 65c336f | HIGH-likely |
| H017 | HIGH | N2/ni-hanshite.json | title_ja | 反/はん | 〜に反して 「〜に反して」は、予想・期待・規則などと反対の結果や行動になることを表す。例と | 反して | R1+R2-confirmed-truncation+R8+R5 | 採用 = 反して/はんして，已修正 commit 65c336f | HIGH-likely |
| H018 | HIGH | N2/ni-kagiri.json | title_ja | 限/かぎ | に限り 「Nに限り」は、対象をその範囲だけに限定する表現です。案内や規則でよく使われます | 限り | R1+R2-confirmed-truncation+R8+R5 | 採用 = 限り/かぎり，已修正 commit 65c336f | HIGH-likely |
| H019 | HIGH | N2/ni-saishite.json | title_ja | 際/さい | 〜に際して 「〜に際して」は、特別な出来事や公式な行動を行う時に、という意味を表す。例とし | 際して | R1+R2-confirmed-truncation+R8 | 採用 = 際して/さいして，已修正 commit 65c336f | HIGH-likely |
| H020 | HIGH | N2/ni-shitagatte.json | key_terms | 合/あ | AにしたがってB」は、Aが変化すると、それに合わせてBも変化することを表します。 - 年を取るにしたがって、体力が落ちます - | 合わせて | R1+R2-confirmed-truncation+R8 | 採用 = 合わせて/あわせて，已修正 commit 65c336f | HIGH-likely |
| H021 | HIGH | N2/ni-tsurete.json | key_terms | 変/か | につれて 「AにつれてB」は、Aが変わると一緒にBも変わることを表します。「にしたがって」と近いですが、自然に少しずつ変 | 変わると | R1+R2-confirmed-truncation+R8 | 採用 = 変わると/かわると，已修正 commit 65c336f | HIGH-likely |
| H022 | HIGH | N2/nikanshite.json | title_ja | 関/かん | 〜について／〜に関して 「〜について」と「〜に関して」は、どちらも「〜をテーマにして」という意味です。 | 関して | R1+R2-confirmed-truncation+R8 | 採用 = 関して/かんして，已修正 commit 65c336f | HIGH-likely |
| H023 | HIGH | N2/ta-totan.json | key_terms | 話/はな | がったとたんめまいがした」のように、後件には話し手が意図しない出来事が来る。自分の意志で計画して行う後件には使いにくい。 | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H024 | HIGH | N2/warini.json | key_terms | 忙/いそが | の実際の程度が合わないことを表す。例として「忙しいわりに元気だ」「値段のわりに品質がいい」のように、性質・量・値段などと結果を比べ | 忙しいわりに | R1+R2-confirmed-truncation+R8 | 採用 = 忙しいわりに/いそがしいわりに，已修正 commit 65c336f | HIGH-likely |
| H025 | HIGH | N2/wo-tooshite.json | title_ja | 通/とお | を通して 「Nを通して」は、手段・媒介、または期間全体を表します。 - 友人を通して知 | 通して | R1+R2-confirmed-truncation+R8 | 採用 = 通して/とおして，已修正 commit 65c336f | HIGH-likely |
| H026 | HIGH | N3/contrast-noni.json | key_terms | 話/はな | に、期待と違ってBだ」という逆接を表します。話し手の不満、意外さ、残念な気持ちが出やすい表現です。 形：普通形＋のに。な形容詞と | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H027 | HIGH | N3/hazuganai.json | key_terms | 話/はな | です。「こんなに簡単なはずがない」のように、話し手の確信を表します。 - 彼がうそをつくはずがありません - この時間に店が開い | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H028 | HIGH | N3/kagiri.json | title_ja | 限/かぎ | 限り 「限り」は、ある状態や条件が続く範囲内で物事が成り立つことを表す表現である。動詞 | 限り | R1+R2-confirmed-truncation+R8 | 採用 = 限り/かぎり，已修正 commit 65c336f | HIGH-likely |
| H029 | HIGH | N3/koto-ni-natte-iru.json | key_terms | 取/と | ている 「ことになっている」は、規則、予定、取り決めなどによって、ある行為や状態が決まっていることを表す表現である。動詞の辞書形ま | 取り | R1+R2-confirmed-truncation+R8 | 採用 = 取り/とり，已修正 commit 65c336f | HIGH-likely |
| H030 | HIGH | N3/ni-chigainai.json | key_terms | 話/はな | し、「必ずそうだ」「間違いなくそうだ」という話し手の強い結論を示す。 | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H031 | HIGH | N3/ni-chigainai.json | title_ja | 違/ちが | に違いない 「に違いない」は、根拠に基づいて強く確信している判断を表す表現である。普通形 | 違いない | R1+R2-confirmed-truncation+R8+R5 | 採用 = 違いない/ちがいない，已修正 commit 65c336f | HIGH-likely |
| H032 | HIGH | N3/okage-de.json | key_terms | 助/たす | 「おかげで」は、よい結果をもたらした原因や助けを表す表現である。普通形または「名詞 + の」に接続し、感謝やよい評価を伴って用い | 助けを | R1+R2-confirmed-truncation+R8 | 採用 = 助けを/たすけを，已修正 commit 65c336f | HIGH-likely |
| H033 | HIGH | N3/sei-de.json | key_terms | 話/はな | ある。普通形または「名詞 + の」に接続し、話し手がその原因を否定的に捉えたり、責任を置いたりする場合に用いられる。 | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H034 | HIGH | N3/tabi-ni.json | key_terms | 繰/く | 」は、ある出来事が起こるごとに、別の出来事も繰り返し起こることを表す表現である。動詞の辞書形または「名詞 + の」に接続し、習慣的 | 繰り | R1+R2-confirmed-truncation+R8 | 採用 = 繰り/くり，已修正 commit 65c336f | HIGH-likely |
| H035 | HIGH | N3/tabi-ni.json | key_terms | 返/かえ | 、ある出来事が起こるごとに、別の出来事も繰り返し起こることを表す表現である。動詞の辞書形または「名詞 + の」に接続し、習慣的・規 | 返し | R1+R2-confirmed-truncation+R8 | 採用 = 返し/かえし，已修正 commit 65c336f | HIGH-likely |
| H036 | HIGH | N3/te-bakari-iru.json | key_terms | 話/はな | 表す表現である。動詞のて形に接続し、しばしば話し手の不満や批判の気持ちを伴う。 | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H037 | HIGH | N3/teiku-tekuru.json | key_terms | 話/はな | す。「Vてくる」は、過去から今までの変化や、話し手の方へ近づく動きを表します。 - これから寒くなっていきます - 少しずつ日本 | 話し | R1+R2-confirmed-truncation+R8 | 採用 = 話し/はなし，已修正 commit 65c336f | HIGH-likely |
| H038 | HIGH | N5/counter-basic.json | title_ja | 人/にん | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は | 人います | R1+R2-confirmed-truncation+R8+R3+R5+R6+R7 | 採用 = 人います/にんいます，已修正 commit 65c336f | HIGH-likely |

### 1.2 HIGH — On/Kun ambiguity needing context (R3 with both readings plausible)
| id | level | entry | field | key_term | reading | alt-reading | context-hint | codex-suggested-correction | confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| H038 | HIGH | N5/counter-basic.json | title_ja | 人 | にん | ひと/じん/にん | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は | 採用 = 人います/にんいます，已修正 commit 65c336f | HIGH-likely |
| M007 | MEDIUM | N1/ja-arumai-shi.json | key_terms | 子 | こ | native-check | いし、または「わけじゃあるまいし」。 例：子どもじゃあるまいし、そのくらい自分で判断しなさい。／初めて会う人じゃあるまいし、そん | 建議改為 子ども/こども（source: 子どもじゃあるまいし） | HIGH-likely |
| M019 | MEDIUM | N1/wo-kawakiri-ni.json | key_terms | 後 | あと | native-check | 味核：ある出来事を最初のきっかけとして、その後に同種の出来事が連続して広がることを表す。単なる開始時点ではなく、後続の展開が予想さ | 建議改為 その後/そのあと（source: その後に） | MEDIUM-likely |
| M021 | MEDIUM | N1/wo-oite.json | title_ja | 他 | た | ほか | 〜をおいて（他にない） 意味核：ある役割・条件に最もふさわしいものはそれ以外にない、と強く限定する | 建議改為 他にない/ほかにない（source: 他にない） | HIGH-likely |
| M028 | MEDIUM | N2/ta-totan.json | key_terms | 手 | て | native-check | たとたんめまいがした」のように、後件には話し手が意図しない出来事が来る。自分の意志で計画して行う後件には使いにくい。 | 建議改為 話し手/はなして（source: 話し手） | HIGH-likely |
| M030 | MEDIUM | N2/ue-de.json | title_ja | 上 | うえ | native-check | 上で 「Vた上で」は、先に必要なことをしてから次のことをする意味です。「Nの上で」は、 | 建議改為 上で/うえで（source: Vた上で） | HIGH-likely |
| M031 | MEDIUM | N2/ue-wa.json | title_ja | 上 | うえ | native-check | 〜上は 意味核：前件が決まった以上、後件は当然そうするべきだ、または避けられないことを表 | 建議改為 上は/うえは（source: 〜上は） | HIGH-likely |
| M037 | MEDIUM | N3/contrast-noni.json | key_terms | 手 | て | native-check | 期待と違ってBだ」という逆接を表します。話し手の不満、意外さ、残念な気持ちが出やすい表現です。 形：普通形＋のに。な形容詞と名詞 | 建議改為 話し手/はなして（source: 話し手の不満） | HIGH-likely |
| M040 | MEDIUM | N3/hazuganai.json | key_terms | 手 | て | native-check | 。「こんなに簡単なはずがない」のように、話し手の確信を表します。 - 彼がうそをつくはずがありません - この時間に店が開いてい | 建議改為 話し手/はなして（source: 話し手の確信） | HIGH-likely |
| M053 | MEDIUM | N3/ni-chigainai.json | key_terms | 手 | て | native-check | 「必ずそうだ」「間違いなくそうだ」という話し手の強い結論を示す。 | 建議改為 話し手/はなして（source: 話し手の強い結論） | HIGH-likely |
| M056 | MEDIUM | N3/sei-de.json | key_terms | 手 | て | native-check | 。普通形または「名詞 + の」に接続し、話し手がその原因を否定的に捉えたり、責任を置いたりする場合に用いられる。 | 建議改為 話し手/はなして（source: 話し手がその原因） | HIGH-likely |
| M061 | MEDIUM | N3/te-bakari-iru.json | key_terms | 手 | て | native-check | 表現である。動詞のて形に接続し、しばしば話し手の不満や批判の気持ちを伴う。 | 建議改為 話し手/はなして（source: 話し手の不満） | HIGH-likely |
| M063 | MEDIUM | N3/teiku-tekuru.json | key_terms | 手 | て | native-check | 「Vてくる」は、過去から今までの変化や、話し手の方へ近づく動きを表します。 - これから寒くなっていきます - 少しずつ日本語が | 建議改為 話し手/はなして（source: 話し手の方へ） | HIGH-likely |
| M064 | MEDIUM | N3/tokoro.json | key_terms | 中 | ちゅう | native-check | 帰ったところです 前につく形で、直前・進行中・直後が変わります。 | 建議改為 進行中/しんこうちゅう（source: 進行中） | HIGH-likely |
| M071 | MEDIUM | N3/youni-suru.json | key_terms | 心 | こころ | native-check | 日読むようにしています」のように、続けている心がけにも使います。 - 忘れないようにします - 毎朝早く起きるようにしています | 建議改為 心がけ/こころがけ（source: 心がけ） | HIGH-likely |
| M077 | MEDIUM | N5/number-counter-hon.json | title_ja | 本 | ほん | native-check | 〜本 「本」は、ペン、瓶、木、傘など、細長い物を数える助数詞である。数字によって読み方が | NEEDS-NATIVE: 候選A=本/ほん for bare counter heading; 候選B=ぽん/ぼん after specific numerals | NEEDS-NATIVE |

### 1.3 MEDIUM — Single-kanji tokens (R1)
| id | level | entry | field | key_term | reading | source-context-snippet | codex-suggested-correction | confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| M001 | MEDIUM | N1/bekarazu.json | key_terms | 不 | ふ | 〜べからず／〜べからざる 意味核：禁止・不許可・強い否定的義務を表す。掲示・規則・標語などで「してはならない」と硬く示す。 | 建議改為 不許可/ふきょか（source: 不許可） | HIGH-likely |
| M002 | MEDIUM | N1/de-are-de-are.json | key_terms | 差 | さ | らの場合でも後件が同じであることを表す。条件差を超えた原則を示す。 接続：名詞A＋であれ＋名詞B＋であれ。な形容詞語幹A＋であれ | 建議改為 条件差/じょうけんさ（source: 条件差） | HIGH-likely |
| M003 | MEDIUM | N1/ga-hayai-ka.json | key_terms | 差 | さ | 時性を強く示し、前件と後件の間にほとんど時間差がない。 接続：形：動詞辞書形＋が早いか。前件には瞬間的に完了する動作を置き、後件 | 建議改為 時間差/じかんさ（source: 時間差） | HIGH-likely |
| M004 | MEDIUM | N1/ga-hayai-ka.json | key_terms | 速 | はや | 、次の動作がただちに起こることを表す。反応の速さ、切迫感、機械的な即時性を強く示し、前件と後件の間にほとんど時間差がない。 接続 | 建議改為 速さ/はやさ（source: 反応の速さ） | HIGH-likely |
| M005 | MEDIUM | N1/gotoki.json | key_terms | 夢 | ゆめ | 表現です。 使い分け： - ごとき＋名詞：夢のごとき時間 - ごとく＋動詞・形容詞：風のごとく走る - ごとし：文末で「〜のよう | 保持現狀正確: 夢/ゆめ（source: 夢のごとき時間） | HIGH-likely |
| M006 | MEDIUM | N1/gotoki.json | key_terms | 風 | かぜ | ：夢のごとき時間 - ごとく＋動詞・形容詞：風のごとく走る - ごとし：文末で「〜のようだ」 古風で硬い響きがあるため、日常会話 | 保持現狀正確: 風/かぜ（source: 風のごとく走る） | HIGH-likely |
| M007 | MEDIUM | N1/ja-arumai-shi.json | key_terms | 子 | こ | いし、または「わけじゃあるまいし」。 例：子どもじゃあるまいし、そのくらい自分で判断しなさい。／初めて会う人じゃあるまいし、そん | 建議改為 子ども/こども（source: 子どもじゃあるまいし） | HIGH-likely |
| M008 | MEDIUM | N1/ka-ina-ka.json | title_ja | 否 | いな | 〜か否か 意味核：ある事柄が成立するかしないかを、硬く二択として示す。判断・検討・可否の対 | 建議改為 否か/いなか（source: 〜か否か） | HIGH-likely |
| M009 | MEDIUM | N1/katagata.json | key_terms | 礼 | れい | に、別の目的も兼ねて行うことを表す。挨拶・お礼・報告など改まった訪問や連絡で多い。 接続：名詞＋かたがた。する名詞は「する」を省 | 建議改為 お礼/おれい（source: お礼・報告） | HIGH-likely |
| M010 | MEDIUM | N1/kirai-ga-aru.json | key_terms | 癖 | くせ | 「〜きらいがある」は、よくない傾向や心配な癖がある、という意味です。人や組織の性質を少し批判的に言うときに使います。 形：動詞 | 保持現狀正確: 癖/くせ（source: 心配な癖） | HIGH-likely |
| M011 | MEDIUM | N1/nagara-ni.json | key_terms | 昔 | むかし | 成立することを表す。特に「生まれながらに」「昔ながらに」のように、自然な属性や昔からの状態がそのまま続くことを述べる。 接続：形 | 建議改為 昔ながらに/むかしながらに（source: 昔ながらに） | HIGH-likely |
| M012 | MEDIUM | N1/nagara-ni.json | key_terms | 涙 | なみだ | 慣用的には「生まれながらに」「昔ながらに」「涙ながらに」などが多い。 例：生まれながらに才能を備える／昔ながらに製法を守る／涙な | 建議改為 涙ながらに/なみだながらに（source: 涙ながらに） | HIGH-likely |
| M013 | MEDIUM | N1/nimokakawarazu.json | key_terms | 書 | か | 意味は近いですが、「〜にもかかわらず」はより書き言葉で、客観的に対比を述べる感じがあります。 | 建議改為 書き言葉/かきことば（source: 書き言葉） | HIGH-likely |
| M014 | MEDIUM | N1/nimokakawarazu.json | key_terms | 雨 | あめ | と名詞は「である」を使う形が硬いです。 - 雨が降っているにもかかわらず - 危険であるにもかかわらず - 休日であるにもかかわら | 保持現狀正確: 雨/あめ（source: 雨が降っている） | HIGH-likely |
| M015 | MEDIUM | N1/taru-mono.json | key_terms | 姿 | すがた | 多くの場合、後ろには義務・心構え・あるべき姿が続きます。少し古風で、演説・論説・強い主張に向いています。 | 保持現狀正確: 姿/すがた（source: あるべき姿） | HIGH-likely |
| M016 | MEDIUM | N1/toittemo-kagonai.json | key_terms | 町 | まち | 医学を変えたといっても過言ではない - この町は観光で成り立っているといっても過言ではない 強く評価するときに使います。話し言葉 | 保持現狀正確: 町/まち（source: この町） | HIGH-likely |
| M017 | MEDIUM | N1/tomonaku.json | key_terms | 音 | おと | 疑問詞とともに用いる場合は「どこからともなく音が聞こえる」「誰ともなく拍手が起こる」のように、出どころや主体がはっきりしないことを | 保持現狀正確: 音/おと（source: 音が聞こえる） | HIGH-likely |
| M018 | MEDIUM | N1/towaie.json | key_terms | 春 | はる | とはいえ。名詞やな形容詞にも使えます。 - 春とはいえ、まだ寒い - 安いとはいえ、品質は悪くない - 専門家とはいえ、すべてを知 | 保持現狀正確: 春/はる（source: 春とはいえ） | HIGH-likely |
| M019 | MEDIUM | N1/wo-kawakiri-ni.json | key_terms | 後 | あと | 味核：ある出来事を最初のきっかけとして、その後に同種の出来事が連続して広がることを表す。単なる開始時点ではなく、後続の展開が予想さ | 建議改為 その後/そのあと（source: その後に） | MEDIUM-likely |
| M020 | MEDIUM | N1/wo-oite.json | key_terms | 彼 | かれ | として評価する名詞を取る。 例：この任務は彼をおいてほかにない／開催地は京都をおいて他にない／解決策は対話をおいてほかにない | 保持現狀正確: 彼/かれ（source: この任務は彼をおいて） | HIGH-likely |
| M021 | MEDIUM | N1/wo-oite.json | title_ja | 他 | た | 〜をおいて（他にない） 意味核：ある役割・条件に最もふさわしいものはそれ以外にない、と強く限定する | 建議改為 他にない/ほかにない（source: 他にない） | HIGH-likely |
| M022 | MEDIUM | N1/ya-inaya.json | title_ja | 否 | いな | 〜や否や 意味核：前の動作が成立した直後に、次の出来事が間を置かず起こることを表す。前件を | 建議改為 否や/いなや（source: 〜や否や） | HIGH-likely |
| M023 | MEDIUM | N2/monoda.json | key_terms | 驚 | おどろ | ろ、よく公園で遊んだものだ 3. 強い感心や驚き：よく我慢したものだ 形は普通形＋ものだです。な形容詞は「静かなものだ」のように | 建議改為 驚き/おどろき（source: 驚き） | HIGH-likely |
| M024 | MEDIUM | N2/monono-formal.json | key_terms | 異 | こと | として認めながら、そこから予想される結果とは異なる後件を述べる。硬めの逆接表現。 接続：動詞普通形＋ものの。い形容詞＋ものの。な | 建議改為 異なる/ことなる（source: 異なる後件） | HIGH-likely |
| M025 | MEDIUM | N2/monono-formal.json | key_terms | 認 | みと | 〜ものの 意味核：前件を事実として認めながら、そこから予想される結果とは異なる後件を述べる。硬めの逆接表現。 接続：動 | 建議改為 認めながら/みとめながら（source: 認めながら） | HIGH-likely |
| M026 | MEDIUM | N2/ni-motozuite.json | title_ja | 基 | もと | 〜に基づいて 「〜に基づいて」は、資料・事実・規則などを根拠として行動や判断をすることを表 | 建議改為 基づいて/もとづいて（source: 〜に基づいて） | HIGH-likely |
| M027 | MEDIUM | N2/ni-oujite.json | title_ja | 応 | おう | 〜に応じて 「〜に応じて」は、相手・状況・量などに合わせて対応や内容が変わることを表す。例 | 建議改為 応じて/おうじて（source: 〜に応じて） | HIGH-likely |
| M028 | MEDIUM | N2/ta-totan.json | key_terms | 手 | て | たとたんめまいがした」のように、後件には話し手が意図しない出来事が来る。自分の意志で計画して行う後件には使いにくい。 | 建議改為 話し手/はなして（source: 話し手） | HIGH-likely |
| M029 | MEDIUM | N2/ta-totan.json | key_terms | 雨 | あめ | が起きることを表す。例として「外に出たとたん雨が降り出した」「立ち上がったとたんめまいがした」のように、後件には話し手が意図しない | 保持現狀正確: 雨/あめ（source: 雨が降り出した） | HIGH-likely |
| M030 | MEDIUM | N2/ue-de.json | title_ja | 上 | うえ | 上で 「Vた上で」は、先に必要なことをしてから次のことをする意味です。「Nの上で」は、 | 建議改為 上で/うえで（source: Vた上で） | HIGH-likely |
| M031 | MEDIUM | N2/ue-wa.json | title_ja | 上 | うえ | 〜上は 意味核：前件が決まった以上、後件は当然そうするべきだ、または避けられないことを表 | 建議改為 上は/うえは（source: 〜上は） | HIGH-likely |
| M032 | MEDIUM | N2/wo-chuushin-ni.json | key_terms | 駅 | えき | して範囲や活動が広がることを表す。例として「駅を中心に店が集まっている」「若者を中心に人気がある」のように使う。中心点を示し、その | 保持現狀正確: 駅/えき（source: 駅を中心に） | HIGH-likely |
| M033 | MEDIUM | N3/conditional-ba.json | key_terms | 五 | ご | 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲めば、書く→書けば - い形容詞：高い→高ければ - 否定：〜なけれ | 建議改為 五段/ごだん（source: 五段動詞） | HIGH-likely |
| M034 | MEDIUM | N3/conditional-ba.json | key_terms | 段 | だん | いう論理的な条件を表します。 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲めば、書く→書けば - い形容詞：高い | NEEDS-NATIVE: 候選A=一段/いちだん（source first hit）; 候選B=五段/ごだん（same entry also lists 五段） | NEEDS-NATIVE |
| M035 | MEDIUM | N3/conditional-ba.json | key_terms | 的 | てき | Aという条件が成り立つとBになる、という論理的な条件を表します。 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲め | 建議改為 論理的/ろんりてき（source: 論理的な条件） | HIGH-likely |
| M036 | MEDIUM | N3/contrast-noni.json | key_terms | 形 | がた | 意外さ、残念な気持ちが出やすい表現です。 形：普通形＋のに。な形容詞と名詞は「な」を入れます。 - 勉強したのに - 高いのに | 建議改為 普通形/ふつうけい（source: 普通形＋のに） | HIGH-likely |
| M037 | MEDIUM | N3/contrast-noni.json | key_terms | 手 | て | 期待と違ってBだ」という逆接を表します。話し手の不満、意外さ、残念な気持ちが出やすい表現です。 形：普通形＋のに。な形容詞と名詞 | 建議改為 話し手/はなして（source: 話し手の不満） | HIGH-likely |
| M038 | MEDIUM | N3/dokoroka.json | key_terms | 形 | がた | は反対の状態であることを表す表現である。普通形や名詞に接続し、「それどころではなく」という強い対比を示す。 | 建議改為 普通形/ふつうけい（source: 普通形や名詞） | HIGH-likely |
| M039 | MEDIUM | N3/hazuganai.json | key_terms | 性 | せい | 「はずがない」は、理由から考えて、その可能性はないと強く言う表現です。「こんなに簡単なはずがない」のように、話し手の確信を表しま | 建議改為 可能性/かのうせい（source: 可能性はない） | HIGH-likely |
| M040 | MEDIUM | N3/hazuganai.json | key_terms | 手 | て | 。「こんなに簡単なはずがない」のように、話し手の確信を表します。 - 彼がうそをつくはずがありません - この時間に店が開いてい | 建議改為 話し手/はなして（source: 話し手の確信） | HIGH-likely |
| M041 | MEDIUM | N3/hodo.json | key_terms | 形 | がた | の基準を表す表現である。名詞または動詞の辞書形に接続し、「それくらい」「その程度まで」という意味を表す。また、「動詞ば形 + 動詞 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M042 | MEDIUM | N3/kagiri.json | key_terms | 形 | がた | 事が成り立つことを表す表現である。動詞の辞書形、ている形、ない形、または「名詞 + の」に接続し、「その条件が続いている間は」「そ | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M043 | MEDIUM | N3/kamoshirenai.json | key_terms | 形 | がた | るが、はっきり分からないことを表します。普通形の後に続き、丁寧に言う時は「かもしれません」になります。 - 明日は雨が降るかもし | 建議改為 普通形/ふつうけい（source: 普通形の後） | HIGH-likely |
| M044 | MEDIUM | N3/kamoshirenai.json | key_terms | 性 | せい | かもしれない 「かもしれない」は、可能性があるが、はっきり分からないことを表します。普通形の後に続き、丁寧に言う時は「かもし | 建議改為 可能性/かのうせい（source: 可能性がある） | HIGH-likely |
| M045 | MEDIUM | N3/kawari-ni.json | key_terms | 形 | がた | 表現である。「名詞 + の」または動詞の辞書形に接続し、「その代わりとして」「その反面」という意味で用いられる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M046 | MEDIUM | N3/koto-ni-natte-iru.json | key_terms | 決 | き | る 「ことになっている」は、規則、予定、取り決めなどによって、ある行為や状態が決まっていることを表す表現である。動詞の辞書形または | 建議改為 取り決め/とりきめ（source: 取り決め） | HIGH-likely |
| M047 | MEDIUM | N3/kotoni-naru.json | key_terms | 形 | がた | ことになる 「V辞書形/ない形 + ことになる」は、自分だけの意志ではなく、予定や決まりとしてそう決まった | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M048 | MEDIUM | N3/kotoni-naru.json | key_terms | 決 | き | とになる」は、自分だけの意志ではなく、予定や決まりとしてそう決まったことを表します。結果としてその状態になる場合にも使います。 | 建議改為 決まり/きまり（source: 予定や決まり） | HIGH-likely |
| M049 | MEDIUM | N3/kotoni-suru.json | key_terms | 形 | がた | ことにする 「V辞書形/ない形 + ことにする」は、自分の意志で決めたことを表します。習慣として決めている | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M050 | MEDIUM | N3/monono.json | key_terms | 形 | がた | る結果や対照的な内容を続ける表現である。普通形に接続し、書き言葉で用いられることが多い。 | 建議改為 普通形/ふつうけい（source: 普通形に接続） | HIGH-likely |
| M051 | MEDIUM | N3/monono.json | key_terms | 書 | か | 的な内容を続ける表現である。普通形に接続し、書き言葉で用いられることが多い。 | 建議改為 書き言葉/かきことば（source: 書き言葉） | HIGH-likely |
| M052 | MEDIUM | N3/ni-chigainai.json | key_terms | 形 | がた | て強く確信している判断を表す表現である。普通形に接続し、「必ずそうだ」「間違いなくそうだ」という話し手の強い結論を示す。 | 建議改為 普通形/ふつうけい（source: 普通形に接続） | HIGH-likely |
| M053 | MEDIUM | N3/ni-chigainai.json | key_terms | 手 | て | 「必ずそうだ」「間違いなくそうだ」という話し手の強い結論を示す。 | 建議改為 話し手/はなして（source: 話し手の強い結論） | HIGH-likely |
| M054 | MEDIUM | N3/okage-de.json | key_terms | 形 | がた | をもたらした原因や助けを表す表現である。普通形または「名詞 + の」に接続し、感謝やよい評価を伴って用いられる。 | 建議改為 普通形/ふつうけい（source: 普通形または） | HIGH-likely |
| M055 | MEDIUM | N3/sei-de.json | key_terms | 形 | がた | 望ましくない結果の原因を表す表現である。普通形または「名詞 + の」に接続し、話し手がその原因を否定的に捉えたり、責任を置いたりす | 建議改為 普通形/ふつうけい（source: 普通形または） | HIGH-likely |
| M056 | MEDIUM | N3/sei-de.json | key_terms | 手 | て | 。普通形または「名詞 + の」に接続し、話し手がその原因を否定的に捉えたり、責任を置いたりする場合に用いられる。 | 建議改為 話し手/はなして（source: 話し手がその原因） | HIGH-likely |
| M057 | MEDIUM | N3/souda-hearsay.json | key_terms | 形 | がた | そうだ（伝聞） 「普通形 + そうだ」は、人から聞いた情報やニュースなどを伝える表現です。 - 明日は雨が | 建議改為 普通形/ふつうけい（source: 普通形 + そうだ） | HIGH-likely |
| M058 | MEDIUM | N3/tabi-ni.json | key_terms | 形 | がた | り返し起こることを表す表現である。動詞の辞書形または「名詞 + の」に接続し、習慣的・規則的な反復を述べる際に用いられる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M059 | MEDIUM | N3/tabi-ni.json | key_terms | 的 | てき | の辞書形または「名詞 + の」に接続し、習慣的・規則的な反復を述べる際に用いられる。 | 建議改為 習慣的/しゅうかんてき（source: 習慣的・規則的） | HIGH-likely |
| M060 | MEDIUM | N3/tameni-purpose.json | key_terms | 形 | がた | ために（目的） 「V辞書形 + ために」「Nのために」は、目的を表します。意志を持って行う行動によく使います。 | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M061 | MEDIUM | N3/te-bakari-iru.json | key_terms | 手 | て | 表現である。動詞のて形に接続し、しばしば話し手の不満や批判の気持ちを伴う。 | 建議改為 話し手/はなして（source: 話し手の不満） | HIGH-likely |
| M062 | MEDIUM | N3/teiku-tekuru.json | key_terms | 動 | うご | く・てくる 「Vていく」は、今から先へ変化や動きが続くことを表します。「Vてくる」は、過去から今までの変化や、話し手の方へ近づく動 | 建議改為 動き/うごき（source: 動きが続く） | HIGH-likely |
| M063 | MEDIUM | N3/teiku-tekuru.json | key_terms | 手 | て | 「Vてくる」は、過去から今までの変化や、話し手の方へ近づく動きを表します。 - これから寒くなっていきます - 少しずつ日本語が | 建議改為 話し手/はなして（source: 話し手の方へ） | HIGH-likely |
| M064 | MEDIUM | N3/tokoro.json | key_terms | 中 | ちゅう | 帰ったところです 前につく形で、直前・進行中・直後が変わります。 | 建議改為 進行中/しんこうちゅう（source: 進行中） | HIGH-likely |
| M065 | MEDIUM | N3/wake-niwa-ikanai.json | key_terms | 形 | がた | てはいけないことを表す表現である。動詞の辞書形に接続し、単なる能力ではなく、事情による不可を述べる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M066 | MEDIUM | N3/wake-niwa-ikanai.json | key_terms | 的 | てき | けにはいかない 「わけにはいかない」は、社会的、道徳的、または状況上の理由により、ある行為をすることができない、またはしてはいけな | NEEDS-NATIVE: 候選A=社会的/しゃかいてき; 候選B=道徳的/どうとくてき（both in source） | NEEDS-NATIVE |
| M067 | MEDIUM | N3/youda.json | key_terms | 書 | か | 誰か来たようです - 夢のようです 少し書き言葉寄りで、「みたいだ」より丁寧・硬めです。 | 建議改為 書き言葉/かきことば（source: 書き言葉寄り） | HIGH-likely |
| M068 | MEDIUM | N3/youni-goal.json | key_terms | 形 | がた | ように（目標・変化） 「V辞書形/ない形 + ように」は、できる状態を目標にする表現です。意志で直接コントロールしに | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M069 | MEDIUM | N3/youni-naru.json | key_terms | 形 | がた | ようになる 「V辞書形 + ようになる」は、前はできなかったことができる状態に変わることを表します。「Vな | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M070 | MEDIUM | N3/youni-suru.json | key_terms | 形 | がた | ようにする 「V辞書形/ない形 + ようにする」は、ある状態になるように努力したり、気をつけたりすることを | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M071 | MEDIUM | N3/youni-suru.json | key_terms | 心 | こころ | 日読むようにしています」のように、続けている心がけにも使います。 - 忘れないようにします - 毎朝早く起きるようにしています | 建議改為 心がけ/こころがけ（source: 心がけ） | HIGH-likely |
| M072 | MEDIUM | N5/counter-basic.json | title_ja | 枚 | まい | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は一般 | 保持現狀正確: 枚/まい（source: 数え方 title + 一枚） | HIGH-likely |
| M073 | MEDIUM | N5/kore-sore-are.json | title_ja | 系 | けい | こ・そ・あ系 「こ・そ・あ・ど」は、話し手や聞き手からの距離を表す指示語の体系である。「これ」は | 保持現狀正確: 系/けい（source: こ・そ・あ系） | HIGH-likely |
| M074 | MEDIUM | N5/masu-form.json | title_ja | 形 | けい | ます形 「ます形」は、動詞をていねいに言う形です。 今のことや未来のことを言うときは「〜 | 建議改為 ます形/ますけい（source: ます形） | HIGH-likely |
| M075 | MEDIUM | N5/nani-doko-dare.json | title_ja | 何 | なに | 何・どこ・誰 基本的な疑問詞には、「何」「どこ」「誰」「いつ」「どう」などがある。「何 | 保持現狀正確: 何/なに（title/context is bare 疑問詞; counters may read なん separately） | MEDIUM-likely |
| M076 | MEDIUM | N5/nani-doko-dare.json | title_ja | 誰 | だれ | 何・どこ・誰 基本的な疑問詞には、「何」「どこ」「誰」「いつ」「どう」などがある。「何」は物や内 | 保持現狀正確: 誰/だれ（source: 誰） | HIGH-likely |
| M077 | MEDIUM | N5/number-counter-hon.json | title_ja | 本 | ほん | 〜本 「本」は、ペン、瓶、木、傘など、細長い物を数える助数詞である。数字によって読み方が | NEEDS-NATIVE: 候選A=本/ほん for bare counter heading; 候選B=ぽん/ぼん after specific numerals | NEEDS-NATIVE |
| M078 | MEDIUM | N5/te-form.json | title_ja | 形 | けい | て形 「て形」は、動詞をほかの表現につなげるための大事な形です。 使い方： - 動作を | 建議改為 て形/てけい（source: て形） | HIGH-likely |
| M079 | MEDIUM | N5/verb-te-form-connection.json | title_ja | 形 | けい | て形の接続 動詞のて形は、複数の動作を順番に述べるときに使う。「V1て、V2」は「V1を | 建議改為 て形/てけい（source: て形の接続） | HIGH-likely |

### 1.4 MEDIUM — Stem-suspect (R2 not auto-confirmed)
| id | level | entry | field | key_term | reading | source-context-snippet | codex-suggested-correction | confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| M009 | MEDIUM | N1/katagata.json | key_terms | 礼 | れい | に、別の目的も兼ねて行うことを表す。挨拶・お礼・報告など改まった訪問や連絡で多い。 接続：名詞＋かたがた。する名詞は「する」を省 | 建議改為 お礼/おれい（source: お礼・報告） | HIGH-likely |
| M018 | MEDIUM | N1/towaie.json | key_terms | 春 | はる | とはいえ。名詞やな形容詞にも使えます。 - 春とはいえ、まだ寒い - 安いとはいえ、品質は悪くない - 専門家とはいえ、すべてを知 | 保持現狀正確: 春/はる（source: 春とはいえ） | HIGH-likely |
| M027 | MEDIUM | N2/ni-oujite.json | title_ja | 応 | おう | 〜に応じて 「〜に応じて」は、相手・状況・量などに合わせて対応や内容が変わることを表す。例 | 建議改為 応じて/おうじて（source: 〜に応じて） | HIGH-likely |
| M030 | MEDIUM | N2/ue-de.json | title_ja | 上 | うえ | 上で 「Vた上で」は、先に必要なことをしてから次のことをする意味です。「Nの上で」は、 | 建議改為 上で/うえで（source: Vた上で） | HIGH-likely |
| M031 | MEDIUM | N2/ue-wa.json | title_ja | 上 | うえ | 〜上は 意味核：前件が決まった以上、後件は当然そうするべきだ、または避けられないことを表 | 建議改為 上は/うえは（source: 〜上は） | HIGH-likely |
| M039 | MEDIUM | N3/hazuganai.json | key_terms | 性 | せい | 「はずがない」は、理由から考えて、その可能性はないと強く言う表現です。「こんなに簡単なはずがない」のように、話し手の確信を表しま | 建議改為 可能性/かのうせい（source: 可能性はない） | HIGH-likely |
| M044 | MEDIUM | N3/kamoshirenai.json | key_terms | 性 | せい | かもしれない 「かもしれない」は、可能性があるが、はっきり分からないことを表します。普通形の後に続き、丁寧に言う時は「かもし | 建議改為 可能性/かのうせい（source: 可能性がある） | HIGH-likely |
| M064 | MEDIUM | N3/tokoro.json | key_terms | 中 | ちゅう | 帰ったところです 前につく形で、直前・進行中・直後が変わります。 | 建議改為 進行中/しんこうちゅう（source: 進行中） | HIGH-likely |
| M072 | MEDIUM | N5/counter-basic.json | title_ja | 枚 | まい | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は一般 | 保持現狀正確: 枚/まい（source: 数え方 title + 一枚） | HIGH-likely |
| M073 | MEDIUM | N5/kore-sore-are.json | title_ja | 系 | けい | こ・そ・あ系 「こ・そ・あ・ど」は、話し手や聞き手からの距離を表す指示語の体系である。「これ」は | 保持現狀正確: 系/けい（source: こ・そ・あ系） | HIGH-likely |
| M074 | MEDIUM | N5/masu-form.json | title_ja | 形 | けい | ます形 「ます形」は、動詞をていねいに言う形です。 今のことや未来のことを言うときは「〜 | 建議改為 ます形/ますけい（source: ます形） | HIGH-likely |
| M078 | MEDIUM | N5/te-form.json | title_ja | 形 | けい | て形 「て形」は、動詞をほかの表現につなげるための大事な形です。 使い方： - 動作を | 建議改為 て形/てけい（source: て形） | HIGH-likely |
| M079 | MEDIUM | N5/verb-te-form-connection.json | title_ja | 形 | けい | て形の接続 動詞のて形は、複数の動作を順番に述べるときに使う。「V1て、V2」は「V1を | 建議改為 て形/てけい（source: て形の接続） | HIGH-likely |

### 1.5 MEDIUM — Compound-tokenizer-split (R5)
| id | level | entry | field | key_term | reading | source-context-snippet | codex-suggested-correction | confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| M001 | MEDIUM | N1/bekarazu.json | key_terms | 不 | ふ | 〜べからず／〜べからざる 意味核：禁止・不許可・強い否定的義務を表す。掲示・規則・標語などで「してはならない」と硬く示す。 | 建議改為 不許可/ふきょか（source: 不許可） | HIGH-likely |
| M002 | MEDIUM | N1/de-are-de-are.json | key_terms | 差 | さ | らの場合でも後件が同じであることを表す。条件差を超えた原則を示す。 接続：名詞A＋であれ＋名詞B＋であれ。な形容詞語幹A＋であれ | 建議改為 条件差/じょうけんさ（source: 条件差） | HIGH-likely |
| M003 | MEDIUM | N1/ga-hayai-ka.json | key_terms | 差 | さ | 時性を強く示し、前件と後件の間にほとんど時間差がない。 接続：形：動詞辞書形＋が早いか。前件には瞬間的に完了する動作を置き、後件 | 建議改為 時間差/じかんさ（source: 時間差） | HIGH-likely |
| M004 | MEDIUM | N1/ga-hayai-ka.json | key_terms | 速 | はや | 、次の動作がただちに起こることを表す。反応の速さ、切迫感、機械的な即時性を強く示し、前件と後件の間にほとんど時間差がない。 接続 | 建議改為 速さ/はやさ（source: 反応の速さ） | HIGH-likely |
| M006 | MEDIUM | N1/gotoki.json | key_terms | 風 | かぜ | ：夢のごとき時間 - ごとく＋動詞・形容詞：風のごとく走る - ごとし：文末で「〜のようだ」 古風で硬い響きがあるため、日常会話 | 保持現狀正確: 風/かぜ（source: 風のごとく走る） | HIGH-likely |
| M008 | MEDIUM | N1/ka-ina-ka.json | title_ja | 否 | いな | 〜か否か 意味核：ある事柄が成立するかしないかを、硬く二択として示す。判断・検討・可否の対 | 建議改為 否か/いなか（source: 〜か否か） | HIGH-likely |
| M009 | MEDIUM | N1/katagata.json | key_terms | 礼 | れい | に、別の目的も兼ねて行うことを表す。挨拶・お礼・報告など改まった訪問や連絡で多い。 接続：名詞＋かたがた。する名詞は「する」を省 | 建議改為 お礼/おれい（source: お礼・報告） | HIGH-likely |
| M019 | MEDIUM | N1/wo-kawakiri-ni.json | key_terms | 後 | あと | 味核：ある出来事を最初のきっかけとして、その後に同種の出来事が連続して広がることを表す。単なる開始時点ではなく、後続の展開が予想さ | 建議改為 その後/そのあと（source: その後に） | MEDIUM-likely |
| M027 | MEDIUM | N2/ni-oujite.json | title_ja | 応 | おう | 〜に応じて 「〜に応じて」は、相手・状況・量などに合わせて対応や内容が変わることを表す。例 | 建議改為 応じて/おうじて（source: 〜に応じて） | HIGH-likely |
| M031 | MEDIUM | N2/ue-wa.json | title_ja | 上 | うえ | 〜上は 意味核：前件が決まった以上、後件は当然そうするべきだ、または避けられないことを表 | 建議改為 上は/うえは（source: 〜上は） | HIGH-likely |
| M033 | MEDIUM | N3/conditional-ba.json | key_terms | 五 | ご | 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲めば、書く→書けば - い形容詞：高い→高ければ - 否定：〜なけれ | 建議改為 五段/ごだん（source: 五段動詞） | HIGH-likely |
| M034 | MEDIUM | N3/conditional-ba.json | key_terms | 段 | だん | いう論理的な条件を表します。 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲めば、書く→書けば - い形容詞：高い | NEEDS-NATIVE: 候選A=一段/いちだん（source first hit）; 候選B=五段/ごだん（same entry also lists 五段） | NEEDS-NATIVE |
| M035 | MEDIUM | N3/conditional-ba.json | key_terms | 的 | てき | Aという条件が成り立つとBになる、という論理的な条件を表します。 形： - 一段動詞：食べる→食べれば - 五段動詞：飲む→飲め | 建議改為 論理的/ろんりてき（source: 論理的な条件） | HIGH-likely |
| M036 | MEDIUM | N3/contrast-noni.json | key_terms | 形 | がた | 意外さ、残念な気持ちが出やすい表現です。 形：普通形＋のに。な形容詞と名詞は「な」を入れます。 - 勉強したのに - 高いのに | 建議改為 普通形/ふつうけい（source: 普通形＋のに） | HIGH-likely |
| M038 | MEDIUM | N3/dokoroka.json | key_terms | 形 | がた | は反対の状態であることを表す表現である。普通形や名詞に接続し、「それどころではなく」という強い対比を示す。 | 建議改為 普通形/ふつうけい（source: 普通形や名詞） | HIGH-likely |
| M039 | MEDIUM | N3/hazuganai.json | key_terms | 性 | せい | 「はずがない」は、理由から考えて、その可能性はないと強く言う表現です。「こんなに簡単なはずがない」のように、話し手の確信を表しま | 建議改為 可能性/かのうせい（source: 可能性はない） | HIGH-likely |
| M041 | MEDIUM | N3/hodo.json | key_terms | 形 | がた | の基準を表す表現である。名詞または動詞の辞書形に接続し、「それくらい」「その程度まで」という意味を表す。また、「動詞ば形 + 動詞 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M042 | MEDIUM | N3/kagiri.json | key_terms | 形 | がた | 事が成り立つことを表す表現である。動詞の辞書形、ている形、ない形、または「名詞 + の」に接続し、「その条件が続いている間は」「そ | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M043 | MEDIUM | N3/kamoshirenai.json | key_terms | 形 | がた | るが、はっきり分からないことを表します。普通形の後に続き、丁寧に言う時は「かもしれません」になります。 - 明日は雨が降るかもし | 建議改為 普通形/ふつうけい（source: 普通形の後） | HIGH-likely |
| M044 | MEDIUM | N3/kamoshirenai.json | key_terms | 性 | せい | かもしれない 「かもしれない」は、可能性があるが、はっきり分からないことを表します。普通形の後に続き、丁寧に言う時は「かもし | 建議改為 可能性/かのうせい（source: 可能性がある） | HIGH-likely |
| M045 | MEDIUM | N3/kawari-ni.json | key_terms | 形 | がた | 表現である。「名詞 + の」または動詞の辞書形に接続し、「その代わりとして」「その反面」という意味で用いられる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M046 | MEDIUM | N3/koto-ni-natte-iru.json | key_terms | 決 | き | る 「ことになっている」は、規則、予定、取り決めなどによって、ある行為や状態が決まっていることを表す表現である。動詞の辞書形または | 建議改為 取り決め/とりきめ（source: 取り決め） | HIGH-likely |
| M047 | MEDIUM | N3/kotoni-naru.json | key_terms | 形 | がた | ことになる 「V辞書形/ない形 + ことになる」は、自分だけの意志ではなく、予定や決まりとしてそう決まった | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M049 | MEDIUM | N3/kotoni-suru.json | key_terms | 形 | がた | ことにする 「V辞書形/ない形 + ことにする」は、自分の意志で決めたことを表します。習慣として決めている | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M050 | MEDIUM | N3/monono.json | key_terms | 形 | がた | る結果や対照的な内容を続ける表現である。普通形に接続し、書き言葉で用いられることが多い。 | 建議改為 普通形/ふつうけい（source: 普通形に接続） | HIGH-likely |
| M052 | MEDIUM | N3/ni-chigainai.json | key_terms | 形 | がた | て強く確信している判断を表す表現である。普通形に接続し、「必ずそうだ」「間違いなくそうだ」という話し手の強い結論を示す。 | 建議改為 普通形/ふつうけい（source: 普通形に接続） | HIGH-likely |
| M054 | MEDIUM | N3/okage-de.json | key_terms | 形 | がた | をもたらした原因や助けを表す表現である。普通形または「名詞 + の」に接続し、感謝やよい評価を伴って用いられる。 | 建議改為 普通形/ふつうけい（source: 普通形または） | HIGH-likely |
| M055 | MEDIUM | N3/sei-de.json | key_terms | 形 | がた | 望ましくない結果の原因を表す表現である。普通形または「名詞 + の」に接続し、話し手がその原因を否定的に捉えたり、責任を置いたりす | 建議改為 普通形/ふつうけい（source: 普通形または） | HIGH-likely |
| M057 | MEDIUM | N3/souda-hearsay.json | key_terms | 形 | がた | そうだ（伝聞） 「普通形 + そうだ」は、人から聞いた情報やニュースなどを伝える表現です。 - 明日は雨が | 建議改為 普通形/ふつうけい（source: 普通形 + そうだ） | HIGH-likely |
| M058 | MEDIUM | N3/tabi-ni.json | key_terms | 形 | がた | り返し起こることを表す表現である。動詞の辞書形または「名詞 + の」に接続し、習慣的・規則的な反復を述べる際に用いられる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M059 | MEDIUM | N3/tabi-ni.json | key_terms | 的 | てき | の辞書形または「名詞 + の」に接続し、習慣的・規則的な反復を述べる際に用いられる。 | 建議改為 習慣的/しゅうかんてき（source: 習慣的・規則的） | HIGH-likely |
| M060 | MEDIUM | N3/tameni-purpose.json | key_terms | 形 | がた | ために（目的） 「V辞書形 + ために」「Nのために」は、目的を表します。意志を持って行う行動によく使います。 | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M064 | MEDIUM | N3/tokoro.json | key_terms | 中 | ちゅう | 帰ったところです 前につく形で、直前・進行中・直後が変わります。 | 建議改為 進行中/しんこうちゅう（source: 進行中） | HIGH-likely |
| M065 | MEDIUM | N3/wake-niwa-ikanai.json | key_terms | 形 | がた | てはいけないことを表す表現である。動詞の辞書形に接続し、単なる能力ではなく、事情による不可を述べる。 | 建議改為 辞書形/じしょけい（source: 動詞の辞書形） | HIGH-likely |
| M066 | MEDIUM | N3/wake-niwa-ikanai.json | key_terms | 的 | てき | けにはいかない 「わけにはいかない」は、社会的、道徳的、または状況上の理由により、ある行為をすることができない、またはしてはいけな | NEEDS-NATIVE: 候選A=社会的/しゃかいてき; 候選B=道徳的/どうとくてき（both in source） | NEEDS-NATIVE |
| M068 | MEDIUM | N3/youni-goal.json | key_terms | 形 | がた | ように（目標・変化） 「V辞書形/ない形 + ように」は、できる状態を目標にする表現です。意志で直接コントロールしに | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M069 | MEDIUM | N3/youni-naru.json | key_terms | 形 | がた | ようになる 「V辞書形 + ようになる」は、前はできなかったことができる状態に変わることを表します。「Vな | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M070 | MEDIUM | N3/youni-suru.json | key_terms | 形 | がた | ようにする 「V辞書形/ない形 + ようにする」は、ある状態になるように努力したり、気をつけたりすることを | 建議改為 辞書形/じしょけい（source: V辞書形） | HIGH-likely |
| M072 | MEDIUM | N5/counter-basic.json | title_ja | 枚 | まい | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は一般 | 保持現狀正確: 枚/まい（source: 数え方 title + 一枚） | HIGH-likely |
| M073 | MEDIUM | N5/kore-sore-are.json | title_ja | 系 | けい | こ・そ・あ系 「こ・そ・あ・ど」は、話し手や聞き手からの距離を表す指示語の体系である。「これ」は | 保持現狀正確: 系/けい（source: こ・そ・あ系） | HIGH-likely |
| M077 | MEDIUM | N5/number-counter-hon.json | title_ja | 本 | ほん | 〜本 「本」は、ペン、瓶、木、傘など、細長い物を数える助数詞である。数字によって読み方が | NEEDS-NATIVE: 候選A=本/ほん for bare counter heading; 候選B=ぽん/ぼん after specific numerals | NEEDS-NATIVE |
| M079 | MEDIUM | N5/verb-te-form-connection.json | title_ja | 形 | けい | て形の接続 動詞のて形は、複数の動作を順番に述べるときに使う。「V1て、V2」は「V1を | 建議改為 て形/てけい（source: て形の接続） | HIGH-likely |

### 1.6 LOW — Other R-classes（R4, R6, R7, R8）
| id | level | entry | field | key_term | reading | classes | source-context-snippet |
| --- | --- | --- | --- | --- | --- | --- | --- |
| L001 | LOW | N1/de-are-de-are.json | key_terms | 候補 | こうほ | R7 | 〜であれ〜であれ 意味核：二つ以上の候補を挙げ、どちらの場合でも後件が同じであることを表す。条件差を超えた原則を示す。 接 |
| L002 | LOW | N1/de-are-de-are.json | key_terms | 後件 | こうけん | R7 | 味核：二つ以上の候補を挙げ、どちらの場合でも後件が同じであることを表す。条件差を超えた原則を示す。 接続：名詞A＋であれ＋名詞B＋ |
| L003 | LOW | N1/ka-ina-ka.json | key_terms | 二択 | にたく | R6 | 意味核：ある事柄が成立するかしないかを、硬く二択として示す。判断・検討・可否の対象を明確にする。 接続：動詞普通形＋か否か。い形容 |
| L004 | LOW | N1/te-yamanai.json | key_terms | 敬愛 | けいあい | R7 | 表現である。「成功を願ってやまない」「恩師を敬愛してやまない」のように、好意・祈り・惜別などを改まった調子で述べる。日常の一時的な感 |
| L005 | LOW | N1/wo-kawakiri-ni.json | key_terms | 出来事 | できごと | R6 | りに（して）／〜を皮切りとして 意味核：ある出来事を最初のきっかけとして、その後に同種の出来事が連続して広がることを表す。単なる開始時 |
| L006 | LOW | N1/wo-kawakiri-ni.json | key_terms | 東京 | とうきょう | R4 | 開の最初に位置づけられる名詞を取る。 例：東京公演を皮切りに全国を巡回する／首脳会談を皮切りとして協議が進む／新制度の導入を皮切り |
| L007 | LOW | N1/ya-inaya.json | key_terms | 出来事 | できごと | R6 | 否や 意味核：前の動作が成立した直後に、次の出来事が間を置かず起こることを表す。前件を境に状況が急に動く感じがあり、叙述には勢いと緊張 |
| L008 | LOW | N2/hanmen.json | key_terms | 側面 | そくめん | R7 | 意味核：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所の対比に多い。 接続：動詞普通形＋反面。い形容詞 |
| L009 | LOW | N2/hanmen.json | key_terms | 反対 | はんたい | R7 | 〜反面 意味核：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所の対比に多い。 接続：動詞普通形＋ |
| L010 | LOW | N2/hanmen.json | key_terms | 同時 | どうじ | R7 | ：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所の対比に多い。 接続：動詞普通形＋反面。い形容詞＋反面 |
| L011 | LOW | N2/hanmen.json | key_terms | 対照 | たいしょう | R7 | 〜反面 意味核：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所の対比に多い。 接続：動詞普通形＋反面。 |
| L012 | LOW | N2/hanmen.json | key_terms | 物事 | ものごと | R7 | 〜反面 意味核：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所の対比に多い。 接 |
| L013 | LOW | N2/hanmen.json | title_ja | 反面 | はんめん | R7 | 〜反面 意味核：一つの物事に、前件とは反対・対照的な側面が同時にあることを表す。長所と短所 |
| L014 | LOW | N2/ni-saishite.json | key_terms | 出来事 | できごと | R6 | 〜に際して 「〜に際して」は、特別な出来事や公式な行動を行う時に、という意味を表す。例として「出発に際して注意を受けた」「契約 |
| L015 | LOW | N2/ni-sakidatte.json | key_terms | 出来事 | できごと | R6 | 〜に先立って 「〜に先立って」は、重要な出来事や公式な行動の前に何かを行うことを表す。例として「開会に先立って注意事項を説明する」 |
| L016 | LOW | N2/nishitewa.json | key_terms | 予想 | よそう | R7 | 「〜にしては」は、前に示した基準から考えると予想と違う、という評価を表す。例として「新人にしては落ち着いている」「六月にしては涼しい |
| L017 | LOW | N2/nishitewa.json | key_terms | 六月 | ろくがつ | R7 | す。例として「新人にしては落ち着いている」「六月にしては涼しい」のように、立場・年齢・時期などを基準にして述べる。客観的な基準とのず |
| L018 | LOW | N2/nishitewa.json | key_terms | 年齢 | ねんれい | R7 | いる」「六月にしては涼しい」のように、立場・年齢・時期などを基準にして述べる。客観的な基準とのずれを述べる表現で、単なる比較よりも意 |
| L019 | LOW | N2/nishitewa.json | key_terms | 新人 | しんじん | R7 | ると予想と違う、という評価を表す。例として「新人にしては落ち着いている」「六月にしては涼しい」のように、立場・年齢・時期などを基準に |
| L020 | LOW | N2/nishitewa.json | key_terms | 時期 | じき | R7 | 「六月にしては涼しい」のように、立場・年齢・時期などを基準にして述べる。客観的な基準とのずれを述べる表現で、単なる比較よりも意外さが |
| L021 | LOW | N2/nishitewa.json | key_terms | 立場 | たちば | R7 | 着いている」「六月にしては涼しい」のように、立場・年齢・時期などを基準にして述べる。客観的な基準とのずれを述べる表現で、単なる比較よ |
| L022 | LOW | N2/nishitewa.json | key_terms | 評価 | ひょうか | R7 | 前に示した基準から考えると予想と違う、という評価を表す。例として「新人にしては落ち着いている」「六月にしては涼しい」のように、立場・ |
| L023 | LOW | N2/nitomonatte.json | key_terms | 売上 | うりあげ | R7 | 伴って交通量も増えた」「気温が上がるに伴って売上が伸びる」のように、社会的・数量的な変化によく使われる。書き言葉で、二つの変化の連動 |
| L024 | LOW | N2/nitomonatte.json | key_terms | 社会 | しゃかい | R7 | 気温が上がるに伴って売上が伸びる」のように、社会的・数量的な変化によく使われる。書き言葉で、二つの変化の連動を客観的に述べる。 |
| L025 | LOW | N2/ta-totan.json | key_terms | 出来事 | できごと | R6 | がした」のように、後件には話し手が意図しない出来事が来る。自分の意志で計画して行う後件には使いにくい。 |
| L026 | LOW | N2/te-hajimete.json | key_terms | 出来事 | できごと | R6 | 〜てはじめて 「〜てはじめて」は、ある経験や出来事の後で、それまで分からなかったことに気づくことを表す。例として「失敗してはじめて準備 |
| L027 | LOW | N2/tsutsu.json | key_terms | 勉強 | べんきょう | R7 | ります。 1. 同時進行：音楽を聞きつつ、勉強する 2. 逆接：悪いと知りつつ、嘘をついた 3. 変化の途中：状況は改善しつつある |
| L028 | LOW | N2/tsutsu.json | key_terms | 同時 | どうじ | R7 | 現です。主に三つの意味があります。 1. 同時進行：音楽を聞きつつ、勉強する 2. 逆接：悪いと知りつつ、嘘をついた 3. 変化の |
| L029 | LOW | N2/tsutsu.json | key_terms | 進行 | しんこう | R7 | す。主に三つの意味があります。 1. 同時進行：音楽を聞きつつ、勉強する 2. 逆接：悪いと知りつつ、嘘をついた 3. 変化の途中 |
| L030 | LOW | N2/tsutsu.json | key_terms | 音楽 | おんがく | R7 | に三つの意味があります。 1. 同時進行：音楽を聞きつつ、勉強する 2. 逆接：悪いと知りつつ、嘘をついた 3. 変化の途中：状況 |
| L031 | LOW | N2/wakeda.json | key_terms | 学生 | がくせい | R7 | わけだ。な形容詞は「上手なわけだ」、名詞は「学生な／学生というわけだ」のようにします。 例：彼は10年も日本にいた。日本語が上手な |
| L032 | LOW | N2/wakeda.json | key_terms | 日本語 | にほんご | R6+R7+R4 | うにします。 例：彼は10年も日本にいた。日本語が上手なわけだ。 この文では、「10年日本にいた」という理由から、「日本語が上手」と |
| L033 | LOW | N2/wakeda.json | key_terms | 理由 | りゆう | R7 | だ。 この文では、「10年日本にいた」という理由から、「日本語が上手」という結論が自然だと言っています。 「はずだ」は一般的な予想 |
| L034 | LOW | N2/wo-hajime.json | key_terms | 京都 | きょうと | R4 | 。例として「社長をはじめ全社員が参加した」「京都をはじめ多くの都市で行われた」のように使う。後件には複数の人・物・場所の広がりが続く |
| L035 | LOW | N2/wo-hajime.json | key_terms | 代表 | だいひょう | R7 | 〜をはじめ 「〜をはじめ」は、代表例を一つ挙げ、そのほかにも同じ種類のものがあることを表す。例として「社長をはじめ全社 |
| L036 | LOW | N2/wo-hajime.json | key_terms | 種類 | しゅるい | R7 | じめ」は、代表例を一つ挙げ、そのほかにも同じ種類のものがあることを表す。例として「社長をはじめ全社員が参加した」「京都をはじめ多くの |
| L037 | LOW | N2/wo-keiki-ni.json | key_terms | 出来事 | できごと | R6 | 〜を契機に 「〜を契機に」は、ある出来事がきっかけとなって、その後に大きな変化や行動が起こることを表す。例として「留学を契機 |
| L038 | LOW | N2/wo-tooshite.json | key_terms | 全体 | ぜんたい | R7 | て 「Nを通して」は、手段・媒介、または期間全体を表します。 - 友人を通して知りました - 一年を通して暖かいです 人や活動を |
| L039 | LOW | N2/wo-tooshite.json | key_terms | 友人 | ゆうじん | R7 | 段・媒介、または期間全体を表します。 - 友人を通して知りました - 一年を通して暖かいです 人や活動を媒介にする場合と、時間の |
| L040 | LOW | N2/wo-tooshite.json | key_terms | 媒介 | ばいかい | R7 | を通して 「Nを通して」は、手段・媒介、または期間全体を表します。 - 友人を通して知りました - 一年を通して暖かいで |
| L041 | LOW | N2/wo-tooshite.json | key_terms | 手段 | しゅだん | R7 | を通して 「Nを通して」は、手段・媒介、または期間全体を表します。 - 友人を通して知りました - 一年を通して暖 |
| L042 | LOW | N2/wo-tooshite.json | key_terms | 期間 | きかん | R7 | 通して 「Nを通して」は、手段・媒介、または期間全体を表します。 - 友人を通して知りました - 一年を通して暖かいです 人や活 |
| L043 | LOW | N2/wo-tooshite.json | key_terms | 活動 | かつどう | R7 | りました - 一年を通して暖かいです 人や活動を媒介にする場合と、時間の範囲全体を言う場合があります。 |
| L044 | LOW | N3/conditional-nara.json | key_terms | 東京 | とうきょう | R4 | いなら - 動詞：行くなら／行くのなら 「東京に行くなら、チケットを先に買ったほうがいい」のように、BがAより前に起きても使えます |
| L045 | LOW | N3/sae-ba.json | key_terms | 名詞 | めいし | R7 | が満たされれば十分だ、という意味を表します。名詞には「Nさえあれば」、動詞には「Vます形 + さえすれば」の形がよく使われます。  |
| L046 | LOW | N3/sae-ba.json | key_terms | 条件 | じょうけん | R7 | さえ〜ば 「さえ〜ば」は、その一つの条件が満たされれば十分だ、という意味を表します。名詞には「Nさえあれば」、動詞には「Vま |
| L047 | LOW | N3/tabi-ni.json | key_terms | 出来事 | できごと | R6 | たびに 「たびに」は、ある出来事が起こるごとに、別の出来事も繰り返し起こることを表す表現である。動詞の辞書形または「 |
| L048 | LOW | N3/teiku-tekuru.json | key_terms | 日本語 | にほんご | R6+R4 | - これから寒くなっていきます - 少しずつ日本語が分かってきました |
| L049 | LOW | N4/ato-de.json | key_terms | 出来事 | できごと | R6 | 後で 「後で」は、ある動作や出来事が終わってから何かをすることを表します。動詞は「た形 + 後で」を使います。名詞の場 |
| L050 | LOW | N4/demo-particle.json | key_terms | 一つ | ひとつ | R7 | の「でも」は、「名詞 + でも」の形で、例を一つ軽く示したり、やわらかい提案をしたりする表現である。また、「そのようなものでも」とい |
| L051 | LOW | N4/demo-particle.json | key_terms | 助詞 | じょし | R7 | でも 助詞の「でも」は、「名詞 + でも」の形で、例を一つ軽く示したり、やわらかい提案をしたり |
| L052 | LOW | N4/demo-particle.json | key_terms | 名詞 | めいし | R7 | でも 助詞の「でも」は、「名詞 + でも」の形で、例を一つ軽く示したり、やわらかい提案をしたりする表現である。また |
| L053 | LOW | N4/demo-particle.json | key_terms | 提案 | ていあん | R7 | も」の形で、例を一つ軽く示したり、やわらかい提案をしたりする表現である。また、「そのようなものでも」という意味で、意外な対象を含める |
| L054 | LOW | N4/demo-particle.json | key_terms | 表現 | ひょうげん | R7 | 一つ軽く示したり、やわらかい提案をしたりする表現である。また、「そのようなものでも」という意味で、意外な対象を含める場合にも用いられ |
| L055 | LOW | N4/koto-ga-dekiru.json | key_terms | 日本語 | にほんご | R6+R4 | できる」は、能力や可能性を表します。 - 日本語を話すことができます - ここで泳ぐことができます 会話では可能形もよく使いますが |
| L056 | LOW | N4/made-ni.json | key_terms | 動作 | どうさ | R7 | までに 「までに」は、ある時点より前に動作を終える必要があることを表します。「三時までに」は三時が期限であるという意味です。「 |
| L057 | LOW | N4/made-ni.json | key_terms | 必要 | ひつよう | R7 |  「までに」は、ある時点より前に動作を終える必要があることを表します。「三時までに」は三時が期限であるという意味です。「まで」は動作 |
| L058 | LOW | N4/made-ni.json | key_terms | 意味 | いみ | R7 | ます。「三時までに」は三時が期限であるという意味です。「まで」は動作や状態が続く終点を表すため、意味が異なります。 |
| L059 | LOW | N4/made-ni.json | key_terms | 時点 | じてん | R7 | までに 「までに」は、ある時点より前に動作を終える必要があることを表します。「三時までに」は三時が期限であるという |
| L060 | LOW | N4/made-ni.json | key_terms | 期限 | きげん | R7 | があることを表します。「三時までに」は三時が期限であるという意味です。「まで」は動作や状態が続く終点を表すため、意味が異なります。 |
| L061 | LOW | N4/mae-ni.json | key_terms | 出来事 | できごと | R6 | 前に 「前に」は、ある動作や出来事より先に何かをすることを表します。動詞は辞書形に接続し、「名詞 + の + 前に」も |
| L062 | LOW | N4/nagara.json | key_terms | 二つ | ふたつ | R7 | ながら 「動詞ます形の語幹 + ながら」は、二つの動作を同時にすることを表します。例えば「音楽を聞きながら勉強します」は聞くことと勉 |
| L063 | LOW | N4/nagara.json | key_terms | 動作 | どうさ | R7 |  「動詞ます形の語幹 + ながら」は、二つの動作を同時にすることを表します。例えば「音楽を聞きながら勉強します」は聞くことと勉強する |
| L064 | LOW | N4/nagara.json | key_terms | 動詞 | どうし | R7 | ながら 「動詞ます形の語幹 + ながら」は、二つの動作を同時にすることを表します。例えば「音楽を聞 |
| L065 | LOW | N4/nagara.json | key_terms | 語幹 | ごかん | R7 | ながら 「動詞ます形の語幹 + ながら」は、二つの動作を同時にすることを表します。例えば「音楽を聞きながら勉強 |
| L066 | LOW | N4/node-reason.json | key_terms | 明日 | あした | R6 |  - 雨が降っているので、行きません - 明日試験なので、勉強します 名詞とな形容詞の前では「なので」になります。 |
| L067 | LOW | N4/shi-reason.json | key_terms | 一つ | ひとつ | R7 | し 「普通形 + し」は、理由や事実を一つ以上挙げるときに使います。「安いし、便利です」のように、複数の理由を並べることが多い |
| L068 | LOW | N4/shi-reason.json | key_terms | 事実 | じじつ | R7 | し 「普通形 + し」は、理由や事実を一つ以上挙げるときに使います。「安いし、便利です」のように、複数の理由を並べること |
| L069 | LOW | N4/shi-reason.json | key_terms | 便利 | べんり | R7 | 実を一つ以上挙げるときに使います。「安いし、便利です」のように、複数の理由を並べることが多いです。名詞とナ形容詞では「だし」を使いま |
| L070 | LOW | N4/shi-reason.json | key_terms | 普通 | ふつう | R7 | し 「普通形 + し」は、理由や事実を一つ以上挙げるときに使います。「安いし、便利です」のよう |
| L071 | LOW | N4/shi-reason.json | key_terms | 理由 | りゆう | R7 | し 「普通形 + し」は、理由や事実を一つ以上挙げるときに使います。「安いし、便利です」のように、複数の理由を並べ |
| L072 | LOW | N4/sorekara.json | key_terms | 出来事 | できごと | R6 | それから 「それから」は、前に述べた出来事の次に起こることを続けて言うときに使います。時間の順序をはっきり示す接続表現です。ま |
| L073 | LOW | N4/ta-koto-ga-aru.json | key_terms | 京都 | きょうと | R4 | がある」は、今までの経験を表します。例えば「京都へ行ったことがあります」は行った経験を表し、「すしを作ったことがあります」は作った経 |
| L074 | LOW | N4/yori-comparison.json | key_terms | 今日 | きょう | R6 |  バスより電車のほうが速いです - 昨日より今日のほうが寒いです 「より」は比べる基準を示します。 |
| L075 | LOW | N5/counter-basic.json | key_terms | 一般 | いっぱん | R7 | 数え方を使います。「りんごが三つあります」は一般的な物、「学生が二人います」は人、「紙が一枚あります」は薄い物です。数える言葉は名詞 |
| L076 | LOW | N5/counter-basic.json | key_terms | 三つ | みっつ | R7 | 人を数えるときに数え方を使います。「りんごが三つあります」は一般的な物、「学生が二人います」は人、「紙が一枚あります」は薄い物です。 |
| L077 | LOW | N5/counter-basic.json | key_terms | 動詞 | どうし | R7 | にも置けますが、「本が二冊あります」のように動詞の前によく置きます。 |
| L078 | LOW | N5/counter-basic.json | key_terms | 名詞 | めいし | R7 | 紙が一枚あります」は薄い物です。数える言葉は名詞の後ろにも前にも置けますが、「本が二冊あります」のように動詞の前によく置きます。 |
| L079 | LOW | N5/counter-basic.json | key_terms | 学生 | がくせい | R7 | す。「りんごが三つあります」は一般的な物、「学生が二人います」は人、「紙が一枚あります」は薄い物です。数える言葉は名詞の後ろにも前に |
| L080 | LOW | N5/counter-basic.json | key_terms | 後ろ | うしろ | R7 | 枚あります」は薄い物です。数える言葉は名詞の後ろにも前にも置けますが、「本が二冊あります」のように動詞の前によく置きます。 |
| L081 | LOW | N5/counter-basic.json | key_terms | 日本語 | にほんご | R6+R7+R4 | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあります」は一般的な物、「 |
| L082 | LOW | N5/counter-basic.json | key_terms | 言葉 | ことば | R7 | 人、「紙が一枚あります」は薄い物です。数える言葉は名詞の後ろにも前にも置けますが、「本が二冊あります」のように動詞の前によく置きます |
| L083 | LOW | N5/counter-basic.json | title_ja | 数え方 | かぞえかた | R7 | 数え方（つ・人・枚） 日本語では、物や人を数えるときに数え方を使います。「りんごが三つあり |
| L084 | LOW | N5/dake-only.json | key_terms | 一つ | ひとつ | R7 | だけ 「だけ」は、数や物や人を一つに限る表現です。「水だけ飲みます」や「日曜日だけ休みです」のように、「それ以外はない |
| L085 | LOW | N5/dake-only.json | key_terms | 日曜日 | にちようび | R7 | を一つに限る表現です。「水だけ飲みます」や「日曜日だけ休みです」のように、「それ以外はない」という意味を表します。名詞の後ろにも、時間 |
| L086 | LOW | N5/dake-only.json | key_terms | 表現 | ひょうげん | R7 | だけ 「だけ」は、数や物や人を一つに限る表現です。「水だけ飲みます」や「日曜日だけ休みです」のように、「それ以外はない」という意 |
| L087 | LOW | N5/date-expression.json | key_terms | 一月 | いちがつ | R7 | 3日」のように、年、月、日の順である。月は「一月、二月」のように読み、四月は「しがつ」、七月は「しちがつ」、九月は「くがつ」と読む。 |
| L088 | LOW | N5/date-expression.json | key_terms | 七月 | しちがつ | R7 | 一月、二月」のように読み、四月は「しがつ」、七月は「しちがつ」、九月は「くがつ」と読む。日は一日から十日までに特別な読み方が多く、十 |
| L089 | LOW | N5/date-expression.json | key_terms | 九月 | くがつ | R7 | 読み、四月は「しがつ」、七月は「しちがつ」、九月は「くがつ」と読む。日は一日から十日までに特別な読み方が多く、十四日は「じゅうよっか |
| L090 | LOW | N5/date-expression.json | key_terms | 二月 | にがつ | R7 | のように、年、月、日の順である。月は「一月、二月」のように読み、四月は「しがつ」、七月は「しちがつ」、九月は「くがつ」と読む。日は一 |
| L091 | LOW | N5/date-expression.json | key_terms | 四月 | しがつ | R7 | の順である。月は「一月、二月」のように読み、四月は「しがつ」、七月は「しちがつ」、九月は「くがつ」と読む。日は一日から十日までに特別 |
| L092 | LOW | N5/date-expression.json | key_terms | 日付 | ひづけ | R7 | 年月日 日付は「年・月・日」を使って表す。順番は「2026年5月3日」のように、年、月、日の順で |
| L093 | LOW | N5/date-expression.json | key_terms | 読み方 | よみかた | R7 | くがつ」と読む。日は一日から十日までに特別な読み方が多く、十四日は「じゅうよっか」、二十日は「はつか」、二十四日は「にじゅうよっか」と |
| L094 | LOW | N5/date-expression.json | key_terms | 順番 | じゅんばん | R7 | 年月日 日付は「年・月・日」を使って表す。順番は「2026年5月3日」のように、年、月、日の順である。月は「一月、二月」のように読 |
| L095 | LOW | N5/date-expression.json | title_ja | 年月日 | ねんがっぴ | R7 | 年月日 日付は「年・月・日」を使って表す。順番は「2026年5月3日」のように、年、月、日 |
| L096 | LOW | N5/e-direction.json | key_terms | 日本 | にほん | R4 | として読むときは「え」と発音します。 - 日本へ行きます - こちらへ来てください 到着点を強く言う「に」と似ていますが、「へ」 |
| L097 | LOW | N5/ga-contrast.json | key_terms | 日本語 | にほんご | R6+R4 | の内容と後ろの内容をやわらかく対比します。「日本語は難しいですが、おもしろいです」や「行きたいですが、時間がありません」のように使いま |
| L098 | LOW | N5/ga-particle.json | key_terms | 日本語 | にほんご | R6+R4 |  疑問詞：だれが来ますか - 好き・わかる：日本語が好きです - 存在：机の上に本があります 「は」は話題、「が」は主語の焦点です。 |
| L099 | LOW | N5/ikura-nanbon.json | key_terms | 値段 | ねだん | R7 | いくら・何本 値段を尋ねるときは「いくら」を使う。数を尋ねるときは「何」に助数詞を付ける。細長い物を尋 |
| L100 | LOW | N5/ikura-nanbon.json | key_terms | 助数詞 | じょすうし | R7 | は「いくら」を使う。数を尋ねるときは「何」に助数詞を付ける。細長い物を尋ねる場合は「何本」を使い、「ペンは何本ありますか」のように言う |
| L101 | LOW | N5/ikura-nanbon.json | key_terms | 場合 | ばあい | R7 | きは「何」に助数詞を付ける。細長い物を尋ねる場合は「何本」を使い、「ペンは何本ありますか」のように言う。助数詞は数える物の形や種類に |
| L102 | LOW | N5/ikura-nanbon.json | title_ja | 何本 | なんぼん | R7 | いくら・何本 値段を尋ねるときは「いくら」を使う。数を尋ねるときは「何」に助数詞を付ける。細長い |
| L103 | LOW | N5/ka-question.json | key_terms | 日本語 | にほんご | R6+R4 | ます。 - 学生ですか - 行きますか 日本語の疑問文では、文末を上げて読むことが多いです。 |
| L104 | LOW | N5/kara-made.json | key_terms | 一つ | ひとつ | R7 | ます - 駅から学校まで歩きます どちらか一つだけでも使えます。 |
| L105 | LOW | N5/kara-made.json | key_terms | 場所 | ばしょ | R7 | 始まり、「まで」は終わりを表します。時間にも場所にも使えます。 - 9時から5時まで働きます - 駅から学校まで歩きます どちら |
| L106 | LOW | N5/kara-made.json | key_terms | 学校 | がっこう | R7 |  - 9時から5時まで働きます - 駅から学校まで歩きます どちらか一つだけでも使えます。 |
| L107 | LOW | N5/kara-made.json | key_terms | 時間 | じかん | R7 | から」は始まり、「まで」は終わりを表します。時間にも場所にも使えます。 - 9時から5時まで働きます - 駅から学校まで歩きます  |
| L108 | LOW | N5/masenka-invitation.json | key_terms | 明日 | あした | R6 | です。「いっしょに昼ご飯を食べませんか」や「明日、映画を見ませんか」のように使います。形は否定ですが、意味は「いっしょにしませんか」 |
| L109 | LOW | N5/masu-negative.json | key_terms | 今日 | きょう | R6 | ていねいに言う形です。「肉を食べません」や「今日は行きません」のように、しないことを表します。過去の否定は「ませんでした」で、「昨日 |
| L110 | LOW | N5/mo-particle.json | key_terms | 否定 | ひてい | R7 | - 私も学生です - コーヒーも飲みます 否定文では「誰も」「何も」のように使い、「一人もいない」「何も食べない」という意味になり |
| L111 | LOW | N5/mo-particle.json | key_terms | 学生 | がくせい | R7 | 〜も同じように」という意味です。 - 私も学生です - コーヒーも飲みます 否定文では「誰も」「何も」のように使い、「一人もいな |
| L112 | LOW | N5/nagara-simultaneous.json | key_terms | 二つ | ふたつ | R7 | ながら 「ながら」は、二つの動作を同時にすることを表します。動詞のます形から「ます」を取って、「音楽を聞きなが |
| L113 | LOW | N5/nagara-simultaneous.json | key_terms | 動作 | どうさ | R7 | ながら 「ながら」は、二つの動作を同時にすることを表します。動詞のます形から「ます」を取って、「音楽を聞きながら勉強 |
| L114 | LOW | N5/nagara-simultaneous.json | key_terms | 動詞 | どうし | R7 | 」は、二つの動作を同時にすることを表します。動詞のます形から「ます」を取って、「音楽を聞きながら勉強します」や「歩きながら話します」 |
| L115 | LOW | N5/ni-time.json | key_terms | 今日 | きょう | R6+R7 |  7時に起きます - 月曜日に行きます 「今日」「明日」「毎日」のように、相対的な時間や習慣を表す言葉には普通「に」を付けません。 |
| L116 | LOW | N5/ni-time.json | key_terms | 動作 | どうさ | R7 | に（時間） 「に」は、動作が起こる時間を示します。数字を含むはっきりした時間によく使います。 - 7時に起き |
| L117 | LOW | N5/ni-time.json | key_terms | 数字 | すうじ | R7 | 間） 「に」は、動作が起こる時間を示します。数字を含むはっきりした時間によく使います。 - 7時に起きます - 月曜日に行きます  |
| L118 | LOW | N5/ni-time.json | key_terms | 明日 | あした | R6 | 起きます - 月曜日に行きます 「今日」「明日」「毎日」のように、相対的な時間や習慣を表す言葉には普通「に」を付けません。 |
| L119 | LOW | N5/ni-time.json | key_terms | 月曜日 | げつようび | R7 | によく使います。 - 7時に起きます - 月曜日に行きます 「今日」「明日」「毎日」のように、相対的な時間や習慣を表す言葉には普通 |
| L120 | LOW | N5/no-possessive.json | key_terms | 日本語 | にほんご | R6+R4 | 係するB」という意味です。 - 私の本 - 日本語の先生 - 机の上 日本語では「の」を使う関係が多いので、名詞を説明するときの基本 |
| L121 | LOW | N5/tai-desire.json | key_terms | 日本 | にほん | R4 | します。動詞のます形から「ます」を取って、「日本へ行きたいです」や「水を飲みたいです」のように使います。「たい」はい形容詞のように、 |
| L122 | LOW | N5/time-expression.json | key_terms | 十分 | じゅうぶん | R7 | は「時」と「分」を使って表す。「三時」「三時十分」のように、時間の後に分を置く。四時は「よじ」、七時は「しちじ」、九時は「くじ」と読 |
| L123 | LOW | N5/time-expression.json | key_terms | 午前 | ごぜん | R7 | 六分、八分、十分は「ぷん」で読むことが多い。午前は昼の前、午後は昼の後を表す。 |
| L124 | LOW | N5/time-expression.json | key_terms | 時刻 | じこく | R7 | 時刻 時刻は「時」と「分」を使って表す。「三時」「三時十分」のように、時間の後に分を置く |
| L125 | LOW | N5/time-expression.json | key_terms | 時間 | じかん | R7 | を使って表す。「三時」「三時十分」のように、時間の後に分を置く。四時は「よじ」、七時は「しちじ」、九時は「くじ」と読む。分は「ぷん」 |
| L126 | LOW | N5/time-expression.json | title_ja | 時刻 | じこく | R7 | 時刻 時刻は「時」と「分」を使って表す。「三時」「三時十分」のように、時間の後に分を置く |
| L127 | LOW | N5/wa-particle.json | key_terms | 今日 | きょう | R6 | です」という形です。 - 私は学生です - 今日は暑いです 「が」との違いが大切です。「は」は話題やテーマを出します。「が」は主語 |

### 1.7 R-title-uncovered（title_ja 含漢字但無 furigana key_term 涵蓋）
No R-title-uncovered hits. All kanji-bearing titles are covered by `annotations.furigana.title_ja` and/or `key_terms`.

## Section 2 — Native-Reviewer Second-Pass (TO BE FILLED BY HUMAN)

Per the requested native-review workflow, this section is for a human native reviewer. Use the row `id` values from Section 1 and fill verdict/correction/notes.

### 2.1 HIGH Rows
| id | entry | field | key_term | native-verdict (✓/✗/?) | corrected | notes |
| --- | --- | --- | --- | --- | --- | --- |
| H001 | N1/atte-no.json | key_terms | 支/ささ |  |  |  |
| H002 | N1/bakoso.json | key_terms | 厳/きび |  |  |  |
| H003 | N1/ga-hayai-ka.json | title_ja | 早/はや |  |  |  |
| H004 | N1/gotoki.json | key_terms | 使/つか |  |  |  |
| H005 | N1/gotoki.json | key_terms | 分/わ |  |  |  |
| H006 | N1/gotoki.json | key_terms | 走/はし |  |  |  |
| H007 | N1/kagiri-da.json | title_ja | 限/かぎ |  |  |  |
| H008 | N1/te-yamanai.json | key_terms | 惜/お |  |  |  |
| H009 | N1/te-yamanai.json | key_terms | 愛/あい |  |  |  |
| H010 | N1/te-yamanai.json | key_terms | 願/ねが |  |  |  |
| H011 | N1/towaie.json | key_terms | 寒/さむ |  |  |  |
| H012 | N1/zujimai.json | key_terms | 聞/き |  |  |  |
| H013 | N1/zuniwa-irarenai.json | key_terms | 笑/わら |  |  |  |
| H014 | N2/ageku.json | key_terms | 至/いた |  |  |  |
| H015 | N2/dokoro-dewa-nai.json | key_terms | 厳/きび |  |  |  |
| H016 | N2/nai-koto-wa-nai.json | key_terms | 控/ひか |  |  |  |
| H017 | N2/ni-hanshite.json | title_ja | 反/はん |  |  |  |
| H018 | N2/ni-kagiri.json | title_ja | 限/かぎ |  |  |  |
| H019 | N2/ni-saishite.json | title_ja | 際/さい |  |  |  |
| H020 | N2/ni-shitagatte.json | key_terms | 合/あ |  |  |  |
| H021 | N2/ni-tsurete.json | key_terms | 変/か |  |  |  |
| H022 | N2/nikanshite.json | title_ja | 関/かん |  |  |  |
| H023 | N2/ta-totan.json | key_terms | 話/はな |  |  |  |
| H024 | N2/warini.json | key_terms | 忙/いそが |  |  |  |
| H025 | N2/wo-tooshite.json | title_ja | 通/とお |  |  |  |
| H026 | N3/contrast-noni.json | key_terms | 話/はな |  |  |  |
| H027 | N3/hazuganai.json | key_terms | 話/はな |  |  |  |
| H028 | N3/kagiri.json | title_ja | 限/かぎ |  |  |  |
| H029 | N3/koto-ni-natte-iru.json | key_terms | 取/と |  |  |  |
| H030 | N3/ni-chigainai.json | key_terms | 話/はな |  |  |  |
| H031 | N3/ni-chigainai.json | title_ja | 違/ちが |  |  |  |
| H032 | N3/okage-de.json | key_terms | 助/たす |  |  |  |
| H033 | N3/sei-de.json | key_terms | 話/はな |  |  |  |
| H034 | N3/tabi-ni.json | key_terms | 繰/く |  |  |  |
| H035 | N3/tabi-ni.json | key_terms | 返/かえ |  |  |  |
| H036 | N3/te-bakari-iru.json | key_terms | 話/はな |  |  |  |
| H037 | N3/teiku-tekuru.json | key_terms | 話/はな |  |  |  |
| H038 | N5/counter-basic.json | title_ja | 人/にん |  |  |  |

### 2.2 MEDIUM Rows
| id | entry | field | key_term | native-verdict (✓/✗/?) | corrected | notes |
| --- | --- | --- | --- | --- | --- | --- |
| M001 | N1/bekarazu.json | key_terms | 不/ふ |  |  |  |
| M002 | N1/de-are-de-are.json | key_terms | 差/さ |  |  |  |
| M003 | N1/ga-hayai-ka.json | key_terms | 差/さ |  |  |  |
| M004 | N1/ga-hayai-ka.json | key_terms | 速/はや |  |  |  |
| M005 | N1/gotoki.json | key_terms | 夢/ゆめ |  |  |  |
| M006 | N1/gotoki.json | key_terms | 風/かぜ |  |  |  |
| M007 | N1/ja-arumai-shi.json | key_terms | 子/こ |  |  |  |
| M008 | N1/ka-ina-ka.json | title_ja | 否/いな |  |  |  |
| M009 | N1/katagata.json | key_terms | 礼/れい |  |  |  |
| M010 | N1/kirai-ga-aru.json | key_terms | 癖/くせ |  |  |  |
| M011 | N1/nagara-ni.json | key_terms | 昔/むかし |  |  |  |
| M012 | N1/nagara-ni.json | key_terms | 涙/なみだ |  |  |  |
| M013 | N1/nimokakawarazu.json | key_terms | 書/か |  |  |  |
| M014 | N1/nimokakawarazu.json | key_terms | 雨/あめ |  |  |  |
| M015 | N1/taru-mono.json | key_terms | 姿/すがた |  |  |  |
| M016 | N1/toittemo-kagonai.json | key_terms | 町/まち |  |  |  |
| M017 | N1/tomonaku.json | key_terms | 音/おと |  |  |  |
| M018 | N1/towaie.json | key_terms | 春/はる |  |  |  |
| M019 | N1/wo-kawakiri-ni.json | key_terms | 後/あと |  |  |  |
| M020 | N1/wo-oite.json | key_terms | 彼/かれ |  |  |  |
| M021 | N1/wo-oite.json | title_ja | 他/た |  |  |  |
| M022 | N1/ya-inaya.json | title_ja | 否/いな |  |  |  |
| M023 | N2/monoda.json | key_terms | 驚/おどろ |  |  |  |
| M024 | N2/monono-formal.json | key_terms | 異/こと |  |  |  |
| M025 | N2/monono-formal.json | key_terms | 認/みと |  |  |  |
| M026 | N2/ni-motozuite.json | title_ja | 基/もと |  |  |  |
| M027 | N2/ni-oujite.json | title_ja | 応/おう |  |  |  |
| M028 | N2/ta-totan.json | key_terms | 手/て |  |  |  |
| M029 | N2/ta-totan.json | key_terms | 雨/あめ |  |  |  |
| M030 | N2/ue-de.json | title_ja | 上/うえ |  |  |  |
| M031 | N2/ue-wa.json | title_ja | 上/うえ |  |  |  |
| M032 | N2/wo-chuushin-ni.json | key_terms | 駅/えき |  |  |  |
| M033 | N3/conditional-ba.json | key_terms | 五/ご |  |  |  |
| M034 | N3/conditional-ba.json | key_terms | 段/だん | ✓ confirmed (mechanical-equivalent) | 五段/ごだん | PM analysis 2026-05-09: same lexeme class as existing 一段/いちだん entry — mechanical merge per F1 pattern |
| M035 | N3/conditional-ba.json | key_terms | 的/てき |  |  |  |
| M036 | N3/contrast-noni.json | key_terms | 形/がた |  |  |  |
| M037 | N3/contrast-noni.json | key_terms | 手/て |  |  |  |
| M038 | N3/dokoroka.json | key_terms | 形/がた |  |  |  |
| M039 | N3/hazuganai.json | key_terms | 性/せい |  |  |  |
| M040 | N3/hazuganai.json | key_terms | 手/て |  |  |  |
| M041 | N3/hodo.json | key_terms | 形/がた |  |  |  |
| M042 | N3/kagiri.json | key_terms | 形/がた |  |  |  |
| M043 | N3/kamoshirenai.json | key_terms | 形/がた |  |  |  |
| M044 | N3/kamoshirenai.json | key_terms | 性/せい |  |  |  |
| M045 | N3/kawari-ni.json | key_terms | 形/がた |  |  |  |
| M046 | N3/koto-ni-natte-iru.json | key_terms | 決/き |  |  |  |
| M047 | N3/kotoni-naru.json | key_terms | 形/がた |  |  |  |
| M048 | N3/kotoni-naru.json | key_terms | 決/き |  |  |  |
| M049 | N3/kotoni-suru.json | key_terms | 形/がた |  |  |  |
| M050 | N3/monono.json | key_terms | 形/がた |  |  |  |
| M051 | N3/monono.json | key_terms | 書/か |  |  |  |
| M052 | N3/ni-chigainai.json | key_terms | 形/がた |  |  |  |
| M053 | N3/ni-chigainai.json | key_terms | 手/て |  |  |  |
| M054 | N3/okage-de.json | key_terms | 形/がた |  |  |  |
| M055 | N3/sei-de.json | key_terms | 形/がた |  |  |  |
| M056 | N3/sei-de.json | key_terms | 手/て |  |  |  |
| M057 | N3/souda-hearsay.json | key_terms | 形/がた |  |  |  |
| M058 | N3/tabi-ni.json | key_terms | 形/がた |  |  |  |
| M059 | N3/tabi-ni.json | key_terms | 的/てき |  |  |  |
| M060 | N3/tameni-purpose.json | key_terms | 形/がた |  |  |  |
| M061 | N3/te-bakari-iru.json | key_terms | 手/て |  |  |  |
| M062 | N3/teiku-tekuru.json | key_terms | 動/うご |  |  |  |
| M063 | N3/teiku-tekuru.json | key_terms | 手/て |  |  |  |
| M064 | N3/tokoro.json | key_terms | 中/ちゅう |  |  |  |
| M065 | N3/wake-niwa-ikanai.json | key_terms | 形/がた |  |  |  |
| M066 | N3/wake-niwa-ikanai.json | key_terms | 的/てき | ✓ confirmed (mechanical-equivalent) | 社会的/しゃかいてき + 道徳的/どうとくてき | PM analysis 2026-05-09: 的 suffix split from 漢語 + 的 compounds — same as F1 stem-truncation class |
| M067 | N3/youda.json | key_terms | 書/か |  |  |  |
| M068 | N3/youni-goal.json | key_terms | 形/がた |  |  |  |
| M069 | N3/youni-naru.json | key_terms | 形/がた |  |  |  |
| M070 | N3/youni-suru.json | key_terms | 形/がた |  |  |  |
| M071 | N3/youni-suru.json | key_terms | 心/こころ |  |  |  |
| M072 | N5/counter-basic.json | title_ja | 枚/まい |  |  |  |
| M073 | N5/kore-sore-are.json | title_ja | 系/けい |  |  |  |
| M074 | N5/masu-form.json | title_ja | 形/けい |  |  |  |
| M075 | N5/nani-doko-dare.json | title_ja | 何/なに |  |  |  |
| M076 | N5/nani-doko-dare.json | title_ja | 誰/だれ |  |  |  |
| M077 | N5/number-counter-hon.json | title_ja | 本/ほん | ✓ confirmed (keep) | (keep ほん) | PM analysis 2026-05-09: ほん is the citation form for the bare counter title; prose explicitly enumerates ぽん/ぼん/ほん variants by numeral — current annotation is教材-standard |
| M078 | N5/te-form.json | title_ja | 形/けい |  |  |  |
| M079 | N5/verb-te-form-connection.json | title_ja | 形/けい |  |  |  |

## Summary Stats
- HIGH count: 38
- MEDIUM count: 79
- LOW count: 127
- 已被 F1 修正案例: 4（皮切り / 先立って / 伴って / 心構え）
- 預估 native-review 工作量（HIGH+MEDIUM 條目數）: 117
- R-title-uncovered count: 0
- Risk rows by level: N1=42, N2=58, N3=56, N4=26, N5=62
- Furigana pairs by level: N1=327, N2=305, N3=267, N4=286, N5=312
- Class hit counts: R1=117, R2-confirmed-truncation=38, R8=38, R5=48, R3=16, R6=30, R7=99, R2=13, R4=14
