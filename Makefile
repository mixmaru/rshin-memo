build:
	cd ./cui_app && go build

exec: build
	cd ./cui_app && ./cui_app

clean:
	rm -f ./cui_app/cui_app

# docker-build:
# 	cd ./development && \
# 		docker-compose up -d --build
# 
# docker-clean:
# 	cd ./development && \
# 		docker-compose stop && docker-compose rm -f
# 
# docker-rebuild: docker-clean docker-build
# 
# test:
# 	cd ./development && docker-compose exec dev go test ./... -count=1 -cover -p 1
# 
# sql-migration-up-test:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate up -env='test'"
# 
# sql-migration-up-development:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate up"
# 
# sql-migration-up-production:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate up -env='production'"
# 
# sql-migration-down-test:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate down -env='test'"
# 
# sql-migration-down-development:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate down"
# 
# sql-migration-down-production:
# 	cd ./development && docker-compose exec dev bash -c "cd sql_migration/ && sql-migrate down -env='production'"

