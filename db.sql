CREATE SCHEMA "qon";

DROP TABLE IF EXISTS "qon"."app";
DROP TABLE IF EXISTS "qon"."trigger";

-- CREATE EXTENSION POSTGIS

-- App Table
-- CREATE TABLE "qon"."app" (
--    "app_id" CHARACTER(128) PRIMARY KEY
-- );

-- Qon Trigger Table
CREATE TABLE "qon"."trigger" (
    "id" BIGSERIAL PRIMARY KEY,
    "app_id" CHARACTER(128),    -- REFERENCES "qon"."app" (app_id), -- identifies the group
    "identifier" TEXT NOT NULL, -- resource id
    "coords" GEOGRAPHY (POINT, 4326) NOT NULL,
    "trigger_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'utc'),
    "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'utc'),
    "expires_at" TIMESTAMP DEFAULT 'epoch' -- not using interval b
);

-- Create spatial index
CREATE INDEX "qon_trigger_gix"
ON "qon"."trigger"
USING gist
(coords);
