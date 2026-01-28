import re


'''
Example Code,

CONFIG:
    START: q0
    ACCEPT: done
    REJECT: fail
MACROS:
    DEF seek_blank:
        q0, 0 -> 1, R, q0
        q0, _ -> _, S, RETURN
MAIN:
    q0, 1 -> 0, R, CALL seek_blank
'''

# Define Token Types
class TokenType:
    SECTION   = 'SECTION'   # CONFIG:, MACROS:, MAIN:
    KEYWORD   = 'KEYWORD'   # START:, ACCEPT:, REJECT:, DEF, CALL, RETURN
    ID        = 'ID'        # Identifiers (q0, my_macro)
    SYMBOL    = 'SYMBOL'    # 0, 1, *, _
    DIRECTION = 'DIR'       # L, R, S
    ARROW     = 'ARROW'     # ->
    COMMA     = 'COMMA'     # ,
    COLON     = 'COLON'     # :
    NEWLINE   = 'NEWLINE'   # \n
    EOF       = 'EOF'       # End of File

class Token:
    def __init__(self, type, value, line):
        self.type = type
        self.value = value
        self.line = line
    
    def __repr__(self):
        return f"Token({self.type}, '{self.value}', Line {self.line})"

class Lexer:
    def __init__(self, source_code):
        self.source = source_code
        self.tokens = []
        self.current_line = 1

        # Regex Patterns (Order is CRITICAL)
        self.rules = [
           
            (TokenType.SECTION,   r'(CONFIG:|MACROS:|MAIN:)'),

          
            (TokenType.KEYWORD,   r'(START:|ACCEPT:|REJECT:)'), 

          
            (TokenType.KEYWORD,   r'\b(DEF|CALL|RETURN)\b'),

            
            (TokenType.ARROW,     r'->'),
            (TokenType.COMMA,     r','),
            (TokenType.COLON,     r':'),

            
            (TokenType.DIRECTION, r'\b(L|R|S)\b'),

            # Identifiers (Must come BEFORE Symbol to catch 'q0')
            # IDs can only start with letters not '_'
            (TokenType.ID,        r'\b[a-zA-Z][a-zA-Z0-9_]*\b'), 
            
            # Symbols (Digits, wildcard, underscore)
            (TokenType.SYMBOL,    r'[0-9*_]'), 
            
            # Formatting & Skippables
            (TokenType.NEWLINE,   r'\n'),
            ('SKIP',              r'[ \t]+'),       # Spaces/Tabs
            ('COMMENT',           r'//.*'),         # Comments
            ('MISMATCH',          r'.'),            # Error catcher
        ]

    def tokenize(self):
        pos = 0
        while pos < len(self.source):
            match = None
            for token_type, pattern in self.rules:
                regex = re.compile(pattern)
                match = regex.match(self.source, pos)
                if match:
                    text = match.group(0)
                    
                    if token_type == TokenType.NEWLINE:
                        self.current_line += 1
                    elif token_type == 'SKIP' or token_type == 'COMMENT':
                        pass
                    elif token_type == 'MISMATCH':
                        # Unexpected Token
                        print(f"Lexer Error: Unexpected character '{text}' at line {self.current_line}")
                        return []
                    else:
                        self.tokens.append(Token(token_type, text, self.current_line))
                    
                    pos = match.end()
                    break
            
            if not match:
                print("Fatal Lexing Error: Logic Halt")
                break
        
        self.tokens.append(Token(TokenType.EOF, '', self.current_line))
        return self.tokens

