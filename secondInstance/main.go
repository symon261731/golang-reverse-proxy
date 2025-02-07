package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/package/types"
	"main/package/utils"
	"net/http"
	"strconv"
)

var PORT = ":4010"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, "invalid get request", http.StatusNotFound)
			return
		}
		mockJsonData, err := utils.GetJsonFile()

		if err != nil {
			log.Println("Произошла ошибка при чтении json")
		}

		log.Println(mockJsonData)

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

		users, readFileError := utils.GetJsonFile()

		if readFileError != nil {
			log.Println("Произошла ошибка при парсинге json файла")
			log.Println(readFileError)

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

		errorSetData := utils.SetDataInJson(users)

		if errorSetData != nil {
			log.Println(errorSetData)
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

		var users, getJsonFileError = utils.GetJsonFile()

		if getJsonFileError != nil {
			log.Println(getJsonFileError)
			return
		}

		newId := len(users) + 1
		newUser := types.User{Id: newId, Name: bodyData.Name, Age: bodyData.Age, Friends: make([]types.UserFriends, 0)}

		users[strconv.Itoa(newId)] = newUser

		errorSetData := utils.SetDataInJson(users)

		if errorSetData != nil {
			log.Println(errorSetData)
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

		var fileJsonData, unmarshallErr = utils.GetJsonFile()

		if unmarshallErr != nil {
			log.Println("Произошла ошибка при парсинге файла")
			log.Println(unmarshallErr)
			return
		}

		if entry, ok := fileJsonData[id]; ok {
			log.Println(entry.Friends)

			byteJson, errMarshal := json.Marshal(entry.Friends)

			if errMarshal != nil {
				log.Println(errMarshal)
				return
			}

			writer.WriteHeader(http.StatusOK)
			writer.Write(byteJson)
		} else {
			log.Println("Пользователя с таким id не существует")
			writer.WriteHeader(http.StatusNotFound)
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

		var fileJsonData, unmarshallErr = utils.GetJsonFile()

		if unmarshallErr != nil {
			log.Println("Произошла ошибка при парсинге файла")
			log.Println(unmarshallErr)
			return
		}

		resultedFileJsonData := fileJsonData

		//! Очистка других пользователей от удаляемого пользователя
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

		errSetData := utils.SetDataInJson(resultedFileJsonData)

		if errSetData != nil {
			log.Println(errSetData)
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

		var users, getJsonError = utils.GetJsonFile()
		if getJsonError != nil {
			log.Println("Произошла ошибка при парсинге json")
			log.Println(getJsonError)

			return
		}

		var neededUser types.User

		if entry, ok := users[userId]; ok {
			neededUser = entry
		}

		neededUser.Age = newAge
		users[userId] = neededUser

		errorSetData := utils.SetDataInJson(users)

		if errorSetData != nil {
			log.Println(errorSetData)
			return

		}

	})

	log.Printf("Веб-сервер запущен на http://127.0.0.1%s", PORT)
	err := http.ListenAndServe(PORT, r)

	log.Fatal(err)
}
