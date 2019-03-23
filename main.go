package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jnatalzia/gralhund-bot/commands"
	"github.com/jnatalzia/gralhund-bot/giphy"
	"github.com/jnatalzia/gralhund-bot/resizer"
	"github.com/jnatalzia/gralhund-bot/utils"
)

var token string
var RATING = "pg-13"
var DEBUG = strings.ToLower(os.Getenv("DEBUG")) == "true"
var GUILD_ID = "554402794711416838"

var redisClient = utils.RedisClient

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func checkDebug() {
	if !DEBUG {
		DEBUG = strings.ToLower(os.Getenv("DEBUG")) == "t"
	}
	fmt.Println("Debug is " + strconv.FormatBool(DEBUG))
}

var giphyClient = giphy.NewClient(&giphy.ClientOptions{})

var imageStorePath = os.Getenv("IMAGEPATH")

var imageResizer *resizer.Resizer

func main() {
	checkDebug()
	if imageStorePath == "" {
		imageStorePath = "/tmp"
	}
	imageResizer = resizer.NewResizer(imageStorePath)

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
	botName := "gralhund"
	if DEBUG == true {
		botName = "gralhund-test"
	}
	if strings.HasPrefix(lowerContent, botName+" ") {
		var trimmedMessage = lowerContent[9:]

		// Find the channel that the message came from.
		_, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// do something cool
		re := regexp.MustCompile("gif me ?([a-z0-9 ]+)?")
		var gifRequest string
		for _, match := range re.FindAllStringSubmatch(trimmedMessage, -1) {
			gifRequest = match[1]
			gifUrl, _ := getGif(gifRequest)
			s.ChannelMessageSend(m.ChannelID, gifUrl)
		}

		re = regexp.MustCompile("make emoji from \"([0-9A-Za-z:/._\\-~%&?]+)\" with name \"([a-z_][0-9a-z_]+)\"")
		for _, match := range re.FindAllStringSubmatch(trimmedMessage, -1) {
			fmt.Println("Adding emoji")
			urlPath := match[1]
			newName := match[2]

			p, err := imageResizer.DownloadImage(urlPath)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was an issue downloading your image :(")
				return
			}
			baseSixFourData, err := imageResizer.ResizeImage(p)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was an issue resizing your image: "+err.Error())
				return
			}
			_, err = s.GuildEmojiCreate(GUILD_ID, newName, baseSixFourData, []string{})

			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was an issue creating your image :(")
				return
			}

			s.ChannelMessageSend(m.ChannelID, "Emoji created! 🎉")
		}

		reMatch, _ := regexp.Match("ping", []byte(trimmedMessage))
		if reMatch {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		reMatch, _ = regexp.Match("help", []byte(trimmedMessage))
		if reMatch {
			hi := commands.ListHelp()

			s.ChannelMessageSend(m.ChannelID, "Gralhund knows the following tricks")
			s.ChannelMessageSend(m.ChannelID, "<term> denotes your input, ? means that portion is optional")
			for _, element := range hi {
				// element is the element from someSlice for where we are
				s.ChannelMessageSend(m.ChannelID, "`"+element.Command+"`"+": "+element.Description)
			}
		}

		re = regexp.MustCompile("give ([0-9]+) points to <@([0-9]+)>")
		for _, match := range re.FindAllStringSubmatch(trimmedMessage, -1) {
			fmt.Println("Attempting to give user points")
			numPoints := match[1]
			username := match[2]

			intPointVal, _ := strconv.Atoi(numPoints)
			message, err := commands.GivePointsToUser(username, intPointVal, m.Author.ID)

			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was a problem giving that user points. "+err.Error())
			}

			s.ChannelMessageSend(m.ChannelID, message)
		}

		re = regexp.MustCompile("take ([0-9]+) points from <@([0-9]+)>")
		for _, match := range re.FindAllStringSubmatch(trimmedMessage, -1) {
			fmt.Println("Attempting to take away user points")
			numPoints := match[1]
			username := match[2]

			intPointVal, _ := strconv.Atoi(numPoints)
			message, err := commands.TakePointsFromUser(username, intPointVal, m.Author.ID)

			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was a problem taking that user's points. "+err.Error())
			}

			s.ChannelMessageSend(m.ChannelID, message)
		}

		reMatch, _ = regexp.Match("show point leaderboard", []byte(trimmedMessage))
		if reMatch {
			leaderboard, err := commands.GetPointLeaderBoard(s)

			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "There was an issue generating the leaderboard. "+err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Current Leaderboard: ")

			for _, entry := range leaderboard {
				s.ChannelMessageSend(m.ChannelID, entry.Username+": "+strconv.Itoa(entry.Points)+" points")
			}
		}
	}
}

func getGif(keyword string) (string, error) {
	var g *giphy.Gif
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
