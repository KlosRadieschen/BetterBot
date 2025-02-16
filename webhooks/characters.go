package webhooks

import "BetterScorch/database"

type Character struct {
	OwnerID    string
	Name       string
	AvatarLink string
	Brackets   string
}

func AddCharacter(character Character) error {
	err := database.Insert("Character",
		&database.DBValue{Name: "pk_ownerID", Value: character.OwnerID},
		&database.DBValue{Name: "name", Value: character.Name},
		&database.DBValue{Name: "avatar", Value: character.AvatarLink},
		&database.DBValue{Name: "brackets", Value: character.Brackets},
	)
	return err
}

func RemoveCharacter(character string) error {
	// Implementation goes here
	return nil
}
