package models

import (
	"testing"
)

// go test -v ./models
func TestGetLibrarySeats(t *testing.T) {
	library, err := GetLibraryAvailable()
	if err != nil || !library.Success {
		t.Errorf("Check your HTML connection. %v", library.Success)
		return
	}

	for _, lib := range library.Data.List {
		t.Logf("%v %v", lib.Name, lib.IsActive)
		t.Logf("\n잔여 좌석: %v\n사용 좌석: %v", lib.Available, lib.Occupied)
	}
}
