#Docker postgres example
#https://stackoverflow.com/questions/46981073/how-to-run-postgres-in-a-docker-alpine-linux-container
#https://dotsandbrackets.com/docker-health-check-ru/
#https://habr.com/ru/articles/578744/


# # Разработчики официального образа PostgreSQL естественно предусмотрели этот момент и предоставили 
# нам специальную точку входа для инициализации базы данных - docker-entrypoint-initdb.d. Любые *.sql 
# или *.sh файлы в этом каталоге будут рассматриваться как скрипты для инициализации БД. Здесь есть несколько нюансов:
# # -- если БД уже была проинициализирована ранее, то никакие изменения к ней применяться не будут;
# # -- если в каталоге присутствует несколько файлов, то они будут отсортированы по имени с использованием 
# текущей локали (по умолчанию en_US.utf8).

FROM postgres:10.0-alpine


USER postgres

RUN chmod 0700 /var/lib/postgresql/data &&\
    initdb /var/lib/postgresql/data &&\
    echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf &&\
    pg_ctl start &&\
    psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'main'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE main" &&\
    psql -c "ALTER USER postgres WITH ENCRYPTED PASSWORD 'postgres';" &&\
    #create test users and Database autotest
    psql -c "CREATE DATABASE praktikum;" &&\
    #create test users for project
    psql -c "CREATE DATABASE short;" &&\
    psql -c "CREATE USER short WITH ENCRYPTED PASSWORD '1';" &&\
    psql -c "GRANT ALL PRIVILEGES ON DATABASE short TO short;" 
    
EXPOSE 5432

HEALTHCHECK --interval=5s --timeout=10s --retries=3 CMD curl -sS 127.0.0.1 || exit 1