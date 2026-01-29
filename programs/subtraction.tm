// Unary subtraction (absolute value), expects 111011 -> 1

CONFIG:
    START: q0
    ACCEPT: success
    REJECT: fail

MAIN:
    q0, 1 -> _, R, q1
    q0, 0 -> _, S, success
    
    q1, 1 -> 1, R, q1
    q1, 0 -> 0, R, q1
    q1, _ -> _, L, q2

    q2, 1 -> _, L, q3
    q2, 0 -> 1, S, success

    q3, 1 -> 1, L, q3
    q3, 0 -> 0, L, q3
    q3, _ -> _, R, q0
