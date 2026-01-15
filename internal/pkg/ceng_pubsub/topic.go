package ceng_pubsub

/*
PubSubTopic represents a topic name where events are published and
consumed by different modules. Each topic must contain only events
related to a specific entity domain.
*/
type PubSubTopic string

/*
List of available topics.
*/
const (
	TopicMenuCategoryV1 PubSubTopic = "topic/v1/menu-category"
	TopicMenuItemV1     PubSubTopic = "topic/v1/menu-item"
	TopicMenuOptionV1   PubSubTopic = "topic/v1/menu-option"
	TopicTableV1        PubSubTopic = "topic/v1/table"
)
