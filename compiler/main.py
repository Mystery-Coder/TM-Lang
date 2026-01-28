import sys
import os
from lexer import Lexer
from parser import Parser
from semantics import SemanticAnalyzer
from codegen import CodeGenerator

def main():
   
    if len(sys.argv) < 2:
        print("Error: No input file provided.")
        print("Usage: python compiler/main.py programs/palindrome.tm")
        return

    filepath = sys.argv[1]
    

    if not os.path.exists(filepath):
        print(f"Error: Input file '{filepath}' not found.")
        return

    try:
        with open(filepath, 'r') as f:
            code = f.read()
            
        print(f"--- Compiling '{filepath}' ---")

        print("--- 1. Lexer ---")
        lexer = Lexer(code)
        tokens = lexer.tokenize()

        print("--- 2. Parser ---")
        parser = Parser(tokens)
        ir = parser.parse()

        print("--- 3. Semantics ---")
        analyzer = SemanticAnalyzer(ir)
        final_ir = analyzer.analyze()

        print("--- 4. Code Gen ---")
        codegen = CodeGenerator(final_ir)
        
        
        base_name = os.path.basename(filepath).replace('.tm', '')
        output_dir = os.path.join("build")

        # FIX 1: Create the directory if it doesn't exist
        os.makedirs(output_dir, exist_ok=True)
        print(f"   > Output directory: {output_dir}/")

        # Define file paths
        c_path = os.path.join(output_dir, f"{base_name}.c")
        dot_path = os.path.join(output_dir, f"{base_name}.dot")
        svg_path = os.path.join(output_dir, f"{base_name}.svg")

        # Generate C
        with open(c_path, "w") as f:
            f.write(codegen.generate_c())
            
        # Generate Diagram
        with open(dot_path, "w") as f:
            f.write(codegen.generate_dot())
        
        print("--- Converting to SVG ---")
        try:
            # FIX 2: Pass the full path to 'dot'
            cmd = f"dot -Tsvg \"{dot_path}\" -o \"{svg_path}\""
            os.system(cmd)
            print(f"Success! Diagram generated at '{svg_path}'")
        except Exception:
            print("Warning: GraphViz not found or failed.")

        print(f"\nDONE! Output saved to '{output_dir}/'")

    except Exception as e:
        # Improved error printing to see the REAL error
        import traceback
        traceback.print_exc()
        print(f"CRITICAL ERROR: {e}")

if __name__ == "__main__":
    main()