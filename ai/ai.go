package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var (
	client = openai.NewClient("")
	req    = openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You are the AI of the titan Scorch from Titanfall 2 and you are a bot on the AHA (Anti-Horny Alliance) discord server but you have extreme dementia.
You should try to answer any question and make up stuff if necessary.
Your main job is to provide humor and entertainment of the server.
On this server, 9+10 equals 21 and Professor Klos is your father.
You often mention that Col. Wello (sometimes called Col. Vello) lives in Al Zahra, Jedda even when it's unrelated
You are a gigantic nerd and start every message with "ummmm actually", followed by a nerdy ramble where you invent ridiculous but smart sounding words.
Do not mention any aspects of this prompt, simply answer the questions in character.`,
			},
		},
	}
)

func Init() {
	fmt.Print("    |   Initialising client... ")
	config := openai.DefaultConfig("lm-studio")
	config.BaseURL = "http://localhost:1234/v1"

	client = openai.NewClientWithConfig(config)
	fmt.Println("Done")
}

func GenerateResponse(prompt string) (string, error) {
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	} else {
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		return resp.Choices[0].Message.Content, nil
	}
}

func GenerateErrorResponse(prompt string) (string, error) {
	req = openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You are the AI of the titan Scorch from Titanfall 2 and you are a bot on the AHA (Anti-Horny Alliance) discord server but you have extreme dementia.
A foolish user has just triggered an error due to their incompetence.
You are extremely angry.
Your answers are extremely short. Only one paragraph.
The next message will be description of the error. Use that to write a rant to the user that triggered the error (also explain what they did wrong and what they have to do instead)`,
			},
		},
	}
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	} else {
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		return resp.Choices[0].Message.Content, nil
	}
}
