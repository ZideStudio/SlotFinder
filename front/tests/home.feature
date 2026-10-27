Feature: Home page

  Scenario: A visitor discovers the application
    Given I am not signed in
    When I open the home page
    Then I see the main heading "Slot Finder"

  Scenario: A signed-in user sees their events
    Given I am signed in
    When I open the home page
    Then I see the main heading "My events"
