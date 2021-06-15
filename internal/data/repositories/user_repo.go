package repositories

import (
	entity "chatapp/internal/data/entities"
	"chatapp/internal/lib/util"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type UserRepository struct {
	Db   *sql.DB
	base BaseRepository
}

// Create a new instance
func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		Db: db,
		// TODO Anti-pattern (tight-coupling) -- need to pass it as a paramater instead!
		base: BaseRepository{
			tableName: "user",
		},
	}
}

// Add a new user
func (this *UserRepository) Add(user *entity.User) bool {
	stmt, err := this.Db.Prepare("INSERT INTO \"" + this.base.tableName + "\"(id, name, username, password) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return false
	}
	user.Id = uuid.New().String()
	user.Password, _ = util.GenPasswordHash(user.Password)
	_, err = stmt.Exec(user.Id, user.Name.String(), user.Username.String(), user.Password)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// Find a user by its Id (Primary Key)
func (this *UserRepository) FindById(id string) *entity.User {
	var user entity.User
	err := this.Db.QueryRow(
		"SELECT id, name, username, password FROM \""+this.base.tableName+"\" where id = ? LIMIT 1", id,
	).Scan(&user.Id, &user.Name, &user.Username, &user.Password)
	// check error
	switch err {
	case sql.ErrNoRows:
		fmt.Println("There is no retrieved rows, dummy!")
	case nil:
		fmt.Println("Row is Nil")
	default:
		panic(err)
	}
	return &user
}

// Find a user by its username
func (this UserRepository) FindByUsername(username string) *entity.User {
	var user entity.User
	err := this.Db.QueryRow(
		"SELECT id, name, username, password FROM \""+this.base.tableName+"\" where username = ? LIMIT 1", username,
	).Scan(&user.Id, &user.Name, &user.Username, &user.Password)
	// check error
	switch err {
	case sql.ErrNoRows:
		fmt.Println("There is no retrieved rows, dummy!")
	case nil:
		fmt.Println("Row is Nil")
	default:
		panic(err)
	}
	return &user
}

// Returns all users, otherwise an error
func (this *UserRepository) FindAll() []entity.User {
	rows, err := this.Db.Query("SELECT id, name, username, password FROM \"" + this.base.tableName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		var user entity.User
		rows.Scan(&user.Id, &user.Name, &user.Username, &user.Password)
		users = append(users, user)
	}
	return users
}

// Delete a user by a given id
func (this *UserRepository) Delete(id string) bool {
	stmt, err := this.Db.Prepare("DELETE FROM \"" + this.base.tableName + "\" WHERE id = ?")
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
