# archGap-server

Server side of my own chat app.

## Description

`archGap-server` is the backend server for the `archGap` chat application. This server is built using Go and handles all server-side operations, including user authentication, message handling, and real-time communication.

## Features

- User authentication and authorization
- Real-time messaging
- Group chat support
- Message history
- User presence and status updates

## Requirements

- Go 1.19 or later

## Installation

1. Clone the repository:

```bash
git clone https://github.com/archroid/archGap-server.git
```

2. Navigate to the project directory:

```bash
cd archGap-server
```

3. Install dependencies:

```bash
go mod download
```

## Configuration

Create a `.env` file in the root directory and configure the following environment variables:

```env
DB_HOST=your_database_host
DB_PORT=your_database_port
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name

JWT_SECRET=your_jwt_secret_key
```

## Running the Server

Start the server with the following command:

```bash
go run main.go
```

The server will start running on the configured port. You can change the port by setting the `PORT` environment variable in the `.env` file.

## API Endpoints

### Authentication

- `POST /register` - Register a new user
- `POST /login` - Authenticate a user and return a JWT

### Messages

- `POST /messages` - Send a new message
- `GET /messages` - Retrieve message history

### Users

- `GET /users` - Get a list of users
- `GET /users/:id` - Get details of a specific user

### WebSocket

- Connect to `/ws` for real-time messaging

## Contributing

Feel free to submit issues and enhancement requests. Please follow the [contributing guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or inquiries, please reach out to [archroid](https://github.com/archroid).