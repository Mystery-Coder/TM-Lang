
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define TAPE_SIZE 20000
#define HEAD_START 10000

/* --- STATE MAP --- 
// done: 0
// fail: 1
// move_left_end_3_q: 2
// move_left_end_4_q: 3
// move_right_end_1_q: 4
// move_right_end_2_q: 5
// q1: 6
// q2: 7
// start: 8
*/

int current_state = 8;
int ACCEPT_STATE = 0;
int REJECT_STATE = 1;

char tape[TAPE_SIZE];
int head = HEAD_START;

void print_tape() {
    printf("\r[ ");
    for(int i = head - 10; i <= head + 10; i++) {
        if(i == head) printf("[%c]", tape[i]);
        else printf(" %c ", tape[i]);
    }
    printf(" ] State: %d  ", current_state);
    fflush(stdout); 
}

int main() {
    memset(tape, '_', TAPE_SIZE);
    
    printf("Enter Input: ");
    char input[100];
    scanf("%s", input);
    
    for(int i=0; i<strlen(input); i++) {
        tape[head + i] = input[i];
    }

    printf("\n--- RUNNING ---\n");

    while(1) {
        print_tape();
        system("sleep 0.7");

        if (current_state == ACCEPT_STATE) { printf("\n\nACCEPTED!\n"); return 0; }
        if (current_state == REJECT_STATE) { printf("\n\nREJECTED!\n"); return 1; }

        char read_val = tape[head];
        int matched = 0;

        switch(current_state) {
            case 8:
                if (read_val == '0') {
                    tape[head] = '_';
                    head++;
                    current_state = 4;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '_';
                    head++;
                    current_state = 5;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    
                    current_state = 0;
                    matched = 1;
                }
                break;
            case 4:
                if (read_val == '0') {
                    tape[head] = '0';
                    head++;
                    current_state = 4;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '1';
                    head++;
                    current_state = 4;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    head--;
                    current_state = 6;
                    matched = 1;
                }
                break;
            case 5:
                if (read_val == '0') {
                    tape[head] = '0';
                    head++;
                    current_state = 5;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '1';
                    head++;
                    current_state = 5;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    head--;
                    current_state = 7;
                    matched = 1;
                }
                break;
            case 6:
                if (read_val == '0') {
                    tape[head] = '_';
                    head--;
                    current_state = 2;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '1';
                    
                    current_state = 1;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    
                    current_state = 0;
                    matched = 1;
                }
                break;
            case 2:
                if (read_val == '0') {
                    tape[head] = '0';
                    head--;
                    current_state = 2;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '1';
                    head--;
                    current_state = 2;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    head++;
                    current_state = 8;
                    matched = 1;
                }
                break;
            case 7:
                if (read_val == '1') {
                    tape[head] = '_';
                    head--;
                    current_state = 3;
                    matched = 1;
                }
                else if (read_val == '0') {
                    tape[head] = '0';
                    
                    current_state = 1;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    
                    current_state = 0;
                    matched = 1;
                }
                break;
            case 3:
                if (read_val == '0') {
                    tape[head] = '0';
                    head--;
                    current_state = 3;
                    matched = 1;
                }
                else if (read_val == '1') {
                    tape[head] = '1';
                    head--;
                    current_state = 3;
                    matched = 1;
                }
                else if (read_val == '_') {
                    tape[head] = '_';
                    head++;
                    current_state = 8;
                    matched = 1;
                }
                break;

        } // End Switch
        
        if (!matched) {
            printf("\n\nCRASH: State %d has no rule for char '%c'\n", current_state, read_val);
            return 1;
        }
    }
}
