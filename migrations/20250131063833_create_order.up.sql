CREATE TABLE public.orders (
	id serial4 NOT NULL,
	customer_id int4 NOT NULL,
	order_total int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	canceled_at timestamp NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id),
	CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES public.customers(id)
);