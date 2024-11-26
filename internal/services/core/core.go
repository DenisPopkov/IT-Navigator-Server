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

type PoetProvider interface {
	Poets(ctx context.Context, userId int64) ([]models.Poet, error)
}

type ArticleProvider interface {
	Articles(ctx context.Context, userId int64) ([]models.Article, error)
}

type AuthorProvider interface {
	Authors(ctx context.Context, userId int64) ([]models.Author, error)
}

type Core struct {
	log             *slog.Logger
	userProvider    UserProvider
	poetProvider    PoetProvider
	articleProvider ArticleProvider
	authorProvider  AuthorProvider
	tokenTTL        time.Duration
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	poetProvider PoetProvider,
	articleProvider ArticleProvider,
	authorProvider AuthorProvider,
	tokenTTL time.Duration,
) *Core {
	return &Core{
		log:             log,
		userProvider:    userProvider,
		poetProvider:    poetProvider,
		articleProvider: articleProvider,
		authorProvider:  authorProvider,
		tokenTTL:        tokenTTL,
	}
}

func (c *Core) GetAuthorHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetAuthorHandler"

	uid, ok := r.Context().Value("uid").(int64)
	if !ok {
		http.Error(w, "UID not found in context", http.StatusInternalServerError)
		return
	}

	authors, err := c.authorProvider.Authors(r.Context(), uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}
}

func (c *Core) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetArticlesHandler"

	uid, ok := r.Context().Value("uid").(int64)
	if !ok {
		http.Error(w, "UID not found in context", http.StatusInternalServerError)
		return
	}

	articles, err := c.articleProvider.Articles(r.Context(), uid)
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

func (c *Core) GetPoetsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "core.GetPoetsHandler"

	uid, ok := r.Context().Value("uid").(int64)
	if !ok {
		http.Error(w, "UID not found in context", http.StatusInternalServerError)
		return
	}

	poets, err := c.poetProvider.Poets(r.Context(), uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(poets); err != nil {
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
