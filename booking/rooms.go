package booking

import (
	"errors"
	"math/rand/v2"
	"slices"
)

type Room struct {
	Id    int
	Title string
}

var Rooms []Room = []Room{
	{Id: 1, Title: "Room A"},
	{Id: 2, Title: "Room B"},
	{Id: 3, Title: "Room C"},
}

func GetRoomById(id int) (*Room, error) {
	idx := slices.IndexFunc(Rooms, func(room Room) bool { return room.Id == id })
	if idx == -1 {
		return nil, errors.New("Unknown room id")
	}
	return &Rooms[idx], nil
}

func RemoveRoomById(id int) error {
	idx := slices.IndexFunc(Rooms, func(room Room) bool { return room.Id == id })
	if idx == -1 {
		return errors.New("Unknown room id")
	}
	// Remove element by index
	Rooms = append(Rooms[:idx], Rooms[idx+1:]...)
	// Todo Remove bookings for that room
	return nil
}

func AddRoom(title string) (*Room, error) {
	newRoom := Room{rand.Int(), title}
	Rooms = append(Rooms, newRoom)
	return &newRoom, nil
}
