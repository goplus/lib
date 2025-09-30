#!/bin/bash
set -e

# .ld,.s current not need generate at emb/device/avr,they are in llgo

declare_avr_config() {
    name="avr"
    repo="https://github.com/avr-rust/avr-mcu"
    lib_path="lib/avr"
    git_hash="6624554c02b237b23dc17d53e992bf54033fc228"
    tasks=(
        "lib/avr/packs/atmega"
        "lib/avr/packs/tiny"
    )
    ignore_list="avr.go avr_tiny85.go"
    target="emb/device/avr"
    generator="gen-device-avr"
}

declare_esp_config() {
    name="esp"
    repo="https://github.com/cmsis-svd/cmsis-svd-data"
    lib_path="lib/cmsis-svd"
    git_hash="05a9562ec59b87945a8d7177a4b08b7aa2f2fd58"
    tasks=(
        "-source=https://github.com/posborne/cmsis-svd/tree/master/data/Espressif-Community -interrupts=software lib/cmsis-svd/data/Espressif-Community/"
        "-source=https://github.com/posborne/cmsis-svd/tree/master/data/Espressif -interrupts=software lib/cmsis-svd/data/Espressif/"
    )
    ignore_list=""
    target="emb/device/esp"
    generator="gen-device-svd"
}

declare_nrf_config() {
    name="nrf"
    repo="https://github.com/NordicSemiconductor/nrfx"
    lib_path="lib/nrfx"
    git_hash="d779b49fc59c7a165e7da1d7cd7d57b28a059f16"
    tasks=(
        "-source=https://github.com/NordicSemiconductor/nrfx/tree/master/mdk lib/nrfx/mdk/"
    )
    ignore_list="README.markdown"
    target="emb/device/nrf"
    generator="gen-device-svd"
}

declare_nxp_config() {
    name="nxp"
    repo="https://github.com/cmsis-svd/cmsis-svd-data"
    lib_path="lib/cmsis-svd"
    git_hash="05a9562ec59b87945a8d7177a4b08b7aa2f2fd58"
    tasks=(
        "-source=https://github.com/posborne/cmsis-svd/tree/master/data/NXP lib/cmsis-svd/data/NXP/"
    )
    ignore_list="hardfault.go mimxrt1062_clock.go mimxrt1062_hardfault.go mimxrt1062_mpu.go"
    target="emb/device/nxp"
    generator="gen-device-svd"
}

declare_sam_config() {
    name="sam"
    repo="https://github.com/cmsis-svd/cmsis-svd-data"
    lib_path="lib/cmsis-svd"
    git_hash="05a9562ec59b87945a8d7177a4b08b7aa2f2fd58"
    tasks=(
        "-source=https://github.com/posborne/cmsis-svd/tree/master/data/Atmel lib/cmsis-svd/data/Atmel/"
    )
    ignore_list="atsamd51x-bitfields.go atsame5x-bitfields.go"
    target="emb/device/sam"
    generator="gen-device-svd"
}

declare_sifive_config() {
    name="sifive"
    repo="https://github.com/cmsis-svd/cmsis-svd-data"
    lib_path="lib/cmsis-svd"
    git_hash="05a9562ec59b87945a8d7177a4b08b7aa2f2fd58"
    tasks=(
        "-source=https://github.com/posborne/cmsis-svd/tree/master/data/SiFive-Community -interrupts=software lib/cmsis-svd/data/SiFive-Community/"
    )
    ignore_list=""
    target="emb/device/sifive"
    generator="gen-device-svd"
}

# List of device configuration functions
DEVICE_CONFIGS=(
    "declare_avr_config"
    "declare_esp_config"
    "declare_nrf_config"
    "declare_nxp_config"
    "declare_sam_config"
    "declare_sifive_config"
)

# Generate device files from all tasks to specified directory
generate_device_files() {
    local device_name="$1"
    shift
    local generator="${@: -2:1}"
    local target_dir="${@: -1}"
    # Remove generator and target_dir from args to get tasks
    set -- "${@:1:$#-2}"
    local tasks_array=("$@")

    echo "[$device_name] Processing ${#tasks_array[@]} tasks..."
    for task in "${tasks_array[@]}"; do
        echo "[$device_name] Processing task: $task"
        "$generator" $task "$target_dir/"
        echo "[$device_name] Generated to $target_dir"
    done
    echo "[$device_name] All tasks processed successfully"

    # Format Go files in target directory
    echo "[$device_name] Formatting Go files in $target_dir..."
    (cd "$target_dir" && GO111MODULE=off ${GO:-go} fmt .)
    echo "[$device_name] Go formatting completed"
}

