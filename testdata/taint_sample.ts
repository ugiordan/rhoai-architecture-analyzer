import express from 'express';

const app = express();

// User input flows to response (XSS)
app.get('/greet', (req, res) => {
    const name = req.query.name;
    res.send(`<h1>Hello ${name}</h1>`);
});

// User input flows through helper to eval
app.get('/calc', (req, res) => {
    const expr = req.query.expr;
    const result = evaluate(expr);
    res.json({ result });
});

function evaluate(expression: string): any {
    return eval(expression);
}

// No taint: pure computation
function add(a: number, b: number): number {
    return a + b;
}
