Feature: Server Tests
  Scenario: Test GET Request for url /store
    When I send GET request to "/store"
    Then The response for url "/store" with request method "GET" should be 404

  Scenario: Test POST Request for url /store
    When I send POST request to "/store" with payload "id: 10"
    Then The response for url "/store" with request method "POST" should be 404

  Scenario: Test POST Request for url /store with incorrect properties
    When I send POST request to "/store" with payload ""
    Then The response for url "/store" with request method "POST" should be 404

  Scenario: Test GET request for url /store/{id} with false id
    When I send GET request to "/store/100"
    Then The response for url "/store/100" with request method "GET" should be 501

  Scenario: Test PUT request for url /store/{id}
    When I send PUT request to "/store/20" with payload "name: Pet shop"
    Then The response for url "/store/20" with request method "PUT" should be 501

  Scenario: Test PUT request for url /store/{id} improper payload
    When I send PUT request to "/store/45" with payload ""
    Then The response for url "/store/45" with request method "PUT" should be 501

  Scenario: Test DELETE request for url /store/{id}
    When I send DELETE request to "/store/40"
    Then The response for url "/store/40" with request method "DELETE" should be 501