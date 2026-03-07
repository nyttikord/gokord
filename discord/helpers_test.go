package discord

import (
	"testing"
	"time"
)

func TestSnowflakeTimestamp(t *testing.T) {
	// #discordgo channel ID :)
	parsedTimestamp := SnowflakeTimestamp(155361364909621248)
	correctTimestamp := time.Date(2016, time.March, 4, 17, 10, 35, 869*1000000, time.UTC)
	if !parsedTimestamp.Equal(correctTimestamp) {
		t.Errorf("parsed time incorrect: got %v, want %v", parsedTimestamp, correctTimestamp)
	}
}
