package giteaapiclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	configuration "xqledger/gitoperator/configuration"
	utils "xqledger/gitoperator/utils"
	"github.com/google/uuid"
	"gopkg.in/resty.v1"
	kafka "github.com/segmentio/kafka-go"
)

// https://try.gitea.io/api/swagger#/

const componentMessage = "Gitea API Client"
const apiFilePath = "%s/repos/%s/%s/contents/%s"
var config = configuration.GlobalConfiguration
var client *resty.Client
var kafkaWriter *kafka.Writer



func getAPIClient() *resty.Client{
	if client != nil {
		return client
	}
	client = resty.New()
	client.SetTimeout(time.Duration(config.Gitserver.Strategy.Timeout) * time.Millisecond)
	client.SetHeaders(map[string]string{
        "Content-Type": "application/json",
        "User-Agent": "GitOperator",
		"Authorization": config.Gitserver.Authtoken,
    })
	return client
}

func getFile(filename string) (ContentsResponse, error) {
	methodMessage := "getFile"
	var contentResponse ContentsResponse
	c := getAPIClient()
	resp, err := c.R().
		Get(fmt.Sprintf(apiFilePath, config.Gitserver.Url, config.Gitserver.Username, filename + ".json"))
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, utils.Record_new_git_written_fail)
		return contentResponse, err
	}
	if resp.StatusCode() != 200 {
		err := errors.New(fmt.Sprintf("Unexpected status code, expected %d, got %d instead", 200, resp.StatusCode()))
		utils.PrintLogError(err, componentMessage, methodMessage, "Unexpected response")
		return contentResponse, err
	}
	
	unmarshalErr := json.Unmarshal(resp.Body(), &contentResponse)

	if unmarshalErr != nil {
		utils.PrintLogError(unmarshalErr, componentMessage, methodMessage, "Response body not processable")
		return contentResponse,unmarshalErr
	}
	return contentResponse, nil
}

func getNameFromEmail(email string) (string, error) {
	at := strings.LastIndex(email, "@")
    if at >= 0 {
        username, _ := email[:at], email[at+1:]
		return username, nil
    } else {
		err := errors.New(fmt.Sprintf("Error: %s is an invalid email address\n", email))
		return "", err
	}
}

func getValidContentFromEventContent(content string) (string, error) {
	var prettyJSON bytes.Buffer
	jsonErr := json.Indent(&prettyJSON, []byte(content), "", "\t")
	if jsonErr != nil {
		return "", jsonErr
	}
	return string(prettyJSON.Bytes()), nil
}

func CreateFileInRepo(event *utils.RecordEvent) error{
	methodMessage := "CreateFileInRepo" 
	var payload CreateFileOptions
	var identity Identity
	username, emailErr := getNameFromEmail(event.User)
	if emailErr != nil {
		identity.Name = event.User
	} else {
		identity.Name = username
	}
	var session = "master"
	if len(event.Session) > 0 {
		session = event.Session
		// optional: check if the branch exists
	}
	identity.Email = event.User
	payload.Author = identity
	payload.Branch = session
	payload.Committer = identity
	payload.Message = event.Id
	payload.New_branch = ""
	payload.Signoff =  true

	validContent, validJsonErr := getValidContentFromEventContent(event.RecordContent)
	if validJsonErr != nil {
		utils.PrintLogError(validJsonErr, componentMessage, methodMessage, "Error getting valid JSON record from content for file creation")
		return validJsonErr
	}
	payload.Content = base64.StdEncoding.EncodeToString([]byte(validContent))

	c := getAPIClient()
	_, err := c.R().
		SetBody(payload).
		Post(fmt.Sprintf(apiFilePath, config.Gitserver.Url, config.Gitserver.Username, event.DBName, event.Id + ".json"))
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, utils.Record_new_git_written_fail)
		return err
	}
	go utils.PrintLogInfo(componentMessage, methodMessage, utils.Record_new_git_written_ok)
	go SendMessageToTopic(event)
	return nil
}

