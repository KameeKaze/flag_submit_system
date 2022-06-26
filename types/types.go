package types

type Flag struct{
	Id    int
	Flag  string
	Score int
}

type User struct{
	Id       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Flags   []*Flag
}

type Scoreboard struct{
	Username string `json:"username`
	Score    int    `json:"score"`
}