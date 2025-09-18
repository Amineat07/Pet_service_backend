-- name: CreateUser :one
INSERT INTO public.users (firstname, lastname, email, password, isCustomer, isServiceProvider,isAdmin)
VALUES ($1, $2, $3, $4,$5,$6, $7)
RETURNING id, firstname, lastname, email, isCustomer,isServiceProvider,isAdmin;

-- name: GetUserByEmail :one
SELECT id, firstname, lastname, email, password,isCustomer,isServiceProvider,isAdmin
FROM users
WHERE email = $1;

-- name: GetRolebyID :one
SELECT id, isAdmin, isCustomer, isServiceProvider
FROM users
WHERE id = $1;


-- name: CheckEmail :one
SELECT email FROM public.users WHERE email = $1;


-- name: UpsertServices :one
INSERT INTO public.services (
    provider_id, pet_sitting, dog_walking, pet_day_care, pet_grooming, pet_training, pet_massage
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (provider_id) DO UPDATE
SET pet_sitting  = EXCLUDED.pet_sitting,
    dog_walking  = EXCLUDED.dog_walking,
    pet_day_care = EXCLUDED.pet_day_care,
    pet_grooming = EXCLUDED.pet_grooming,
    pet_training = EXCLUDED.pet_training,
    pet_massage  = EXCLUDED.pet_massage
RETURNING *;


-- name: GetServicesByProvider :one
SELECT provider_id, pet_sitting, dog_walking, pet_day_care, pet_grooming, pet_training, pet_massage
FROM public.services
WHERE provider_id = $1;





