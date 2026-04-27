# testdata/complexity_sample.py

# Complexity 1
def simple_func():
    return 42

# Complexity 6: if + and + elif + for + if + base (or not counted)
def complex_func(x, items):
    if x > 0 and x < 100:
        return 1
    elif x >= 100:
        return 2
    for item in items:
        if item or x:
            pass
    return 0

# Complexity 4: while + except + for + base
def loop_func(n):
    while n > 0:
        n -= 1
    try:
        result = 1 / n
    except ZeroDivisionError:
        result = 0
    for i in range(n):
        pass
    return result

# Complexity 3: if + comprehension-if + base
def comprehension_func(items):
    if items:
        return [x for x in items if x > 0]
    return []
