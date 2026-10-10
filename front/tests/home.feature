Feature: Home page
  Scenario: Open the application and find the title
    Given the application is started with the mocked API
    When the user opens the home page
    Then the page contains the title "Slot Finder"
