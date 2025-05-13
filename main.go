package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// resolveAt  は JSON パース後のデータ構造を走査し、値が文字列かつ
// 先頭に "@" が付いている場合、そのパスが指すファイルを読み込んで
// 文字列の中身をファイル内容に置き換える。
func resolveAt(v any) (any, error) {
	switch vv := v.(type) {
	case string:
		if strings.HasPrefix(vv, "@") {
			// @ の後をパスとして展開 (環境変数も解決)
			path := os.ExpandEnv(strings.TrimPrefix(vv, "@"))
			// パスが相対だった場合は cwd からの相対
			path, _ = filepath.Abs(path)
			bytes, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file '%s': %w", path, err)
			}
			return string(bytes), nil
		}
		return vv, nil
	case map[string]any:
		for k, val := range vv {
			resolved, err := resolveAt(val)
			if err != nil {
				return nil, err
			}
			vv[k] = resolved
		}
		return vv, nil
	case []any:
		for i, val := range vv {
			resolved, err := resolveAt(val)
			if err != nil {
				return nil, err
			}
			vv[i] = resolved
		}
		return vv, nil
	default:
		return vv, nil
	}
}

func main() {
	inputArg := flag.String("input", "", "入力データ（JSON文字列または @file_path）")
	templatePath := flag.String("template", "", "テンプレートファイルのパス")
	flag.Parse()

	if *inputArg == "" || *templatePath == "" {
		fmt.Fprintln(os.Stderr, "Usage: render -input <json|@file> -template <template-file>")
		os.Exit(1)
	}

	// 入力文字列の取得
	var inputData string
	if strings.HasPrefix(*inputArg, "@") {
		// @FILE の場合: JSON 全体をファイルから読む
		path := strings.TrimPrefix(*inputArg, "@")
		bytes, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read input file: %v\n", err)
			os.Exit(1)
		}
		inputData = string(bytes)
	} else {
		inputData = *inputArg
	}

	// JSON をパース
	var data map[string]any
	if err := json.Unmarshal([]byte(inputData), &data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse JSON: %v\n", err)
		os.Exit(1)
	}

	// 値の中に @ があった場合はファイル展開
	resolved, err := resolveAt(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	data = resolved.(map[string]any)

	// テンプレートを読み込んで実行
	tmpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse template: %v\n", err)
		os.Exit(1)
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute template: %v\n", err)
		os.Exit(1)
	}
}
