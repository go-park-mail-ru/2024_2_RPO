-- Роль для просмотра продовой базы дистанционно
CREATE ROLE "3kybika" WITH LOGIN;

GRANT CONNECT ON DATABASE current_database() TO "3kybika";
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "3kybika";
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT SELECT ON TABLES TO "3kybika";

-- Сервисная роль для работы приложения с базой данных
CREATE ROLE tarasovxx WITH LOGIN;

GRANT CONNECT ON DATABASE current_database() TO tarasovxx;
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
