package main

type event struct {
	EventType     string `json:"event_type"`
	VideoID       string `json:"video_id"`
	VideoCategory string `json:"video_category"`
	ChannelID     string `json:"channel_id"`
	IsSubscribed  bool   `json:"is_subscirbed"`
	SearchTerm    string `json:"search_term"`
}

type categoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type channel struct {
	ChannelID      string          `json:"_id"`
	CategoryCounts []categoryCount `json:"categories"`
}

type eventModel struct {
	EventType     string `bson:"event_type"`
	VideoID       string `bson:"video_id"`
	VideoCategory string `bson:"video_category"`
	ChannelID     string `bson:"channel_id"`
	IsSubscribed  bool   `bson:"is_subscirbed"`
	SearchTerm    string `bson:"search_term"`
}

type categoryModel struct {
	Category string `bson:"category"`
	Count    int    `bson:"count"`
}

type channelModel struct {
	ID         string          `bson:"_id"`
	Categories []categoryModel `bson:"categories"`
}
