# JS-016 Jisho JLPT 等級稽核（2026-05-03）

資料來源：`server/data/corpus/vocab/N1.jsonl` 至 `N5.jsonl`；目前等級依檔名判定。Jisho 等級取 `https://jisho.org/api/v1/search/words?keyword=<headword>` 第一筆結果的 `jlpt` 陣列；空陣列或無第一筆結果記為 `unlisted`。

> 注意：任務描述寫 590 筆，但本次實際讀取 N1-N5 JSONL 共 690 筆。

**摘要**

| 指標 | 數量 |
|---|---:|
| 全部 vocab rows | 690 |
| 非人工覆蓋、已查 Jisho rows | 671 |
| 等級相符 | 581 |
| 等級差異 | 81 |
| Jisho 無收錄 / 無 JLPT tag | 9 |
| 已人工覆蓋 | 19 |

## 等級相符

| 目前等級 | Headword | Reading | Jisho tag | Jisho URL |
|---|---|---|---|---|
| N1 | あくどい | あくどい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%82%E3%81%8F%E3%81%A9%E3%81%84) |
| N1 | あやふや | あやふや | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%82%E3%82%84%E3%81%B5%E3%82%84) |
| N1 | いっそ | いっそ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%84%E3%81%A3%E3%81%9D) |
| N1 | 一切 | いっさい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E5%88%87) |
| N1 | 一律 | いちりつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E5%BE%8B) |
| N1 | 一括 | いっかつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%8B%AC) |
| N1 | 一挙に | いっきょに | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%8C%99%E3%81%AB) |
| N1 | 一概に | いちがいに | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%A6%82%E3%81%AB) |
| N1 | 一連 | いちれん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E9%80%A3) |
| N1 | 上位 | じょうい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E4%BD%8D) |
| N1 | 上回る | うわまわる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E5%9B%9E%E3%82%8B) |
| N1 | 不可欠 | ふかけつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8D%E5%8F%AF%E6%AC%A0) |
| N1 | 不当 | ふとう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8D%E5%BD%93) |
| N1 | 不振 | ふしん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8D%E6%8C%AF) |
| N1 | 不況 | ふきょう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8D%E6%B3%81) |
| N1 | 与党 | よとう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8E%E5%85%9A) |
| N1 | 世論 | よろん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%96%E8%AB%96) |
| N1 | 両立 | りょうりつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%A1%E7%AB%8B) |
| N1 | 中枢 | ちゅうすう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%AD%E6%9E%A2) |
| N1 | 主体 | しゅたい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%BB%E4%BD%93) |
| N1 | 主導 | しゅどう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%BB%E5%B0%8E) |
| N1 | 主権 | しゅけん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%BB%E6%A8%A9) |
| N1 | 予め | あらかじめ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%88%E3%82%81) |
| N1 | 交付 | こうふ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A4%E4%BB%98) |
| N1 | 交渉 | こうしょう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A4%E6%B8%89) |
| N1 | 今更 | いまさら | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E6%9B%B4) |
| N1 | 任務 | にんむ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%BB%E5%8B%99) |
| N1 | 任命 | にんめい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%BB%E5%91%BD) |
| N1 | 企画 | きかく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%81%E7%94%BB) |
| N1 | 会見 | かいけん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E8%A6%8B) |
| N1 | 会談 | かいだん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E8%AB%87) |
| N1 | 伝達 | でんたつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9D%E9%81%94) |
| N1 | 伴う | ともなう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%B4%E3%81%86) |
| N1 | 余地 | よち | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%99%E5%9C%B0) |
| N1 | 使命 | しめい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%BF%E5%91%BD) |
| N1 | 依存 | いぞん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BE%9D%E5%AD%98) |
| N1 | 侵略 | しんりゃく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BE%B5%E7%95%A5) |
| N1 | 促す | うながす | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BF%83%E3%81%99) |
| N1 | 促進 | そくしん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BF%83%E9%80%B2) |
| N1 | 保守 | ほしゅ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BF%9D%E5%AE%88) |
| N1 | 信任 | しんにん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BF%A1%E4%BB%BB) |
| N1 | 個別 | こべつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%8B%E5%88%A5) |
| N1 | 値する | あたいする | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%A4%E3%81%99%E3%82%8B) |
| N1 | 偏見 | へんけん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%81%8F%E8%A6%8B) |
| N1 | 停滞 | ていたい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%81%9C%E6%BB%9E) |
| N1 | 健全 | けんぜん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%81%A5%E5%85%A8) |
| N1 | 偽造 | ぎぞう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%81%BD%E9%80%A0) |
| N1 | 優位 | ゆうい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%84%AA%E4%BD%8D) |
| N1 | 優先 | ゆうせん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%84%AA%E5%85%88) |
| N1 | 免除 | めんじょ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%8D%E9%99%A4) |
| N1 | 入手 | にゅうしゅ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%A5%E6%89%8B) |
| N1 | 公募 | こうぼ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AC%E5%8B%9F) |
| N1 | 公認 | こうにん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AC%E8%AA%8D) |
| N1 | 共存 | きょうぞん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%B1%E5%AD%98) |
| N1 | 内閣 | ないかく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%85%E9%96%A3) |
| N1 | 円滑 | えんかつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%86%E6%BB%91) |
| N1 | 再建 | さいけん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%8D%E5%BB%BA) |
| N1 | 冒頭 | ぼうとう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%92%E9%A0%AD) |
| N1 | 処分 | しょぶん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%A6%E5%88%86) |
| N1 | 判決 | はんけつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%A4%E6%B1%BA) |
| N1 | 到達 | とうたつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%B0%E9%81%94) |
| N1 | 制定 | せいてい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%B6%E5%AE%9A) |
| N1 | 制約 | せいやく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%B6%E7%B4%84) |
| N1 | 制裁 | せいさい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%B6%E8%A3%81) |
| N1 | 削減 | さくげん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%89%8A%E6%B8%9B) |
| N1 | 前提 | ぜんてい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%89%8D%E6%8F%90) |
| N1 | 労る | いたわる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8A%B4%E3%82%8B) |
| N1 | 効率 | こうりつ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8A%B9%E7%8E%87) |
| N1 | 勧告 | かんこく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8B%A7%E5%91%8A) |
| N1 | 卑しい | いやしい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%91%E3%81%97%E3%81%84) |
| N1 | 単調 | たんちょう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%98%E8%AA%BF) |
| N1 | 危ぶむ | あやぶむ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%B1%E3%81%B6%E3%82%80) |
| N1 | 営む | いとなむ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%96%B6%E3%82%80) |
| N1 | 妥協 | だきょう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A6%A5%E5%8D%94) |
| N1 | 委託 | いたく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A7%94%E8%A8%97) |
| N1 | 小売 | こうり | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%B0%8F%E5%A3%B2) |
| N1 | 心強い | こころづよい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BF%83%E5%BC%B7%E3%81%84) |
| N1 | 心細い | こころぼそい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BF%83%E7%B4%B0%E3%81%84) |
| N1 | 志す | こころざす | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BF%97%E3%81%99) |
| N1 | 意向 | いこう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%84%8F%E5%90%91) |
| N1 | 意図 | いと | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%84%8F%E5%9B%B3) |
| N1 | 意気込む | いきごむ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%84%8F%E6%B0%97%E8%BE%BC%E3%82%80) |
| N1 | 憧れ | あこがれ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%86%A7%E3%82%8C) |
| N1 | 打開 | だかい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%89%93%E9%96%8B) |
| N1 | 把握 | はあく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8A%8A%E6%8F%A1) |
| N1 | 挑む | いどむ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8C%91%E3%82%80) |
| N1 | 捗る | はかどる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8D%97%E3%82%8B) |
| N1 | 採択 | さいたく | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8E%A1%E6%8A%9E) |
| N1 | 採算 | さいさん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8E%A1%E7%AE%97) |
| N1 | 敢えて | あえて | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%95%A2%E3%81%88%E3%81%A6) |
| N1 | 断然 | だんぜん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%96%AD%E7%84%B6) |
| N1 | 最善 | さいぜん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9C%80%E5%96%84) |
| N1 | 有りのまま | ありのまま | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9C%89%E3%82%8A%E3%81%AE%E3%81%BE%E3%81%BE) |
| N1 | 案じる | あんじる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%A1%88%E3%81%98%E3%82%8B) |
| N1 | 浅ましい | あさましい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B5%85%E3%81%BE%E3%81%97%E3%81%84) |
| N1 | 滞る | とどこおる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%BB%9E%E3%82%8B) |
| N1 | 焦る | あせる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%84%A6%E3%82%8B) |
| N1 | 異論 | いろん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%95%B0%E8%AB%96) |
| N1 | 着手 | ちゃくしゅ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%9D%80%E6%89%8B) |
| N1 | 著しい | いちじるしい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%91%97%E3%81%97%E3%81%84) |
| N1 | 購買 | こうばい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%B3%BC%E8%B2%B7) |
| N1 | 遂げる | とげる | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%82%E3%81%92%E3%82%8B) |
| N1 | 運営 | うんえい | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8B%E5%96%B6) |
| N1 | 運用 | うんよう | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8B%E7%94%A8) |
| N1 | 配慮 | はいりょ | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%85%8D%E6%85%AE) |
| N1 | 閲覧 | えつらん | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%96%B2%E8%A6%A7) |
| N1 | 陰気 | いんき | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%99%B0%E6%B0%97) |
| N1 | 鮮やか | あざやか | jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%AE%AE%E3%82%84%E3%81%8B) |
| N2 | うろうろ | うろうろ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%86%E3%82%8D%E3%81%86%E3%82%8D) |
| N2 | ぎっしり | ぎっしり | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8E%E3%81%A3%E3%81%97%E3%82%8A) |
| N2 | こっそり | こっそり | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%93%E3%81%A3%E3%81%9D%E3%82%8A) |
| N2 | ご苦労様 | ごくろうさま | jlpt-n2, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%94%E8%8B%A6%E5%8A%B4%E6%A7%98) |
| N2 | しゃがむ | しゃがむ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%97%E3%82%83%E3%81%8C%E3%82%80) |
| N2 | すっきり | すっきり | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%99%E3%81%A3%E3%81%8D%E3%82%8A) |
| N2 | ずらす | ずらす | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9A%E3%82%89%E3%81%99) |
| N2 | ずらり | ずらり | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9A%E3%82%89%E3%82%8A) |
| N2 | せっせと | せっせと | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9B%E3%81%A3%E3%81%9B%E3%81%A8) |
| N2 | そそっかしい | そそっかしい | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%81%9D%E3%81%A3%E3%81%8B%E3%81%97%E3%81%84) |
| N2 | その他 | そのほか | jlpt-n2, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%81%AE%E4%BB%96) |
| N2 | どうせ | どうせ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%81%86%E3%81%9B) |
| N2 | どっと | どっと | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%81%A3%E3%81%A8) |
| N2 | のろのろ | のろのろ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AE%E3%82%8D%E3%81%AE%E3%82%8D) |
| N2 | ふわふわ | ふわふわ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%B5%E3%82%8F%E3%81%B5%E3%82%8F) |
| N2 | ぶつぶつ | ぶつぶつ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%B6%E3%81%A4%E3%81%B6%E3%81%A4) |
| N2 | まあまあ | まあまあ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%BE%E3%81%82%E3%81%BE%E3%81%82) |
| N2 | まごまご | まごまご | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%BE%E3%81%94%E3%81%BE%E3%81%94) |
| N2 | めっきり | めっきり | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%81%E3%81%A3%E3%81%8D%E3%82%8A) |
| N2 | アクセント | アクセント | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%82%AF%E3%82%BB%E3%83%B3%E3%83%88) |
| N2 | アンテナ | アンテナ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%B3%E3%83%86%E3%83%8A) |
| N2 | イコール | イコール | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A4%E3%82%B3%E3%83%BC%E3%83%AB) |
| N2 | ウール | ウール | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A6%E3%83%BC%E3%83%AB) |
| N2 | エチケット | エチケット | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%83%81%E3%82%B1%E3%83%83%E3%83%88) |
| N2 | エプロン | エプロン | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%83%97%E3%83%AD%E3%83%B3) |
| N2 | オイル | オイル | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%82%A4%E3%83%AB) |
| N2 | オルガン | オルガン | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%83%AB%E3%82%AC%E3%83%B3) |
| N2 | オーケストラ | オーケストラ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%83%BC%E3%82%B1%E3%82%B9%E3%83%88%E3%83%A9) |
| N2 | オートメーション | オートメーション | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%83%BC%E3%83%88%E3%83%A1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3) |
| N2 | カセット | カセット | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%82%BB%E3%83%83%E3%83%88) |
| N2 | カバー | カバー | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%90%E3%83%BC) |
| N2 | カラー | カラー | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%A9%E3%83%BC) |
| N2 | カロリー | カロリー | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%AD%E3%83%AA%E3%83%BC) |
| N2 | カーブ | カーブ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%BC%E3%83%96) |
| N2 | ガム | ガム | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AC%E3%83%A0) |
| N2 | キャンパス | キャンパス | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AD%E3%83%A3%E3%83%B3%E3%83%91%E3%82%B9) |
| N2 | ギャング | ギャング | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AE%E3%83%A3%E3%83%B3%E3%82%B0) |
| N2 | クリーニング | クリーニング | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%8B%E3%83%B3%E3%82%B0) |
| N2 | クーラー | クーラー | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AF%E3%83%BC%E3%83%A9%E3%83%BC) |
| N2 | コック | コック | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%83%E3%82%AF) |
| N2 | コレクション | コレクション | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%AC%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3) |
| N2 | コンクール | コンクール | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%B3%E3%82%AF%E3%83%BC%E3%83%AB) |
| N2 | コンセント | コンセント | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%B3%E3%82%BB%E3%83%B3%E3%83%88) |
| N2 | コース | コース | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%BC%E3%82%B9) |
| N2 | コーラス | コーラス | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%BC%E3%83%A9%E3%82%B9) |
| N2 | サイレン | サイレン | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%82%A4%E3%83%AC%E3%83%B3) |
| N2 | サラリーマン | サラリーマン | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%83%A9%E3%83%AA%E3%83%BC%E3%83%9E%E3%83%B3) |
| N2 | サンプル | サンプル | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%83%B3%E3%83%97%E3%83%AB) |
| N2 | サークル | サークル | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%83%BC%E3%82%AF%E3%83%AB) |
| N2 | シャッター | シャッター | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%A3%E3%83%83%E3%82%BF%E3%83%BC) |
| N2 | ショップ | ショップ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%A7%E3%83%83%E3%83%97) |
| N2 | シリーズ | シリーズ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%AA%E3%83%BC%E3%82%BA) |
| N2 | シーツ | シーツ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%BC%E3%83%84) |
| N2 | ジャーナリスト | ジャーナリスト | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B8%E3%83%A3%E3%83%BC%E3%83%8A%E3%83%AA%E3%82%B9%E3%83%88) |
| N2 | スカーフ | スカーフ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%AB%E3%83%BC%E3%83%95) |
| N2 | スタート | スタート | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%BF%E3%83%BC%E3%83%88) |
| N2 | ステージ | ステージ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%86%E3%83%BC%E3%82%B8) |
| N2 | ストップ | ストップ | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%88%E3%83%83%E3%83%97) |
| N2 | スピーカー | スピーカー | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%94%E3%83%BC%E3%82%AB%E3%83%BC) |
| N2 | スマート | スマート | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%9E%E3%83%BC%E3%83%88) |
| N2 | スライド | スライド | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%A9%E3%82%A4%E3%83%89) |
| N2 | セメント | セメント | jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%BB%E3%83%A1%E3%83%B3%E3%83%88) |
| N3 | お喋り | おしゃべり | jlpt-n3, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%96%8B%E3%82%8A) |
| N3 | がっかり | がっかり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8C%E3%81%A3%E3%81%8B%E3%82%8A) |
| N3 | きちんと | きちんと | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8D%E3%81%A1%E3%82%93%E3%81%A8) |
| N3 | きつい | きつい | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8D%E3%81%A4%E3%81%84) |
| N3 | ぐっすり | ぐっすり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%90%E3%81%A3%E3%81%99%E3%82%8A) |
| N3 | こんなに | こんなに | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%93%E3%82%93%E3%81%AA%E3%81%AB) |
| N3 | さっぱり | さっぱり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%95%E3%81%A3%E3%81%B1%E3%82%8A) |
| N3 | ざっと | ざっと | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%96%E3%81%A3%E3%81%A8) |
| N3 | ずっと | ずっと | jlpt-n3, jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9A%E3%81%A3%E3%81%A8) |
| N3 | そっくり | そっくり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%81%A3%E3%81%8F%E3%82%8A) |
| N3 | たった | たった | jlpt-n3, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9F%E3%81%A3%E3%81%9F) |
| N3 | たっぷり | たっぷり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9F%E3%81%A3%E3%81%B7%E3%82%8A) |
| N3 | ちゃんと | ちゃんと | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A1%E3%82%83%E3%82%93%E3%81%A8) |
| N3 | どんなに | どんなに | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%82%93%E3%81%AA%E3%81%AB) |
| N3 | にっこり | にっこり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AB%E3%81%A3%E3%81%93%E3%82%8A) |
| N3 | のんびり | のんびり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AE%E3%82%93%E3%81%B3%E3%82%8A) |
| N3 | ぴったり | ぴったり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%B4%E3%81%A3%E3%81%9F%E3%82%8A) |
| N3 | ぼんやり | ぼんやり | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%BC%E3%82%93%E3%82%84%E3%82%8A) |
| N3 | エネルギー | エネルギー | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%83%8D%E3%83%AB%E3%82%AE%E3%83%BC) |
| N3 | エンジン | エンジン | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%83%B3%E3%82%B8%E3%83%B3) |
| N3 | オフィス | オフィス | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%83%95%E3%82%A3%E3%82%B9) |
| N3 | キャプテン | キャプテン | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AD%E3%83%A3%E3%83%97%E3%83%86%E3%83%B3) |
| N3 | クラシック | クラシック | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AF%E3%83%A9%E3%82%B7%E3%83%83%E3%82%AF) |
| N3 | グループ | グループ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B0%E3%83%AB%E3%83%BC%E3%83%97) |
| N3 | ケース | ケース | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B1%E3%83%BC%E3%82%B9) |
| N3 | コーチ | コーチ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%BC%E3%83%81) |
| N3 | ゴール | ゴール | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B4%E3%83%BC%E3%83%AB) |
| N3 | サービス | サービス | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%83%BC%E3%83%93%E3%82%B9) |
| N3 | スタイル | スタイル | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%BF%E3%82%A4%E3%83%AB) |
| N3 | スピーチ | スピーチ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%94%E3%83%BC%E3%83%81) |
| N3 | チャンス | チャンス | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%81%E3%83%A3%E3%83%B3%E3%82%B9) |
| N3 | チーム | チーム | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%81%E3%83%BC%E3%83%A0) |
| N3 | トンネル | トンネル | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%88%E3%83%B3%E3%83%8D%E3%83%AB) |
| N3 | ドラマ | ドラマ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%89%E3%83%A9%E3%83%9E) |
| N3 | パスポート | パスポート | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%91%E3%82%B9%E3%83%9D%E3%83%BC%E3%83%88) |
| N3 | パーセント | パーセント | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%91%E3%83%BC%E3%82%BB%E3%83%B3%E3%83%88) |
| N3 | ブレーキ | ブレーキ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%96%E3%83%AC%E3%83%BC%E3%82%AD) |
| N3 | プラス | プラス | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%97%E3%83%A9%E3%82%B9) |
| N3 | プラン | プラン | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%97%E3%83%A9%E3%83%B3) |
| N3 | プロ | プロ | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%97%E3%83%AD) |
| N3 | マイク | マイク | jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9E%E3%82%A4%E3%82%AF) |
| N4 | あんな | あんな | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%82%E3%82%93%E3%81%AA) |
| N4 | おかしい | おかしい | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E3%81%8B%E3%81%97%E3%81%84) |
| N4 | お土産 | おみやげ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%9C%9F%E7%94%A3) |
| N4 | お嬢さん | おじょうさん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%AC%A2%E3%81%95%E3%82%93) |
| N4 | お宅 | おたく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%AE%85) |
| N4 | お礼 | おれい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E7%A4%BC) |
| N4 | お祝い | おいわい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E7%A5%9D%E3%81%84) |
| N4 | お見舞い | おみまい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E8%A6%8B%E8%88%9E%E3%81%84) |
| N4 | けんか | けんか | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%91%E3%82%93%E3%81%8B) |
| N4 | この間 | このあいだ | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%93%E3%81%AE%E9%96%93) |
| N4 | この頃 | このごろ | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%93%E3%81%AE%E9%A0%83) |
| N4 | ご主人 | ごしゅじん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%94%E4%B8%BB%E4%BA%BA) |
| N4 | ご存じ | ごぞんじ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%94%E5%AD%98%E3%81%98) |
| N4 | ご覧になる | ごらんになる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%94%E8%A6%A7%E3%81%AB%E3%81%AA%E3%82%8B) |
| N4 | すっかり | すっかり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%99%E3%81%A3%E3%81%8B%E3%82%8A) |
| N4 | すっと | すっと | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%99%E3%81%A3%E3%81%A8) |
| N4 | すると | すると | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%99%E3%82%8B%E3%81%A8) |
| N4 | そろそろ | そろそろ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%82%8D%E3%81%9D%E3%82%8D) |
| N4 | そんな | そんな | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%82%93%E3%81%AA) |
| N4 | そんなに | そんなに | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%82%93%E3%81%AA%E3%81%AB) |
| N4 | だから | だから | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A0%E3%81%8B%E3%82%89) |
| N4 | どんどん | どんどん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%82%93%E3%81%A9%E3%82%93) |
| N4 | はっきり | はっきり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AF%E3%81%A3%E3%81%8D%E3%82%8A) |
| N4 | アクセサリー | アクセサリー | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%82%AF%E3%82%BB%E3%82%B5%E3%83%AA%E3%83%BC) |
| N4 | アジア | アジア | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%82%B8%E3%82%A2) |
| N4 | アナウンサー | アナウンサー | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%8A%E3%82%A6%E3%83%B3%E3%82%B5%E3%83%BC) |
| N4 | アフリカ | アフリカ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%95%E3%83%AA%E3%82%AB) |
| N4 | アメリカ | アメリカ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%A1%E3%83%AA%E3%82%AB) |
| N4 | アルコール | アルコール | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%AB%E3%82%B3%E3%83%BC%E3%83%AB) |
| N4 | アルバイト | アルバイト | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%AB%E3%83%90%E3%82%A4%E3%83%88) |
| N4 | エスカレーター | エスカレーター | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%82%B9%E3%82%AB%E3%83%AC%E3%83%BC%E3%82%BF%E3%83%BC) |
| N4 | オートバイ | オートバイ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AA%E3%83%BC%E3%83%88%E3%83%90%E3%82%A4) |
| N4 | カーテン | カーテン | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%BC%E3%83%86%E3%83%B3) |
| N4 | ガソリン | ガソリン | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AC%E3%82%BD%E3%83%AA%E3%83%B3) |
| N4 | ガソリンスタンド | ガソリンスタンド | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AC%E3%82%BD%E3%83%AA%E3%83%B3%E3%82%B9%E3%82%BF%E3%83%B3%E3%83%89) |
| N4 | ガラス | ガラス | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AC%E3%83%A9%E3%82%B9) |
| N4 | コンサート | コンサート | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%B3%E3%82%B5%E3%83%BC%E3%83%88) |
| N4 | サンダル | サンダル | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B5%E3%83%B3%E3%83%80%E3%83%AB) |
| N4 | ジャム | ジャム | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B8%E3%83%A3%E3%83%A0) |
| N4 | スクリーン | スクリーン | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3) |
| N4 | ステレオ | ステレオ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%86%E3%83%AC%E3%82%AA) |
| N4 | ステーキ | ステーキ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%86%E3%83%BC%E3%82%AD) |
| N4 | スーツ | スーツ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%BC%E3%83%84) |
| N4 | スーツケース | スーツケース | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%BC%E3%83%84%E3%82%B1%E3%83%BC%E3%82%B9) |
| N4 | タイプ | タイプ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%BF%E3%82%A4%E3%83%97) |
| N4 | テキスト | テキスト | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88) |
| N4 | パソコン | パソコン | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%91%E3%82%BD%E3%82%B3%E3%83%B3) |
| N4 | ビル | ビル | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%93%E3%83%AB) |
| N4 | レジ | レジ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%AC%E3%82%B8) |
| N4 | 一度 | いちど | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E5%BA%A6) |
| N4 | 一杯 | いっぱい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%9D%AF) |
| N4 | 一生懸命 | いっしょうけんめい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E7%94%9F%E6%87%B8%E5%91%BD) |
| N4 | 丁寧 | ていねい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%81%E5%AF%A7) |
| N4 | 上がる | あがる | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E3%81%8C%E3%82%8B) |
| N4 | 下げる | さげる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B%E3%81%92%E3%82%8B) |
| N4 | 下宿 | げしゅく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B%E5%AE%BF) |
| N4 | 下着 | したぎ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B%E7%9D%80) |
| N4 | 不便 | ふべん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8D%E4%BE%BF) |
| N4 | 世界 | せかい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%96%E7%95%8C) |
| N4 | 両方 | りょうほう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%A1%E6%96%B9) |
| N4 | 中学校 | ちゅうがっこう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%AD%E5%AD%A6%E6%A0%A1) |
| N4 | 久しぶり | ひさしぶり | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%85%E3%81%97%E3%81%B6%E3%82%8A) |
| N4 | 乗り換える | のりかえる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%97%E3%82%8A%E6%8F%9B%E3%81%88%E3%82%8B) |
| N4 | 乗り物 | のりもの | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%97%E3%82%8A%E7%89%A9) |
| N4 | 乾く | かわく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%BE%E3%81%8F) |
| N4 | 予定 | よてい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%88%E5%AE%9A) |
| N4 | 予約 | よやく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%88%E7%B4%84) |
| N4 | 予習 | よしゅう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%88%E7%BF%92) |
| N4 | 事務所 | じむしょ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8B%E5%8B%99%E6%89%80) |
| N4 | 事故 | じこ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8B%E6%95%85) |
| N4 | 亡くなる | なくなる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A1%E3%81%8F%E3%81%AA%E3%82%8B) |
| N4 | 交通 | こうつう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A4%E9%80%9A) |
| N4 | 人口 | じんこう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%BA%E5%8F%A3) |
| N4 | 人形 | にんぎょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%BA%E5%BD%A2) |
| N4 | 今夜 | こんや | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E5%A4%9C) |
| N4 | 今度 | こんど | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E5%BA%A6) |
| N4 | 仕方 | しかた | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%95%E6%96%B9) |
| N4 | 代わり | かわり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%A3%E3%82%8F%E3%82%8A) |
| N4 | 以上 | いじょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%A5%E4%B8%8A) |
| N4 | 以下 | いか | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%A5%E4%B8%8B) |
| N4 | 以内 | いない | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%A5%E5%86%85) |
| N4 | 以外 | いがい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%A5%E5%A4%96) |
| N4 | 会場 | かいじょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E5%A0%B4) |
| N4 | 会話 | かいわ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E8%A9%B1) |
| N4 | 会議 | かいぎ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E8%AD%B0) |
| N4 | 会議室 | かいぎしつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E8%AD%B0%E5%AE%A4) |
| N4 | 伝える | つたえる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9D%E3%81%88%E3%82%8B) |
| N4 | 似る | にる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%BC%E3%82%8B) |
| N4 | 住所 | じゅうしょ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%8F%E6%89%80) |
| N4 | 例えば | たとえば | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BE%8B%E3%81%88%E3%81%B0) |
| N4 | 倍 | ばい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%8D) |
| N4 | 倒れる | たおれる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%92%E3%82%8C%E3%82%8B) |
| N4 | 値段 | ねだん | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%A4%E6%AE%B5) |
| N4 | 僕 | ぼく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%83%95) |
| N4 | 億 | おく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%84%84) |
| N4 | 優しい | やさしい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%84%AA%E3%81%97%E3%81%84) |
| N4 | 先輩 | せんぱい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%88%E8%BC%A9) |
| N4 | 光 | ひかり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%89) |
| N4 | 光る | ひかる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%89%E3%82%8B) |
| N4 | 公務員 | こうむいん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AC%E5%8B%99%E5%93%A1) |
| N4 | 具合 | ぐあい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%B7%E5%90%88) |
| N4 | 内 | うち | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%85) |
| N4 | 再来月 | さらいげつ | jlpt-n4, jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%8D%E6%9D%A5%E6%9C%88) |
| N4 | 再来週 | さらいしゅう | jlpt-n4, jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%8D%E6%9D%A5%E9%80%B1) |
| N4 | 写す | うつす | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%99%E3%81%99) |
| N4 | 冷える | ひえる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%B7%E3%81%88%E3%82%8B) |
| N4 | 冷房 | れいぼう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%B7%E6%88%BF) |
| N4 | 凄い | すごい | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%84%E3%81%84) |
| N4 | 別 | べつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%A5) |
| N4 | 別れる | わかれる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%A5%E3%82%8C%E3%82%8B) |
| N4 | 医学 | いがく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8C%BB%E5%AD%A6) |
| N4 | 危険 | きけん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%B1%E9%99%BA) |
| N4 | 原因 | げんいん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8E%9F%E5%9B%A0) |
| N4 | 受付 | うけつけ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%97%E4%BB%98) |
| N4 | 合う | あう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%88%E3%81%86) |
| N4 | 国際 | こくさい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%BD%E9%9A%9B) |
| N4 | 壊す | こわす | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A3%8A%E3%81%99) |
| N4 | 売り場 | うりば | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A3%B2%E3%82%8A%E5%A0%B4) |
| N4 | 安全 | あんぜん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%AE%89%E5%85%A8) |
| N4 | 安心 | あんしん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%AE%89%E5%BF%83) |
| N4 | 客 | きゃく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%AE%A2) |
| N4 | 屋上 | おくじょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%B1%8B%E4%B8%8A) |
| N4 | 工場 | こうじょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%B7%A5%E5%A0%B4) |
| N4 | 工業 | こうぎょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%B7%A5%E6%A5%AD) |
| N4 | 帰り | かえり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%B8%B0%E3%82%8A) |
| N4 | 彼 | かれ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BD%BC) |
| N4 | 彼女 | かのじょ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BD%BC%E5%A5%B3) |
| N4 | 怒る | おこる | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%80%92%E3%82%8B) |
| N4 | 急 | きゅう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%80%A5) |
| N4 | 急ぐ | いそぐ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%80%A5%E3%81%90) |
| N4 | 意見 | いけん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%84%8F%E8%A6%8B) |
| N4 | 打つ | うつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%89%93%E3%81%A4) |
| N4 | 掛ける | かける | jlpt-n4, jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%8E%9B%E3%81%91%E3%82%8B) |
| N4 | 故障 | こしょう | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%95%85%E9%9A%9C) |
| N4 | 教会 | きょうかい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%95%99%E4%BC%9A) |
| N4 | 教育 | きょういく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%95%99%E8%82%B2) |
| N4 | 景色 | けしき | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%99%AF%E8%89%B2) |
| N4 | 暮れる | くれる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9A%AE%E3%82%8C%E3%82%8B) |
| N4 | 枝 | えだ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9E%9D) |
| N4 | 校長 | こうちょう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%A0%A1%E9%95%B7) |
| N4 | 機会 | きかい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%A9%9F%E4%BC%9A) |
| N4 | 気持ち | きもち | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B0%97%E6%8C%81%E3%81%A1) |
| N4 | 決して | けっして | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B1%BA%E3%81%97%E3%81%A6) |
| N4 | 決まる | きまる | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B1%BA%E3%81%BE%E3%82%8B) |
| N4 | 決める | きめる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B1%BA%E3%82%81%E3%82%8B) |
| N4 | 海岸 | かいがん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B5%B7%E5%B2%B8) |
| N4 | 消しゴム | けしゴム | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%B6%88%E3%81%97%E3%82%B4%E3%83%A0) |
| N4 | 火事 | かじ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%81%AB%E4%BA%8B) |
| N4 | 着物 | きもの | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%9D%80%E7%89%A9) |
| N4 | 石 | いし | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%9F%B3) |
| N4 | 研究 | けんきゅう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%A0%94%E7%A9%B6) |
| N4 | 研究室 | けんきゅうしつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%A0%94%E7%A9%B6%E5%AE%A4) |
| N4 | 科学 | かがく | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%A7%91%E5%AD%A6) |
| N4 | 空気 | くうき | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%A9%BA%E6%B0%97) |
| N4 | 空港 | くうこう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%A9%BA%E6%B8%AF) |
| N4 | 競争 | きょうそう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%AB%B6%E4%BA%89) |
| N4 | 答え | こたえ | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%AD%94%E3%81%88) |
| N4 | 簡単 | かんたん | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%B0%A1%E5%8D%98) |
| N4 | 終わり | おわり | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%B5%82%E3%82%8F%E3%82%8A) |
| N4 | 美しい | うつくしい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%BE%8E%E3%81%97%E3%81%84) |
| N4 | 聞こえる | きこえる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%81%9E%E3%81%93%E3%81%88%E3%82%8B) |
| N4 | 腕 | うで | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%85%95) |
| N4 | 致す | いたす | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%87%B4%E3%81%99) |
| N4 | 興味 | きょうみ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%88%88%E5%91%B3) |
| N4 | 草 | くさ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%8D%89) |
| N4 | 落とす | おとす | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%90%BD%E3%81%A8%E3%81%99) |
| N4 | 行う | おこなう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%A1%8C%E3%81%86) |
| N4 | 裏 | うら | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%A3%8F) |
| N4 | 見物 | けんぶつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%A6%8B%E7%89%A9) |
| N4 | 講義 | こうぎ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%AC%9B%E7%BE%A9) |
| N4 | 謝る | あやまる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%AC%9D%E3%82%8B) |
| N4 | 警察 | けいさつ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%AD%A6%E5%AF%9F) |
| N4 | 贈り物 | おくりもの | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%B4%88%E3%82%8A%E7%89%A9) |
| N4 | 赤ん坊 | あかんぼう | jlpt-n4, jlpt-n2 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%B5%A4%E3%82%93%E5%9D%8A) |
| N4 | 起こす | おこす | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%B5%B7%E3%81%93%E3%81%99) |
| N4 | 踊る | おどる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%B8%8A%E3%82%8B) |
| N4 | 近所 | きんじょ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%BF%91%E6%89%80) |
| N4 | 遊び | あそび | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8A%E3%81%B3) |
| N4 | 運転手 | うんてんしゅ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8B%E8%BB%A2%E6%89%8B) |
| N4 | 遠慮 | えんりょ | jlpt-n4, jlpt-n3 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%A0%E6%85%AE) |
| N4 | 選ぶ | えらぶ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%B8%E3%81%B6) |
| N4 | 郊外 | こうがい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%83%8A%E5%A4%96) |
| N4 | 鏡 | かがみ | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%8F%A1) |
| N4 | 集まる | あつまる | jlpt-n4, jlpt-n1 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%9B%86%E3%81%BE%E3%82%8B) |
| N4 | 集める | あつめる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%9B%86%E3%82%81%E3%82%8B) |
| N4 | 雲 | くも | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%9B%B2) |
| N4 | 飾る | かざる | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%A3%BE%E3%82%8B) |
| N4 | 首 | くび | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%A6%96) |
| N4 | 高校 | こうこう | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%AB%98%E6%A0%A1) |
| N4 | 高校生 | こうこうせい | jlpt-n4 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%AB%98%E6%A0%A1%E7%94%9F) |
| N5 | ええ | ええ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%88%E3%81%88) |
| N5 | お兄さん | おにいさん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%85%84%E3%81%95%E3%82%93) |
| N5 | お姉さん | おねえさん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%A7%89%E3%81%95%E3%82%93) |
| N5 | お弁当 | おべんとう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%BC%81%E5%BD%93) |
| N5 | お母さん | おかあさん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E6%AF%8D%E3%81%95%E3%82%93) |
| N5 | お父さん | おとうさん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E7%88%B6%E3%81%95%E3%82%93) |
| N5 | お茶 | おちゃ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E8%8C%B6) |
| N5 | お菓子 | おかし | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E8%8F%93%E5%AD%90) |
| N5 | お金 | おかね | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E9%87%91) |
| N5 | お風呂 | おふろ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E9%A2%A8%E5%91%82) |
| N5 | こんな | こんな | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%93%E3%82%93%E3%81%AA) |
| N5 | さあ | さあ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%95%E3%81%82) |
| N5 | どうぞ | どうぞ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%81%86%E3%81%9E) |
| N5 | どうも | どうも | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%81%86%E3%82%82) |
| N5 | もう | もう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%82%E3%81%86) |
| N5 | もっと | もっと | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%82%E3%81%A3%E3%81%A8) |
| N5 | アパート | アパート | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%91%E3%83%BC%E3%83%88) |
| N5 | エレベーター | エレベーター | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%83%AC%E3%83%99%E3%83%BC%E3%82%BF%E3%83%BC) |
| N5 | カップ | カップ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%83%E3%83%97) |
| N5 | カメラ | カメラ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%A1%E3%83%A9) |
| N5 | カレンダー | カレンダー | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AB%E3%83%AC%E3%83%B3%E3%83%80%E3%83%BC) |
| N5 | ギター | ギター | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AE%E3%82%BF%E3%83%BC) |
| N5 | クラス | クラス | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%AF%E3%83%A9%E3%82%B9) |
| N5 | コート | コート | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B3%E3%83%BC%E3%83%88) |
| N5 | シャツ | シャツ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%A3%E3%83%84) |
| N5 | シャワー | シャワー | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%A3%E3%83%AF%E3%83%BC) |
| N5 | スカート | スカート | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%AB%E3%83%BC%E3%83%88) |
| N5 | ストーブ | ストーブ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%88%E3%83%BC%E3%83%96) |
| N5 | スプーン | スプーン | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%97%E3%83%BC%E3%83%B3) |
| N5 | スポーツ | スポーツ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%9D%E3%83%BC%E3%83%84) |
| N5 | セーター | セーター | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%BB%E3%83%BC%E3%82%BF%E3%83%BC) |
| N5 | タクシー | タクシー | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%BF%E3%82%AF%E3%82%B7%E3%83%BC) |
| N5 | テスト | テスト | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%86%E3%82%B9%E3%83%88) |
| N5 | テレビ | テレビ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%86%E3%83%AC%E3%83%93) |
| N5 | テーブル | テーブル | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB) |
| N5 | テープ | テープ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%86%E3%83%BC%E3%83%97) |
| N5 | デパート | デパート | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%87%E3%83%91%E3%83%BC%E3%83%88) |
| N5 | トイレ | トイレ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%88%E3%82%A4%E3%83%AC) |
| N5 | ドア | ドア | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%89%E3%82%A2) |
| N5 | ナイフ | ナイフ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8A%E3%82%A4%E3%83%95) |
| N5 | ニュース | ニュース | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8B%E3%83%A5%E3%83%BC%E3%82%B9) |
| N5 | ネクタイ | ネクタイ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8D%E3%82%AF%E3%82%BF%E3%82%A4) |
| N5 | ノート | ノート | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8E%E3%83%BC%E3%83%88) |
| N5 | ハンカチ | ハンカチ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8F%E3%83%B3%E3%82%AB%E3%83%81) |
| N5 | バス | バス | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%90%E3%82%B9) |
| N5 | バター | バター | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%90%E3%82%BF%E3%83%BC) |
| N5 | パン | パン | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%91%E3%83%B3) |
| N5 | パーティー | パーティー | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%91%E3%83%BC%E3%83%86%E3%82%A3%E3%83%BC) |
| N5 | フィルム | フィルム | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%95%E3%82%A3%E3%83%AB%E3%83%A0) |
| N5 | フォーク | フォーク | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%95%E3%82%A9%E3%83%BC%E3%82%AF) |
| N5 | プール | プール | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%97%E3%83%BC%E3%83%AB) |
| N5 | ベッド | ベッド | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%99%E3%83%83%E3%83%89) |
| N5 | ペット | ペット | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9A%E3%83%83%E3%83%88) |
| N5 | ホテル | ホテル | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9B%E3%83%86%E3%83%AB) |
| N5 | ボールペン | ボールペン | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9C%E3%83%BC%E3%83%AB%E3%83%9A%E3%83%B3) |
| N5 | ポケット | ポケット | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9D%E3%82%B1%E3%83%83%E3%83%88) |
| N5 | ポスト | ポスト | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9D%E3%82%B9%E3%83%88) |
| N5 | ラジオ | ラジオ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%A9%E3%82%B8%E3%82%AA) |
| N5 | レコード | レコード | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%AC%E3%82%B3%E3%83%BC%E3%83%89) |
| N5 | レストラン | レストラン | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%AC%E3%82%B9%E3%83%88%E3%83%A9%E3%83%B3) |
| N5 | 一 | いち | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80) |
| N5 | 一つ | ひとつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E3%81%A4) |
| N5 | 一月 | ひとつき | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%9C%88) |
| N5 | 一緒 | いっしょ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E7%B7%92) |
| N5 | 七つ | ななつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%83%E3%81%A4) |
| N5 | 万 | まん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%87) |
| N5 | 万年筆 | まんねんひつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%87%E5%B9%B4%E7%AD%86) |
| N5 | 丈夫 | じょうぶ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%88%E5%A4%AB) |
| N5 | 三つ | みっつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%89%E3%81%A4) |
| N5 | 上 | うえ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A) |
| N5 | 上手 | じょうず | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E6%89%8B) |
| N5 | 上着 | うわぎ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E7%9D%80) |
| N5 | 下 | した | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B) |
| N5 | 下さい | ください | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B%E3%81%95%E3%81%84) |
| N5 | 下手 | へた | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8B%E6%89%8B) |
| N5 | 両親 | りょうしん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%A1%E8%A6%AA) |
| N5 | 並ぶ | ならぶ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%A6%E3%81%B6) |
| N5 | 並べる | ならべる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%A6%E3%81%B9%E3%82%8B) |
| N5 | 乗る | のる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%97%E3%82%8B) |
| N5 | 九つ | ここのつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%9D%E3%81%A4) |
| N5 | 二 | に | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8C) |
| N5 | 二つ | ふたつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8C%E3%81%A4) |
| N5 | 五 | ご | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%94) |
| N5 | 五つ | いつつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%94%E3%81%A4) |
| N5 | 交差点 | こうさてん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A4%E5%B7%AE%E7%82%B9) |
| N5 | 交番 | こうばん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%A4%E7%95%AA) |
| N5 | 人 | ひと | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%BA) |
| N5 | 今 | いま | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A) |
| N5 | 今年 | ことし | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E5%B9%B4) |
| N5 | 今晩 | こんばん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E6%99%A9) |
| N5 | 今月 | こんげつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E6%9C%88) |
| N5 | 今朝 | けさ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E6%9C%9D) |
| N5 | 今週 | こんしゅう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E9%80%B1) |
| N5 | 仕事 | しごと | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%95%E4%BA%8B) |
| N5 | 休み | やすみ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%91%E3%81%BF) |
| N5 | 休む | やすむ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%91%E3%82%80) |
| N5 | 会社 | かいしゃ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E7%A4%BE) |
| N5 | 低い | ひくい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%8E%E3%81%84) |
| N5 | 住む | すむ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%8F%E3%82%80) |
| N5 | 作文 | さくぶん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%9C%E6%96%87) |
| N5 | 使う | つかう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%BF%E3%81%86) |
| N5 | 便利 | べんり | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BE%BF%E5%88%A9) |
| N5 | 借りる | かりる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%80%9F%E3%82%8A%E3%82%8B) |
| N5 | 傘 | かさ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%82%98) |
| N5 | 働く | はたらく | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%83%8D%E3%81%8F) |
| N5 | 元気 | げんき | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%83%E6%B0%97) |
| N5 | 兄 | あに | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%84) |
| N5 | 兄弟 | きょうだい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%84%E5%BC%9F) |
| N5 | 先 | さき | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%88) |
| N5 | 先月 | せんげつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%88%E6%9C%88) |
| N5 | 先生 | せんせい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%88%E7%94%9F) |
| N5 | 先週 | せんしゅう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%88%E9%80%B1) |
| N5 | 入る | はいる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%A5%E3%82%8B) |
| N5 | 入れる | いれる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%A5%E3%82%8C%E3%82%8B) |
| N5 | 全部 | ぜんぶ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%A8%E9%83%A8) |
| N5 | 八 | はち | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AB) |
| N5 | 八つ | やっつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AB%E3%81%A4) |
| N5 | 八百屋 | やおや | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AB%E7%99%BE%E5%B1%8B) |
| N5 | 公園 | こうえん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AC%E5%9C%92) |
| N5 | 六つ | むっつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AD%E3%81%A4) |
| N5 | 写真 | しゃしん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%99%E7%9C%9F) |
| N5 | 冬 | ふゆ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%AC) |
| N5 | 冷たい | つめたい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%B7%E3%81%9F%E3%81%84) |
| N5 | 冷蔵庫 | れいぞうこ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%B7%E8%94%B5%E5%BA%AB) |
| N5 | 出す | だす | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%BA%E3%81%99) |
| N5 | 出る | でる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%BA%E3%82%8B) |
| N5 | 出口 | でぐち | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%BA%E5%8F%A3) |
| N5 | 切る | きる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%87%E3%82%8B) |
| N5 | 切手 | きって | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%87%E6%89%8B) |
| N5 | 切符 | きっぷ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%87%E7%AC%A6) |
| N5 | 初めて | はじめて | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%9D%E3%82%81%E3%81%A6) |
| N5 | 動物 | どうぶつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8B%95%E7%89%A9) |
| N5 | 北 | きた | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8C%97) |
| N5 | 医者 | いしゃ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8C%BB%E8%80%85) |
| N5 | 千 | せん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%83) |
| N5 | 午前 | ごぜん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%88%E5%89%8D) |
| N5 | 午後 | ごご | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%88%E5%BE%8C) |
| N5 | 半 | はん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%8A) |
| N5 | 半分 | はんぶん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%8A%E5%88%86) |
| N5 | 南 | みなみ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%97) |
| N5 | 危ない | あぶない | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%B1%E3%81%AA%E3%81%84) |
| N5 | 卵 | たまご | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%B5) |
| N5 | 厚い | あつい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8E%9A%E3%81%84) |
| N5 | 去年 | きょねん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8E%BB%E5%B9%B4) |
| N5 | 友達 | ともだち | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%8B%E9%81%94) |
| N5 | 取る | とる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%96%E3%82%8B) |
| N5 | 口 | くち | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%A3) |
| N5 | 古い | ふるい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%A4%E3%81%84) |
| N5 | 台所 | だいどころ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%B0%E6%89%80) |
| N5 | 右 | みぎ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%B3) |
| N5 | 同じ | おなじ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%8C%E3%81%98) |
| N5 | 名前 | なまえ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%8D%E5%89%8D) |
| N5 | 向こう | むこう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%91%E3%81%93%E3%81%86) |
| N5 | 吸う | すう | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%B8%E3%81%86) |
| N5 | 吹く | ふく | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%90%B9%E3%81%8F) |
| N5 | 呼ぶ | よぶ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%91%BC%E3%81%B6) |
| N5 | 咲く | さく | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%92%B2%E3%81%8F) |
| N5 | 問題 | もんだい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%95%8F%E9%A1%8C) |
| N5 | 喫茶店 | きっさてん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%96%AB%E8%8C%B6%E5%BA%97) |
| N5 | 四つ | よっつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%9B%E3%81%A4) |
| N5 | 困る | こまる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%B0%E3%82%8B) |
| N5 | 図書館 | としょかん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%B3%E6%9B%B8%E9%A4%A8) |
| N5 | 国 | くに | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%BD) |
| N5 | 土曜日 | どようび | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9C%9F%E6%9B%9C%E6%97%A5) |
| N5 | 地下鉄 | ちかてつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9C%B0%E4%B8%8B%E9%89%84) |
| N5 | 地図 | ちず | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9C%B0%E5%9B%B3) |
| N5 | 声 | こえ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A3%B0) |
| N5 | 売る | うる | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A3%B2%E3%82%8B) |
| N5 | 夏 | なつ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%8F) |
| N5 | 夏休み | なつやすみ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%8F%E4%BC%91%E3%81%BF) |
| N5 | 夕方 | ゆうがた | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%95%E6%96%B9) |
| N5 | 夕飯 | ゆうはん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%95%E9%A3%AF) |
| N5 | 外 | そと | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%96) |
| N5 | 外国 | がいこく | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%96%E5%9B%BD) |
| N5 | 外国人 | がいこくじん | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%96%E5%9B%BD%E4%BA%BA) |
| N5 | 多い | おおい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%9A%E3%81%84) |
| N5 | 大きい | おおきい | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%A7%E3%81%8D%E3%81%84) |
| N5 | 大きな | おおきな | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%A7%E3%81%8D%E3%81%AA) |
| N5 | 大丈夫 | だいじょうぶ | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%A7%E4%B8%88%E5%A4%AB) |
| N5 | 大人 | おとな | jlpt-n5 | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%A7%E4%BA%BA) |

