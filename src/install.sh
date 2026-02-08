#!/bin/bash
# Greenlight Installer v1.0.0
# TDD-first development system for Claude Code
#
# Usage:
#   bash install.sh --global     Install to ~/.claude/ (all projects)
#   bash install.sh --local      Install to ./.claude/ (this project only)
#   bash install.sh --uninstall  Remove from specified location
#   bash install.sh --check      Verify installation

set -euo pipefail

VERSION="1.0.0"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

print_header() {
    echo ""
    echo -e "${GREEN}Greenlight v${VERSION}${NC}"
    echo "   TDD-first development for Claude Code"
    echo ""
}

print_usage() {
    echo "Usage: bash install.sh [--global | --local] [--uninstall | --check]"
    echo ""
    echo "  --global, -g     Install to ~/.claude/ (all projects)"
    echo "  --local, -l      Install to ./.claude/ (current project only)"
    echo "  --uninstall      Remove Greenlight from install location"
    echo "  --check          Verify existing installation"
    echo "  --help, -h       Show this help"
    echo ""
}

# Validate source files exist before doing anything
validate_source() {
    local missing=0

    if [ ! -d "$SCRIPT_DIR/commands/gl" ]; then
        echo -e "${RED}Error: commands/gl/ directory not found in $SCRIPT_DIR${NC}"
        missing=1
    fi

    if [ ! -d "$SCRIPT_DIR/agents" ]; then
        echo -e "${RED}Error: agents/ directory not found in $SCRIPT_DIR${NC}"
        missing=1
    fi

    if [ ! -f "$SCRIPT_DIR/CLAUDE.md" ]; then
        echo -e "${RED}Error: CLAUDE.md not found in $SCRIPT_DIR${NC}"
        missing=1
    fi

    if [ "$missing" -eq 1 ]; then
        echo ""
        echo "Make sure you're running install.sh from the Greenlight source directory."
        exit 1
    fi
}

# Count files to be installed
count_files() {
    local dir="$1"
    local pattern="$2"
    find "$dir" -name "$pattern" -type f 2>/dev/null | wc -l | tr -d ' '
}

