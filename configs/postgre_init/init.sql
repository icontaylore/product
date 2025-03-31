ALTER SYSTEM SET wal_level = logical;

CREATE PUBLICATION dbz_publication FOR ALL TABLES;

SELECT pg_create_logical_replication_slot('debezium_slot', 'pgoutput');