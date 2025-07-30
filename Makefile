tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down

mocks ::
	@echo "Generating mocks..."
	@which mockgen > /dev/null || (echo "mockgen not found, please install it using 'go install go.uber.org/mock/mockgen@latest'" && exit 1)
	mockgen -source=app/repositories/products_repository.go -destination=app/repositories/mocks/products_mock.go -package=repo_mock

test-ginkgo ::
	@echo "Running Ginkgo tests..."
	@which ginkgo > /dev/null || (echo "Ginkgo not found, please install it using 'go install github.com/onsi/ginkgo/v2/ginkgo@latest'" && exit 1)
	@ginkgo -v -r