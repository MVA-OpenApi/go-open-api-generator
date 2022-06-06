Feature: pet store get request is made with url /store

  Scenario: user makes a get request with path /store
    When user makes a get request with the url /store
    Then the server responds with 200 status code
    And return a list of all stores

Feature: pet store post request is made from user with url /store

  Scenario: user makes a post request with path /store
  Creates a new store from the user input form
    When user makes a post request with path /status
    Then the server responds with 200 status code if the information was correct

Feature: respond messages of the server when user wants to retrieve store from the system

  Scenario: user makes a get request with a specific id
    When user make request with url path /store/{id}
    Then server responds with 200
    And sends the store back to the user

  Scenario: user sends an invalid id to the server
    When the id is invalid
    And the request is a get request
    Then the server sends 400 back to the user

  Scenario: user sends a non-existing id to the server
    When the store with the id does not exist
    And the request is a get request
    Then the server sends 404 back to the user

Feature: respond messages and their correctness when user wants to delete entry

  Scenario: user wants to delete entry from the stores
    When the user sends a delete request to the server
    And the id is correct
    Then the server should send 200 back to the user

  Scenario: user sends an incorrect id to the store
    When the user sends a delete request to the server
    And the id is invalid
    Then the server sends 400 back to the user

  Scenario: user sends a valid id which does not exist
    When the user sends a delete request
    And the id does not exist
    Then the server sends 404 back to the user

Feature: respond messages and their correctness when user wants to update entry

  Scenario: user wants to update an entry
    When the user sends an update request to the server
    And the id is valid
    And the id exists
    Then the server replies wth 200

  Scenario: user sends an invalid id in order to update an entry
    When the user sends a put request
    And the id is not valid
    Then the server replies with 400

  Scenario: user sends a non existing id
    When the user sends a put request
    And the id is valid
    And the id does not exist
    Then the server replies with 404