-- 0012_feedback_n1.sql
-- Generic deterministic feedback templates for the first N1 grammar anchor
-- points. Specific error classes can be added later as the classifier rules
-- mature.

INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
VALUES
('zuniwa-irarenai', 'generic', '這題練的是「〜ずにはいられない」，表示情緒或情勢強到「忍不住不做／不能不做」。注意用動詞ない形去ない＋ず，する 要變成 せず。', 'curated', 'CC-BY-SA-4.0'),
('nimokakawarazu', 'generic', '這題練的是「〜にもかかわらず」，正式書面語，表示「儘管 A，卻 B」。不要直接用口語的「のに」取代。', 'curated', 'CC-BY-SA-4.0'),
('kirai-ga-aru', 'generic', '這題練的是「〜きらいがある」，表示「有不太好的傾向」。通常帶輕微批判，名詞前要用「の」。', 'curated', 'CC-BY-SA-4.0'),
('bakoso', 'generic', '這題練的是「〜ばこそ」，表示「正因為 A 才 B」。重點是ば形＋こそ；名詞、な形容詞常用「であればこそ」。', 'curated', 'CC-BY-SA-4.0'),
('towaie', 'generic', '這題練的是「〜とはいえ」，表示「雖說 A，但是 B」。先承認前項，再提出限制、補充或反向事實。', 'curated', 'CC-BY-SA-4.0')
ON CONFLICT(grammar_point, error_class) DO UPDATE SET
  body_zh=excluded.body_zh,
  source=excluded.source,
  license=excluded.license;
