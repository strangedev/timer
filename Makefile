BINARY_NAME=worktimer
SERVICEFILE_NAME=worktimer.service

PHONY: clean hash

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} main.go

hash: build
	sha256sum worktimer > worktimer.sha256

clean:
	go clean

install:
	install -m 0755 -t /usr/local/bin ${BINARY_NAME}
	install -m 0644 -t /etc/systemd/user ${SERVICEFILE_NAME} 

uninstall:
	rm /usr/local/bin/${BINARY_NAME}
	-systemctl --user disable --now ${SERVICEFILE_NAME}
	rm /etc/systemd/user/${SERVICEFILE_NAME}

redeploy-service:
	systemctl --user daemon-reload
	systemctl --user enable --now worktimer
	systemctl --user restart worktimer
