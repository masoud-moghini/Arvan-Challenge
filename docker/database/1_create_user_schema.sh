#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "arvan_queue" <<-EOSQL

  CREATE TABLE public.users
  (
      id         text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_customers_trgr BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_users_trgr BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  GRANT USAGE ON SCHEMA public TO arvan_db_admin;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO arvan_db_admin;
EOSQL
