// Takes a string, 1011 -> 1101

CONFIG:
    START: q0
    ACCEPT: success
    REJECT: fail

MAIN:
    q0, 1 -> 1, R, q0
    q0, 0 -> 0, R, q0
    q0, _ -> _, L, q1

    q1, 1 -> X, R, q2
    q1, 0 -> X, R, q3
    q1, X -> X, L, q1
    q1, _ -> _, R, q5
  

    q2, X -> X, R, q2
    q2, 0 -> 0, R, q2
    q2, 1 -> 1, R, q2
    q2, _ -> 1, L, q4

    q3, X -> X, R, q3
    q3, 0 -> 0, R, q3
    q3, 1 -> 1, R, q3
    q3, _ -> 0, L, q4

    q4, X -> X, L, q1
    q4, 1 -> 1, L, q4
    q4, 0 -> 0, L, q4

    q5, X -> _, R, q5
    q5, 1 -> 1, S, success
    q5, 0 -> 0, S, success