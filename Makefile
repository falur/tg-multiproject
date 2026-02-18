.PHONY: build run test lint setup-server deploy deploy-local status logs restart

# Local development
build:
	go build -o bin/tg-multiproject ./cmd/bot

run:
	go run ./cmd/bot

test:
	go test -v ./...

lint:
	golangci-lint run ./...

# Server
setup-server:
	ansible-playbook -i deploy/ansible/inventory.ini deploy/ansible/setup.yml

deploy:
	ansible-playbook -i deploy/ansible/inventory.ini deploy/ansible/deploy.yml

deploy-local:
	ansible-playbook -i deploy/ansible/inventory.ini deploy/ansible/deploy.yml -e deploy_method=rsync

status:
	ssh bot-server "systemctl status tg-multiproject"

logs:
	ssh bot-server "journalctl -u tg-multiproject -f"

restart:
	ssh bot-server "sudo systemctl restart tg-multiproject"
