CONFIG:
    START: q0
    ACCEPT: success
    REJECT: fail

MAIN:
    q0, 0 -> 0, R, q0
    q0, 1 -> 1, R, q0
    q0, _ -> _, L, q1
    

    q1, 0 -> 0, L, q1
    q1, 1 -> 1, L, q2

    q2, 0 -> 1, L, q2
    q2, 1 -> 0, L, q2
    q2, _ -> _, S, success