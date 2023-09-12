AR_PROGRAMMER=
include ./config.mk

SRCS := $(shell find $(AR_SKETCH) -type f -name '*.ino' -or -name '*.h' -or -name '*.c' -or -name '*.cpp')

upload: build/$(AR_SKETCH).ino.hex
	arduino-cli upload -b $(AR_FQBN) -p $(AR_PORT) $(AR_PROGRAMMER) --input-dir "$(CURDIR)/build"

compile: build/$(AR_SKETCH).ino.hex

monitor:
	arduino-cli monitor -p $(AR_PORT) -c $(AR_MONITOR_BAUD)

details:
	arduino-cli board details --fqbn=$(AR_FQBN)

clean:
	rm -rf build

build/$(AR_SKETCH).ino.hex: $(SRCS)
	arduino-cli compile -b $(AR_FQBN) --build-path "$(CURDIR)/build" $(AR_SKETCH)

.PHONY: upload compile serial clean

