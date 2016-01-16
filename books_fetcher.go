package main

import (
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/unixpickle/gscrape"
)

// BooksFetcher fetches a user's Play Books library periodically, caching
// results when possible.
type BooksFetcher struct {
	email    string
	password string

	cacheLock       sync.RWMutex
	cacheExpiration time.Time
	cache           []gscrape.BookInfo
	session         *gscrape.Session

	thumbnailCacheLock sync.RWMutex
	thumbnailCache     map[string][]byte
}

func NewBooksFetcher(email, password string) *BooksFetcher {
	return &BooksFetcher{email: email, password: password}
}

func (b *BooksFetcher) Books() ([]gscrape.BookInfo, error) {
	b.cacheLock.RLock()
	if !b.cacheExpired() {
		defer b.cacheLock.RUnlock()
		return b.cache, nil
	}
	b.cacheLock.RUnlock()

	if err := b.updateCache(); err != nil {
		return nil, err
	}

	b.cacheLock.RLock()
	defer b.cacheLock.RUnlock()
	return b.cache, nil
}

func (b *BooksFetcher) Thumbnail(bookId string) ([]byte, error) {
	b.cacheLock.RLock()
	if b.cacheExpired() {
		b.cacheLock.RUnlock()
		if err := b.updateCache(); err != nil {
			return nil, err
		}
		b.cacheLock.RLock()
	}
	defer b.cacheLock.RUnlock()

	b.thumbnailCacheLock.RLock()
	thumbnail, ok := b.thumbnailCache[bookId]
	b.thumbnailCacheLock.RUnlock()
	if ok {
		return thumbnail, nil
	}

	var imageURL string
	for _, book := range b.cache {
		if book.ID == bookId {
			imageURL = book.ImageLinks.Thumbnail
		}
	}
	if imageURL == "" {
		return nil, errors.New("book ID is not in the user's library")
	}
	resp, err := b.session.Get(imageURL + "&usc=0&w=300")
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b.thumbnailCacheLock.Lock()
	b.thumbnailCache[bookId] = result
	b.thumbnailCacheLock.Unlock()
	return result, nil
}

func (b *BooksFetcher) updateCache() (err error) {
	b.cacheLock.Lock()
	defer b.cacheLock.Unlock()

	b.session = nil
	b.thumbnailCache = map[string][]byte{}
	b.cache = nil

	// Multiple Books() calls might try to update the cache simultaneously.
	if !b.cacheExpired() {
		return nil
	}

	client := gscrape.NewSession()
	playBooks, err := client.AuthPlayBooks(b.email, b.password)
	if err != nil {
		return
	}

	books, errChan := playBooks.MyBooks(gscrape.AllBookSources)
	b.cache = make([]gscrape.BookInfo, 0)
	for book := range books {
		b.cache = append(b.cache, book)
	}
	if err = <-errChan; err != nil {
		return
	}

	b.cacheExpiration = time.Now().Add(time.Minute * 15)
	b.session = client
	return nil
}

func (b *BooksFetcher) cacheExpired() bool {
	// If it's January 1st, 1970, then checking b.cache might be necessary.
	return time.Now().After(b.cacheExpiration) || b.cache == nil
}
