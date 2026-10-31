Feature: Home page

  @unauthenticated
  Scenario: A visitor discovers the application
    Given I open the home page
    Then I see the main heading "Slot Finder"

  Scenario: A signed-in user sees their events
    Given I open the home page
    Then I see the main heading "My events"
