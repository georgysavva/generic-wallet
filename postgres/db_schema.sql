CREATE TABLE public.accounts
(
    id text PRIMARY KEY NOT NULL,
    balance float DEFAULT 0 NOT NULL,
    currency text DEFAULT 'USD' NOT NULL
);

CREATE TABLE public.payments
(
    id serial PRIMARY KEY NOT NULL,
    account_id text NOT NULL,
    to_account_id text,
    from_account_id text,
    amount float NOT NULL,
    direction text NOT NULL,
    CONSTRAINT payments_accounts_id_fk FOREIGN KEY (account_id) REFERENCES public.accounts (id) ON DELETE CASCADE,
    CONSTRAINT payments_accounts_id_fk_2 FOREIGN KEY (to_account_id) REFERENCES public.accounts (id) ON DELETE CASCADE,
    CONSTRAINT payments_accounts_id_fk_3 FOREIGN KEY (from_account_id) REFERENCES public.accounts (id) ON DELETE CASCADE
);

INSERT INTO accounts
VALUES ('alice', 100.0, 'USD'),
       ('bob', 100.0, 'USD'),
       ('mark', 100.0, 'USD'),
       ('john', 100.0, 'USD'),
       ('kate_in_europe', 100.0, 'EUR');
