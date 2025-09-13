CREATE TABLE users (
  id        BIGSERIAL PRIMARY KEY,
  firstname TEXT       NOT NULL,
  lastname  TEXT       NOT NULL,
  email     TEXT       UNIQUE NOT NULL,
  password  TEXT       NOT NULL
);