## 等級差異

| 目前等級 | Headword | Reading | Jisho | 建議 | Citation |
|---|---|---|---|---|---|
| N1 | いつの間にか | いつのまにか | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%84%E3%81%A4%E3%81%AE%E9%96%93%E3%81%AB%E3%81%8B) |
| N1 | いよいよ | いよいよ | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%84%E3%82%88%E3%81%84%E3%82%88) |
| N1 | 一々 | いちいち | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E3%80%85) |
| N1 | 呆れる | あきれる | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%91%86%E3%82%8C%E3%82%8B) |
| N1 | 当てはまる | あてはまる | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%BD%93%E3%81%A6%E3%81%AF%E3%81%BE%E3%82%8B) |
| N1 | 慌ただしい | あわただしい | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%85%8C%E3%81%9F%E3%81%A0%E3%81%97%E3%81%84) |
| N1 | 曖昧 | あいまい | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9B%96%E6%98%A7) |
| N1 | 相変わらず | あいかわらず | N2 | N1 -> N2（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%9B%B8%E5%A4%89%E3%82%8F%E3%82%89%E3%81%9A) |
| N1 | 兎に角 | とにかく | N3 | N1 -> N3（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%8E%E3%81%AB%E8%A7%92) |
| N1 | 所謂 | いわゆる | N3 | N1 -> N3（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%89%80%E8%AC%82) |
| N1 | 未だ | いまだ | N3 | N1 -> N3（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%9C%AA%E3%81%A0) |
| N1 | 苛々 | いらいら | N3 | N1 -> N3（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E8%8B%9B%E3%80%85) |
| N2 | スリッパ | スリッパ | N5 | N2 -> N5（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%AA%E3%83%83%E3%83%91) |
| N3 | ずれる | ずれる | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9A%E3%82%8C%E3%82%8B) |
| N3 | ほっと | ほっと | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%BB%E3%81%A3%E3%81%A8) |
| N3 | アンケート | アンケート | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A2%E3%83%B3%E3%82%B1%E3%83%BC%E3%83%88) |
| N3 | ショック | ショック | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B7%E3%83%A7%E3%83%83%E3%82%AF) |
| N3 | ストレス | ストレス | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%83%88%E3%83%AC%E3%82%B9) |
| N3 | タイトル | タイトル | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%BF%E3%82%A4%E3%83%88%E3%83%AB) |
| N3 | デザイン | デザイン | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%87%E3%82%B6%E3%82%A4%E3%83%B3) |
| N3 | データ | データ | N1 | N3 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%87%E3%83%BC%E3%82%BF) |
| N3 | うっかり | うっかり | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%86%E3%81%A3%E3%81%8B%E3%82%8A) |
| N3 | お洒落 | おしゃれ | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E6%B4%92%E8%90%BD) |
| N3 | お辞儀 | おじぎ | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E8%BE%9E%E5%84%80) |
| N3 | そっと | そっと | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9D%E3%81%A3%E3%81%A8) |
| N3 | ぶつかる | ぶつかる | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%B6%E3%81%A4%E3%81%8B%E3%82%8B) |
| N3 | インタビュー | インタビュー | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%93%E3%83%A5%E3%83%BC) |
| N3 | スケジュール | スケジュール | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%B9%E3%82%B1%E3%82%B8%E3%83%A5%E3%83%BC%E3%83%AB) |
| N3 | トレーニング | トレーニング | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%88%E3%83%AC%E3%83%BC%E3%83%8B%E3%83%B3%E3%82%B0) |
| N3 | バランス | バランス | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%90%E3%83%A9%E3%83%B3%E3%82%B9) |
| N3 | マイナス | マイナス | N2 | N3 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9E%E3%82%A4%E3%83%8A%E3%82%B9) |
| N4 | 写る | うつる | N2 | N4 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%86%99%E3%82%8B) |
| N4 | けが | けが | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%91%E3%81%8C) |
| N4 | ごみ | ごみ | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%94%E3%81%BF) |
| N4 | 中々 | なかなか | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%AD%E3%80%85) |
| N4 | 事 | こと | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8B) |
| N4 | 全然 | ぜんぜん | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%A8%E7%84%B6) |
| N4 | 嬉しい | うれしい | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%AC%89%E3%81%97%E3%81%84) |
| N4 | 案内 | あんない | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E6%A1%88%E5%86%85) |
| N4 | 経験 | けいけん | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%B5%8C%E9%A8%93) |
| N4 | 運動 | うんどう | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8B%E5%8B%95) |
| N4 | 運転 | うんてん | N3 | N4 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E9%81%8B%E8%BB%A2) |
| N4 | 米 | こめ | N5 | N4 -> N5（Jisho 較低階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E7%B1%B3) |
| N5 | ああ | ああ | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%82%E3%81%82) |
| N5 | あの | あの | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%82%E3%81%AE) |
| N5 | お手洗い | おてあらい | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E6%89%8B%E6%B4%97%E3%81%84) |
| N5 | 一日 | いちにち | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%97%A5) |
| N5 | 丁度 | ちょうど | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%81%E5%BA%A6) |
| N5 | 三 | さん | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%89) |
| N5 | 二人 | ふたり | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8C%E4%BA%BA) |
| N5 | 六 | ろく | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%85%AD) |
| N5 | 前 | まえ | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%89%8D) |
| N5 | 可愛い | かわいい | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8F%AF%E6%84%9B%E3%81%84) |
| N5 | 塩 | しお | N1 | N5 -> N1（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A1%A9) |
| N5 | ペン | ペン | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%9A%E3%83%B3) |
| N5 | 一昨日 | おととい | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E6%98%A8%E6%97%A5) |
| N5 | 会う | あう | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BC%9A%E3%81%86) |
| N5 | 作る | つくる | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%9C%E3%82%8B) |
| N5 | 出かける | でかける | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%87%BA%E3%81%8B%E3%81%91%E3%82%8B) |
| N5 | 分かる | わかる | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%88%86%E3%81%8B%E3%82%8B) |
| N5 | 勤める | つとめる | N2 | N5 -> N2（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8B%A4%E3%82%81%E3%82%8B) |
| N5 | いいえ | いいえ | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%84%E3%81%84%E3%81%88) |
| N5 | お腹 | おなか | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E8%85%B9) |
| N5 | では | では | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A7%E3%81%AF) |
| N5 | でも | でも | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A7%E3%82%82) |
| N5 | どんな | どんな | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%A9%E3%82%93%E3%81%AA) |
| N5 | はい | はい | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AF%E3%81%84) |
| N5 | 一人 | ひとり | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E4%BA%BA) |
| N5 | 一番 | いちばん | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E7%95%AA) |
| N5 | 七 | しち | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%83) |
| N5 | 中 | なか | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%AD) |
| N5 | 九 | きゅう | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B9%9D) |
| N5 | 二十歳 | はたち | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BA%8C%E5%8D%81%E6%AD%B3) |
| N5 | 今日 | きょう | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BB%8A%E6%97%A5) |
| N5 | 体 | からだ | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%93) |
| N5 | 余り | あまり | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%BD%99%E3%82%8A) |
| N5 | 十 | じゅう | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%81) |
| N5 | 四 | し | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%9B%9B) |
| N5 | 多分 | たぶん | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%9A%E5%88%86) |
| N5 | 夜 | よる | N3 | N5 -> N3（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%9C) |
| N5 | 上げる | あげる | N4 | N5 -> N4（Jisho 較高階） | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%8A%E3%81%92%E3%82%8B) |

