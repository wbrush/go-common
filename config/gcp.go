package config

import "strings"

type GCP struct {
	GoogleApplicationCredentials string `json:"google_application_credentials"`
	ProjectID                    string `json:"gcp_project_id"`
	ServicePubTopic              string `json:"service_pub_topic"`
	SubscriptionTopics           string `json:"subscription_topics"` //delimiter ","
}

func (dp GCP) GetSubscriptionTopicsList() []string {
	return strings.Split(dp.SubscriptionTopics, ",")
}

func (dp GCP) Validate() error {
	return nil
}
