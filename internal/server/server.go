package server

import (
	"bytes"
	"context"
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
	"github.com/zerok/samara/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var tracerName = "github.com/zerok/samara/internal/server"
var tracer = otel.Tracer(tracerName)

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
	mux.HandleFunc("/api/v1/thread", srv.handleGetThread)
	srv.mux = mux
	return srv
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "ServeHTTP")
	defer span.End()
	span.SetAttributes(semconv.HTTPRequestMethodKey.String(r.Method), semconv.URLFull(r.URL.String()))
	srv.mux.ServeHTTP(w, r.WithContext(ctx))
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

func (srv *Server) handleGetThread(w http.ResponseWriter, r *http.Request) {
	var threadURI string
	threadURI = r.URL.Query().Get("uri")
	ctx, span := tracer.Start(r.Context(), "handleGetThread")
	defer span.End()
	span.SetAttributes(telemetry.ThreadURIKey.String(threadURI))
	if !srv.isAllowedAccount(threadURI) {
		span.SetStatus(codes.Ok, "not allowed account")
		http.Error(w, "not allowed root account", http.StatusBadRequest)
		return
	}
	if !threadURIPattern.MatchString(threadURI) {
		span.SetStatus(codes.Ok, "invalid uri")
		http.Error(w, "invalid uri", http.StatusBadRequest)
		return
	}

	// TODO: normalize thread URI
	thread, err := srv.getCachedThread(ctx, fmt.Sprintf("thread:%s", threadURI), func(ctx context.Context) (*bsky.FeedGetPostThread_Output, error) {
		return bsky.FeedGetPostThread(ctx, srv.client, 5, 0, threadURI)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "backend request failed")
		http.Error(w, "backend request failed", http.StatusInternalServerError)
		return
	}
	output, html, err := RenderThread(ctx, thread.Thread.FeedDefs_ThreadViewPost, 0)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "rendering failed")
		http.Error(w, "rendering failed", http.StatusInternalServerError)
		return
	}

	// Depending on HTMX vs. raw XHR, render that whole thing as single HTML or
	// as JSON
	span.SetStatus(codes.Ok, "")
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, string(html))
		return
	}
	w.Header().Add("Content-Type", "text/json")
	json.NewEncoder(w).Encode(output)

}

func (srv *Server) getCachedThread(ctx context.Context, key string, fn func(context.Context) (*bsky.FeedGetPostThread_Output, error)) (*bsky.FeedGetPostThread_Output, error) {
	ctx, span := tracer.Start(ctx, "getCachedThread")
	defer span.End()
	result, hit := srv.threadCache.Get(key)
	span.SetAttributes(telemetry.ThreadCacheHitKey.Bool(hit))
	if hit {
		srv.logger.Debug("thread cache hit", "key", key)
		return result, nil
	}
	srv.logger.Debug("thread cache miss", "key", key)
	result, err := fn(ctx)
	if err != nil {
		return nil, err
	}
	srv.threadCache.Add(key, result)
	return result, err
}

type ThreadRenderingData struct {
	URI                string                `json:"uri"`
	PostID             string                `json:"postID"`
	Text               string                `json:"text"`
	RenderedText       template.HTML         `json:"-"`
	Replies            []ThreadRenderingData `json:"replies"`
	RenderedReplies    []template.HTML       `json:"-"`
	CreatedAt          string                `json:"createdAt"`
	Level              int                   `json:"level"`
	AuthorHandle       string                `json:"authorHandle"`
	AuthorDID          string                `json:"authorDID"`
	AuthorAvatar       string                `json:"authorAvatar"`
	ExternalEmbedURI   string                `json:"externalEmbedURI"`
	ExternalEmbedTitle string                `json:"externalEmbedTitle"`
}

func RenderThread(ctx context.Context, thread *bsky.FeedDefs_ThreadViewPost, level int) (*ThreadRenderingData, template.HTML, error) {
	replies := make([]ThreadRenderingData, 0, len(thread.Replies))
	renderedReplies := make([]template.HTML, 0, len(thread.Replies))
	for _, reply := range thread.Replies {
		r, rendered, err := RenderThread(ctx, reply.FeedDefs_ThreadViewPost, level+1)
		if err != nil {
			return nil, "", err
		}
		replies = append(replies, *r)
		renderedReplies = append(renderedReplies, rendered)
	}
	post := thread.Post.Record.Val.(*bsky.FeedPost)
	markdown := goldmark.New()
	var textOutput bytes.Buffer
	if err := markdown.Convert([]byte(post.Text), &textOutput); err != nil {
		return nil, "", err
	}

	data := ThreadRenderingData{
		URI:             thread.Post.Uri,
		PostID:          "",
		RenderedText:    template.HTML(textOutput.String()),
		Text:            post.Text,
		Replies:         replies,
		RenderedReplies: renderedReplies,
		Level:           level,
		CreatedAt:       post.CreatedAt,
		AuthorHandle:    thread.Post.Author.Handle,
		AuthorDID:       thread.Post.Author.Did,
	}
	if thread.Post.Author.Avatar != nil {
		data.AuthorAvatar = strings.Replace(*thread.Post.Author.Avatar, "avatar", "avatar_thumbnail", -1)
	}
	parsedURI, err := util.ParseAtUri(thread.Post.Uri)
	if err != nil {
		return nil, "", err
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
	html, err := renderThreadToHTML(ctx, data)
	return &data, html, err
}

func renderThreadToHTML(ctx context.Context, thread ThreadRenderingData) (template.HTML, error) {
	tmpl := template.Must(template.New("thread").Parse(postTmpl))
	var output bytes.Buffer
	if err := tmpl.Execute(&output, thread); err != nil {
		return "", err
	}
	return template.HTML(output.String()), nil
}

var postTmpl = `
<div class="bsky-feed-thread bsky-feed-thread--lvl{{ .Level }}">
	{{ if gt .Level 0 }}
	<div class="bsky-feed-post">
	<div class="bsky-feed-post__avatar">
		<a href="https://bsky.app/profile/{{ .AuthorHandle }}" class="bsky-author-handle"><img src="{{ .AuthorAvatar }}" /></a>
	</div>
	<div class="bsky-feed-post__content">
		{{ .RenderedText }}
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
		{{ range .RenderedReplies }}
		{{ . }}
		{{ end }}
	</div>
	{{ end }}
</div>
`
