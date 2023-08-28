package storage

import (
	"avito-task/internal/model"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type UserSegmentsJSON struct {
	UserId        int64         `json:"user_id"`
	SegmentTitles []interface{} `json:"segment_titles"`
}

type UserSegmentsPostgresStorage struct {
	Db *sql.DB
}

func (us *UserSegmentsPostgresStorage) GetUsersSegmentsInCSVFile(conn *sql.DB) error {
	rows, err := conn.Query(`SELECT u.name, s.title, created_at, deleted_at FROM users_segments INNER JOIN segments s ON s.id = users_segments.segment_id INNER JOIN users u ON u.id = users_segments.user_id`)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var allUsersSegments []model.AllUsersSegments
	for rows.Next() {
		var userSegment model.AllUsersSegments
		if err := rows.Scan(&userSegment.UserName, &userSegment.SegmentTitles, &userSegment.CreatedAt, &userSegment.DeletedAt); err != nil {
			if err.Error() == "sql: Scan error on column index 3, name \"deleted_at\": converting NULL to string is unsupported" {
				userSegment.DeletedAt = "NULL"
			} else {
				log.Fatalln(err)
			}
		}
		allUsersSegments = append(allUsersSegments, userSegment)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	err = writeToCsv(allUsersSegments)
	if err != nil {
		return err
	}
	fmt.Println("SUCCESS: The data.csv file is created, filled with data and added to the root directory")
	return nil
}

func writeToCsv(allUsersSegments []model.AllUsersSegments) error {
	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)

	w := csv.NewWriter(file)
	err = w.WriteAll(convertStructInDoubleSlice(allUsersSegments))
	if err != nil {
		return err
	}
	return nil
}

func convertStructInDoubleSlice(allUsersSegments []model.AllUsersSegments) [][]string {
	var result [][]string
	mainRow := []string{"USER", "SEGMENT", "CREATED_AT", "DELETED_AT"}
	result = append(result, mainRow)
	for _, userSegment := range allUsersSegments {
		row := []string{userSegment.UserName, userSegment.SegmentTitles, userSegment.CreatedAt, userSegment.DeletedAt}
		result = append(result, row)
	}
	return result
}

func (us *UserSegmentsPostgresStorage) GetUserSegments(dataJSON []byte, conn *sql.DB) ([]byte, error) {
	var data UserSegmentsJSON
	err := json.Unmarshal(dataJSON, &data)
	rows, err := conn.Query(`SELECT user_id, s.title FROM users_segments INNER JOIN segments s ON s.id = users_segments.segment_id WHERE user_id = $1`, data.UserId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var userSegments model.UserSegments

	for rows.Next() {
		var userId int64
		var segmentTitle string
		if err := rows.Scan(&userId, &segmentTitle); err != nil {
			log.Fatal(err)
		}
		userSegments.UserId = userId
		userSegments.SegmentsTitles = append(userSegments.SegmentsTitles, segmentTitle)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	userSegmentsJSON, err := json.Marshal(userSegments)
	if err != nil {
		return nil, err
	}

	return userSegmentsJSON, nil
}

func (us *UserSegmentsPostgresStorage) AddUserToSegments(dataJSON []byte, conn *sql.DB) error {
	var userSegments UserSegmentsJSON
	if err := json.Unmarshal(dataJSON, &userSegments); err != nil {
		return err
	}
	for i, _ := range userSegments.SegmentTitles {
		_, err := conn.Exec(`
        INSERT INTO users_segments (user_id, created_at, segment_id)
        SELECT $1, $2, id FROM segments WHERE segments.title = $3
        AND NOT EXISTS (
            SELECT 1 FROM users_segments 
            WHERE user_id = $1 
            AND segment_id = (SELECT id FROM segments WHERE segments.title = $3)
        )`, userSegments.UserId, time.Now(), userSegments.SegmentTitles[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UserSegmentsPostgresStorage) AddUserToSegmentsForWhile(dataJSON []byte, duration time.Duration, conn *sql.DB) error {
	var userSegments UserSegmentsJSON
	if err := json.Unmarshal(dataJSON, &userSegments); err != nil {
		return err
	}

	t := time.NewTicker(duration)
	defer t.Stop()

	for i, _ := range userSegments.SegmentTitles {
		_, err := conn.Exec(`
        INSERT INTO users_segments (user_id, created_at, duration, segment_id)
        SELECT $1, $2, $3, id FROM segments WHERE segments.title = $4
        AND NOT EXISTS (
            SELECT 1 FROM users_segments 
            WHERE user_id = $1 
            AND segment_id = (SELECT id FROM segments WHERE segments.title = $4)
        )`, userSegments.UserId, time.Now(), duration.String(), userSegments.SegmentTitles[i])
		if err != nil {
			return err
		}
	}

	for {
		select {
		case <-t.C:
			err := us.DeleteUserSegments(dataJSON, conn)
			if err != nil {
				log.Fatalln(err)
			}
			return nil
		}
	}
}

func (us *UserSegmentsPostgresStorage) AddSomeUsersToSegment(segmentJSON []byte, percentOfUsers float64, conn *sql.DB) error {
	var segment SegmentJSON
	if err := json.Unmarshal(segmentJSON, &segment); err != nil {
		return err
	}

	rows, err := conn.Query(`SELECT count(id) FROM users`)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	var count float64
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatalln(err)
		}
	}

	numberOfUsersByPercent := getNumberOfUsersByPercent(count, percentOfUsers)

	for i := 1; i <= numberOfUsersByPercent; i++ {
		result, err := conn.Exec(`
        INSERT INTO users_segments (user_id, created_at, segment_id)
        SELECT $1, $2, id FROM segments WHERE segments.title = $3
        AND EXISTS(SELECT 1 FROM users WHERE users.id = $1)
        AND NOT EXISTS (
            SELECT 1 FROM users_segments 
            WHERE user_id = $1 
            AND segment_id = (SELECT id FROM segments WHERE segments.title = $3)
        )`, i, time.Now(), segment.Title)
		if err != nil {
			return err
		}
		code, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if code == 0 {
			if numberOfUsersByPercent == int(count) {
				return errors.New("AddSomeUsersToSegment: there is no number of all users you need")
			} else {
				numberOfUsersByPercent++
			}
		}
	}
	return nil
}

func getNumberOfUsersByPercent(count, percentOfUsers float64) int {
	return int(math.Round((count * percentOfUsers) / 100.0))
}

func (us *UserSegmentsPostgresStorage) DeleteUserSegments(dataJSON []byte, conn *sql.DB) error {
	var userSegments UserSegmentsJSON
	if err := json.Unmarshal(dataJSON, &userSegments); err != nil {
		return err
	}
	for i, _ := range userSegments.SegmentTitles {
		_, err := conn.Exec(`DELETE FROM users_segments WHERE user_id = $1 AND segment_id = (SELECT id FROM segments WHERE segments.title = $2)`, userSegments.UserId, userSegments.SegmentTitles[i])
		if err != nil {
			return err
		}
	}
	return nil
}
