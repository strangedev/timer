BINARY_NAME=worktimer
SERVICEFILE_NAME=worktimer.service

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} main.go

clean:
	go clean

install:
	install -m 0755 -t /usr/local/bin ${BINARY_NAME}
	install -m 0644 -t /etc/systemd/user ${SERVICEFILE_NAME} 

uninstall:
	rm /usr/local/bin/${BINARY_NAME}
	-systemctl --user disable --now ${SERVICEFILE_NAME}
	rm /etc/systemd/user/${SERVICEFILE_NAME}
