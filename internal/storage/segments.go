package storage

import (
	"avito-task/internal/model"
	"database/sql"
	"encoding/json"
	"log"
)

type SegmentJSON struct {
	Id    int64  `json:"segment_id"`
	Title string `json:"segment_title"`
}

type SegmentsPostgresStorage struct {
	Db *sql.DB
}

func (s *SegmentsPostgresStorage) GetAllSegments(conn *sql.DB) ([]byte, error) {

	rows, err := conn.Query(`SELECT * FROM segments`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var allSegments []model.Segment
	for rows.Next() {
		var segment model.Segment

		if err := rows.Scan(&segment.Id, &segment.Title); err != nil {
			log.Fatal(err)
		}
		allSegments = append(allSegments, segment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	segmentsJSON, err := json.Marshal(allSegments)
	if err != nil {
		return nil, err
	}

	return segmentsJSON, nil
}

func (s *SegmentsPostgresStorage) StoreSegment(segmentJSON []byte, conn *sql.DB) error {
	var segment SegmentJSON
	err := json.Unmarshal(segmentJSON, &segment)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`INSERT INTO segments (title) SELECT $1 WHERE NOT EXISTS (SELECT * FROM segments WHERE title = $1::varchar(255))`, segment.Title)

	if err != nil {
		return err
	}
	return nil
}

func (s *SegmentsPostgresStorage) UpdateSegment(segmentJSON []byte, conn *sql.DB) error {
	var segment SegmentJSON
	err := json.Unmarshal(segmentJSON, &segment)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`UPDATE segments SET title = $2 WHERE id = $1`, segment.Id, segment.Title)
	if err != nil {
		return err
	}
	return nil
}

func (s *SegmentsPostgresStorage) DeleteSegment(segmentTitleJSON []byte, conn *sql.DB) error {
	var segment SegmentJSON
	err := json.Unmarshal(segmentTitleJSON, &segment)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`DELETE FROM segments WHERE title = $1`, segment.Title)
	if err != nil {
		return err
	}
	return nil
}
