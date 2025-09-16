#!/usr/bin/env python3
import json
import uuid
from datetime import datetime
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

class MockDirtCloudHandler(BaseHTTPRequestHandler):
    # Sample metadata store
    metadata_store = [
        {
            "id": "meta-001",
            "path": "app/config/database_url",
            "value": "postgresql://user:pass@localhost:5432/mydb",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z"
        },
        {
            "id": "meta-002", 
            "path": "app/features/new_ui_enabled",
            "value": "true",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z"
        },
        {
            "id": "meta-003",
            "path": "app/config/api_key",
            "value": "secret-api-key-123",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z"
        }
    ]

    def do_POST(self):
        if self.path == '/v1/projects':
            self.handle_create_project()
        elif self.path == '/v1/metadata':
            self.handle_create_metadata()
        else:
            self.send_error(404)
    
    def do_GET(self):
        if self.path.startswith('/v1/projects/'):
            project_id = self.path.split('/')[-1]
            self.handle_get_project(project_id)
        elif self.path.startswith('/v1/metadata/'):
            metadata_id = self.path.split('/')[-1]
            self.handle_get_metadata(metadata_id)
        elif self.path.startswith('/v1/metadata'):
            self.handle_list_metadata()
        else:
            self.send_error(404)
    
    def do_DELETE(self):
        if self.path.startswith('/v1/projects/'):
            project_id = self.path.split('/')[-1]
            self.handle_delete_project(project_id)
        elif self.path.startswith('/v1/metadata/'):
            metadata_id = self.path.split('/')[-1]
            self.handle_delete_metadata(metadata_id)
        else:
            self.send_error(404)
    
    def handle_create_project(self):
        content_length = int(self.headers.get('Content-Length', 0))
        if content_length > 0:
            post_data = self.rfile.read(content_length)
            try:
                data = json.loads(post_data.decode('utf-8'))
            except json.JSONDecodeError:
                data = {}
        else:
            data = {}
        
        project = {
            "id": str(uuid.uuid4()),
            "name": data.get("name", "default-project"),
            "created_at": datetime.now().isoformat() + "Z",
            "updated_at": datetime.now().isoformat() + "Z"
        }
        
        self.send_response(201)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(project).encode())
    
    def handle_get_project(self, project_id):
        project = {
            "id": project_id,
            "name": "test-project",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z"
        }
        
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(project).encode())
    
    def handle_delete_project(self, project_id):
        self.send_response(204)
        self.end_headers()

if __name__ == '__main__':
    server = HTTPServer(('localhost', 8080), MockDirtCloudHandler)
    print("Mock DirtCloud server running on http://localhost:8080")
    server.serve_forever()
