package database

/**
TODO: Overhaul how messages are stored to facilitate searching by timestamp
Create channel structure (replace contacts with it)
*/

import (
	"chat-app/repository"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisGetRequest[T interface{}](rdb redisDatabaseProvider, key string, res T) (T, error) {
	response, err := rdb.client.Get(context.Background(), key).Result()

	if err != nil {
		return res, repository.ErrorFailedToGet
	}

	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return res, repository.ErrorJSON
	}

	return res, nil
}

func RedisExistsRequest(rdb redisDatabaseProvider, key string) (bool, error) {
	response, err := rdb.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, repository.ErrorFailedToGet
	}

	return response == 1, nil
}

/* USERS */

func (rdb redisDatabaseProvider) CreateUser(username string, password string) (repository.User, error) {
	user := repository.User{
		Timestamp: time.Now().UnixMilli(),
		Username:  username,
		Password:  password,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return user, repository.ErrorJSON
	}
	err = rdb.client.Set(context.Background(), generateUsernameKey(username), string(userJSON), 0).Err()

	if err != nil {
		return user, repository.ErrorFailedToSet
	}

	return user, nil
}

func (rdb redisDatabaseProvider) GetUser(username string) (repository.User, error) {
	return RedisGetRequest[repository.User](rdb, generateUsernameKey(username), repository.User{})
}

func (rdb redisDatabaseProvider) UserExists(username string) (bool, error) {
	return RedisExistsRequest(rdb, generateUsernameKey(username))
}

func (rdb redisDatabaseProvider) CheckPassword(username string, inputtedPassword string) (bool, error) {
	//Get the password for the username
	//Get(key)
	user, err := rdb.GetUser(username)

	if err != nil {
		return false, repository.ErrorFailedToGet
	}

	return user.Password == inputtedPassword, nil
}

/* CHANNELS */

func (rdb redisDatabaseProvider) CreateChannel(members []string) (repository.Channel, error) {
	channel := repository.Channel{
		ID:      rdb.snowflakeProvider.GenerateSnowflake(),
		Members: members,
	}

	channelJSON, err := json.Marshal(channel)
	if err != nil {
		return channel, repository.ErrorJSON
	}

	err = rdb.client.Set(context.Background(), generateChannelKey(channel.ID), channelJSON, 0).Err()
	if err != nil {
		return channel, repository.ErrorFailedToSet
	}

	timestamp := (float64)(time.Now().UnixMilli())

	for _, member := range members {
		err = rdb.UpdateChannels(member, channel.ID, timestamp)
		if err != nil {
			return channel, repository.ErrorFailedToSet
		}
	}

	return channel, nil
}

func (rdb redisDatabaseProvider) GetChannel(channelId string) (repository.Channel, error) {
	return RedisGetRequest[repository.Channel](rdb, generateChannelKey(channelId), repository.Channel{})
}

func (rdb redisDatabaseProvider) ChannelExists(channelId string) (bool, error) {
	return RedisExistsRequest(rdb, generateChannelKey(channelId))
}

/* CHANNELS */

func (rdb redisDatabaseProvider) GetChannels(username string, start, stop int) ([]repository.Channel, error) {
	zRangeArgs := redis.ZRangeArgs{
		Key:   generateChannelsKey(username),
		Start: start,
		Stop:  stop,
		Rev:   true,
	}

	response, err := rdb.client.ZRangeArgsWithScores(context.Background(), zRangeArgs).Result()
	if err != nil {
		return nil, repository.ErrorFailedToGet
	}

	channels := make([]repository.Channel, len(response))

	for i, result := range response {
		channel, err := rdb.GetChannel(result.Member.(string))
		if err != nil {
			return nil, repository.ErrorFailedToGet
		}
		lastMessage, err := rdb.SearchMessages(channel.ID, -1, -1)
		if err != nil {
			return nil, repository.ErrorFailedToGet
		}
		if len(lastMessage) > 0 {
			channel.LastActivity = lastMessage[0]
		}
		channels[i] = channel
	}

	return channels, nil
}

func (rdb redisDatabaseProvider) UpdateChannels(username, channelID string, timestamp float64) error {
	zEntry := redis.Z{Score: timestamp, Member: channelID}

	//ZAdd(setKey, ZEntry)
	err := rdb.client.ZAdd(context.Background(), generateChannelsKey(username), zEntry).Err()

	if err != nil {
		return repository.ErrorFailedToSet
	}

	return nil
}

/* MESSAGES */

func (rdb redisDatabaseProvider) CreateMessage(from string, channelID string, content string) (string, error) {
	message := repository.Message{
		ID:        rdb.snowflakeProvider.GenerateSnowflake(),
		ChannelID: channelID,
		From:      from,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return "", repository.ErrorJSON
	}

	//Set the JSON in the database
	//JSON.SET(key, path ($ is root), json)
	// err = rdb.client.Do(context.Background(), "JSON.SET", generateMessageKey(message.ID), "$", string(messageJSON)).Err()
	err = rdb.client.Set(context.Background(), generateMessageKey(message.ID), string(messageJSON), 0).Err()
	if err != nil {
		return "", repository.ErrorFailedToSet
	}

	channel, err := rdb.GetChannel(channelID)
	if err != nil {
		return "", err
	}

	//Update contact list activity of all users
	for _, member := range channel.Members {
		err = rdb.UpdateChannels(member, channelID, float64(message.Timestamp))
		if err != nil {
			return "", err
		}
	}

	//Update chat sorted set
	zEntry := redis.Z{
		Score:  float64(message.Timestamp),
		Member: message.ID,
	}

	err = rdb.client.ZAdd(context.Background(), generateChannelListKey(channelID), zEntry).Err()
	if err != nil {
		return "", err
	}

	return message.ID, nil
}

func (rdb redisDatabaseProvider) GetMessage(messageID string) (repository.Message, error) {
	return RedisGetRequest[repository.Message](rdb, generateMessageKey(messageID), repository.Message{})
}

func (rdb redisDatabaseProvider) SearchMessages(channelID string, start, stop int) ([]repository.Message, error) {
	zRangeArgs := redis.ZRangeArgs{
		Key:   generateChannelListKey(channelID),
		Start: start,
		Stop:  stop,
	}

	response, err := rdb.client.ZRangeArgs(context.Background(), zRangeArgs).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]repository.Message, len(response))
	for i, messageID := range response {
		message, err := rdb.GetMessage(messageID)
		if err != nil {
			return nil, err
		}
		messages[i] = message
	}

	return messages, nil
}
