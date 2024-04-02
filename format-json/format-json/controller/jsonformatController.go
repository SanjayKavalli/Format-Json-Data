package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"formatjsondata/services"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonFormatController struct {
	jsonservice     services.IJsonFormatService
	constantservice services.Iconstantservices
}

func JsonControllerCtor(jsonservice services.IJsonFormatService, constantservice services.Iconstantservices) *JsonFormatController {
	return &JsonFormatController{
		jsonservice:     jsonservice,
		constantservice: constantservice,
	}
}

type JSONData map[string]interface{}

func (js *JsonFormatController) JsonFormat(ctx *gin.Context) {

	dataChannel := make(chan JSONData)
	go js.worker(dataChannel)

	fmt.Print("Jsonformat Controller starts")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading json body in JsonController ")
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	var data JSONData
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Println("Error while unmarshaling json body in JsonController ")
		log.Println(err)
		ctx.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if js.jsonservice == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JSON service not initialized"})
		return
	}
	dataChannel <- data
	res, _ := js.sendDataToWorker(data)
	ctx.JSON(http.StatusOK, res)

}
func (js *JsonFormatController) worker(dataChannel <-chan JSONData) {
	for {
		select {
		case data := <-dataChannel:

			sendToWebhook(data)
		}
	}
}

func (js *JsonFormatController) sendDataToWorker(data JSONData) (*services.ResponseBody, error) {
	processedData, _ := js.jsonservice.ProcessJson(data)
	fmt.Printf("Received data: %+v\n", data)
	fmt.Printf("processedData data: %+v\n", processedData)
	response, err := sendToWebhook(processedData)
	if err != nil {
		log.Println("sendDataToWorker err ", err)
	}
	return response, nil
}
func sendToWebhook(data JSONData) (*services.ResponseBody, error) {
	var response *services.ResponseBody
	processedDatabyte, _ := json.Marshal(data)

	if err := json.Unmarshal(processedDatabyte, &response); err != nil {
		fmt.Println("Error sending data to webhook:", err)
		return nil, err
	}
	Databyte, _ := json.Marshal(response)
	req, err := http.NewRequest("POST", "https://webhook.site/9501e80a-4b57-463f-a8de-4f8d3481253d", bytes.NewBuffer(Databyte))
	if err != nil {
		fmt.Println("Error sending data to webhook:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("respose:", response)

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	status := "failed"
	if resp.StatusCode == http.StatusOK {
		status = "delivered"
	}

	log.Println(status)

	if status == "failed" {
		return nil, errors.New(status)
	}

	fmt.Println("Webhook Response Body:", resp.Body)
	fmt.Println("Data sent to webhook successfully")
	fmt.Println(resp)
	defer resp.Body.Close()
	return response, nil

}
