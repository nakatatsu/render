# render

簡単なレンダリングツール。

## USAGE

```bash
render -input <json> -template <template-file>
```

jsonは必ず{"key": "value"}の形式とする。valueの先頭に@がある場合はfile_pathと見なし、読みとった上でレンダリングする。

## コマンド例

```bash
render -input '{"style": "@./sample/guidelines.md", "purpose": "IPv6オンリーのVPCを作りたい"}' -template ./sample/prompt.tpl
```