## Jisho 無收錄

| 目前等級 | Headword | Reading | Jisho tag | Jisho URL |
|---|---|---|---|---|
| N2 | さっさと | さっさと | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%95%E3%81%A3%E3%81%95%E3%81%A8) |
| N2 | せめて | せめて | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%9B%E3%82%81%E3%81%A6) |
| N2 | インキ | インキ | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A4%E3%83%B3%E3%82%AD) |
| N3 | イメージ | イメージ | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A4%E3%83%A1%E3%83%BC%E3%82%B8) |
| N4 | お子さん | おこさん | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%8A%E5%AD%90%E3%81%95%E3%82%93) |
| N4 | エアコン | エアコン | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%82%A2%E3%82%B3%E3%83%B3) |
| N4 | ハンバーグ | ハンバーグ | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%83%8F%E3%83%B3%E3%83%90%E3%83%BC%E3%82%B0) |
| N5 | ない | ない | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%AA%E3%81%84) |
| N5 | もしもし | もしもし | unlisted | [Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%82%E3%81%97%E3%82%82%E3%81%97) |

## 已人工覆蓋

以下 rows 命中 `server/data/corpus/jlpt-overrides.jsonl`，依要求只列入本節，不重新分類。

| 目前檔案等級 | Override 等級 | Headword | Reading | Reason |
|---|---|---|---|---|
| N4 | N4 | お代わり | おかわり | common everyday dining word; external JLPT overlay placed it too high |
| N4 | N4 | お先に | おさきに | common departure formula; external JLPT overlay placed it too high |
| N4 | N4 | お参り | おまいり | common everyday/culture word; external JLPT overlay placed it too high |
| N4 | N4 | お大事に | おだいじに | common health greeting; external JLPT overlay placed it too high |
| N4 | N4 | お待たせしました | おまたせしました | common service formula; external JLPT overlay placed it too high |
| N4 | N4 | お手伝いさん | おてつだいさん | common household word; external JLPT overlay placed it too high |
| N5 | N5 | うん | うん | basic casual yes response; external JLPT overlay placed it too high |
| N5 | N5 | お休み | おやすみ | basic greeting; external JLPT overlay placed it too high |
| N5 | N5 | お帰り | おかえり | basic home greeting; external JLPT overlay placed it too high |
| N5 | N5 | お邪魔します | おじゃまします | basic visit greeting; external JLPT overlay placed it too high |
| N5 | N5 | お願いします | おねがいします | basic request formula; external JLPT overlay placed it too high |
| N5 | N5 | そう | そう | basic conversational demonstrative/response; external JLPT overlay placed it too high |
| N5 | N5 | ケーキ | ケーキ | basic everyday food loanword; external JLPT overlay placed it too high |
| N5 | N5 | サラダ | サラダ | basic everyday food loanword; external JLPT overlay placed it too high |
| N5 | N5 | サンドイッチ | サンドイッチ | basic everyday food loanword; external JLPT overlay placed it too high |
| N5 | N5 | テニス | テニス | basic sport loanword; external JLPT overlay placed it too high |
| N5 | N5 | ピアノ | ピアノ | basic instrument loanword; external JLPT overlay placed it too high |
| N5 | N5 | プレゼント | プレゼント | basic everyday gift word; external JLPT overlay placed it too high |
| N5 | N5 | ベル | ベル | basic everyday object word; external JLPT overlay placed it too high |

