// E2pects input as 1101110, 223 = 6 , 1101110111111, a*b
CONFIG:
    START: start
    ACCEPT: success
    REJECT: fail

MACROS:
    DEF move_to_b:
        q, 1 -> 1, R, q
        q, 0 -> 0, R, RETURN
    DEF move_to_end:
        q, 1 -> 1, R, q
        q, 0 -> 0, R, q
        q, _ -> _, L, RETURN
    DEF move_to_b_start:
        q, 1 -> 1, L, q
        q, 0 -> 0, L, p
        p, 1 -> 1, L, p
        p, 2 -> 2, R, RETURN
    DEF rewrite_b:
        q, 2 -> 1, L, q
        q, 0 -> 0, L, RETURN

MAIN:
    start, 1 -> _, R, CALL move_to_b -> q1
    start, 0 -> 0, R, success
    start, _ -> _, R, start

    q1, 1 -> 2, R, CALL move_to_end -> q2
    q1, 0 -> 0, L, CALL rewrite_b -> start

    q2, _ -> 1, L, CALL move_to_b_start -> q1
    q2, 1 -> 1, R, q2
    q2, 0 -> 0, R, q2