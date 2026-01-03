#!/bin/bash

# Build Site Script
# Generates a site from a template and site configuration

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATES_DIR="$SCRIPT_DIR/templates"
SITES_DIR="$SCRIPT_DIR/sites"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to build a site
build_site() {
    local site_name=$1
    local site_dir="$SITES_DIR/$site_name"
    local config_file="$site_dir/config.json"

    if [ ! -f "$config_file" ]; then
        log_error "Config file not found: $config_file"
        return 1
    fi

    log_info "Building site: $site_name"

    # Read config
    local template=$(jq -r '.template' "$config_file")
    local enabled=$(jq -r '.enabled' "$config_file")

    if [ "$enabled" != "true" ]; then
        log_warn "Site $site_name is disabled. Skipping."
        return 0
    fi

    local template_dir="$TEMPLATES_DIR/$template"
    if [ ! -d "$template_dir" ]; then
        log_error "Template not found: $template"
        return 1
    fi

    log_info "Using template: $template"

    # Create public directory
    local public_dir="$site_dir/public"
    mkdir -p "$public_dir"

    # Copy template files
    log_info "Copying template files..."
    cp -r "$template_dir"/*.html "$public_dir/" 2>/dev/null || true
    cp -r "$template_dir"/*.js "$public_dir/" 2>/dev/null || true
    cp -r "$template_dir"/*.css "$public_dir/" 2>/dev/null || true

    # Copy branding.json
    if [ -f "$site_dir/branding.json" ]; then
        log_info "Copying branding configuration..."
        cp "$site_dir/branding.json" "$public_dir/"
    else
        log_warn "No branding.json found for $site_name"
    fi

    # Copy config.json (for client-side access if needed)
    cp "$config_file" "$public_dir/"

    # Copy assets
    if [ -d "$site_dir/assets" ]; then
        log_info "Copying site assets..."
        cp -r "$site_dir/assets"/* "$public_dir/" 2>/dev/null || true
    fi

    # Set proper permissions
    chmod -R 644 "$public_dir"/* 2>/dev/null || true
    find "$public_dir" -type d -exec chmod 755 {} \; 2>/dev/null || true

    log_info "✅ Site built successfully: $site_name"
    log_info "   Output: $public_dir"
    echo ""
}

# Main execution
if [ $# -eq 0 ]; then
    # Build all sites
    log_info "Building all sites..."
    echo ""

    for site_dir in "$SITES_DIR"/*; do
        if [ -d "$site_dir" ]; then
            site_name=$(basename "$site_dir")
            build_site "$site_name"
        fi
    done
else
    # Build specific site
    build_site "$1"
fi

log_info "🎉 Build complete!"
