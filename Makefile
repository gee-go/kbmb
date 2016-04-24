start:
	docker-compose up --build -d
	docker-compose scale worker=3

restart:
	docker-compose up --no-deps --build -d worker

stop:
	docker-compose down