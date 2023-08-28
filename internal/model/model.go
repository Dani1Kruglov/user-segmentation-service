package model

type User struct {
	Id   int64
	Name string
}

type Segment struct {
	Id    int64
	Title string
}

type UserSegments struct {
	UserId         int64    `json:"user_id"`
	SegmentsTitles []string `json:"segment_titles"`
}

type AllUsersSegments struct {
	UserName      string `json:"user_names"`
	SegmentTitles string `json:"segment_titles"`
	CreatedAt     string `json:"created_at"`
	DeletedAt     string `json:"deleted_at"`
}
