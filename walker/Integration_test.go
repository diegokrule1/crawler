package walker

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func writeConttent(w http.ResponseWriter, filePath string) {
	w.Header().Add("Content-type", "text/html")
	b, err := os.ReadFile(fmt.Sprintf("./responses/%s.html", filePath))
	if err != nil {
		log.Printf("could not read file for path %s", filePath)
		w.WriteHeader(400)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Printf("could not write http response for path %s", filePath)
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Path in mock %s\n", r.URL.Path)
	if r.URL.Path != "/" && r.URL.Path != "/gb" && r.URL.Path != "/sec" && r.URL.Path != "/third" {
		w.WriteHeader(401)
		return
	}
	if r.URL.Path == "/" {
		writeConttent(w, "init")
	} else {
		writeConttent(w, r.URL.Path)
	}

}

func TestConsumerProducer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(urlHandler))
	defer server.Close()
	err := os.Setenv("POSTGRES_USER", "postgres")
	if err != nil {
		t.Fatalf("Could not set postgres user env %v", err)
	}
	err = os.Setenv("POSTGRES_PASSWORD", "mysecretpassword")
	if err != nil {
		t.Fatalf("Could not set postgres passeword env %v", err)
	}
	err = os.Setenv("POSTGRES_DB", "crawler")
	if err != nil {
		t.Fatalf("Could not set postgres db env %v", err)
	}
	err = os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("Could not set postgres host	 env %v", err)
	}

	err = os.Setenv("POSTGRES_PORT", "5432")
	if err != nil {
		t.Fatalf("Could not set postgres port env %v", err)
	}
	consumer, producer, err := Init()
	if err != nil {
		t.Fatalf("could not create app %v", err)
		return
	}
	go consumer.Consume()
	serverUrl := server.URL
	parsedUrl, _ := url.Parse(serverUrl)
	producer.Produce(parsedUrl.Scheme, parsedUrl.Host, "/", nil)
	<-producer.KillChan
	total, err := consumer.Repo.countTerminated()
	if err != nil {
		t.Fatalf("could not check total records processed %v", err)
	}

	if total != 4 {
		t.Fatalf("Expecting 4 records processed. Found %d", total)
	}
}
