.PHONY: all clean rtcgw worker

all: rtcgw worker

clean:
	rm -f rtcgw-go workers/workers

rtcgw:
	go build

worker:
	go build -o workers/workers ./workers

run-server: rtcgw
	./rtcgw

run-worker: worker
	./workers/workers