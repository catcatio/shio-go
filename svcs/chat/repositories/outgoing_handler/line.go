package outgoing_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// APIEndpoint constants
const (
	APIEndpointBase = "https://api.line.me"

	APIEndpointPushMessage           = "/v2/bot/message/push"
	APIEndpointReplyMessage          = "/v2/bot/message/reply"
	APIEndpointMulticast             = "/v2/bot/message/multicast"
	APIEndpointGetMessageContent     = "/v2/bot/message/%s/content"
	APIEndpointGetMessageQuota       = "/v2/bot/message/quota"
	APIEndpointLeaveGroup            = "/v2/bot/group/%s/leave"
	APIEndpointLeaveRoom             = "/v2/bot/room/%s/leave"
	APIEndpointGetProfile            = "/v2/bot/profile/%s"
	APIEndpointGetGroupMemberProfile = "/v2/bot/group/%s/member/%s"
	APIEndpointGetRoomMemberProfile  = "/v2/bot/room/%s/member/%s"
	APIEndpointGetGroupMemberIDs     = "/v2/bot/group/%s/members/ids"
	APIEndpointGetRoomMemberIDs      = "/v2/bot/room/%s/members/ids"
	APIEndpointCreateRichMenu        = "/v2/bot/richmenu"
	APIEndpointGetRichMenu           = "/v2/bot/richmenu/%s"
	APIEndpointListRichMenu          = "/v2/bot/richmenu/list"
	APIEndpointDeleteRichMenu        = "/v2/bot/richmenu/%s"
	APIEndpointGetUserRichMenu       = "/v2/bot/user/%s/richmenu"
	APIEndpointLinkUserRichMenu      = "/v2/bot/user/%s/richmenu/%s"
	APIEndpointUnlinkUserRichMenu    = "/v2/bot/user/%s/richmenu"
	APIEndpointSetDefaultRichMenu    = "/v2/bot/user/all/richmenu/%s"
	APIEndpointDefaultRichMenu       = "/v2/bot/user/all/richmenu"   // Get: GET / Delete: DELETE
	APIEndpointDownloadRichMenuImage = "/v2/bot/richmenu/%s/content" // Download: GET / Upload: POST
	APIEndpointUploadRichMenuImage   = "/v2/bot/richmenu/%s/content" // Download: GET / Upload: POST
	APIEndpointBulkLinkRichMenu      = "/v2/bot/richmenu/bulk/link"
	APIEndpointBulkUnlinkRichMenu    = "/v2/bot/richmenu/bulk/unlink"

	APIEndpointGetAllLIFFApps = "/liff/v1/apps"
	APIEndpointAddLIFFApp     = "/liff/v1/apps"
	APIEndpointUpdateLIFFApp  = "/liff/v1/apps/%s/view"
	APIEndpointDeleteLIFFApp  = "/liff/v1/apps/%s"

	APIEndpointLinkToken = "/v2/bot/user/%s/linkToken"

	APIEndpointGetMessageDelivery = "/v2/bot/message/delivery/%s"
)

// LineClient type
type LineClient struct {
	channelSecret string
	channelToken  string
	endpointBase  *url.URL     // default APIEndpointBase
	httpClient    *http.Client // default http.DefaultClient
}

func NewLineClient(channelSecret string, channelToken string) *LineClient {
	u, err := url.ParseRequestURI(APIEndpointBase)
	if err != nil {
		panic(err)
	}

	return &LineClient{
		channelSecret: channelSecret,
		channelToken:  channelToken,
		httpClient:    http.DefaultClient,
		endpointBase:  u,
	}
}

// ClientOption type
type ClientOption func(*LineClient) error

