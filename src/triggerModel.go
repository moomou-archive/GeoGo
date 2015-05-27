package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type Trigger struct {
	db *sql.DB
}

type trigger struct {
	AppId      string     `json:"appId"`
	Identifier string     `json:"identifier"`
	CreatedAt  *time.Time `json:"createdAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	Coords     []string   `json:"coords"` //lat, lon
}

func (t *Trigger) add(triggers *[]trigger) error {
	query := `INSERT INTO qon.trigger (app_id, identifier, coords, expires_at) VALUES `

	for _, v := range *triggers {
		expiresAt := "NULL"
		if v.ExpiresAt != nil {
			expiresAt = fmt.Sprintf(`'%s'`, v.ExpiresAt.Format(time.RFC3339))
		}
		query += fmt.Sprintf(`('%s', '%s', ST_GeographyFromText('SRID=4326;POINT(%s)'), %s)`,
			strings.TrimSpace(v.AppId),
			strings.TrimSpace(v.Identifier),
			fmt.Sprintf("%s %s", strings.TrimSpace(v.Coords[1]), strings.TrimSpace(v.Coords[0])),
			expiresAt,
		) + ",\n"
	}

	query = strings.TrimRight(strings.TrimSpace(query), ",")
	log.Printf("ADD Query: ", query)

	if _, err := t.db.Query(query); err != nil {
		log.Printf("Err:", err)
		return err
	}

	return nil
}

func (t *Trigger) remove(appId string, identifier string) error {
	query := `DELETE FROM qon.trigger WHERE app_id=$1 AND identifier=$2;`

	if _, err := t.db.Query(query, appId, identifier); err != nil {
		return err
	}

	return nil
}

func (t *Trigger) updateIdentifier(appId string, oldIdentifier string, newIdentifier string) error {
	query := `UPDATE qon.trigger SET identifier=$1 WHERE app_id=$2 AND identifier=$3;`

	if _, err := t.db.Query(query, newIdentifier, appId, oldIdentifier); err != nil {
		return err
	}

	return nil
}

func (t *Trigger) findNearBy(triggerInfo *trigger, radius int64, unit string) ([]trigger, error) {

	radiusMeter := radius

	switch {
	case unit == "km":
		{
			radiusMeter *= 1000
		}
	}

	// PostGIS expects lon lat, in that order
	query :=
		`SELECT DISTINCT ON (app_id, identifier)
                app_id,
                identifier,
                ST_X(coords::geometry),
                ST_Y(coords::geometry),
                created_at,
                expires_at
        FROM qon.trigger
            WHERE ST_DWITHIN(
                Geography(coords),
                Geography(ST_MakePoint($1, $2)),
                $3
            )
            AND (expires_at IS NULL OR expires_at >= 'now()' OR
			     date_trunc('hour', expires_at) = date_trunc('hour', TIMESTAMP 'epoch'))
            AND app_id = $4
        `
	// Optinally filter for a specific item.
	if triggerInfo.Identifier != "" {
		query += fmt.Sprintf("AND identifier = '%s';", triggerInfo.Identifier)
	} else {
		query += ";"
	}

	qCoords := (*triggerInfo).Coords
	qlon, qlat := qCoords[1], qCoords[0]

	log.Printf(query)
	log.Printf("Query Params", qlon, qlat, strconv.FormatInt(radiusMeter, 10), triggerInfo.AppId)

	rows, err := t.db.Query(query,
		qlon,
		qlat,
		strconv.FormatInt(radiusMeter, 10),
		triggerInfo.AppId,
	)

	if err != nil {
		return nil, err
	}

	var (
		appId       string
		identifier  string
		lat         string
		lon         string
		createdAt   *time.Time
		expiresAt   *time.Time
		ntCreatedAt pq.NullTime
		ntExpiresAt pq.NullTime
	)

	resultTrigger := make([]trigger, 100)[1:2]
	coords := make([]string, 2)
	rowLength := 1

	for rows.Next() {
		// internal always lon = ST_X, lat = ST_Y
		err := rows.Scan(&appId, &identifier, &lon, &lat, &ntCreatedAt, &ntExpiresAt)

		// external coords always lat and lon, in that order
		coords = []string{lat, lon}

		if err != nil {
			return nil, err
		}

		if ntCreatedAt.Valid {
			createdAt = &ntCreatedAt.Time
		} else {
			createdAt = nil
		}

		if ntExpiresAt.Valid {
			expiresAt = &ntExpiresAt.Time
		} else {
			expiresAt = nil
		}

		resultTrigger = append(resultTrigger, trigger{
			AppId:      appId,
			Identifier: identifier,
			CreatedAt:  createdAt,
			ExpiresAt:  expiresAt,
			Coords:     coords,
		})

		rowLength += 1
	}

	return resultTrigger[1:rowLength], nil
}

func newTrigger(db *sql.DB) *Trigger {
	return &Trigger{db: db}
}
