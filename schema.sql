create table entries
(
    entry_id          bigserial primary key,
    account_id        bigint,
    previous_entry_id bigint,
    balance           bigint,
    data              text
);

create index ix_account_id_entry_id
    on entries (account_id, entry_id)
    -- Use a low fillfactor to cause more page splits
    with (fillfactor = 10);

INSERT INTO entries (account_id, previous_entry_id, balance, data)
SELECT s,
       NULL,
       0,
       NULL
FROM generate_series(1, 10000) AS s(i);


explain select max(entry_id) from entries where account_id = 500;