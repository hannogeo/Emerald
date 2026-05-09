import re
content = """/ Emerald Language sample script
/ Declare variables
var message - "Hello, Emerald Language!"
var version - 1.0
var is_working - True
var year - 2024

/
   This is a multiline comment
   demonstrating the new syntax.
/

/ Print variables
print(var.message)
print(var.version)
print(var.is_working)
print(var.year)"""

pattern = r'("(?:\\.|[^"\\])*")|(\'(?:\\.|[^\'\\])*\')|(/\s*?\n[\s\S]*?/\s*?(?:\n|$))|(/.*)'

def _comment_replacer(match):
    if match.group(3):
        print(f"Matched multiline: {repr(match.group(3))}")
        return "\n" * match.group(3).count("\n")
    if match.group(4):
        print(f"Matched single: {repr(match.group(4))}")
        return ""
    return match.group(0)

new_content = re.sub(pattern, _comment_replacer, content)
print("--- Result ---")
for i, line in enumerate(new_content.splitlines(), 1):
    print(f"{i}: {repr(line)}")
