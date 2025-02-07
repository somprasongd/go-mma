CREATE TABLE public.customers (
	id serial4 NOT NULL,
	email text NOT NULL,
	credit_limit int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT customers_pkey PRIMARY KEY (id),
	CONSTRAINT customers_unique UNIQUE (email)
);