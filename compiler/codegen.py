import json

class CodeGenerator:
    def __init__(self, ir):
        self.meta = ir['meta']
        self.transitions = ir['transitions']
    
    def generate_c(self):
        print("--- Generating C Code (Verbose Mode) ---")
        
        # 1. Collect States
        states = set()
        states.add(self.meta['start'])
        states.add(self.meta['accept'])
        states.add(self.meta['reject'])
        for t in self.transitions:
            states.add(t['src'])
            states.add(t['next'])
            
        state_list = sorted(list(states))
        state_map = { name: i for i, name in enumerate(state_list) }
        
        state_comment = "\n".join([f"// {name}: {i}" for name, i in state_map.items()])

        # C template
        c_code = f"""
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define TAPE_SIZE 20000
#define HEAD_START 10000

/* --- STATE MAP --- 
{state_comment}
*/

int current_state = {state_map[self.meta['start']]};
int ACCEPT_STATE = {state_map[self.meta['accept']]};
int REJECT_STATE = {state_map[self.meta['reject']]};

char tape[TAPE_SIZE];
int head = HEAD_START;

void print_tape() {{
    printf("\\r[ ");
    for(int i = head - 10; i <= head + 10; i++) {{
        if(i == head) printf("[%c]", tape[i]);
        else printf(" %c ", tape[i]);
    }}
    printf(" ] State: %d  ", current_state);
    fflush(stdout); 
}}

int main() {{
    memset(tape, '_', TAPE_SIZE);
    
    printf("Enter Input: ");
    char input[100];
    scanf("%s", input);
    
    for(int i=0; i<strlen(input); i++) {{
        tape[head + i] = input[i];
    }}

    printf("\\n--- RUNNING ---\\n");

    while(1) {{
        print_tape();
        system("sleep 0.3");

        if (current_state == ACCEPT_STATE) {{ printf("\\n\\nACCEPTED!\\n"); return 0; }}
        if (current_state == REJECT_STATE) {{ printf("\\n\\nREJECTED!\\n"); return 1; }}

        char read_val = tape[head];
        int matched = 0;

        switch(current_state) {{
"""

        #  Generate Logic
        grouped = {}
        for t in self.transitions:
            src = state_map[t['src']]
            if src not in grouped: grouped[src] = []
            grouped[src].append(t)

        for state_id, rules in grouped.items():
            c_code += f"            case {state_id}:\n"
            
            first = True
            for rule in rules:
                # SIMPLE LOGIC: Explicit Check Only
                condition = f"read_val == '{rule['read']}'"
                
                prefix = "if" if first else "else if"
                first = False
                
                move_code = "head++;" if rule['dir'] == 'R' else ("head--;" if rule['dir'] == 'L' else "")
                
                next_id = state_map[rule['next']]

                c_code += f"""                {prefix} ({condition}) {{
                    tape[head] = '{rule['write']}';
                    {move_code}
                    current_state = {next_id};
                    matched = 1;
                }}
"""
            c_code += "                break;\n"

        c_code += """
        } // End Switch
        
        if (!matched) {
            printf("\\n\\nCRASH: State %d has no rule for char '%c'\\n", current_state, read_val);
            return 1;
        }
    }
}
"""
        return c_code

    # GRAPHVIZ GENERATOR
    def generate_dot(self):
        print("--- Generating Diagram ---")
        dot = "digraph TuringMachine {\n"
        dot += '    rankdir=LR;\n'
        dot += '    node [shape = circle];\n'
        
        # 1. Special Shapes
        dot += f'    "{self.meta["accept"]}" [shape = doublecircle, color=green];\n'
        dot += f'    "{self.meta["reject"]}" [shape = doublecircle, color=red];\n'
        
        # 2. Entry Point
        dot += '    entry [shape = point];\n'
        dot += f'    entry -> "{self.meta["start"]}";\n'

        # 3. Group Transitions by (Source -> Destination)
        edges = {} 
        
        for t in self.transitions:
            key = (t['src'], t['next'])
            
            # Clean Labels
            r_lbl = "BLANK" if t['read'] == '_' else t['read']
            w_lbl = "BLANK" if t['write'] == '_' else t['write']
            label = f"{r_lbl} / {w_lbl}, {t['dir']}"
            
            if key not in edges:
                edges[key] = []
            edges[key].append(label)

        # 4. Generate Merged Edges
        for (src, dst), labels in edges.items():
            # Join all labels with a newline character for GraphViz
            combined_label = "\\n".join(labels)
            dot += f'    "{src}" -> "{dst}" [label = "{combined_label}"];\n'

        dot += "}\n"
        return dot