package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCLISyndication(t *testing.T) {
	out := captureStdout(func() {
		handleSyndication([]string{"match", "darlington-new-nuclear-project-unit-1"})
	})
	if !strings.Contains(out, "Institutional Syndication Consortium") {
		t.Fatalf("expected syndication output, got: %s", out)
	}
	if !strings.Contains(out, "Audit Hash:") {
		t.Fatalf("expected audit hash in output, got: %s", out)
	}

	outIndig := captureStdout(func() {
		handleSyndication([]string{"indigenous", "Ring of Fire Test Corridor"})
	})
	if !strings.Contains(outIndig, "Multi-Nation Indigenous Equity Syndicate") {
		t.Fatalf("expected indigenous syndicate output, got: %s", outIndig)
	}
	if !strings.Contains(outIndig, "Marten Falls First Nation") {
		t.Fatalf("expected Marten Falls in output, got: %s", outIndig)
	}
}

func TestCLIOfftake(t *testing.T) {
	out := captureStdout(func() {
		handleOfftake([]string{"list"})
	})
	if !strings.Contains(out, "Commercial Offtake & Clean Power Purchase Agreements") {
		t.Fatalf("expected offtake table, got: %s", out)
	}
	if !strings.Contains(out, "Amazon Web Services") {
		t.Fatalf("expected AWS offtake, got: %s", out)
	}
}

func TestCLICorridor(t *testing.T) {
	outRoute := captureStdout(func() {
		handleCorridor([]string{"route", "--mode", "pipeline"})
	})
	if !strings.Contains(outRoute, "Linear Right-of-Way (RoW) Pathfinding Evaluation") {
		t.Fatalf("expected route evaluation, got: %s", outRoute)
	}
	if !strings.Contains(outRoute, "CLEAN_HYDROGEN_PIPELINE") {
		t.Fatalf("expected clean hydrogen pipeline type, got: %s", outRoute)
	}

	outPort := captureStdout(func() {
		handleCorridor([]string{"port"})
	})
	if !strings.Contains(outPort, "Strategic Maritime Gateways") {
		t.Fatalf("expected port table, got: %s", outPort)
	}
	if !strings.Contains(outPort, "PORT_OF_PRINCE_RUPERT") {
		t.Fatalf("expected Prince Rupert port, got: %s", outPort)
	}
}

func TestCLIFinance(t *testing.T) {
	outSim := captureStdout(func() {
		handleFinance([]string{"simulate", "darlington-new-nuclear-project-unit-1", "--runs", "500"})
	})
	if !strings.Contains(outSim, "Stochastic Project Finance Monte Carlo Simulation") {
		t.Fatalf("expected finance simulation, got: %s", outSim)
	}
	if !strings.Contains(outSim, "Synthetic Credit Rating:") {
		t.Fatalf("expected credit rating, got: %s", outSim)
	}

	outTax := captureStdout(func() {
		handleFinance([]string{"cleantax", "darlington-new-nuclear-project-unit-1"})
	})
	if !strings.Contains(outTax, "Clean Economy Tax Credit & CCfD Underwriting") {
		t.Fatalf("expected clean tax output, got: %s", outTax)
	}
	if !strings.Contains(outTax, "CLEAN_ELECTRICITY_ITC_15") {
		t.Fatalf("expected clean electricity ITC, got: %s", outTax)
	}
}

func TestCLIExportMemoAndGeoJSON(t *testing.T) {
	outMemo := captureStdout(func() {
		handleExport([]string{"memo", "darlington-new-nuclear-project-unit-1"})
	})
	if !strings.Contains(outMemo, "MEMORANDUM_TO_CABINET") {
		t.Fatalf("expected memorandum to cabinet, got: %s", outMemo)
	}

	outGeo := captureStdout(func() {
		handleExport([]string{"geojson", "darlington-new-nuclear-project-unit-1"})
	})
	if !strings.Contains(outGeo, "FeatureCollection") {
		t.Fatalf("expected geojson FeatureCollection, got: %s", outGeo)
	}
}
