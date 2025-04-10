package linkfylink

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/yantology/linkfy/pkg/customerror"
)

type linkfyLinkPostgres struct {
	db *sql.DB
}

func NewLinkfyLinkPostgres(db *sql.DB) LinkfyLinkDBInterface {
	return &linkfyLinkPostgres{db: db}
}

// CreateLinks creates multiple links for a linkfy profile with specified order
func (l *linkfyLinkPostgres) CreateLinks(linkfy_id string, links []*LinkfyLinkCreated) *customerror.CustomError {
	log.Printf("Creating links for linkfy_id: %s with %d links", linkfy_id, len(links))

	tx, err := l.db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return customerror.NewCustomError(err, "Failed to begin transaction", 500)
	}

	deleteQuery := "DELETE FROM linkfy_links WHERE linkfy_id = $1"
	log.Printf("Executing delete query: %s with linkfy_id: %s", deleteQuery, linkfy_id)

	_, err = tx.Exec(deleteQuery, linkfy_id)
	if err != nil {
		log.Printf("Error deleting existing links: %v", err)
		tx.Rollback()
		return customerror.NewCustomError(err, "Failed to delete existing links", 500)
	}

	insertQuery := "INSERT INTO linkfy_links (linkfy_id, name, name_url, icons_url, display_order) VALUES ($1, $2, $3, $4, $5)"
	log.Printf("Preparing insert query: %s", insertQuery)

	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		tx.Rollback()
		return customerror.NewCustomError(err, "Failed to prepare statement", 500)
	}
	defer stmt.Close()

	for i, link := range links {
		log.Printf("Inserting link %d: {name: %s, nameURL: %s, iconsURL: %s, order: %d}",
			i, link.Name, link.NameURL, link.IconsURL, i)

		_, err = stmt.Exec(linkfy_id, link.Name, link.NameURL, link.IconsURL, i)
		if err != nil {
			log.Printf("Error inserting link: %v", err)
			tx.Rollback()
			return customerror.NewCustomError(err, "Failed to insert link", 500)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return customerror.NewCustomError(err, "Failed to commit transaction", 500)
	}

	log.Printf("Successfully created %d links for linkfy_id: %s", len(links), linkfy_id)
	return nil
}

// GetLinkByLinkfyID retrieves links by their linkfy_id
func (l *linkfyLinkPostgres) GetLinkByLinkfyID(linkfy_id string) ([]*LinkfyLink, *customerror.CustomError) {
	query := `
		SELECT id, name, name_url, icons_url, created_at 
		FROM linkfy_links 
		WHERE linkfy_id = $1
		ORDER BY display_order ASC
	`
	log.Printf("Executing query: %s with linkfy_id: %s", query, linkfy_id)

	rows, err := l.db.Query(query, linkfy_id)
	if err != nil {
		log.Printf("Error querying links: %v", err)
		return nil, customerror.NewCustomError(err, "Failed to query links", 500)
	}
	defer rows.Close()

	links := []*LinkfyLink{}
	for rows.Next() {
		link := &LinkfyLink{}
		var id uuid.UUID
		err := rows.Scan(&id, &link.Name, &link.NameURL, &link.IconsURL, &link.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, customerror.NewCustomError(err, "Failed to scan link row", 500)
		}
		link.ID = id

		log.Printf("Retrieved link: %+v", link)
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, customerror.NewCustomError(err, "Error iterating over rows", 500)
	}

	log.Printf("Successfully retrieved %d links for linkfy_id: %s", len(links), linkfy_id)
	return links, nil
}
