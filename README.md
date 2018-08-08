purchase-api
=====

> purchase-api is a simple example of rest api built using Go.

It receive some purchase request, call a external service, store data and return the external content to original client

## Prerequisites

purchase-api needs the Postgres database.

## API Reference

> At bellow are the exposed services:
> Create new purchase

URI : /purchase
Method: POST
Payload example:

    Header:
        Bearer 12345
    Body:
        {
            “amount”: 20.32
        }

Response:
    {
        voucher-code:FHGCHEEH12343
    }

Response codes:
- 201 - If purchase was created
- 400 - If problems was found on request payload
- 403 - If the head parameter Bearer was not found

> Get all purchase stored

URI : /purchase
Method: GET
Payload example:

    Header:
        Bearer 12345

Response:
    [
        {
            "id": "1",
            "external_id": "6af5fa8f-e3e9-46d5-ad3c-239a36f4a395",
            "voucher-code": "FHGCHEEH12343",
            "amount": 435.02
        }
    ]

Response codes:
- 201 - If purchase was created
- 400 - If problems was found on request payload
- 403 - If the head parameter Bearer was not found
