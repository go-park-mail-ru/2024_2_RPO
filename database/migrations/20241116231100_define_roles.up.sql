-- Роль для просмотра продовой базы дистанционно
DO
$$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE  rolname = '3kybika') THEN

        REASSIGN OWNED BY "3kybika" TO postgres;
        DROP OWNED BY "3kybika";
    END IF;
END
$$;
REASSIGN OWNED BY "3kybika" TO postgres;
DROP OWNED BY "3kybika";
DROP ROLE IF EXISTS "3kybika";
CREATE ROLE "3kybika" WITH LOGIN;

DO $$
BEGIN
    EXECUTE format('GRANT CONNECT ON DATABASE %I TO %I', current_database(), '3kybika');
END
$$;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "3kybika";
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT SELECT ON TABLES TO "3kybika";

-- Сервисная роль для работы приложения с базой данных
DO
$$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE  rolname = 'tarasovxx') THEN

        REASSIGN OWNED BY tarasovxx TO postgres;
        DROP OWNED BY tarasovxx;
    END IF;
END
$$;
DROP ROLE IF EXISTS tarasovxx;
CREATE ROLE tarasovxx WITH LOGIN;

DO $$
BEGIN
    EXECUTE format('GRANT CONNECT ON DATABASE %I TO %I', current_database(), 'tarasovxx');
END
$$;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tarasovxx;
DO $$ DECLARE
    seq RECORD;
BEGIN
    FOR seq IN SELECT c.relname, n.nspname
                FROM pg_class c
                JOIN pg_namespace n ON c.relnamespace = n.oid
                WHERE c.relkind = 'S'
LOOP
    EXECUTE format('GRANT ALL PRIVILEGES ON SEQUENCE %I.%I TO tarasovxx;', seq.nspname, seq.relname);
END LOOP;
END $$;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO tarasovxx;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO tarasovxx;
