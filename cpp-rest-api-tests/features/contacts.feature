Feature: Contact Management API
  As a user of the contact management API
  I want to perform CRUD operations and query contacts
  So that I can manage contact information effectively

  Scenario: Create a new contact
    Given the API is running
    When I send a POST request to "/records" with contact details
      """
      {
        "first_name": "John",
        "last_name": "Doe",
        "phone": "1234567890",
        "email": "john@example.com",
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "zip": "12345"
      }
      """
    Then the response status code should be 201
    And the response should contain the contact ID

  Scenario: Read an existing contact
    Given a contact exists with ID {lastCreatedID}
    When I send a GET request to "/records/{lastCreatedID}"
    Then the response status code should be 200
    And the response should contain the contact details
      """
      {
        "id": {lastCreatedID},
        "first_name": "John",
        "last_name": "Doe",
        "middle_name": "",
        "phone": "1234567890",
        "email": "john@example.com",
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "zip": "12345"
      }
      """

  Scenario: Update an existing contact
    Given a contact exists with ID {lastCreatedID}
    When I send a PUT request to "/records/{lastCreatedID}" with updated details
      """
      {
        "first_name": "Jane",
        "email": "jane@example.com"
      }
      """
    Then the response status code should be 200
    And the response should contain the updated contact details
      """
      {
        "id": {lastCreatedID},
        "first_name": "Jane",
        "last_name": "Doe",
        "middle_name": "",
        "phone": "1234567890",
        "email": "jane@example.com",
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "zip": "12345"
      }
      """

  Scenario: Delete a contact
    Given a contact exists with ID {lastCreatedID}
    When I send a DELETE request to "/records/{lastCreatedID}"
    Then the response status code should be 204
    And a subsequent GET request to "/records/{lastCreatedID}" should return 404

  Scenario: Query contacts by first name and phone
    Given a contact exists with first name "John" and phone "1234567890"
    When I send a GET request to "/records?first_name=John&phone=123"
    Then the response status code should be 200
    And the response should contain a list with the contact
      """
      [
        {
          "id": {lastCreatedID},
          "first_name": "John",
          "last_name": "Doe",
          "middle_name": "",
          "phone": "1234567890",
          "email": "john@example.com",
          "street": "123 Main St",
          "city": "Anytown",
          "state": "CA",
          "zip": "12345"
        }
      ]
      """

  Scenario: Attempt to create a contact with missing required fields
    Given the API is running
    When I send a POST request to "/records" with contact details
      """
      {
        "first_name": "John"
      }
      """
    Then the response status code should be 201
