Feature: Contact Management API
  Background:
    Given the API is running
    And the database should be empty
  Scenario: Create a new contact
    When I send a POST request to "/records" with contact details:
      """
      {
        "first_name": "John",
        "last_name": "Doe",
        "phone": "1234567890",
        "email": "john@example.com"
      }
      """
    Then the response status code should be 201
    And the response should contain "John"
  Scenario: Retrieve all contacts
    Given I have created 2 contacts
    When I send a GET request to "/records"
    Then the response status code should be 200
    And the response should contain 2 contacts
  Scenario: Retrieve contact by ID
    Given I have created a contact with ID 1
    When I send a GET request to "/records/{lastCreatedID}"
    Then the response status code should be 200
    And the response should contain "John"
  Scenario: Update existing contact
    Given I have created a contact with ID 1
    When I send a PUT request to "/records/{lastCreatedID}" with updated details:
      """
      {
        "first_name": "Jane"
      }
      """
    Then the response status code should be 200
    And the response should contain "Jane"
  Scenario: Delete existing contact
    Given I have created a contact with ID 1
    When I send a DELETE request to "/records/{lastCreatedID}"
    Then the response status code should be 204
  Scenario: Query by phone area code
    Given I have created a contact with phone "1234567890"
    When I send a GET request to "/records?phone=123"
    Then the response status code should be 200
    And the response should contain 1 contact

