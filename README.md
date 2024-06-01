## Stellar Journal

Service that allows you to store and view the daily image and metadata from the NASA APOD API.


## Installation

1. Clone the repository
2. Set all the environment variables from docker-compose.yml file(APP_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, CONFIG_PATH)
3. Add yaml file with the following structure to the CONFIG_PATH directory like this:
```yaml
env: your environment (local, dev, prod)
http_server:
  host: 0.0.0.0:8087
  read_timeout: 4s
  write_timeout: 4s
  idle_timeout: 60s
ctx_timeout: 5s
storage:
  db_uri: "user=your_user dbname=your_db_name sslmode=disable password=your_pass host=postgresql"
nasa_api:
  host: "https://api.nasa.gov"
  token: "your_token" // you can get it from https://api.nasa.gov/
```

4. Run docker-compose up

## Usage

1. Go to http://localhost:8123/journal to see the list of images and metadata
2. Go to http://localhost:8123/journal/{date} to see the image and metadata for the specific date(date format: YYYY-MM-DD)