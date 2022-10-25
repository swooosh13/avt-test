create table balance
(
    user_id bigint primary key,
    amount  bigint
);

create table transactions
(
    user_id    bigint      not null,
    service_id bigint,
    order_id   bigint,
    amount     bigint      not null,
    from_info       varchar(45) not null,
    to_info         varchar(45) not null,
    created_at timestamp default now()
);

create table reservation
(
    user_id bigint primary key,
    amount bigint
);

-- перевод между балансами пользователей
create procedure transfer(from_user_id bigint, to_user_id bigint, amount_t bigint)
    language plpgsql
as
$$
begin
    -- если хоть один из пользователей не существует
    if (select user_id from balance where user_id = from_user_id) is null or (select user_id from balance where user_id = to_user_id) is null  then
        RAISE EXCEPTION 'Invalid sender or receiver';
    END IF;

    -- если баланс меньше
    IF (select amount from balance where user_id = from_user_id) < amount_t THEN
        RAISE EXCEPTION 'Not enough money';
    END IF;

    update balance set amount = balance.amount - amount_t  where user_id = from_user_id;
    update balance set amount = balance.amount + amount_t where user_id = to_user_id;

    INSERT INTO transactions (user_id, service_id, order_id, amount, from_info, to_info)
    VALUES (from_user_id, null, null, amount_t, 'user-' || from_user_id, 'user-' || to_user_id);
end;
$$;

-- резервирование средств
create procedure reserve_from_user(from_user_id bigint, service_id_t bigint, order_id_t bigint, amount_t bigint)
    language plpgsql
as $$
begin
    -- если пользователь не существует
    if (select user_id from balance where user_id = from_user_id) is null then
        RAISE EXCEPTION 'Invalid user';
    END IF;

    IF (select amount from balance where user_id = from_user_id) < amount_t THEN
        RAISE EXCEPTION 'Not enough money';
    END IF;

    update balance set amount = balance.amount - amount_t  where user_id = from_user_id;
    INSERT INTO reservation (user_id, amount) VALUES (from_user_id, amount_t)
    ON CONFLICT (user_id) DO update SET amount = reservation.amount + amount_t;

    INSERT INTO transactions (user_id, service_id, order_id, amount, from_info, to_info)
    VALUES (from_user_id, service_id_t, order_id_t, amount_t, 'balance-' || from_user_id, 'reservation-' || from_user_id);
end;
$$;


-- возврат средств на баланс
create procedure reserve_to_user(to_user_id bigint, service_id_t bigint, order_id_t bigint, amount_t bigint)
    language plpgsql
as $$
begin
    if (select user_id from reservation where user_id = to_user_id) is null then
        RAISE EXCEPTION 'Invalid user';
    END IF;

    IF (select amount from reservation where user_id = to_user_id) < amount_t THEN
        RAISE EXCEPTION 'Not enough money';
    END IF;

    update reservation set amount = reservation.amount - amount_t  where user_id = to_user_id;

    INSERT INTO balance (user_id, amount) VALUES (to_user_id, amount_t)
    ON CONFLICT (user_id) DO update SET amount = balance.amount + amount_t;

    INSERT INTO transactions (user_id, service_id, order_id, amount, from_info, to_info)
    VALUES (to_user_id, service_id_t, order_id_t, amount_t, 'reservation-' || to_user_id, 'balance-'|| to_user_id);
end;
$$;

-- признание выручки
create procedure remove_reserve_from_user(from_user_id bigint, service_id_t bigint, order_id_t bigint, amount_t bigint)
    language plpgsql
as $$
begin
    -- если пользователь не существует
    if (select user_id from reservation where user_id = from_user_id) is null then
        RAISE EXCEPTION 'Invalid sender';
    END IF;

    -- если баланс меньше
    IF (select amount from reservation  where user_id = from_user_id) < amount_t THEN
        RAISE EXCEPTION 'Not enough money';
    END IF;

    update reservation set amount = reservation.amount - amount_t  where user_id = from_user_id;

    INSERT INTO transactions (user_id, service_id, order_id, amount, from_info, to_info)
    VALUES (from_user_id, service_id_t, order_id_t, amount_t, 'reservation-' || from_user_id, 'payment');
end;
$$;
