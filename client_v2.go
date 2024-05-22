package openai

type AssistantVersion string

const AssistantVersionV2 AssistantVersion = "v2"

type optionFnc func(*ClientConfig)

func WithAssistantVersion(version AssistantVersion) optionFnc {
	return func(config *ClientConfig) {
		config.AssistantVersion = string(version)
	}
}

type V2Client struct {
	*Client
}

// NewClient creates new OpenAI API client.
func NewClientV2(authToken string, opts ...optionFnc) *V2Client {
	config := DefaultConfig(authToken)

	opts = append(opts, WithAssistantVersion(AssistantVersionV2))

	for _, opt := range opts {
		opt(&config)
	}

	return &V2Client{
		Client: NewClientWithConfig(config),
	}
}
