# error 模块

## error

`src/error.rs`

```rust
use std::fmt;

use axum::{
    Json,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    pub status: String,
    pub message: String,
}

impl fmt::Display for ErrorResponse {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", serde_json::to_string(&self).unwrap())
    }
}

#[derive(Debug, PartialEq)]
pub enum ErrorMessage {
    EmpltyPassword,
    ExceededMaxPasswordLength(usize),
    PasswordTooShort(usize),
    InvalidHashFormat,
    HashingError,
    InvaildToken,
    ServerError,
    WrongCredentials,
    EmailExit,
    UserNoLongExist,
    TokenNotProvided,
    PermissionDenied,
}

impl ToString for ErrorMessage {
    fn to_string(&self) -> String {
        self.to_str().to_owned()
    }
}

impl ErrorMessage {
    fn to_str(&self) -> String {
        match self {
            ErrorMessage::ServerError => "Server Error, Please Try Again later".to_string(),
            ErrorMessage::WrongCredentials => "Email or password is incorrect".to_string(),
            ErrorMessage::EmailExit => "A user with this email already exists".to_string(),
            ErrorMessage::UserNoLongExist => {
                "User belonging to this token no longer exists".to_string()
            }
            ErrorMessage::EmpltyPassword => "Password cannot be empty".to_string(),
            ErrorMessage::HashingError => "Error while hashing password".to_string(),
            ErrorMessage::InvalidHashFormat => "Invalid password hash format".to_string(),
            ErrorMessage::ExceededMaxPasswordLength(max_length) => {
                format!("Password must not be more than {} characters", max_length)
            }
            ErrorMessage::PasswordTooShort(min_length) => {
                format!("Password must be at least {} characters", min_length)
            }
            ErrorMessage::InvaildToken => "Authentication token is invalid or expired".to_string(),
            ErrorMessage::TokenNotProvided => {
                "You are not logged in, please log in to continue".to_string()
            }
            ErrorMessage::PermissionDenied => {
                "You are not authorized to perform this action".to_string()
            }
        }
    }
}

#[derive(Debug, Clone)]
pub struct HttpError {
    pub message: String,
    pub status: StatusCode,
}

impl HttpError {
    pub fn _new(message: impl Into<String>, status: StatusCode) -> Self {
        Self {
            message: message.into(),
            status,
        }
    }

    pub fn server_error(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
            status: StatusCode::INTERNAL_SERVER_ERROR,
        }
    }

    pub fn bad_request(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
            status: StatusCode::BAD_REQUEST,
        }
    }

    pub fn unauthorized(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
            status: StatusCode::UNAUTHORIZED,
        }
    }

    pub fn unique_constraint_violation(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
            status: StatusCode::CONFLICT,
        }
    }

    pub fn into_http_response(self) -> Response {
        let json_response = Json(ErrorResponse {
            message: self.message,
            status: "fail".to_string(),
        });
        (self.status, json_response).into_response()
    }
}

impl fmt::Display for HttpError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(
            f,
            "HttpError: message: {}, status: {}",
            self.message, self.status
        )
    }
}

impl std::error::Error for HttpError {}

impl IntoResponse for HttpError {
    fn into_response(self) -> Response {
        self.into_http_response()
    }
}
```

## main 入口

`src/main.rs`

```rust
mod constants;
mod models;
mod config;
mod error; // [!code ++]

fn main() {
    println!("Hello, world!");
}
```