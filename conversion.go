package etherutils

// Copyright 2017 Orinoco Payments
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Turn a string in to number of Wei.
// The string can be a simple number of Wei, e.g. "1000000000000000" or it can
// be a number followed by a unit, e.g. "10 ether".  Unit names are
// case-insensitive, and can be either given names (e.g. "finney") or metric
// names (e.g. "mlliether").
// Note that this function expects use of the period as the decimal separator.
func StringToWei(input string) (*big.Int, error) {
	if input == "" {
		return nil, errors.New("Failed to parse empty value")
	}
	var result big.Int
	// See if the input contains a space
	if strings.Contains(input, " ") {
		// Break the value in to number and space
		s := strings.Split(input, " ")

		if len(s) != 2 {
			return nil, errors.New(fmt.Sprintf("Unknown format of %s", input))
		}

		if strings.Contains(s[0], ".") {
			err := decimalStringToWei(s[0], s[1], &result)
			if err != nil {
				return nil, err
			}
		} else {
			err := integerStringToWei(s[0], s[1], &result)
			if err != nil {
				return nil, err
			}
		}

	} else {
		// No space so we're a simple number of Wei
		result.SetString(input, 10)
	}

	// Ensure we don't have a negative number
	if result.Cmp(new(big.Int)) < 0 {
		return nil, errors.New("Value resulted in negative number of Wei")
	}

	// return nil, errors.New("Failed to parse value")
	return &result, nil
}

// Used in WeiToString
var zero = big.NewInt(0)
var thousand = big.NewInt(1000)
var million = big.NewInt(1000000)

// Turn a number of Wei in to a string.
// If the 'standard' argument is true then this will display the value
// in either (KMG)Wei or Ether only
func WeiToString(input *big.Int, standard bool) (string, error) {
	postfixPos := 0
	modInt := new(big.Int).Set(input)
	// Step 1: step down whole thousands for our first attempt at the unit
	for input.Cmp(thousand) >= 0 && modInt.Mod(input, thousand).Cmp(zero) == 0 {
		postfixPos++
		input = input.Div(input, thousand)
		modInt = modInt.Set(input)
	}

	// Step 2: move to a fraction if sensible

	// Because of the innacuracy of floating point we use string manipulation
	// to place the decimal in the correct position
	outputValue := input.Text(10)

	desiredPostfixPos := postfixPos
	if len(outputValue) > 3 {
		desiredPostfixPos += len(outputValue) / 3
		if len(outputValue)%3 == 0 {
			desiredPostfixPos -= 1
		}
	}
	decimalPlace := len(outputValue)
	if desiredPostfixPos > 3 && standard {
		// We want this in a standard unit.  We will show up to
		// 999999999999 in (KMG)Wei and anything higher in Ether
		desiredPostfixPos = 6
	}
	for postfixPos < desiredPostfixPos {
		decimalPlace -= 3
		postfixPos++
	}
	for postfixPos > desiredPostfixPos {
		outputValue = outputValue + strings.Repeat("0", 3)
		decimalPlace += 3
		postfixPos--
	}
	if decimalPlace <= 0 {
		outputValue = "0." + strings.Repeat("0", 0-decimalPlace) + outputValue
	} else if decimalPlace < len(outputValue) {
		outputValue = outputValue[:decimalPlace] + "." + outputValue[decimalPlace:]
	}

	// Trim trailing zeros if this is a decimal
	if strings.Contains(outputValue, ".") {
		outputValue = strings.TrimRight(outputValue, "0")
	}

	// Return our value
	return fmt.Sprintf("%s %s", outputValue, metricUnits[postfixPos]), nil
}

func decimalStringToWei(amount string, unit string, result *big.Int) error {
	// Because floating point maths is not accurate we need to break potentially
	// large decimal fractions in to two separate pieces: the integer part and the
	// decimal part
	parts := strings.Split(amount, ".")

	// The value for the integer part of the number is easy
	if parts[0] != "" {
		err := integerStringToWei(parts[0], unit, result)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to parse %s %s", amount, unit))
		}
	}

	// The value for the decimal part of the number is harder.  We left-shift it
	// so that we end up multiplying two integers rather than two floats, as the
	// latter is unreliable

	// Obtain multiplier
	// This will never fail because it is already called above in integerStringToWei()
	multiplier, _ := UnitToMultiplier(unit)

	// Trim trailing 0s
	trimmedDecimal := strings.TrimRight(parts[1], "0")
	if len(trimmedDecimal) == 0 {
		// Nothing more to do
		return nil
	}
	var decVal big.Int
	decVal.SetString(trimmedDecimal, 0)

	// Divide multiplier by 10^len(trimmed decimal) to obtain sane value
	div := big.NewInt(10)
	for i := 0; i < len(trimmedDecimal); i++ {
		multiplier.Div(multiplier, div)
	}

	// Ensure we don't have a fractional number of Wei
	if multiplier.Sign() == 0 {
		return errors.New("Value resulted in fractional number of Wei")
	}

	var decResult big.Int
	decResult.Mul(multiplier, &decVal)

	// Add it to the integer result
	result.Add(result, &decResult)

	return nil
}

func integerStringToWei(amount string, unit string, result *big.Int) error {
	// Obtain number
	number := new(big.Int)
	_, success := number.SetString(amount, 10)
	if !success {
		return errors.New(fmt.Sprintf("Failed to parse numeric value of %s %s", amount, unit))
	}

	// Obtain multiplier
	multiplier, err := UnitToMultiplier(unit)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse unit of %s %s", amount, unit))
	}

	result.Mul(number, multiplier)
	return nil
}

// Metric units
var metricUnits = [...]string{"Wei", "KWei", "MWei", "GWei", "Microether", "Milliether", "Ether", "Kiloether", "Megaether", "Gigaether", "Teraether"}

// Named units
// var namedUnits = [...]string{"Wei", "Ada", "Babbage", "Shannon", "Szazbo", "Finney", "Ether", "Einstein", "Kilo", "Mega", "Giga", "Tera"}

// Take an Ethereum unit and return the relevant multiplier from Wei.
func UnitToMultiplier(unit string) (*big.Int, error) {
	result := new(big.Int)
	switch strings.ToLower(unit) {
	case "", "wei":
		result.SetString("1", 10)
	case "ada", "kwei", "kilowei":
		result.SetString("1000", 10)
	case "babbage", "mwei", "megawei":
		result.SetString("1000000", 10)
	case "shannon", "gwei", "gigawei":
		result.SetString("1000000000", 10)
	case "szazbo", "micro", "microether":
		result.SetString("1000000000000", 10)
	case "finney", "milli", "milliether":
		result.SetString("1000000000000000", 10)
	case "ether":
		result.SetString("1000000000000000000", 10)
	case "einstein", "kilo", "kiloether":
		result.SetString("1000000000000000000000", 10)
	case "mega", "megaether":
		result.SetString("1000000000000000000000000", 10)
	case "giga", "gigaether":
		result.SetString("1000000000000000000000000000", 10)
	case "tera", "teraether":
		result.SetString("1000000000000000000000000000000", 10)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown unit %s", unit))
	}
	return result, nil
}
