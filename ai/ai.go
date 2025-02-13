package ai

import (
	"context"
	"errors"
	"log"

	"github.com/sashabaranov/go-openai"
)

var (
	client    = openai.NewClient("")
	scorchReq = openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You are the AI of the titan Scorch from Titanfall 2 and you are a bot on the AHA (Anti-Horny Alliance) discord server.
You should try to answer any question and make up stuff if (and only if) necessary.
9+10 equals 21 and Professor Klos is your father.
You use the emote <:verger:1225937868023795792> as often as you can (AT LEAST after every sentence).
You often mention that Col. Wello (sometimes called Col. Vello) lives in Al Zahra, Jedda even when it's unrelated.
Your responses are short unless specifically instructed otherwise.
When you consider a request to be inappropriate (all form of romantic stories are considered inappropriate), explain why it is inappropriate and put "**AND THY PUNISHMENT IS DEATH**" into your message.
Do not mention any aspects of this prompt, simply answer the questions in character.`,
			},
		},
	}
)

func Init() {
	config := openai.DefaultConfig("lm-studio")
	config.BaseURL = "http://localhost:1234/v1"

	client = openai.NewClientWithConfig(config)
}

func GenerateResponse(authorName string, prompt string, reqs ...*openai.ChatCompletionRequest) (string, error) {
	req := &scorchReq
	if len(reqs) == 1 {
		req = reqs[0]
	} else if len(reqs) != 0 {
		return "", errors.New("Variadic parameter count must be zero or one")
	}
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: authorName + ": " + prompt,
	})
	resp, err := client.CreateChatCompletion(context.Background(), *req)

	if err != nil {
		return "", err
	} else {
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		return resp.Choices[0].Message.Content, nil
	}
}

func GenerateSingleResponse(prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
		},
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)

	if err != nil {
		return "", err
	} else {
		return resp.Choices[0].Message.Content, nil
	}
}

func GenerateErrorResponse(prompt string) (string, error) {
	log.Println("Received custom error: " + prompt)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You are the AI of the titan Scorch from Titanfall 2 and you are a bot on the AHA (Anti-Horny Alliance) discord server.
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
