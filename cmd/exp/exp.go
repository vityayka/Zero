package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Sslmode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.Database, cfg.Sslmode)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "vityayka",
		Password: "azazaz1488",
		Database: "zero",
		Sslmode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT NOT NULL,
			description TEXT
		);

		CREATE TABLE IF NOT EXISTS tweets (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			parent_id INT,
			content TEXT
		);

		CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
			tweet_id INT NOT NULL,
			user_id INT NOT NULL
		);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created")

	// var uname string = "asdfasfd"
	// var uemail string = "asdfas12367@email.com"
	// result := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", uname, uemail)

	// var id int
	// result.Scan(&id)

	// fmt.Println(id)

	// if err != nil {
	// 	panic(err)
	// }

	// userID := 6
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Generated order #%d", i)
	// 	_, err = db.Exec(`
	// 		INSERT INTO orders (user_id, amount, description)
	// 		VALUES ($1, $2, $3)
	// 	`, userID, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Orders generated")
	// }

	// for i := 1; i <= 5; i++ {
	// 	desc := fmt.Sprintf("Generated tweet #%d", i)
	// 	_, err = db.Exec(`
	// 		INSERT INTO tweets (user_id, content)
	// 		VALUES ($1, $2)
	// 	`, userID, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// 	fmt.Println("Tweets generated")

	tweetId := 1
	for userID := 5; userID <= 10; userID++ {
		_, err = db.Exec(`
			INSERT INTO likes (user_id, tweet_id)
			VALUES ($1, $2)
		`, userID, tweetId)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Likes generated")

	type Order struct {
		Id          int
		UserID      int
		Amount      int
		Description string
	}

	type Tweet struct {
		Id       int
		UserID   int
		ParentId int
		Content  string
	}

	type TweetLike struct {
		Id      int
		UserID  int
		TweetID int
	}

	// var orders []Order

	// rows, err := db.Query("select id, amount, description from orders where user_id = $1", userID)

	// for rows.Next() {
	// 	var order Order
	// 	order.UserID = userID
	// 	rows.Scan(&order.Id, &order.Amount, &order.Description)
	// 	orders = append(orders, order)
	// }
	// var tweets []Tweet

	// rows, _ := db.Query("select id, content from tweets where user_id = $1", userID)

	// for rows.Next() {
	// 	var tweet Tweet
	// 	tweet.UserID = userID
	// 	rows.Scan(&tweet.Id, &tweet.Content)
	// 	tweets = append(tweets, tweet)
	// }

	// fmt.Printf("Tweets: %+v", tweets)

	var likes []TweetLike

	rows, _ := db.Query("select id, user_id from likes where tweet_id = $1", tweetId)

	for rows.Next() {
		var like TweetLike
		like.TweetID = tweetId
		rows.Scan(&like.Id, &like.UserID)
		likes = append(likes, like)
	}

	fmt.Printf("Likes: %+v", likes)
}
