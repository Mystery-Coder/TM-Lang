from lexer import TokenType

class Parser:
    def __init__(self, tokens):
        self.tokens = tokens
        self.pos = 0
        self.current_token = self.tokens[0]
        
        # The Output Structure (Intermediate Representation)
        self.ir = {
            "meta": {},      # Stores START, ACCEPT, REJECT
            "macros": {},    # Stores macro definitions
            "main": []       # Stores main logic
        }

   
    def advance(self):
        self.pos += 1
        if self.pos < len(self.tokens):
            self.current_token = self.tokens[self.pos]
        else:
            self.current_token = None

   
    def consume(self, token_type):
        if self.current_token and self.current_token.type == token_type:
            val = self.current_token.value
            self.advance()
            return val
        else:
            current_val = self.current_token.type if self.current_token else "EOF"
            raise Exception(f"Syntax Error: Expected {token_type} but got {current_val} at Line {self.current_token.line if self.current_token else 'End'}")

    # Main Parser
    def parse(self):
       
        if self.current_token.type == TokenType.SECTION and self.current_token.value == "CONFIG:":
            self.parse_config()
        else:
            raise Exception("Syntax Error: Code must start with 'CONFIG:' section")
        
       
        if self.current_token.type == TokenType.SECTION and self.current_token.value == "MACROS:":
            self.parse_macros()

        if self.current_token.type == TokenType.SECTION and self.current_token.value == "MAIN:":
            self.parse_main()
        else:
            raise Exception("Syntax Error: Missing 'MAIN:' section")
            
        return self.ir

    
    def parse_config(self):
        print("Parsing Config...")
        self.consume(TokenType.SECTION) # Expect 'CONFIG:'
        
        # Expects 3 specific lines: START, ACCEPT, REJECT
        for _ in range(3):
            key = self.consume(TokenType.KEYWORD) # e.g. 'START:'
            val = self.consume(TokenType.ID)      # e.g. 'q0'
            
            if key == "START:": self.ir["meta"]["start"] = val
            elif key == "ACCEPT:": self.ir["meta"]["accept"] = val
            elif key == "REJECT:": self.ir["meta"]["reject"] = val

    def parse_macros(self):
        print("Parsing Macros...")
        self.consume(TokenType.SECTION)
 
        # Loop until we hit the 'MAIN:' section

        while self.current_token.type != TokenType.SECTION:
            self.consume(TokenType.KEYWORD) # Expect 'DEF'
            macro_name = self.consume(TokenType.ID)
            self.consume(TokenType.COLON)
            
            transitions = []
            # Keep parsing transitions until we hit another DEF or the MAIN section
            # We check if current token is an ID (start of transition) 
            while self.current_token.type == TokenType.ID:
                transitions.append(self.parse_transition())
            
            self.ir["macros"][macro_name] = transitions

    def parse_main(self):
        print("Parsing Main Logic...")
        self.consume(TokenType.SECTION)

        while self.current_token.type != TokenType.EOF:
            self.ir["main"].append(self.parse_transition())

    # --- The Core Logic: Parsing a single line ---
    def parse_transition(self):
        # Source State
        src = self.consume(TokenType.ID)
        self.consume(TokenType.COMMA)
        
        # Read Symbol
        read = self.consume(TokenType.SYMBOL)
        self.consume(TokenType.ARROW)
        
        # Write Symbol
        write = self.consume(TokenType.SYMBOL)
        self.consume(TokenType.COMMA)
        
        # Direction
        direction = self.consume(TokenType.DIRECTION)
        self.consume(TokenType.COMMA)
        
        # Target (The Logic Change is Here)
        target = {}
        
        if self.current_token.type == TokenType.KEYWORD:
            kw = self.consume(TokenType.KEYWORD)
            
            if kw == "CALL":
                # Logic: CALL <macro_name> -> <return_state>
                macro_name = self.consume(TokenType.ID)
                
                # Check for the return arrow
                self.consume(TokenType.ARROW)
                ret_state = self.consume(TokenType.ID)
                
                target = {
                    "type": "CALL", 
                    "name": macro_name, 
                    "ret": ret_state 
                }
                
            elif kw == "RETURN":
                # Logic: RETURN (Inside a macro, just ends)
                target = {"type": "RETURN"}
                
        else:
            # Logic: <state_name> (Standard GOTO)
            target_name = self.consume(TokenType.ID)
            target = {"type": "GOTO", "name": target_name}

        # Return the Transition Object
        return {
            "src": src,
            "read": read,
            "write": write,
            "dir": direction,
            "target": target
        }