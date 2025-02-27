# Overview
This service is an intermediary between multiple systems for client and lab result data exchange. The service integrates:

- **eCHIS:** A system that sends client data to eCBSS (a DHIS2-based system) through this service.

- **LabXpert:** A system that submits lab results, which are forwarded to eCBSS.

The service ensures secure data transmission and validation before forwarding information to eCBSS.
## Authentication

This API supports two authentication methods:

- **Basic Authentication** (Username & Password)
- **Token Authentication** (Token in the Authorization header: Token: \<user-token\>)

Clients must provide valid credentials in the request headers to access the endpoints.

## Endpoints

### 1. Get User Token

**Endpoint:** `POST /api/getToken`

**Description:** This endpoint retrieves a user's authentication token.

**Request Headers:**

```
Authorization: Basic <base64-encoded-credentials>
Content-Type: application/json
```

**Request Body:**

```json
{
  "username": "user@example.com",
  "password": "securepassword"
}
```

**Response:**

```json
{
  "token": "your-authentication-token"
}
```

---

### 2. Generate New Token

**Endpoint:** `POST /api/generateToken`

**Description:** This endpoint generates a new authentication token.

**Request Headers:**

```
Authorization: Basic <base64-encoded-credentials>
Content-Type: application/json
```

**Response:**

```json
{
  "message": "New token generated successfully",
  "token": "newly-generated-authentication-token"
}
```

---


### 3. eCHIS integration with eCBSS

**Endpoint:** `POST /api/clients`

**Description:** This endpoint registers a new client from the eCHIS system into eCBSS.

**Request Headers:**

```
Authorization: Basic <base64-encoded-credentials>  
OR  
Authorization: Bearer <your-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "echis_patient_id": "1234567890", // Mandatory
  "national_identification_number": "",
  "name": "",
  "patient_phone": "",
  "patient_age_in_years": "",
  "patient_age_in_months": "",
  "patient_age_in_days": "",
  "patient_gender": "",
  "client_category": "",
  "facility_id": "",
  "facility_dhis2_id": "", // Mandatory
  "patient_category": "",
  "cough": "",
  "fever": "",
  "weight_loss": "",
  "excessive_night_sweat": "",
  "is_on_tb_treatment": "",
  "poor_weight_gain": ""
}
```

**Response:**

```json
{
  "message": "client queued for saving to DHIS2"
}
```

---

### 4. LabXpert integration with eCBSS
Here the service submits results from the LabXpert system to the eCBSS system.

**Endpoint:** `POST /api/results`

**Description:** This endpoint submits test results for a patient.

**Request Headers:**

```
Authorization: Basic <base64-encoded-credentials>  
OR  
Authorization: Bearer <your-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "patient_id": "1234567890",
  "mtb": "DETECTED HIGH",
  "rr": "DETECTED",
  "results_data": "2025-01-27 13:08:27",
  "facility_id": "FvewOonC8lS",
  "dhis2_id": ""
}
```

**Response:**

```json
{
  "message": "results queued for saving to DHIS2"
}
```

## Error Responses

For all endpoints, the API returns standard HTTP status codes. Below are common responses:

| Status Code               | Description                            |
| ------------------------- | -------------------------------------- |
| 200 OK                    | Request successful                     |
| 201 Created               | Resource successfully created          |
| 400 Bad Request           | Invalid request payload                |
| 401 Unauthorized          | Invalid authentication credentials     |
| 403 Forbidden             | Insufficient permissions               |
| 500 Internal Server Error | Server encountered an unexpected error |

## Notes

- Ensure that authentication credentials are valid before making requests.
- Required fields must be provided to avoid errors.
- The response format is JSON.

