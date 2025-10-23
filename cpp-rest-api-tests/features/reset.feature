Feature: Reset Endpoint

  Scenario: Reset clears all records and resets ID counter
    Given the API is running
    And a contact exists with ID 1
    And a contact exists with first name "Jane" and phone "5559876543"
    When I send a DELETE request to "/reset"
    Then the response status code should be 204
    And the database should be empty
