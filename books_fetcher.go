package main

import (
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

func (b *BooksFetcher) updateCache() (err error) {
	b.cacheLock.Lock()
	defer b.cacheLock.Unlock()

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
	return nil
}

func (b *BooksFetcher) cacheExpired() bool {
	return time.Now().After(b.cacheExpiration)
}
