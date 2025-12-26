#!/bin/bash

# Update Swagger documentation - Simple and reliable approach
# Uses buf-generated OpenAPI specs directly

set -e

ROOT_DIR=$(pwd)

echo "🔄 Updating Swagger documentation from proto files..."
echo ""

# Step 1: Generate all code including OpenAPI
echo "1️⃣  Running buf generate..."
buf generate

echo ""
echo "2️⃣  Updating swagger HTML files..."

# Function to update swagger HTML for each service
update_swagger_html() {
    local service_name=$1
    local swagger_internal="services/${service_name}/internal/server/swagger.html"
    local swagger_docs="services/${service_name}/docs/swagger.html"
    
    # Check if proto file exists for this service
    local proto_file="api/${service_name}/v1/${service_name}.proto"
    if [ ! -f "${proto_file}" ]; then
        echo "⚠️  ${service_name}: proto file not found, skipping"
        return
    fi
    
    # Generate OpenAPI spec for this specific service using protoc directly
    echo "📝 Generating OpenAPI for ${service_name}..."
    
    # Use protoc with buf's modules
    buf generate --path "${proto_file}" --template <(cat <<EOF
version: v2
plugins:
  - local: protoc-gen-openapi
    out: /tmp
    opt:
      - naming=proto
      - output_mode=source_relative
EOF
)
    
    if [ ! -f "/tmp/api/${service_name}/v1/${service_name}.openapi.yaml" ] && [ ! -f "/tmp/${service_name}.openapi.yaml" ] && [ ! -f "/tmp/openapi.yaml" ]; then
        echo "❌ Failed to generate OpenAPI for ${service_name}"
        return
    fi
    
    # Find the generated file
    local openapi_file=""
    for possible_path in "/tmp/api/${service_name}/v1/${service_name}.openapi.yaml" "/tmp/${service_name}.openapi.yaml" "/tmp/openapi.yaml"; do
        if [ -f "${possible_path}" ]; then
            openapi_file="${possible_path}"
            break
        fi
    done
    
    if [ -z "${openapi_file}" ]; then
        echo "❌ Could not find generated OpenAPI file for ${service_name}"
        return
    fi
    
    # Convert YAML to JSON and update HTML files
    python3 -c "
import yaml
import json
import sys

# Read YAML spec
with open('${openapi_file}', 'r') as f:
    spec = yaml.safe_load(f)

# Output as formatted JSON
print(json.dumps(spec, indent=2))
" > /tmp/${service_name}-spec.json
    
    # Update both swagger HTML files
    for html_file in "${swagger_internal}" "${swagger_docs}"; do
        if [ ! -f "${html_file}" ]; then
            continue
        fi
        
        # Replace spec in HTML
        python3 - "${html_file}" /tmp/${service_name}-spec.json <<'PYEOF'
import sys
import json

html_file = sys.argv[1]
spec_file = sys.argv[2]

# Read new spec
with open(spec_file, 'r') as f:
    new_spec = f.read().strip()

# Read HTML
with open(html_file, 'r') as f:
    content = f.read()

# Find spec: { ... },
spec_marker = 'spec: '
spec_pos = content.find(spec_marker)
if spec_pos == -1:
    sys.exit(1)

# Find closing brace
parse_start = spec_pos + len(spec_marker)
brace_count = 0
pos = parse_start
in_string = False
found_opening = False

while pos < len(content):
    char = content[pos]
    
    if char == '"' and (pos == 0 or content[pos-1] != '\\\\'):
        in_string = not in_string
    elif not in_string:
        if char == '{':
            if not found_opening:
                found_opening = True
            else:
                brace_count += 1
        elif char == '}':
            if brace_count == 0:
                break
            brace_count -= 1
    pos += 1

spec_end = pos + 1

# Replace
new_content = content[:spec_pos] + 'spec: ' + new_spec + ',' + content[spec_end:]

# Write
with open(html_file, 'w') as f:
    f.write(new_content)
PYEOF
        
        if [ $? -eq 0 ]; then
            echo "   ✅ Updated $(basename $(dirname ${html_file}))/$(basename ${html_file})"
        fi
    done
    
    # Cleanup
    rm -f "${openapi_file}" /tmp/${service_name}-spec.json
}

# Update each service
for service in inventory user product; do
    update_swagger_html "${service}"
done

echo ""
echo "🎉 Done!"
echo "💡 Swagger files have been updated. Rebuild your services to embed the changes:"
echo "   make build"
