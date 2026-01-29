# TM-Lang Specification

## 1. Introduction

TM-Lang is a Domain Specific Language (DSL) designed for defining Turing Machines. It abstracts complex state transition tuples into a readable, modular syntax supporting Macros. The compiler targets C for simulation.

## 2. File Structure

A source file (.tm) is processed sequentially and consists of three required sections:

- **CONFIG**: Defines entry and exit states.
- **MACROS**: Defines reusable subroutines (can be empty).
- **MAIN**: Defines the primary transition logic.

## 3. Lexical Specification

| Token   | Pattern                   | Description                    |
| ------- | ------------------------- | ------------------------------ |
| SECTION | (CONFIG:\|MACROS:\|MAIN:) | Section Headers                |
| KEYWORD | START:, ACCEPT:, REJECT:  | Configuration Keys             |
| LOGIC   | DEF, CALL, RETURN         | Macro Logic                    |
| ID      | [a-zA-Z][a-zA-Z0-9_]\*    | State or Macro Names           |
| SYMBOL  | 0, 1, \_                  | Tape Alphabet                  |
| DIR     | L, R, S                   | Directions (Left, Right, Stay) |
| ARROW   | ->                        | Transition Operator            |

## 4. Syntax & Grammar

### 4.1 Configuration

Must appear at the top. Defines the specific state names for machine lifecycle events.

```
CONFIG:
    START:  <id>
    ACCEPT: <id>
    REJECT: <id>
```

### 4.2 Macro Definitions (DEF)

Macros act as subroutines. They are expanded inline by the compiler.

- **Return**: Use the RETURN keyword to exit the macro.

```
MACROS:
    DEF <macro_name>:
        <src_state>, <read_sym> -> <write_sym>, <dir>, <next_state>
        <src_state>, <read_sym> -> <write_sym>, <dir>, RETURN
```

### 4.3 Main Logic (MAIN)

The entry point of the machine.

- **Standard Transition**: Jump from one state to another.
- **Macro Call**: Execute a macro and specify where to jump upon return.

Syntax:

```
MAIN:
    // Standard GOTO
    <state>, <read> -> <write>, <dir>, <next_state>

    // Macro Call (Must specify return state)
    <state>, <read> -> <write>, <dir>, CALL <macro_name> -> <return_state>
```

## 5. Compiler Semantics

### 5.1 Macro Expansion

The compiler performs a "Search and Replace" strategy for macros:

- **Injection**: The CALL line is replaced by the body of the macro.
- **Scoping**: Internal states of the macro are renamed to avoid collisions (e.g., q0 becomes macro_id_q0).
- **Linking**:
    - The CALL transition connects to the Macro's Start State.
    - The Macro's RETURN transitions connect to the return_state specified in the call.

### 5.3 Code Generation Examples

| Logic    | Code in .tm         | C Code Output                                |
| -------- | ------------------- | -------------------------------------------- |
| Specific | "q0, 0 -> 1, R, q1" | if (read_val == '0') { tape[head]='1'; ... } |

## 6. Example Program (Binary Incrementer)

```
CONFIG:
    START: start
    ACCEPT: done
    REJECT: fail
MACROS:
    DEF move_end:
        q0, 0 -> 0, R, q0
        q0, 1 -> 1, R, q0
        q0, _ -> _, L, RETURN
MAIN:
    // Go to end of number
    start, 0 -> 0, S, CALL move_end -> add
    start, 1 -> 1, S, CALL move_end -> add
    start, _ -> _, S, CALL move_end -> add

    // Add 1 logic
    add, 0 -> 1, S, done       // 0->1, Finished
    add, 1 -> 0, L, add        // 1->0, Carry Left
    add, _ -> 1, S, done       // Overflow
```

The compiler generates C code to simulate a Turing Machine and a GraphViz dot file with svg of the State Diagram

## 7. Requirements

- **Python 3.x**: The compiler is written in Python and requires Python 3 to run.
- **GraphViz**: Required for generating SVG diagrams from the dot files. Install from [graphviz.org](https://graphviz.org/download/).
- **C Compiler**: A C compiler (e.g., GCC) is needed to compile the generated `simulation.c` file into an executable.

---

# Todo List

- [ ] Rewrite compiler in Go
- [ ] Sandbox in Nextjs
