package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/yuin/goldmark"
)

var threadURIPattern = regexp.MustCompile(`at://[^/]+/app.bsky.feed.post/[a-z0-9]+`)

type Server struct {
	cfg         Configuration
	mux         *http.ServeMux
	client      *xrpc.Client
	threadCache *expirable.LRU[string, *bsky.FeedGetPostThread_Output]
	logger      *slog.Logger
}

type Favorite struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"displayName,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
}

func New(cfg Configuration) *Server {
	threadCache := expirable.NewLRU[string, *bsky.FeedGetPostThread_Output](10, nil, time.Minute*1)

	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(nil)
	}

	srv := &Server{
		cfg:         cfg,
		client:      cfg.Client,
		threadCache: threadCache,
		logger:      logger,
	}

	if srv.client == nil {
		srv.client = &xrpc.Client{
			Host: "https://public.api.bsky.app",
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/favorited_by", srv.handleGetFavoritedBy)
	mux.HandleFunc("/api/v1/thread", srv.handleGetThread)
	srv.mux = mux
	return srv
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}

func (srv *Server) isAllowedAccount(threadURI string) bool {
	for _, did := range srv.cfg.AllowedRootAccountDIDs {
		if strings.HasPrefix(threadURI, fmt.Sprintf("at://%s/", did)) {
			return true
		}
	}
	for _, handle := range srv.cfg.AllowedRootAccounts {
		if strings.HasPrefix(threadURI, fmt.Sprintf("at://%s/", handle)) {
			return true
		}
	}
	return false
}

func (srv *Server) handleGetFavoritedBy(w http.ResponseWriter, r *http.Request) {
	var threadURI string
	threadURI = r.URL.Query().Get("uri")
	if !srv.isAllowedAccount(threadURI) {
		http.Error(w, "not allowed root account", http.StatusBadRequest)
		return
	}
	if !threadURIPattern.MatchString(threadURI) {
		http.Error(w, "invalid uri", http.StatusBadRequest)
		return
	}

	var cursor string
	result := make([]Favorite, 0, 10)
	for {
		output, err := bsky.FeedGetLikes(r.Context(), srv.client, "", cursor, 10, threadURI)
		if err != nil {
			http.Error(w, "backend request failed", http.StatusInternalServerError)
			return
		}
		for _, like := range output.Likes {
			fav := Favorite{}
			if like.Actor == nil {
				continue
			}
			if like.Actor.Avatar != nil {
				fav.Avatar = *like.Actor.Avatar
			}
			fav.DID = like.Actor.Did
			fav.Handle = like.Actor.Handle
			if like.Actor.DisplayName != nil {
				fav.DisplayName = *like.Actor.DisplayName
			}
			result = append(result, fav)
		}
		if output.Cursor == nil {
			break
		}
		cursor = *output.Cursor
	}

	w.Header().Add("Content-Type", "text/json")
	json.NewEncoder(w).Encode(result)

}

func (srv *Server) handleGetThread(w http.ResponseWriter, r *http.Request) {
	var threadURI string
	threadURI = r.URL.Query().Get("uri")
	if !srv.isAllowedAccount(threadURI) {
		http.Error(w, "not allowed root account", http.StatusBadRequest)
		return
	}
	if !threadURIPattern.MatchString(threadURI) {
		http.Error(w, "invalid uri", http.StatusBadRequest)
		return
	}

	// TODO: normalize thread URI
	thread, err := srv.getCachedThread(fmt.Sprintf("thread:%s", threadURI), func() (*bsky.FeedGetPostThread_Output, error) {
		return bsky.FeedGetPostThread(r.Context(), srv.client, 5, 0, threadURI)
	})
	if err != nil {
		http.Error(w, "backend request failed", http.StatusInternalServerError)
		return
	}
	output, err := RenderThread(thread.Thread.FeedDefs_ThreadViewPost, 0)
	if err != nil {
		http.Error(w, "rendering failed", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, string(output))

}

func (srv *Server) getCachedThread(key string, fn func() (*bsky.FeedGetPostThread_Output, error)) (*bsky.FeedGetPostThread_Output, error) {
	result, hit := srv.threadCache.Get(key)
	if hit {
		srv.logger.Debug("thread cache hit", "key", key)
		return result, nil
	}
	srv.logger.Debug("thread cache miss", "key", key)
	result, err := fn()
	if err != nil {
		return nil, err
	}
	srv.threadCache.Add(key, result)
	return result, err
}

type ThreadRenderingData struct {
	URI                string
	PostID             string
	Text               template.HTML
	Replies            []template.HTML
	CreatedAt          string
	Level              int
	AuthorHandle       string
	AuthorDID          string
	AuthorAvatar       string
	ExternalEmbedURI   string
	ExternalEmbedTitle string
}

func RenderThread(thread *bsky.FeedDefs_ThreadViewPost, level int) (template.HTML, error) {
	renderedReplies := make([]template.HTML, 0, len(thread.Replies))
	for _, reply := range thread.Replies {
		r, err := RenderThread(reply.FeedDefs_ThreadViewPost, level+1)
		if err != nil {
			return "", err
		}
		renderedReplies = append(renderedReplies, r)
	}
	post := thread.Post.Record.Val.(*bsky.FeedPost)
	markdown := goldmark.New()
	var textOutput bytes.Buffer
	if err := markdown.Convert([]byte(post.Text), &textOutput); err != nil {
		return "", err
	}

	data := ThreadRenderingData{
		URI:          thread.Post.Uri,
		PostID:       "",
		Text:         template.HTML(textOutput.String()),
		Replies:      renderedReplies,
		Level:        level,
		CreatedAt:    post.CreatedAt,
		AuthorHandle: thread.Post.Author.Handle,
		AuthorDID:    thread.Post.Author.Did,
	}
	if thread.Post.Author.Avatar != nil {
		data.AuthorAvatar = strings.Replace(*thread.Post.Author.Avatar, "avatar", "avatar_thumbnail", -1)
	}
	parsedURI, err := util.ParseAtUri(thread.Post.Uri)
	if err != nil {
		return "", err
	}
	data.PostID = parsedURI.Rkey
	if post.Embed != nil && post.Embed.EmbedExternal != nil {
		if post.Embed.EmbedExternal.External.Uri != "" {
			data.ExternalEmbedURI = post.Embed.EmbedExternal.External.Uri
		}
		if post.Embed.EmbedExternal.External.Title != "" {
			data.ExternalEmbedURI = post.Embed.EmbedExternal.External.Title
		}
	}
	tmpl := template.Must(template.New("thread").Parse(postTmpl))
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

var postTmpl = `
<div class="bsky-feed-thread bsky-feed-thread--lvl{{ .Level }}">
	{{ if gt .Level 0 }}
	<div class="bsky-feed-post">
	<div class="bsky-feed-post__avatar">
		<a href="https://bsky.app/profile/{{ .AuthorHandle }}" class="bsky-author-handle"><img src="{{ .AuthorAvatar }}" /></a>
	</div>
	<div class="bsky-feed-post__content">
		{{ .Text }}
		{{ if .ExternalEmbedURI }}
		<div class="bsky-feed-post__embed">
			<a href="{{ .ExternalEmbedURI }}" rel="no-follow">{{ .ExternalEmbedURI }}</a>
		</div>
		{{ end }}
		<a class="bsky-feed-post__date" href="https://bsky.app/profile/{{ .AuthorHandle }}/post/{{ .PostID }}">{{ .CreatedAt }}</a>
	</div>
	</div>
	{{ end }}
	{{ if gt (len .Replies) 0 }}
	<div class="bsky-feed-post__replies">
		{{ range .Replies }}
		{{ . }}
		{{ end }}
	</div>
	{{ end }}
</div>
`
