package main

import (
	"avito-task/internal/database"
	"avito-task/internal/storage"
	"log"
)

func main() {
	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	//users := WorkWithUsers{Connect: &storage.UsersPostgresStorage{Db: db}}
	/*
			работа с пользователями
		    Данные для примера:
		    userName := "user"
			userId = 1
			var user = model.User{Id: userId, Name: userName}
	*/
	//segments := WorkWithSegments{Connect: &storage.SegmentsPostgresStorage{Db: db}}
	/*
		работа с сегментами
		Данные для примера:
		dataJSON := []byte(`{"segment_id": 13, "segment_title": "AVITO_VOICE_MESSAGES" }`)
	*/
	usersSegments := WorkWithUserSegments{Connect: &storage.UserSegmentsPostgresStorage{Db: db}}
	/*
				работа с сегментами конкретных пользователей
				Данные для примера:
				dataJSON := []byte(`{ "user_id": 2, "segment_titles": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}`)
				Для доп задания 3: dataJSON := []byte(`{"segment_title": "AVITO_VOICE_MESSAGES"}`)
				percentOfUsers := 50.0
			    Для доп задания 2: dataJSON := []byte(`{ "user_id": 2, "segment_titles": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}`)
		        duration := time.Duration(10) * time.Second
		        Для доп задания 1 никакие данные не нужны
	*/
	err = usersSegments.Connect.GetUsersSegmentsInCSVFile(db)
	if err != nil {
		log.Fatalln(err)
	}
}

type WorkWithUsers struct {
	Connect *storage.UsersPostgresStorage
}

type WorkWithSegments struct {
	Connect *storage.SegmentsPostgresStorage
}

type WorkWithUserSegments struct {
	Connect *storage.UserSegmentsPostgresStorage
}
