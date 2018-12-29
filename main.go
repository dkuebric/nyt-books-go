package main

import (
	"fmt"
	"log"

	"github.com/dkuebric/nyt-books-go/nytbooks"
)

func main() {
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
}
