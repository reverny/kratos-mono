#!/bin/bash

# Update Swagger documentation from generated HTTP pb.go files
# This is more reliable than parsing proto files directly

set -e

ROOT_DIR=$(pwd)

echo "🔄 Updating Swagger documentation..."
echo ""

# First, ensure we have the latest generated code
echo "1️⃣  Generating latest proto code..."
buf generate

echo ""
echo "2️⃣  Extracting OpenAPI specs from generated code..."

# Function to update swagger for a service
update_service_swagger() {
    local service_name=$1
    local service_dir="services/${service_name}"
    local swagger_html="${service_dir}/internal/server/swagger.html"
    local http_pb_go="gen/go/api/${service_name}/v1/${service_name}_http.pb.go"
    
    if [ ! -f "${http_pb_go}" ]; then
        echo "⚠️  ${service_name}: HTTP pb.go not found, skipping"
        return
    fi
    
    if [ ! -f "${swagger_html}" ]; then
        echo "⚠️  ${service_name}: swagger.html not found, skipping"
        return
    fi
    
    echo "📝 Updating ${service_name}..."
    
    # Generate OpenAPI spec from HTTP pb.go using Go program
    go run scripts/gen-swagger-from-pb.go "${http_pb_go}" "${service_name}" > /tmp/openapi-${service_name}.json
    
    if [ $? -ne 0 ]; then
        echo "❌ Failed to generate spec for ${service_name}"
        return
    fi
    
    # Update swagger.html with new spec
    python3 - <<PYEOF
import json
import re

# Read generated spec
with open('/tmp/openapi-${service_name}.json', 'r') as f:
    spec_json = f.read()

# Read current swagger HTML
with open('${swagger_html}', 'r') as f:
    content = f.read()

# Find and replace the spec object
spec_marker = 'spec: '
spec_pos = content.find(spec_marker)
if spec_pos == -1:
    print("Error: Could not find 'spec: ' in HTML")
    exit(1)

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
                spec_end = pos + 1
                break
            brace_count -= 1
    
    pos += 1

# Replace spec
new_content = content[:spec_pos] + 'spec: ' + spec_json.strip() + ',' + content[spec_end:]

# Write updated HTML
with open('${swagger_html}', 'w') as f:
    f.write(new_content)

# Also update docs/ version
import shutil
shutil.copy('${swagger_html}', '${service_dir}/docs/swagger.html')

print("✅ Updated ${service_name}")
PYEOF
}

# Update all services
for service in inventory user product; do
    update_service_swagger "${service}"
done

echo ""
echo "🎉 Swagger update complete!"
