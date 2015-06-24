GEOGO
=====

What is it
==========
A simple Go based web service that associates location (currently lat, lon)
with a user specified unique identifier that maps to other services.

For example, say you run a social network where every user has an opaque user id.

You can use GEOGO to store latitude and longitude information for each user id and use GEOGO to find users
near each other.

Listens on port 3003.

Dependencies
============
Uses Postgres + PostGIS database for geolocation. Schema is located in db.sql.

If you are on Mac, http://postgresapp.com/ is the easiest way to get started with PostGIS already installed.

Otherwise, checkout http://postgis.net/install for installation instructions.

Methods
=====
GET: Returns a list of ids near the lat lon within the radius.

/trigger?lat=x&lon=y&radius=z&unit=m|meter

DELETE: Remove a list of ids from the geolocation

/trigger?ids=[x,y,z]

POST: Add a list of ids to the geolocation

/trigger

    [
        {
            expiresAt: `null` or string - ISOString8601 of the expire time
            identifier: string - object id,
            appId: string - identifies application,
            coords: array - [latitude, longitude]
        },
        ...
    ]
