#!/bin/bash

# Simple Swagger Update - No Python parsing, use protoc-gen-openapi properly
set -e

echo "🔄 Updating Swagger documentation..."
echo ""

# Generate OpenAPI specs using buf for each service separately
generate_openapi_for_service() {
    local service=$1
    local proto_path="api/${service}/v1/${service}.proto"
    
    if [ ! -f "${proto_path}" ]; then
        return 1
    fi
    
    echo "📝 Generating OpenAPI for ${service}..."
    
    # Run protoc-gen-openapi for this service only
    buf generate --path "${proto_path}"
    
    # The output should be in gen/openapi/
    return 0
}

# Update swagger HTML with generated spec
update_swagger_html() {
    local service=$1
    local openapi_yaml="gen/openapi/api/${service}/v1/${service}.openapi.yaml"
    
    # Check multiple possible locations
    for possible in "gen/openapi/${service}.openapi.yaml" "gen/openapi/openapi.yaml" "gen/openapi/api/${service}/v1/${service}.openapi.yaml"; do
        if [ -f "${possible}" ]; then
            openapi_yaml="${possible}"
            break
        fi
    done
    
    if [ ! -f "${openapi_yaml}" ]; then
        echo "⚠️  OpenAPI spec not found for ${service}"
        return 1
    fi
    
    echo "✅ Found spec: ${openapi_yaml}"
    
    # Convert to JSON and update HTML
    # ... rest of update logic
}

# Main
buf generate

echo ""
echo "Checking generated OpenAPI files..."
find gen/openapi -type f

echo ""
echo "🎉 Done!"
