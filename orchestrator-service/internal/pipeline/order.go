package pipeline

import (
	"fmt"
	"sort"

	"go.yaml.in/yaml/v3"
)

func JobNamesInOrder(data []byte) ([]string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("pipeline: yaml unmarshal: %w", err)
	}
	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("pipeline: empty document")
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("pipeline: root must be a mapping")
	}
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value != "jobs" {
			continue
		}
		jobsNode := root.Content[i+1]
		if jobsNode.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("pipeline: jobs must be a mapping")
		}
		names := make([]string, 0, len(jobsNode.Content)/2)
		for j := 0; j < len(jobsNode.Content); j += 2 {
			names = append(names, jobsNode.Content[j].Value)
		}
		return names, nil
	}
	return nil, fmt.Errorf("pipeline: no jobs section")
}

func ExecutionOrder(data []byte) ([]string, error) {
	pl, err := LoadBytes(data)
	if err != nil {
		return nil, err
	}
	docOrder, err := JobNamesInOrder(data)
	if err != nil {
		return nil, err
	}

	rank := make(map[string]int, len(docOrder))
	for i, name := range docOrder {
		rank[name] = i
	}

	inDegree := make(map[string]int, len(pl.Jobs))
	dependents := make(map[string][]string, len(pl.Jobs))
	for name, job := range pl.Jobs {
		inDegree[name] = len(job.Needs)
		for _, dep := range job.Needs {
			dependents[dep] = append(dependents[dep], name)
		}
	}

	ready := make([]string, 0, len(pl.Jobs))
	for _, name := range docOrder {
		if inDegree[name] == 0 {
			ready = append(ready, name)
		}
	}

	result := make([]string, 0, len(pl.Jobs))
	for len(ready) > 0 {
		sort.Slice(ready, func(i, j int) bool {
			return rank[ready[i]] < rank[ready[j]]
		})
		current := ready[0]
		ready = ready[1:]
		result = append(result, current)

		for _, dep := range dependents[current] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				ready = append(ready, dep)
			}
		}
	}

	if len(result) != len(pl.Jobs) {
		return nil, fmt.Errorf("pipeline: cycle detected in job dependencies")
	}
	return result, nil
}