# Install Greenlight
do_install() {
    local install_dir="$1"
    local scope="$2"

    validate_source

    # Check if already installed — offer upgrade path
    if [ -f "$install_dir/commands/gl/init.md" ]; then
        local existing_version="unknown"
        if [ -f "$install_dir/.greenlight-version" ]; then
            existing_version=$(cat "$install_dir/.greenlight-version")
        fi
        echo -e "${YELLOW}Greenlight already installed (version: $existing_version)${NC}"
        echo ""
        read -p "Overwrite with v${VERSION}? [y/N]: " confirm
        case "$confirm" in
            [yY]|[yY][eE][sS]) ;;
            *) echo "Cancelled."; exit 0 ;;
        esac
        echo ""
    fi

    # Create directories
    mkdir -p "$install_dir/commands/gl"
    mkdir -p "$install_dir/agents"

    # Copy commands
    local cmd_count
    cmd_count=$(count_files "$SCRIPT_DIR/commands/gl" "*.md")
    cp "$SCRIPT_DIR/commands/gl/"*.md "$install_dir/commands/gl/"
    echo -e "  ${GREEN}+${NC} Commands: $cmd_count installed"

    # Copy agents
    local agent_count
    agent_count=$(count_files "$SCRIPT_DIR/agents" "*.md")
    cp "$SCRIPT_DIR/agents/"*.md "$install_dir/agents/"
    echo -e "  ${GREEN}+${NC} Agents: $agent_count installed"

    # Copy references (if they exist)
    if [ -d "$SCRIPT_DIR/references" ]; then
        mkdir -p "$install_dir/references"
        cp "$SCRIPT_DIR/references/"*.md "$install_dir/references/"
        local ref_count
        ref_count=$(count_files "$SCRIPT_DIR/references" "*.md")
        echo -e "  ${GREEN}+${NC} References: $ref_count installed"
    fi

    # Copy templates (if they exist)
    if [ -d "$SCRIPT_DIR/templates" ]; then
        mkdir -p "$install_dir/templates"
        cp "$SCRIPT_DIR/templates/"*.md "$install_dir/templates/"
        local tmpl_count
        tmpl_count=$(count_files "$SCRIPT_DIR/templates" "*.md")
        echo -e "  ${GREEN}+${NC} Templates: $tmpl_count installed"
    fi

    # Handle CLAUDE.md — never overwrite without consent
    local claude_target
    if [ "$scope" = "global" ]; then
        claude_target="$HOME/.claude/CLAUDE.md"
    else
        claude_target="./CLAUDE.md"
    fi

    if [ -f "$claude_target" ]; then
        echo ""
        echo -e "  ${YELLOW}!${NC} Existing CLAUDE.md found at $claude_target"
        echo "    Options:"
        echo "    1) Keep existing (Greenlight standards saved to CLAUDE_GREENLIGHT.md)"
        echo "    2) Replace with Greenlight CLAUDE.md (backup created)"
        echo "    3) Append Greenlight standards to existing"
        read -p "  Choice [1/2/3]: " claude_choice
        case "$claude_choice" in
            2)
                cp "$claude_target" "${claude_target}.backup.$(date +%Y%m%d%H%M%S)"
                cp "$SCRIPT_DIR/CLAUDE.md" "$claude_target"
                echo -e "  ${GREEN}+${NC} CLAUDE.md replaced (backup created)"
                ;;
            3)
                echo "" >> "$claude_target"
                echo "# --- Greenlight Engineering Standards ---" >> "$claude_target"
                echo "" >> "$claude_target"
                cat "$SCRIPT_DIR/CLAUDE.md" >> "$claude_target"
                echo -e "  ${GREEN}+${NC} Greenlight standards appended to CLAUDE.md"
                ;;
            *)
                cp "$SCRIPT_DIR/CLAUDE.md" "${claude_target%.md}_GREENLIGHT.md"
                echo -e "  ${GREEN}+${NC} Standards saved to CLAUDE_GREENLIGHT.md"
                ;;
        esac
    else
        cp "$SCRIPT_DIR/CLAUDE.md" "$claude_target"
        echo -e "  ${GREEN}+${NC} CLAUDE.md installed"
    fi

    # Write version file for future upgrade detection
    echo "$VERSION" > "$install_dir/.greenlight-version"

    echo ""
    echo -e "${GREEN}Greenlight v${VERSION} installed ($scope)${NC}"
    echo ""
    echo -e "  Commands:   ${CYAN}/gl:help${NC}     — see all commands"
    echo -e "  Start:      ${CYAN}/gl:init${NC}     — new project"
    echo -e "  Brownfield: ${CYAN}/gl:map${NC}      — analyse existing code first"
    echo ""
}

# Uninstall Greenlight
do_uninstall() {
    local install_dir="$1"
    local scope="$2"

    if [ ! -d "$install_dir/commands/gl" ] && [ ! -f "$install_dir/agents/gl-architect.md" ]; then
        echo -e "${YELLOW}Greenlight not found in $install_dir${NC}"
        exit 0
    fi

    echo "This will remove:"
    [ -d "$install_dir/commands/gl" ] && echo "  - $install_dir/commands/gl/"
    ls "$install_dir/agents/gl-"*.md 2>/dev/null && echo "  - Greenlight agents from $install_dir/agents/"
    [ -d "$install_dir/references" ] && echo "  - $install_dir/references/"
    [ -d "$install_dir/templates" ] && echo "  - $install_dir/templates/"
    [ -f "$install_dir/.greenlight-version" ] && echo "  - $install_dir/.greenlight-version"
    echo ""
    echo "CLAUDE.md will NOT be removed (may contain your customizations)."
    echo ""

    read -p "Continue? [y/N]: " confirm
    case "$confirm" in
        [yY]|[yY][eE][sS]) ;;
        *) echo "Cancelled."; exit 0 ;;
    esac

    rm -rf "$install_dir/commands/gl"
    rm -f "$install_dir/agents/gl-"*.md
    rm -rf "$install_dir/references"
    rm -rf "$install_dir/templates"
    rm -f "$install_dir/.greenlight-version"

    echo -e "${GREEN}Greenlight removed from $install_dir${NC}"
}

