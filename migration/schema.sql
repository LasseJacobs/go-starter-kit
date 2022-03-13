
CREATE SCHEMA IF NOT EXISTS content;

CREATE TABLE IF NOT EXISTS stories (
   id           varchar(32) CONSTRAINT pkey PRIMARY KEY,
   title        varchar(40) NOT NULL,
   author       varchar(40) DEFAULT 'unknown',
   votes        int DEFAULT 0,
   url          varchar(2048) NOT NULL,
   origin_date  timestamp DEFAULT current_timestamp,
   transaction_date timestamp DEFAULT current_timestamp
);