package nytbooks

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

const fixtureGetListResponse = `{
  "status": "OK",
  "copyright": "Copyright (c) 2018 The New York Times Company.  All Rights Reserved.",
  "num_results": 15,
  "last_modified": "2018-12-26T23:38:02-05:00",
  "results": [
    {
      "list_name": "Hardcover Fiction",
      "display_name": "Hardcover Fiction",
      "bestsellers_date": "2018-12-22",
      "published_date": "2019-01-06",
      "rank": 1,
      "rank_last_week": 1,
      "weeks_on_list": 9,
      "asterisk": 0,
      "dagger": 0,
      "amazon_product_url": "https://www.amazon.com/Reckoning-Novel-John-Grisham-ebook/dp/B079DBS447?tag=NYTBS-20",
      "isbns": [
        {
          "isbn10": "0385544154",
          "isbn13": "9780385544153"
        },
        {
          "isbn10": "0385544162",
          "isbn13": "9780385544160"
        },
        {
          "isbn10": "052563925X",
          "isbn13": "9780525639251"
        },
        {
          "isbn10": "0525639292",
          "isbn13": "9780525639299"
        },
        {
          "isbn10": "0385544170",
          "isbn13": "9780385544177"
        }
      ],
      "book_details": [
        {
          "title": "THE RECKONING",
          "description": "A decorated World War II veteran shoots and kills a pastor.",
          "contributor": "by John Grisham",
          "author": "John Grisham",
          "contributor_note": "",
          "price": 0,
          "age_group": "",
          "publisher": "Doubleday",
          "primary_isbn13": "9780385544153",
          "primary_isbn10": "0385544154"
        }
      ],
      "reviews": [
        {
          "book_review_link": "",
          "first_chapter_link": "",
          "sunday_review_link": "",
          "article_chapter_link": ""
        }
      ]
    }
	]
}`

/* A little HTTP testing inspiration from http://hassansin.github.io/Unit-Testing-http-client-in-Go */
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func getMockClient(fn RoundTripFunc) *Client {
	c := NewClient("mon_key")
	c.httpClient = NewTestHTTPClient(fn)
	return c
}

// TestEndpointFromOpts tests assembly of query params for API calls
func TestEndpointFromOpts(t *testing.T) {
	c := getMockClient(nil)
	url := c.endpointFromOpts("https://foo.bar/baz.json", map[string]string{"test": "a"})
	expected := "https://foo.bar/baz.json?api-key=mon_key&test=a"
	if url.String() != expected {
		t.Errorf("URL not as expected: expected=%s constructed=%s\n", expected, url)
	}
}

// TestGetBestSellers tests both request assembly and response deserialization
func TestGetBestSellers(t *testing.T) {
	fn := func(req *http.Request) *http.Response {
		expected := "https://api.nytimes.com/svc/books/v3/lists.json?api-key=mon_key&list=hardcover-fiction"
		if req.URL.String() != expected {
			t.Errorf("URL not as expected: expected=%s constructed=%s\n", expected, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(fixtureGetListResponse)),
			Header:     make(http.Header),
		}
	}
	c := getMockClient(fn)
	books, err := c.GetBestSellers("hardcover-fiction", nil)
	if err != nil {
		t.Errorf("Error response from GetBestSellers: %s", err)
	} else {
		if books[0].ListName != "Hardcover Fiction" {
			t.Errorf("Unexpected response from GetBestSellers: %+v", books)
		}
	}
}
