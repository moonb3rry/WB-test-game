create table users(
                      user_id serial primary key,
                      username varchar(50),
                      password varchar(200),
                      user_role varchar(1)
);

create table customers(
                          user_id int,
                          capital int
);

create table loaders(
                        user_id int,
                        max_weight int,
                        alcoholism bool,
                        fatigue int,
                        wage int
);

create table tasks(
                      task_id serial primary key,
                      task_name varchar(100),
                      weight int,
                      status bool,
                      customer_id int
);

create table assigned_loaders(
                                 loader_id int,
                                 task_id int
);

create table game(
                     game_id serial primary key,
                     user_id int,
                     game_result bool
);

alter table customers add constraint user_fk foreign key (user_id) references users(user_id);

alter table loaders add constraint user_fk foreign key (user_id) references users(user_id);

alter table tasks add constraint customer_fk foreign key (customer_id) references users(user_id);

alter table assigned_loaders add constraint loader_fk foreign key (loader_id) references users(user_id);
alter table assigned_loaders add constraint task_fk foreign key (task_id) references tasks(task_id);

alter table game add constraint user_fk foreign key (user_id) references users(user_id);