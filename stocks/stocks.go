package stocks

import (
	"BetterScorch/database"
	"fmt"
	"strconv"
	"strings"
)

func Trade(user, company string, amount int, buying bool) (int, error) {
	if !buying {
		amount = -amount
	}

	keys := []*database.DBValue{
		{Name: "pkfk_user_owns", Value: user},
		{Name: "pkfk_company_belongsTo", Value: company},
	}

	current, err := database.Get("Shareholder", []string{"amount"}, keys...)
	if err != nil {
		return -1, err
	}

	intCurrent := 0
	if len(current) > 0 && len(current[0]) > 0 {
		intCurrent, _ = strconv.Atoi(current[0][0])
	}

	newValue := intCurrent + amount
	if newValue < 0 {
		return -1, fmt.Errorf("cannot hold negative shares")
	}

	err = database.InsertOrUpdate("Shareholder", keys, &database.DBValue{
		Name:  "amount",
		Value: strconv.Itoa(newValue),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1452") && strings.Contains(err.Error(), "fk_user_owns") {
			return -1, fmt.Errorf("You are not registered in the economy. Please use /entereconomy first")
		} else if strings.HasPrefix(err.Error(), "Error 1452") && strings.Contains(err.Error(), "fk_company_belongsTo") {
			return -1, fmt.Errorf("Please enter a valid company name")
		} else {
			return -1, err
		}
	}

	return newValue, ModifyCompanyValue(company, amount)
}

func ModifyCompanyValue(name string, delta int) error {
	// Fetch the company's current value
	current, err := database.Get("Company", []string{"value"}, &database.DBValue{
		Name:  "pk_name",
		Value: name,
	})
	if err != nil {
		return err
	}

	intCurrent, _ := strconv.Atoi(current[0][0])

	// Update the company's value
	newValue := intCurrent + delta

	// If the new value is negative, throw an error
	if newValue < 0 {
		return fmt.Errorf("company value cannot be negative")
	}

	// Get all shareholder amounts for this company
	shareholders, err := database.Get("Shareholder", []string{"pkfk_user_owns", "amount"}, &database.DBValue{
		Name:  "pkfk_company_belongsTo",
		Value: name,
	})
	if err != nil {
		return err
	}

	// Calculate the total of all current shareholder amounts
	var totalAmount int
	for _, shareholder := range shareholders {
		amt, err := strconv.Atoi(shareholder[1]) // assuming amount is the second field
		if err != nil {
			panic(err)
		}
		totalAmount += amt
	}

	// If the total amount is zero, no need to do scaling
	if totalAmount == 0 {
		return fmt.Errorf("cannot scale amounts when total shareholder amount is zero")
	}

	// Calculate the scaling factor (new value / old total)
	scalingFactor := float64(newValue) / float64(intCurrent)
	fmt.Println(scalingFactor)

	// Update each shareholder's amount based on the scaling factor
	for _, shareholder := range shareholders {
		user := shareholder[0] // user ID is the first column
		amt, _ := strconv.Atoi(shareholder[1])

		// Calculate new amount for this shareholder
		newAmount := int(float64(amt) * scalingFactor)

		// Update the shareholder's amount in the database
		err := database.Update("Shareholder",
			[]*database.DBValue{{Name: "pkfk_user_owns", Value: user}, {Name: "pkfk_company_belongsTo", Value: name}},
			&database.DBValue{Name: "amount", Value: strconv.Itoa(newAmount)},
		)
		if err != nil {
			return err
		}
	}

	// Now update the company value
	return database.Update("Company",
		[]*database.DBValue{{Name: "pk_name", Value: name}},
		&database.DBValue{Name: "value", Value: strconv.Itoa(newValue)},
	)
}

func ModifyBalance(user string, amount int) (int, error) {
	currentVal, err := database.Get("ScorchCoin", []string{"balance"}, &database.DBValue{Name: "pk_user", Value: user})
	if err != nil {
		if strings.Contains(err.Error(), "index out of range") {
			return -1, fmt.Errorf("You are not registered in the economy. Please use /entereconomy first")
		}
		return -1, err
	}

	intVal, err := strconv.Atoi(currentVal[0][0])
	if amount > intVal {
		return -1, fmt.Errorf("amount higher than current balance")
	} else if err != nil {
		return -1, fmt.Errorf("You are not registered in the economy. Please use /entereconomy first")
	}

	return amount + intVal, database.Update("ScorchCoin", []*database.DBValue{{Name: "pk_user", Value: user}}, &database.DBValue{Name: "balance", Value: strconv.Itoa(amount + intVal)})

}

func Enter(user string) error {
	return database.Insert("ScorchCoin", &database.DBValue{Name: "pk_user", Value: user}, &database.DBValue{Name: "balance", Value: "420"})
}
