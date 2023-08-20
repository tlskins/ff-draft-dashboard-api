package parsers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ayush6624/go-chatgpt"

	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func CalcPlayerReport(player *t.Player, client *chatgpt.Client) (report *t.PlayerReport, err error) {
	ctx := context.Background()
	report = &t.PlayerReport{
		Id:        player.Id,
		CreatedAt: time.Now(),
		Position:  player.Position,
		Name:      player.Name,
	}

	if player.EspnPlayerOutlook == "" {
		return
	}

	question := fmt.Sprintf(
		`Please summarize the following nfl player news with regard to their productivity in the upcoming season. Summaries should be in 5-12 word bullet points organized under 2 sections separated by two new lines \n\n. First section "Pros:" and then "Cons:": %s`,
		player.EspnPlayerOutlook,
	)

	res, err := client.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: question,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var respContent string
	if len(res.Choices) > 0 {
		respContent = res.Choices[0].Message.Content
	}
	respParts := strings.Split(respContent, "\n\n")
	if len(respParts) > 0 {
		report.Pros = respParts[0]
	}
	if len(respParts) > 1 {
		report.Cons = respParts[1]
	}

	return
}
