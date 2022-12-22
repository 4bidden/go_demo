package db

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	DataVendorID           int
	VendorFeedID           int
	VerticalID             int
	ShortCode              string `gorm:"type:varchar(255)"`
	OptInID                int
	CampaignID             int
	MessageID              int
	ScheduledTime          time.Time
	StateCode              string `gorm:"type:varchar(255)"`
	PostalCode             string `gorm:"type:varchar(255)"`
	Timezone               string `gorm:"type:varchar(255)"`
	CarrierID              int
	PhoneNumber            string `gorm:"type:varchar(255)"`
	MessageUID             int
	AggregatorID           int
	MessageSample          string `gorm:"type:text"`
	SampleMessageVersionID int
}

// Generate a random Message object
func (m *Message) RandomMessage() *Message {
	message := &Message{
		DataVendorID:           rand.Int(),
		VendorFeedID:           rand.Int(),
		VerticalID:             rand.Int(),
		ShortCode:              randString(10),
		OptInID:                rand.Int(),
		CampaignID:             rand.Int(),
		MessageID:              rand.Int(),
		ScheduledTime:          time.Now().Add(time.Duration(rand.Intn(100)) * time.Hour),
		StateCode:              randString(2),
		PostalCode:             randString(5),
		Timezone:               randString(3),
		CarrierID:              rand.Int(),
		PhoneNumber:            randString(10),
		MessageUID:             rand.Int(),
		AggregatorID:           rand.Int(),
		MessageSample:          randString(100),
		SampleMessageVersionID: rand.Int(),
	}
	return message
}

// Generate a random string with a given length
func randString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
