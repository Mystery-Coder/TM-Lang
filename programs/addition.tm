// Turing Machine for unary addition, expects 110111 -> 11111

CONFIG:
    START: start
    ACCEPT: success
    REJECT: fail


MAIN:
    start, 1 -> 1, R, start
    start, 0 -> 1, R, q1

    q1, 1 -> 1, R, q1
    q1, _ -> _, L, q2
    q2, 1 -> _, S, success