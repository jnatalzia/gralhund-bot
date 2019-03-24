package commands

type helpInfo struct {
	Command     string
	Description string
}

var allHelpInfo = []helpInfo{
	helpInfo{Command: "help", Description: "Displays this list."},
	helpInfo{Command: "ping", Description: "Returns a `pong` if Gralhund is up and running."},
	helpInfo{Command: "gif me <search_term>?", Description: "Returns a gif either random or based on the provided search term."},
	helpInfo{Command: "make emoji from \"<image_url>\" with name \"<emoji_name>\"", Description: "Creates an emoji based on the provided image link. Supports jpeg and png."},
	helpInfo{Command: "give <pts> points to @<username>", Description: "Gives points to that user."},
	helpInfo{Command: "take <pts> points from @<username>", Description: "Takes points from that user."},
	helpInfo{Command: "show point leaderboard", Description: "Shows top points."},
}

func ListHelp() []helpInfo {
	return allHelpInfo
}
