Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/ght-gen-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: Version with no flags
    Given a build of ght-gen
    When I run `bin/ght-gen-int-test --version`
    Then the output should contain:
      """""

      ght-gen version development-0.1.0
      """""