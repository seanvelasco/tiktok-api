package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func getVideo(postID string) (io.ReadCloser, error) {
	// This is a temporary implementation; only using tikwm as a placeholder for the service

	u := fmt.Sprintf("https://www.tikwm.com/video/media/wmplay/%s.mp4", postID)

	res, err := http.Get(u)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func getPost(userID string, postID string) (OEmbed, error) {
	u, err := url.Parse("https://www.tiktok.com/oembed")
	if err != nil {
		return OEmbed{}, err
	}

	q := u.Query()
	q.Set("url", fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", userID, postID))
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())

	if err != nil {
		return OEmbed{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return OEmbed{}, errors.New(res.Status)
	}

	var oembed OEmbed

	if err := json.NewDecoder(res.Body).Decode(&oembed); err != nil {
		return OEmbed{}, err
	}

	return oembed, nil
}

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
		return GetCommentResponse{}, errors.New("Failed to build request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:129.0) Gecko/20100101 Firefox/129.0")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return GetCommentResponse{}, errors.New("Failed to get comments")
	}

	defer res.Body.Close()

	var comments GetCommentResponse

	if err := json.NewDecoder(res.Body).Decode(&comments); err != nil {
		return GetCommentResponse{}, errors.New("Failed to decode JSON")
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
		body, _ := io.ReadAll(res.Body)
		fmt.Println("Raw response:", string(body))
		return GetReplyResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return replies, nil
}

func getThread(postID string) ([]Comment, error) {
	var comments []Comment
	offset := 0
	limit := 50

	for {
		post, err := getComments(postID, limit, offset)

		if err != nil {
			break
		}

		if post.Comments == nil {
			break
		}

		comments = append(comments, post.Comments...)

		if post.HasMore != 1 {
			break
		}
		offset += limit
	}

	var wg sync.WaitGroup

	for i := range comments {
		if comments[i].ReplyCount > 0 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				replies, err := getReplies(postID, comments[i].ID, comments[i].ReplyCount)
				fmt.Println("REPLY", i)
				if err != nil {
					log.Fatal("Failed to get replies", err)
					return
				}
				if replies.Comments != nil {
					comments[i].Replies = replies.Comments
				}
			}(i)
		}
	}

	wg.Wait()

	return comments, nil
}
