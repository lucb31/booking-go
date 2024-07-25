package booking

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Room struct {
	Id    int64
	Title string
}

type RoomScan struct {
	Id    int64
	Title sql.NullString
}

func RoomFromScan(s *RoomScan) Room {
	return Room{Id: s.Id, Title: s.Title.String}
}

type RoomsRepository interface {
	Migrate() error
	SeedTestData() error
	Create(room Room) (*Room, error)
	GetAll() ([]*Room, error)
	Delete(id int64) error
	GetById(id int64) (*Room, error)
}

type RoomsRepositorySQLite struct {
	db *sqlx.DB
}

func NewRoomsRepositorySQLite(db *sqlx.DB) *RoomsRepositorySQLite {
	return &RoomsRepositorySQLite{db}
}

func (r *RoomsRepositorySQLite) Migrate() error {
	query := `
CREATE TABLE IF NOT EXISTS room (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL
); `
	_, err := r.db.Exec(query)
	return err
}

func (r *RoomsRepositorySQLite) SeedTestData() error {
	query := ` INSERT INTO room (title) VALUES ("Test room"); `
	_, err := r.db.Exec(query)
	return err
}

func (r *RoomsRepositorySQLite) Create(room Room) (*Room, error) {
	query := ` INSERT INTO room ( title ) VALUES (?); `
	rows, err := r.db.Exec(query, room.Title)
	if room.Id, err = rows.LastInsertId(); err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomsRepositorySQLite) GetAll() ([]*Room, error) {
	query := `
	SELECT
		id
		title
	FROM
		room;
`
	rows, err := r.db.Queryx(query)
	rooms := []*Room{}
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var scan RoomScan
		err = rows.StructScan(&scan)
		if err := rows.StructScan(&scan); err != nil {
			return rooms, err
		}
		room := RoomFromScan(&scan)
		rooms = append(rooms, &room)
	}
	return rooms, nil
}

func (r *RoomsRepositorySQLite) Delete(id int64) error {
	return fmt.Errorf("Missing implementation: Delete room")
}

func (r *RoomsRepositorySQLite) GetById(id int64) (*Room, error) {
	return nil, fmt.Errorf("Missing implementation: Get room by ID")
}
