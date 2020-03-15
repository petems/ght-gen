Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/ght-gen-int-test ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of ght-gen
    When I run `bin/ght-gen-int-test` interactively
    And I type "bob"
    And I type "password"
    Then the output should contain: 
      """
      Github Username: 
      Github Password: 
      """
  
  Scenario:
    Given a build of ght-gen
    When I run `bin/ght-gen-int-test` interactively
    And I type "bob"
    And I type "password"
    Then the output should contain: 
      """"
      Github Username: 
      Github Password: 
      """"