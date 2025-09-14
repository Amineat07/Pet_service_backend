-- name: CreateUser :one
INSERT INTO public.users (firstname, lastname, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, firstname, lastname, email;

-- name: GetUserByEmail :one
SELECT id, firstname, lastname, email, password
FROM users
WHERE email = $1;

-- name: CheckEmail :one
SELECT email FROM public.users WHERE email = $1;





