-- 0022_grammar_slug_dedup.sql
--
-- Full rename map (11 pairs):
-- | Old level | Old slug              | New level | New slug              |
-- | N3        | hazuda                | N4        | hazu-da               |
-- | N3        | hazuganai             | N4        | hazu-ga-nai           |
-- | N3        | kamoshirenai          | N4        | kamo-shirenai         |
-- | N3        | teshimau              | N4        | te-shimau             |
-- | N4        | mono-da               | N3        | mono-da-norm          |
-- | N2        | monoda                | N2        | mono-da-emotion       |
-- | N4        | wake-da               | N3        | wake-da-result        |
-- | N2        | wakeda                | N2        | wake-da-nuance        |
-- | N3        | monono                | N2        | mono-no               |
-- | N2        | monono-formal         | N2        | mono-no               |
-- | N5        | nagara-simultaneous   | N4        | nagara                |

-- Rollback plan (inverse of this migration): reverse every UPDATE/DELETE path
-- or restore from corpus seed, which is the source-of-truth and will recreate
-- canonical row keys.

-- ------------------------------------------------------------
-- grammar_point
-- ------------------------------------------------------------
-- Strategy notes:
-- - Non-conflict pairs: pure UPDATE.
-- - Conflict pairs: update only when target slug/level is absent, otherwise
--   DELETE source to merge into the existing target row.
-- PK/UNIQUE risk is only on grammar_point.slug.

-- Non-conflicting moves (no target conflicts expected):
UPDATE grammar_point
SET slug='hazu-ga-nai', jlpt_level='N4'
WHERE slug='hazuganai' AND jlpt_level='N3';

UPDATE grammar_point
SET slug='mono-da-norm', jlpt_level='N3'
WHERE slug='mono-da' AND jlpt_level='N4';

UPDATE grammar_point
SET slug='mono-da-emotion', jlpt_level='N2'
WHERE slug='monoda' AND jlpt_level='N2';

UPDATE grammar_point
SET slug='wake-da-result', jlpt_level='N3'
WHERE slug='wake-da' AND jlpt_level='N4';

UPDATE grammar_point
SET slug='wake-da-nuance', jlpt_level='N2'
WHERE slug='wakeda' AND jlpt_level='N2';

-- Conflict merge targets: apply deletes + guarded updates last.
-- 1) hazuda -> hazu-da (N3 -> N4) collides with existing N4/hazu-da
UPDATE grammar_point
SET slug='hazu-da', jlpt_level='N4'
WHERE slug='hazuda' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='hazu-da' AND jlpt_level='N4'
  );
DELETE FROM grammar_point
WHERE slug='hazuda' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='hazu-da' AND jlpt_level='N4'
  );

-- 2) kamoshirenai -> kamo-shirenai (N3 -> N4)
UPDATE grammar_point
SET slug='kamo-shirenai', jlpt_level='N4'
WHERE slug='kamoshirenai' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='kamo-shirenai' AND jlpt_level='N4'
  );
DELETE FROM grammar_point
WHERE slug='kamoshirenai' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='kamo-shirenai' AND jlpt_level='N4'
  );

-- 3) teshimau -> te-shimau (N3 -> N4)
UPDATE grammar_point
SET slug='te-shimau', jlpt_level='N4'
WHERE slug='teshimau' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='te-shimau' AND jlpt_level='N4'
  );
DELETE FROM grammar_point
WHERE slug='teshimau' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='te-shimau' AND jlpt_level='N4'
  );

-- 4) monono -> mono-no (N3 -> N2)
UPDATE grammar_point
SET slug='mono-no', jlpt_level='N2'
WHERE slug='monono' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no' AND jlpt_level='N2'
  );
DELETE FROM grammar_point
WHERE slug='monono' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no' AND jlpt_level='N2'
  );

-- 5) monono-formal -> mono-no (N2 -> N2), same destination as #4
UPDATE grammar_point
SET slug='mono-no', jlpt_level='N2'
WHERE slug='monono-formal' AND jlpt_level='N2'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no' AND jlpt_level='N2'
  );
DELETE FROM grammar_point
WHERE slug='monono-formal' AND jlpt_level='N2'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no' AND jlpt_level='N2'
  );

-- 6) nagara-simultaneous -> nagara (N5 -> N4)
UPDATE grammar_point
SET slug='nagara', jlpt_level='N4'
WHERE slug='nagara-simultaneous' AND jlpt_level='N5'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='nagara' AND jlpt_level='N4'
  );
DELETE FROM grammar_point
WHERE slug='nagara-simultaneous' AND jlpt_level='N5'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='nagara' AND jlpt_level='N4'
  );

-- ------------------------------------------------------------
-- question (grammar_point+level references)
-- ------------------------------------------------------------
-- Strategy notes:
-- - Non-conflict rows: pure UPDATE filtered by content_type='grammar'.
-- - Conflict slug targets are merged conservatively: prefer existing target rows,
--   deleting source rows only when needed.

-- Non-conflicting moves:
UPDATE question
SET grammar_point='hazu-ga-nai', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='hazuganai' AND jlpt_level='N3';

UPDATE question
SET grammar_point='mono-da-norm', jlpt_level='N3'
WHERE content_type='grammar' AND grammar_point='mono-da' AND jlpt_level='N4';

