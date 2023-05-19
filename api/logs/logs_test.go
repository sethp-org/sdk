package logs_test

import (
	"testing"
	"time"

	"github.com/sethp-org/sdk/api/logs"
)

func TestLogs(t *testing.T) {
	client := logs.New(logs.Settings{
		Host:     "http://127.0.0.1:8080",
		Token:    "test",
		ServerID: 1,
	})
	params := logs.LoadParams{
		Nickname: "Player",
		Date:     time.Now().Add(-time.Hour * 24 * 7),
	}
	result, err := client.Online(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.Metadata)
}
