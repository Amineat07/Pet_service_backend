CREATE TABLE public.users (
  id               BIGSERIAL PRIMARY KEY,
  firstname        TEXT       NOT NULL,
  lastname         TEXT       NOT NULL,
  email            TEXT       UNIQUE NOT NULL,
  isCustomer       BOOLEAN    NOT NULL,
  isServiceProvider BOOLEAN   NOT NULL,
  isAdmin BOOLEAN NOT NULL,
  password         TEXT       NOT NULL
);

CREATE TABLE public.services (
    provider_id BIGINT PRIMARY KEY REFERENCES public.users(id),
    pet_sitting BOOLEAN NOT NULL DEFAULT false,
    dog_walking BOOLEAN NOT NULL DEFAULT false,
    pet_day_care BOOLEAN NOT NULL DEFAULT false,
    pet_grooming BOOLEAN NOT NULL DEFAULT false,
    pet_training BOOLEAN NOT NULL DEFAULT false,
    pet_massage BOOLEAN NOT NULL DEFAULT false
);



