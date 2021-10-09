package main

//User attributes
type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

//Post attributes
type Posts struct {
	Caption     string `json:"caption" bson:"caption"`
	Url         string `json:"url" bson:"url"`
	CurrentTime string `json:"currentTime" bson:"currentTime"`
	UserID      string `json:"userID" bson:"userID"`
}
