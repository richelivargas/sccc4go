package cloudconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type CloudConfigClient struct {
	host        string
	application string
	profile     string
	branch      string
	format      Format
	basicAuth   *basicAuthInfo
	raw         *map[string]interface{}
}

type basicAuthInfo struct {
	username string
	password string
}

func NewClient(host string, application string, profile string, opts ...ClientOption) (*CloudConfigClient, error) {

	if host == "" {
		return nil, errors.New("server is required")
	}

	if application == "" {
		return nil, errors.New("application is required")
	}

	if profile == "" {
		return nil, errors.New("a base profile is required")
	}

	client := &CloudConfigClient{
		host:        host,
		application: application,
		profile:     profile,
		format:      JSONFormat,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (client *CloudConfigClient) url() string {

	url := client.host

	if client.branch != "" {
		url = fmt.Sprintf("%s/%s", url, client.branch)
	}

	url = fmt.Sprintf("%s/%s-%s.%s", url, client.application, client.profile, client.format)

	return url
}

func (d *CloudConfigClient) Fetch() (io.ReadCloser, error) {

	u := d.url()
	if _, err := url.Parse(u); err != nil {
		return nil, errors.Wrapf(err, "invalid config url [%s]", u)
	}

	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	if d.basicAuth != nil {
		request.SetBasicAuth(d.basicAuth.username, d.basicAuth.password)
	}

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "config resolution failed")
	}

	if response.StatusCode > 299 {
		return nil, errors.Errorf("error do get config data: [%v] - %v", response.StatusCode, response.Status)
	}

	return response.Body, nil
}

func (d *CloudConfigClient) Decode(v interface{}) error {

	reader, err := d.Fetch()
	if err != nil {
		return err
	}

	if d.format == JSONFormat {
		return json.NewDecoder(reader).Decode(v)
	}

	return yaml.NewDecoder(reader).Decode(v)
}

func (d *CloudConfigClient) Raw() (map[string]interface{}, error) {
	var raw map[string]interface{}
	if err := d.Decode(&raw); err != nil {
		return nil, err
	}
	d.raw = &raw
	return raw, nil
}

func (d *CloudConfigClient) Get(keys ...string) interface{} {

	if d.raw == nil {
		if _, err := d.Raw(); err != nil {
			log.Fatal(err)
		}
	}

	var value interface{} = *d.raw
	for _, key := range keys {
		if value == nil {
			return value
		}
		value = value.(map[string]interface{})[key]
	}

	return value
}
