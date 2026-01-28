CONFIG:
    START: start
    ACCEPT: done
    REJECT: fail

MACROS:
    DEF move_right_end:
        q, 0 -> 0, R, q
        q, 1 -> 1, R, q
        q, _ -> _, L, RETURN
    DEF move_left_end:
        q, 0 -> 0, L, q
        q, 1 -> 1, L, q
        q, _ -> _, R, RETURN

MAIN:
    start, 0 -> _, R, CALL move_right_end -> q1
    start, 1 -> _, R, CALL move_right_end -> q2
    start, _ -> _, S, done
    
    q1, 0 -> _, L, CALL move_left_end -> start
    q1, 1 -> 1, S, fail
    q1, _ -> _, S, done

    q2, 1 -> _, L, CALL move_left_end -> start
    q2, 0 -> 0, S, fail
    q2, _ -> _, S, done