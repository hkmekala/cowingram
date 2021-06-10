package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	CowinApi     = "https://cdn-api.co-vin.in/api/v2/"
	SessionsPath = "appointment/sessions/public/findByPin"
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

func CheckAvailability(c *fiber.Ctx) error {

	currentTime := time.Now()
	today := currentTime.Format("02-01-2006")
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)
	pin := body["pin"].(string)
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

	c.Send([]byte(response.Sessions[0].Address))

	defer res.Body.Close()

	return nil
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/pin", CheckAvailability)

	app.Listen(":3000")
}
