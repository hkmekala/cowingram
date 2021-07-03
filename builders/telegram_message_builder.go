package builders

import (
	"cowingram/constants"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func BuildTelegramPayload() {

}

func SendMessage(ChatId int64, Text string) bool {
	Endpoint := constants.TelegramApi + os.Getenv("TELEGRAM_API_KEY") + constants.TelegramMessagePath
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
