package entity

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Comments []Comment `json:"comments"`
}

type CreatePostPayload struct {
	Title     string   `json:"title" binding:"required,max=100"`
	Content   string   `json:"content" binding:"required,max=300"`
	Tags      []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title     *string   `json:"title" binding:"omitempty,max=100"`
	Content   *string   `json:"content" binding:"omitempty,max=300"`
}