CREATE USER short WITH ENCRYPTED PASSWORD '1';
CREATE USER praktikum WITH ENCRYPTED PASSWORD 'praktikum';
CREATE DATABASE short;
CREATE DATABASE praktikum;
GRANT ALL PRIVILEGES ON DATABASE short TO short;
GRANT ALL PRIVILEGES ON DATABASE praktikum TO short;

-- need for migrations (issue https://github.com/golang-migrate/migrate/issues/826)
ALTER DATABASE short OWNER TO short;
ALTER DATABASE praktikum OWNER TO short;