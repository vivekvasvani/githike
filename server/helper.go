package server

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
)

type Select struct {
	Actions []struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		SelectedOptions []struct {
			Value string `json:"value"`
		} `json:"selected_options"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		Text        string `json:"text"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Text       string `json:"text"`
				Type       string `json:"type"`
				DataSource string `json:"data_source"`
				Options    []struct {
					Text  string `json:"text"`
					Value string `json:"value"`
				} `json:"options"`
			} `json:"actions"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
}

type Button struct {
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		Text        string `json:"text"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Text       string `json:"text"`
				Type       string `json:"type"`
				DataSource string `json:"data_source"`
				Options    []struct {
					Text  string `json:"text"`
					Value string `json:"value"`
				} `json:"options"`
			} `json:"actions"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
}

func GetPayload(payloadPath string) string {
	if payloadPath != "" {
		dir, _ := os.Getwd()
		templateData, _ := ioutil.ReadFile(dir + "/payloads/" + payloadPath)
		return string(templateData)
	} else {
		return ""
	}
}

func SubstParams(sessionMap []string, textData string) string {
	for i, value := range sessionMap {
		if strings.ContainsAny(textData, "${"+strconv.Itoa(i)) {
			textData = strings.Replace(textData, "${"+strconv.Itoa(i)+"}", value, -1)
		}
	}
	return textData
}

func SetErrorResponse(ctx *fasthttp.RequestCtx, statusCode, statusType, statusMessage string, httpStatus int) {
	log.Println(statusCode, statusType, statusMessage)
	var response Response
	response.Status.StatusCode = statusCode
	response.Status.StatusType = statusType
	response.Status.Message = statusMessage
	glog.Infoln("Error Reponse " + ToJsonString(response))
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(response))
	ctx.SetStatusCode(httpStatus)
}

func SetSuccessResponse(ctx *fasthttp.RequestCtx, statusCode, statusType, statusMessage string, httpStatus int, data interface{}) {
	var response Response
	response.Status.StatusCode = statusCode
	response.Status.StatusType = statusType
	response.Status.Message = statusMessage
	response.Data = data
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(response))
	glog.Infoln("Success Reponse " + ToJsonString(response))
	ctx.SetStatusCode(httpStatus)
}
