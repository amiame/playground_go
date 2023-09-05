example_usecase:
	go run example_usecase/main/main.go

doc:
	@docker build -t my-pkg-error-stack-golang:local -f Dockerfile.golang . &> logs.txt
	@docker run -d --rm -p 127.0.0.1:6060:6060 --name doc my-pkg-error-stack-golang:local godoc -http 0.0.0.0:6060 &> logs.txt
	@echo "open http://127.0.0.1:6060/pkg/myPkg/errors to see document"
	@rm logs.txt

clean:
	@docker container stop godoc &> errors.txt || true
	@docker image rm my-pkg-error-stack-golang:local &> errors.txt || true
	@rm errors.txt

.PHONY: example_usecase
