package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Provider capabilities
type Capabilities struct {
	Kline        bool
	Spot         bool
	Instrument   bool
	Profile      bool
	Financial    bool
	Announcement bool
}

var providerNames = map[string]string{
	"eastmoney":   "EastMoney (东方财富)",
	"sina":        "Sina (新浪财经)",
	"tencent":     "Tencent (腾讯财经)",
	"xueqiu":      "Xueqiu (雪球)",
	"tushare":     "Tushare (Tushare Pro)",
	"cninfo":      "Cninfo (巨潮资讯)",
	"sse":         "SSE (上交所)",
	"szse":        "SZSE (深交所)",
	"bse":         "BSE (北交所)",
	"binance":     "Binance (币安)",
	"eastmoneyhk": "EastMoney HK (东财港股)",
	"yahoo":       "Yahoo Finance",
}

func main() {
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	adaptersDir := filepath.Join(rootDir, "adapters")
	entries, err := os.ReadDir(adaptersDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading adapters dir: %v\n", err)
		os.Exit(1)
	}

	providers := make(map[string]*Capabilities)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == "mocks" {
			continue
		}

		caps := &Capabilities{}
		providers[name] = caps

		files, _ := os.ReadDir(filepath.Join(adaptersDir, name))
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			fname := f.Name()
			if strings.HasSuffix(fname, "_test.go") {
				continue
			}

			switch fname {
			case "kline.go":
				caps.Kline = true
			case "spot.go":
				caps.Spot = true
			case "instrument.go":
				caps.Instrument = true
			case "profile.go":
				caps.Profile = true
			case "financial.go":
				caps.Financial = true
			case "announcement.go", "news.go":
				caps.Announcement = true
			}
		}
	}

	table := generateMarkdownTable(providers)

	updateFile(filepath.Join(rootDir, "README.md"), table)
	updateFile(filepath.Join(rootDir, "README_EN.md"), table)

	fmt.Println("Documentation updated successfully.")
}

func generateMarkdownTable(providers map[string]*Capabilities) string {
	var buf bytes.Buffer

	// Header
	buf.WriteString("| 数据源 (Provider) | K线 (Kline) | 实时行情 (Spot) | 证券列表 (Instrument) | 证券详情 (Profile) | 财务数据 (Financial) | 公告资讯 (News) |\n")
	buf.WriteString("| :--- | :---: | :---: | :---: | :---: | :---: | :---: |\n")

	// Sort providers by name for consistency
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		caps := providers[name]
		displayName := providerNames[name]
		if displayName == "" {
			displayName = "**" + strings.Title(name) + "**"
		} else {
			// Format as **Key** (Desc)
			parts := strings.SplitN(displayName, " (", 2)
			if len(parts) == 2 {
				displayName = "**" + parts[0] + "** (" + parts[1]
			} else {
				displayName = "**" + displayName + "**"
			}
		}

		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
			displayName,
			checkMark(caps.Kline),
			checkMark(caps.Spot),
			checkMark(caps.Instrument),
			checkMark(caps.Profile),
			checkMark(caps.Financial),
			checkMark(caps.Announcement),
		))
	}

	return buf.String()
}

func checkMark(supported bool) string {
	if supported {
		return "✅"
	}
	return "-"
}

func updateFile(path string, content string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %s: %v\n", path, err)
		return
	}

	startMarker := "<!-- START_SUPPORTED_TABLE -->"
	endMarker := "<!-- END_SUPPORTED_TABLE -->"

	sData := string(data)
	startIdx := strings.Index(sData, startMarker)
	endIdx := strings.Index(sData, endMarker)

	if startIdx == -1 || endIdx == -1 {
		fmt.Printf("Markers not found in %s, skipping update.\n", path)
		return
	}

	newContent := sData[:startIdx+len(startMarker)] + "\n\n" + content + "\n" + sData[endIdx:]

	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Could not write file %s: %v\n", path, err)
	} else {
		fmt.Printf("Updated %s\n", path)
	}
}
