// Copyright 2013 Joshua Marsh. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gae

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"github.com/icub3d/urls"
)

const (
	// The Kind for URLs.
	urlKind = "URL"

	// The Kind for Logs.
	logKind = "Log"

	// The Kind for Statistics.
	statsKind = "Stats"
)

// DataStore implements the urls.DataStore interface
type DataStore struct {
	cxt appengine.Context
}

// NewDataStore creates a new datastore with the given context.
func NewDataStore(cxt appengine.Context) *DataStore {
	return &DataStore{
		cxt: cxt,
	}
}

// CountURLs implements the urls.DataStore interface.
func (ds *DataStore) CountURLs() (int, error) {
	q := datastore.NewQuery(urlKind)
	return q.Count(ds.cxt)
}

// GetURLs implements the urls.DataStore interface.
func (ds *DataStore) GetURLs(limit, offset int) ([]*urls.URL, error) {
	q := datastore.NewQuery(urlKind).Order("-Created").
		Offset(offset).Limit(limit)

	us := make([]*urls.URL, 0, limit)
	_, err := q.GetAll(ds.cxt, &us)
	return us, err
}

// GetURL implements the urls.DataStore interface.
func (ds *DataStore) GetURL(id string) (*urls.URL, error) {
	var u urls.URL

	if item, err := memcache.Get(ds.cxt, id); err != nil {
		// When we don't find one, we should get it from the datastore.
		key := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

		err := datastore.Get(ds.cxt, key, &u)
		if err == datastore.ErrNoSuchEntity {
			return nil, urls.ErrNotFound
		} else if err != nil {
			ds.cxt.Errorf("failed to get %v in datastore: %v", id, err)
			return nil, err
		}

		// Try to save it to the datastore.
		data, err := json.Marshal(u)
		if err != nil {
			return nil, err
		}
		item = &memcache.Item{
			Key:   id,
			Value: data,
		}

		if err := memcache.Set(ds.cxt, item); err != nil {
			ds.cxt.Errorf("failed to set %v in memcache: %v", id, err)
		}

	} else {
		// We got it from memcache.
		err := json.Unmarshal(item.Value, &u)
		if err != nil {
			return nil, err
		}
	}

	return &u, nil
}

// DeleteURL implements the urls.DataStore interface.
func (ds *DataStore) DeleteURL(id string) error {
	key := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	// Delete the logs.
	iter := datastore.NewQuery(logKind).Ancestor(key).Run(ds.cxt)
	var l urls.Log
	for k, err := iter.Next(&l); err == nil; k, err = iter.Next(&l) {
		datastore.Delete(ds.cxt, k)
	}

	// Delete the stats.
	skey := datastore.NewKey(ds.cxt, statsKind, "", urls.ShortToInt(id), nil)
	datastore.Delete(ds.cxt, skey)

	return datastore.Delete(ds.cxt, key)
}

// PutURL implements the urls.DataStore interface.
func (ds *DataStore) PutURL(u *urls.URL) (string, error) {

	// We may need to create an ID.
	if u.Short == "" {
		i, _, err := datastore.AllocateIDs(ds.cxt, urlKind, nil, 1)
		if err != nil {
			return "", err
		}

		u.Short = urls.IntToShort(i)
	}

	// Get the key
	key := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(u.Short), nil)

	_, err := datastore.Put(ds.cxt, key, u)
	if err != nil {
		return "", err
	}

	return u.Short, nil
}

// I guess these things need to be stored as a struct.
type statData struct {
	Data []byte
}

// GetStatistics implements the urls.DataStore interface.
func (ds *DataStore) GetStatistics(id string) (*urls.Statistics, error) {
	stats := urls.NewStatistics(id)

	s := statData{Data: []byte{}}

	key := datastore.NewKey(ds.cxt, statsKind, "", urls.ShortToInt(id), nil)
	err := datastore.Get(ds.cxt, key, &s)
	if err == datastore.ErrNoSuchEntity {
		return nil, urls.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(s.Data, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// PutStatistics implements the urls.DataStore interface.
func (ds *DataStore) PutStatistics(stats *urls.Statistics) error {
	// Not sure how to do maps in datastore, so I'm simply doing a json
	// encoded value.
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	key := datastore.NewKey(ds.cxt, statsKind, "",
		urls.ShortToInt(stats.Short), nil)

	s := statData{Data: data}
	_, err = datastore.Put(ds.cxt, key, &s)

	return err
}

// LogClick implements the urls.DataStore interface.
func (ds *DataStore) LogClick(l *urls.Log) error {
	pkey := datastore.NewKey(ds.cxt, urlKind, "",
		urls.ShortToInt(l.Short), nil)
	key := datastore.NewIncompleteKey(ds.cxt, logKind, pkey)

	_, err := datastore.Put(ds.cxt, key, l)

	return err
}

// CountLogs implements the urls.DataStore interface.
func (ds *DataStore) CountLogs(id string) (int, error) {
	pkey := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	q := datastore.NewQuery(logKind).Ancestor(pkey)
	return q.Count(ds.cxt)
}

// GetLogs implements the urls.DataStore interface.
func (ds *DataStore) GetLogs(id string, limit, offset int) ([]*urls.Log,
	error) {

	pkey := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	q := datastore.NewQuery(logKind).Ancestor(pkey).Order("Created").
		Offset(offset).Limit(limit)

	us := make([]*urls.Log, 0, limit)
	_, err := q.GetAll(ds.cxt, &us)
	return us, err
}
