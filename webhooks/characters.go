package webhooks

import (
	"BetterScorch/database"
	"BetterScorch/sender"
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Character struct {
	OwnerID    string
	Name       string
	AvatarLink string
	Brackets   string
}

var CharacterBuffer = make(map[string][]Character)

func AddCharacter(ownerID string, character Character) error {
	CharacterBuffer[ownerID] = append(CharacterBuffer[ownerID], character)

	_, err := database.Insert("Character",
		&database.DBValue{Name: "pk_ownerID", Value: character.OwnerID},
		&database.DBValue{Name: "name", Value: character.Name},
		&database.DBValue{Name: "avatar", Value: character.AvatarLink},
		&database.DBValue{Name: "brackets", Value: character.Brackets},
	)
	return err
}

func CheckAndUseCharacters(s *discordgo.Session, m *discordgo.MessageCreate) error {
	for _, character := range CharacterBuffer[m.Author.ID] {
		match, err := matchesTupperPattern(m.Content, character.Brackets)
		if err != nil {
			return err
		} else if match {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			cleaned, err := extractTupperContent(m.Content, character.Brackets)
			if err != nil {
				return err
			} else if m.MessageReference != nil {
				sender.SendCharacterReply(s, m, cleaned, character.Name, character.AvatarLink)
			} else {
				sender.SendCharacterMessage(s, m, cleaned, character.Name, character.AvatarLink)
			}
		}
	}

	return nil
}

func RetrieveCharacters() {
	// Reset/clear the buffer first to ensure clean state
	CharacterBuffer = make(map[string][]Character)

	// Get all characters from persistent storage
	characters, err := database.GetAll("Character")
	if err != nil {
		panic(err)
	}

	// Rebuild the buffer from database records
	for _, row := range characters {
		char := Character{
			OwnerID:    row[0], // ownerID
			Name:       row[1], // name
			AvatarLink: row[2], // avatar
			Brackets:   row[3], // brackets
		}
		CharacterBuffer[char.OwnerID] = append(CharacterBuffer[char.OwnerID], char)
	}
}

func RemoveCharacter(userID string, characterName string) (bool, error) {
	values := []*database.DBValue{
		{
			Name:  "pk_ownerID",
			Value: userID,
		},
		{
			Name:  "name",
			Value: characterName,
		},
	}

	affected, err := database.Remove("Character", values...)

	if err != nil {
		return false, err
	}

	CharacterBuffer[userID] = slices.DeleteFunc(CharacterBuffer[userID], func(c Character) bool {
		return c.Name == characterName
	})

	return affected > 0, err
}

// List all characters owned by the user from characterBuffer
func ListCharacters(userID string) []Character {
	if characters, ok := CharacterBuffer[userID]; ok {
		return characters
	}

	return []Character{}
}

// matchesTupperPattern returns true if the message follows the bracket template pattern.
func matchesTupperPattern(message, bracketTemplate string) (bool, error) {
	re, err := buildTupperRegex(bracketTemplate)
	if err != nil {
		return false, err
	}
	return re.MatchString(message), nil
}

// extractTupperContent extracts the user content from the message that follows the bracket template.
// It returns an error if the message does not match the expected pattern.
func extractTupperContent(message, bracketTemplate string) (string, error) {
	re, err := buildTupperRegex(bracketTemplate)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(message)
	if len(matches) < 2 {
		return "", errors.New("message does not match the defined bracket template")
	}

	// matches[1] holds the captured content.
	return matches[1], nil
}

func buildTupperRegex(bracketTemplate string) (*regexp.Regexp, error) {
	if !strings.Contains(bracketTemplate, "text") {
		return nil, errors.New("bracket template must contain the literal 'text'")
	}

	// Escape any regex metacharacters in the template.
	escapedTemplate := regexp.QuoteMeta(bracketTemplate)
	// Replace the escaped "text" with a non-greedy capture group.
	pattern := strings.Replace(escapedTemplate, "text", "(.*?)", 1)
	// Optionally anchor the pattern to match the entire message.
	pattern = "^" + pattern + "$"

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re, nil
}
