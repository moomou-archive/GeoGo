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
Uses Postgres database for geolocation. Schema is located in db.sql.

Methods
=====
GET: Returns a list of ids near the lat lon within the radius.
/trigger?lat=x&lon=y&radius=z&unit=m|meter

POST: Add a list of ids to the geolocation
/trigger?lat=x&lon=y
{
    ids: [1, 2, 3]
}

DELETE: Remove a list of ids from the geolocation
/trigger?lat=x&lon=y&ids=[x,y,z]

Note
====
First ever Go project :)
