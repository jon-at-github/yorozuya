GO_VERSION := 1.20
ARCH := amd64
GO_BINARY_LATEST := go$(GO_VERSION).linux-$(ARCH).tar.gz

setup: install-go init-go build

install-go:
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xvf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}.bashrc

upgrade-go:
	sudo rm -rf /usr/local/go
	wget "https://go.dev/dl/${GO_BINARY_LATEST}"
	sudo  tar -C /usr/local -xzf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

build:
	go build -o api cmd/main.go