package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

// Структура для хранения данных пользователя
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Глобальная мапа для хранения пользователей
var users = make(map[string]*User)

func main() {
	db, _ := sql.Open("sqlite3", "./gopher.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people (vod TEXT)")
	statement.Exec()
	//_, err := db.Exec("DELETE FROM people")
	//if err != nil {
	//panic(err)
	//}
	statement, _ = db.Prepare("INSERT INTO people (vod) VALUES (?)")
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	//http.HandleFunc("/test", )
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//1

		var data map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		v, ok := data["message"].(string)
		if !ok {
			http.Error(w, "invalid data type", http.StatusBadRequest)
			return
		}
		str := v
		statement.Exec(str)
		rows, _ := db.Query("SELECT vod FROM people")
		fmt.Println("Предыдущии запросы:")
		var str2 string
		for rows.Next() {
			rows.Scan(&str2)
			fmt.Println(str2)
		}
		fmt.Printf("Последнее выражение: %v", str)
		//2
		fmt.Println()
		i := infixToPostfix(str)
		result := evaluatePostfix(i)
		fmt.Printf("Ответ: %v\n", result)

		response1 := map[string]interface{}{
			"message": "OK",
			"result":  result,
		}

		// Преобразование ответа в JSON
		jsonResponse1, err := json.Marshal(response1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse1)
		return

	})
	http.ListenAndServe("localhost:8080", nil)
}

func precedence(operator rune) int {
	if operator == '*' || operator == '/' {
		return 2
	} else if operator == '+' || operator == '-' {
		return 1
	}
	return 0
}

func infixToPostfix(expression string) string {
	var result string
	var stack []rune
	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			result += string(char) + " "
		case char == '(':
			stack = append(stack, char)
		case char == ')':
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				result += string(stack[len(stack)-1]) + " "
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		default:
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(char) {
				result += string(stack[len(stack)-1]) + " "
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		}
	}
	for len(stack) > 0 {
		result += string(stack[len(stack)-1]) + " "
		stack = stack[:len(stack)-1]
	}
	return result
}
func evaluatePostfix(postfix string) int {
	var stack []int
	for _, char := range postfix {
		if char >= '0' && char <= '9' {
			digit, _ := strconv.Atoi(string(char))
			stack = append(stack, digit)
		} else if char == '+' {
			op1 := stack[len(stack)-2]
			op2 := stack[len(stack)-1]
			stack[len(stack)-2] = op1 + op2
			stack = stack[:len(stack)-1]
		} else if char == '-' {
			op1 := stack[len(stack)-2]
			op2 := stack[len(stack)-1]
			stack[len(stack)-2] = op1 - op2
			stack = stack[:len(stack)-1]
		} else if char == '*' {
			op1 := stack[len(stack)-2]
			op2 := stack[len(stack)-1]
			stack[len(stack)-2] = op1 * op2
			stack = stack[:len(stack)-1]
		}
	}
	return stack[0]
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Разрешены только POST-запросы", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка, существует ли пользователь с таким логином
	if _, ok := users[user.Login]; ok {
		http.Error(w, "Пользователь с таким логином уже существует", http.StatusBadRequest)
		return
	}

	// Добавление пользователя в мапу
	users[user.Login] = &user
	fmt.Printf("Пользователь зарегистрирован: %s\n", user.Login)

	// Возвращение OK ответа
	w.WriteHeader(http.StatusOK)
}

// Обработчик запроса входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Разрешены только POST-запросы", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя и пароля
	if storedUser, ok := users[user.Login]; ok && storedUser.Password == user.Password {
		// Генерация токена (простой пример, в реальном случае используйте более надежную систему)

		const token = "super_secret_signature"
		now := time.Now()
		token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"name": "user_name",
			"nbf": now.Add(time.Minute).Unix(),
			"exp": now.Add(5 * time.Minute).Unix(),
			"iat": now.Unix(),
		})

		tokenString, err := token1.SignedString([]byte(token)) // Подписание токена с использованием секретного ключа
		if err != nil {
			panic(err) // Обработка ошибки при подписи токена
		}

		// Вывод в консоль строки токена
		// Формирование ответа

		response := map[string]interface{}{
			"message": "OK",
			"token":   tokenString,
		}

		// Преобразование ответа в JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		return
	}

	http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
}
