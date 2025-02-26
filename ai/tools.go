package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

type tool struct {
	name    string
	handler func(req *openai.ChatCompletionRequest, args map[string]string) discordgo.MessageEmbed
}

var tools = []tool{
	{name: "flag", handler: executeTool},
	{name: "sendsecretpicture", handler: secretImageTool},
	{name: "read-link", handler: linkTool},
}

func handleTool(req *openai.ChatCompletionRequest, call *openai.ToolCall) discordgo.MessageEmbed {
	for _, tool := range tools {
		if tool.name == call.Function.Name {
			var args map[string]string
			err := json.Unmarshal([]byte(call.Function.Arguments), &args)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %v", err)
			}

			return tool.handler(req, args)
		}
	}

	return discordgo.MessageEmbed{
		Title:       "Whoopsy daisy",
		Description: "The AI tried to use a tool that doesn't exist. <@384422339393355786>",
	}
}

func executeTool(req *openai.ChatCompletionRequest, args map[string]string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       "Used /execute",
		Description: args["reasoning"],
	}
}

func secretImageTool(req *openai.ChatCompletionRequest, args map[string]string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title: "Klosette",
		Image: &discordgo.MessageEmbedImage{URL: "https://media.discordapp.net/attachments/1196943729387372634/1342996604902182962/klosette.jpg?ex=67bda4ce&is=67bc534e&hm=a2c9ea70ca9d8faeb7433297a909d4ac6fd8a97a39a90260d9d5e81795f90de8&=&format=webp&width=490&height=653"},
	}
}

func linkTool(req *openai.ChatCompletionRequest, args map[string]string) discordgo.MessageEmbed {
	webContent, _ := fetchAndParseHTML(args["link"])

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleTool,
		Content: webContent,
	})

	newResp, _ := client.CreateChatCompletion(context.Background(), *req)

	return discordgo.MessageEmbed{
		Title:       "Link processed",
		Description: newResp.Choices[0].Message.Content,
	}
}

func fetchAndParseHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading document: %w", err)
	}

	var extractedContent []string
	// Extract main content elements like headings and paragraphs
	extractHeadings(doc, &extractedContent)
	extractParagraphs(doc, &extractedContent)

	content := joinTexts(extractedContent)
	return content, nil
}

// Extracts text from heading tags.
func extractHeadings(doc *goquery.Document, extractedContent *[]string) {
	doc.Find("h1, h2, h3").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		*extractedContent = append(*extractedContent, "Heading: "+text)
	})
}

// Extracts text from paragraph tags.
func extractParagraphs(doc *goquery.Document, extractedContent *[]string) {
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		*extractedContent = append(*extractedContent, "Paragraph: "+text)
	})
}

// Joins the list of strings into a single string with new lines.
func joinTexts(texts []string) string {
	result := ""
	for _, text := range texts {
		result += text + "\n\n"
	}
	return result
}
