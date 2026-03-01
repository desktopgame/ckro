# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

ckroはGoで開発されたTUIマークダウンエディターで、ObsidianとCursorに影響を受けています。コーディング用ではなく、日記や記事の執筆にフォーカスしています。最大の特徴は、リッチテキストの表示と編集を同時に行えること（プレビューモードが不要）です。独自のMarkdown方言は **litemark** と呼ばれています。

## 言語

日本語で応答すること。ソースコード中のコメントは常に英語で記述すること。

## ビルド・実行

```bash
go build ./cmd/ckro        # ビルド
go test ./internal/tui     # テスト実行（シナリオベース、testdata/参照）
ckro                       # エディター起動
ckro -d ./project          # 指定ディレクトリで起動
```

## アーキテクチャ

### モジュール境界

- `cmd/ckro/` — アプリケーション固有の処理。`internal/` に依存する。
- `internal/` — 使いまわしのきく処理。**`cmd/ckro/` に依存してはいけない。**

### MVCテキスト編集コア

```
Model (internal/tui/model/)         — Document, Element, Buffer
View  (internal/tui/view/)          — TextViewインターフェイス、描画ロジック
Presenter (internal/tui/presenter/) — モデルを視覚表現に変換
Control (internal/tui/)             — TextBox、レイアウトウィジェット（Box, Grid, Stack, Frame, Center）
```

### 二重カーソル座標系

設計の核心：カーソルは常に可視かつ編集可能な部分を指し、不可視の装飾（太字の`*`など）を指してはいけない。

- **モデル座標**: (行, 先頭バイト位置, バイト長) — 常に編集可能な内容を指す
- **ビュー座標**: 単一の符号なし整数（0〜全TextViewの`MoveLength()`の和）— 可視カーソル位置にマッピング

`TextView`インターフェイスの`MoveLeft/Right/Up/Down`が不可視装飾をスキップするカーソル移動を処理する。`ConvertModel`でビュー→モデル座標変換、`ConvertViewLocalPos`でモデル→ビュー座標変換を行う。

### コントロール階層

全UIは`Control`インターフェイスを実装: `Traverse()`, `Update()`, `Draw()`, `MinimumSize()`, `Layout()`, `IsFlexibleWidth()`, `IsFlexibleHeight()`。

- **Tile** — TextBox + Presenterをラップする基本コントロール（ツリー、リスト、行番号、スクロールバーに使用）
- **Box** — 一方向レイアウト（水平/垂直）、フレキシブル属性でサイズ分配
- **Window** — Z方向のスタック管理とモーダルダイアログのフォーカス制御

### LLM連携

OpenAI互換API（`openai-go`）+ Model Context Protocol（`go-sdk`）を使用。LLMの応答は特殊な構文ブロック `{{{ BOT: }}}`, `{{{ CALL: }}}`, `{{{ TOOL: }}}` でエディターに挿入される。

起動ディレクトリの `ckro.yaml` で設定（ApiKey, BaseUrl, Model, SystemPrompt）。

### Litemarkパーサー

`internal/tui/extensions/litemark/` — Elementツリーを生成する独自Markdownパーサー。各要素型に対応するTextViewが装飾付きの描画を担当（見出し、太字、コードブロック、テーブルなど）。

## テスト

`testdata/*.txt` にカスタムシナリオ言語で定義:

```
%%
幅 高さ エンジン(PlainText|StyledText)
%%
命令列 (TYPE, MOVE_LEFT, BYTE_POS など)
%%
期待するテキスト
%%
```

主な命令: `TYPE "text"`, `TYPELN "text"`, `REMOVE_CHAR`, `MOVE_LEFT/RIGHT/UP/DOWN`, `SELECTION_START/END`, `BYTE_POS 行 列`, `VAR "変数名" "中身"`。

## 主要依存関係

| パッケージ | 役割 |
|-----------|------|
| `gdamore/tcell/v2` | ターミナル描画バックエンド（tviewは不使用） |
| `openai/openai-go/v2` | LLM APIクライアント |
| `modelcontextprotocol/go-sdk` | LLMツール連携のためのMCPプロトコル |
| `rivo/uniseg` | Unicodeグラフェムクラスター処理 |

## 制約事項

- 横方向のスクロールは非サポート
- リッチテキストがウィンドウ幅を超える場合、プレーンテキストにフォールバックして折り返す
- ルートのTextViewは行単位で空間を占める
