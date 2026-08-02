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
--
-- Preservation policy:
-- 1) question rows always rekey grammar_point + jlpt_level (id is independent PK;
--    never DELETE questions; attempts reference question.id with ON DELETE CASCADE).
-- 2) feedback_template rekeys grammar_point; on (grammar_point, error_class)
--    collision keep the destination row and drop only the obsolete source row.
-- 3) read_log merges counters when both old and new slugs exist, then drops source.
-- 4) grammar_point rekeys when destination slug is free; when destination already
--    exists (corpus already seeded canonical row), reassign grammar_example rows
--    from the source id onto the destination id, then drop only the obsolete
--    source grammar_point row (after questions/feedback/read_log are rekeyed).
--
-- Rollback / recovery (NOT "git revert alone"):
-- - Schema migrations are forward-only; git revert of app code does not reverse DB.
-- - Preferred recovery: restore a pre-migration SQLite backup, or rebuild DB via
--   `make seed` / `make seed-corpus` from L1 corpus (source of truth for slugs).
-- - Operator runbook: stop server → copy japanese-site.sqlite backup → start old
--   binary if needed. Documented in DECISIONS.md (2026-05-18 JS-114a entry).

-- ------------------------------------------------------------
-- 1) question — always rekey; never delete
-- ------------------------------------------------------------
UPDATE question
SET grammar_point='hazu-da', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='hazuda' AND jlpt_level='N3';

UPDATE question
SET grammar_point='hazu-ga-nai', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='hazuganai' AND jlpt_level='N3';

UPDATE question
SET grammar_point='kamo-shirenai', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='kamoshirenai' AND jlpt_level='N3';

UPDATE question
SET grammar_point='te-shimau', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='teshimau' AND jlpt_level='N3';

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

UPDATE question
SET grammar_point='mono-no', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='monono' AND jlpt_level='N3';

UPDATE question
SET grammar_point='mono-no', jlpt_level='N2'
WHERE content_type='grammar' AND grammar_point='monono-formal' AND jlpt_level='N2';

UPDATE question
SET grammar_point='nagara', jlpt_level='N4'
WHERE content_type='grammar' AND grammar_point='nagara-simultaneous' AND jlpt_level='N5';

-- ------------------------------------------------------------
-- 2) feedback_template — conflict-safe rekey
-- ------------------------------------------------------------
-- Map each old slug → new slug. When destination already has the same
-- error_class, keep destination and delete the obsolete source row only.

-- hazuda → hazu-da
UPDATE feedback_template
SET grammar_point='hazu-da'
WHERE grammar_point='hazuda'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='hazu-da' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='hazuda';

-- hazuganai → hazu-ga-nai
UPDATE feedback_template
SET grammar_point='hazu-ga-nai'
WHERE grammar_point='hazuganai'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='hazu-ga-nai' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='hazuganai';

-- kamoshirenai → kamo-shirenai
UPDATE feedback_template
SET grammar_point='kamo-shirenai'
WHERE grammar_point='kamoshirenai'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='kamo-shirenai' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='kamoshirenai';

-- teshimau → te-shimau
UPDATE feedback_template
SET grammar_point='te-shimau'
WHERE grammar_point='teshimau'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='te-shimau' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='teshimau';

-- mono-da → mono-da-norm
UPDATE feedback_template
SET grammar_point='mono-da-norm'
WHERE grammar_point='mono-da'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='mono-da-norm' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='mono-da';

-- monoda → mono-da-emotion
UPDATE feedback_template
SET grammar_point='mono-da-emotion'
WHERE grammar_point='monoda'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='mono-da-emotion' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='monoda';

-- wake-da → wake-da-result
UPDATE feedback_template
SET grammar_point='wake-da-result'
WHERE grammar_point='wake-da'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='wake-da-result' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='wake-da';

-- wakeda → wake-da-nuance
UPDATE feedback_template
SET grammar_point='wake-da-nuance'
WHERE grammar_point='wakeda'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='wake-da-nuance' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='wakeda';

-- monono → mono-no
UPDATE feedback_template
SET grammar_point='mono-no'
WHERE grammar_point='monono'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='mono-no' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='monono';

-- monono-formal → mono-no
UPDATE feedback_template
SET grammar_point='mono-no'
WHERE grammar_point='monono-formal'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='mono-no' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='monono-formal';

-- nagara-simultaneous → nagara
UPDATE feedback_template
SET grammar_point='nagara'
WHERE grammar_point='nagara-simultaneous'
  AND NOT EXISTS (
    SELECT 1 FROM feedback_template d
    WHERE d.grammar_point='nagara' AND d.error_class=feedback_template.error_class
  );
DELETE FROM feedback_template WHERE grammar_point='nagara-simultaneous';

-- ------------------------------------------------------------
-- 3) read_log — rekey or merge counters then drop source
-- ------------------------------------------------------------
-- hazuda → hazu-da
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

-- hazuganai → hazu-ga-nai
UPDATE read_log
SET slug='hazu-ga-nai'
WHERE content_type='grammar' AND slug='hazuganai'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazu-ga-nai'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuganai', 'hazu-ga-nai')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuganai', 'hazu-ga-nai')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('hazuganai', 'hazu-ga-nai')
      )
WHERE content_type='grammar' AND slug='hazu-ga-nai'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazuganai'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='hazuganai'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='hazu-ga-nai'
  );

-- kamoshirenai → kamo-shirenai
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

-- teshimau → te-shimau
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

