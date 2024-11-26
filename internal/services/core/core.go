package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sso/internal/domain/models"
	"time"
)

type UserProvider interface {
	DeleteUser(ctx context.Context, userId int64) error
	GetUser(ctx context.Context, userId int64) (models.UserData, error)
}

type CourseProvider interface {
	Courses(ctx context.Context) ([]models.Course, error)
}

type ArticleProvider interface {
	Articles(ctx context.Context) ([]models.Article, error)
}

type FeedProvider interface {
	Feeds(ctx context.Context) ([]models.Feed, error)
}

type Core struct {
	log             *slog.Logger
	userProvider    UserProvider
	courseProvider  CourseProvider
	articleProvider ArticleProvider
	feedProvider    FeedProvider
	tokenTTL        time.Duration
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	courseProvider CourseProvider,
	articleProvider ArticleProvider,
	feedProvider FeedProvider,
	tokenTTL time.Duration,
) *Core {
	return &Core{
		log:             log,
		userProvider:    userProvider,
		courseProvider:  courseProvider,
		articleProvider: articleProvider,
		feedProvider:    feedProvider,
		tokenTTL:        tokenTTL,
	}
}

func (c *Core) GetFeedHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetFeedHandler"

	feeds, err := c.feedProvider.Feeds(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}
}

func (c *Core) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetArticlesHandler"

	articles, err := c.articleProvider.Articles(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}
}

func (c *Core) GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetCoursesHandler"

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	courses, err := c.courseProvider.Courses(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}
}

func (c *Core) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.DeleteUserHandler"

	uid, ok := r.Context().Value("uid").(int64)
	if !ok {
		http.Error(w, "UID not found in context", http.StatusInternalServerError)
		return
	}

	err := c.userProvider.DeleteUser(r.Context(), uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Core) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetUserHandler"

	uid, ok := r.Context().Value("uid").(int64)
	if !ok {
		http.Error(w, "UID not found in context", http.StatusInternalServerError)
		return
	}

	_, err := c.userProvider.GetUser(r.Context(), uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
