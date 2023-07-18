PROJECT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

TARGET_DIR=target

BINARY_DIR=${TARGET_DIR}/bin
BINARY_PATH=${BINARY_DIR}/app
TEST_BINARY_PATH=${BINARY_DIR}/test-app

TEST_DIR=${TARGET_DIR}/test
TEST_REPORT_DIR=${TEST_DIR}/reports
TEST_API_REPORT=${TEST_REPORT_DIR}/api-tests-report.out

pwd:
	@echo $(PROJECT_DIR)

clean: ## cleans target dir
	@rm -rf $(TARGET_DIR)

bin: ## build the binary for local run
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) cmd/service/main.go

run: bin ## runs the binary
	@$(BINARY_PATH)

test/bin: ## build the binary for local run
	@mkdir -p $(BINARY_DIR)
	@go test -c cmd/test-service/main_test.go -tags=api_test -cover -covermode=count -o $(TEST_BINARY_PATH) -coverpkg=./...; \

test/api: test/bin ## runs the binary for local api testing and outputs cover result
	@mkdir -p $(TEST_REPORT_DIR)
# 	start the test binary and say where it should output the coverage file
	@$(TEST_BINARY_PATH) -test.coverprofile $(TEST_API_REPORT)
	@echo "\n\ncoverage report\n"
#	display the coverage	
	@go tool cover -func=$(TEST_API_REPORT)