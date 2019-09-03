TARGET_DIR=./target
TESTER_DIR=tester
BIN_NAME=seatgeek-be-challenge

.PHONY: clean
clean:
	rm -rf $(TARGET_DIR)
	mkdir -p $(TARGET_DIR)/src


.PHONY: clean
build: 
	cd $(TESTER_DIR) &&	go build  -o ../$(TARGET_DIR)/$(BIN_NAME) .

.PHONY: package
package: clean build
	cp -a $(TESTER_DIR) $(TARGET_DIR)/src
	cp INSTRUCTIONS.md $(TARGET_DIR)
	cd $(TARGET_DIR) &&	tar cvzf ../$(BIN_NAME).tar.gz *

