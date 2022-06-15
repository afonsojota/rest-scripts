create table input_file
(
    id     int auto_increment primary key,
    name   varchar(100)  not null,
    owner  varchar(30)   null,
    status varchar(10) default 'PENDING'         null
);

create table send_ticket
(
    ticket_id     bigint                              not null,
    solution_id int         default 0                 null,
    site_id     varchar(10)                           null,
    status      varchar(10) default 'PENDING'         null,
    timestamp   timestamp   default CURRENT_TIMESTAMP null,
    file_id     int                                 null,
    primary key (ticket_id, timestamp),
    constraint send_ticket_input_file_ID_fk
        foreign key (file_id) references input_file (ID)
);

create table send_errors
(
    ticket_id   bigint                              not null,
    timestamp timestamp default CURRENT_TIMESTAMP not null,
    error     varchar(500)                        null,
    PRIMARY KEY (ticket_id, timestamp)
);

