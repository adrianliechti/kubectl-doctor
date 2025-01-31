package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
)

var namespace string

var openaiURL string
var openaiToken string
var openaiModel string

func init() {
	url := os.Getenv("OPENAI_BASE_URL")

	if url == "" {
		url = "https://api.openai.com/v1/"
	}

	token := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")

	if model == "" {
		model = "gpt-4o-mini"
	}

	flag.StringVar(&namespace, "namespace", "", "namespace to analyze")

	flag.StringVar(&openaiURL, "url", url, "openai endpoint")
	flag.StringVar(&openaiToken, "token", token, "openai api key")
	flag.StringVar(&openaiModel, "model", model, "openai model")

	flag.Parse()
}

func main() {
	backend := "LLM"

	configAI := ai.AIConfiguration{
		DefaultProvider: backend,

		Providers: []ai.AIProvider{
			{
				Name: backend,

				BaseURL:  openaiURL,
				Password: openaiToken,

				Model: openaiModel,

				MaxTokens:   2048,
				Temperature: 0.7,
				TopP:        0.5,
				TopK:        50,
			},
		},
	}

	viper.Set("ai", configAI)

	var (
		explain         = true
		output          = "text" // json
		filters         = []string{}
		language        = "english"
		nocache         = false
		labelSelector   = ""
		anonymize       = false
		maxConcurrency  = 10
		withDoc         = false
		interactiveMode = false
		customAnalysis  = false
		customHeaders   = []string{}
		withStats       = false
	)

	config, err := analysis.NewAnalysis(
		backend,
		language,
		filters,
		namespace,
		labelSelector,
		nocache,
		explain,
		maxConcurrency,
		withDoc,
		interactiveMode,
		customHeaders,
		withStats,
	)

	if err != nil {
		panic(err)
	}

	defer config.Close()

	if customAnalysis {
		config.RunCustomAnalysis()
	}

	config.RunAnalysis()

	if explain {
		if err := config.GetAIResults(output, anonymize); err != nil {
			panic(err)
		}
	}

	output_data, err := config.PrintOutput(output)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(output_data))
}
