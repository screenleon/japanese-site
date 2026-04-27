-- 0004_feedback_seed.sql
-- Seed feedback templates for the deterministic grader. Templates are keyed
-- by (grammar_point, error_class). The grader matches the user's mistake to
-- one of these classes; if no specific match, it falls back to 'generic'.
--
-- Templates here are the *floor* — additional rows can be added by L1 corpus
-- loaders later if we want them git-curated.

INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
VALUES
  ('te-form', 'used-masu',
   '你寫了ます形，但題目要的是「て形」。「て形」用來連接兩個動作或句子（A 然後 B），不是 ます 形。把語幹改成對應的て形：例如 食べます → 食べて、飲みます → 飲んで。',
   'curated', 'CC-BY-SA-4.0'),
  ('te-form', 'used-dictionary',
   '你寫了原形（辭書形），但題目要的是「て形」。原形不能直接接後續動作。把原形改成て形：例えば 食べる → 食べて、行く → 行って（不規則）。',
   'curated', 'CC-BY-SA-4.0'),
  ('te-form', 'wrong-conjugation',
   '你選的形式不是標準的て形。檢查語尾規則：う/つ/る → って，む/ぶ/ぬ → んで，く → いて（行く例外為行って），ぐ → いで，す → して。',
   'curated', 'CC-BY-SA-4.0'),
  ('te-form', 'generic',
   '這題答案是「て形」。標準變化：一段去る加て；五段依語尾規則（う/つ/る→って、む/ぶ/ぬ→んで、く→いて、ぐ→いで、す→して）。',
   'curated', 'CC-BY-SA-4.0'),

  ('masu-form', 'used-plain',
   '你寫了原形，但題目要的是「ます形」（禮貌體）。原形 → ます形：一段去る加ます（食べる→食べます）；五段把う段改成い段加ます（飲む→飲みます）。',
   'curated', 'CC-BY-SA-4.0'),
  ('masu-form', 'wrong-tense',
   '時態不對。檢查上下文：昨日 → ました（過去肯定）/ ませんでした（過去否定）；明日/今 → ます / ません。',
   'curated', 'CC-BY-SA-4.0'),
  ('masu-form', 'generic',
   '這題答案是「ます形」。一段：去る加ます；五段：う段→い段加ます；不規則：する→します、来る→来ます。',
   'curated', 'CC-BY-SA-4.0'),

  ('copula-desu', 'used-da-instead',
   '在禮貌句裡用了普通形「だ」。書面/正式場合用「です」。',
   'curated', 'CC-BY-SA-4.0'),
  ('copula-desu', 'wrong-tense',
   '注意時態：です（現在）/ でした（過去）/ ではありません（現在否定）/ ではありませんでした（過去否定）。',
   'curated', 'CC-BY-SA-4.0'),
  ('copula-desu', 'generic',
   '這裡需要繫詞「です」或其變化形。です接在名詞/な形容詞後；不接動詞。',
   'curated', 'CC-BY-SA-4.0'),

  ('wa-particle', 'used-ga',
   '這裡應該用「は」標示主題，不是「が」標示主格。一般談論已知主題（A 是 B）用は；強調是「誰」做的或回答疑問詞時用が。',
   'curated', 'CC-BY-SA-4.0'),
  ('wa-particle', 'used-wo',
   '「を」是受詞助詞，不是主題助詞。談論主題時用「は」。',
   'curated', 'CC-BY-SA-4.0'),
  ('wa-particle', 'generic',
   '這裡需要主題助詞「は」（發音 wa，但寫成 hiragana 的 は）。標示「我們現在要談的東西」。',
   'curated', 'CC-BY-SA-4.0'),

  ('ga-particle', 'used-wa',
   '這裡應該用「が」。如果是回答疑問詞主語（誰、何）、強調主語、或存在句（〜があります），都用が，不是は。',
   'curated', 'CC-BY-SA-4.0'),
  ('ga-particle', 'used-wo',
   '「を」是受詞助詞，這裡需要主格助詞「が」。',
   'curated', 'CC-BY-SA-4.0'),
  ('ga-particle', 'generic',
   '這裡需要主格助詞「が」。常見場景：疑問詞主語、強調主語、存在句（〜があります）。',
   'curated', 'CC-BY-SA-4.0');
