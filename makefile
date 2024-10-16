SHELL := /bin/bash
BINARY_NAME := YfApi
SERVICE_NAME := yf-api
SSH_USER := root
DEV_SERVER := 192.168.77.112
REMOTE_PATH := /data/$(SERVICE_NAME)
SYSTEMD_CONF_NAME := $(SERVICE_NAME).service
SYSTEMD_CONF_PATH := /etc/systemd/system
#max编译
mac-scp: build_mac_to_linux
	@echo "Deploying to server: $(DEV_SERVER)"
	@echo "Remove old binary file"
	ssh $(SSH_USER)@$(DEV_SERVER) "rm -f $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Sending new binary file to dev server"
	scp $(BINARY_NAME) $(SSH_USER)@$(DEV_SERVER):$(REMOTE_PATH)
	@echo "Restart service"
	ssh $(SSH_USER)@$(DEV_SERVER) "systemctl restart $(SERVICE_NAME)"
	@echo "Done!"
build_mac_to_linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)
#windows编译
deploy_prod_windows_to_linux: build_windows_to_linux
	@echo "Deploying to server: $(DEV_SERVER)"
	@echo "Remove old binary file"
	ssh $(SSH_USER)@$(DEV_SERVER) "rm -f $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Sending new binary file to dev server"
	scp $(BINARY_NAME) $(SSH_USER)@$(DEV_SERVER):$(REMOTE_PATH)
	ssh $(SSH_USER)@$(DEV_SERVER) "chmod +x $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Restart service"
	ssh $(SSH_USER)@$(DEV_SERVER) "systemctl restart $(SERVICE_NAME)"
	@echo "Done!"
build_windows_to_linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o $(BINARY_NAME)

.PHONY: mac-scp build_mac_to_linux deploy_prod_windows_to_linux


swag:
	swag init --parseDependency --parseDepth=6 --exclude ./core,./internal/,./util/  -o ./docs

dev:
	go run main.go  -f ./config/config.yaml

.PHONY: swag dev
