Feature: Home page
  Scenario: Open the application and find the title
    Given the application is started with the mocked API
    When the user opens the home page
    Then the page contains the title "Slot Finder"

  Scenario: Override the mocked auth status to 200
    Given the application is started with auth status 200
    Then the page contains the dashboard heading "My events"
