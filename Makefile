SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS = -ec

PGDATA := ~/pg/pgdata
PGBIN := ~/pg/pgbin/bin

pg.reset_pg_data: pg.stop
	rm -rf $(PGDATA)/*
	$(PGBIN)/initdb -D $(PGDATA) -U postgres
	echo "log_error_verbosity = terse" >> $(PGDATA)/postgresql.conf
	echo "log_min_error_statement = 'FATAL'" >> $(PGDATA)/postgresql.conf


pg.stop:
	@echo "Stopping PostgreSQL..."
	$(PGBIN)/pg_ctl -D $(PGDATA) stop -m immediate || true

pg.start:
	@echo "Starting PostgreSQL..."
	$(PGBIN)/pg_ctl -D $(PGDATA) start

pg.start_in_foreground:
	@echo "Starting PostgreSQL in foreground..."
	$(PGBIN)/postgres -D $(PGDATA)

create_schema_and_seed_data:
	@echo "Applying schema..."
	$(PGBIN)/psql -U postgres -d postgres -f schema.sql
