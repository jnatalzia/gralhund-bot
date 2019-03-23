package commands

type helpInfo struct {
	Command     string
	Description string
}

var allHelpInfo = []helpInfo{
	helpInfo{Command: "help", Description: "Displays this list."},
	helpInfo{Command: "ping", Description: "Returns a `pong` if Gralhund is up and running."},
	helpInfo{Command: "gif me <search_term>?", Description: "Returns a gif either random or based on the provided search term."},
	helpInfo{Command: "make emoji from \"<image_url>\" with name \"<emoji_name>\"", Description: "Creates a gif based on the provided image link. Supports jpeg and png."},
}

func ListHelp() []helpInfo {
	return allHelpInfo
}
