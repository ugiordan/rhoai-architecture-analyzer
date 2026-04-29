// testdata/complexity_sample.rs

// Complexity 1
fn simple_func() -> i32 {
    42
}

// Complexity 4: if + else if + for + &&
fn complex_func(x: i32, items: &[String]) -> i32 {
    if x > 0 && x < 100 {
        return 1;
    } else if x >= 100 {
        return 2;
    }
    for item in items {
        println!("{}", item);
    }
    0
}

// Complexity 5: match with 4 arms
fn match_func(op: &str) -> i32 {
    match op {
        "add" => 1,
        "sub" => 2,
        "mul" => 3,
        _ => 0,
    }
}

// Complexity 4: while + loop + if + ||
fn loop_func(mut n: i32) -> i32 {
    while n > 0 {
        n -= 1;
    }
    loop {
        n += 1;
        if n > 10 || n < -10 {
            break;
        }
    }
    n
}
