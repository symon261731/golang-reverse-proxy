package utils

import (
	"encoding/json"
	"log"
	"main/package/types"
	"os"
)

func FilterSliceById(friendList []types.UserFriends, id int) []types.UserFriends {
	var filteredSlice []types.UserFriends

	for _, friend := range friendList {
		if friend.Id != id {
			filteredSlice = append(filteredSlice, friend)
		}
	}

	return filteredSlice
}

var filepath = "../mockJson/mockJson.json"

func GetJsonFile() (types.UserListMap, error) {

	jsonFile, err := os.ReadFile(filepath)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var users types.UserListMap
	marshalError := json.Unmarshal(jsonFile, &users)

	if marshalError != nil {
		log.Println(marshalError)
		return nil, marshalError
	}

	return users, nil
}

func SetDataInJson(data types.UserListMap) error {

	decodeData, err := json.Marshal(data)
	if err != nil {
		log.Println("Произошла ошибка при форматировании данных")
		return err
	}

	errorWriteFile := os.WriteFile(filepath, decodeData, 0644)

	if errorWriteFile != nil {
		log.Println("Произошла ошибка при записи в файл")
		return err
	}

	return nil
}
