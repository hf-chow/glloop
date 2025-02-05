package llm

import (
	"context"
	"net/http"

	"github.com/hf-chow/glloop/internal/config"
	"github.com/tmc/langchaingo/documentloaders"
	emb "github.com/tmc/langchaingo/embeddings/huggingface"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
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

func addToVectorStore(docs []schema.Document, model string) error {
	vstore, err := getVectorStore(model)
	if err != nil {
		return err
	}
	_, err = vstore.AddDocuments(context.Background(), docs)
	return nil
}

func getDocsFromUrl(src string) ([]schema.Document, error) {
	resp, err := http.Get(src)
	if err != nil {
		return []schema.Document{}, err
	}

	defer resp.Body.Close()

	docs, err := documentloaders.NewHTML(resp.Body).LoadAndSplit(
		context.Background(), textsplitter.NewRecursiveCharacter(),
	)
	if err != nil {
		return []schema.Document{}, err
	}

	return docs, nil
}
