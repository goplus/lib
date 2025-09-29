#!/bin/bash
set -e

# .ld,.s current not need generate at emb/device/avr,they are in llgo

# Device configurations (key-value style for better readability)
declare_avr_config() {
    name="avr"
    repo="https://github.com/avr-rust/avr-mcu"
    lib_path="lib/avr"
    git_hash="6624554c02b237b23dc17d53e992bf54033fc228"
    packs_path=("packs/atmega" "packs/tiny")
    ignore_list="avr.go avr_tiny85.go"
    target="emb/device/avr"
    generator="gen-device-avr"
}

# List of device configuration functions
DEVICE_CONFIGS=("declare_avr_config")

# Generate devices for a specific device configuration
generate_devices() {
    local config_func="$1"
    local mode="$2"

    # Call the configuration function to set variables
    $config_func

    # Validate configuration
    if [[ -z "$repo" || -z "$lib_path" || -z "$git_hash" || -z "$generator" || -z "$target" || ${#packs_path[@]} -eq 0 ]]; then
        echo "Error: Missing configuration for device '$name'"
        return 1
    fi

    # Use target directory from configuration
    local target_dir="$target"
    if [ "$mode" = "generate" ]; then
        echo "[$name] Running in generation mode - will clean and generate to $target_dir"
    else
        echo "[$name] Running in default mode - will generate to $target_dir"
    fi

    # Clone and setup repository
    echo "[$name] Setting up repository..."
    rm -rf "$lib_path"
    git clone "$repo" "$lib_path"
    cd "$lib_path"
    git checkout "$git_hash"
    cd ../../

    # Create target directory
    mkdir -p "$target_dir"

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

    # Generate device files from all packs
    echo "[$name] Generating devices from ${#packs_path[@]} packs..."
    for pack in "${packs_path[@]}"; do
        echo "[$name] Processing pack: $pack"
        "$generator" "$lib_path/$pack" "$target_dir/"
        echo "[$name] Generated devices from $lib_path/$pack to $target_dir"
    done
    echo "[$name] All packs processed successfully"
}

# Parse command line arguments
MODE=""
if [ "$1" = "generate" ]; then
    MODE="generate"
else
    MODE="default"
fi

# Process all enabled devices
for config_func in "${DEVICE_CONFIGS[@]}"; do
    generate_devices "$config_func" "$MODE"
done

echo "All device generations completed!"