# Check installation
do_check() {
    local install_dir="$1"
    local scope="$2"
    local issues=0

    echo "Checking $scope installation at $install_dir..."
    echo ""

    # Version
    if [ -f "$install_dir/.greenlight-version" ]; then
        local ver
        ver=$(cat "$install_dir/.greenlight-version")
        echo -e "  Version: ${CYAN}$ver${NC}"
        if [ "$ver" != "$VERSION" ]; then
            echo -e "  ${YELLOW}Update available: v${VERSION}${NC}"
        fi
    else
        echo -e "  ${YELLOW}No version file — may be pre-v1 install${NC}"
        issues=$((issues + 1))
    fi

    # Commands
    local cmd_count
    cmd_count=$(count_files "$install_dir/commands/gl" "*.md")
    if [ "$cmd_count" -gt 0 ]; then
        echo -e "  Commands: ${GREEN}$cmd_count found${NC}"
    else
        echo -e "  Commands: ${RED}none found${NC}"
        issues=$((issues + 1))
    fi

    # Agents
    local agent_count=0
    if [ -d "$install_dir/agents" ]; then
        agent_count=$(ls "$install_dir/agents/gl-"*.md 2>/dev/null | wc -l | tr -d ' ')
    fi
    if [ "$agent_count" -gt 0 ]; then
        echo -e "  Agents: ${GREEN}$agent_count found${NC}"
    else
        echo -e "  Agents: ${RED}none found${NC}"
        issues=$((issues + 1))
    fi

    # References
    local ref_count
    ref_count=$(count_files "$install_dir/references" "*.md")
    if [ "$ref_count" -gt 0 ]; then
        echo -e "  References: ${GREEN}$ref_count found${NC}"
    else
        echo -e "  References: ${YELLOW}none (optional)${NC}"
    fi

    # CLAUDE.md
    local claude_target
    if [ "$scope" = "global" ]; then
        claude_target="$HOME/.claude/CLAUDE.md"
    else
        claude_target="./CLAUDE.md"
    fi
    if [ -f "$claude_target" ]; then
        echo -e "  CLAUDE.md: ${GREEN}found${NC}"
    else
        echo -e "  CLAUDE.md: ${RED}not found${NC}"
        issues=$((issues + 1))
    fi

    echo ""
    if [ "$issues" -eq 0 ]; then
        echo -e "${GREEN}Installation looks good.${NC}"
    else
        echo -e "${YELLOW}$issues issue(s) found. Consider reinstalling with: bash install.sh --${scope}${NC}"
    fi
}

# Parse arguments
INSTALL_DIR=""
SCOPE=""
ACTION="install"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --global|-g)
            INSTALL_DIR="$HOME/.claude"
            SCOPE="global"
            shift
            ;;
        --local|-l)
            INSTALL_DIR="./.claude"
            SCOPE="local"
            shift
            ;;
        --uninstall)
            ACTION="uninstall"
            shift
            ;;
        --check)
            ACTION="check"
            shift
            ;;
        --help|-h)
            print_header
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# If no location specified, ask
if [ -z "$INSTALL_DIR" ]; then
    print_header
    echo "Where do you want to install?"
    echo "  1) Global  (~/.claude/ — all projects)"
    echo "  2) Local   (./.claude/ — this project only)"
    read -p "Choice [1/2]: " choice
    case "$choice" in
        1) INSTALL_DIR="$HOME/.claude"; SCOPE="global" ;;
        2) INSTALL_DIR="./.claude"; SCOPE="local" ;;
        *) echo -e "${RED}Invalid choice${NC}"; exit 1 ;;
    esac
    echo ""
fi

print_header

case "$ACTION" in
    install)   do_install "$INSTALL_DIR" "$SCOPE" ;;
    uninstall) do_uninstall "$INSTALL_DIR" "$SCOPE" ;;
    check)     do_check "$INSTALL_DIR" "$SCOPE" ;;
esac
