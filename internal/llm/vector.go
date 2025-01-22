package llm

import (
	"context"
	"fmt"

	"github.com/hf-chow/glloop/internal/config"
	emb "github.com/tmc/langchaingo/embeddings/huggingface"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

type Huggingface struct {
	Model 			string
	Task 			string
	StripNewLines 	bool
	BatchSize 		int
}

func getHuggingFaceEmbeddingModel() {}

func getVectorStore() (vectorstores.VectorStore, error) {

	model := "llama3.2"

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error when reading config: %v\n", err)
	}

	embeddingsModel, err := emb.NewHuggingface()
	embeddingsModel.Model = model

	vstore, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(cfg.DBURL),
		pgvector.WithEmbedder(embeddingsModel),
	)
	return vstore, nil
}

