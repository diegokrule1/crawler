package walker

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Walker struct {
	Logger *zap.Logger
}

func (w *Walker) Walk(url string, defaultDomain string) ([]string, error) {
	w.Logger.Info(fmt.Sprintf("Invoking url %s", url))
	resp, err := http.Get(url)
	if err != nil {
		w.Logger.Error("Could not invoke url", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	responseContentType := resp.Header.Get("Content-Type")
	if strings.Contains(responseContentType, "text/html") {
		return w.analyzeResponse(resp, defaultDomain)
	}
	return nil, nil
}

func (w *Walker) analyzeResponse(response *http.Response, defaultDomain string) ([]string, error) {
	bts, err := io.ReadAll(response.Body)

	if err != nil {
		w.Logger.Error("could not read http response", zap.Error(err))
		return nil, err
	}

	w.Logger.Info("Response", zap.String("resp", string(bts)))

	return search(string(bts), defaultDomain), nil
}
