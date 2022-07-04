Feature: Test Server
Scenario: Test GET Request for url "/store/{id}"
  When I send GET request to "/store/100" with payload ""
  Then The response for url "/store/100" with request method "GET" should be 404