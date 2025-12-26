#!/usr/bin/env python3
"""
Generate Swagger/OpenAPI spec from proto files
"""

import re
import json
import sys
from pathlib import Path

def parse_proto_file(proto_path):
    """Parse a proto file and extract service and message definitions"""
    with open(proto_path, 'r') as f:
        content = f.read()
    
    # Extract package name
    package_match = re.search(r'package\s+([\w.]+);', content)
    package = package_match.group(1) if package_match else ""
    
    # Extract service name - use better pattern to handle nested braces
    service_start = re.search(r'service\s+(\w+)\s*\{', content)
    if not service_start:
        return None
    
    service_name = service_start.group(1)
    
    # Find matching closing brace for service block
    start_pos = service_start.end()
    brace_count = 1
    pos = start_pos
    
    while pos < len(content) and brace_count > 0:
        if content[pos] == '{':
            brace_count += 1
        elif content[pos] == '}':
            brace_count -= 1
        pos += 1
    
    service_content = content[start_pos:pos-1]
    
    # Extract RPC methods with HTTP annotations
    methods = []
    # Updated pattern to handle multiline RPC definitions
    rpc_pattern = r'rpc\s+(\w+)\s*\(([^)]+)\)\s+returns\s+\(([^)]+)\)\s*\{([^}]+)\}'
    
    for match in re.finditer(rpc_pattern, service_content, re.DOTALL):
        method_name = match.group(1)
        request_type = match.group(2).strip()
        response_type = match.group(3).strip()
        options = match.group(4)
        
        # Extract HTTP method and path directly (simpler and more reliable)
        method_match = re.search(r'(get|post|put|delete|patch):\s*"([^"]+)"', options, re.IGNORECASE)
        if method_match:
            http_method = method_match.group(1).lower()
            http_path = method_match.group(2)
            
            methods.append({
                'name': method_name,
                'request': request_type,
                'response': response_type,
                'http_method': http_method,
                'http_path': http_path,
                'has_body': 'body:' in options or 'body :' in options
            })
    
    # Extract messages
    messages = {}
    message_pattern = r'message\s+(\w+)\s*\{([^}]+)\}'
    
    for match in re.finditer(message_pattern, content):
        msg_name = match.group(1)
        msg_content = match.group(2)
        
        # Extract fields
        fields = {}
        field_pattern = r'([\w.]+)\s+(\w+)\s*=\s*\d+;'
        
        for field_match in re.finditer(field_pattern, msg_content):
            field_type = field_match.group(1)
            field_name = field_match.group(2)
            
            # Map proto types to OpenAPI types
            type_mapping = {
                'string': {'type': 'string'},
                'int32': {'type': 'integer', 'format': 'int32'},
                'int64': {'type': 'integer', 'format': 'int64'},
                'bool': {'type': 'boolean'},
                'double': {'type': 'number', 'format': 'double'},
                'float': {'type': 'number', 'format': 'float'},
            }
            
            if field_type in type_mapping:
                fields[field_name] = type_mapping[field_type]
            else:
                # Reference to another message
                fields[field_name] = {'$ref': f'#/components/schemas/{field_type}'}
        
        messages[msg_name] = {
            'type': 'object',
            'properties': fields
        }
    
    return {
        'service': service_name,
        'package': package,
        'methods': methods,
        'messages': messages
    }

def generate_swagger_spec(proto_data, service_name):
    """Generate OpenAPI/Swagger spec from parsed proto data"""
    
    paths = {}
    
    for method in proto_data['methods']:
        path = method['http_path']
        http_method = method['http_method']
        
        # Build operation
        operation = {
            'summary': method['name'],
            'operationId': method['name'],
            'responses': {
                '200': {
                    'description': 'Successful response',
                    'content': {
                        'application/json': {
                            'schema': {'$ref': f"#/components/schemas/{method['response']}"}
                        }
                    }
                }
            },
            'security': [{'BearerAuth': []}]
        }
        
        # Add request body if needed
        if method['has_body'] and http_method in ['post', 'put', 'patch']:
            operation['requestBody'] = {
                'required': True,
                'content': {
                    'application/json': {
                        'schema': {'$ref': f"#/components/schemas/{method['request']}"}
                    }
                }
            }
        
        # Add path parameters
        path_params = re.findall(r'\{(\w+)\}', path)
        if path_params:
            operation['parameters'] = [
                {
                    'name': param,
                    'in': 'path',
                    'required': True,
                    'schema': {'type': 'integer' if param == 'id' else 'string'}
                }
                for param in path_params
            ]
        
        if path not in paths:
            paths[path] = {}
        paths[path][http_method] = operation
    
    spec = {
        'openapi': '3.0.3',
        'info': {
            'title': f'{service_name.title()} Service API',
            'version': '1.0.0',
            'description': f'API documentation for {service_name} service'
        },
        'servers': [
            {'url': 'http://localhost', 'description': 'Local server'}
        ],
        'paths': paths,
        'components': {
            'schemas': proto_data['messages'],
            'securitySchemes': {
                'BearerAuth': {
                    'type': 'http',
                    'scheme': 'bearer',
                    'bearerFormat': 'JWT'
                }
            }
        },
        'security': [{'BearerAuth': []}]
    }
    
    return spec

def main():
    if len(sys.argv) < 3:
        print("Usage: proto-to-swagger.py <proto_file> <service_name>")
        sys.exit(1)
    
    proto_file = sys.argv[1]
    service_name = sys.argv[2]
    
    proto_data = parse_proto_file(proto_file)
    if not proto_data:
        print(f"Error: Could not parse {proto_file}")
        sys.exit(1)
    
    spec = generate_swagger_spec(proto_data, service_name)
    print(json.dumps(spec, indent=2))

if __name__ == '__main__':
    main()
