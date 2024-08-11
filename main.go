package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func getComments(postID string, count int, cursor int) (GetCommentResponse, error) {
	u, err := url.Parse("https://www.tiktok.com/api/comment/list/")
	if err != nil {
		return GetCommentResponse{}, err
	}

	q := u.Query()
	q.Set("aweme_id", postID)
	q.Set("count", strconv.Itoa(count))
	q.Set("cursor", strconv.Itoa(cursor))
	q.Set("os", "mac")
	q.Set("region", "US")
	q.Set("screen_height", "900")
	q.Set("screen_width", "1440")
	q.Set("X-Bogus", "DFSzsIVOxrhAN9fbtfB5EX16ZwHH")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return GetCommentResponse{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:129.0) Gecko/20100101 Firefox/129.0")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return GetCommentResponse{}, err
	}

	defer res.Body.Close()

	var comments GetCommentResponse

	if err := json.NewDecoder(res.Body).Decode(&comments); err != nil {
		return GetCommentResponse{}, err
	}

	return comments, nil
}

func getReplies(postID string, commentID string, count int) (GetReplyResponse, error) {
	u, err := url.Parse("https://www.tiktok.com/api/comment/list/reply")
	if err != nil {
		return GetReplyResponse{}, err
	}

	q := u.Query()
	q.Set("item_id", postID)
	q.Set("comment_id", commentID)
	q.Set("count", strconv.Itoa(count))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return GetReplyResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return GetReplyResponse{}, err
	}

	defer res.Body.Close()

	var replies GetReplyResponse

	if err := json.NewDecoder(res.Body).Decode(&replies); err != nil {
		return GetReplyResponse{}, err
	}

	return replies, nil
}

func getThread(postID string) []Comment {
	var comments []Comment
	offset := 0
	limit := 50

	for {
		post, err := getComments(postID, limit, offset)
		if err != nil {
			log.Fatal(err)
		}
		if post.Comments == nil || post.HasMore != 1 {
			break
		}
		comments = append(comments, post.Comments...)
		offset += limit
	}

	var wg sync.WaitGroup
	for i, comment := range comments {
		wg.Add(1)
		if comment.ReplyCount > 0 {
			go func(i int, comment Comment) {
				defer wg.Done()
				replies, err := getReplies(postID, comment.ID, comment.ReplyCount)
				if err != nil {
					log.Fatal(err)
				}
				comments[i].Replies = replies.Comments
			}(i, comment)
		}
	}
	wg.Wait()

	return comments
}

func main() {

}
