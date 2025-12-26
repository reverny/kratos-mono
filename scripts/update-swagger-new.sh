#!/bin/bash

# Update Swagger documentation from generated OpenAPI specs
# Uses grpc-gateway generated .swagger.json files

set -e

ROOT_DIR=$(pwd)

echo "🔄 Updating Swagger documentation..."
echo ""

# Step 1: Generate all code including OpenAPI specs
echo "1️⃣  Running buf generate..."
buf generate

if [ $? -ne 0 ]; then
    echo "❌ buf generate failed"
    exit 1
fi

echo ""
echo "2️⃣  Updating swagger HTML files..."

# Function to update swagger HTML for a service
update_service_swagger() {
    local service_name=$1
    local swagger_json="gen/openapi/${service_name}/v1/${service_name}.swagger.json"
    local swagger_internal="services/${service_name}/internal/server/swagger.html"
    local swagger_docs="services/${service_name}/docs/swagger.html"
    
    # Check if generated swagger.json exists
    if [ ! -f "${swagger_json}" ]; then
        echo "⚠️  ${service_name}: swagger.json not found at ${swagger_json}, skipping"
        return
    fi
    
    # Check if swagger HTML files exist
    if [ ! -f "${swagger_internal}" ]; then
        echo "⚠️  ${service_name}: swagger.html not found, skipping"
        return
    fi
    
    echo "📝 Updating ${service_name}..."
    
    # Update both swagger HTML files
    for html_file in "${swagger_internal}" "${swagger_docs}"; do
        if [ ! -f "${html_file}" ]; then
            continue
        fi
        
        # Use Python to replace the spec in HTML
        python3 - "${html_file}" "${swagger_json}" <<'PYEOF'
import sys
import json
import re

html_file = sys.argv[1]
spec_file = sys.argv[2]

# Read the generated spec
with open(spec_file, 'r') as f:
    new_spec = json.load(f)

# Read current HTML
with open(html_file, 'r') as f:
    content = f.read()

# Find 'spec: ' marker
spec_marker = 'spec: '
spec_pos = content.find(spec_marker)
if spec_pos == -1:
    print(f"Error: Could not find '{spec_marker}' in {html_file}", file=sys.stderr)
    sys.exit(1)

# Find the closing brace of the spec object
parse_start = spec_pos + len(spec_marker)
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
    
    if char == '\\':
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
                spec_end = pos + 1
                break
            brace_count -= 1
    
    pos += 1

if not found_opening or pos >= len(content):
    print(f"Error: Could not find matching closing brace for spec in {html_file}", file=sys.stderr)
    sys.exit(1)

# Check if there's a comma after the closing brace
next_pos = spec_end
while next_pos < len(content) and content[next_pos] in ' \n\t\r':
    next_pos += 1

has_comma = next_pos < len(content) and content[next_pos] == ','

# Build new content
new_spec_json = json.dumps(new_spec, indent=2)
new_content = content[:spec_pos] + 'spec: ' + new_spec_json

if has_comma:
    new_content += content[spec_end:]
else:
    new_content += ',' + content[spec_end:]

# Write updated HTML
with open(html_file, 'w') as f:
    f.write(new_content)

print(f"   ✅ Updated {html_file}")
PYEOF
        
        if [ $? -ne 0 ]; then
            echo "   ❌ Failed to update ${html_file}"
        fi
    done
}

# Update all services - auto-detect from gen/openapi directory
echo "🔍 Detecting services..."
SERVICES=$(find gen/openapi -maxdepth 1 -type d ! -path gen/openapi -exec basename {} \; | grep -v "^common$" | sort)

if [ -z "$SERVICES" ]; then
    echo "⚠️  No services found in gen/openapi/"
    echo "   Make sure you have run 'buf generate' first"
    exit 1
fi

echo "   Found: ${SERVICES}"
echo ""

for service in ${SERVICES}; do
    update_service_swagger "${service}"
done

echo ""
echo "🎉 Swagger update complete!"
echo "💡 Swagger HTML files have been updated. Rebuild services to embed changes:"
echo "   make build"
