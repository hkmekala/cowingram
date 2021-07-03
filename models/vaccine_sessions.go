package models

type CowinResponse struct {
	Sessions []Session
}

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
