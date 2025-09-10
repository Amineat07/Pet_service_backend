CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  firstname text      NOT NULL,
  lastname  text      NOT NULL,
  email     text      NOT NULL,
  password  text       NOT NULL
);