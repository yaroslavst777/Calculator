package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"unicode"
)

/*
1. Мы определяем структуру `RequestData`, которая описывает данные, ожидаемые в JSON-формате.
2. Создаем обработчик `calculateHandler`, который обрабатывает POST-запросы. Он проверяет метод запроса, декодирует JSON из тела запроса и выполняет некоторую логику обработки.
3. Используем `http.HandleFunc` для связывания URL-адреса `/api/v1/calculate` с нашим обработчиком.
4. Запускаем HTTP-сервер на порту 8080.
*/
type RequestData struct {
	Expression string `json:"expression"`
}

func Calc(expression string) (float64, error) {
	var operators []rune
	var operands []float64

	priority := func(op rune) int {
		if op == '*' || op == '/' {
			return 2
		} else if op == '+' || op == '-' {
			return 1
		}
		return 0
	}

	calculate := func() error {
		if len(operators) == 0 {
			return errors.New("invalid expression")
		}

		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		if len(operands) < 2 {
			return errors.New("invalid expression")
		}

		b := operands[len(operands)-1]
		operands = operands[:len(operands)-1]
		a := operands[len(operands)-1]
		operands = operands[:len(operands)-1]

		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return errors.New("division by zero")
			}
			result = a / b
		default:
			return errors.New("invalid operator")
		}

		operands = append(operands, result)
		return nil
	}

	var numStr string
	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		if unicode.IsDigit(char) || char == '.' {
			numStr += string(char)
		} else {
			if numStr != "" {
				num, err := strconv.ParseFloat(numStr, 64)
				if err != nil {
					return 0, fmt.Errorf("invalid number: %s", numStr)
				}
				operands = append(operands, num)
				numStr = ""
			}

			if char == '+' || char == '-' || char == '*' || char == '/' {
				for len(operators) > 0 && priority(operators[len(operators)-1]) >= priority(char) {
					if err := calculate(); err != nil {
						return 0, err
					}
				}
				operators = append(operators, char)
			} else if char == '(' {
				operators = append(operators, char)
			} else if char == ')' {
				for len(operators) > 0 && operators[len(operators)-1] != '(' {
					if err := calculate(); err != nil {
						return 0, err
					}
				}
				if len(operators) == 0 || operators[len(operators)-1] != '(' {
					return 0, errors.New("mismatched parentheses")
				}
				operators = operators[:len(operators)-1]
			} else {
				return 0, fmt.Errorf("invalid character: %c", char)
			}
		}
	}

	if numStr != "" {
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", numStr)
		}
		operands = append(operands, num)
	}

	for len(operators) > 0 {
		if err := calculate(); err != nil {
			return 0, err
		}
	}

	if len(operands) != 1 {
		return 0, errors.New("invalid expression")
	}

	return operands[0], nil
}

func makeResponse(w http.ResponseWriter, statusCode int, answer float64) {
	// Отправка JSON-ответа с ошибкой
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) //код ответа

	var response map[string]interface{}

	//Status code 200
	if statusCode == http.StatusOK {
		response = map[string]interface{}{
			"result": answer,
		}
	}

	//Status code 422
	if statusCode == http.StatusUnprocessableEntity {
		response = map[string]interface{}{
			"error": "Expression is not valid",
		}
	}

	//Status code 500
	if statusCode == http.StatusInternalServerError {
		response = map[string]interface{}{
			"error": "Internal server error",
		}
	}
	json.NewEncoder(w).Encode(response)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		makeResponse(w, http.StatusInternalServerError, 0)
		return
	}

	var requestData RequestData

	data := make([]byte, 1024)        // Создадим буфер для чтения данных в него
	num, errRead := r.Body.Read(data) // Прочитаем данные в буфер
	defer r.Body.Close()
	if errRead != nil && errRead != io.EOF {
		makeResponse(w, http.StatusInternalServerError, 0)
		return
	}

	data = data[:num]

	errUnmarshal := json.Unmarshal(data, &requestData)

	if errUnmarshal != nil {
		makeResponse(w, http.StatusInternalServerError, 0)
	}

	// Получение значения expression из формы
	expression := requestData.Expression

	answer, errCalc := Calc(expression)
	if errCalc != nil {
		makeResponse(w, http.StatusUnprocessableEntity, 0)
		return
	}

	// Отправка JSON-ответа
	makeResponse(w, http.StatusOK, answer)
}

func main() {
	http.Handle("/api/v1/calculate", http.HandlerFunc(CalcHandler))
	http.ListenAndServe(":8080", nil)
}
