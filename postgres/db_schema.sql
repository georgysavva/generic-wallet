CREATE TABLE public.accounts
(
    id text PRIMARY KEY NOT NULL,
    balance float DEFAULT 0 NOT NULL,
    currency text DEFAULT 'USD' NOT NULL
);

CREATE TABLE public.payments
(
    account text NOT NULL,
    to_account text,
    from_account text,
    amount float NOT NULL,
    direction text NOT NULL,
    CONSTRAINT payments_accounts_id_fk FOREIGN KEY (account) REFERENCES public.accounts (id) ON DELETE CASCADE,
    CONSTRAINT payments_accounts_id_fk_2 FOREIGN KEY (to_account) REFERENCES public.accounts (id) ON DELETE CASCADE,
    CONSTRAINT payments_accounts_id_fk_3 FOREIGN KEY (from_account) REFERENCES public.accounts (id) ON DELETE CASCADE
);