CHIRPY
An API service for chirping and managing chirpers. Chirp away!

HOW TO INSTALL
Before running this application, you will need to install and setup:

Environment variables

- `JWT_SECRET`: A secret key used for generating JSON Web Tokens
- `POLKA_API_KEY`: The API Key for Polka

Pull dependencies with go get.
Load your environment variables into an .env file.

HOW TO RUN
Once everything is set up, you can run the service using go run main.go.

FEATURES
/app, /app/* - Fileserver endpoint which hits the middleware incrementing metrics.
/healthz - Health check endpoint
/reset - Endpoint to reset the state of the DB
/revoke - Endpoint to revoke JWT access
/refresh - Refresh JWT access
/login - Login endpoint
User management POST /users - Create a new user, PUT /users - Update user data
Manage polka webhooks - POST /polka/webhooks
POST /chirps, GET /chirps, /chirps/{chirpID} - Endpoints to create, retrieve and manage chirps
/admin/metrics - View app metrics

DEBUG MODE
This application has a debug mode. Run go run main.go -debug to enable debug mode. This will also reset the database.

NOTES
Make sure all your environment variables are loaded correctly
The database used is a JSON database connected at runtime, the location is database.json

HAPPY CHIRPING!
