package once

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/wvoliveira/once/configs"
	"github.com/wvoliveira/once/internal/pkg/util"
)

// Limite de upload (ex: 10MB) para segurança
const MaxUploadSize = 10 << 20

// Estrutura para resposta JSON
type UploadResponse struct {
	ID        string `json:"id"`
	Link      string `json:"link"`
	ExpiresAt string `json:"expires_at"`
}

type TransportHTTP interface {
	Start()
	HTTPUpload(http.ResponseWriter, *http.Request)
	HTTPDownload(http.ResponseWriter, *http.Request)
}

type transportHTTP struct {
	config     configs.Config
	service    Service
	httpServer *http.Server
}

func NewTransportHTTP(config configs.Config, svc Service) TransportHTTP {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(httprate.LimitByIP(10, time.Second))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPServerPort),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	handler := transportHTTP{
		service:    svc,
		config:     config,
		httpServer: server,
	}

	r.Post("/api/upload", handler.HTTPUpload)
	r.Post("/api/download", handler.HTTPDownload)

	return handler
}

func (s transportHTTP) Start() {
	// Canal para graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("HTTP listening", "port", s.config.HTTPServerPort)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Mantém o web server até o sinal de "parar" chegar até a aplicação.
	// O sinal pode ser o SIGINT (ctrl+c, desenvolvimento local) ou o SIGTERM (comando kill do linux)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	slog.Info("shutting down HTTP server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("error to shutdown nicely", "error", err.Error())
	}

	slog.Info("bye!")
}

func (t transportHTTP) HTTPUpload(w http.ResponseWriter, r *http.Request) {
	// Segurança: Limita tamanho do body
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// Parse do Multipart Form
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "Arquivo muito grande (Max 10MB)", http.StatusBadRequest)
		return
	}

	// Pega o arquivo do form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Campo 'file' é obrigatório", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Lê o conteúdo para memória (Para arquivos gigantes, usaríamos io.Copy direto pro disco)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erro ao ler arquivo", http.StatusInternalServerError)
		return
	}

	// Pega o TTL (opcional, default 1h)
	ttlSeconds, _ := strconv.Atoi(r.FormValue("ttl"))
	if ttlSeconds <= 0 {
		ttlSeconds = 3600
	}
	ttl := time.Duration(ttlSeconds) * time.Second

	// Gera ID e Salva
	id := util.GenerateID()
	err = t.service.Upload(id, header.Filename, fileBytes, ttl)
	if err != nil {
		http.Error(w, "Erro ao salvar", http.StatusInternalServerError)
		return
	}

	// Retorna JSON
	resp := UploadResponse{
		ID:        id,
		Link:      fmt.Sprintf("http://localhost:8080/files/%s", id),
		ExpiresAt: time.Now().Add(ttl).Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (t transportHTTP) HTTPDownload(w http.ResponseWriter, r *http.Request) {
	// Pega o ID da URL (Go 1.22 feature)
	id := r.PathValue("id")

	// Chama a lógica de visualização única
	data, meta, err := t.service.Download(id)
	if err != nil {
		// Se deu erro, assumimos 404 para não vazar se o arquivo existia ou não
		http.Error(w, "Arquivo não encontrado ou link expirado", http.StatusNotFound)
		return
	}

	// Configura Headers para Download
	// Content-Disposition: attachment força o download em vez de abrir no navegador
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", meta.FileName))
	w.Header().Set("Content-Type", http.DetectContentType(data)) // Ou meta.MimeType se você salvou
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	// Escreve os bytes
	w.Write(data)
}
