# Order Notification System

This project implements an order notification system using Golang. It consists of a REST API for order creation and a WebSocket server for real-time notifications to the kitchen/admin side when new orders are placed.

## Project Structure

```
order-notification-system
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── api
│   │   └── order.go     # REST API for handling orders
│   ├── websocket
│   │   └── handler.go   # WebSocket connection management
│   ├── models
│   │   └── order.go     # Order model definition
│   └── utils
│       └── notifier.go   # Utility functions for notifications
├── go.mod                # Module dependencies
├── go.sum                # Module dependency checksums
└── README.md             # Project documentation
```
  
## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd order-notification-system
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

## Usage

### Creating an Order

To create a new order, send a POST request to the `/order` endpoint with the following JSON body:

```json
{
  "item": "Pizza",
  "quantity": 2
}
```

### WebSocket Notifications

The kitchen/admin can connect to the WebSocket server to receive real-time notifications about new orders. The WebSocket server will broadcast notifications whenever a new order is created.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.