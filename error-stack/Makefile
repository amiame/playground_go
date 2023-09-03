run:
	go run main/main.go

check-package-vulnerabilities:
	docker build -t error-stack:local .
	docker scout quickview error-stack:local
