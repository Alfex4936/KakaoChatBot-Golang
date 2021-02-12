package models

// SimpleText for Kakao Response
type SimpleText struct {
	Template struct {
		Outputs struct {
			SimpleText struct {
				Text string `json:"text"`
			} `json:"simpleText"`
		} `json:"outputs"`
	} `json:"template"`
	Version string `json:"version"`
}

// KakaoJSON request main
type KakaoJSON struct {
	Action struct {
		ID          string `json:"id"`
		ClientExtra struct {
		} `json:"clientExtra"`
		DetailParams map[string]interface{} `json:"detailParams"`
		Name         string                 `json:"name"`
		Params       map[string]interface{} `json:"params"`
	} `json:"action"`
	Bot struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"bot"`
	Contexts []interface{} `json:"contexts"`
	Intent   struct {
		ID    string `json:"id"`
		Extra struct {
			Reason struct {
				Code    int64  `json:"code"`
				Message string `json:"message"`
			} `json:"reason"`
		} `json:"extra"`
		Name string `json:"name"`
	} `json:"intent"`
	UserRequest struct {
		Block struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"block"`
		Lang   string `json:"lang"`
		Params struct {
			IgnoreMe bool   `json:"ignoreMe,string"`
			Surface  string `json:"surface"`
		} `json:"params"`
		Timezone string `json:"timezone"`
		User     struct {
			ID         string `json:"id"`
			Properties struct {
				BotUserKey  string `json:"botUserKey"`
				BotUserKey2 string `json:"bot_user_key"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"user"`
		Utterance string `json:"utterance"`
	} `json:"userRequest"`
}

// BuildSimpleText ...
func BuildSimpleText(msg string) SimpleText {
	stext := SimpleText{Version: "2.0"}
	stext.Template.Outputs.SimpleText.Text = msg
	return stext
}
