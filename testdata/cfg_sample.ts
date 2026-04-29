function ifElse(x: number): string {
    if (x > 0) {
        return "positive";
    } else {
        return "non-positive";
    }
}

function earlyReturn(data: string | null): string | null {
    if (!data) {
        return null;
    }
    const result = data.trim();
    return result;
}

function forLoop(items: string[]): string {
    let result = "";
    for (const item of items) {
        result = item;
    }
    return result;
}

function switchCase(op: string): number {
    let result: number;
    switch (op) {
        case "add":
            result = 1;
            break;
        case "sub":
            result = 2;
            break;
        default:
            result = 0;
    }
    return result;
}

function nestedIfInFor(items: string[]): number {
    let count = 0;
    for (const item of items) {
        if (item.length > 3) {
            count++;
        }
    }
    return count;
}

function linearFunction(): number {
    const x = 1;
    const y = x + 2;
    return y;
}

function emptyFunction(): void {
}
