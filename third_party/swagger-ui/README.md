# Readme

URL: https://github.com/swagger-api/swagger-ui/tree/master/dist
License: Apache License, Version 2.0
License File: LICENSE

## Description
Files from dist folder of https://github.com/swagger-api/swagger-ui.
These are dependency-free collection of HTML, Javascript, and CSS assets that
dynamically generate beautiful documentation and sandbox from a
Swagger-compliant API.
Instructions on how to use these:
https://github.com/swagger-api/swagger-ui#how-to-use-it

## Local Modifications
- Updated the url to "../../swagger-spec" as per instructions at:
https://github.com/swagger-api/swagger-ui#how-to-use-it
- Set supportedSubmitMethods: [] in index.html to remove "Try it out" buttons.
- Set validatorUrl to null, to disable validation (which works only for swagger
  2.0) as per https://github.com/swagger-api/swagger-ui#parameters

LICENSE file has been created for compliance purposes.
Not included in original distribution.
