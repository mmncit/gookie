package store

import (
	"database/sql"
	"errors"
)

// Event represents a top‐level event record.
type Event struct {
	ID          int          `json:"id"`
	UserID      int          `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartTime   string       `json:"start_time"` // RFC3339
	EndTime     string       `json:"end_time"`   // RFC3339
	Location    string       `json:"location"`
	Entries     []EventEntry `json:"entries"`
}

// EventEntry represents a “sub‐entry” of an Event (each with its own type, times, etc.).
type EventEntry struct {
	ID          int    `json:"id"`
	EventTypeID int    `json:"event_type_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"` // RFC3339
	EndTime     string `json:"end_time"`   // RFC3339
	Location    string `json:"location"`
}

// PostgresEventStore wraps a *sql.DB and implements EventStore.
type PostgresEventStore struct {
	db *sql.DB
}

// NewPostgresEventStore constructs a new store backed by Postgres.
func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

// EventStore defines the interface for storing/retrieving Events.
type EventStore interface {
	CreateEvent(*Event) (*Event, error)
	GetEventByID(id int64) (*Event, error)
	UpdateEvent(*Event) error
	DeleteEvent(id int64) error
	GetEventOwner(id int64) (int, error)
}

// CreateEvent inserts an Event and its associated EventEntry rows transactionally.
func (pg *PostgresEventStore) CreateEvent(evt *Event) (*Event, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1) Insert into events table, returning the new ID.
	insertMain := `
		INSERT INTO events (user_id, title, description, start_time, end_time, location)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRow(
		insertMain,
		evt.UserID,
		evt.Title,
		evt.Description,
		evt.StartTime,
		evt.EndTime,
		evt.Location,
	).Scan(&evt.ID)
	if err != nil {
		return nil, err
	}

	// 2) Insert each EventEntry
	for i := range evt.Entries {
		entry := &evt.Entries[i]
		insertEntry := `
			INSERT INTO event_entries
				(event_id, event_type_id, title, description, start_time, end_time, location)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		err = tx.QueryRow(
			insertEntry,
			evt.ID,
			entry.EventTypeID,
			entry.Title,
			entry.Description,
			entry.StartTime,
			entry.EndTime,
			entry.Location,
		).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return evt, nil
}

// GetEventByID retrieves an Event by its ID, along with its entries.
func (pg *PostgresEventStore) GetEventByID(id int64) (*Event, error) {
	evt := &Event{}
	queryMain := `
		SELECT id, user_id, title, description, start_time, end_time, location
		FROM events
		WHERE id = $1
	`
	err := pg.db.QueryRow(queryMain, id).Scan(
		&evt.ID,
		&evt.UserID,
		&evt.Title,
		&evt.Description,
		&evt.StartTime,
		&evt.EndTime,
		&evt.Location,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Now load associated entries (ordered by id)
	entryQuery := `
		SELECT id, event_type_id, title, description, start_time, end_time, location
		FROM event_entries
		WHERE event_id = $1
		ORDER BY id
	`
	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry EventEntry
		err = rows.Scan(
			&entry.ID,
			&entry.EventTypeID,
			&entry.Title,
			&entry.Description,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Location,
		)
		if err != nil {
			return nil, err
		}
		evt.Entries = append(evt.Entries, entry)
	}
	return evt, nil
}

// UpdateEvent replaces the Event’s fields and re‐inserts all its entries.
func (pg *PostgresEventStore) UpdateEvent(evt *Event) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1) Ensure the event exists and update its columns.
	updateMain := `
		UPDATE events
		SET title = $1,
		    description = $2,
		    start_time = $3,
		    end_time = $4,
		    location = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`
	res, err := tx.Exec(
		updateMain,
		evt.Title,
		evt.Description,
		evt.StartTime,
		evt.EndTime,
		evt.Location,
		evt.ID,
	)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return sql.ErrNoRows
	}

	// 2) Delete all existing entries for that event.
	_, err = tx.Exec(`DELETE FROM event_entries WHERE event_id = $1`, evt.ID)
	if err != nil {
		return err
	}

	// 3) Re‐insert each entry from evt.Entries
	for i := range evt.Entries {
		entry := &evt.Entries[i]
		insertEntry := `
			INSERT INTO event_entries
				(event_id, event_type_id, title, description, start_time, end_time, location)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		err = tx.QueryRow(
			insertEntry,
			evt.ID,
			entry.EventTypeID,
			entry.Title,
			entry.Description,
			entry.StartTime,
			entry.EndTime,
			entry.Location,
		).Scan(&entry.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteEvent removes the Event row (and cascades entries if FOREIGN KEY is ON DELETE CASCADE).
func (pg *PostgresEventStore) DeleteEvent(id int64) error {
	del, err := pg.db.Exec(`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return err
	}
	ra, err := del.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetEventOwner returns the user_id who owns this event.
func (pg *PostgresEventStore) GetEventOwner(eventID int64) (int, error) {
	var userID int
	err := pg.db.QueryRow(
		`SELECT user_id FROM events WHERE id = $1`,
		eventID,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}
