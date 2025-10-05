-- name: CreateUser :one
INSERT INTO public.users (firstname, lastname, email, password, isCustomer, isServiceProvider,isAdmin)
VALUES ($1, $2, $3, $4,$5,$6, $7)
RETURNING id, firstname, lastname, email, isCustomer,isServiceProvider,isAdmin;

-- name: GetUserByEmail :one
SELECT id, firstname, lastname, email, password,isCustomer,isServiceProvider,isAdmin
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT *
FROM public.users
WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM public.users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;;

-- name: UpdateUser :exec
UPDATE public.users SET firstname=$2,lastname=$3,email=$4,password=$5,updated_at = now()
WHERE id=$1;

-- name: DeleteUser :exec
DELETE FROM public.users
WHERE id = $1;

-- name: DeleteServices :exec
DELETE FROM public.services WHERE provider_id = $1;

-- name: GetProviders :many
SELECT * FROM users WHERE isServiceProvider = true LIMIT $1 OFFSET $2;

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

-- name: GetServices :many
SELECT * FROM public.services;

-- name: GetServiceByProviderID :one
SELECT * FROM public.services WHERE provider_id = $1;

-- name: UpdateServices :exec
UPDATE public.services SET pet_sitting = $2 ,dog_walking= $3,pet_day_care=$4,pet_grooming=$5,pet_training=$6,pet_massage=$7
WHERE provider_id =$1; 

-- name: MakeReservation :one
INSERT INTO public.booked_service (customer_id, provider_id, service_type, start_time, end_time)
SELECT $1, $2, $3, $4, $5
WHERE NOT EXISTS (
    SELECT 1
    FROM public.booked_service
    WHERE customer_id = $1
      AND service_type = $3
      AND start_time < $5
      AND end_time > $4
)
RETURNING id, customer_id, provider_id, service_type, start_time, end_time;







