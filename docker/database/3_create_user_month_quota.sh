#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "arvan_queue" <<-EOSQL

  CREATE TABLE public.users_month_quota
  (
      id         text NOT NULL,
      user_id    text NOT NULL
      quota      bigint NOT NULL
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
      CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
        REFERENCES users(id)
  );

  CREATE TRIGGER created_at_customers_trgr BEFORE UPDATE ON public.users_month_quota FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_users_trgr BEFORE UPDATE ON public.users_month_quota FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

EOSQL
