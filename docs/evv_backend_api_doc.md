# 📘 API Endpoint & Structure Specification: Mini EVV Logger (Caregiver Shift Tracker)

This document outlines all backend API endpoints with precise request and response body structures, based strictly on the provided case study and updated requirements.

---

## ✅ API Endpoint: `GET /schedules`

**Purpose**: Retrieve all caregiver visit schedules along with their current status, associated tasks, and interaction metadata.

### 🔸 Response:

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

## ✅ API Endpoint: `GET /schedules/today`

**Purpose**: Retrieve today's caregiver schedules.

### 🔸 Response:

Same structure as `GET /schedules`, filtered by the current date.

---

## ✅ API Endpoint: `GET /schedules/:id`

**Purpose**: Retrieve detailed information for a specific schedule, including associated tasks and visit status.

### 🔸 Response:

```json
{
  "id": "uuid",
  "client_user_id": "uuid",
  "assigned_user_id": "uuid",
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

## ✅ API Endpoint: `POST /schedules`

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

### 🔸 Response:

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

## ✅ API Endpoint: `POST /schedules/:id/start`

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

### 🔸 Response:

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

## ✅ API Endpoint: `POST /schedules/:id/end`

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

### 🔸 Response:

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

## ✅ API Endpoint: `POST /tasks/:taskId/update`

**Purpose**: Update the status of a specific care task.

### 🔸 Request Body:

```json
{
  "status": "not_completed",
  "done": false,
  "feedback": "Patient was asleep"
}
```

### 🔸 Response:

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

## ✅ Updated `User` Table Schema

```go
type User struct {
  ID           uuid.UUID `gorm:"primaryKey"`
  UserName     string    `gorm:"column:user_name;unique"`
  Email        string    `gorm:"unique"`
  FirstName    string    `gorm:"column:first_name"`
  LastName     string    `gorm:"column:last_name"`
  Status       bool      `gorm:"column:status"`
  HashPassword string    `gorm:"column:hash_password"`
  Role         string    `gorm:"column:role"` // caregiver, admin, client
  Location     struct {
    HouseNumber string  `json:"house_number"`
    Street      string  `json:"street"`
    City        string  `json:"city"`
    State       string  `json:"state"`
    Pincode     string  `json:"pincode"`
    Lat         float64 `json:"lat"`
    Long        float64 `json:"long"`
  } `gorm:"embedded;embeddedPrefix:location_"`
  CreatedAt    time.Time `gorm:"autoCreateTime:milli"`
  UpdatedAt    time.Time `gorm:"autoUpdateTime:milli"`
}
```

> 📝 The `User` table contains location information as a structured object. The schedule now references only `client_user_id` and avoids duplicating user information.

---

## 🧪 Standard Error Response Format

```json
{
  "error": "Validation failed",
  "details": [
    "timestamp is required",
    "location.lat must be between -90 and 90",
    "feedback is required when done is false"
  ]
}
```

---

> 🧾 This document reflects the latest structural and naming conventions, including: `scheduled_slot`, `checkin`/`checkout`, embedded task feedback, geolocation tracking, and user reference integration.
