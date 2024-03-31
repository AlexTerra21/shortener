package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"go.uber.org/zap"
)

// Типы запросов для сжатия
const compressContent = "application/json, text/html"

// Проверка, следует ли сжимать контент
func IsCompress(content string) bool {
	return strings.Contains(compressContent, content)
}

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// Инициализация структуры compressWriter
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w), //
	}
}

// Реализация метода Header интерфейса ResponseWriter
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Реализация метода Write интерфейса ResponseWriter
func (c *compressWriter) Write(p []byte) (int, error) {
	if IsCompress(c.Header().Get("Content-Type")) {
		return c.zw.Write(p)
	} else {
		return c.w.Write(p)
	}
}

// Реализация метода WriteHeader интерфейса ResponseWriter
func (c *compressWriter) WriteHeader(statusCode int) {
	if IsCompress(c.Header().Get("Content-Type")) {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if IsCompress(c.Header().Get("Content-Type")) {
		return c.zw.Close()
	} else {
		return nil
	}
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// Инициализация Reader с функцией декомпрессии
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Декмпрессует поток байт
func (c *compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

// Закрытие Reader
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Middleware для компрессии/декомпрессии
func WithCompress(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				logger.Log().Error("Error read encoded body", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		logger.Log().Debug("WithCompress",
			zap.String("Content-Type", r.Header.Get("Content-Type")),
			zap.String("Accept-Encoding", r.Header.Get("Accept-Encoding")),
			zap.String("Content-Encoding", r.Header.Get("Content-Encoding")),
		)

		h.ServeHTTP(ow, r)
	}
}
