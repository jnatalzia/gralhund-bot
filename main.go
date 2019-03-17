package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	// "strconv"
	"flag"
	// "io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	// "time"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string
var buffer = make([][]byte, 0)
var RATING string = "pg-13"

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

var giphyClient = NewClient(&ClientOptions{})

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

func main() {

	if token == "" {
		fmt.Println("No token provided. Please run: dndbot -t <bot token>")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("DND Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	var lowerContent = strings.ToLower(m.Content)

	if strings.HasPrefix(lowerContent, "gralhund ") {
		var trimmedMessage = lowerContent[9:]
		
		// Find the channel that the message came from.
		_, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// do something cool
		re := regexp.MustCompile("gif me ?([a-z0-9 ]+)?")
		// groupNames := re.SubexpNames()
		var gifRequest string
		for _, match := range re.FindAllStringSubmatch(trimmedMessage, -1) {
			gifRequest = match[1]
			gifUrl, _ := getGif(gifRequest)
			s.ChannelMessageSend(m.ChannelID, gifUrl)
		}
	}
}

func getGif(keyword string) (string, error) {
	var g *Gif
	var err error
	if keyword == "" {
		g, err = giphyClient.RandomGif(RATING)
	} else {
		g, err = giphyClient.TranslateGif(keyword, RATING)
	}

	if err != nil {
		// Could not find channel.
		fmt.Println("ERROR!!!")
		return "", err
	}
	return g.URL, nil
}