// New returns a new bot client instance.
func New(channelSecret, channelToken string, options ...ClientOption) (*LineClient, error) {
	if channelSecret == "" {
		return nil, errors.New("missing channel secret")
	}
	if channelToken == "" {
		return nil, errors.New("missing channel access token")
	}
	c := &LineClient{
		channelSecret: channelSecret,
		channelToken:  channelToken,
		httpClient:    http.DefaultClient,
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.endpointBase == nil {
		u, err := url.ParseRequestURI(APIEndpointBase)
		if err != nil {
			return nil, err
		}
		c.endpointBase = u
	}
	return c, nil
}

// WithHTTPClient function
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *LineClient) error {
		client.httpClient = c
		return nil
	}
}

// WithEndpointBase function
func WithEndpointBase(endpointBase string) ClientOption {
	return func(client *LineClient) error {
		u, err := url.ParseRequestURI(endpointBase)
		if err != nil {
			return err
		}
		client.endpointBase = u
		return nil
	}
}

func (client *LineClient) url(endpoint string) string {
	u := *client.endpointBase
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (client *LineClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+client.channelToken)
	req.Header.Set("User-Agent", "LINE-ME")
	if ctx != nil {
		res, err := client.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
			}
		}

		return res, err
	}
	return client.httpClient.Do(req)

}

func (client *LineClient) get(ctx context.Context, endpoint string, query url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", client.url(endpoint), nil)
	if err != nil {
		return nil, err
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	return client.do(ctx, req)
}

func (client *LineClient) post(ctx context.Context, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", client.url(endpoint), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return client.do(ctx, req)
}

func (client *LineClient) put(ctx context.Context, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", client.url(endpoint), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return client.do(ctx, req)
}

func (client *LineClient) delete(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", client.url(endpoint), nil)
	if err != nil {
		return nil, err
	}
	return client.do(ctx, req)
}

func closeResponse(res *http.Response) error {
	defer func() { _ = res.Body.Close() }()
	_, err := io.Copy(ioutil.Discard, res.Body)
	return err
}

func (client *LineClient) ReplyRaw(ctx context.Context, replyToken string, messages string) error {
	msg := fmt.Sprintf(`{"replyToken": "%s", "messages": %s}`, replyToken, messages)
	buf := bytes.NewBufferString(msg)
	res, err := client.post(ctx, APIEndpointReplyMessage, buf)
	if err != nil {
		return err
	}
	defer func() { _ = closeResponse(res) }()
	return checkResponse(res)
}

func (client *LineClient) SendRaw(ctx context.Context, recipient string, messages string) error {
	msg := fmt.Sprintf(`{"to": "%s", "messages": %s}`, recipient, messages)
	buf := bytes.NewBufferString(msg)
	res, err := client.post(ctx, APIEndpointPushMessage, buf)
	if err != nil {
		return err
	}
	defer closeResponse(res)
	return checkResponse(res)
}

func (client *LineClient) TryReplyRaw(ctx context.Context, replyToken string, recipient string, messages string) error {
	if replyToken != "" {
		if err := client.ReplyRaw(ctx, replyToken, messages); err == nil {
			return nil
		}
	}

	return client.SendRaw(ctx, recipient, messages)
}

func checkResponse(res *http.Response) error {
	if isSuccess(res.StatusCode) {
		return nil
	}
	decoder := json.NewDecoder(res.Body)
	result := errorResponse{}
	if err := decoder.Decode(&result); err != nil {
		return &apiError{
			Code: res.StatusCode,
		}
	}
	return &apiError{
		Code:     res.StatusCode,
		Response: &result,
	}
}

func isSuccess(code int) bool {
	return code/100 == 2
}

type errorResponse struct {
	Message string                `json:"message"`
	Details []errorResponseDetail `json:"details"`
}

type errorResponseDetail struct {
	Message  string `json:"message"`
	Property string `json:"property"`
}

// APIError type
type apiError struct {
	Code     int
	Response *errorResponse
}

// Error method
func (e *apiError) Error() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "linebot: APIError %d ", e.Code)
	if e.Response != nil {
		_, _ = fmt.Fprintf(&buf, "%s", e.Response.Message)
		for _, d := range e.Response.Details {
			_, _ = fmt.Fprintf(&buf, "\n[%s] %s", d.Property, d.Message)
		}
	}
	return buf.String()
}
