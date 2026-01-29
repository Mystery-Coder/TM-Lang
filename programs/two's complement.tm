CONFIG:
    START: q0
    ACCEPT: success
    REJECT: fail

MACROS:
    DEF move_to_end:
        q, 0 -> 0, R, q
        q, 1 -> 1, R, q
        q, _ -> _, L, RETURN


MAIN:
    q0, 0 -> 0, R, CALL move_to_end -> q1
    q0, 1 -> 1, R, CALL move_to_end -> q1
    

    q1, 0 -> 0, L, q1
    q1, 1 -> 1, L, q2

    q2, 0 -> 1, L, q2
    q2, 1 -> 0, L, q2
    q2, _ -> _, S, success