use std::process::Command;
use actix_web::{web, HttpResponse};

// User input flows to command execution
async fn taint_command(query: web::Query<std::collections::HashMap<String, String>>) -> HttpResponse {
    let cmd = query.get("cmd").unwrap();
    let output = Command::new(cmd).output().unwrap();
    HttpResponse::Ok().body(format!("{:?}", output))
}

// No taint: pure computation
fn compute(x: i32, y: i32) -> i32 {
    x + y
}

// User input to file access
async fn taint_file(query: web::Query<std::collections::HashMap<String, String>>) -> HttpResponse {
    let path = query.get("path").unwrap();
    let content = std::fs::read_to_string(path).unwrap();
    HttpResponse::Ok().body(content)
}
