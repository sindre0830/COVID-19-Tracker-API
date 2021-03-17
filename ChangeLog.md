# Changelog for assignment-2

[0.0.0] setup: Initial commit

[0.0.1] formatting: Removed old URL from README

[0.0.2] setup: Added GITIGNORE to ignore executable files

[0.1.0] development: Added structure for corona cases

[0.1.1] formatting: Changed confirmed/recoverd fields in structure for corona cases to int32

[0.1.2] formatting: Fixed naming mistake in structure for corona cases

[0.2.0] development: Added get request for all corona cases

[0.2.1] formatting: Renamed coronaCases structure to cases

[0.2.2] formatting: Renamed casesRaw structure to casesTotal

[0.3.0] development: Added error catching when creating new request

[0.4.0] development: Added structure for case history

[0.5.0] development: Added request for case history

[0.6.0] development: Added get and addCases for case history

[0.6.1] formatting: Added same get/req structure in casesTotal as in casesHistory for consistency

[0.7.0] development: Added getTotal and getHistory in cases

[0.8.0] development: Added get in cases to call on either getTotal or getHistory

[0.9.0] development: Added simple I/O handler in cases

[0.10.0] development: Changed handler in cases to handle HTTP requests

[0.11.0] development: Added error handler (from assignment-1)

[0.12.0] development: Added date validation (from assignment-1)

[0.13.0] development: Added URL parsing in cases

[0.14.0] development: Added proper error handling when getting cases and outputting to user

[0.15.0] development: Added error handling if object is empty

[0.16.0] development: Changed LifeExpectancy field to pointer of interface to deal with Italy edge-case

[0.16.1] formatting: Fixed version mistake in the last three commits

[0.17.0] development: Added functionality to convert country input to correct syntax and validating if it's empty

[0.17.1] formatting: Split structure files into struct-header and struct-functionality files for readability

[0.17.2] formatting: Split structure files into struct-header and struct-functionality files for readability

[0.18.0] development: Changed status code from internal server error to bad request after checking ValidateCountry

[0.19.0] testing: Added cases handler testing

[0.20.0] development: Added proper error handling if service timesout

[0.21.0] development: Added more status handling when getting cases

[0.21.1] formatting: Added comments and headers

[0.21.2] formatting: Modified reciever names to be more descriptive

[0.22.0] development: Added function to handle URL parsing

[0.23.0] development: Fixed Cases handler after last commit

**Milestone Reached:** Cases endpoint has been implement

[1.0.0] setup: Initial commit of policy structure

[1.0.1] setup: Initial commit of policyHistory structure

[1.0.2] setup: Initial commit of policyCurrent structure

[1.0.3] setup: Initial commit of countryCode structure

[1.1.0] development: Added status code check in requestData to return statuscode even if no errors are thrown

[1.2.0] development: Renamed countryCode structure to CountryNameDetails and added get, req and testing

[1.3.0] development: Added Get, req, isEmpty and testing to PolicyCurrent

[1.4.0] development: Increased timeout variable to 4 seconds since some REST services are slow

[1.5.0] development: Added Get, req, isEmpty and testing to PolicyHistory

[1.5.1] formatting: Moved structures to the same file as their functions

[1.6.0] development: Added decreaseTime function in PolicyHistory

[1.6.1] formatting: Added comments in decreaseTime

[1.7.0] development: Moved object declaration in testing functions to always start with an empty object

[1.8.0] development: Added trend calculation

[1.9.0] development: Added get in Policy structure

[1.10.0] development: Added getCurrent and update in Policy structure

[1.11.0] development: Added getHistory in Policy structure
