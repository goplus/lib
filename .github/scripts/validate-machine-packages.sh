#!/bin/bash

# Machine package validation script
# Validates machine package compilation compatibility across platforms based on TinyGo smoketest targets

set -e

# Statistics variables
TOTAL_TARGETS=0
PASSED_TARGETS=0
FAILED_TARGETS=0
FAILED_LIST=()

# Test case directory
TEST_DIR="_demo/emptycheck"
TEST_FILE="$TEST_DIR/main.go"

echo "🚀 Starting Machine package compilation compatibility validation"
echo "Test file: $TEST_FILE"
echo ""

# Create test directory and file
mkdir -p "$TEST_DIR"
cat > "$TEST_FILE" << 'EOF'
package main

import (
    _ "github.com/goplus/lib/emb/machine"
)

//export handleHardFault
func handleHardFault() {}

//export handleInterrupt
func handleInterrupt() {}

//export Reset_Handler
func Reset_Handler() {}

func main(){

}
EOF

echo "Created test file: $TEST_FILE"

# Function to validate a single target
validate_target() {
    local target=$1
    local description=$2
    
    TOTAL_TARGETS=$((TOTAL_TARGETS + 1))
    
    echo -n "[$TOTAL_TARGETS] Testing $target"
    if [ -n "$description" ]; then
        echo -n " ($description)"
    fi
    echo -n "... "
    
    # Run llgo build test
    if llgo build -tags nogc -target="$target" "$TEST_FILE" 2>&1; then
        echo "✅ PASS"
        PASSED_TARGETS=$((PASSED_TARGETS + 1))
        return 0
    else
        echo "❌ FAIL"
        FAILED_TARGETS=$((FAILED_TARGETS + 1))
        FAILED_LIST+=("$target")
        return 1
    fi
}

echo "📋 Starting target validation one by one..."
echo ""

# Define all test targets list
TARGETS=(
    "pca10040-s132v6"
    "microbit"
    "microbit-s110v8"
    "microbit-v2"
    "microbit-v2-s113v7"
    "microbit-v2-s140v7"
    "nrf52840-mdk"
    "btt-skr-pico"
    "pca10031"
    "reelboard"
    "pca10056"
    "pca10059"
    "bluemicro840"
    "itsybitsy-m0"
    "feather-m0"
    "trinket-m0"
    "gemma-m0"
    "circuitplay-express"
    "circuitplay-bluefruit"
    "clue-alpha"
    # "gameboy-advance"
    # some picolibc symbol missed

    "grandcentral-m4"
    "itsybitsy-m4"
    "feather-m4"
    "matrixportal-m4"
    "pybadge"
    "metro-m4-airlift"
    "pyportal"
    "particle-argon"
    "particle-boron"
    "particle-xenon"
    "pinetime"
    "x9pro"
    "pca10056-s140v7"
    "pca10059-s140v7"
    "reelboard-s140v7"
    "wioterminal"
    "pygamer"
    "xiao"
    "rak4631"
    "feather-nrf52840"
    "feather-nrf52840-sense"
    "itsybitsy-nrf52840"
    "qtpy"
    "teensy41"
    "teensy40"
    "teensy36"
    "p1am-100"
    "atsame54-xpro"
    "feather-m4-can"
    "arduino-nano33"
    "arduino-mkrwifi1010"
    "pico"
    "nano-33-ble"
    "nano-rp2040"
    "feather-rp2040"
    "qtpy-rp2040"
    "kb2040"
    "macropad-rp2040"
    "badger2040"
    "badger2040-w"
    "tufty2040"
    "thingplus-rp2040"
    "xiao-rp2040"
    "waveshare-rp2040-zero"
    "challenger-rp2040"
    "trinkey-qt2040"
    "gopher-badge"
    "ae-rp2040"
    "thumby"
    "pico2"
    "tiny2350"
    "pico-plus2"
    # "metro-rp2350"
    # TODO: llgo sync this 0.39.0 target
    "waveshare-rp2040-tiny"
    "nrf52840-s140v6-uf2-generic"
    "bluepill"
    "feather-stm32f405"
    "lgt92"
    "nucleo-f103rb"
    "nucleo-f722ze"
    "nucleo-l031k6"
    "nucleo-l432kc"
    "nucleo-l476rg"
    "nucleo-l552ze"
    "nucleo-wl55jc"
    "stm32f4disco"
    "stm32f4disco-1"
    "stm32f469disco"
    "lorae5"
    "swan"
    "mksnanov3"
    # "stm32l0x1"
    # TODO: llgo sync this 0.39.0 target

    # "atmega328pb"
    # "atmega1284p"
    # "arduino"
    # "arduino-leonardo"
    # "arduino-mega1280"
    # "arduino-nano"
    # "attiny1616"
    # "digispark"
    # error: ran out of registers during register allocation
    # 2 errors generated.
    # panic: export object of internal/abi failed: exit status 1


    
    # "esp32-mini32"
    # ld.lld: error: section 'text' will not fit in region 'iram_seg': overflowed by 13157 bytes
    # binaray too large https://github.com/goplus/llgo/issues/1317

    "esp32c3-supermini"
    # "nodemcu"
    # compile picolibcnewlib/libc/tinystdio/puts.c fail

    "esp-c3-32s-kit"
    "qtpy-esp32c3"
    "m5stamp-c3"
    "xiao-esp32c3"
    "esp32-c3-devkit-rust-1"
    "esp32c3-12f"
    "makerfabs-esp32c3spi35"

    # "hifive1b"
    # "maixbit"
    # "tkey" 
    # picolibc symbol notfound like vfprintf,stdout,srandom ...

    "elecrow-rp2040"
    "elecrow-rp2350"
    "hw-651"
    "hw-651-s110v8"
    "wasm"
    "wasm-unknown"
)

# Execute validation in loop
for target in "${TARGETS[@]}"; do
    validate_target "$target" || true
done

# ==================== Final Statistics ====================
echo "📊 Validation Results Statistics"
echo "========================================"
echo "Total test targets: $TOTAL_TARGETS"
echo "Passed count: $PASSED_TARGETS"
echo "Failed count: $FAILED_TARGETS"

if [ $FAILED_TARGETS -gt 0 ]; then
    echo ""
    echo "❌ Failed targets:"
    for target in "${FAILED_LIST[@]}"; do
        echo "  - $target"
    done
    echo ""
    echo "⚠️  There are compilation failures, please check related issues"
    exit 1
else
    echo ""
    echo "🎉 All targets validated successfully! Machine package compatibility is good"
    exit 0
fi
