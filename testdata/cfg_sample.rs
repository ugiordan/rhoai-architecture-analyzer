fn if_else(x: i32) -> &'static str {
    if x > 0 {
        return "positive";
    } else {
        return "non-positive";
    }
}

fn early_return(data: Option<String>) -> Option<String> {
    if data.is_none() {
        return None;
    }
    let result = data.unwrap();
    Some(result)
}

fn for_loop(items: Vec<String>) -> String {
    let mut result = String::new();
    for item in items {
        result = item;
    }
    result
}

fn match_case(op: &str) -> i32 {
    let result = match op {
        "add" => 1,
        "sub" => 2,
        _ => 0,
    };
    result
}

fn nested_if_in_for(items: Vec<String>) -> i32 {
    let mut count = 0;
    for item in items {
        if item.len() > 3 {
            count += 1;
        }
    }
    count
}

fn linear_function() -> i32 {
    let x = 1;
    let y = x + 2;
    y
}

fn empty_function() {
}
