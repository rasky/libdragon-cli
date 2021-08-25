BUILD_DIR=build
include n64.mk

src = main.c

all: game.z64

game.dfs: $(wildcard filesystem/*)
game.z64: N64_ROM_TITLE="Game"
game.z64: game.dfs 
game.elf: $(src:%.c=$(BUILD_DIR)/%.o)

clean:
	rm -f $(BUILD_DIR)/* game.dfs game.elf game.z64

-include $(wildcard $(BUILD_DIR)/*.d)

.PHONY: all clean
