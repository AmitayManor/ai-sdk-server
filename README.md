# AI SDK Server

A robust Go-based server application that powers the AI SDK Client Library. This server provides authentication, text generation, and request management services through a RESTful API interface. The server integrates with Supabase for data storage and user management, and uses Hugging Face's inference API for AI text generation.

## Architecture Overview

The server is built using a modular architecture with clear separation of concerns:

```
ai-sdk-server/
├── config/              # Configuration management
├── handlers/            # Request handlers
├── middleware/          # HTTP middleware components
├── models/             # Data models and types
└── supabase/           # Supabase edge functions
    └── functions/
        └── hf-inference/  # Hugging Face inference handler
```

### Key Components

The server consists of several key components working together:

1. **Main Server (Go)**
   - Handles HTTP requests and routing
   - Manages authentication and authorization
   - Coordinates with Supabase for data persistence
   - Processes model generation requests

2. **Edge Functions (Deno/TypeScript)**
   - Handles AI model inference
   - Processes text generation requests
   - Communicates with Hugging Face API

3. **Supabase Integration**
   - User authentication and management
   - Request history storage
   - Real-time updates

## Prerequisites

To run or develop the server, you'll need:

- Go 1.23 or higher
- Deno runtime
- Supabase account
- Hugging Face API key
- Docker (for containerization)
- Koyeb account (for deployment)

## Environment Configuration

Create a `.env` file based on `.env.example` with the following variables:

```env
SUPABASE_URL=""              # Your Supabase project URL
SUPABASE_ANON_KEY=""         # Supabase anonymous key
SUPABASE_SERVICE_ROLE_KEY="" # Supabase service role key
SUPABASE_JWT_SECRET=""       # JWT secret for token validation
DATABASE_URL=""              # PostgreSQL connection string
STORAGE_URL=""              # Storage URL for Supabase
STORAGE_REGION=""           # Storage region
STORAGE_ACCESS_KEY_ID=""    # Storage access key
STORAGE_SECRET_ACCESS_KEY="" # Storage secret key
ALLOWED_ORIGINS=""          # CORS allowed origins
PORT=""                     # Server port (default: 3000)
HF_API_TOKEN=""            # Hugging Face API token
```

## Local Development

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ai-sdk-server.git
cd ai-sdk-server
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up Supabase:
   - Create a new Supabase project
   - Set up the required tables (see Database Schema section)
   - Configure the edge functions

4. Run the server:
```bash
go run main.go
```

## Database Schema

The server requires the following Supabase table structure:

```sql
-- Model requests table
CREATE TABLE model_requests (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    user_id UUID REFERENCES auth.users(id),
    status TEXT,
    input_data JSONB,
    output_data JSONB,
    error_msg TEXT,
    token_used INTEGER,
    token_count INTEGER,
    processing_time INTERVAL,
    model_type TEXT
);

-- Add necessary indexes
CREATE INDEX idx_model_requests_user_id ON model_requests(user_id);
CREATE INDEX idx_model_requests_status ON model_requests(status);
```

## API Endpoints

### Authentication

```
POST /api/auth/signup
POST /auth/signin
```

### Model Requests

```
POST /api/requests    # Create new generation request
GET /api/requests     # List user's requests
```

### Request/Response Examples

#### Signup Request
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Text Generation Request
```json
{
  "model_type": "text2text",
  "input_data": {
    "prompt": "Write a story about..."
  }
}
```

## Deployment

The server is deployed on Koyeb using Docker containerization.

### Docker Build

```bash
docker build -t ai-sdk-server .
```

### Koyeb Deployment

1. Push your Docker image to a registry
2. Create a new Koyeb app
3. Configure environment variables in Koyeb
4. Deploy using the Koyeb CLI or web interface

```bash
# Using Koyeb CLI
koyeb app init ai-sdk-server
koyeb service create main --docker your-registry/ai-sdk-server:latest
```

### Deployment Configuration

Configure the following on Koyeb:

- Environment variables from your `.env` file
- Port mapping (default: 3000)
- Resource allocation (recommended: 512MB RAM minimum)
- Health check endpoint: `/health`

## Edge Functions

The Supabase edge functions handle AI model inference. To deploy them:

1. Navigate to the functions directory:
```bash
cd supabase/functions
```

2. Deploy using Supabase CLI:
```bash
supabase functions deploy hf-inference
```

### Edge Function Configuration

Set the following secrets in Supabase:

```bash
supabase secrets set HF_API_TOKEN=your_huggingface_token
supabase secrets set MODEL_TIMEOUT=90000
supabase secrets set MAX_RETRIES=3
```

## Monitoring and Logging

The server implements comprehensive logging:

- Request/response logging via Fiber middleware
- Error logging with stack traces
- Edge function execution logs in Supabase
- Koyeb deployment logs

Access logs through:
- Koyeb dashboard
- Supabase dashboard
- Server stdout/stderr

## Security

The server implements several security measures:

- JWT-based authentication
- Input validation middleware
- CORS protection
- Rate limiting
- Secure environment variable handling

## Error Handling

The server provides detailed error responses:

```json
{
  "error": "Error description",
  "detail": "Detailed error information",
  "code": "ERROR_CODE"
}
```

Common error codes and their meanings:
- `AUTH_ERROR`: Authentication failed
- `VALIDATION_ERROR`: Invalid input
- `MODEL_ERROR`: AI model processing failed
- `SERVER_ERROR`: Internal server error

## Performance Considerations

- The server uses connection pooling for database operations
- Long-running operations are handled asynchronously
- Edge functions have configurable timeouts
- Rate limiting is implemented for API endpoints

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

Please ensure your code follows the existing style and includes appropriate tests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and support:
1. Check existing GitHub issues
2. Create a new issue with:
   - Detailed description
   - Steps to reproduce
   - Environment information
   - Relevant logs

## Acknowledgments

- Supabase team for the authentication and database infrastructure
- Hugging Face for the AI model inference API
- Koyeb for hosting services
