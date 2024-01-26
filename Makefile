SOURCES := $(wildcard *.go)

bin/checkmark: $(SOURCES)
	go build -o bin/checkmark .

clean:
	rm -rf bin

.PHONY: default
