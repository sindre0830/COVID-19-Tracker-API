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

[1.12.0] development: Added early version of Policy Handler

[1.13.0] development: Fixed stringency output from PolicyHistory

[1.13.1] formatting: Renamed structure types from data to their relevant name

[1.14.0] development: Added error catching if alphacode isn't found

[1.14.1] formatting: Added policy and cases packages and moved testing to the same folder as their package

[1.15.0] development: Added LimitDecimals in functionality package

[1.15.1] formatting: Modified comments and error messages in policy package

[1.15.2] formatting: Modified comments and error messages in cases package

[1.16.0] development: Added testing to CasesHistory and CasesTotal

[1.17.0] development: Fixed array size in policy testing package

[1.18.0] development: Added testing to Policy Handler

[1.18.1] formatting: Modified format of error output in CMD for readability

[1.19.0] development: Added dictionary to deal with country edge cases

**Milestone Reached:** Policy endpoint has been implement

[2.0.0] setup: Initial commit of diagnosis structure

[2.1.0] development: Added update function to diagnosis

[2.2.0] development: Added request function to diagnosis

[2.3.0] development: Added get function to diagnosis

[2.4.0] development: Added Handler function to diagnosis

[2.5.0] development: Changed debug package to use structure functionality for consistency

[2.6.0] development: Moved CountryNameDetails structure to it's own package and changed status code to not found if country name is not found

**Milestone Reached:** Diagnosis endpoint has been implement

[3.0.0] setup: Initial commit of notification structure

[3.0.1] setup: Initial commit of notificationGetAll structure

[3.0.2] setup: Initial commit of notificationGetOne structure

[3.0.3] setup: Initial commit of webhook file in api package

[3.1.0] development: Early version of MethodHandler with POST and GET functionality (copy from teachor for testing purposes)

[3.2.0] development: Early version of ServiceHandler with POST functionality (copy from teachor for testing purposes)

[3.2.1] formatting: Restructuring notification package

[3.3.0] development: Added feedback structure and error handling in notification package

[3.4.0] development: Added RandomString function in the fun package

[3.5.0] development: Added ID in feedback message

[3.5.1] formatting: Restructuring notification package

[3.6.0] development: Added country validation in Notification POST request

[3.7.0] development: Added DELETE request in Notification

[3.8.0] development: Removed 'field' attribute from Notification

[3.9.0] development: Added GET all and GET specific requests in Notification

[3.10.0] development: Added simple URL validation in Notification POST request

[3.11.0] development: Added field validation in Notification POST request

[3.12.0] development: Added trigger validation in Notification POST request

[3.13.0] development: Added webhook length in Diagnosis Registered

[3.14.0] development: Added method handler for 'cases' package

[3.15.0] development: Added method handler for 'policy' package

[3.16.0] development: Added method handler for 'diag' package

[3.17.0] development: Requesting data will now give proper status code when timedout

[3.17.1] formatting: Moved code out from casesHistory into cases to keep consistency

[3.18.0] development: Removed start- and end date from CasesHistory

[3.19.0] development: Removed garbage and added better error messages in Cases

[3.20.0] development: Renamed global error struct and moved dynamic 'reason' in 'debug' package

[3.21.0] development: Removed 'Updated' field from Policy structure

[3.22.0] development: Fixed latency in Policy data and added fall back if StringencyActual isn't valid

[3.23.0] development: Removed comprehensive checks in policy to give stringency -1 instead of error message

[3.24.0] development: Fixed testing of policy after last commit

**Milestone Reached:** Local deployment of webhooks implemented

[4.0.1] setup: Removed garbage file

[4.0.2] setup: Added firebase key to gitignore

[4.0.3] setup: Added firebase to go.mod and go.sum

[4.1.0] development: Added firebase package and moved firebase init to it's own setup function

[4.2.0] development: Added 'Add' functionality to 'firebase' package

[4.3.0] development: POST request in Notification adds data to firebase

[4.4.0] development: Added simple Get function to firbase package

[4.4.1] setup: Moved firebase package content into Notification package because of import cycle

[4.5.0] development: Fixed database Get function by setting notifications to be local

[4.6.0] development: Added Delete function to Database

[4.7.0] development: Added Schedule loop

[4.8.0] development: Fixed Policy to only remove days if it's within the buffer of 10 days

[4.9.0] development: Changed Notifications to be global again to be used in the Schedule loop

[4.10.0] development: Notifications is now only updated at startup and during changes in database

[4.11.0] development: Firebase only updates ID now on second interaction in Add

[4.12.0] development: Added URL calling depending on schedule

[4.13.0] development: Added functionality to ON_TIMEOUT and ON_CHANGE

[4.14.0] development: Added proper status code when requesting an empty list of webhooks

[4.14.1] setup: Removed garbage

[4.14.2] formatting: Reformatted Diagnosis and added comments

[4.14.3] formatting: Added comments to method handler for diag endpoint

[4.14.4] formatting: Reformatted 'countryinfo' package and added comments

[4.14.5] formatting: Reformatted 'debug' package and added comments

[4.14.6] formatting: Reformatted 'dict' package and added comments

[4.14.7] formatting: Reformatted 'fun' package and added comments
