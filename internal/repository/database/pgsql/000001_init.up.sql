CREATE TABLE IF NOT EXISTS fioapi (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL,
     surname TEXT NOT NULL,
     patronymic TEXT,
     age INTEGER,
     gender TEXT,
     nationality JSON
);

CREATE INDEX IF NOT EXISTS idx_name_surname ON fioapi (name, surname);