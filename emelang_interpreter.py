#!/usr/bin/env python3
"""
Emerald Language (emelang/emlg) interpreter.

Supported statements:
  var <name> <value>          / declare a variable with a literal value
  print(<expr>)               / print a value (literal or previously declared var)

Comments:
  // this is a single line comment

Literals:
  - strings:  "hello"  or  'world'
  - booleans: True  False
  - numbers:  123  45.67

Variables are stored in a simple global dictionary and must be accessed
via the 'var.' prefix (e.g., var.my_variable).

Usage:
  emelang script.emlg
"""

import re
import sys

# ----------------------------------------------------------------------
# Helpers: evaluation and expressions
# ----------------------------------------------------------------------
def _eval_token(tok: str, env: dict):
    tok = tok.strip()
    # variable reference
    if tok.startswith("var."):
        ref_name = tok[4:]
        if ref_name in env:
            return env[ref_name]
        else:
            raise NameError(f"Variable '{ref_name}' not defined")

    # booleans
    if tok.lower() == "true":
        return True
    if tok.lower() == "false":
        return False

    # numbers (int or float)
    if re.fullmatch(r"-?\d+(\.\d+)?", tok):
        return float(tok) if "." in tok else int(tok)

    # quoted strings
    if tok.startswith('"') and tok.endswith('"'):
        return tok[1:-1]
    if tok.startswith("'") and tok.endswith("'"):
        return tok[1:-1]

    return None # Not a simple token

def _evaluate_expression(expr: str, env: dict):
    # Support: math (+, -, *, /) and comparisons (=, !=, <, >, =<, =>)
    # 1. Substitute variables with their literal values for easier parsing
    def _sub_var(m):
        val = _eval_token(m.group(0), env)
        return str(val) if val is not None else m.group(0)

    # Simple binary operation support
    op_match = re.search(r"\s*(==|!=|<=|>=|=<|=>|=|<|>|\+|\-|\*|\/)\s*", expr)
    if not op_match:
        res = _eval_token(expr, env)
        if res is None: raise SyntaxError(f"Invalid expression: {expr}")
        return res
        
    # If we have operators, let's split and evaluate
    # Note: This simple evaluator doesn't handle precedence yet, just left-to-right
    # For a real language, we'd want a proper parser.
    
    # We'll use Python's eval safely by replacing var.name with env[name]
    safe_expr = expr
    for var_name in sorted(env.keys(), key=len, reverse=True):
        placeholder = f"var.{var_name}"
        if placeholder in safe_expr:
            val = env[var_name]
            # Replace with a unique string that we can replace with the actual value safely
            safe_expr = safe_expr.replace(placeholder, str(val))
            
    # Normalize operators to Python syntax
    pkg_expr = safe_expr.replace(" = ", " == ").replace(" => ", " >= ").replace(" =< ", " <= ")
    # Handle single '=' if it wasn't caught
    if " = " in pkg_expr and " == " not in pkg_expr:
        pkg_expr = pkg_expr.replace(" = ", " == ")
        
    try:
        # Use a restricted eval
        return eval(pkg_expr, {"__builtins__": None}, {})
    except Exception as e:
        raise SyntaxError(f"Error evaluating expression {expr!r}: {e}")

def _get_indent(line: str) -> int:
    if not line.strip(): return -1
    return len(line) - len(line.lstrip())

def _skip_block(lines, start_idx, base_indent):
    idx = start_idx
    while idx < len(lines):
        line, _ = lines[idx]
        if not line.strip():
            idx += 1
            continue
        if _get_indent(line) <= base_indent:
            return idx
        idx += 1
    return idx

def _ensure_indent(lines, current_idx, base_indent, line_no):
    """Ensure the next non-empty line has more indentation than base_indent."""
    idx = current_idx + 1
    while idx < len(lines):
        line, next_line_no = lines[idx]
        if not line.strip():
            idx += 1
            continue
        if _get_indent(line) <= base_indent:
            raise SyntaxError(f"IndentationError on line {next_line_no}: expected an indented block after statement ending in ':'")
        return # Found a properly indented line
    raise SyntaxError(f"IndentationError: expected an indented block at end of file")


