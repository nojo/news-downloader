package main

import "github.com/mediocregopher/radix.v2/redis"

var client *redis.Client

// a better design might be to use a hash to store both the filename and the content.  Would be easier to keep
// in sync that way.  There's a small chance of an error when one or the other has been set.  That's hard
// to recover from since there's no correlation between the file names and the content.
const filenameSetName = "FILENAMES_XML"
const xmlListName = "NEWS_XML"

// InitRedisStorage creates a connection to Redis
func InitRedisStorage(hostNamePort string) (err error) {
	client, err = redis.Dial("tcp", hostNamePort)
	return err
}

// IsFileUploaded checks whether a file has already been uploaded. We keep a separate set of filenames of the values in the array so lookup is quick.
func IsFileUploaded(fileName string) (bool, error) {
	if exists, err := client.Cmd("SISMEMBER", filenameSetName, fileName).Int(); exists == 1 {
		return true, err
	} else {
		return false, err
	}
}

// AddFileToList adds a file to the list of files and to the tracking set of filenames
func AddFileToList(fileName string, fileContents string) error {
	if err := client.Cmd("RPUSH", xmlListName, fileContents).Err; err != nil {
		return err
	}
	if err := client.Cmd("SADD", filenameSetName, fileName).Err; err != nil {
		return err
	}
	return nil
}
