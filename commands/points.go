package commands

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/jnatalzia/gralhund-bot/utils"
)

var MAX_POINTS_PER_DAY = 20

func formatPointKey(userID string, messageGuildID string) string {
	return userID + "-" + messageGuildID + "__POINTS"
}

func formatTimeKey(time string, userID string, messageGuildID string) string {
	return userID + "-" + messageGuildID + "__POINTS_GIVEN__" + time
}

// TODO: Cleanup duplicate logic around checking points given today (reuse the value when the transaction is complete to
//    avoid taking user points before the redis commits the change
func GivePointsToUser(userID string, numPoints int, authorID string, messageGuildID string) (message string, err error) {
	timeNow := time.Now().UTC()
	timeString := strconv.Itoa(timeNow.Year()) + "-" + timeNow.Month().String() + "-" + strconv.Itoa(timeNow.Day())
	timeKey := formatTimeKey(timeString, authorID, messageGuildID)

	timeValue, timeErr := utils.RedisClient.Get(timeKey).Result()

	givenPoints := 0
	if timeErr == nil {
		fmt.Println("User has changed " + timeValue + " points today.")
		givenPoints, _ = strconv.Atoi(timeValue)
	}

	if givenPoints+numPoints > MAX_POINTS_PER_DAY {
		return "", errors.New("You are attempting to change more than the maximum allotted " + strconv.Itoa(MAX_POINTS_PER_DAY) + " points per day. You have added/removed " + strconv.Itoa(givenPoints) + " today.")
	}

	formattedUserKey := formatPointKey(userID, messageGuildID)
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

	if timeErr == redis.Nil {
		fmt.Println("User: " + userID + " has not given points today. Adding.")
		utils.RedisClient.Set(timeKey, numPoints, 0)
	} else if timeErr != nil {
		return "", timeErr
	} else {
		intPoints, _ := strconv.Atoi(timeValue)
		fmt.Println("Current points given for user: ", timeValue)
		newPoints := intPoints + numPoints
		fmt.Println("New points given for user: ", newPoints)
		utils.RedisClient.Set(timeKey, newPoints, 0)
	}

	return "Points successfully awarded", nil
}

func TakePointsFromUser(userID string, numPoints int, authorID string, messageGuildID string) (message string, err error) {
	timeNow := time.Now().UTC()
	timeString := strconv.Itoa(timeNow.Year()) + "-" + timeNow.Month().String() + "-" + strconv.Itoa(timeNow.Day())
	timeKey := formatTimeKey(timeString, authorID, messageGuildID)

	timeValue, timeErr := utils.RedisClient.Get(timeKey).Result()

	givenPoints := 0
	if timeErr == nil {
		fmt.Println("User has changed " + timeValue + " points today.")
		givenPoints, _ = strconv.Atoi(timeValue)
	}

	if givenPoints+numPoints > MAX_POINTS_PER_DAY {
		return "", errors.New("You are attempting to change more than the maximum allotted " + strconv.Itoa(MAX_POINTS_PER_DAY) + " points per day. You have added/removed " + strconv.Itoa(givenPoints) + " today.")
	}

	formattedUserKey := formatPointKey(userID, messageGuildID)
	value, err := utils.RedisClient.Get(formattedUserKey).Result()

	if err == redis.Nil {
		fmt.Println("Points do not exist for user: " + userID + ". Adding.")
		utils.RedisClient.Set(formattedUserKey, -numPoints, 0)
	} else if err != nil {
		return "", err
	} else {
		intPoints, _ := strconv.Atoi(value)
		fmt.Println("Current point value for user: ", value)
		newPoints := intPoints - numPoints
		fmt.Println("New point value for user: ", newPoints)
		utils.RedisClient.Set(formattedUserKey, newPoints, 0)
	}

	if timeErr == redis.Nil {
		fmt.Println("User: " + userID + " has not changed points today. Adding entry.")
		utils.RedisClient.Set(timeKey, numPoints, 0)
	} else if timeErr != nil {
		return "", timeErr
	} else {
		intPoints, _ := strconv.Atoi(timeValue)
		fmt.Println("Current points changed for user: ", timeValue)
		newPoints := intPoints + numPoints
		fmt.Println("New points changed for user: ", newPoints)
		utils.RedisClient.Set(timeKey, newPoints, 0)
	}

	return "Points successfully removed", nil
}

type LeaderBoardEntry struct {
	Username string
	Points   int
}

type pointList []LeaderBoardEntry

func (p pointList) Len() int           { return len(p) }
func (p pointList) Less(i, j int) bool { return p[i].Points < p[j].Points }
func (p pointList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func GetPointLeaderBoard(s *discordgo.Session, messageGuildID string) (pointList, error) {
	allKeys, _, err := utils.RedisClient.Scan(0, "*"+messageGuildID+"__POINTS", 1000).Result()
	if err != nil {
		return nil, err
	}

	result := make(pointList, len(allKeys))
	for idx, k := range allKeys {
		value, _ := utils.RedisClient.Get(k).Result()
		userid := strings.Split(k, "-")[0]
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
	sort.Sort(sort.Reverse(result))
	return result, nil
}