# ----------------------------------------------------------------------
# Main execution
# ----------------------------------------------------------------------
def _run_file(path: str):
    import os
    if not os.path.exists(path):
        print(f"Error: The file '{path}' could not be found.")
        sys.exit(1)
    
    # Remove comments using a robust state machine
    def _find_comment_start(line):
        in_double = False
        in_single = False
        escaped = False
        for i, char in enumerate(line):
            if escaped:
                escaped = False
                continue
            if char == "\\":
                escaped = True
                continue
            if char == '"' and not in_single:
                in_double = not in_double
            elif char == "'" and not in_double:
                in_single = not in_single
            elif char == "/" and not in_double and not in_single:
                return i
        return -1

    with open(path, "r", encoding="utf-8") as f:
        lines = f.readlines()

    clean_lines = []
    i = 0
    in_string = False
    string_char = ""
    escaped = False
    
    result = []
    idx = 0
    in_string = False
    string_char = ""
    escaped = False
    full_content = "".join(lines)
    
    while idx < len(full_content):
        char = full_content[idx]
        
        if escaped:
            result.append(char)
            escaped = False
            idx += 1
            continue
            
        if char == "\\":
            result.append(char)
            escaped = True
            idx += 1
            continue
            
        if (char == '"' or char == "'") and not in_string:
            in_string = True
            string_char = char
            result.append(char)
            idx += 1
            continue
            
        if in_string:
            if char == string_char:
                in_string = False
            result.append(char)
            idx += 1
            continue
        
        # Check for comments
        if full_content[idx:idx+2] == "//":
            # Start of a single-line comment
            # Skip until newline
            newline_pos = full_content.find("\n", idx)
            if newline_pos != -1:
                idx = newline_pos # Pointer now at \n
            else:
                idx = len(full_content)
            continue
            
        result.append(char)
        idx += 1
        
    content = "".join(result)

    logical_lines = []
    buffer = ""
    start_line_no = 1
    raw_lines = content.splitlines()
    for i, line in enumerate(raw_lines, 1):
        trimmed = line.rstrip()
        if trimmed.endswith(","):
            # Check if next line exists and is NOT a statement start
            next_line_exists = (i < len(raw_lines))
            next_line_start = raw_lines[i].strip() if next_line_exists else ""
            is_keyword = any(next_line_start.startswith(kw) for kw in ["if ", "elif ", "else:", "var ", "print("])
            
            if next_line_exists and not is_keyword:
                buffer += trimmed[:-1] + " "
                continue
            else:
                # Treat as a statement with a trailing comma
                logical_lines.append((buffer + trimmed[:-1], start_line_no))
                buffer = ""
                start_line_no = i + 1
        else:
            logical_lines.append((buffer + line, start_line_no))
            buffer = ""
            start_line_no = i + 1

    env = {}
    funcs = {} # { name: [ (line, line_no), ... ] }
    
    def _execute_lines(lines):
        nonlocal env, funcs
        i = 0
        last_if_executed = False  # Track if any branch in current if/elif/else chain executed
        
        while i < len(lines):
            line, line_no = lines[i]
            indent = _get_indent(line)
            trimmed = line.strip()
            
            if not trimmed:
                i += 1
                continue

            # --------------------------------------------------------------
            # 0️⃣ Block Conditionals: if / elif / else
            # --------------------------------------------------------------
            
            if trimmed.startswith("if ") and trimmed.endswith(":"):
                expr = trimmed[3:-1].strip()
                cond = _evaluate_expression(expr, env)
                if cond:
                    _ensure_indent(lines, i, indent, line_no)
                    last_if_executed = True
                    i += 1
                else:
                    last_if_executed = False
                    i = _skip_block(lines, i + 1, indent)
                continue

            if trimmed.startswith("elif ") and trimmed.endswith(":"):
                if last_if_executed:
                    i = _skip_block(lines, i + 1, indent)
                else:
                    expr = trimmed[5:-1].strip()
                    cond = _evaluate_expression(expr, env)
                    if cond:
                        _ensure_indent(lines, i, indent, line_no)
                        last_if_executed = True
                        i += 1
                    else:
                        i = _skip_block(lines, i + 1, indent)
                continue
                
            if trimmed == "else:":
                if last_if_executed:
                    i = _skip_block(lines, i + 1, indent)
                else:
                    _ensure_indent(lines, i, indent, line_no)
                    i += 1
                continue

            # --------------------------------------------------------------
            # 0.5️⃣ Functions: func <name>:  and  func.<call>
            # --------------------------------------------------------------
            
            if trimmed.startswith("func ") and trimmed.endswith(":"):
                func_name = trimmed[5:-1].strip()
                _ensure_indent(lines, i, indent, line_no)
                block_start = i + 1
                block_end = _skip_block(lines, block_start, indent)
                funcs[func_name] = lines[block_start : block_end]
                i = block_end
                continue
                
            if trimmed.startswith("func."):
                call_name = trimmed[5:].strip()
                if call_name in funcs:
                    _execute_lines(funcs[call_name])
                else:
                    raise NameError(f"Function '{call_name}' not defined on line {line_no}")
                i += 1
                continue

            # --------------------------------------------------------------
            # 1️⃣ Variable declaration:  var <name> <value>
            # --------------------------------------------------------------
            var_match = re.match(r"^var\s+([A-Za-z_]\w*)\s+(.+)$", trimmed)
            if var_match:
                name, value_expr = var_match.groups()
                var_value = _evaluate_expression(value_expr.strip(), env)
                env[name] = var_value
                i += 1
                continue

            # --------------------------------------------------------------
            # 2️⃣ Print statement:  print(<expr>)
            # --------------------------------------------------------------
            print_match = re.match(r"^print\(\s*(.+?)\s*\)$", trimmed)
            if print_match:
                expr = print_match.group(1).strip()
                try:
                    result = _evaluate_expression(expr, env)
                except SyntaxError:
                    result = _eval_token(expr, env)
                print(result)
                i += 1
                continue

            raise SyntaxError(f"Invalid syntax on line {line_no}: {trimmed!r}")

    # Kick off main execution
    _execute_lines(logical_lines)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python emelang_interpreter.py <script.emlg>")
        sys.exit(1)
    _run_file(sys.argv[1])