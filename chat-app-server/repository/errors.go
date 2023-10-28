package repository

import "errors"

/* DATABASE */
var ErrorJSON error = errors.New("internal json error")
var ErrorFailedToSet error = errors.New("failed to set/add to a key in database")
var ErrorFailedToGet error = errors.New("failed to get key from database")

/* HTTP SERVER */
var ErrorFailedToDecode error = errors.New("failed to decode bad JSON object")
var ErrorUsernameDoesntExist error = errors.New("one of the usernames doesnt exist")
var ErrorUsernameNotUnique error = errors.New("username is not unique")
var ErrorInvalidCredentials error = errors.New("invalid credentials")
var ErrorChannelDoesntExist error = errors.New("channel doesnt exist")
