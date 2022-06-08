Feature: Server Tests
  Scenario: Test GET Request for url /store
    When I send "GET" request to "/store"
    Then The response should be 200

  Scenario: Test POST Request for url /store
    When I send "POST" request to "/store" with payload "id: 10"
    Then The response should be 200

  Scenario: Test POST Request for url /store with incorrect properties
    When I send "POST" request to "/store" with payload ""
    Then The response should be 400

  Scenario: Test GET request for url /store/{id}
    Given I registered store with 20
    When I send "GET" request to "/store/20"
    Then The response should be 200

  Scenario: Test GET request for url /store/{id} with false id
    When I send "GET" request to "/store/100"
    Then The response should be 404

  Scenario: Test PUT request for url /store/{id}
    When I send "PUT" request to "/store/20" with payload "name: Pet shop"
    Then The response should be 200

  Scenario: Test PUT request for url /store/{id}
    When I send "PUT" request to "/store/45"
    Then The response should be 404

  Scenario: Test PUT request for url /store/{id} with false property
    When I send "PUT" request to "store/45" with payload "football"
    Then The response should be 400
