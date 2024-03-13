package walker

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
)

type Producer struct {
	Logger *zap.Logger

	PageChan chan string

	Repository UrlRepository

	KillChan chan bool
}

func (p *Producer) Produce(scheme, host, path string, parentId *string) {
	id := uuid.New().String()
	total, err := p.Repository.CreateUrl(scheme, host, path, id, parentId)
	if err != nil {
		log.Fatalf("could not persist url %v", err)
	}
	fmt.Printf("Total inserted %d", total)
	if total > 0 {
		p.PageChan <- id
	}
}
