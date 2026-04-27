# Attribution

japanese-site uses several public language-learning datasets. All are
imported at build time from `server/data/external.lock` (not committed to
this repo). When you run `make bootstrap` or `make seed-all`, the seed
pipeline downloads the latest pinned versions and verifies their sha256.

## Imported datasets

### JMdict (jmdict-simplified)

- Upstream: https://github.com/scriptin/jmdict-simplified
- Original source: https://www.edrdg.org/jmdict/edict_doc.html
- License: **CC-BY-SA-4.0** (inherited from JMdict)
- Used for: `vocab` table (~22.5k common words)
- Attribution: This product uses JMdict files. These files are the property
  of the Electronic Dictionary Research and Development Group, and are used
  in conformance with the Group's licence
  (https://www.edrdg.org/edrdg/licence.html).

### KANJIDIC2 (via jmdict-simplified)

- Upstream: https://github.com/scriptin/jmdict-simplified
- Original source: https://www.edrdg.org/wiki/index.php/KANJIDIC_Project
- License: **CC-BY-SA-4.0** (inherited from KANJIDIC2)
- Used for: `kanji` table (~10.4k kanji)
- Attribution: This product uses the KANJIDIC dictionary file. This file is
  the property of the Electronic Dictionary Research and Development Group,
  and is used in conformance with the Group's licence.

### Tatoeba

- Upstream: https://downloads.tatoeba.org/exports/per_language/jpn/
- License: **CC-BY-2.0-FR**
- Used for: `sentence` table (~248k Japanese sentences)
- Attribution: Sentences from the Tatoeba project (https://tatoeba.org).

### JLPT vocabulary overlay

- Upstream: https://github.com/jamsinclair/open-anki-jlpt-decks
- License: **MIT**
- Used for: `vocab.jlpt_level` overlay (N5–N1)

## Curated content

`server/data/corpus/grammar/**` is original content authored for this project.
Released under **CC-BY-SA-4.0** to remain compatible with the imported
dictionary data.

If you reuse content from `corpus/grammar/`, please credit:
"japanese-site (https://github.com/screenleon/japanese-site) — CC-BY-SA-4.0"

## Source code

`LICENSE` (root) covers the source code of japanese-site (Apache-2.0).
Curated content (above) and imported datasets (above) keep their own licenses.
