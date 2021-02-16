package models

type H map[string]interface{}

// Make a basicCard
// template := gin.H{"outputs": []gin.H{{
// 	"basicCard": gin.H{
// 		"buttons":   []gin.H{{"action": "webLink", "label": "날씨 홈페이지 열기", "webLinkUrl": models.NaverWeather}},
// 		"thumbnail": gin.H{"imageUrl": ""},
// 		"title":     fmt.Sprintf("현재 수원 영통구 날씨: %s", weather.CurrentTemp),
// 		"description": fmt.Sprintf("(해)<br>현재 %s<br>최저, 최고 온도: %s, %s<br>낮, 밤 강수 확률: %s, %s<br>미세먼지: %s<br>초미세먼지: %s<br>자외선: %s",
// 			weather.CurrentStatus,
// 			weather.MinTemp, weather.MaxTemp,
// 			weather.RainDay, weather.RainNight,
// 			weather.FineDust, weather.UltraDust, weather.UV),
// 	},
// }}}
// basicCard := gin.H{"version": "2.0", "template": template}

// SimpleText for Kakao Response
type SimpleText struct {
	Template struct {
		Outputs []struct {
			SimpleText Text `json:"simpleText"`
		} `json:"outputs"`
	} `json:"template"`
	Version string `json:"version"`
}

// Text for SimpleText
type Text struct {
	Text string `json:"text"`
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
func BuildSimpleText(msg string) *SimpleText {
	stext := &SimpleText{Version: "2.0"}

	var temp []struct {
		SimpleText Text `json:"simpleText"`
	}
	simpleText := Text{Text: msg}

	text := struct {
		SimpleText Text `json:"simpleText"`
	}{SimpleText: simpleText}

	temp = append(temp, text)

	stext.Template.Outputs = temp
	return stext
}

func BuildQuickReply(msg, label string) H {
	return H{"messageText": msg, "aciton": "message", "label": label}
}
