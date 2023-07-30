# Weather API

This is a simple Weather API built with Golang that provides weather data for various cities and allows users to register, login, and keep track of their weather search history.
Live API URL : https://api.duranz.in

## Functionality

The Weather API provides the following endpoints:

1. **POST /api/login**

   - Description: Authenticate a user and return a JWT token.
   - Body: JSON object with `username` and `password`.
   - Returns: A JWT token in the `Authorization` header and in the response body as `{"token": "JWT_TOKEN"}`.

2. **POST /api/register**

   - Description: Register a new user.
   - Body: JSON object with `username`, `password`, and `birth_date`.
   - Returns: A JWT token in the `Authorization` header.

3. **GET /api/weather?city={city_name}**

   - Description: Fetch weather data for a given city.
   - Query parameters: `city` - the city name to get the weather for.
   - Returns: A JSON object with the weather data for the given city.

4. **GET /api/history**

   - Description: Fetch the logged-in user's weather search history.
   - Returns: A JSON array of the user's past weather searches.


5. **DELETE /api/history/delete?weatherID={weatherID}**

   - Description: Delete a specific weather search history record for the logged-in user.
   - Query parameters: `weatherID` - the ID of the weather history record to delete.
   - Returns: A success message if the deletion was successful.

6. **DELETE /api/history/bulkdelete**

   - Description: Delete multiple weather search history records for the logged-in user.
   - Body: JSON array of `weatherID`s to delete.
   - Returns: A success message if the deletions were successful.

## Setup Instructions

To run the Weather API on your machine, follow these instructions:

1. Ensure you have Golang (1.18+) installed on your system. If not, download and install it from the official Golang website: https://golang.org/

2. Clone this repository to your local machine using Git:

   ```bash
   git clone https://github.com/KunalDuran/weather-api.git
   ```

3. Change to the project directory:

   ```bash
   cd weather-api
   ```

4. Initialize the project and download dependencies using Go modules:

   ```bash
   go mod tidy
   ```

5. Create a `.env` file in the project root directory with the following variables:

   ```plaintext
   DB_USER=mysql_database_user
   DB_PASS=mysql_database_password
   DB_HOST=mysql_database_host
   DB_PORT=mysql_database_port
   DB_NAME=weather
   API_KEY=your_openweathermap_API_key
   ```

   Replace the values with your database credentials and the API key you obtained for accessing weather data (e.g., from OpenWeather API).

6. Build the application:

   ```bash
   go build
   ```

7. Run the application:

   ```bash
   ./weather-api
   ```

8. The API will be running at `http://localhost:8080`. You can now use API endpoints as described in the "Functionality" section above.

## Authentication

The Weather API uses JWT (JSON Web Tokens) for authentication. When a user logs in or registers, a JWT token is generated and returned, which should be included in the `Authorization` header for subsequent requests to protected endpoints.

## Database

This API uses MySQL as the Database.
Creation of Database and Tables is done automatically by the API.

## Conclusion

The Weather API allows users to register, log in, fetch weather data for cities, and manage their weather search history. We integrated this API into our Weather application available on https://github.com/KunalDuran/weather-reactjs