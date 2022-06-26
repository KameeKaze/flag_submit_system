package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/karak1974/flag_submit_system/db"
	"github.com/karak1974/flag_submit_system/types"
	"github.com/karak1974/flag_submit_system/utils"
)

const (
	PORT = ":8000"
)


func RoutesHandler() {
	// init mux handler
	r := mux.NewRouter()

	// define routes
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/register", RegisterHandler).Methods("GET")
	r.HandleFunc("/register", RegisterUser).Methods("POST")
	r.HandleFunc("/submit", SubmitHandler).Methods("GET")
	r.HandleFunc("/submit", SubmitFlag).Methods("POST")
	r.HandleFunc("/scoreboard", GetScoreboard).Methods("POST")
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("./html/")))

	//start http server
	fmt.Println("Running on http://localhost" + PORT)
	http.ListenAndServe(PORT, r)
}

// Three handlers for the static files
func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "html/index.html")
}

func RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "html/register.html")
}

func SubmitHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "html/submit.html")
}

// Register a new user if it's avaliable
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//create user  from post body
	user := &types.User{} //username, token
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error("Error reading body from post request")
	}
	json.Unmarshal(body, &user)
	user.Token = utils.GenerateToken()


	// connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer database.Close()

	if user.Username == "" {
		fmt.Fprintf(w, utils.MsgParser("empty"))
	// check if user exist and not empty
	}else if database.CheckUserExist(user.Username) {
		// create user
		err = database.AddUser(user.Username, user.Token)
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		//Return the creds to the user and log in the terminal
		fmt.Fprintf(w, utils.MsgParser(user.Username + "\nToken "+user.Token))
		utils.Logger.Info("New account | Username: " + user.Username)
	}else { // user already exist
		fmt.Fprintf(w, utils.MsgParser(user.Username+" taken"))
	}
}

// Update the user with the flag's score
func SubmitFlag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get post data from body
	var postdata map[string]string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error("Error reading body from post request")
	}
	json.Unmarshal(body, &postdata) 

	// connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer database.Close()

	// get user and check if token is valid
	user, err := database.GetUser(postdata["token"])
	if err != nil{
		fmt.Fprintf(w, utils.MsgParser("Invalid token"))
		return
	} 
	user.Flags, err = database.GetUserFlags(user.Username)
	flag, err := database.GetFlag(postdata["flag"])

	//check if flag is valid
	if flag.Flag != "" {
		//check if flag already submited
		for _, f := range user.Flags {
			if flag.Flag == f.Flag {
				fmt.Fprintf(w, utils.MsgParser("You already submited this flag"))
				return
			}
		}
		// update the users solved flags
		err = database.UpdateUsersFlags(user.Token, flag)
		fmt.Fprintf(w, utils.MsgParser("Flag successfuly submited"))
	} else {
		fmt.Fprintf(w, utils.MsgParser("Invalid flag"))
	}

}

// Get scoreboard for the main page
func GetScoreboard(w http.ResponseWriter, r *http.Request) {
	var scoreboard []types.Scoreboard
	// connect to database
	database, err := db.ConnectDB()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer database.Close()

	score, err := database.GetScoreboard()
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	for _, user := range score {
		scoreboard = append(scoreboard, types.Scoreboard{Username: user.Username, Score: user.Score})
	}

	j, _ := json.Marshal(scoreboard)
	fmt.Fprintf(w, string(j))
}

