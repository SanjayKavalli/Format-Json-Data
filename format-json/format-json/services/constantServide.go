package services

type Iconstantservices interface {
	GetDictionary() map[string]string
}
type Constantservicestruct struct{}

// GetDictionary implements Iconstantservices.

func ConstantserviceCtor() *Constantservicestruct {
	return &Constantservicestruct{}
}

type ResponseBody struct {
	Event            string                 `json:"event"`
	Event_type       string                 `json:"event_type"`
	App_id           string                 `json:"app_id"`
	Uer_id           string                 `json:"user_id"`
	Message_id       string                 `json:"message_id"`
	Page_title       string                 `json:"page_title"`
	Page_url         string                 `json:"page_url"`
	Browser_language string                 `json:"browser_language"`
	Screen_size      string                 `json:"screen_size"`
	Attributes       map[string]interface{} `json:"attributes"`
	Traits           map[string]Trait       `json:"traits"`
}

// Declare Dictionary
var Data = map[string]string{
	"ev":  "event",
	"et":  "event_type",
	"id":  "app_id",
	"uid": "user_id",
	"mid": "message_id",
	"t":   "page_title",
	"p":   "page_url",
	"l":   "browser_language",
	"sc":  "screen_size",
}

func (cs Constantservicestruct) GetDictionary() map[string]string {

	return Data
}
