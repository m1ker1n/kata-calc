package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	const (
		minOperand = 1
		maxOperand = 10
	)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ToUpper(text)

		// tokens := strings.Split(text, " ") checks only " ", there are more white space characters
		tokens := strings.Fields(text)
		if len(tokens) != 3 {
			panic(errors.New("must be binary operation"))
		}

		operand1, operand1System, err := parseOperand(tokens[0])
		if err != nil {
			panic(err)
		}
		if operand1 < minOperand || operand1 > maxOperand {
			panic(fmt.Errorf("operand1=%d must be in boundaries [%d; %d]", operand1, minOperand, maxOperand))
		}

		operand2, operand2System, err := parseOperand(tokens[2])
		if err != nil {
			panic(err)
		}
		if operand1 < minOperand || operand1 > maxOperand {
			panic(fmt.Errorf("operand2=%d must be in boundaries [%d; %d]", operand2, minOperand, maxOperand))
		}

		if operand1System != operand2System {
			panic(errors.New("operands should have same numeral systems"))
		}

		operation := tokens[1]

		operationFn, ok := availableOperations[operation]
		if !ok {
			panic(errors.New("operation not found"))
		}

		result, err := operationFn(operand1, operand2)
		if err != nil {
			panic(err)
		}

		var strResult string
		switch operand1System {
		case Arabic:
			strResult = strconv.Itoa(result)
		case Roman:
			strResult, err = convertToRoman(result)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(strResult)
	}
}

type binaryOperation func(v1, v2 int) (int, error)

var availableOperations = map[string]binaryOperation{
	"+": func(v1, v2 int) (int, error) {
		return v1 + v2, nil
	},
	"-": func(v1, v2 int) (int, error) {
		return v1 - v2, nil
	},
	"*": func(v1, v2 int) (int, error) {
		return v1 * v2, nil
	},
	"/": func(v1, v2 int) (int, error) {
		if v2 == 0 {
			return 0, errors.New("could not divide by 0")
		}
		return v1 / v2, nil
	},
}

type numeralSystem int

const (
	Undefined numeralSystem = iota
	Arabic
	Roman
)

func parseOperand(s string) (int, numeralSystem, error) {
	arabic, err := strconv.Atoi(s)
	if err == nil {
		return arabic, Arabic, nil
	}

	roman, err := parseRoman(s)
	if err == nil {
		return roman, Roman, nil
	}

	return 0, Undefined, errors.New("undefined numeral system")
}

var (
	romanToArabic = []struct {
		str string
		val int
	}{
		{"M", 1000},
		{"CM", 900},
		{"D", 500},
		{"CD", 400},
		{"C", 100},
		{"XC", 90},
		{"L", 50},
		{"XL", 40},
		{"X", 10},
		{"IX", 9},
		{"V", 5},
		{"IV", 4},
		{"I", 1},
	}
	maxRomanNumber = 3999
)

func parseRoman(s string) (int, error) {
	var result int

	for i, prevDigit, sameDigitInRow := 0, maxRomanNumber, 0; i < len(s); {
		transformationFound := false
		for _, transformation := range romanToArabic {
			if !strings.HasPrefix(s[i:], transformation.str) {
				continue
			}
			digit := transformation.val
			if prevDigit < digit {
				return 0, fmt.Errorf("wrong digit position: %d stays before %d", prevDigit, digit)
			}
			if prevDigit == digit {
				sameDigitInRow++
			} else {
				sameDigitInRow = 1
			}
			if sameDigitInRow > 3 {
				return 0, fmt.Errorf("there can't be more than 3 equal digits %d in row", digit)
			}

			prevDigit = digit
			result += digit
			i += len(transformation.str)
			transformationFound = true
			break
		}
		if !transformationFound {
			return 0, fmt.Errorf("unknown digits in %s", s)
		}
	}

	return result, nil
}

func convertToRoman(v int) (string, error) {
	if v <= 0 || v > maxRomanNumber {
		return "", fmt.Errorf("%d is beyond the Roman boundaries (0, 3999]", v)
	}
	result := ""
	for v > 0 {
		for _, transformation := range romanToArabic {
			if transformation.val <= v {
				result += transformation.str
				v -= transformation.val
				break
			}
		}
	}
	return result, nil
}
