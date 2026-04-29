// testdata/complexity_sample.ts

// Complexity 1
function simpleFunc(): number {
  return 42;
}

// Complexity 5: if + else if + for + && + ternary
function complexFunc(x: number, items: string[]): number {
  if (x > 0 && x < 100) {
    return 1;
  } else if (x >= 100) {
    return 2;
  }
  for (const item of items) {
    console.log(item);
  }
  return x > 0 ? 1 : 0;
}

// Complexity 5: switch with 4 cases
function switchFunc(op: string): number {
  switch (op) {
    case "add":
      return 1;
    case "sub":
      return 2;
    case "mul":
      return 3;
    default:
      return 0;
  }
}

// Complexity 4: while + do + catch + base (|| not counted)
function loopFunc(n: number): number {
  while (n > 0) {
    n--;
  }
  do {
    n++;
  } while (n < 0);
  try {
    return 1 / n;
  } catch (e) {
    return n || 0;
  }
}
