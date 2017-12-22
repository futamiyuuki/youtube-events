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

type CategoryModel struct {
	Category string `bson:"category"`
	Count    int    `bson:"count"`
}

type ChannelModel struct {
	ID         string          `bson:"_id"`
	Categories []CategoryModel `bson:"categories"`
}
