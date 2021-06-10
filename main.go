package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CowinApi            = "https://cdn-api.co-vin.in/api/v2/"
	SessionsPath        = "appointment/sessions/public/findByPin"
	TelegramApi         = "https://api.telegram.org/bot"
	TelegramMessagePath = "/sendMessage"
)

type Session struct {
	SessionId              string   `json:"session_id"`
	CenterID               string   `json:"center_id"`
	Name                   string   `json:"name"`
	Address                string   `json:"address"`
	StateName              string   `json:"state_name"`
	DistrictName           string   `json:"district_name"`
	Pincode                int32    `json:"pincode"`
	From                   string   `json:"from"`
	To                     string   `json:"to"`
	Lat                    string   `json:"lat"`
	Long                   string   `json:"long"`
	FeeType                string   `json:"fee_type"`
	Date                   string   `json:"date"`
	AvailableCapacityDose1 uint64   `json:"available_capacity_dose1"`
	AvailableCapacityDoes2 uint64   `json:"available_capacity_does_2"`
	Fee                    string   `json:"fee"`
	MinAgeLimit            uint16   `json:"min_age_limit"`
	Vaccine                string   `json:"vaccine"`
	Slots                  []string `json:"slots"`
}

type CowinResponse struct {
	Sessions []Session
}

type TelegramPayload struct {
	UpdateId uint64       `json:"update_id"`
	Message  MessageField `json:"message"`
}

type MessageField struct {
	MessageId uint64    `json:"message_id"`
	From      FromField `json:"from"`
	Chat      ChatField `json:"chat"`
	Date      uint64    `json:"date"`
	Text      string    `json:"text"`
}

type FromField struct {
	Id           uint64 `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type ChatField struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

func sendMessage(ChatId int64, Text string) bool {
	Endpoint := TelegramApi + os.Getenv("TELEGRAM_API_KEY") + TelegramMessagePath
	client := &http.Client{Timeout: 10 * time.Second}
	request, _ := http.NewRequest("GET", Endpoint, nil)
	query := request.URL.Query()
	query.Add("chat_id", strconv.FormatInt(ChatId, 10))
	query.Add("text", Text)
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)

	fmt.Println(response)

	if err != nil {
		panic(err)
		return false
	}

	return true
}

func CheckAvailability(c *fiber.Ctx) error {

	currentTime := time.Now()
	today := currentTime.Format("02-01-2006")

	var payload TelegramPayload
	body := string(c.Body())
	fmt.Println(body)
	err := json.Unmarshal([]byte(body), &payload)

	if err != nil {
		panic(err)
	}

	pin := payload.Message.Text
	if _, err := strconv.Atoi(pin); err != nil && len(pin) != 6 {
		sendMessage(int64(payload.Message.Chat.Id), "Please send a valid pin number.")
		return c.Send([]byte("Please send valid pincode"))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, _ := http.NewRequest("GET", CowinApi+SessionsPath, nil)
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

	var response CowinResponse

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

	responseData := buffer.String()
	sendMessage(int64(payload.Message.Chat.Id), responseData)
	c.Send([]byte(responseData))

	defer res.Body.Close()

	return nil
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/pin", CheckAvailability)

	if strings.Compare(os.Getenv("PRODUCTION"), "true") == 0 {
		app.Listen(":8080")
	} else {
		app.Listen(":3000")
	}
}
