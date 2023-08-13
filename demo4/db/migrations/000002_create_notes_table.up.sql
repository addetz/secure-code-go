BEGIN;
CREATE TABLE IF NOT EXISTS notes
(
   id VARCHAR (50) PRIMARY KEY,
   username VARCHAR(50)  REFERENCES users (username),
   noteText VARCHAR (500) NOT NULL
);
COMMIT;