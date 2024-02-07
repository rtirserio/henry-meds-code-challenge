# Henry Meds Coding Challenge
https://henrymeds.notion.site/Reservation-Backend-v3-1e5c24f700b846f19b173f5e18c4ebc5

## Running
* Requires go (tested using version 1.21.6). Can be installed however you'd like as long as running the `go` command is available
* from root dir run `go run main.go` (this will install dependencies if they are missing and run the server)
* default server ruus on localhost:8080

## API
* `POST /providers/{id}/availability`

    Create an availability for the provider ID
    ### Request
    ```
        {
            start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
        }
    ```
    ### Response

    ```
        {
            provider_id: string -- The ID from the URL path
            start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
        }

    ```

* `GET /provider/availability`

    Gets the availability for all providers split into 15 min intervals
    ### Response

    ```
        [
            {
                provider_id: string -- The ID from the URL path
                start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
                end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            }
        ]

    ```
* `/client/{id}/appointment`

    Allows client to create an appointment

    ### Request
    ```
        {
            provider_id: string -- The ID of the provider (from the availability selected)
            start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z) (from the availability selected)
            end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z) (from the availability selected)
        }
    ```
    ### Response

    ```
        {
            id: string -- appointment ID
            provider_id: string -- ID of the provider
            client_id": string -- ID of the client from the URL path
            start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            status: string -- pending
            expires: datetime -- Time when appointment expires
        }

    ```

* `client/{{clientID}}/appointment/{{appointmentID}}/confirm`

    Allows client to confirm an appointment

    ### Request Path
    ```
        {
            clientID: string -- The ID of the client for the appointment 
            appointmentID: string -- The ID of the appointment 
        }
    ```
    ### Response

    ```
        {
            id: string -- appointment ID
            provider_id: string -- ID of the provider
            client_id": string -- ID of the client from the URL path
            start: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            end: datetime -- in the format YYYY-MM-DDTHH:MM:SSZ (ex. 2024-03-01T00:00:00Z)
            status: string -- confirmed
            expires: datetime -- Time when appointment expires
        }

    ```

## Notes / Points for discussion
* Uses an in memory database that is reset on each server restart
* Lots of functions in the "database" layer that would be better in a service layer
* There are no tests and there should be
* Confirming an appointment should validate the client ID in the URL matches the client on the appointment (ideally would handle this via auth means so only the client for that appointment can confirm)
* Errors are not mapped to the 'correct' codes and the responses are lacking in error details when they occur.