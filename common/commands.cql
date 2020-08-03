create keyspace demo with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table demo.message(id int,email text, title text, content text, magic_number text, PRIMARY KEY(id));


