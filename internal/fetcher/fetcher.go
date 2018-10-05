package fetcher

import (
	"sync"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/fetcher/services"
	"github.com/objque/musicmash/internal/fetcher/services/itunes"
	"github.com/objque/musicmash/internal/fetcher/services/yandex"
	"github.com/objque/musicmash/internal/log"
)

func getServices() []services.Service {
	fetchers := []services.Service{}
	for _, store := range config.Config.Stores {
		switch store.Name {
		case "itunes":
			fetchers = append(fetchers, itunes.NewService(store.URL, store.Meta["token"]))
		case "yandex":
			fetchers = append(fetchers, yandex.NewService(store.URL))
		}
	}
	return fetchers
}

func fetchFromServices(services []services.Service) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(services))

	// fetch from all services
	for i := range services {
		go services[i].FetchAndSave(&wg)
	}

	return &wg
}

func Fetch() {
	fetchFromServices(getServices()).Wait()

	// run callback
	log.Info("All stores were fetched")
}
