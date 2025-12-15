CREATE TABLE public.users (
  id               BIGSERIAL PRIMARY KEY,
  firstname        TEXT       NOT NULL,
  lastname         TEXT       NOT NULL,
  email            TEXT       UNIQUE NOT NULL,
  isCustomer       BOOLEAN    NOT NULL,
  isServiceProvider BOOLEAN   NOT NULL,
  isAdmin BOOLEAN NOT NULL,
  password         TEXT       NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
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

ALTER TABLE public.services
ADD COLUMN pet_sitting_price NUMERIC(10,2) DEFAULT 0,
ADD COLUMN dog_walking_price NUMERIC(10,2) DEFAULT 0,
ADD COLUMN pet_day_care_price NUMERIC(10,2) DEFAULT 0,
ADD COLUMN pet_grooming_price NUMERIC(10,2) DEFAULT 0,
ADD COLUMN pet_training_price NUMERIC(10,2) DEFAULT 0,
ADD COLUMN pet_massage_price NUMERIC(10,2) DEFAULT 0;


CREATE TABLE public.booked_service (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES public.users(id),
    provider_id BIGINT NOT NULL REFERENCES public.users(id),
    service_type TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


