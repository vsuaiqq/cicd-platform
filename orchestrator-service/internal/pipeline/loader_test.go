package pipeline

import "testing"

func TestLoadBytes_validPipeline(t *testing.T) {
	data := []byte(`
name: demo
on:
  push:
    branches: [main]
env:
  APP: hello
jobs:
  lint:
    image: golang:1.25
    steps:
      - name: vet
        run: go vet ./...
  build:
    needs: [lint]
    steps:
      - name: build
        run: go build ./...
`)
	pl, err := LoadBytes(data)
	if err != nil {
		t.Fatalf("LoadBytes: %v", err)
	}
	if pl.Name != "demo" {
		t.Fatalf("name = %q", pl.Name)
	}
	if pl.Env["APP"] != "hello" {
		t.Fatalf("env APP = %q", pl.Env["APP"])
	}
	if len(pl.Jobs) != 2 {
		t.Fatalf("jobs = %d", len(pl.Jobs))
	}
	if pl.Jobs["build"].Needs[0] != "lint" {
		t.Fatalf("build needs = %v", pl.Jobs["build"].Needs)
	}
}

func TestLoadBytes_rejectsUnknownNeed(t *testing.T) {
	data := []byte(`
jobs:
  deploy:
    needs: [missing]
    steps:
      - name: deploy
        run: echo deploy
`)
	_, err := LoadBytes(data)
	if err == nil {
		t.Fatal("expected error for unknown need")
	}
}

func TestLoadBytes_rejectsEmptyJobs(t *testing.T) {
	data := []byte(`name: empty`)
	_, err := LoadBytes(data)
	if err == nil {
		t.Fatal("expected error for empty jobs")
	}
}

func TestLoadBytes_performanceGate(t *testing.T) {
	data := []byte(`
jobs:
  load-test:
    image: golang:1.25
    steps:
      - name: load
        run: ./scripts/load-test.sh
  performance-gate:
    needs: [load-test]
    performance_gate:
      source_job: load-test
      baseline:
        min_samples: 3
`)
	pl, err := LoadBytes(data)
	if err != nil {
		t.Fatalf("LoadBytes: %v", err)
	}
	if !pl.Jobs["performance-gate"].IsPerformanceGate() {
		t.Fatal("expected performance gate job")
	}
}

func TestLoadBytes_rejectsGateWithoutNeed(t *testing.T) {
	data := []byte(`
jobs:
  load-test:
    steps:
      - name: load
        run: echo load
  performance-gate:
    performance_gate:
      source_job: load-test
`)
	_, err := LoadBytes(data)
	if err == nil {
		t.Fatal("expected error when source job not in needs")
	}
}

func TestBranchAllowed(t *testing.T) {
	pl := &Pipeline{
		On: &Trigger{
			Push: &PushTrigger{Branches: []string{"main", "develop"}},
		},
	}
	if !pl.BranchAllowed("main") {
		t.Fatal("main should be allowed")
	}
	if pl.BranchAllowed("feature/x") {
		t.Fatal("feature/x should be rejected")
	}
	if !(&Pipeline{}).BranchAllowed("anything") {
		t.Fatal("empty trigger should allow any branch")
	}
}
