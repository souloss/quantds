package main

import (
	"bytes"
	"fmt"
	"net/url"
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
	Markets      []string // Supported markets: CN, HK, US, Crypto
}

var providerNames = map[string]string{
	"eastmoney":   "EastMoney",
	"sina":        "Sina",
	"tencent":     "Tencent",
	"xueqiu":      "Xueqiu",
	"tushare":     "Tushare",
	"cninfo":      "Cninfo",
	"sse":         "SSE",
	"szse":        "SZSE",
	"bse":         "BSE",
	"binance":     "Binance",
	"eastmoneyhk": "EastMoneyHK",
	"yahoo":       "Yahoo",
}

// Manual mapping of provider supported markets (simplest way without parsing complex AST)
var providerMarkets = map[string][]string{
	"eastmoney":   {"A股"},
	"sina":        {"A股", "港股"},
	"tencent":     {"A股", "港股", "美股"},
	"xueqiu":      {"A股", "港股", "美股"},
	"tushare":     {"A股"},
	"cninfo":      {"A股"},
	"sse":         {"A股"},
	"szse":        {"A股"},
	"bse":         {"A股"},
	"binance":     {"Crypto"},
	"eastmoneyhk": {"港股"},
	"yahoo":       {"A股", "港股", "美股"},
}

// Beta status mapping
var betaProviders = map[string]bool{
	"xueqiu": true,
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
		caps.Markets = providerMarkets[name]

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
	badges := generateBadges(providers)

	updateFile(filepath.Join(rootDir, "README.md"), "<!-- START_SUPPORTED_TABLE -->", "<!-- END_SUPPORTED_TABLE -->", table)
	updateFile(filepath.Join(rootDir, "README_EN.md"), "<!-- START_SUPPORTED_TABLE -->", "<!-- END_SUPPORTED_TABLE -->", table)

	updateFile(filepath.Join(rootDir, "README.md"), "<!-- START_STATUS_BADGES -->", "<!-- END_STATUS_BADGES -->", badges)
	updateFile(filepath.Join(rootDir, "README_EN.md"), "<!-- START_STATUS_BADGES -->", "<!-- END_STATUS_BADGES -->", badges)

	// Generate SVG
	if err := generateSVG(providers, filepath.Join(rootDir, "docs", "supported_sources.svg")); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating SVG: %v\n", err)
	}

	fmt.Println("Documentation updated successfully.")
}

func generateSVG(providers map[string]*Capabilities, path string) error {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)

	// Constants
	rowHeight := 30
	headerHeight := 40
	colWidth := 140
	firstColWidth := 180
	fontSize := 14
	headerFontSize := 14
	padding := 10

	cols := []string{"Provider", "Kline", "Spot", "Instrument", "Profile", "Financial", "News"}
	width := firstColWidth + (len(cols)-1)*colWidth
	height := headerHeight + len(names)*rowHeight

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" font-family="Arial, sans-serif">`, width, height))
	buf.WriteString(`<style>text { dominant-baseline: middle; text-anchor: middle; } .left { text-anchor: start; }</style>`)

	// Background
	buf.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, width, height))

	// Header
	buf.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#f6f8fa"/>`, width, headerHeight))
	for i, col := range cols {
		x := 0
		if i > 0 {
			x = firstColWidth + (i-1)*colWidth
		}

		w := colWidth
		if i == 0 {
			w = firstColWidth
		}

		textX := x + w/2
		anchor := ""
		if i == 0 {
			textX = x + padding
			anchor = `class="left"`
		}

		buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="%d" font-weight="bold" fill="#24292e" %s>%s</text>`,
			textX, headerHeight/2, headerFontSize, anchor, col))
	}

	// Rows
	for i, name := range names {
		y := headerHeight + i*rowHeight
		caps := providers[name]

		// Alternating background
		if i%2 == 1 {
			buf.WriteString(fmt.Sprintf(`<rect y="%d" width="%d" height="%d" fill="#fcfcfc"/>`, y, width, rowHeight))
		}

		// Border bottom
		buf.WriteString(fmt.Sprintf(`<line x1="0" y1="%d" x2="%d" y2="%d" stroke="#eaecef" stroke-width="1"/>`, y+rowHeight, width, y+rowHeight))

		// Provider Name
		displayName := providerNames[name]
		if displayName == "" {
			displayName = strings.Title(name)
		} else {
			// Strip (Description) for SVG to save space if needed, or keep it
			// Keeping it for now
		}

		buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="%d" fill="#24292e" class="left">%s</text>`,
			padding, y+rowHeight/2, fontSize, displayName))

		// Checks
		checks := []bool{caps.Kline, caps.Spot, caps.Instrument, caps.Profile, caps.Financial, caps.Announcement}
		for j, checked := range checks {
			x := firstColWidth + j*colWidth
			text := "-"
			color := "#6a737d"
			if checked {
				text = "✔"
				color = "#28a745"
			}
			buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="%d" fill="%s">%s</text>`,
				x+colWidth/2, y+rowHeight/2, fontSize, color, text))
		}
	}

	buf.WriteString(`</svg>`)

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func generateBadges(providers map[string]*Capabilities) string {
	var buf bytes.Buffer

	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		displayName := providerNames[name]
		if displayName == "" {
			displayName = strings.Title(name)
		}

		isBeta := betaProviders[name]

		var label, message, color string
		label = displayName

		if isBeta {
			message = "🟡 Beta"
			color = "yellow"
		} else {
			markets := providerMarkets[name]
			if len(markets) == 0 {
				message = "✓ Ready"
			} else {
				message = "✓ " + strings.Join(markets, " | ")
			}
			color = "brightgreen"
		}

		// URL encode message
		message = strings.ReplaceAll(message, "-", "--") // Escape dash
		message = strings.ReplaceAll(message, "_", "__") // Escape underscore
		// Simple replacements for common chars in badges, but better to use url query usually.
		// Shields.io path format: /badge/<LABEL>-<MESSAGE>-<COLOR>
		// We need to query escape the message part mostly.

		// Manual encoding for shields.io path style
		// Spaces -> %20
		// | -> %7C
		encodedMsg := url.PathEscape(message)
		// Shields.io treats spaces as %20.

		badgeURL := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s",
			url.PathEscape(label), encodedMsg, color)

		buf.WriteString(fmt.Sprintf("![](%s) ", badgeURL))
	}

	return buf.String()
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
			displayName = "**" + displayName + "**"
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

func updateFile(path string, startMarker, endMarker, content string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %s: %v\n", path, err)
		return
	}

	sData := string(data)
	startIdx := strings.Index(sData, startMarker)
	endIdx := strings.Index(sData, endMarker)

	if startIdx == -1 || endIdx == -1 {
		fmt.Printf("Markers %s...%s not found in %s, skipping update.\n", startMarker, endMarker, path)
		return
	}

	newContent := sData[:startIdx+len(startMarker)] + "\n\n" + content + "\n\n" + sData[endIdx:]

	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Could not write file %s: %v\n", path, err)
	} else {
		fmt.Printf("Updated %s\n", path)
	}
}
