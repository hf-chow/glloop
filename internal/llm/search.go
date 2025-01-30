package llm

import (
	"context"

	"github.com/hf-chow/glloop/internal/config"
)

func rag(model string) error {
	vstore, err := getVectorStore(model)
	if err != nil {
		return err
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	return nil
}


func semanticSearch(query string, model string, maxResults int) error {
	vstore, err := getVectorStore(model)
	if err != nil {
		return err
	}

	results, err := vstore.SimilaritySearch(
		context.Background(),
		query,
		maxResults,
	)
	if err != nil {
		return err
	}
	return nil

}
