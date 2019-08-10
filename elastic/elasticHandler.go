package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
	"go-delic-products/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PostElastic struct {
	post model.Post
}

func (p *PostElastic) Save(post *model.Post, c elasticsearch.Client) (string, error) {

	jsonPost, _ := json.Marshal(post)
	request := esapi.IndexRequest{
		Index:      p.post.GetIndexName(),
		DocumentID: uuid.New().String(),
		Body:       strings.NewReader(string(jsonPost)),
		Refresh:    "true",
	}

	res, err := request.Do(context.Background(), &c)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return "", errors.New("errors during the response")
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		}
		return r["_id"].(string), nil
	}

}

func (p *PostElastic) FindById(id string, c elasticsearch.Client) (*esapi.Response, error) {
	req := esapi.GetRequest{Index: p.post.GetIndexName(), DocumentID: id}

	res, err := req.Do(context.Background(), &c)

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		log.Fatal("error parsing response")
	}

	return res, nil
}

func (p *PostElastic) FindByCriteria(criteria io.Reader) (string, error) {
	url := "http://localhost:9200/shared_post/_search"

	//r := "{\n    \"query\": {\n        \"match\" : {\n            \"title\" : \"FBI with coffee\"\n        }\n    }\n}"
	//payload := strings.NewReader(r)

	request, _ := ioutil.ReadAll(criteria)

	query := string(request)
	//fmt.Printf("%v", r)
	//fmt.Printf("%v", string(request))

	req, _ := http.NewRequest("POST", url, strings.NewReader(query))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "f7a22433-456b-4e81-87c4-67e959e2a034")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body), nil
}
