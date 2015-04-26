package elastic

import (
  "net/url"
  "errors"
  "encoding/json"
  "time"
  "log"

  "github.com/gliderlabs/registrator/bridge"
  elasticapi "github.com/olivere/elastic"
)

type Container struct {
  Name      string `json:"container_name"`
  Action    string `json:"action"`
  Message   string `json:"message"`
  Timestamp string `json:"@timestamp"`
}

func init() {
  bridge.Register(new(Factory), "elastic")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
  urls := "http://127.0.0.1:9200"

  if uri.Host != "" {
    urls = "http://"+uri.Host
  }

  client, err := elasticapi.NewClient(elasticapi.SetURL(urls))
  if err != nil {
    log.Fatal("elastic: ", uri.Scheme)
  }

  return &ElasticAdapter{client: client}
}

type ElasticAdapter struct {
  client   *elasticapi.Client
}

func (r *ElasticAdapter) Ping() error {
  status := r.client.IsRunning()

  if !status {
    return errors.New("client is not Running")
  }

  return nil
}

func (r *ElasticAdapter) Register(service *bridge.Service) error {
  serviceAsJson, err := json.Marshal(service)
  if err != nil {
    return err
  }

  timestamp := time.Now().Local().Format("2006-01-02T15:04:05.000Z07:00")

  // Add a document to the index
  container := Container{Name: service.Name, Action: "start", Message: string(serviceAsJson), Timestamp: timestamp}
  _, err = r.client.Index().
    Index("containers").
    Type("audit").
    BodyJson(container).
    Timestamp(timestamp).
    Do()
  if err != nil {
    return err
  }
  return nil
}

func (r *ElasticAdapter) Deregister(service *bridge.Service) error {
  serviceAsJson, err := json.Marshal(service)
  if err != nil {
    return err
  }

  timestamp := time.Now().Local().Format("2006-01-02T15:04:05.000Z07:00")

  // Add a document to the index
  container := Container{Name: service.Name, Action: "stop", Message: string(serviceAsJson), Timestamp: timestamp}
  _, err = r.client.Index().
    Index("containers").
    Type("audit").
    BodyJson(container).
    Timestamp(timestamp).
    Do()
  if err != nil {
    return err
  }
  return nil
}

func (r *ElasticAdapter) Refresh(service *bridge.Service) error {
  return nil
}
