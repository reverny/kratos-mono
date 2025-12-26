#!/bin/bash

# Update Swagger documentation from proto files
# This script:
# 1. Parses proto files directly
# 2. Generates OpenAPI specs using Python
# 3. Updates swagger.html files for each service

set -e

ROOT_DIR=$(pwd)

# Function to get proto file for a service
get_proto_file() {
    case "$1" in
        inventory)
            echo "api/inventory/v1/inventory.proto"
            ;;
        user)
            echo "api/user/v1/user.proto"
            ;;
        product)
            echo "api/product/v1/product.proto"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Function to update swagger HTML for a service
update_service_swagger() {
    local service_name=$1
    local service_dir="services/${service_name}"
    local swagger_html="${service_dir}/internal/server/swagger.html"
    local proto_file=$(get_proto_file "${service_name}")
    
    if [ -z "${proto_file}" ]; then
        return
    fi
    
    if [ ! -f "${swagger_html}" ]; then
        echo "⚠️  Warning: ${swagger_html} not found, skipping ${service_name}"
        return
    fi
    
    if [ ! -f "${proto_file}" ]; then
        echo "⚠️  Warning: ${proto_file} not found, skipping ${service_name}"
        return
    fi
    
    echo "📝 Updating swagger for ${service_name} from ${proto_file}..."
    
    # Generate OpenAPI spec from proto file and update HTML in one Python script
    temp_file=$(python3 - <<PYEOF
import sys
import re
sys.path.insert(0, 'scripts')

# Import the proto-to-swagger functionality
import subprocess
result = subprocess.run(
    ['python3', 'scripts/proto-to-swagger.py', '${proto_file}', '${service_name}'],
    capture_output=True,
    text=True
)

if result.returncode != 0:
    print(f"Error generating spec: {result.stderr}", file=sys.stderr)
    sys.exit(1)

spec_json = result.stdout

# Read current swagger HTML
with open('${swagger_html}', 'r') as f:
    content = f.read()

# Find and replace the spec object - it's a JSON object after "spec: "
# We need to replace only the spec object, keeping everything after it
import json

# Find where spec starts
spec_marker = 'spec: '
spec_pos = content.find(spec_marker)
if spec_pos == -1:
    print("Error: Could not find 'spec: ' in HTML", file=sys.stderr)
    sys.exit(1)

# Start parsing after "spec: "
parse_start = spec_pos + len(spec_marker)

# Find the matching closing brace for the spec object
brace_count = 0
pos = parse_start
in_string = False
escape_next = False
found_opening = False

while pos < len(content):
    char = content[pos]
    
    if escape_next:
        escape_next = False
        pos += 1
        continue
    
    if char == '\\\\':
        escape_next = True
        pos += 1
        continue
    
    if char == '"':
        in_string = not in_string
    elif not in_string:
        if char == '{':
            if not found_opening:
                found_opening = True
            else:
                brace_count += 1
        elif char == '}':
            if brace_count == 0:
                # Found the closing brace of spec object
                spec_end = pos + 1
                break
            brace_count -= 1
    
    pos += 1

if pos >= len(content):
    print("Error: Could not find matching closing brace for spec", file=sys.stderr)
    sys.exit(1)

# Replace: keep everything before "spec: ", add new spec with comma, keep everything after the closing brace
# Check if there's already a comma after the closing brace
next_char_pos = spec_end
while next_char_pos < len(content) and content[next_char_pos] in ' \\n\\t\\r':
    next_char_pos += 1

# Skip the comma if it exists (we'll add it back)
if next_char_pos < len(content) and content[next_char_pos] == ',':
    after_spec = content[next_char_pos:]  # Keep comma and everything after
else:
    after_spec = ',' + content[spec_end:]  # Add comma

new_content = content[:spec_pos] + 'spec: ' + spec_json.strip() + after_spec

# Write to temp file first
import tempfile
with tempfile.NamedTemporaryFile(mode='w', delete=False, suffix='.html') as f:
    f.write(new_content)
    temp_path = f.name

print(temp_path)
PYEOF
)
    
    if [ -n "${temp_file}" ] && [ -f "${temp_file}" ]; then
        # Backup original
        cp "${swagger_html}" "${swagger_html}.bak"
        
        # Replace with updated version
        mv "${temp_file}" "${swagger_html}"
        
        # Also update the docs/ version
        cp "${swagger_html}" "${service_dir}/docs/swagger.html"
        
        echo "✅ Updated ${service_name} swagger"
    else
        echo "❌ Failed to update ${service_name} swagger"
    fi
}

# Find all services
echo "🔍 Scanning for services..."
services=$(find services -maxdepth 1 -type d ! -path services -exec basename {} \;)

for service in ${services}; do
    proto_file=$(get_proto_file "${service}")
    if [ -n "${proto_file}" ]; then
        update_service_swagger "${service}"
    fi
done

echo ""
echo "🎉 Swagger update complete!"
echo "💡 Remember to rebuild your services to embed the updated swagger:"
echo "   make build"
