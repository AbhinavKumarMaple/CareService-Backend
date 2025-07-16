# API Endpoints and Curl Commands for Mini EVV Logger (Caregiver Shift Tracker)

**Base URL:** `http://localhost:8080`

---

## ✅ API Endpoint: `GET /v1/schedules`

**Purpose**: Retrieve all caregiver visit schedules along with their current status, associated tasks, and interaction metadata.

### 🔸 Curl Command:

```bash
curl -X GET http://localhost:8080/v1/schedules
```

### 🔸 Expected Response (example):

```json
[
  {
    "id": "uuid",
    "client_user_id": "uuid",
    "assigned_user_id": "uuid",
    "scheduled_slot": {
      "from": "2025-07-15T09:00:00Z",
      "to": "2025-07-15T11:00:00Z"
    },
    "visit_status": "upcoming",
    "checkin_time": null,
    "checkout_time": null,
    "checkin_location": { "lat": null, "long": null },
    "checkout_location": { "lat": null, "long": null },
    "tasks": [
      {
        "id": "uuid",
        "title": "Give medication",
        "description": "Administer prescribed morning medication",
        "status": "pending",
        "done": null,
        "feedback": null
      }
    ],
    "service_note": null
  }
]
```

---

## ✅ API Endpoint: `GET /v1/schedules/today`

**Purpose**: Retrieve today's caregiver schedules.

### 🔸 Curl Command:

```bash
curl -X GET "http://localhost:8080/v1/schedules/today?client_user_id=YOUR_CLIENT_USER_ID_HERE"
```

### 🔸 Expected Response:

Same structure as `GET /v1/schedules`, filtered by the current date.

---

## ✅ API Endpoint: `GET /v1/schedules/:id`

**Purpose**: Retrieve detailed information for a specific schedule, including associated tasks and visit status.

### 🔸 Curl Command:

```bash
curl -X GET http://localhost:8080/v1/schedules/YOUR_SCHEDULE_ID_HERE
```

### 🔸 Expected Response (example):

```json
{
  "id": "uuid",
  "client_user_id": "uuid",
  "scheduled_slot": {
    "from": "2025-07-15T08:30:00Z",
    "to": "2025-07-15T10:30:00Z"
  },
  "visit_status": "completed",
  "checkin_time": "2025-07-15T08:35:00Z",
  "checkout_time": "2025-07-15T10:25:00Z",
  "checkin_location": {
    "lat": 28.6135,
    "long": 77.2091
  },
  "checkout_location": {
    "lat": 28.6137,
    "long": 77.2089
  },
  "tasks": [
    {
      "id": "uuid",
      "title": "Assist with bathing",
      "description": "Help patient with morning bath and grooming",
      "status": "not_completed",
      "done": false,
      "feedback": "Patient refused bathing"
    }
  ],
  "service_note": "Client was calm and cooperative."
}
```

---

## ✅ API Endpoint: `POST /v1/schedules`

**Purpose**: Create a new caregiver visit schedule with associated tasks.

### 🔸 Request Body:

```json
{
  "client_user_id": "uuid",
  "assigned_user_id": "uuid",
  "scheduled_slot": {
    "from": "2025-07-15T09:00:00Z",
    "to": "2025-07-15T11:00:00Z"
  },
  "tasks": [
    {
      "title": "Give medication",
      "description": "Administer prescribed morning medication"
    },
    {
      "title": "Prepare meal",
      "description": "Prepare lunch according to dietary restrictions"
    }
  ]
}
```

### 🔸 Curl Command:

```bash
curl -X POST http://localhost:8080/v1/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "client_user_id": "YOUR_CLIENT_USER_ID_HERE",
    "assigned_user_id": "YOUR_ASSIGNED_USER_ID_HERE",
    "scheduled_slot": {
      "from": "2025-07-15T09:00:00Z",
      "to": "2025-07-15T11:00:00Z"
    },
    "tasks": [
      {
        "title": "Give medication",
        "description": "Administer prescribed morning medication"
      },
      {
        "title": "Prepare meal",
        "description": "Prepare lunch according to dietary restrictions"
      }
    ]
  }'
```

### 🔸 Expected Response:

