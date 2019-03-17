package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var buffer = make([][]byte, 0)

type ClientOptions struct {
	ApiKey      string
	ApiEndpoint string
	HttpClient  *http.Client
}

type Client struct {
	apiKey      string
	apiEndpoint string
	httpClient  *http.Client
}

func NewClient(co *ClientOptions) *Client {

	client := &Client{}

	// set default api key if not set
	if co.ApiKey == "" {
		client.apiKey = "dc6zaTOxFJmzC"
	} else {
		client.apiKey = co.ApiKey
	}

	// set default endpoint if not set (mostly used for overriding the server
	// url during test runs)
	if co.ApiEndpoint == "" {
		client.apiEndpoint = "https://api.giphy.com/v1"
	} else {
		client.apiEndpoint = strings.TrimRight(co.ApiEndpoint, "/")
	}

	// set default http client if not set. Useful in situations where you need
	// special behaviour or aren't able to use a standard `http.Client`
	// instance (like on appengine).
	if co.HttpClient == nil {
		client.httpClient = &http.Client{}
	} else {
		client.httpClient = co.HttpClient
	}

	return client
}

func (c *Client) makeRequest(suffix string, qs *url.Values, ds interface{}) error {

	// inject configured API key into url
	qs.Set("api_key", c.apiKey)

	// execute HTTP request
	u := fmt.Sprintf("%s/%s?%s", c.apiEndpoint, suffix, qs.Encode())
	resp, err := c.httpClient.Get(u)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	// unmarshal HTTP response as JSON into the provided data structure
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(ds); err != nil {
		return err
	}

	return nil
}

// translateCommon provides common search functionality for both GIF and
// sticker translate endpoints.
func (c *Client) translateCommon(q string, rating string, urlFragment string) (*Gif, error) {

	// check that a query string was actually passed in
	if len(q) < 1 {
		err := errors.New("`q` parameter must not be empty.")
		return nil, err
	}

	// check that the given rating is valid
	if err := validateRating(rating); err != nil {
		return nil, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	qs.Set("s", q)
	if rating != "" {
		qs.Set("rating", rating)
	}

	// construct and execute the HTTP request
	sr := &singleResult{}
	if err := c.makeRequest(urlFragment, qs, sr); err != nil {
		return nil, err
	}

	return sr.Data, nil
}

// TranslateGif is prototype endpoint for using Giphy as a translation engine
// for a GIF dialect. The translate API draws on search, but uses the Giphy
// "special sauce" to handle translating from one vocabulary to another. In
// this case, words and phrases to GIFs. Returns a single GIF from the Giphy
// API.
func (c *Client) TranslateGif(q string, rating string) (*Gif, error) {
	return c.translateCommon(q, rating, "gifs/translate")
}

// get a random gif
func (c *Client) randomCommon(rating string, urlFragment string) (*Gif, error) {

	// check that the given rating is valid
	if err := validateRating(rating); err != nil {
		return nil, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	if rating != "" {
		qs.Set("rating", rating)
	}

	// construct and execute the HTTP request
	sr := &singleResult{}
	if err := c.makeRequest(urlFragment, qs, sr); err != nil {
		return nil, err
	}

	return sr.Data, nil
}

// RandomGif is prototype endpoint for using Giphy as a translation engine
// for a GIF dialect.
func (c *Client) RandomGif(rating string) (*Gif, error) {
	return c.randomCommon(rating, "gifs/random")
}

// validateRating checks if the given string matches an allowed value for the
// giphy API's `rating` parameter.
func validateRating(r string) error {

	r = strings.ToLower(r)
	switch r {
	case "y", "g", "pg", "pg-13", "r", "":
		return nil
	}

	fmtString := "\"%s\" is not a valid value for `rating`, must be one of \"y\", \"g\", \"pg\", \"pg-13\", \"r\" or \"\""
	return fmt.Errorf(fmtString, r)
}
