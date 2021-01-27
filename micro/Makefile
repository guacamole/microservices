check_install:
	which swagger || go -u get github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger.yaml --scan-models

