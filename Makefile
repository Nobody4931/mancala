SRC_FILES := main.go game.go minimax.go
OUT_DIR := bin
OUT_FILE := $(OUT_DIR)/mancala

RMDIR := rm -rfd

ifeq ($(OS),Windows_NT)
	OUT_FILE := $(OUT_FILE).exe
	RMDIR := rmdir /Q /S
endif

all: build

build: $(SRC_FILES)
	go build -o $(OUT_FILE) $^

clean:
	$(RMDIR) $(OUT_DIR)
