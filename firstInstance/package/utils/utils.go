package utils

import "main/package/types"

func FilterSliceById(friendList []types.UserFriends, id int) []types.UserFriends {
	var filteredSlice []types.UserFriends

	for _, friend := range friendList {
		if friend.Id != id {
			filteredSlice = append(filteredSlice, friend)
		}
	}

	return filteredSlice
}
