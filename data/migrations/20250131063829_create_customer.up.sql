create table public.customers (
  id serial PRIMARY KEY,
  email text NOT NULL,
  credit_limit integer NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);