-- mono-da → mono-da-norm
UPDATE read_log
SET slug='mono-da-norm'
WHERE content_type='grammar' AND slug='mono-da'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-da-norm'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('mono-da', 'mono-da-norm')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('mono-da', 'mono-da-norm')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('mono-da', 'mono-da-norm')
      )
WHERE content_type='grammar' AND slug='mono-da-norm'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-da'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='mono-da'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-da-norm'
  );

-- monoda → mono-da-emotion
UPDATE read_log
SET slug='mono-da-emotion'
WHERE content_type='grammar' AND slug='monoda'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-da-emotion'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monoda', 'mono-da-emotion')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monoda', 'mono-da-emotion')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('monoda', 'mono-da-emotion')
      )
WHERE content_type='grammar' AND slug='mono-da-emotion'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='monoda'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='monoda'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='mono-da-emotion'
  );

-- wake-da → wake-da-result
UPDATE read_log
SET slug='wake-da-result'
WHERE content_type='grammar' AND slug='wake-da'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wake-da-result'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wake-da', 'wake-da-result')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wake-da', 'wake-da-result')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wake-da', 'wake-da-result')
      )
WHERE content_type='grammar' AND slug='wake-da-result'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wake-da'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='wake-da'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wake-da-result'
  );

-- wakeda → wake-da-nuance
UPDATE read_log
SET slug='wake-da-nuance'
WHERE content_type='grammar' AND slug='wakeda'
  AND NOT EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wake-da-nuance'
  );
UPDATE read_log
SET read_count = (
        SELECT COALESCE(SUM(read_count), 0)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wakeda', 'wake-da-nuance')
      ),
    first_read_at = (
        SELECT MIN(first_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wakeda', 'wake-da-nuance')
      ),
    last_read_at = (
        SELECT MAX(last_read_at)
        FROM read_log
        WHERE content_type='grammar' AND slug IN ('wakeda', 'wake-da-nuance')
      )
WHERE content_type='grammar' AND slug='wake-da-nuance'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wakeda'
  );
DELETE FROM read_log
WHERE content_type='grammar' AND slug='wakeda'
  AND EXISTS (
    SELECT 1 FROM read_log WHERE content_type='grammar' AND slug='wake-da-nuance'
  );

-- monono → mono-no
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

-- monono-formal → mono-no
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

-- nagara-simultaneous → nagara
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

-- ------------------------------------------------------------
-- 4) grammar_point — rekey when free; drop obsolete source only on collision
-- ------------------------------------------------------------
-- Non-conflicting renames (destination slug not expected to pre-exist as old key):
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

-- Conflict-aware (destination may already exist from corpus seed).
-- Before dropping a source grammar_point row, reassign its grammar_example
-- children to the destination id so learner-visible examples are preserved.

-- hazuda → hazu-da
UPDATE grammar_point
SET slug='hazu-da', jlpt_level='N4'
WHERE slug='hazuda' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='hazu-da'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='hazu-da' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='hazuda' AND jlpt_level='N3'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='hazu-da');
DELETE FROM grammar_point
WHERE slug='hazuda' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='hazu-da'
  );

-- kamoshirenai → kamo-shirenai
UPDATE grammar_point
SET slug='kamo-shirenai', jlpt_level='N4'
WHERE slug='kamoshirenai' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='kamo-shirenai'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='kamo-shirenai' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='kamoshirenai' AND jlpt_level='N3'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='kamo-shirenai');
DELETE FROM grammar_point
WHERE slug='kamoshirenai' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='kamo-shirenai'
  );

-- teshimau → te-shimau
UPDATE grammar_point
SET slug='te-shimau', jlpt_level='N4'
WHERE slug='teshimau' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='te-shimau'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='te-shimau' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='teshimau' AND jlpt_level='N3'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='te-shimau');
DELETE FROM grammar_point
WHERE slug='teshimau' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='te-shimau'
  );

-- monono → mono-no
UPDATE grammar_point
SET slug='mono-no', jlpt_level='N2'
WHERE slug='monono' AND jlpt_level='N3'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='mono-no' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='monono' AND jlpt_level='N3'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='mono-no');
DELETE FROM grammar_point
WHERE slug='monono' AND jlpt_level='N3'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no'
  );

-- monono-formal → mono-no
UPDATE grammar_point
SET slug='mono-no', jlpt_level='N2'
WHERE slug='monono-formal' AND jlpt_level='N2'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='mono-no' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='monono-formal' AND jlpt_level='N2'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='mono-no');
DELETE FROM grammar_point
WHERE slug='monono-formal' AND jlpt_level='N2'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='mono-no'
  );

-- nagara-simultaneous → nagara
UPDATE grammar_point
SET slug='nagara', jlpt_level='N4'
WHERE slug='nagara-simultaneous' AND jlpt_level='N5'
  AND NOT EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='nagara'
  );
UPDATE grammar_example
SET grammar_point_id = (SELECT id FROM grammar_point WHERE slug='nagara' LIMIT 1)
WHERE grammar_point_id IN (
    SELECT id FROM grammar_point WHERE slug='nagara-simultaneous' AND jlpt_level='N5'
  )
  AND EXISTS (SELECT 1 FROM grammar_point WHERE slug='nagara');
DELETE FROM grammar_point
WHERE slug='nagara-simultaneous' AND jlpt_level='N5'
  AND EXISTS (
    SELECT 1 FROM grammar_point WHERE slug='nagara'
  );
