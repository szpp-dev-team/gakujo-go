package gakujo

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/szpp-dev-team/gakujo-api/model"
)

func (c *Client) Bot(classNooticeRow []model.ClassNoticeRow) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	_, _, _ = api.PostMessage(
		"botてすと",
		slack.MsgOptionText(fmt.Sprintf("%v\n", classNooticeRow), false),
	)
}
