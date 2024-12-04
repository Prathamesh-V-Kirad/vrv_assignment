# Task Management System with RBAC

A secure task management application demonstrating **Role-Based Access Control (RBAC)**, authentication, and authorization using **React.js** and the **Go Fiber** framework.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Database Model](#database-model)
- [Project Setup](#project-setup)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
- [Security Implementation](#security-implementation)
- [Testing](#testing)
- [Future Scope](#future-scope)
- [Important Note](#important-note)

## Features

### Core Features
- Complete user authentication flow
- Role-based access control (Admin, Manager, User)
- Task management with CRUD operations
- Protected routes and API endpoints
- Activity logging and monitoring
- Real-time validations

### Security Features
- JWT-based authentication 
- Password hashing with bcrypt
- Rate limiting
- CORS protection
- Secure headers implementation
- Input validation and sanitization
- Role-based middleware checks

## Tech Stack

### Backend
- Go 
- Fiber web framework
- MongoDB
- JWT for authentication
- bcrypt for password hashing
- CORS middleware

### Frontend
- React 
- React Router 
- Axios
- ShadCN
- JWT token management
- Protected routes

## Database Model
The Below model is built keeping in mind. So we can add additonal roles , permission and change name of the roles.
This is a modular setup which means a Admin dashboard can be implemented easily without need to manually can the backend role authorization functions.
By only introducing certain apis for changing roles and permissions in DB.

![Database_model](https://github.com/user-attachments/assets/d5b050f6-920c-49ce-ba2f-4c4268d2350c)


## Project Setup

### Prerequisites

Before setting up the project, ensure you have the following installed:

- **Go** v1.21 or higher
- **Node.js** v16 or higher
- **MongoDB** v6.0 or higher (for local setup) or a **MongoDB Atlas** account
- **Git**

### Backend Setup

1. **Clone Repository**

   Clone the repository to your local machine and navigate to the backend folder:

   ```bash
   git clone https://github.com/Prathamesh-V-Kirad/vrv_assignment
   cd vrv_assignment/backend
  
2. **Install Dependencies**

   Run the following command to install Go dependencies:

   ```bash
   go mod tidy
   ```

3. **Stepup env**

   ```bash
   PORT = 8000
   MONGODB_URI = MONGODB_URI
   JWT_SECRET_KEY = JWT_SECRET_KEY
   ```

4. **Start the server**

   ```bash
   go run main.go
   ```

### Frontend Setup

1. **Clone Repository**

   Clone the repository to your local machine and navigate to the frontend folder:
  
    ```bash
    git clone https://github.com/Prathamesh-V-Kirad/vrv_assignment
    cd vrv_assignment/frontend
    ```

2. **Install Dependencies**

    Install the required dependencies for the frontend:

    ```bash
    Copy code
    npm install
    ```
    
3. **Run the Frontend**

    Start the frontend development server:

    ```bash
    Copy code
    npm start
    ```
    The frontend will be available on http://localhost:5173.


### Security Implementation
1. **JWT Authentication**
The backend uses JWT for authentication. Tokens are generated after a successful login and are required for accessing protected routes.
The JWT is stored in a secure HTTP-only cookie.
2. **Role-Based Access Control (RBAC)**
The backend enforces role-based access control (Admin, Manager, User).
Each role has specific permissions to create, update, and delete tasks.
3. **Password Hashing**
Passwords are securely hashed using bcrypt before being stored in the database.
4. **JWT Algorithm**
Used HSA256 Asymmetric Algorithm for more safety.
5 **Middlewares**
Used Multiple Middles to implement CSP , CSRF , etc header policies
csrf - Disables keep-alive while helmet-sets CSP (Content Security Policy Headers) and XSS.
Eg:
    ```bash
    app.Use(csrf.New())
    app.Use(helmet.New())
    ```

### Testing 

**Important Note**
Here We have intialized DB already with roles and permissions. Due to limited time a Admin Dashboard was not created, 
Hence , the first registered user is treated as admin and the rest of the registrations are treated as users.
The code is modular. Where we have permissions , roles which can be easily used to make a admin dashboard .

https://github.com/user-attachments/assets/6b3725ee-32a1-49c1-9e56-c4cb21145df3



### Future Scope
Give below are the ways to make app more secure.
- Adding OTP verification for registrations
- Adding Refresh Token Functionality
- Adding encrpytion to payload being sent (Currently plain/text)
- Implementing Strict Policies.
- Making code more Modular separating Middlewares.
- Adding Schema Validation (Zod) to Server Side Also for extra security.

### Important Note
- HTTPS and Certain Headers are commented because of localhost testing. They will be implemented while hosting.
- OTP verification/ Verification of mail is an important feature while registration of email .
  It can be implemented but due to localhost constraints no Twilio , etc are setup to ease up the process
    
