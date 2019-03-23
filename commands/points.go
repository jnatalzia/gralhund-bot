package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/jnatalzia/gralhund-bot/utils"
)

func formatPointKey(userID string) string {
	return userID + "__POINTS"
}

func GivePointsToUser(userID string, numPoints int) (message string, err error) {
	formattedUserKey := formatPointKey(userID)
	value, err := utils.RedisClient.Get(formattedUserKey).Result()

	if err == redis.Nil {
		fmt.Println("Points do not exist for user: " + userID + ". Adding.")
		utils.RedisClient.Set(formattedUserKey, numPoints, 0)
	} else if err != nil {
		return "", err
	} else {
		intPoints, _ := strconv.Atoi(value)
		fmt.Println("Current point value for user: ", value)
		newPoints := intPoints + numPoints
		fmt.Println("New point value for user: ", newPoints)
		utils.RedisClient.Set(formattedUserKey, newPoints, 0)
	}

	return "Points successfully awarded", nil
}

type LeaderBoardEntry struct {
	Username string
	Points   int
}

func GetPointLeaderBoard(s *discordgo.Session) ([]LeaderBoardEntry, error) {
	allKeys, _, err := utils.RedisClient.Scan(0, "*__POINTS", 1000).Result()
	if err != nil {
		return nil, err
	}

	result := make([]LeaderBoardEntry, len(allKeys))
	for idx, k := range allKeys {
		value, _ := utils.RedisClient.Get(k).Result()
		userid := strings.Split(k, "__")[0]
		user, err := s.User(userid)

		if err != nil {
			return nil, errors.New("There was an issue finding that user")
		}

		pointCount, _ := strconv.Atoi(value)

		result[idx] = LeaderBoardEntry{
			Username: user.Username,
			Points:   pointCount,
		}
	}

	// sort.Ints(result)
	return result, nil
}
