goose -dir internal/migrate/postgres/migrations postgres "postgresql://postgres:postgres@127.0.0.1:5436/mersi?sslmode=disable" down
goose -dir internal/migrate/postgres/migrations create new_accepted_table sql
scp -r ./dist administrator@route:/home/administrator/apps/mersi
npx vite-bundle-visualizer

для realm я решил делать динамические формы (для создания новых инструментов)

export DOCKER_API_VERSION=1.44
unset DOCKER_API_VERSION
