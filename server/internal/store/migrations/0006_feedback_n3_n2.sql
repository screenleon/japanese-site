-- 0006_feedback_n3_n2.sql
-- Feedback templates for the new N3 + N2 grammar points.
-- Each grammar point has at least 'generic' so the grader always has
-- something to say; family-confusion error classes give more pointed help.

INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
VALUES
  -- ── N3: conditional-tara ───────────────────────────────────────────────
  ('conditional-tara', 'used-ba',
   '你用了「〜ば」形，但題目語意需要「〜たら」。「〜ば」強調邏輯成立，B 子句不能含意志/命令；當 B 是「請打電話/想去」這類意志時，要用「〜たら」。把語幹改成た形＋ら：例如 着く→着いたら、行く→行ったら。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-tara', 'used-nara',
   '你用了「〜なら」，但這題不是承接對方話題的條件。「〜たら」表示時間/一般條件（A 後 B 才發生）；「〜なら」是「以對方提到的事為前提」。改成た形＋ら。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-tara', 'used-to',
   '你用了「〜と」，但 B 子句含意志/邀請/命令時不能用「〜と」。在這種情境用「〜たら」（た形＋ら）。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-tara', 'wrong-conjugation',
   '形態錯了。「〜たら」是動詞**た形＋ら**：食べる→食べたら、行く→行ったら（不規則）、来る→来たら、する→したら。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-tara', 'generic',
   '答案是「〜たら」（動詞た形＋ら）。萬用條件句，且 B 子句可以含意志/命令。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N3: conditional-ba ────────────────────────────────────────────────
  ('conditional-ba', 'used-tara',
   '你用了「〜たら」，但本題語意是邏輯型條件（A 必然導致 B）。「〜ば」更貼切：一段動詞→れば（食べれば）、五段→え段＋ば（押せば）、い形容詞→ければ（高ければ）。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-ba', 'wrong-conjugation',
   '「〜ば」變化規則：一段動詞去る加れば；五段把う段改成え段＋ば；い形容詞 い→ければ；な形容詞/名詞用「〜なら」。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-ba', 'generic',
   '答案是「〜ば」形。記住 B 子句不能含意志/命令（「駅に着けば、電話してください」是錯的，要用「着いたら」）。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N3: conditional-nara ──────────────────────────────────────────────
  ('conditional-nara', 'used-tara',
   '你用了「〜たら」。但題目情境是「以對方提到的事/既知事實為前提」（給定話題下的條件），這時用「〜なら」更自然。動詞辭書形/名詞/形容詞直接接 なら。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-nara', 'used-ba',
   '你用了「〜ば」。「〜ば」強調邏輯關係；「〜なら」強調「以…為前提」。本題接話題用「〜なら」更貼切。',
   'curated', 'CC-BY-SA-4.0'),
  ('conditional-nara', 'generic',
   '答案是「〜なら」。用於承接對方話題或設定前提。名詞/な形容詞直接接「なら」（不加 だ）；動詞辭書形＋（の）なら。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N3: concession-temo ────────────────────────────────────────────────
  ('concession-temo', 'used-tara',
   '你用了「〜たら」（順接），但本題是讓步（即使…也…）。要用「〜ても」：て形＋も。例如 降る→降っても、高い→高くても、名詞→でも。',
   'curated', 'CC-BY-SA-4.0'),
  ('concession-temo', 'used-noni',
   '你用了「〜のに」（既成事實的轉折），但本題是假設讓步。要用「〜ても」：假設 A 也會 B。',
   'curated', 'CC-BY-SA-4.0'),
  ('concession-temo', 'wrong-conjugation',
   '「〜ても」是 て形＋も。動詞：食べる→食べても；い形容詞：高い→高くても；な形容詞・名詞：暇でも、学生でも；否定：行かなくても。',
   'curated', 'CC-BY-SA-4.0'),
  ('concession-temo', 'generic',
   '答案是「〜ても」（て形＋も）。表「即使 A 也 B」，A 是假設。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N3: contrast-noni ─────────────────────────────────────────────────
  ('contrast-noni', 'used-temo',
   '你用了「〜ても」（假設讓步），但本題 A 是已成事實，B 是不滿/意外的事實。要用「〜のに」：普通形＋のに（名詞/な形容詞接「な」）。',
   'curated', 'CC-BY-SA-4.0'),
  ('contrast-noni', 'used-kedo',
   '你用了「〜けど／〜が」中性轉折，但本題的轉折帶有情緒（不滿/意外/失望），用「〜のに」更貼切。',
   'curated', 'CC-BY-SA-4.0'),
  ('contrast-noni', 'generic',
   '答案是「〜のに」。普通形＋のに（名詞/な形容詞要加「な」：日曜日**な**のに、静か**な**のに）。表示帶情緒的轉折。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N2: wakeda ─────────────────────────────────────────────────────────
  ('wakeda', 'used-hazu',
   '你用了「〜はずだ」。「はずだ」基於常識的「應該」；「わけだ」基於具體前提推結論「難怪/原來如此」。本題有具體前提，「〜わけだ」更貼切。',
   'curated', 'CC-BY-SA-4.0'),
  ('wakeda', 'used-koto',
   '你用了「〜ことだ」（建議/應該）。本題語意是「結論/推得」，要用「〜わけだ」。',
   'curated', 'CC-BY-SA-4.0'),
  ('wakeda', 'generic',
   '答案是「〜わけだ」。從前提推結論（「難怪…」「也就是說…」）。注意名詞「〜なわけだ／〜というわけだ」、な形容詞「〜なわけだ」。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N2: monoda ─────────────────────────────────────────────────────────
  ('monoda', 'used-koto',
   '你用了「〜ことだ」。「〜ものだ」表「普遍真理／懷舊／感慨」；「〜ことだ」是給人的「建議/應該」。語意不同。',
   'curated', 'CC-BY-SA-4.0'),
  ('monoda', 'used-beki',
   '你用了「〜べきだ」（強烈應該）。本題語意是「本來就如此/常態」，要用「〜ものだ」。',
   'curated', 'CC-BY-SA-4.0'),
  ('monoda', 'generic',
   '答案是「〜ものだ」。三大用法：①普遍真理（子供は泣くものだ）②懷舊（よく遊んだものだ）③感慨（よく我慢したものだ）。注意名詞接「なもの」、な形容詞接「なもの」。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N2: dokoroka ──────────────────────────────────────────────────────
  ('dokoroka', 'used-bakarika',
   '你用了「〜ばかりか」。「ばかりか」是「不只 A，還 B（同方向加強）」；「どころか」是「哪是 A，還 B（反向或更極端）」。本題語意反向，要用「〜どころか」。',
   'curated', 'CC-BY-SA-4.0'),
  ('dokoroka', 'used-kurai',
   '你用了「〜くらい／〜ぐらい」。它是程度比較；「〜どころか」是強烈反向否定。本題要強調「哪裡是…根本是反向」，用「どころか」。',
   'curated', 'CC-BY-SA-4.0'),
  ('dokoroka', 'generic',
   '答案是「〜どころか」。表「哪裡是 A，反而 B」或「別說 A，連 B 都…」。直接接名詞、辭書形、形容詞詞幹。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N2: nikanshite ────────────────────────────────────────────────────
  ('nikanshite', 'used-nitotte',
   '你用了「〜にとって」。「にとって」表「對…來說（立場/觀點）」；本題是表「主題/關於」，用「〜について」或「〜に関して」。',
   'curated', 'CC-BY-SA-4.0'),
  ('nikanshite', 'used-nitaishite',
   '你用了「〜に対して」。「に対して」表「對於…（行為對象）」；本題是「關於主題」，用「〜について／に関して」。',
   'curated', 'CC-BY-SA-4.0'),
  ('nikanshite', 'generic',
   '答案是「〜について」或「〜に関して」（兩者皆表「關於」）。書面/正式語境較多用「〜に関して／に関する」；口語用「〜について／についての」。',
   'curated', 'CC-BY-SA-4.0'),

  -- ── N2: tsutsu ────────────────────────────────────────────────────────
  ('tsutsu', 'used-nagara',
   '你用了「〜ながら」。意思接近，但「〜つつ」更**書面/文藝**，且能表「漸進」（〜つつある）。N2 考點看你是否能辨識書面用法。',
   'curated', 'CC-BY-SA-4.0'),
  ('tsutsu', 'used-teiru',
   '你用了「〜ている」。「〜ている」表持續狀態；「〜つつある」表**漸進變化中**（緩慢、書面）。本題語意是漸進，要用「〜つつある」。',
   'curated', 'CC-BY-SA-4.0'),
  ('tsutsu', 'wrong-conjugation',
   '「〜つつ」接動詞**ます形**（去ます）：聞きます→聞きつつ、知ります→知りつつ、改善します→改善しつつ。',
   'curated', 'CC-BY-SA-4.0'),
  ('tsutsu', 'generic',
   '答案是「〜つつ」。動詞ます形＋つつ。三大用法：①一邊 A 一邊 B（書面的 ながら）②雖然 A 卻 B（〜つつも）③漸進中（〜つつある）。',
   'curated', 'CC-BY-SA-4.0');
