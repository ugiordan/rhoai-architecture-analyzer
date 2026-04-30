import subprocess
from flask import request, render_template

def taint_direct():
    """User input flows directly to SQL."""
    query = request.args.get("q")
    db.execute(query)

def taint_via_call():
    """User input flows through helper to subprocess."""
    cmd = request.args.get("cmd")
    run_command(cmd)

def run_command(cmd):
    subprocess.run(cmd, shell=True)

def taint_template():
    """User input flows to template render."""
    name = request.args.get("name")
    return render_template("hello.html", name=name)

def no_taint(x, y):
    """No user input, no taint."""
    return x + y
