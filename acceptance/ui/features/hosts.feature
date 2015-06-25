@hosts
Feature: Host Management
  In order to use Control Center
  As a CC admin user
  I want to manage hosts

  @login-required @emptyHostsPage
  Scenario: View empty Hosts page
    When I am on the hosts page
    Then I should see "Hosts Map"
      And I should see "Name"
      And I should see "Active"
      And I should see "Resource Pool"
      And I should see "Memory"
      And I should see "RAM Commitment"
      And I should see "No Data Found"
      And I should see "Showing 0 Results"

  @login-required
  Scenario: View Add Host dialog
    When I am on the hosts page
      And I click the Add-Host button
    Then I should see the Add Host dialog
      And I should see "Host and port"
      And I should see the Host and port field
      And I should see "Resource Pool ID"
      And I should see the Resource Pool ID field
      And I should see "RAM Commitment"
      And I should see the RAM Commitment field

  @login-required @emptyHostsPage
  Scenario: Add an invalid host with an invalid name
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with "bogushost"
      And I fill in the Resource Pool field with the default resource pool
      And I fill in the RAM Commitment field with the default RAM commitment
      And I click "Add Host"
    Then I should see "Error"
      And I should see "Bad Request"
      And I should see "No Data Found"
      And I should see "Showing 0 Results"

  @login-required @emptyHostsPage
  Scenario: Add an invalid host with an invalid Resource Pool field
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with the default host name
      And I fill in the RAM Commitment field with the default RAM commitment
      And I click "Add Host"
    Then I should see "Error"
      And I should see "Bad Request"
      And I should see "No Data Found"
      And I should see "Showing 0 Results"

  @login-required @emptyHostsPage
  Scenario: Add an invalid host with an invalid RAM Commitment field
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with the default host name
      And I fill in the Resource Pool field with the default resource pool
      And I fill in the RAM Commitment field with "invalidentry"
      And I click "Add Host"
    Then I should see "Error"
      And I should see "Bad Request"
      And I should see "No Data Found"
      And I should see "Showing 0 Results"

  @login-required @emptyHostsPage
  Scenario: Fill in the hosts dialog and cancel
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with the default host name
      And I fill in the Resource Pool field with the default resource pool
      And I fill in the RAM Commitment field with the default RAM commitment
      And I click "Cancel"
    Then I should see "No Data Found"
      And I should see "Showing 0 Results"
      And I should not see "Success"

  @login-required @emptyHostsPage
  Scenario: Add an valid host
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with the default host name
      And I fill in the Resource Pool field with the default resource pool
      And I fill in the RAM Commitment field with the default RAM commitment
      And I click "Add Host"
    Then I should see "Success"
      And I should see "roei-dev" in the "Name" column
      And I should see "default" in the "Resource Pool" column
      And I should see "Showing 1 Result"

  @login-required @defaultHostPage
  Scenario: Add another valid host
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with "vagrant:4979"
      And I fill in the Resource Pool field with the default resource pool
      And I fill in the RAM Commitment field with "0%"
      And I click "Add Host"
    Then I should see "Success"
      And I should see "roei-dev" in the "Name" column
      And I should see "default" in the "Research Pool" column
      And I should see "vagrant" in the "Name" column
      And I should see "default" in the "Research Pool" column
      And I should see "Showing 2 Results"

  @login-required
  Scenario: Test ascending name sort
    When I am on the hosts page
      And I sort by "Name" in ascending order
    Then the "Name" column should be sorted in ascending order

  @login-required
  Scenario: Test descending name sort
    When I am on the hosts page
      And I sort by "Name" in descending order
    Then the "Name" column should be sorted in descending order

  @login-required
  Scenario: Test ascending status sort
    When I am on the hosts page
      And I sort by "Active" in ascending order
    Then the "Active" column should be sorted with active hosts on the bottom

  @login-required
  Scenario: Test descending status sort
    When I am on the hosts page
      And I sort by "Active" in descending order
    Then the "Active" column should be sorted with active hosts on top

  @login-required
  Scenario: Test descending resource pool sort
    When I am on the hosts page
      And I sort by "Resource Pool" in descending order
    Then the "Resource Pool" column should be sorted in descending order

  @login-required
  Scenario: Test ascending resource pool sort
    When I am on the hosts page
      And I sort by "Resource Pool" in ascending order
    Then the "Resource Pool" column should be sorted in ascending order

  @login-required
  Scenario: Test descending memory sort
    When I am on the hosts page
      And I sort by "Memory" in descending order
    Then the "Memory" column should be sorted in descending order

  @login-required
  Scenario: Test ascending memory sort
    When I am on the hosts page
      And I sort by "Memory" in ascending order
    Then the "Memory" column should be sorted in ascending order

  @login-required
  Scenario: Test ascending CPU cores sort
    When I am on the hosts page
      And I sort by "CPU Cores" in ascending order
    Then the "CPU Cores" column should be sorted in ascending order

  @login-required
  Scenario: Test descending CPU cores sort
    When I am on the hosts page
      And I sort by "CPU Cores" in descending order
    Then the "CPU Cores" column should be sorted in descending order

  @login-required
  Scenario: Test ascending kernel version sort
    When I am on the hosts page
      And I sort by "Kernel Version" in ascending order
    Then the "Kernel Version" column should be sorted in ascending order

  @login-required
  Scenario: Test descending kernel version sort
    When I am on the hosts page
      And I sort by "Kernel Version" in descending order
    Then the "Kernel Version" column should be sorted in descending order

  @login-required
  Scenario: Test ascending CC release sort
    When I am on the hosts page
      And I sort by "CC Release" in ascending order
    Then the "CC Release" column should be sorted in ascending order

  @login-required
  Scenario: Test descending CC release sort
    When I am on the hosts page
      And I sort by "CC Release" in descending order
    Then the "CC Release" column should be sorted in descending order

  @login-required @defaultHostPage
  Scenario: Add a duplicate host
    When I am on the hosts page
      And I click the Add-Host button
      And I fill in the Host Name field with "172.17.42.1:4979"
      And I fill in the Resource Pool field with "default"
      And I fill in the RAM Commitment field with "50%"
      And I click "Add Host"
    Then I should see "Error"
      And I should see "Internal Server Error: host already exists"

  @login-required @defaultHostPage
  Scenario: Remove a host
    When I am on the hosts page
      And I click "Delete"
    Then I should see "This action will permanently delete the host"
    When I click "Remove Host"
    Then I should see "Removed host"
      And I should see "No Data Found"
      And I should see "Showing 0 Results"

  @login-required
  Scenario: View Hosts Map
    When I am on the hosts page
      And I click "Hosts Map"
    Then I should see "By RAM"
      And I should see "By CPU"
      And I should not see "Active"