## 遷移建議

- 非人工覆蓋 rows：671；等級差異：81；差異率：12.07%。
- Jisho 無收錄 / 無 JLPT tag：9；占非人工覆蓋 1.34%。
- 主要差異方向：N5->N3: 19、N5->N1: 11、N3->N2: 10、N4->N3: 10、N1->N2: 8、N3->N1: 8、N5->N2: 7、N1->N3: 4、N2->N5: 1、N4->N2: 1。
- Recommendation：建議建立小批次 migration：先人工審核所有「等級差異」列，確認 Jisho 第一筆結果 sense 無誤後再調整檔案歸屬；無收錄項目維持現狀或補 override reason。

---

Footer self-verify：

- Cross-source：非人工覆蓋 671 rows 均有 Jisho API lookup record（unique headwords: 671）。
- 5-item OK re-check：
  - OK: エスカレーター（エスカレーター）目前 N4；Jisho fresh=N4；[Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%82%A8%E3%82%B9%E3%82%AB%E3%83%AC%E3%83%BC%E3%82%BF%E3%83%BC)
  - OK: 一律（いちりつ）目前 N1；Jisho fresh=N1；[Jisho](https://jisho.org/api/v1/search/words?keyword=%E4%B8%80%E5%BE%8B)
  - OK: ええ（ええ）目前 N5；Jisho fresh=N5；[Jisho](https://jisho.org/api/v1/search/words?keyword=%E3%81%88%E3%81%88)
  - OK: 大丈夫（だいじょうぶ）目前 N5；Jisho fresh=N5；[Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%A4%A7%E4%B8%88%E5%A4%AB)
  - OK: 危ない（あぶない）目前 N5；Jisho fresh=N5；[Jisho](https://jisho.org/api/v1/search/words?keyword=%E5%8D%B1%E3%81%AA%E3%81%84)

Generated: 2026-05-03 Asia/Tokyo. Script source: temporary local audit runner, not committed. Network source: Jisho API first-result `jlpt` tags.