```json
{
  "id": "uuid",
  "client_user_id": "uuid",
  "assigned_user_id": "uuid",
  "scheduled_slot": {
    "from": "2025-07-15T09:00:00Z",
    "to": "2025-07-15T11:00:00Z"
  },
  "visit_status": "upcoming",
  "checkin_time": null,
  "checkout_time": null,
  "checkin_location": { "lat": null, "long": null },
  "checkout_location": { "lat": null, "long": null },
  "tasks": [
    {
      "id": "uuid",
      "title": "Give medication",
      "description": "Administer prescribed morning medication",
      "status": "pending",
      "done": null,
      "feedback": null
    },
    {
      "id": "uuid",
      "title": "Prepare meal",
      "description": "Prepare lunch according to dietary restrictions",
      "status": "pending",
      "done": null,
      "feedback": null
    }
  ],
  "service_note": null
}
```

---

## ✅ API Endpoint: `POST /v1/schedules/:id/start`

**Purpose**: Register the caregiver's check-in time and geolocation.

### 🔸 Request Body:

```json
{
  "timestamp": "2025-07-15T08:35:00Z",
  "location": {
    "lat": 28.6135,
    "long": 77.2091
  }
}
```

### 🔸 Curl Command:

```bash
curl -X POST http://localhost:8080/v1/schedules/YOUR_SCHEDULE_ID_HERE/start \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2025-07-15T08:35:00Z",
    "location": {
      "lat": 28.6135,
      "long": 77.2091
    }
  }'
```

### 🔸 Expected Response:

```json
{
  "message": "Check-in recorded successfully",
  "checkin_time": "2025-07-15T08:35:00Z",
  "checkin_location": {
    "lat": 28.6135,
    "long": 77.2091
  }
}
```

---

## ✅ API Endpoint: `POST /v1/schedules/:id/end`

**Purpose**: Register the caregiver's check-out time, geolocation, and visit summary note.

### 🔸 Request Body:

```json
{
  "timestamp": "2025-07-15T10:25:00Z",
  "location": {
    "lat": 28.6137,
    "long": 77.2089
  },
  "service_note": "All activities attempted. Some refused by client."
}
```

### 🔸 Curl Command:

```bash
curl -X POST http://localhost:8080/v1/schedules/YOUR_SCHEDULE_ID_HERE/end \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2025-07-15T10:25:00Z",
    "location": {
      "lat": 28.6137,
      "long": 77.2089
    },
    "service_note": "All activities attempted. Some refused by client."
  }'
```

### 🔸 Expected Response:

```json
{
  "message": "Check-out recorded successfully",
  "checkout_time": "2025-07-15T10:25:00Z",
  "checkout_location": {
    "lat": 28.6137,
    "long": 77.2089
  },
  "service_note": "All activities attempted. Some refused by client."
}
```

---

## ✅ API Endpoint: `POST /v1/tasks/:taskId/update`

**Purpose**: Update the status of a specific care task.

### 🔸 Request Body:

```json
{
  "status": "not_completed",
  "done": false,
  "feedback": "Patient was asleep"
}
```

### 🔸 Curl Command:

```bash
curl -X POST http://localhost:8080/v1/tasks/YOUR_TASK_ID_HERE/update \
  -H "Content-Type: application/json" \
  -d '{
    "status": "not_completed",
    "done": false,
    "feedback": "Patient was asleep"
  }'
```

### 🔸 Expected Response:

```json
{
  "message": "Task updated successfully",
  "task": {
    "id": "uuid",
    "status": "not_completed",
    "done": false,
    "feedback": "Patient was asleep"
  }
}
```

---

## ✅ User Management Endpoints (Authentication Disabled)

These endpoints no longer require an `Authorization` header.

1.  **Create User**

    - **Endpoint:** `POST /v1/user`
    - **Description:** Creates a new user.
    - **Curl Command:**
      ```bash
      curl -X POST http://localhost:8080/v1/user \
        -H "Content-Type: application/json" \
        -d '{
          "user_name": "newuser",
          "email": "newuser@example.com",
          "first_name": "New",
          "last_name": "User",
          "role": "caregiver",
          "location": {
            "house_number": "101",
            "street": "New Street",
            "city": "New City",
            "state": "NS",
            "pincode": "98765",
            "lat": 34.5678,
            "long": -119.8765
          }
        }'
      ```

2.  **Get All Users**

    - **Endpoint:** `GET /v1/user`
    - **Description:** Retrieves a list of all users.
    - **Curl Command:**
      ```bash
      curl -X GET http://localhost:8080/v1/user
      ```

3.  **Get User by ID**

    - **Endpoint:** `GET /v1/user/:id`
    - **Description:** Retrieves a single user by their ID.
    - **Curl Command:**
      ```bash
      curl -X GET http://localhost:8080/v1/user/YOUR_USER_ID_HERE
      ```

