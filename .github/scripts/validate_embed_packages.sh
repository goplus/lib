#!/bin/bash

set -e


# Usage function
usage() {
    echo "Usage: $0 [device|machine]"
    echo ""
    echo "Arguments:"
    echo "  device    Validate device packages"
    echo "  machine   Validate machine packages"
}

# Device validation targets
DEVICE_VALIDATION_TARGETS=(
    # arm
    "device/arm:cortex-m-qemu"

    # arm64
    "device/arm64:esp32"

    # avr
    "device/avr:arduino,arduino-leonardo,atmega1284p,atmega328pb,attiny1616,simavr"

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

# Machine validation targets
MACHINE_VALIDATION_TARGETS=(
)

# Global variables
VALIDATION_TYPE=""
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
PASSED_LIST=()
FAILED_LIST=()

# Logging functions
log_info() {
    echo "[INFO] $1"
}

log_success() {
    echo "[SUCCESS] $1"
}

log_error() {
    echo "[ERROR] $1"
}

# Create test file for device package validation
create_device_test_file() {
    local device_package="$1"
    local target="$2"
    local filename="test_${target//[-.]/_}.go"

    cat > "$filename" << EOF
package main
import _ "github.com/goplus/lib/emb/$device_package"
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

# Parse command line arguments
parse_arguments() {
    if [ $# -eq 0 ]; then
        log_error "No arguments provided"
        usage
        exit 1
    fi

    VALIDATION_TYPE="$1"

    # Validate the validation type
    if [[ "$VALIDATION_TYPE" != "device" && "$VALIDATION_TYPE" != "machine" ]]; then
        log_error "Invalid validation type: $VALIDATION_TYPE"
        log_error "Must be either 'device' or 'machine'"
        usage
        exit 1
    fi
}

# Create test file for machine validation
create_machine_test_file() {
    local test_dir="_demo/emptycheck"
    local test_file="$test_dir/main.go"

    mkdir -p "$test_dir"
    cat > "$test_file" << 'EOF'
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
    log_info "Created test file: $test_file" >&2
    echo "$test_file"
}

# Device package validation function
validate_device_package() {
    local package_spec="$1"
    local package_path="${package_spec%:*}"
    local targets="${package_spec#*:}"

    log_info "Validating device package: $package_path"

    # Parse targets
    IFS=',' read -ra target_array <<< "$targets"

    for target in "${target_array[@]}"; do
        local test_number=$((PASSED_TESTS + FAILED_TESTS + 1))

        echo -n "[$test_number] Testing $package_path with target $target... "

        # Create test file
        local test_file
        test_file=$(create_device_test_file "$package_path" "$target")

        # Run llgo build test
        if llgo build -v -tags nogc -target "$target" "$test_file" 2>&1; then
            echo "PASS"
            rm -f "$test_file"
            ((PASSED_TESTS++))
            PASSED_LIST+=("$package_path -> $target")
        else
            echo "FAIL"
            rm -f "$test_file"
            ((FAILED_TESTS++))
            FAILED_LIST+=("$package_path -> $target")
            return 1
        fi
    done

    return 0
}

# Machine package validation function
validate_machine_target() {
    local target="$1"
    local test_file="$2"
    local test_number=$((PASSED_TESTS + FAILED_TESTS + 1))

    echo -n "[$test_number] Testing $target... "

    # Run llgo build test
    if llgo build -v -tags nogc -target="$target" "$test_file" 2>&1; then
        echo "PASS"
        ((PASSED_TESTS++))
        PASSED_LIST+=("$target")
        return 0
    else
        echo "FAIL"
        ((FAILED_TESTS++))
        FAILED_LIST+=("$target")
        return 1
    fi
}

# Main validation function
run_validation() {
    log_info "Starting $VALIDATION_TYPE package validation"

    if [ "$VALIDATION_TYPE" = "device" ]; then
        log_info "Total device packages to validate: ${#DEVICE_VALIDATION_TARGETS[@]}"

        for package_spec in "${DEVICE_VALIDATION_TARGETS[@]}"; do
            validate_device_package "$package_spec" || true
        done

    elif [ "$VALIDATION_TYPE" = "machine" ]; then
        TOTAL_TESTS=${#MACHINE_VALIDATION_TARGETS[@]}
        log_info "Total machine targets to validate: $TOTAL_TESTS"

        # Create test file for machine validation
        local test_file=$(create_machine_test_file)

        for target in "${MACHINE_VALIDATION_TARGETS[@]}"; do
            validate_machine_target "$target" "$test_file" || true
        done
    fi
}

# Print final results
print_results() {
    local total_tests=$((PASSED_TESTS + FAILED_TESTS))

    echo ""
    echo "======================================"
    echo "         VALIDATION SUMMARY"
    echo "======================================"
    echo "Validation Type: $VALIDATION_TYPE"
    echo "Total Tests: $total_tests"
    echo "Passed: $PASSED_TESTS"
    echo "Failed: $FAILED_TESTS"
    echo "======================================"

    if [ ${#PASSED_LIST[@]} -gt 0 ]; then
        echo ""
        echo "PASSED TARGETS:"
        for target in "${PASSED_LIST[@]}"; do
            echo "  ✓ $target"
        done
    fi

    if [ ${#FAILED_LIST[@]} -gt 0 ]; then
        echo ""
        echo "FAILED TARGETS:"
        for target in "${FAILED_LIST[@]}"; do
            echo "  ✗ $target"
        done
    fi

    echo "======================================"
}


# Main execution
main() {
    parse_arguments "$@"

    # Clean up any existing test files
    rm -f test_*.go

    run_validation
    print_results

    # Clean up test files
    rm -f test_*.go

    if [ $FAILED_TESTS -eq 0 ]; then
        log_success "All validations passed!"
        exit 0
    else
        log_error "Some validations failed!"
        exit 1
    fi
}

# Run main function with all arguments
main "$@"
