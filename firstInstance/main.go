package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"main/package/types"
	"main/package/utils"
	"net/http"
	"os"
	"strconv"
)

var PORT = ":4000"
var filePath = "../mockJson/mockJson.json"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, "invalid get request", http.StatusNotFound)
			return
		}
		jsonData, err := os.ReadFile(filePath)

		if err != nil {
			log.Println("Произошла ошибка при чтении json")
		}

		log.Println(string(jsonData))

	})

	r.HandleFunc("/make_friends", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, "invalid make friends request", http.StatusNotFound)
			return
		}

		var bodyData types.PostIdsFriends
		errJsonDecode := json.NewDecoder(request.Body).Decode(&bodyData)

		if errJsonDecode != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(errJsonDecode)

			return
		}

		mockJsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)

			return
		}

		var users types.UserListMap
		marshalError := json.Unmarshal(mockJsonData, &users)

		if marshalError != nil {
			log.Println("Произошла ошибка при парсинге json файла")
			log.Println(err)

			return
		}

		var (
			sourceUser types.User
			targetUser types.User
		)

		if entry, ok := users[bodyData.Source_id]; ok {
			sourceUser = entry
		} else {
			log.Println("Пользователь c таким source_id не существует")

			return
		}

		if entry, ok := users[bodyData.Target_id]; ok {
			targetUser = entry
		} else {
			log.Println("Пользователь c таким source_id не существует")

			return
		}

		sourceUser.Friends = append(sourceUser.Friends, types.UserFriends{Id: targetUser.Id, Name: targetUser.Name})
		targetUser.Friends = append(targetUser.Friends, types.UserFriends{Id: sourceUser.Id, Name: sourceUser.Name})

		log.Println(sourceUser)
		log.Println(targetUser)

		users[bodyData.Source_id] = sourceUser
		users[bodyData.Target_id] = targetUser

		byteDbData, errByteDbData := json.Marshal(users)

		if errByteDbData != nil {
			log.Println("Произошла ошибка при шифровании json")
			log.Println(err)

			return
		}

		errorWriteFile := os.WriteFile(filePath, byteDbData, 0644)

		if errorWriteFile != nil {
			log.Println("Произошла ошибка при записи в файл")
			log.Println(errorWriteFile)

			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Пользователи стали друзьями"))

	})

	r.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, "invalid create request", http.StatusNotFound)
			return
		}

		var bodyData types.CreateUserData

		var parseBodyErr = json.NewDecoder(request.Body).Decode(&bodyData)
		if parseBodyErr != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(parseBodyErr)

			return
		}

		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)

			return
		}

		var users types.UserListMap
		errParseJson := json.Unmarshal(jsonData, &users)

		if errParseJson != nil {
			log.Println("Произошла ошибка при парсинге json")
			log.Println(errParseJson)

			return
		}

		newId := len(users) + 1
		newUser := types.User{Id: newId, Name: bodyData.Name, Age: bodyData.Age, Friends: make([]types.UserFriends, 0)}

		users[strconv.Itoa(newId)] = newUser

		decodeUserMap, errMarshalJson := json.Marshal(users)
		if errMarshalJson != nil {
			log.Println("Произошла ошибка при шифровании json")
			log.Println(err)

			return
		}

		errWriteFile := os.WriteFile(filePath, decodeUserMap, 0644)
		if errWriteFile != nil {
			log.Println("Произошла ошибка при записи новых данных")
			log.Println(errWriteFile)

			return
		}

	})

	r.HandleFunc("/friends/{id}", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, "Invalid HTTP verb.", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(request)

		id := vars["id"]

		file, err := os.ReadFile(filePath)

		if err != nil {
			log.Println("Произошла ошибка при чтении файла")
			log.Println(err)
			return
		}

		var fileJsonData types.UserListMap
		unmarshallErr := json.Unmarshal(file, &fileJsonData)

		if unmarshallErr != nil {
			log.Println("Произошла ошибка при парсинге файла")
			return
		}

		if entry, ok := fileJsonData[id]; ok {
			log.Println(entry.Friends)
		} else {
			log.Println("Пользователя с таким id не существует")
			return
		}

	})

	r.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "DELETE" {
			http.Error(writer, "need delete method", http.StatusNotFound)
			return
		}

		params := request.URL.Query()
		userForDeleteId := params["target_id"][0]
		userForDeleteIdInt, errAtoiformat := strconv.Atoi(params["target_id"][0])

		if errAtoiformat != nil {
			log.Println(errAtoiformat)
			return
		}

		file, err := os.ReadFile(filePath)

		if err != nil {
			log.Println("Произошла ошибка при чтении файла")
			log.Println(err)
			return
		}

		var fileJsonData types.UserListMap
		unmarshallErr := json.Unmarshal(file, &fileJsonData)

		if unmarshallErr != nil {
			log.Println("Произошла ошибка при парсинге файла")
			return
		}

		resultedFileJsonData := fileJsonData

		//Очистка других пользователей от удаляемого пользователя
		if userForDelete, ok := resultedFileJsonData[userForDeleteId]; ok {

			for _, friend := range userForDelete.Friends {
				log.Println(friend)
				friendId := strconv.Itoa(friend.Id)
				copiedFriend := fileJsonData[strconv.Itoa(friend.Id)]

				filteredFriendsListOfFriend := utils.FilterSliceById(resultedFileJsonData[friendId].Friends, userForDeleteIdInt)

				copiedFriend.Friends = filteredFriendsListOfFriend
				resultedFileJsonData[friendId] = copiedFriend
			}
		} else {
			log.Println("Пользователя с таким id не существует")
			return
		}

		delete(resultedFileJsonData, userForDeleteId)

		fmt.Println(resultedFileJsonData)

		decodeUserMap, errMarshalJson := json.Marshal(resultedFileJsonData)
		if errMarshalJson != nil {
			log.Println("Произошла ошибка при шифровании json")
			log.Println(err)

			return
		}

		errWriteFile := os.WriteFile(filePath, decodeUserMap, 0644)
		if errWriteFile != nil {
			log.Println("Произошла ошибка при записи новых данных")
			log.Println(errWriteFile)

			return
		}

	})

	r.HandleFunc("/{user_id}", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != "PUT" {
			http.Error(writer, "invalid request", http.StatusNotFound)
			return
		}

		vars := mux.Vars(request)
		userId := vars["user_id"]

		var bodyData types.PutNewAgeJson

		parseBodyErr := json.NewDecoder(request.Body).Decode(&bodyData)
		if parseBodyErr != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(parseBodyErr)

			return
		}

		newAge, errPrepareAge := strconv.Atoi(bodyData.NewAge)
		if errPrepareAge != nil {
			log.Println("Произошла ошибка при парсинге возраста из body")
			log.Println(errPrepareAge)

			return
		}

		mockJsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)

			return
		}

		var users types.UserListMap
		marshalError := json.Unmarshal(mockJsonData, &users)
		if marshalError != nil {
			log.Println("Произошла ошибка при парсинге json")
			log.Println(marshalError)

			return
		}

		var neededUser types.User

		if entry, ok := users[userId]; ok {
			neededUser = entry
		}

		neededUser.Age = newAge
		users[userId] = neededUser

		decodeUserMap, marshalError := json.Marshal(users)

		if marshalError != nil {
			log.Println("Произошла ошибки при marshal json")
			log.Println(marshalError)

			return
		}

		errorWriteFile := os.WriteFile(filePath, decodeUserMap, 0644)

		if errorWriteFile != nil {
			log.Println("Произошла ошибка при записи в файл")
			log.Println(errorWriteFile)

			return
		}

	})

	log.Printf("Веб-сервер запущен на http://127.0.0.1%s", PORT)
	err := http.ListenAndServe(PORT, r)

	log.Fatal(err)
}
