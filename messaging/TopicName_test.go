package messaging

import (
	"github.com/wbrush/go-common/config"
	"testing"
)

func TestValidate(t *testing.T) {
	testValue1 := PubSubCmdTopic
	if !testValue1.Validate() {
		t.Error("Validate failed for valid Topic!")
	}
	testValue2 := "testTopicName"
	if PubSubTopic(testValue2).Validate() {
		t.Error("Validate failed for invalid Topic!")
	}
}

func TestGetTopicName(t *testing.T) {
	testValue1 := PubSubCmdTopic
	topic1,err1 := testValue1.GetTopicName(string(config.EnvironmentTypeLocal))
	if err1 != nil || topic1 != string(PubSubCmdTopic)+topicSeparator+string(config.EnvironmentTypeLocal) {
		t.Error("Invalid topic name returned for local env!")
	}

	topic2,err2 := testValue1.GetTopicName("testEnv")
	if err2 == nil || topic2 != ""{
		t.Error("Should have gotten an error about Invalid Environment Name and topic name should be empty!")
	}

}

func TestGetTopicFromName(t *testing.T) {
	testTopicName,_ := PubSubCmdTopic.GetTopicName(string(config.EnvironmentTypeLocal))
	var topic1 PubSubTopic
	topic1.GetTopicFromName(testTopicName)
	if topic1 != PubSubCmdTopic {
		t.Error("Got invalid topic returned from topic name")
	}
}
