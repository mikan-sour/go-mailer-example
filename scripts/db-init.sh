#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USERNAME" \
    --dbname="$POSTGRES_DB"<<-EOSQL

    CREATE TABLE IF NOT EXISTS profanity 
      (
         id SERIAL PRIMARY KEY,
         word VARCHAR(30) NOT NULL, 
         region VARCHAR(2) NOT NULL,
         active BOOL NOT NULL DEFAULT TRUE,
         created timestamp without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'), 
         created_by VARCHAR(30) NOT NULL DEFAULT 'system',
         updated timestamp without time zone,
         updated_by VARCHAR(30) NOT NULL DEFAULT 'system'
      );
    
    COPY profanity(word, region)
    FROM '/var/lib/postgresql/data/bad-words-en.csv'
    DELIMITER ','
    CSV HEADER;

EOSQL