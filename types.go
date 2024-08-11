package main

type OEmbed struct {
	AuthorUniqueID string `json:"author_unique_id"`
	Title          string `json:"title"`
	AuthorName     string `json:"author_name"`
	ThumbnailURL   string `json:"thumbnail_url"`
	HTML           string `json:"html"`
}

type User struct {
	ID        string `json:"uid"`
	Username  string `json:"unique_id"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"signature"`
	Region    string `json:"region"`
	AvatarURI string `json:"avatar_uri"`
	Avatar    string `json:"avatar"`
}

type Comment struct {
	User           User    `json:"user"`
	ID             string  `json:"cid"`
	Created        int64   `json:"create_time"`
	Likes          int     `json:"digg_count"`
	ReplyCount     int     `json:"reply_comment_total"`
	Text           string  `json:"text"`
	LikedByCreator bool    `json:"is_author_digged"`
	Replies        []Reply `json:"replies,omitempty"`
}

type Reply struct {
	User           User   `json:"user"`
	ID             string `json:"cid"`
	Created        int64  `json:"create_time"`
	Likes          int    `json:"digg_count"`
	Text           string `json:"text"`
	LikedByCreator bool   `json:"is_author_digged"`
	ParentComment  string `json:"parent_comment"`
	ParentUser     string `json:"parent_user,omitempty"`
}

type GetCommentResponse struct {
	Comments []Comment `json:"comments"`
	HasMore  int       `json:"has_more"`
}

type GetReplyResponse struct {
	Comments []Reply `json:"comments"`
}
