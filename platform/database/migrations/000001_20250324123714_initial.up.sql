-- public.wallets definition

-- Drop table

-- DROP TABLE public.wallets;

CREATE TABLE public.wallets (
	id varchar(50) NOT NULL,
	balance numeric(20, 2) NOT NULL DEFAULT 0,
	CONSTRAINT wallets_pkey PRIMARY KEY (id)
);

-- public.transactions definition

-- Drop table

-- DROP TABLE public.transactions;

CREATE TABLE public.transactions (
	id serial4 NOT NULL,
	wallet_id varchar(50) NOT NULL,
	amount numeric(20, 2) NOT NULL,
	"type" varchar(10) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT transactions_pkey PRIMARY KEY (id),
	CONSTRAINT transactions_type_check CHECK (((type)::text = ANY ((ARRAY['credit'::character varying, 'debit'::character varying])::text[])))
);

-- public.transactions foreign keys

ALTER TABLE public.transactions ADD CONSTRAINT transactions_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES public.wallets(id);   
