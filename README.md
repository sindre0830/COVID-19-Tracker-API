# [Assignment 2](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Assignment-2)

### Info
- Author: Sindre Eiklid (sindreik@stud.ntnu.no)
    - While the submission is individual, I have discussed the tasks with Rickard Loland and Susanne Edvardsen. We have also helped each other with problems that occurred during development ([rubber ducking](https://en.wikipedia.org/wiki/Rubber_duck_debugging) mostly).
- Root path:
    - Main: http://10.212.143.161:8080/corona/v1
    - Client: http://10.212.143.161:8081/client/v1
- I have used these REST web services to build my service:
    - https://restcountries.eu/
    - https://blog.mmediagroup.fr/post/m-media-launches-covid-19-api/
    - https://covidtracker.bsg.ox.ac.uk/about-api
- You need to be connected to NTNU network with a VPN to run the program. If you want to run it locally, you will have to change the URL variable in the 'dict' package to ```http://localhost```.
- Client Repo: https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/sindre0830/assignment-2/assignment-2-client

### Usage

There are 4 endpoints that you can append to the main root path.

1. Covid-19 Cases per Country

    - Input:
        ```
        Method: GET
        Path: .../country/{:country_name}{?scope=begin_date-end_date}
        ```
        - **{:country_name}** refers to the English name for the country as supported by https://mmediagroup.fr/covid-19.
        - **{?scope=begin_date-end_date}** indicates the scope of the data requested.

    - Output:
        ```go
        type Cases struct {
            Country              string  `json:"country"`
            Continent            string  `json:"continent"`
            Scope                string  `json:"scope"`
            Confirmed            int     `json:"confirmed"`
            Recovered            int     `json:"recovered"`
            PopulationPercentage float64 `json:"population_percentage"`
        }
        ```

    - Example:
        - Input: 
            ```
            Method: GET
            Path: http://10.212.143.161:8080/corona/v1/country/Norway?scope=2020-05-25-2021-02-01
            ```
        - Output:
            ```json
            {
                "country": "Norway",
                "continent": "Europe",
                "scope": "2020-05-25-2021-02-01",
                "confirmed": 54898,
                "recovered": 10271,
                "population_percentage": 0.01
            }
            ```

2. Covid-19 Policy Stringency Trends

    - Input:
        ```
        Method: GET
        Path: .../policy/{:country_name}{?scope=begin_date-end_date}
        ```
        - **{:country_name}** refers to the English name for the country as supported by https://mmediagroup.fr/covid-19.
        - **{?scope=begin_date-end_date}** indicates the scope of the data requested.

    - Output:
        ```go
        type Policy struct {
            Country    string  `json:"country"`
            Scope      string  `json:"scope"`
            Stringency float64 `json:"stringency"`
            Trend	   float64 `json:"trend"`
        }
        ```

    - Example:
        - Input: 
            ```
            Method: GET
            Path: http://10.212.143.161:8080/corona/v1/policy/Norway?scope=2020-05-25-2021-02-01
            ```
        - Output:
            ```json
            {
                "country": "Norway",
                "scope": "2020-05-25-2021-02-01",
                "stringency": 73.15,
                "trend": 14.82
            }
            ```

3. Diagnostics Interface

    - Input:
        ```
        Method: GET
        Path: .../diag/
        ```
        - Outputs status of each API used by the program as well as version, amount of webhooks and service uptime (in seconds).

    - Output:
        ```go
        type Diagnosis struct {
            Mmediagroupapi   int    `json:"mmediagroupapi"`
            Covidtrackerapi  int    `json:"covidtrackerapi"`
            Restcountriesapi int    `json:"restcountriesapi"`
            Registered       int    `json:"registered"`
            Version          string `json:"version"`
            Uptime           int    `json:"uptime"`
        }
        ```

    - Example:
        - Input: 
            ```
            Method: GET
            Path: http://10.212.143.161:8080/corona/v1/diag/
            ```
        - Output:
            ```json
            {
                "mmediagroupapi": 200,
                "covidtrackerapi": 200,
                "restcountriesapi": 200,
                "registered": 0,
                "version": "v1",
                "uptime": 4685
            }
            ```

4. Notification Webhook

    1. Registration of Webhook

        - Input:
            ```
            Method: POST
            Path: .../notifications/
            Body:
            {
                "url": string,
                "timeout": int64,
                "field": string,
                "country": string,
                "trigger": string
            }
            ```
            - url: Where the request will be sent to
            - timeout: Amount of time between each request. Has to be between 15 and 86400(24 hours) seconds
            - field: What data to request. Either 'Confirmed' or 'Stringency'
            - country: Country of interest
            - trigger: 'ON_TIMEOUT' to send data at each timeout. 'ON_CHANGE' to only send data if it has changed.
        - Output:
            ```go
            type Feedback struct {
                StatusCode int    `json:"status_code"`
                Message    string `json:"message"`
                ID		   string `json:"id"`
            }
            ```
        - Example:
            - Input:
                ```
                Method: POST
                Path: http://10.212.143.161:8080/corona/v1/notifications/
                Body:
                {
                    "url": "http://10.212.143.161:8081/client/v1/input/",
                    "timeout": 15,
                    "field": "confirmed",
                    "country": "denmark",
                    "trigger": "on_timeout"
                }
                ```
            - Output:
                ```json
                Status: 201
                Body:
                {
                    "status_code": 201,
                    "message": "Webhook successfully created for 'http://10.212.143.161:8081/client/v1/input/'",
                    "id": "RDlZgDM0JY6ne9wBvBRu"
                }
                ```

    2. View Registered Webhook
        - Input:
            ```
            Method: GET
            Path: .../notifications/{id}
            ```
        - Output:
            ```go
            type Feedback struct {
                StatusCode int    `json:"status_code"`
                Message    string `json:"message"`
                ID		   string `json:"id"`
            }
            ```
        - Example:
            - Input:
                ```
                Method: GET
                Path: http://10.212.143.161:8080/corona/v1/notifications/RDlZgDM0JY6ne9wBvBRu
                ```
            - Output:
                ```json
                Status: 200
                Body:
                {
                    "id": "RDlZgDM0JY6ne9wBvBRu",
                    "url": "http://10.212.143.161:8081/client/v1/input/",
                    "timeout": 15,
                    "information": "confirmed",
                    "country": "denmark",
                    "trigger": "ON_TIMEOUT"
                }
                ```

    3. Deletion of Webhook
        - Input:
            ```
            Method: DELETE
            Path: .../notifications/{:id}
            ```
        - Output:
            ```go
            type Feedback struct {
                StatusCode int    `json:"status_code"`
                Message    string `json:"message"`
                ID		   string `json:"id"`
            }
            ```
        - Example:
            - Input:
                ```
                Method: DELETE
                Path: http://10.212.143.161:8080/corona/v1/notifications/RDlZgDM0JY6ne9wBvBRu
                ``` 
            - Output:
                ```json
                Status: 200
                Body:
                {
                    "status_code": 200,
                    "message": "Webhook successfully deleted",
                    "id": "RDlZgDM0JY6ne9wBvBRu"
                }
                ```

There are 2 endpoints that you can append to the client root path.

1. Testing a webhook encryption
    - Input:
        ```
        Method: POST
        Path: .../input/
        ```
        - This is meant for webhook testing
    - Output
        - Status code and short message
    - Example:
        - See examples of usage in the notification webhook above

2. Log of webhook tests
Will show the last 5 inputs where the newest logs are first.
    - Input:
        ```
        Method: GET
        Path: .../output/
        ```
    - Output:
        ```go
        type Catch struct {
            Time         string      `json:"time"`
            ErrorMessage error       `json:"error_message"`
            RawBody      interface{} `json:"raw_body"`
        }
        ```
    - Example:
        - Input:
            ```
            Method: GET
            Path: http://10.212.143.161:8081/client/v1/output/
            ```
        - Output:
            ```json
            Status: 200
            Body:
            [
                {
                    "time": "2021-03-26 19:45:37",
                    "error_message": null,
                    "raw_body": "{\"country\":\"Denmark\",\"continent\":\"Europe\",\"scope\":\"total\",\"confirmed\":227325,\"recovered\":216287,\"population_percentage\":0.04}"
                },
                {
                    "time": "2021-03-26 19:45:22",
                    "error_message": null,
                    "raw_body": "{\"country\":\"Denmark\",\"continent\":\"Europe\",\"scope\":\"total\",\"confirmed\":227325,\"recovered\":216287,\"population_percentage\":0.04}"
                },
                {
                    "time": "2021-03-26 19:45:07",
                    "error_message": null,
                    "raw_body": "{\"country\":\"Denmark\",\"continent\":\"Europe\",\"scope\":\"total\",\"confirmed\":227325,\"recovered\":216287,\"population_percentage\":0.04}"
                },
                {
                    "time": "2021-03-26 19:44:52",
                    "error_message": null,
                    "raw_body": "{\"country\":\"Denmark\",\"continent\":\"Europe\",\"scope\":\"total\",\"confirmed\":227325,\"recovered\":216287,\"population_percentage\":0.04}"
                },
                {
                    "time": "2021-03-26 19:44:37",
                    "error_message": null,
                    "raw_body": "{\"country\":\"Denmark\",\"continent\":\"Europe\",\"scope\":\"total\",\"confirmed\":227325,\"recovered\":216287,\"population_percentage\":0.04}"
                }
            ]
            ```

## Notes

#### Design Decisions
TBA

#### Structure
TBA

#### Error Handling
TBA

#### Testing
TBA

##### Usage
TBA
