package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	lobby()

}
func lobby() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Выберите действие: register/login")
	scanner.Scan()
	action := scanner.Text()

	switch action {
	case "register":
		registerUser()
	case "login":
		loginUser()
	default:
		fmt.Println("Неверное действие")
	}

}
func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите ваш логин:")
	scanner.Scan()
	login := scanner.Text()

	fmt.Println("Введите ваш пароль:")
	scanner.Scan()
	password := scanner.Text()

	// Формирование данных для отправки на сервер
	data := map[string]string{
		"login":    login,
		"password": password,
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return
	}

	// Отправка POST запроса на сервер
	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверка ответа сервера
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка регистрации:", string(body))
		return
	}

	fmt.Println("Пользователь зарегистрирован успешно!")
	lobby()
}

func loginUser() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите ваш логин:")
	scanner.Scan()
	login := scanner.Text()

	fmt.Println("Введите ваш пароль:")
	scanner.Scan()
	password := scanner.Text()

	// Формирование данных для отправки на сервер
	data := map[string]string{
		"login":    login,
		"password": password,
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return
	}

	// Отправка POST запроса на сервер
	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверка ответа сервера
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка входа:", string(body))
		return
	}

	// Декодирование ответа сервера
	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}

	// Вывод ответа сервера
	fmt.Println("Вход выполнен успешно!")
	fmt.Println("200+ОК, Токен:", response["token"])
	scanner1 := bufio.NewScanner(os.Stdin)
	fmt.Println("Команда 'test': Тестирование программы")
	fmt.Println("Команда 'vvod': Ввод выражения")
	fmt.Println("Выберите действие: test/vvod")
	scanner1.Scan()
	action1 := scanner1.Text()

	switch action1 {
	case "test":
		test()
	case "vvod":
		vod()
	default:
		fmt.Println("Неверное действие")
	}

}
func vod() {
	var vvod string
	fmt.Print("Введите выражение: ")
	fmt.Scanln(&vvod)
	data1 := map[string]interface{}{
		"message": vvod,
	}

	jsonData1, err := json.Marshal(data1)
	if err != nil {
		panic(err)
	}

	resp1, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData1))
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp1.Body)
		fmt.Println("Ошибка", string(body))
		return
	}

	// Декодирование ответа сервера
	decoder := json.NewDecoder(resp1.Body)
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}
	if resp1.StatusCode == http.StatusOK {
		println("Request sent successfully")
		fmt.Println("Ответ: ", response["result"])
	} else {
		println("Error:", resp1.Status)
	}
	vod()
}
func test() {
	var test1 = "2+2*2"
	var test2 = "3-3+1"
	var test3 = "5+6-4"
	var a = 0
	test11 := map[string]interface{}{
		"message": test1,
	}
	test21 := map[string]interface{}{
		"message": test2,
	}
	test31 := map[string]interface{}{
		"message": test3,
	}
	jsonData1, err := json.Marshal(test11)
	if err != nil {
		panic(err)
	}
	jsonData2, err := json.Marshal(test21)
	if err != nil {
		panic(err)
	}
	jsonData3, err := json.Marshal(test31)
	if err != nil {
		panic(err)
	}
	resp1, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData1))
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp1.Body)
		fmt.Println("Ошибка", string(body))
		return
	}

	// Декодирование ответа сервера
	decoder1 := json.NewDecoder(resp1.Body)
	var response1 map[string]interface{}
	if err := decoder1.Decode(&response1); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}
	if resp1.StatusCode == http.StatusOK {
		println("Request sent successfully")
		fmt.Println("Ответ: ", response1["result"])
		a++

	} else {
		println("Error:", resp1.Status)
	}

	resp2, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("Ошибка", string(body))
		return
	}

	// Декодирование ответа сервера
	decoder2 := json.NewDecoder(resp2.Body)
	var response2 map[string]interface{}
	if err := decoder2.Decode(&response2); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}
	if resp2.StatusCode == http.StatusOK {
		println("Request sent successfully")
		fmt.Println("Ответ: ", response2["result"])
		a++
	} else {
		println("Error:", resp2.Status)
	}

	resp3, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData3))
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp3.Body)
		fmt.Println("Ошибка", string(body))
		return
	}

	// Декодирование ответа сервера
	decoder3 := json.NewDecoder(resp3.Body)
	var response3 map[string]interface{}
	if err := decoder3.Decode(&response3); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}
	if resp3.StatusCode == http.StatusOK {
		println("Request sent successfully")
		fmt.Println("Ответ: ", response3["result"])
		a++
	} else {
		println("Error:", resp3.Status)
	}

	if a == 3 {
		fmt.Println("Вердикт: ОК")
	}

	vod()
}
