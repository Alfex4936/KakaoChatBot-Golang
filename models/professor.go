package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const AjouPeople = "https://mportal.ajou.ac.kr/system/phone/selectList.ajax"

type People struct {
	MsgCode     string `json:"msgCode"`
	PhoneNumber []struct {
		BussNm    string `json:"bussNm"`           // 업무명: "XXX학과(공학인증)"
		DeptCd    string `json:"deptCd"`           // "DS01234657"
		DeptNm    string `json:"deptNm"`           // 부서명: "정보통신대학교학팀(팔달관 777-1)"
		Email     string `json:"email"`            // 이메일: "example@ajou.ac.kr"
		Name      string `json:"korNm"`            // 이름(신분): "이름1(직원)" | "이름2(교원)"
		MdfLineNo int64  `json:"mdfLineNo,string"` // "289"
		TelNo     string `json:"telNo"`            // 전화번호: 031-219-"1234"
		UserNo    int64  `json:"userNo,string"`    // "201900000"
	} `json:"phoneNumber"`
}

// GetPeople searches a given name in AjouPeople link
func GetPeople(keyword string) (People, error) {
	var people People

	// POST
	jsonValue, _ := json.Marshal(map[string]string{"keyword": keyword})

	// Disable SSL authentication and post jsonValue
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post((AjouPeople), "application/json;charset=UTF-8", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		return people, err
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &people)
	if err != nil {
		fmt.Println(err)
		return people, err
	}
	// fmt.Println(people.PhoneNumber[0].Name)

	return people, nil
}
