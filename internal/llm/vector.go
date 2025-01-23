package llm

import (
	"context"

	"github.com/hf-chow/glloop/internal/config"
	emb "github.com/tmc/langchaingo/embeddings/huggingface"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)


const defaultModel = "sentence-transformers/all-mpnet-base-v2"

func getEmbeddingModel(model string) (*emb.Huggingface, error) {
	embeddingModel, err := emb.NewHuggingface()
	if err != nil {
		return nil, err
	}
	embeddingModel.Model = model

	return embeddingModel, nil

}

func getVectorStore(model string) (vectorstores.VectorStore, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	embeddingModel, err := getEmbeddingModel(model)
	if err !=nil {
		return nil, err
	}

	vstore, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(cfg.DBURL),
		pgvector.WithEmbedder(embeddingModel),
	)
	return vstore, nil
}

func embedDoc(texts []string, embeddingModelName string) ([][]float32, error) {
	model, err := getEmbeddingModel(embeddingModelName)
	if err != nil {
		return [][]float32{}, err
	}

	vector, err := model.EmbedDocuments(context.Background(), texts)
	if err != nil {
		return [][]float32{}, err
	}
	return vector, nil
}