UPDATE question
SET grammar_point='mono-da-emotion', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='monoda' AND jlpt_level='N2';

UPDATE question
SET grammar_point='wake-da-result', jlpt_level='N3'
WHERE content_type='grammar' AND grammar_point='wake-da' AND jlpt_level='N4';

UPDATE question
SET grammar_point='wake-da-nuance', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='wakeda' AND jlpt_level='N2';

-- Conflict-aware merges: delete source only if destination already exists.
UPDATE question
SET grammar_point='hazu-da', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='hazuda' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='hazu-da' AND jlpt_level='N4'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='hazuda' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='hazu-da' AND jlpt_level='N4'
  );

UPDATE question
SET grammar_point='kamo-shirenai', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='kamoshirenai' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='kamo-shirenai' AND jlpt_level='N4'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='kamoshirenai' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='kamo-shirenai' AND jlpt_level='N4'
  );

UPDATE question
SET grammar_point='te-shimau', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='teshimau' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='te-shimau' AND jlpt_level='N4'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='teshimau' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='te-shimau' AND jlpt_level='N4'
  );

UPDATE question
SET grammar_point='mono-no', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='monono' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='mono-no' AND jlpt_level='N2'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='monono' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='mono-no' AND jlpt_level='N2'
  );

UPDATE question
SET grammar_point='mono-no', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='monono-formal' AND jlpt_level='N2'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='mono-no' AND jlpt_level='N2'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='monono-formal' AND jlpt_level='N2'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='mono-no' AND jlpt_level='N2'
  );

UPDATE question
SET grammar_point='nagara', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='nagara-simultaneous' AND jlpt_level='N5'
  AND NOT EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='nagara' AND jlpt_level='N4'
  );
DELETE FROM question
WHERE content_type='grammar' AND grammar_point='nagara-simultaneous' AND jlpt_level='N5'
  AND EXISTS (
    SELECT 1 FROM question WHERE content_type='grammar' AND grammar_point='nagara' AND jlpt_level='N4'
  );

-- ------------------------------------------------------------
-- read_log (content_type='grammar' only)
-- ------------------------------------------------------------
-- Strategy notes:
-- - PK exists on (content_type, slug).
-- - Non-conflict pairs: pure UPDATE.
-- - Conflict pairs: merge remaining counters and keep existing target row,
--   then delete source.

-- Non-conflicting moves:
UPDATE read_log
SET slug='hazu-ga-nai'
WHERE content_type='grammar' AND slug='hazuganai';

UPDATE read_log
SET slug='mono-da-norm'
WHERE content_type='grammar' AND slug='mono-da';

UPDATE read_log
SET slug='mono-da-emotion'
WHERE content_type='grammar' AND slug='monoda';

UPDATE read_log
SET slug='wake-da-result'
WHERE content_type='grammar' AND slug='wake-da';

UPDATE read_log
SET slug='wake-da-nuance'
WHERE content_type='grammar' AND slug='wakeda';

-- Conflict-aware merges.
UPDATE read_log
SET slug='hazu-da'
WHERE content_type='grammar' AND slug='hazuda'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazu-da'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuda', 'hazu-da')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuda', 'hazu-da')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuda', 'hazu-da')
      )
WHERE content_type='grammar' AND slug='hazu-da'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazuda'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='hazuda'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazu-da'
  );

UPDATE read_log
SET slug='kamo-shirenai'
WHERE content_type='grammar' AND slug='kamoshirenai'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='kamo-shirenai'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('kamoshirenai', 'kamo-shirenai')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('kamoshirenai', 'kamo-shirenai')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('kamoshirenai', 'kamo-shirenai')
      )
WHERE content_type='grammar' AND slug='kamo-shirenai'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='kamoshirenai'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='kamoshirenai'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='kamo-shirenai'
  );

UPDATE read_log
SET slug='te-shimau'
WHERE content_type='grammar' AND slug='teshimau'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='te-shimau'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('teshimau', 'te-shimau')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('teshimau', 'te-shimau')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('teshimau', 'te-shimau')
      )
WHERE content_type='grammar' AND slug='te-shimau'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='teshimau'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='teshimau'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='te-shimau'
  );

UPDATE read_log
SET slug='mono-no'
WHERE content_type='grammar' AND slug='monono'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-no'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono', 'mono-no')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono', 'mono-no')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono', 'mono-no')
      )
WHERE content_type='grammar' AND slug='mono-no'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='monono'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='monono'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-no'
  );

UPDATE read_log
SET slug='mono-no'
WHERE content_type='grammar' AND slug='monono-formal'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-no'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono-formal', 'mono-no')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono-formal', 'mono-no')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monono-formal', 'mono-no')
      )
WHERE content_type='grammar' AND slug='mono-no'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='monono-formal'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='monono-formal'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-no'
  );

UPDATE read_log
SET slug='nagara'
WHERE content_type='grammar' AND slug='nagara-simultaneous'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='nagara'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('nagara-simultaneous', 'nagara')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('nagara-simultaneous', 'nagara')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('nagara-simultaneous', 'nagara')
      )
WHERE content_type='grammar' AND slug='nagara'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='nagara-simultaneous'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='nagara-simultaneous'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='nagara'
  );
