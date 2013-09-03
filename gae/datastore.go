package gae

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/icub3d/urls"
)

const (
	// The Kind for URLs.
	urlKind = "URL"

	// The Kind for Logs.
	logKind = "Log"

	// The Kind for Statistics.
	statKind = "Statistics"
)

// DataStore implements the urls.DataStore interface
type DataStore struct {
	cxt appengine.Context
}

func NewDataStore(cxt appengine.Context) *DataStore {
	return &DataStore{
		cxt: cxt,
	}
}

func (ds *DataStore) CountURLs() (int, error) {
	q := datastore.NewQuery(urlKind)
	return q.Count(ds.cxt)
}

func (ds *DataStore) GetURLs(limit, offset int) ([]*urls.URL, error) {
	q := datastore.NewQuery(urlKind).Order("-Created").
		Offset(offset).Limit(limit)

	us := make([]*urls.URL, 0, limit)
	_, err := q.GetAll(ds.cxt, *us)
	return us, err
}

func (ds *DataStore) GetURL(id string) (*urls.URL, error) {
	key := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	var u urls.URL
	err := datastore.Get(ds.cxt, key, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ds *DataStore) DeleteURL(id string) error {
	key := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	// TODO we need to delete the stats as well as the Logs.

	return datastore.Delete(ds.cxt, key)
}

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

// TODO GetStatistics
func (ds *DataStore) GetStatistics(short string) (*Statistics, error) {
	stats := urls.NewStatistics(short)

	data := make([]byte, 0)
	key := datastore.NewKey(ds.cxt, statsKind, "", urls.ShortToInt(id), nil)
	err := datastore.Get(ds.cxt, key, &data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// TODO PutStatistics
func (ds *DataStore) PutStatistics(stats *Statistics) error {
	// Not sure how to do maps in datastore, so I'm simply doing a json
	// encoded value.
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	key := datastore.NewKey(ds.cxt, statKind, "",
		urls.ShortToInt(stat.Short), nil)

	_, err := datastore.Put(ds.cxt, key, data)

	return err
}

func (ds *DataStore) LogClick(l *urls.Log) error {
	pkey := datastore.NewKey(ds.cxt, urlKind, "",
		urls.ShortToInt(l.Short), nil)
	key := datastore.NewIncompleteKey(ds.cxt, logKind, pkey)

	_, err := datastore.Put(ds.cxt, key, l)

	return err
}

func (ds *DataStore) CountLogs() (int, error) {
	q := datastore.NewQuery(logKind)
	return q.Count(ds.cxt)
}

func (ds *DataStore) GetLogs(id string, limit, offset int) ([]*urls.URL,
	error) {

	pkey := datastore.NewKey(ds.cxt, urlKind, "", urls.ShortToInt(id), nil)

	q := datastore.NewQuery(logKind).Ancestor(pkey).Order("Created").
		Offset(offset).Limit(limit)

	us := make([]*urls.Log, 0, limit)
	_, err := q.GetAll(ds.cxt, *us)
	return us, err
}