# Generate devices for a specific device configuration
generate_devices() {
    local config_func="$1"
    local mode="$2"

    # Call the configuration function to set variables
    $config_func

    # Validate configuration
    if [[ -z "$repo" || -z "$lib_path" || -z "$git_hash" || -z "$generator" || -z "$target" || ${#tasks[@]} -eq 0 ]]; then
        echo "Error: Missing configuration for device '$name'"
        return 1
    fi

    # Use target directory from configuration
    local target_dir="$target"
    case "$mode" in
        generate)
            echo "[$name] Running in generation mode - will clean and generate to $target_dir"
            ;;
        verify)
            echo "[$name] Running in verification mode - will verify generated content matches $target_dir"
            ;;
        *)
            echo "[$name] Running in default mode - will generate to $target_dir"
            ;;
    esac

    # Clone and setup repository
    echo "[$name] Setting up repository..."
    rm -rf "$lib_path"
    git clone "$repo" "$lib_path"
    cd "$lib_path"
    git checkout "$git_hash"
    cd ../../

    # Create target directory
    mkdir -p "$target_dir"

    # Handle verification mode
    if [ "$mode" = "verify" ]; then
        echo "[$name] Starting verification process..."

        # Check if target directory exists and has content
        if [ ! -d "$target_dir" ] || [ -z "$(ls -A "$target_dir" 2>/dev/null)" ]; then
            echo "[$name] Verification failed: target directory is empty or does not exist"
            return 1
        fi

        # Create temporary directory with same structure
        temp_target=".verify/$target"
        echo "[$name] Creating temporary directory: $temp_target"
        mkdir -p "$temp_target"

        # Generate files to temporary directory
        echo "[$name] Generating files to temporary directory..."
        generate_device_files "$name" "${tasks[@]}" "$generator" "$temp_target"

        # Copy ignore list files to temporary directory
        if [ -n "$ignore_list" ]; then
            echo "[$name] Copying ignore list files to temporary directory: $ignore_list"
            for file in $ignore_list; do
                if [ -f "$target_dir/$file" ]; then
                    cp "$target_dir/$file" "$temp_target/"
                    echo "[$name] Copied: $file"
                fi
            done
        fi

        # Compare directories
        echo "[$name] Comparing original directory with generated directory..."
        if diff -u -r "$target_dir" "$temp_target" > /dev/null 2>&1; then
            echo "[$name] Verification passed: generated files match existing files"
            VERIFY_RESULT=0
        else
            echo "[$name] Verification failed: differences found"
            echo "[$name] Detailed differences:"
            diff -u -r "$target_dir" "$temp_target" || true
            VERIFY_RESULT=1
        fi

        return $VERIFY_RESULT
    fi

    # Clean target directory if in generation mode
    if [ "$mode" = "generate" ] && [ -n "$ignore_list" ]; then
        echo "[$name] Cleaning $target_dir (preserving ignore list: $ignore_list)"

        # Create temp directory to store ignored files
        TEMP_DIR=$(mktemp -d)

        # Backup ignored files if they exist
        for file in $ignore_list; do
            if [ -f "$target_dir/$file" ]; then
                cp "$target_dir/$file" "$TEMP_DIR/"
                echo "[$name] Backed up: $file"
            fi
        done

        # Remove all files from target directory
        rm -rf "$target_dir"/*

        # Restore ignored files
        for file in $ignore_list; do
            if [ -f "$TEMP_DIR/$file" ]; then
                cp "$TEMP_DIR/$file" "$target_dir/"
                echo "[$name] Restored: $file"
            fi
        done

        # Clean up temp directory
        rm -rf "$TEMP_DIR"
    fi

    # Generate device files from all tasks
    generate_device_files "$name" "${tasks[@]}" "$generator" "$target_dir"
}

# Parse command line arguments
MODE=""
case "$1" in
    generate)
        MODE="generate"
        ;;
    verify)
        MODE="verify"
        ;;
    *)
        MODE="default"
        ;;
esac

# Initialize verification tracking
VERIFICATION_FAILED=0

# Process all enabled devices
for config_func in "${DEVICE_CONFIGS[@]}"; do
    if ! generate_devices "$config_func" "$MODE"; then
        if [ "$MODE" = "verify" ]; then
            VERIFICATION_FAILED=1
        fi
    fi
done

if [ "$MODE" = "verify" ]; then
    if [ $VERIFICATION_FAILED -eq 0 ]; then
        echo "All device verifications passed"
    else
        echo "Device verification failed: inconsistencies found"
        exit 1
    fi
else
    echo "All device generations completed!"
fi