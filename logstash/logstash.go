package logstash

import (
  "net"
  "net/url"
  "encoding/json"
  "log"
  "strconv"

  "github.com/gliderlabs/registrator/bridge"
  logstashapi "github.com/heatxsink/go-logstash"
)

type Container struct {
  Name      string            `json:"container_name"`
  Action    string            `json:"action"`
  Service   *bridge.Service   `json:"info"`
}

func init() {
  bridge.Register(new(Factory), "logstash")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
  urls := "127.0.0.1:5959"

  if uri.Host != "" {
    urls = uri.Host
  }

  host, port, err := net.SplitHostPort(urls)
  if err != nil {
    log.Fatal("logstash: ", "split error")
  }

  intPort, _ := strconv.Atoi(port)

  client := logstashapi.New(host, intPort, 5000)

  return &LogstashAdapter{client: client}
}

type LogstashAdapter struct {
  client   *logstashapi.Logstash
}

func (r *LogstashAdapter) Ping() error {
  _, err := r.client.Connect()
  if err != nil {
    return err
  }

  return nil
}

func (r *LogstashAdapter) Register(service *bridge.Service) error {
  container := Container{Name: service.Name, Action: "start", Service: service}
  asJson, err := json.Marshal(container)
  if err != nil {
    return err
  }

  _, err = r.client.Connect()
  if err != nil {
    return err
  }

  err = r.client.Writeln(string(asJson))
  if err != nil {
    return err
  }

  return nil
}

func (r *LogstashAdapter) Deregister(service *bridge.Service) error {
  container := Container{Name: service.Name, Action: "stop", Service: service}
  asJson, err := json.Marshal(container)
  if err != nil {
    return err
  }

  _, err = r.client.Connect()
  if err != nil {
    return err
  }

  err = r.client.Writeln(string(asJson))
  if err != nil {
    return err
  }

  return nil
}

func (r *LogstashAdapter) Refresh(service *bridge.Service) error {
  return nil
}
