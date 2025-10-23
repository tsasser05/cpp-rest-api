# Contact Management API

A lightweight C++ RESTful API for managing contact records, built with [Pistache](https://github.com/pistacheio/pistache) and [nlohmann/json](https://github.com/nlohmann/json). The API supports CRUD operations (Create, Read, Update, Delete) and querying of contact data, including Name (first, middle, last), Address (street, city, state, zip), Phone Number, and Email. Data is stored in-memory and does not persist across server restarts.

## Authors

Grok did the work.  I just guided it to see if it would actually make a working API.

## Features

- **CRUD Operations**:
  - **Create**: `POST /records` to add a new contact.
  - **Read**: `GET /records/:id` to retrieve a contact by ID.
  - **Read**: `GET /records/` to retrieve all contacts.
  - **Update**: `PUT /records/:id` to modify a contact.
  - **Delete**: `DELETE /records/:id` to remove a contact.
  - **Reset Endpoint**: `DELETE /reset` clears all records and resets the ID counter to 1.
  
- **Querying**: `GET /records?param=value` supports exact matches on `first_name`, `middle_name`, `last_name`, `street`, `city`, `state`, `zip`, `phone`, and `email`. Phone queries support full numbers or 3-digit area codes. Multiple parameters are combined with AND logic.

- **Data Fields**:
  - `id`: Integer (auto-assigned, starting from 1).
  - `first_name`, `middle_name`, `last_name`, `street`, `city`, `state`, `zip`, `phone`, `email`: Strings (optional).

- **Response Format**: JSON with standard HTTP status codes (201, 200, 400, 404, 204).

## Prerequisites

- **macOS** with Homebrew installed:
  - **Intel Macs**: Homebrew in `/usr/local` (default for Intel-based systems).
  - **Apple Silicon (arm64) Macs**: Homebrew in `/opt/homebrew` (default for M1/M2 systems).
- **Compiler**: `g++` (Clang) with C++17 support, included in Xcode Command Line Tools.
- **Libraries**:
  - Pistache 0.4.26 (C++ REST framework).
  - nlohmann/json (JSON parsing).

## Installation

1. **Install Dependencies via Homebrew**:
   ```bash
   brew install pistache nlohmann-json
   ```

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/tsasser05/cpp-rest-api
   cd cpp-rest-api
   ```

3. **Verify Dependencies**:
   Check libraries and headers based on your system:
   - **Intel Macs**:
     ```bash
     ls /usr/local/lib/libpistache*  # Should show libpistache.a or .dylib
     ls /usr/local/include/nlohmann  # Should show json.hpp
     ```
   - **Apple Silicon (arm64) Macs**:
     ```bash
     ls /opt/homebrew/lib/libpistache*  # Should show libpistache.a or .dylib
     ls /opt/homebrew/include/nlohmann  # Should show json.hpp
     ```

## Build

Compile the API based on your system:
- **Intel Macs**:
  ```bash
  g++ -std=c++17 main.cpp -o api -lpistache -lpthread -I/usr/local/include -L/usr/local/lib
  ```
- **Apple Silicon (arm64) Macs**:
  ```bash
  g++ -std=c++17 main.cpp -o api -lpistache -lpthread -I/opt/homebrew/include -L/opt/homebrew/lib
  ```

## Usage

1. **Run the Server**:
   ```bash
   ./api
   ```
   - The server listens on `http://localhost:8080`.
   - Stop with `Ctrl+C`.

2. **Test with curl**:
   ```bash
   # Create a contact
   curl -X POST http://localhost:8080/records -H "Content-Type: application/json" \
     -d '{"first_name":"John","last_name":"Doe","phone":"1234567890","email":"john@example.com","street":"123 Main St","city":"Anytown","state":"CA","zip":"12345"}'

   # Query by first name and area code
   curl "http://localhost:8080/records?first_name=John&phone=123"

   # Read a contact by ID
   curl http://localhost:8080/records/1

   # Get all contacts with pretty print
   curl -s http://localhost:8080/records/ | jq .
   
   # Update a contact
   curl -X PUT http://localhost:8080/records/1 -H "Content-Type: application/json" \
     -d '{"first_name":"Jane"}'

   # Delete a contact
   curl -X DELETE http://localhost:8080/records/1

   # Delete the database
   curl -X DELETE http://localhost:8080/reset
   ```
3. **Run test automation**:

   - Always run the test automation from a new instance of the API.  Restart the API first.

```
   cd cpp-rest-api-tests
   go test -v
   ```
  

## API Endpoints

- **POST /records**: Create a new contact. Returns 201 with the created record.
  - Body: JSON object with optional fields (e.g., `{"first_name":"John","last_name":"Doe"}`).
  - Error: 400 for invalid JSON.

- **GET /records/:id**: Retrieve a contact by ID. Returns 200 with the record or 404 if not found.

- **PUT /records/:id**: Update a contact. Updates only provided fields. Returns 200 with updated record or 404 if not found.
  - Body: JSON with fields to update.
  - Error: 400 for invalid JSON.

- **DELETE /records/:id**: Delete a contact by ID. Returns 204 on success or 404 if not found.

- **GET /records**: Query contacts. Returns 200 with an array of matching records.
  - Query parameters: `first_name`, `middle_name`, `last_name`, `street`, `city`, `state`, `zip`, `phone`, `email`.
  - Phone: Matches full number or 3-digit area code.
  - Example: `GET /records?first_name=John&phone=123`.

- **DELETE /reset**: Clear out the database.  Returns 204 No content.


## Load Contacts

- load_contacts.sh will generate 100 contacts and insert them into the application.  Use this as you will.


## Notes

- **Data Storage**: In-memory only; data is lost on server restart.
- **Error Handling**: Returns standard HTTP status codes and JSON/text error messages.
- **Performance**: Suitable for small-scale use due to in-memory storage.
- **Future Improvements**: Add file or database persistence, input validation (e.g., phone format), or HTTPS support.

## License

MIT License. See [LICENSE](LICENSE) for details.