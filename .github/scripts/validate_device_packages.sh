#!/bin/bash

# validate_device_packages.sh
# Validates device packages compilation with different LLGO targets
#
# Usage:
#   ./validate_device_packages.sh                      # Run all validations
#   ./validate_device_packages.sh device/arm:cortex-m-qemu  # Run single validation

set -e

# Define ignore list - targets not supported by LLGO
# Reference: https://github.com/goplus/llgo/blob/18e036568de0b8d3f1284b975402e44a5ab7c248/_demo/embed/targetsbuild/build.sh
ignore_list=(
	"atmega1280"
	"atmega2560"
	"atmega328p"
	"atmega32u4"
	"attiny85"
	"fe310"
	"k210"
	"riscv32"
	"riscv64"
	"rp2040"
)

# Configuration: device packages and their validation targets
# Format: "device_package:target1,target2,..."
VALIDATION_TARGETS=(
    # arm
    "device/arm:cortex-m-qemu"
    
    # arm64
    "device/arm64:esp32"
    
    # avr
    "device/avr:arduino"
    
    # esp
    "device/esp:esp32,esp8266"
    
    # gba
    "device/gba:gameboy-advance"
    
    # kendryte
    # riscv64 current not support at llgo
    # "device/kendryte"
    
    # nrf
    "device/nrf:nrf52840,nrf51"
    
    # nxp
    "device/nxp:teensy40,teensy36"
    
    # renesas
    # "device/renesas"

    # riscv picolibc fail
    # "device/riscv:riscv-qemu"
    
    # rp
    "device/rp:rp2350"
    
    # sam
    "device/sam:arduino-zero,atsamd51j19a"
    
    # sifive
    # riscv32 current not support at llgo
    # "device/sifive"
    
    # stm32 - STM32 series microcontrollers
    "device/stm32:stm32f4disco,nucleo-f103rb"
    
    # riscv32 current not support at llgo
    # "device/tkey"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage information
usage() {
    echo "Usage: $0"
    echo ""
    echo "Validates device package compilation with LLGO targets"
    echo "Requires 'llgo' command to be available in PATH"
    exit 1
}

# Create test file for a device package
create_test_file() {
    local device_package="$1"
    local target="$2"
    local filename="test_${target//[-.]/_}.go"
    
    cat > "$filename" << EOF
package main

import _ "github.com/goplus/lib/emb/$device_package"

// Placeholder functions to resolve undefined symbols from LLGO targets
// Similar to github.com/goplus/llgo/_demo/embed/targetsbuild/C/c.go

//export handleHardFault
func handleHardFault() {}

//export handleInterrupt  
func handleInterrupt() {}

//export Reset_Handler
func Reset_Handler() {}

func main() {}
EOF
    
    echo "$filename"
}

# Check if target is in ignore list
check_ignore_list() {
    local target="$1"
    for ignored in "${ignore_list[@]}"; do
        if [[ "$target" == "$ignored" ]]; then
            return 0  # Found in ignore list
        fi
    done
    return 1  # Not in ignore list
}

# Validate single target
validate_target() {
    local device_package="$1"
    local target="$2"
    local current="$3"
    local total="$4"
    
    # Check if target is in ignore list
    if check_ignore_list "$target"; then
        printf "${YELLOW}[$current/$total]${NC} ⚠️  Skipping $device_package with target $target (not supported by LLGO)\n"
        echo "    📋 Target '$target' is in LLGO ignore list"
        echo "    🔗 Reference: https://github.com/goplus/llgo/blob/18e036568de0b8d3f1284b975402e44a5ab7c248/_demo/embed/targetsbuild/build.sh"
        return 0
    fi
    
    printf "${BLUE}[$current/$total]${NC} 🔄 Validating $device_package with target $target..."
    
    # Create test file
    local test_file
    test_file=$(create_test_file "$device_package" "$target")
    
    # Ensure cleanup
    trap "rm -f '$test_file'" EXIT
    
    # Run llgo build with real-time output
    if llgo build -tags nogc -target "$target" "$test_file" 2>&1; then
        printf " ${GREEN}✅ SUCCESS${NC}\n"
        rm -f "$test_file"
        return 0
    else
        printf " ${RED}❌ FAILED${NC}\n"
        echo -e "${RED}Error details for $device_package -> $target:${NC}"
        llgo build -tags nogc -target "$target" "$test_file" 2>&1 | sed 's/^/  /'
        echo ""
        rm -f "$test_file"
        return 1
    fi
}

# Main validation function
validate_all_packages() {
    local success_count=0
    local failure_count=0
    local failures=()
    local target_index=0
    local total_targets=0
    
    # Count total targets
    for entry in "${VALIDATION_TARGETS[@]}"; do
        local targets_str="${entry#*:}"
        if [[ -n "$targets_str" ]]; then
            IFS=',' read -ra targets <<< "$targets_str"
            total_targets=$((total_targets + ${#targets[@]}))
        fi
    done
    
    echo -e "${BLUE}🚀 Starting validation of ${#VALIDATION_TARGETS[@]} device packages ($total_targets total targets)${NC}"
    echo -e "${BLUE}Using LLGO command: llgo${NC}"
    echo ""
    
    # Process each package
    for entry in "${VALIDATION_TARGETS[@]}"; do
        local device_package="${entry%%:*}"
        local targets_str="${entry#*:}"
        
        if [[ -z "$targets_str" ]]; then
            continue
        fi
        
        IFS=',' read -ra targets <<< "$targets_str"
        echo -e "${BLUE}📦 Package: $device_package (${#targets[@]} targets)${NC}"
        
        for target in "${targets[@]}"; do
            target_index=$((target_index + 1))
            
            if validate_target "$device_package" "$target" "$target_index" "$total_targets"; then
                success_count=$((success_count + 1))
            else
                failure_count=$((failure_count + 1))
                failures+=("$device_package -> $target")
            fi
        done
        echo ""
    done
    
    # Summary
    echo -e "${BLUE}📊 VALIDATION SUMMARY${NC}"
    echo "==================="
    echo -e "${GREEN}✅ Successful: $success_count${NC}"
    echo -e "${RED}❌ Failed: $failure_count${NC}"
    
    if [[ $failure_count -gt 0 ]]; then
        echo ""
        echo -e "${RED}💥 FAILURES:${NC}"
        for failure in "${failures[@]}"; do
            echo "  - $failure"
        done
        return 1
    fi
    
    echo ""
    echo -e "${GREEN}🎉 All device package validations passed!${NC}"
    return 0
}

# Validate single target from command line
validate_single() {
    local input="$1"
    
    # Parse format: device/package:target
    if [[ ! "$input" =~ ^([^:]+):([^:]+)$ ]]; then
        echo -e "${RED}❌ Invalid format. Expected: device/package:target${NC}" >&2
        echo -e "${YELLOW}💡 Example: device/arm:cortex-m-qemu${NC}" >&2
        return 1
    fi
    
    local device_package="${BASH_REMATCH[1]}"
    local target="${BASH_REMATCH[2]}"
    
    echo -e "${BLUE}🎯 Single target validation${NC}"
    echo -e "${BLUE}Package: $device_package${NC}"
    echo -e "${BLUE}Target: $target${NC}"
    echo ""
    
    if validate_target "$device_package" "$target" "1" "1"; then
        echo ""
        echo -e "${GREEN}🎉 Single validation passed!${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}❌ Single validation failed${NC}" >&2
        return 1
    fi
}

# Main script
main() {
    # Check if llgo is available
    if ! command -v llgo &>/dev/null; then
        echo -e "${RED}❌ LLGO command not found in PATH${NC}" >&2
        echo -e "${YELLOW}💡 Make sure LLGO is installed and available as 'llgo'${NC}" >&2
        exit 1
    fi
    
    # Clean up any existing test files
    rm -f test_*.go
    
    # Check for single target mode
    if [[ $# -eq 1 ]]; then
        if validate_single "$1"; then
            exit 0
        else
            exit 1
        fi
    elif [[ $# -gt 1 ]]; then
        echo -e "${RED}❌ Too many arguments${NC}" >&2
        echo -e "${YELLOW}💡 Usage: $0 [device/package:target]${NC}" >&2
        exit 1
    fi
    
    # Run full validation
    if validate_all_packages; then
        exit 0
    else
        echo -e "${RED}❌ Validation failed${NC}" >&2
        exit 1
    fi
}

# Cleanup on script exit
cleanup() {
    rm -f test_*.go
}
trap cleanup EXIT

# Run main function with all arguments
main "$@"