4.  **Update User**

    - **Endpoint:** `PUT /v1/user/:id`
    - **Description:** Updates an existing user's information.
    - **Curl Command:**
      ```bash
      curl -X PUT http://localhost:8080/v1/user/YOUR_USER_ID_HERE \
        -H "Content-Type: application/json" \
        -d '{
          "user_name": "updateduser",
          "email": "updateduser@example.com",
          "first_name": "Updated",
          "last_name": "User",
          "role": "client",
          "location": {
            "house_number": "202",
            "street": "Updated Ave",
            "city": "Updated City",
            "state": "UP",
            "pincode": "54321",
            "lat": 35.1234,
            "long": -120.5678
          }
        }'
      ```

5.  **Delete User**

    - **Endpoint:** `DELETE /v1/user/:id`
    - **Description:** Deletes a user by their ID.
    - **Curl Command:**
      ```bash
      curl -X DELETE http://localhost:8080/v1/user/YOUR_USER_ID_HERE
      ```

6.  **Search Users with Pagination**

    - **Endpoint:** `GET /v1/user/search`
    - **Description:** Searches for users with pagination.
    - **Curl Command:**
      ```bash
      curl -X GET "http://localhost:8080/v1/user/search?page=1&limit=10&email_like=example.com"
      ```

7.  **Search Users by Specific Property**
    - **Endpoint:** `GET /v1/user/search-property`
    - **Description:** Searches for users by a specific property.
    - **Curl Command:**
      ```bash
      curl -X GET "http://localhost:8080/v1/user/search-property?property=email&searchText=newuser@example.com"
      ```

---

## ✅ Medicine Management Endpoints (Authentication Disabled)

These endpoints no longer require an `Authorization` header.

1.  **Create Medicine**

    - **Endpoint:** `POST /v1/medicine`
    - **Description:** Creates a new medicine entry.
    - **Curl Command:**
      ```bash
      curl -X POST http://localhost:8080/v1/medicine \
        -H "Content-Type: application/json" \
        -d '{
          "name": "Aspirin",
          "description": "Pain reliever",
          "ean_code": "1234567890123",
          "laboratory": "Bayer"
        }'
      ```

2.  **Get All Medicines**

    - **Endpoint:** `GET /v1/medicine`
    - **Description:** Retrieves a list of all medicines.
    - **Curl Command:**
      ```bash
      curl -X GET http://localhost:8080/v1/medicine
      ```

3.  **Get Medicine by ID**

    - **Endpoint:** `GET /v1/medicine/:id`
    - **Description:** Retrieves a single medicine by its ID.
    - **Curl Command:**
      ```bash
      curl -X GET http://localhost:8080/v1/medicine/YOUR_MEDICINE_ID_HERE
      ```

4.  **Update Medicine**

    - **Endpoint:** `PUT /v1/medicine/:id`
    - **Description:** Updates an existing medicine's information.
    - **Curl Command:**
      ```bash
      curl -X PUT http://localhost:8080/v1/medicine/YOUR_MEDICINE_ID_HERE \
        -H "Content-Type: application/json" \
        -d '{
          "name": "Ibuprofen",
          "description": "Anti-inflammatory",
          "ean_code": "9876543210987",
          "laboratory": "Pfizer"
        }'
      ```

5.  **Delete Medicine**

    - **Endpoint:** `DELETE /v1/medicine/:id`
    - **Description:** Deletes a medicine by its ID.
    - **Curl Command:**
      ```bash
      curl -X DELETE http://localhost:8080/v1/medicine/YOUR_MEDICINE_ID_HERE
      ```

6.  **Search Medicines with Pagination**

    - **Endpoint:** `GET /v1/medicine/search`
    - **Description:** Searches for medicines with pagination.
    - **Curl Command:**
      ```bash
      curl -X GET "http://localhost:8080/v1/medicine/search?page=1&limit=10&name_like=aspirin"
      ```

7.  **Search Medicines by Specific Property**
    - **Endpoint:** `GET /v1/medicine/search-property`
    - **Description:** Searches for medicines by a specific property.
    - **Curl Command:**
      ```bash
      curl -X GET "http://localhost:8080/v1/medicine/search-property?property=name&searchText=aspirin"
      ```

---

## ✅ Health Check Endpoint

This endpoint does not require authentication.

1.  **Health Check**
    - **Endpoint:** `GET /v1/health`
    - **Description:** Checks the health status of the microservice.
    - **Curl Command:**
      ```bash
      curl -X GET http://localhost:8080/v1/health
      ```
    - **Expected Response (example):**
      ```json
      {
        "status": "ok",
        "message": "Service is running"
      }
      ```
