# nyt-books-go

A Golang client for the [New York Times Books API](http://developer.nytimes.com/books_api.json).

Currently supports only the [GET /lists](http://developer.nytimes.com/books_api.json#/Documentation/GET/lists.%7Bformat%7D) endpoint.

## Installing

```
go get github.com/dkuebric/nyt-books-go
```

## Usage:

```go
booksClient := nytbooks.NewClient("126cd3f9d17941c4ba6c1fca6b47c734")

// Get the most recent hardcover fiction best sellers
books, err := booksClient.GetBestSellers("hardcover-fiction", nil)
if err != nil {
	log.Printf("Error getting books %+v\n", err)
}

// Make a shopping list
for book := range books {
	fmt.Printf("%s => %s\n", books[book].BookDetails[0].Title, books[book].AmazonProductURL)
}
```
