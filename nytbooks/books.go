// Package nytbooks provides a client for the NYT Books API:
// http://developer.nytimes.com/books_api.json
package nytbooks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const baseURL = "https://api.nytimes.com/svc"
const booksListsURL = baseURL + "/books/v3/lists.json"

// The Client struct encapsulates the API client
type Client struct {
	key        string
	httpClient *http.Client
}

// NewClient generates a new API client, provided an API key.
func NewClient(key string) *Client {
	return &Client{key, &http.Client{}}
}

// BookRanking describes an entry in a NYT book list (eg. Hardcover Fiction best sellers).
// Structs generated via https://mholt.github.io/json-to-go/
type BookRanking struct {
	ListName         string `json:"list_name"`
	DisplayName      string `json:"display_name"`
	BestSellersDate  string `json:"bestsellers_date"`
	PublishedDate    string `json:"published_date"`
	Rank             int    `json:"rank"`
	RankLastWeek     int    `json:"rank_last_week"`
	WeeksOnList      int    `json:"weeks_on_list"`
	Asterisk         int    `json:"asterisk"`
	Dagger           int    `json:"dagger"`
	AmazonProductURL string `json:"amazon_product_url"`
	Isbns            []struct {
		Isbn10 string `json:"isbn10"`
		Isbn13 string `json:"isbn13"`
	} `json:"isbns"`
	BookDetails []struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		Contributor     string `json:"contributor"`
		Author          string `json:"author"`
		ContributorNote string `json:"contributor_note"`
		Price           int    `json:"price"`
		AgeGroup        string `json:"age_group"`
		Publisher       string `json:"publisher"`
		PrimaryIsbn13   string `json:"primary_isbn13"`
		PrimaryIsbn10   string `json:"primary_isbn10"`
	} `json:"book_details"`
	Reviews []struct {
		BookReviewLink     string `json:"book_review_link"`
		FirstChapterLink   string `json:"first_chapter_link"`
		SundayReviewLink   string `json:"sunday_review_link"`
		ArticleChapterLink string `json:"article_chapter_link"`
	} `json:"reviews"`
}

type booksListResponse struct {
	BookRankings []BookRanking `json:"results"`
	Status       string        `json:"status"`
	NumResults   int           `json:"num_results"`
	LastModified string        `json:"last_modified"`
}

// GetBestSellers gets the most recent list of best sellers for a particular list name.
func (c *Client) GetBestSellers(list string, params map[string]string) ([]BookRanking, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["list"] = list

	url := c.endpointFromOpts(booksListsURL, params)

	raw, err := c.getAPIResponse(url)
	if err != nil {
		return nil, err
	}

	var resp booksListResponse
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		log.Printf("Couldn't unmarshal json: %s\n", raw)
		return nil, err
	}

	return resp.BookRankings, nil
}

func (c *Client) getAPIResponse(url *url.URL) ([]byte, error) {
	resp, err := c.httpClient.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (c *Client) endpointFromOpts(endpointURL string, options map[string]string) *url.URL {
	u, _ := url.ParseRequestURI(endpointURL)
	data := url.Values{}
	data.Add("api-key", c.key)
	for o, v := range options {
		data.Add(o, v)
	}
	u.RawQuery = data.Encode()
	return u
}
