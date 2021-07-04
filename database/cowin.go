package database

import (
	"bytes"
	"cowingram/builders"
	"cowingram/constants"
	"cowingram/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type CowinHttpClient struct {
}

func (Client *CowinHttpClient) GetClient() {

}

func CheckAvailability(c *fiber.Ctx) error {

	currentTime := time.Now()
	today := currentTime.Format("02-01-2006")

	var payload models.TelegramPayload
	body := string(c.Body())
	fmt.Println(body)
	err := json.Unmarshal([]byte(body), &payload)

	if err != nil {
		panic(err)
	}

	pin := payload.Message.Text
	if _, err := strconv.Atoi(pin); err != nil && len(pin) != 6 {
		builders.SendMessage(int64(payload.Message.Chat.Id), "Please send a valid pin number.")
		return c.Send([]byte("Please send valid pincode"))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, _ := http.NewRequest("GET", constants.CowinApi+constants.SessionsPath, nil)
	request.Header.Add("accept", "application/json")
	request.Header.Add("Accept-Language", "hi_IN")
	query := request.URL.Query()
	query.Add("pincode", pin)
	query.Add("date", today)
	request.URL.RawQuery = query.Encode()
	res, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response models.CowinResponse

	json.Unmarshal(responseBody, &response)
	var buffer bytes.Buffer
	if len(response.Sessions) > 0 {
		fmt.Fprintf(&buffer, " ---------")
	}
	for i := 0; i < len(response.Sessions); i++ {

		fmt.Fprintf(&buffer, "\n Name: %s\n Address: %s\n Dose 1: %d\n Dose 2: %d\n Vaccine: %s\n Fee: %s\n ---------",
			response.Sessions[i].Name, response.Sessions[i].Address, response.Sessions[i].AvailableCapacityDose1,
			response.Sessions[i].AvailableCapacityDoes2, response.Sessions[i].Vaccine, response.Sessions[i].FeeType)
	}

	if len(response.Sessions) == 0 {
		fmt.Fprintf(&buffer, "There are currently no vaccine slots "+
			"available to us, please check with https://cowin.gov.in site.")
	}

	responseData := buffer.String()
	builders.SendMessage(int64(payload.Message.Chat.Id), responseData)
	c.Send([]byte(responseData))

	defer res.Body.Close()

	return c.SendStatus(res.StatusCode)
}
