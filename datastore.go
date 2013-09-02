package urls

// DataStore is the interface that any backend datastore should
// implement to be compatible with the handlers.
type DataStore interface {
	// Get the total count of Urls in the system.
	CountURLs() (int, error)

	// Get the next limit urls order by create date (newest first) and offset by the
	// given offset.
	GetURLs(limit, offset int) ([]*URL, error)

	// Get the url with the given short id.
	GetURL(short string) (*URL, error)

	// Remove the given url and it's associated logs and statistics.
	DeleteURL(short string) error

	// Put the given url into the data store. If the short id exists,
	// overwrite it. Otherwise insert a new entry. When inserting, the
	// new short ID should be updated before insertion and it should be
	// returned. You can use the helper functions IntToShort to help
	// convert a unique integer to a representative string.
	PutURL(url *URL) (string, error)

	// Get the statistics for the given short id.
	GetStatistics(short string) (*Statistics, error)

	// Put the given statistics into the data store. If the short id
	// exists, overwrite it. Otherwise insert a new entry.
	PutStatistics(stats *Statistics) error

	// Log a click from the system.
	LogClick(l *Log) error

	// Get the total count of Logs in the system for the given short id.
	CountLogs(short string) (int, error)

	// Get the next limit logs of the given short id sorted by create
	// date (newest first) and offset by the given offset.
	GetLogs(short string, limit, offset int) ([]*Log, error)
}
