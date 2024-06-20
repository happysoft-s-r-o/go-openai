package openai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type V2ThreadRequest struct {
	Messages []V2ThreadMessage `json:"messages,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty"`
}

type V2ThreadMessageRole ThreadMessageRole

const (
	V2ThreadMessageRoleUser      V2ThreadMessageRole = V2ThreadMessageRole(ThreadMessageRoleUser)
	V2ThreadMessageRoleAssistant V2ThreadMessageRole = "assistant"
)

type V2Attachment struct {
	FileID string             `json:"file_id,omitempty"`
	Tools  []V2AttachmentTool `json:"tools,omitempty"`
}

type V2AttachmentTool struct {
	Type ToolType `json:"type"`
}

type V2ThreadMessage struct {
	Role         V2ThreadMessageRole `json:"role"`
	Content      string              `json:"content"`
	Attachments  []V2Attachment      `json:"attachments,omitempty"`
	ToolResource []V2ToolResource    `json:"tool_resources,omitempty"`
	Metadata     map[string]any      `json:"metadata,omitempty"`
}

type V2CodeInterpreter struct {
	FileIDs []string `json:"file_ids"`
}

type V2FileSearch struct {
	VectorStoreIDs []string      `json:"vector_store_ids"`
	VectorStores   []VectorStore `json:"vector_stores"`
}

type V2ToolResource struct {
	CodeInterpreter *V2CodeInterpreter `json:"code_interpreter"`
	FileSearch      *V2FileSearch      `json:"file_search"`
}

type V2Thread struct {
	Thread

	ToolResources V2ToolResource `json:"tool_resources"`
}

// CreateThread creates a new thread.
func (c *V2Client) CreateThread(ctx context.Context, request V2ThreadRequest) (response V2Thread, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(threadsSuffix), withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// RetrieveThread retrieves a thread.
func (c *V2Client) RetrieveThread(ctx context.Context, threadID string) (response V2Thread, err error) {
	urlSuffix := threadsSuffix + "/" + threadID
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

type V2MessageRequest struct {
	MessageRequest
	Role        V2ThreadMessageRole `json:"role"`
	Attachments []V2Attachment      `json:"attachments,omitempty"`
}

type V2IncopmleteDetails struct {
	Reason string `json:"reason"`
}

type V2ContentImageFile struct {
	FileID string `json:"file_id"`
	Detail string `json:"detail"`
}

type V2ContentImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

type V2ContentMessageTextAnnotationFileCitation struct {
	FileID string `json:"file_id"`
	Quote  string `json:"quote"`
}

type V2ContentMessageTextAnnotationFilePath struct {
	FileID string `json:"file_id"`
}

type V2ContentMessageTextAnnotationType string

const (
	V2ContentMessageTextAnnotationTypeFilePath     V2ContentMessageTextAnnotationType = "file_path"
	V2ContentMessageTextAnnotationTypeFileCitation V2ContentMessageTextAnnotationType = "file_citation"
)

type V2ContentMessageTextAnnotation struct {
	Type         V2ContentMessageTextAnnotationType         `json:"type"`
	Text         string                                     `json:"text"`
	Start        int                                        `json:"start_index"`
	End          int                                        `json:"end_index"`
	FileCitation V2ContentMessageTextAnnotationFileCitation `json:"file_citation,omitempty"`
	FilePath     V2ContentMessageTextAnnotationFilePath     `json:"file_path,omitempty"`
}

type V2ContentMessageText struct {
	MessageText
	Annotations []V2ContentMessageTextAnnotation `json:"annotations"`
}

type V2MessageContentType string

const (
	V2MessageContentTypeText      V2MessageContentType = "text"
	V2MessageContentTypeImageFile V2MessageContentType = "image_file"
	V2MessageContentTypeImageUrl  V2MessageContentType = "image_url"
)

type V2MessageContent struct {
	MessageContent
	Type      V2MessageContentType  `json:"type"`
	Text      *V2ContentMessageText `json:"text,omitempty"`
	ImageFile *V2ContentImageFile   `json:"image_file,omitempty"`
	ImageUrl  *V2ContentImageUrl    `json:"image_url,omitempty"`
}

type V2Message struct {
	Message
	Status            string              `json:"status"`
	IncompleteDetails V2IncopmleteDetails `json:"incomplete_details"`
	CompleteAt        int                 `json:"completed_at"`
	IncompleteAt      int                 `json:"incomplete_at"`
	Content           []V2MessageContent  `json:"content"`

	Attachments []V2Attachment `json:"attachments"`
}

// CreateMessage creates a new message.
func (c *V2Client) CreateMessage(ctx context.Context, threadID string, request V2MessageRequest) (msg V2Message, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/%s", threadID, messagesSuffix)
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &msg)
	return
}

type V2MessagesList struct {
	MessagesList
	Messages []V2Message `json:"data"`
}

// ListMessage fetches all messages in the thread.
func (c *V2Client) ListMessage(ctx context.Context, threadID string,
	limit *int,
	order *string,
	after *string,
	before *string,
) (messages V2MessagesList, err error) {
	urlValues := url.Values{}
	if limit != nil {
		urlValues.Add("limit", fmt.Sprintf("%d", *limit))
	}
	if order != nil {
		urlValues.Add("order", *order)
	}
	if after != nil {
		urlValues.Add("after", *after)
	}
	if before != nil {
		urlValues.Add("before", *before)
	}
	encodedValues := ""
	if len(urlValues) > 0 {
		encodedValues = "?" + urlValues.Encode()
	}

	urlSuffix := fmt.Sprintf("/threads/%s/%s%s", threadID, messagesSuffix, encodedValues)
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &messages)
	return
}

type V2RunRequest struct {
	RunRequest
	AdditionalMessages []V2ThreadMessage `json:"additional_messages"`
	TopP               float64           `json:"top_p,omitempty"`
	Stream             bool              `json:"stream,omitempty"`
	ToolChoice         string            `json:"tool_choice,omitempty"`
	ResponseFormat     string            `json:"response_format,omitempty"`
}

type V2RunIncopmleteDetails struct {
	Reason string `json:"reason"`
}

type V2Run struct {
	Run
	IncompleteDetails *V2RunIncopmleteDetails `json:"incomplete_details"`
	TopP              *float64                `json:"top_p"`
	ResponseFormat    *string                 `json:"response_format"`
	ToolChoice        *string                 `json:"tool_choice"`
}

// CreateRun creates a new run.
func (c *V2Client) CreateRun(
	ctx context.Context,
	threadID string,
	request V2RunRequest,
) (response V2Run, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/runs", threadID)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// RetrieveRun retrieves a run.
func (c *V2Client) RetrieveRun(
	ctx context.Context,
	threadID string,
	runID string,
) (response V2Run, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/runs/%s", threadID, runID)
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

type V2SubmitToolOutputsRequest struct {
	SubmitToolOutputsRequest
	Stream bool `json:"stream,omitempty"`
}

// SubmitToolOutputs submits tool outputs.
func (c *V2Client) SubmitToolOutputs(
	ctx context.Context,
	threadID string,
	runID string,
	request V2SubmitToolOutputsRequest) (response V2Run, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/runs/%s/submit_tool_outputs", threadID, runID)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// CancelRun cancels a run.
func (c *V2Client) CancelRun(
	ctx context.Context,
	threadID string,
	runID string) (response V2Run, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/runs/%s/cancel", threadID, runID)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
