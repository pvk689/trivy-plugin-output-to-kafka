package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

var severityRank = map[string]int{
	"UNKNOWN":  0,
	"LOW":      1,
	"MEDIUM":   2,
	"HIGH":     3,
	"CRITICAL": 4,
}

type Report struct {
	Results []Result `json:"Results"`
}

type Result struct {
	Target          string          `json:"Target"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

type Vulnerability struct {
	Severity string `json:"Severity"`
}

func Summarize(data []byte, w io.Writer) error {
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return fmt.Errorf("parse report: %w", err)
	}

	byTarget := make(map[string]map[string]int, len(r.Results))
	totals := make(map[string]int)
	severitySet := make(map[string]struct{})

	for _, res := range r.Results {
		counts, ok := byTarget[res.Target]
		if !ok {
			counts = make(map[string]int)
			byTarget[res.Target] = counts
		}
		for _, v := range res.Vulnerabilities {
			counts[v.Severity]++
			totals[v.Severity]++
			severitySet[v.Severity] = struct{}{}
		}
	}

	severities := make([]string, 0, len(severitySet))
	for s := range severitySet {
		severities = append(severities, s)
	}
	sort.Slice(severities, func(i, j int) bool {
		ri, iok := severityRank[severities[i]]
		rj, jok := severityRank[severities[j]]
		switch {
		case iok && jok:
			return ri < rj
		case iok != jok:
			return iok
		default:
			return severities[i] < severities[j]
		}
	})

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(append([]string{"TARGET"}, severities...), "\t"))

	targets := make([]string, 0, len(byTarget))
	for t := range byTarget {
		targets = append(targets, t)
	}
	sort.Strings(targets)

	total := 0
	for _, t := range targets {
		fmt.Fprintln(tw, row(t, severities, byTarget[t]))
	}
	for _, n := range totals {
		total += n
	}
	fmt.Fprintln(tw, row("TOTAL SEVERITY", severities, totals))

	if err := tw.Flush(); err != nil {
		return err
	}
	fmt.Fprintf(w, "Total: %d\n", total)

	return nil
}

func row(label string, severities []string, counts map[string]int) string {
	cells := make([]string, 0, len(severities)+1)
	cells = append(cells, label)
	for _, s := range severities {
		cells = append(cells, fmt.Sprintf("%d", counts[s]))
	}
	return strings.Join(cells, "\t")
}
