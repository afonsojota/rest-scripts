create table send_ticket
(
    ticket_id     bigint                              not null,
    solution_id int         default 0                 null,
    site_id     varchar(10)                           null,
    status      varchar(10) default 'PENDING'         null,
    timestamp   timestamp   default CURRENT_TIMESTAMP not null,
    file_id     int                                 null,
    primary key (ticket_id, timestamp),
    constraint send_ticket_input_file_ID_fk
        foreign key (file_id) references input_file (ID)
);


create index send_ticket_ticket_id_index
on send_ticket (ticket_id);


create table send_errors
(
    ticket_id   bigint                              not null,
    timestamp timestamp default CURRENT_TIMESTAMP not null,
    error     varchar(500)                        null,
    PRIMARY KEY (ticket_id, timestamp)
);