func UpdateFileInRepo(event *utils.RecordEvent) error{
	methodMessage := "UpdateFileInRepo" 
	var payload UpdateFileOptions
	var identity Identity
	username, emailErr := getNameFromEmail(event.User)
	if emailErr != nil {
		identity.Name = event.User
	} else {
		identity.Name = username
	}
	var session = "master"
	if len(event.Session) > 0 {
		session = event.Session
	}
	identity.Email = event.User
	payload.Author = identity
	payload.Branch = session
	payload.Committer = identity
	payload.Message = event.Id
	payload.New_branch = ""
	payload.Signoff =  true

	currentFile, shaErr := getFile(event.Id)
	if shaErr != nil {
		return shaErr
	}
	payload.Sha = currentFile.Sha

	var commitDate CommitDateOptions
	commitDate.Author = event.User
	commitDate.Committer = event.User
	payload.Dates = commitDate
	payload.From_path = ""

	validContent, validJsonErr := getValidContentFromEventContent(event.RecordContent)
	if validJsonErr != nil {
		utils.PrintLogError(validJsonErr, componentMessage, methodMessage, "Error getting valid JSON record from content for file update")
		return validJsonErr
	}
	payload.Content = base64.StdEncoding.EncodeToString([]byte(validContent))

	c := getAPIClient()
	_, err := c.R().
		SetBody(payload).
		Put(fmt.Sprintf(apiFilePath, config.Gitserver.Url, config.Gitserver.Username, event.DBName, event.Id + ".json"))
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, utils.Record_update_git_written_fail)
		return err
	}
	go utils.PrintLogInfo(componentMessage, methodMessage, utils.Record_update_git_written_ok)
	go SendMessageToTopic(event)
	return nil
}

func DeleteFileInRepo(event *utils.RecordEvent) error{
	methodMessage := "DeleteFileInRepo" 

	var payload DeleteFileOptions
	var identity Identity
	username, emailErr := getNameFromEmail(event.User)
	if emailErr != nil {
		identity.Name = event.User
	} else {
		identity.Name = username
	}
	var session = "master"
	if len(event.Session) > 0 {
		session = event.Session
	}
	identity.Email = event.User
	payload.Author = identity
	payload.Branch = session
	payload.Committer = identity
	payload.Message = event.Id
	payload.New_branch = ""
	payload.Signoff =  true

	c := getAPIClient()
	_, err := c.R().
		SetBody(payload).
		Delete(fmt.Sprintf(apiFilePath, config.Gitserver.Url, config.Gitserver.Username, event.DBName, event.Id + ".json"))
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, utils.Record_delete_git_written_fail)
		return err
	}
	go utils.PrintLogInfo(componentMessage, methodMessage, utils.Record_delete_git_written_ok)
	go SendMessageToTopic(event)
	return nil
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	if kafkaWriter != nil {
		return kafkaWriter
	} else {
		return kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{kafkaURL},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
	}
}

func SendMessageToTopic(event *utils.RecordEvent) {
	methodMessage := "SendMessageToTopic"

	eventAsJSON, err := json.Marshal(event)
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, "Event cannot be marshaled properly after written to Git")
	}

	topic := config.Kafka.Gitactionbacktopic
	broker := config.Kafka.Bootstrapserver
	kafkaWriter := getKafkaWriter(broker, topic)
	utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Message sent to topic '%s'", topic))
	//defer kafkaWriter.Close()

	topicContent := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: []byte(eventAsJSON),

	}

	writeErr := kafkaWriter.WriteMessages(context.Background(), topicContent)
	if writeErr != nil {
		utils.PrintLogError(writeErr, componentMessage, methodMessage, fmt.Sprintf("Error writing message to topic '%s'", topic))
	}
	utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Message sent to topic '%s' successfully", topic))
}


type Identity struct {
	Email   string `json:"email"`
	Name   string `json:"name"`
}

type CreateFileOptions struct {
	Author   Identity `json:"author"` // author and committer are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used)
	Branch   string `json:"branch"` // opt
	Committer Identity `json:"committer"`
	Content string `json:"content"` // base64 encoded
	Message string `json:"message"` // opt
	New_branch string `json:"new_branch"` // opt
	Signoff bool `json:"signoff"`
}

type CommitDateOptions struct {
	Author   string `json:"author"`
	Committer   string `json:"committer"`
}

type UpdateFileOptions struct {
	Author   Identity `json:"author"` // author and committer are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used)
	Branch   string `json:"branch"` // opt
	Committer Identity `json:"committer"`
	Content string `json:"content"` // base64 encoded
	Dates CommitDateOptions `json:"dates"` 
	From_path string `json:"from_path"` 
	Message string `json:"message"` // opt
	New_branch string `json:"new_branch"` // opt
	Sha string `json:"sha"` // sha is the SHA for the file that already exists
	Signoff bool `json:"signoff"`
}

type DeleteFileOptions struct {
	Author   Identity `json:"author"` // author and committer are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used)
	Branch   string `json:"branch"` // opt
	Committer Identity `json:"committer"`
	Dates CommitDateOptions `json:"dates"` 
	Message string `json:"message"` // opt
	New_branch string `json:"new_branch"` // opt
	Sha string `json:"sha"` // sha is the SHA for the file that already exists
	Signoff bool `json:"signoff"`
}

type ContentsResponse struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sha   string `json:"sha"`
	Type string `json:"type"`
	Encoding   string `json:"encoding"`
	Content   string `json:"content"`
	Url   string `json:"url"`
	Git_url   string `json:"git_url"`
	Download_url   string `json:"download_url"`
}
