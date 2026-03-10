package dashboard

import (
	"strings"

	"github.com/PeterHiroshi/cfmon/internal/api"
)

func filterWorkers(workers []api.Worker, text string) []api.Worker {
	if text == "" {
		return workers
	}
	lower := strings.ToLower(text)
	var result []api.Worker
	for _, w := range workers {
		if strings.Contains(strings.ToLower(w.Name), lower) {
			result = append(result, w)
		}
	}
	return result
}

func filterContainers(containers []api.Container, text string) []api.Container {
	if text == "" {
		return containers
	}
	lower := strings.ToLower(text)
	var result []api.Container
	for _, c := range containers {
		if strings.Contains(strings.ToLower(c.Name), lower) {
			result = append(result, c)
		}
	}
	return result
}
