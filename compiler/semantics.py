import json

class SemanticAnalyzer:
    def __init__(self, ir):
        self.ir = ir
        self.expanded_main = [] # The final flat list for C generation
        self.macro_counter = 0  # Counter to ensure unique state names

    def analyze(self):
        print("--- Semantic Analysis & Macro Expansion ---")
        
        # 1. Verification: Check Entry Point
        start_node = self.ir['meta'].get('start')
        if not start_node:
            raise Exception("Semantic Error: No START state defined in CONFIG.")
        
        # 2. Flattening: Process every transition in MAIN
        for trans in self.ir['main']:
            self.process_transition(trans)
            
        return {
            "meta": self.ir['meta'],
            "transitions": self.expanded_main
        }

    def process_transition(self, trans):
        target = trans['target']
        
        # CASE 1: Standard GOTO (Simple copy)
        # If it's a normal jump, we just copy it to the final list.
        if target['type'] == 'GOTO':
            self.expanded_main.append({
                "src": trans['src'],
                "read": trans['read'],
                "write": trans['write'],
                "dir": trans['dir'],     # Matches Parser output key
                "next": target['name']   # The simple next state
            })

        # CASE 2: Macro CALL (The Complex Part)
        elif target['type'] == 'CALL':
            macro_name = target['name']
            return_state = target['ret']
            
            print(f"   > Expanding Macro: '{macro_name}' (Returns to '{return_state}')")
            
            # Fetch the macro definition
            macro_body = self.ir['macros'].get(macro_name)
            if not macro_body:
                raise Exception(f"Semantic Error: Call to undefined macro '{macro_name}'")
                
            # 1. Generate Unique Prefix
            # e.g., If macro is called twice, we get "move_right_1_" and "move_right_2_"
            self.macro_counter += 1
            prefix = f"{macro_name}_{self.macro_counter}_"
            
            # 2. Link Main -> Macro Start
            # We assume the first rule in the macro defines the start state
            macro_start_original = macro_body[0]['src']
            macro_start_renamed = prefix + macro_start_original
            
            # Add the "Bridge" transition
            self.expanded_main.append({
                "src": trans['src'],
                "read": trans['read'],
                "write": trans['write'],
                "dir": trans['dir'],
                "next": macro_start_renamed
            })
            
            # 3. Inject the Macro Body (Renamed)
            for m_trans in macro_body:
                # Rename the source state of this internal transition
                new_src = prefix + m_trans['src']
                
                # Determine where this internal transition goes
                m_target = m_trans['target']
                new_next = ""
                
                if m_target['type'] == 'GOTO':
                    # Internal Jump: Rename it so it stays inside this macro instance
                    new_next = prefix + m_target['name']
                    
                elif m_target['type'] == 'RETURN':
                    # Exit Jump: Connect it to the Return State defined in MAIN
                    new_next = return_state
                    
                elif m_target['type'] == 'CALL':
                    raise Exception("Nested Macros (Macro calling Macro) are not supported in v1.0")

                # Add this internal transition to the final list
                self.expanded_main.append({
                    "src": new_src,
                    "read": m_trans['read'],
                    "write": m_trans['write'],
                    "dir": m_trans['dir'],
                    "next": new_next
                })