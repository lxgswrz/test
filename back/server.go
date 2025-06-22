package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"back/db"

	"github.com/spf13/viper"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/comment/get", corsMiddleware(getCommentsHandler))
	mux.HandleFunc("/comment/add", corsMiddleware(addCommentHandler))
	mux.HandleFunc("/comment/delete", corsMiddleware(deleteCommentHandler))
	return mux
}

// CORS 中间件
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origins := viper.GetStringSlice("cors_origins")
		requestOrigin := r.Header.Get("Origin")
		allowOrigin := ""

		for _, origin := range origins {
			if origin == requestOrigin || origin == "*" {
				allowOrigin = origin
				break
			}
		}

		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	size, _ := strconv.Atoi(query.Get("size"))
	if size == 0 {
		size = 10
	}

	comments, total, err := db.GetComments(page, size)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "获取评论失败: " + err.Error(),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"total":    total,
			"comments": comments,
		},
	})
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无法读取请求体: " + err.Error(),
		})
		return
	}
	defer r.Body.Close()

	var comment db.Comment
	if err := json.Unmarshal(body, &comment); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "JSON解析失败: " + err.Error(),
		})
		return
	}

	if comment.Name == "" || comment.Content == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "名称和内容不能为空",
		})
		return
	}

	newComment, err := db.AddComment(comment.Name, comment.Content)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": newComment,
	})
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	id, err := strconv.Atoi(query.Get("id"))
	if err != nil || id == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的评论ID: " + query.Get("id"),
		})
		return
	}

	err = db.DeleteComment(id)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"code": 0,
	})
}
