API documentation for weather-api:

1. **POST /api/login**
   - Description: Authenticate a user and return a JWT token.
   - Body: JSON object with `username` and `password`.
   - Returns: A JWT token in the `Authorization` header and in the response body as `{"token": "JWT_TOKEN"}`. 

2. **POST /api/register**
   - Description: Register a new user.
   - Body: JSON object with `username`, `password` and `birth_date`.
   - Returns: A JWT token in the `Authorization` header.

3. **GET /api/weather?city={city_name}**
   - Description: Fetch weather data for a given city.
   - Query parameters: `city` - the city name to get the weather for.
   - Returns: A JSON object with the weather data for the given city.

4. **GET /api/history**
   - Description: Fetch the logged-in user's weather search history.
   - Returns: A JSON array of the user's past weather searches.

5. **PUT /api/history/update**
   - Description: Update a specific weather search history record for the logged-in user.
   - Body: JSON object with updated weather history data. The required fields will depend on what fields you decide to allow updates for (see previous discussion).
   - Returns: A success message if the update was successful.

6. **DELETE /api/history/delete?weatherID={weatherID}**
   - Description: Delete a specific weather search history record for the logged-in user.
   - Query parameters: `weatherID` - the ID of the weather history record to delete.
   - Returns: A success message if the deletion was successful.

7. **DELETE /api/history/bulkdelete**
   - Description: Delete multiple weather search history records for the logged-in user.
   - Body: JSON array of `weatherID`s to delete.
   - Returns: A success message if the deletions were successful.
