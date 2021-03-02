package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AjouLibrary 아주대학교 중앙도서관 좌석이용 현황
var AjouLibrary = os.Getenv("AJOU_LIBRARY")

// Library ...
type Library struct {
	Code string `json:"code"`
	Data struct {
		List []struct {
			ID          int64 `json:"id"`
			ActiveTotal int64 `json:"activeTotal"`
			Available   int64 `json:"available"`
			BranchGroup struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			} `json:"branchGroup"`
			DisablePeriod interface{} `json:"disablePeriod"`
			IsActive      bool        `json:"isActive"`
			IsReservable  bool        `json:"isReservable"`
			Name          string      `json:"name"`
			Note          interface{} `json:"note"`
			Occupied      int64       `json:"occupied"`
			RoomTypeID    int64       `json:"roomTypeId"`
			Total         int64       `json:"total"`
		} `json:"list"`
		TotalCount int64 `json:"totalCount"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

//go:noinline
// GetLibraryAvailable ...
func GetLibraryAvailable() (Library, error) {
	var library Library

	// Disable SSL authentication and post jsonValue
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(AjouLibrary)
	if err != nil {
		fmt.Println(err)
		return library, err
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &library)
	if err != nil {
		fmt.Println(err)
		return library, err
	}

	return library, nil
}
