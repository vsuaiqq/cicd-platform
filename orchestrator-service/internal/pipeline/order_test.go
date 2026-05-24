package pipeline

import "testing"

func TestExecutionOrder_respectsNeedsAndYAMLTiebreak(t *testing.T) {
	yaml := []byte(`
jobs:
  build:
    needs: [test]
    steps:
      - name: build
        run: echo build
  lint:
    steps:
      - name: lint
        run: echo lint
  test:
    needs: [lint]
    steps:
      - name: test
        run: echo test
`)
	got, err := ExecutionOrder(yaml)
	if err != nil {
		t.Fatalf("ExecutionOrder: %v", err)
	}
	want := []string{"lint", "test", "build"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestExecutionOrder_parallelJobsKeepDocumentOrder(t *testing.T) {
	yaml := []byte(`
jobs:
  beta:
    steps:
      - name: beta
        run: echo beta
  alpha:
    steps:
      - name: alpha
        run: echo alpha
`)
	got, err := ExecutionOrder(yaml)
	if err != nil {
		t.Fatalf("ExecutionOrder: %v", err)
	}
	want := []string{"beta", "alpha"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
