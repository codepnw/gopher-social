package feed

import "database/sql"

func InitFeedDomain(db *sql.DB) FeedHandler {
	repo := NewFeedRepository(db)
	uc := NewFeedUsecase(repo)
	hdl := NewFeedHandler(uc)

	return hdl
}