package db

import (
	"database/sql"
	"strings"
	"os"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	"github.com/karak1974/flag_submit_system/types"
)

type Database struct {
	db *sql.DB
}

//connect to mysql database
func ConnectDB() (*Database, error) {
	//connect to the database
	db, _ := func() (*sql.DB, error) {
		dbUser := "root"
		dbPass := os.Getenv("DATABASE_PASSWORD")
		dbName := "fss"
		return sql.Open("mysql", dbUser+":"+dbPass+"@(FSS-database:3306)/"+dbName)
	}()

	//create db stuct
	DbHandler := &Database{
		db: db,
	}
	//return error if can't ping database
	err := DbHandler.db.Ping()

	return DbHandler, err
}

func (h *Database) Close() {
	h.db.Close()
}

// get flag [id, flag, score]
func (h *Database) GetFlag(f string) (*types.Flag, error) {
	flag := &types.Flag{}
	err := h.db.QueryRow("SELECT id, flag, score FROM flags WHERE flag=?", f).
		Scan(&flag.Id, &flag.Flag, &flag.Score)
	return flag, err
}

// get all flags in an array
func (h *Database) GetAllFlags() (flags []*types.Flag, err error) {
	rows, err := h.db.Query("SELECT id, flag, score FROM flags")
	if err != nil {
		return
	}
	for rows.Next() {
		flag := &types.Flag{}
		rows.Scan(&flag.Id, &flag.Flag, &flag.Score)
		flags = append(flags, flag)
	}
	return
}

// get user's flags by the token
// todo mysql
func (h *Database) GetUserFlags(username string) (flags []*types.Flag, err error) {
	var flaglist string
	err = h.db.QueryRow("SELECT flags FROM users WHERE username = ?", username).Scan(&flaglist)
	if err != nil {
		return
	}

	for _, id := range strings.Split(flaglist, ",") {
		flag := &types.Flag{}
		h.db.QueryRow("SELECT id, flag, score FROM flags WHERE id=?", id).Scan(&flag.Id, &flag.Flag, &flag.Score)
		flags = append(flags, flag)
	}
	return
}

// Add new flag to the DB
func (h *Database) AddFlag(flag string, score int) error {
	_, err := h.db.Exec("INSERT INTO flags (flag, score) VALUES (?, ?)", flag, score)
	return err
}

// get user [name, token, flags[]]
func (h *Database) GetUser(token string) (user types.User, err error) {
	err = h.db.QueryRow("SELECT username, token FROM users WHERE token = ?", token).
		Scan(&user.Username, &user.Token)
	return
}

// check if user exist
func (h *Database) CheckUserExist(username string) bool {
	var exist string
	h.db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&exist)
	return exist == ""
}

// check if token exist
func (h *Database) CheckToken(token string) bool {
	err := h.db.QueryRow("SELECT token FROM users WHERE token = ?", token)
	return err == nil
}

// Update user's flags and score
func (h *Database) UpdateUsersFlags(token string, flag *types.Flag) error {
	_, err := h.db.Query("UPDATE users SET flags = CONCAT_WS(',', flags, ?) WHERE token = ?", flag.Id, token)
	return err
}

//create user
func (h *Database) AddUser(username, token string) error {
	_, err := h.db.Query("INSERT INTO users(username, token) VALUES (?, ?)", username, token)
	return err
}

func (h *Database) GetScoreboard() (scoreboard []*types.Scoreboard, err error) {
	//Query the scoreboard
	res, err := h.db.Query("SELECT username FROM users")
	if err != nil {
		return
	}
	// append username and flags
	for res.Next() {
		user := &types.Scoreboard{}
		res.Scan(&user.Username)
		flags, _ := h.GetUserFlags(user.Username)
		//sum points of flags
		for _, flag := range flags{
			user.Score += flag.Score
		}
		//append to the scoreboard
		scoreboard = append(scoreboard, user)	
	}
	compareScore := func(i, j int) bool {
		return scoreboard[i].Score > scoreboard[j].Score
	 }
	sort.Slice(scoreboard, compareScore)
	return
}
