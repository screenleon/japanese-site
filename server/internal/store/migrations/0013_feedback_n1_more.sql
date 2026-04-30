-- 0013_feedback_n1_more.sql
-- Generic deterministic feedback templates for the second N1 grammar corpus
-- pass under JS-009.

INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
VALUES
('aru-majiki', 'generic', '這題練的是「〜まじき」，表示以某身分或立場來說「不該有／不符合資格」。形態是動詞辞書形＋まじき＋名詞，常帶強烈批判。', 'curated', 'CC-BY-SA-4.0'),
('gotoki', 'generic', '這題練的是「〜ごとき／〜ごとく／〜ごとし」，是文語、正式的「像…一樣」。名詞前多用「ごとき」，修飾動詞或形容詞多用「ごとく」。', 'curated', 'CC-BY-SA-4.0'),
('taru-mono', 'generic', '這題練的是「〜たるもの」，表示「身為…者」。後面通常接應有的責任、態度或義務，語氣正式而有主張性。', 'curated', 'CC-BY-SA-4.0'),
('toittemo-kagonai', 'generic', '這題練的是「〜といっても過言ではない」，表示「即使說成…也不為過」。用來強調某判斷雖強烈但並非誇張。', 'curated', 'CC-BY-SA-4.0'),
('yogi-naku-sareru', 'generic', '這題練的是「〜を余儀なくされる」，表示因外在情勢而「被迫…／不得不…」。形態是名詞＋を余儀なくされる。', 'curated', 'CC-BY-SA-4.0')
ON CONFLICT(grammar_point, error_class) DO UPDATE SET
  body_zh=excluded.body_zh,
  source=excluded.source,
  license=excluded.license;
