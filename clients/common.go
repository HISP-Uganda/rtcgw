package clients

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Client struct {
	RestClient *resty.Client
	BaseURL    string
}

type Server struct {
	BaseUrl    string `json:"base_url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthToken  string `json:"auth_token"`
	AuthMethod string `json:"auth_method"`
}

func (c *Client) GetResource(resourcePath string, params map[string]string) (*resty.Response, error) {
	request := c.RestClient.R()

	if params != nil {
		request.SetQueryParams(params)
	}

	resp, err := request.Get(resourcePath)
	if err != nil {
		log.WithError(err).Infof("Error when calling `GetResource`: %v", err)
	}
	return resp, err
}
func (c *Client) PostResource(resourcePath string, params map[string]any, data interface{}) (*resty.Response, error) {
	request := c.RestClient.R()
	// Prepare query parameters
	queryParams := url.Values{}
	// XXX: this ensures that all parameters added via -Q to and command are added
	// newParams := config.CombineMaps(params, config.ParamsMap(config.QueryParams))
	// newParams = config.CombineMaps(newParams, config.ParamsMap(strings.Split(config.QueryParamsString, ",")))

	for key, value := range params {
		switch v := value.(type) {
		case string:
			queryParams.Add(key, v)
		case bool:
			if v {
				queryParams.Add(key, "true")
			} else {
				queryParams.Add(key, "false")
			}
		case []string:
			for _, item := range v {
				queryParams.Add(key, item)
			}
		default:
			return nil, fmt.Errorf("unsupported query parameter type for key %s", key)
		}
	}

	// Set the query parameters
	if len(queryParams) > 0 {
		request.SetQueryParamsFromValues(queryParams)
	}

	resp, err := request.
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(resourcePath)
	if err != nil {
		log.Fatalf("Error when calling `PostResource`: %v", err)
	}
	return resp, err
}

func (c *Client) PutResource(resourcePath string, data interface{}) (*resty.Response, error) {
	resp, err := c.RestClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Put(resourcePath)
	if err != nil {
		log.Errorf("Error when calling `PutResource`: %v", err)
	}
	return resp, err
}

func (c *Client) DeleteResource(resourcePath string) (*resty.Response, error) {
	resp, err := c.RestClient.R().
		Delete(resourcePath)
	if err != nil {
		log.Errorf("Error when calling `DeleteResource`: %v", err)
	}
	return resp, err
}

func (c *Client) PatchResource(resourcePath string, data interface{}) (*resty.Response, error) {
	resp, err := c.RestClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Patch(resourcePath)
	if err != nil {
		log.Errorf("Error when calling `PatchResource`: %v", err)
	}
	return resp, err
}
