BINS=webby

UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
  GOOS ?= 'linux'
else
  GOOS ?= 'darwin'
endif


$(BINS): webby.go
	CGO_ENABLED=0 GOOS=$(GOOS) go build -ldflags="-s -w -X main.Version=`git log --pretty=format:'%h' -n 1`" -a -installsuffix cgo -o webby .


clean:
	rm -rf $(BINS)

dockerbuild:
	docker build -f Dockerfile-buildimage -t 024376647576.dkr.ecr.us-east-1.amazonaws.com/webby:`git log --pretty=format:'%h' -n 1` .
	docker tag 024376647576.dkr.ecr.us-east-1.amazonaws.com/webby:`git log --pretty=format:'%h' -n 1` 024376647576.dkr.ecr.us-east-1.amazonaws.com/webby:latest
	

dockerpush:
	docker push 024376647576.dkr.ecr.us-east-1.amazonaws.com/webby:`git log --pretty=format:'%h' -n 1`
	docker push 024376647576.dkr.ecr.us-east-1.amazonaws.com/webby:latest

all: clean $(BINS)
