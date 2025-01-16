package llm

import (
	"context"
	"fmt"

	"github.com/hf-chow/glloop/internal/config"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/embeddings/huggingface"
	"github.com/tmc/langchaingo/llms/huggingface"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

type Huggingface struct {
	Model 			string
	Task 			string
	StripNewLines 	bool
	BatchSize 		int
}

func getVector() (vectorstores.VectorStore, error) {

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error when reading config: %v\n", err)
	}
	
	embeddingsModel, err := huggingface.NewHuggingface(
		//TODO implement token in the config
		huggingface.WithToken(cfg.HuggingFaceToken),
		//Default is GPT2
		huggingface.WithModel(),
	)

	vstore, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(cfg.DBURL),
		pgvector.WithEmbedder(embeddingsModel),
	)
	return vstore, nil
}

