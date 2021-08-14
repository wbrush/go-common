package messaging

import (
	"strings"
	"testing"
)

const (
	testProvider = "UnitTest"
	testMedia    = "disk"
	testAction   = "action1"
	testProperty = "myProperty"
	testId       = "id_unknown"
	testDefSender       = "tbi"
	testSender   =  "something"

)

func TestFilterType_New(t *testing.T) {
	record := ChangeMessage{}
	if len(string(record.Filter)) != 0 {
		t.Error("Invalid Filter for new record!")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	if len(string(record.Filter)) == 0 {
		t.Error("Invalid Filter Setting")
	} else if !strings.HasPrefix(string(record.Filter), testProvider) {
		t.Error("Invalid Provider in Filter")
	} else if !strings.HasSuffix(string(record.Filter), testDefSender) {
		t.Error("Invalid Sender in Filter")
	}
}

func TestFilterType_NewV2(t *testing.T) {
	record := ChangeMessage{}
	if len(string(record.Filter)) != 0 {
		t.Error("Invalid Filter for new record!")
	}

	record.Filter.NewV2(testProvider, testMedia, testAction, testProperty, testId, testSender)
	if len(string(record.Filter)) == 0 {
		t.Error("Invalid Filter Setting")
	} else if !strings.HasPrefix(string(record.Filter), testProvider) {
		t.Error("Invalid Provider in Filter")
	} else if !strings.HasSuffix(string(record.Filter), testSender) {
		t.Error("Invalid Sender in Filter")
	}
}

func TestFilterType_GetProvider(t *testing.T) {
	record := ChangeMessage{}
	provider, err := record.Filter.GetProvider()
	if err == nil {
		t.Error("Didn't Get Error Reading Provider From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	provider, err = record.Filter.GetProvider()
	if err != nil {
		t.Error("Error Reading Provider From Filter: " + err.Error())
	} else if len(string(provider)) == 0 {
		t.Error("Invalid Provider Read from Filter")
	} else if provider != testProvider {
		t.Error("Invalid Provider in Filter")
	}
}

func TestFilterType_GetMedia(t *testing.T) {
	record := ChangeMessage{}
	media, err := record.Filter.GetMedia()
	if err == nil {
		t.Error("Didn't Get Error Reading Media From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	media, err = record.Filter.GetMedia()
	if err != nil {
		t.Error("Error Reading Media From Filter: " + err.Error())
	} else if len(string(media)) == 0 {
		t.Error("Invalid Media Read from Filter")
	} else if media != testMedia {
		t.Error("Invalid Media in Filter")
	}
}

func TestFilterType_GetAction(t *testing.T) {
	record := ChangeMessage{}
	action, err := record.Filter.GetAction()
	if err == nil {
		t.Error("Didn't Get Error Reading Action From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	action, err = record.Filter.GetAction()
	if err != nil {
		t.Error("Error Reading Action From Filter: " + err.Error())
	} else if len(string(action)) == 0 {
		t.Error("Invalid Action Read from Filter")
	} else if action != testAction {
		t.Error("Invalid Action in Filter")
	}
}

func TestFilterType_GetProperty(t *testing.T) {
	record := ChangeMessage{}
	property, err := record.Filter.GetProperty()
	if err == nil {
		t.Error("Didn't Get Error Reading ID From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	property, err = record.Filter.GetProperty()
	if err != nil {
		t.Error("Error Reading Property From Filter: " + err.Error())
	} else if len(string(property)) == 0 {
		t.Error("Invalid Property Read from Filter")
	} else if property != testProperty {
		t.Error("Invalid Property in Filter")
	}
}

func TestFilterType_GetId(t *testing.T) {
	record := ChangeMessage{}
	id, err := record.Filter.GetId()
	if err == nil {
		t.Error("Didn't Get Error Reading ID From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	id, err = record.Filter.GetId()
	if err != nil {
		t.Error("Error Reading ID From Filter: " + err.Error())
	} else if len(string(id)) == 0 {
		t.Error("Invalid ID Read from Filter")
	} else if id != testId {
		t.Error("Invalid ID in Filter")
	}
}

func TestFilterType_GetSender(t *testing.T) {
	record := ChangeMessage{}
	sender, err := record.Filter.GetSender()
	if err == nil {
		t.Error("Didn't Get Error Reading Sender From Empty Filter")
	}

	record.Filter.New(testProvider, testMedia, testAction, testProperty, testId)
	sender, err = record.Filter.GetSender()
	if err != nil {
		t.Error("Error Reading Sender From Filter: " + err.Error())
	} else if len(string(sender)) == 0 {
		t.Error("Invalid Sender Read from Filter")
	} else if sender != testDefSender {
		t.Error("Invalid Sender in Filter")
	}

	record.Filter.NewV2(testProvider, testMedia, testAction, testProperty, testId, testSender)
	sender, err = record.Filter.GetSender()
	if err != nil {
		t.Error("Error Reading Sender From Filter: " + err.Error())
	} else if len(string(sender)) == 0 {
		t.Error("Invalid Sender Read from Filter")
	} else if sender != testSender {
		t.Error("Invalid Sender in Filter")
	}
}
