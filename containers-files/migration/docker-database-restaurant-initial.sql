create or replace function next_user_id()
returns int as $$
declare
  values integer[] := array(select (id) from users order by id);
  length integer = array_length(values, 1);
  i int = 1;
begin
  if values[1] is null then 
  	return 1;
  end if;
  while i < length loop
    if values[i+1] != values[i]+1 then
      return (values[i]+1);
    end if;
      i = i+1;
  end loop;
  return (select max(id)+1 from users);
end;
$$ language plpgsql;

CREATE TABLE users (
   id INTEGER PRIMARY KEY default next_user_id(), 
   name VARCHAR(80) NOT NULL,
	 email VARCHAR(80) UNIQUE NOT NULL,
	 passwd VARCHAR(80) NOT NULL,
	 img BYTEA 
);

create or replace function next_room_id ()
returns int as $$
declare
  values integer[] := array(select (id) from rooms order by id);
  length integer = array_length(values, 1);
  i int = 1;
begin
  if values[1] is null then 
  	return 1;
  end if;
  while i < length loop
    if values[i+1] != values[i]+1 then
      return (values[i]+1);
    end if;
      i = i+1;
  end loop;
  return ((select max(id) from rooms)+1);
end;
$$ language plpgsql;

CREATE TABLE rooms (
	 id integer primary key default next_room_id(), 
   owner INTEGER not null,
   foreign key (owner) references users (id)
);

CREATE OR REPLACE FUNCTION user_new_room()
RETURNS trigger AS $$
BEGIN
  insert into rooms (owner) values (NEW.id);
  RETURN NEW;
END;
$$ language plpgsql;

CREATE OR REPLACE TRIGGER user_new_room_trigger 
AFTER INSERT ON users 
FOR EACH ROW 
EXECUTE PROCEDURE user_new_room();

create table guests (
  inviting_room integer not null,
  user_id integer not null,
  permission_level integer not null default 1,
  FOREIGN KEY(inviting_room) REFERENCES rooms (id),
  FOREIGN KEY(user_id) REFERENCES users (id),
  primary key (inviting_room, user_id)
);

CREATE TABLE invites (
   id serial primary key, 
   target INTEGER NOT NULL,
	 inviting_room INTEGER NOT NULL,
   status text not null default 'not aceppeted',
   permission integer not null default 1,
	 FOREIGN KEY (target) REFERENCES users (id),
   FOREIGN KEY (inviting_room) REFERENCES rooms (id)
);

-- caso um convite seja aceito, atualiza a tabela de convidados
-- create or replace function guest_after_invite()
-- returns trigger as $$
-- declare
--   status text = (select (status) from invites where target = old.target);
-- begin
--   if (status = 'aceppeted') then
--     insert into guests (inviting_room, user_id) values (old.inviting_room, old.target);
--   end if;
--   return new;
-- end;
-- $$ language plpgsql;

-- create or replace trigger guest_after_invite
-- after update on invites
-- for each row
-- execute procedure guest_after_invite();

-- ate aqui tambem parece tudo ok. Farei mais testes amanha

CREATE TABLE product_list (
   name varchar,
   origin_room INTEGER,
   FOREIGN KEY (origin_room) REFERENCES rooms (id),
   PRIMARY KEY(name, origin_room)
);

CREATE TABLE products (
   list_name text not null,
   list_room INTEGER not null,
	 name VARCHAR not null,
	 price DECIMAL NOT NULL,
	 description VARCHAR,
   image bytea,
   FOREIGN KEY (list_name, list_room) REFERENCES product_list (name, origin_room),
   primary key (name, list_room)
);

create or replace function next_tab_number (room_id integer)
returns int as $$
declare
  values integer[] := array(select (number) from tabs where room = room_id order by number);
  length integer = array_length(values, 1);
  i int = 1;
begin
  if values[1] is null then 
  	return 1;
  end if;
  while i < length loop
    if values[i+1] != values[i]+1 then
      return (values[i]+1);
    end if;
      i = i+1;
  end loop;
  return ((select max(number) from tabs where room = room_id)+1);
end;
$$ language plpgsql;

CREATE TABLE tabs (
	 number INTEGER not null, 
   room INTEGER,
   pay_value decimal default 0,
   time_maded time default current_time,
   table_number INTEGER default 0, 
	 PRIMARY KEY (room, number),
   FOREIGN KEY (room) REFERENCES rooms (id)
);

create table requests (
  tab_room integer not null,
  tab_number serial not null,
  product_name text not null,
  product_list integer not null, 
  quantity integer not null,
  PRIMARY KEY(product_name, tab_number, tab_room),
  foreign key (product_name, product_list) references products (name, list_room),
  foreign key (tab_room, tab_number) references tabs (room, number) ON DELETE CASCADE ON UPDATE CASCADE
);

-- create or replace function update_request_tab_number()
-- returns trigger as $$
-- begin
--   if new.number != old.number then
--     update requests set tab_number = new.number where tab_room = old.room and tab_number = new.number;
--   elsif new.room != old.room then
--     raise notice 'its not possible to change the rooms tab number';
--     new.room = old.room;
--   end if;
--   
--   return new;
-- end;
-- $$ language plpgsql;

-- create or replace trigger update_request_tab_number_trigger
-- after update on tabs
-- for each row
-- execute procedure update_request_tab_number();

-- adiciona ao valor final da comando de acordo com N pedidos adicionados a comanda 
-- create or replace function add_to_final_tab_value()
-- returns trigger as $$
-- declare
--   total_value decimal = (select (pay_value) from tabs where number = new.tab_number and room = new.tab_room); 
--   new_value decimal = ((
--     select (price) from products join product_list on product_list.name = products.list_name and product_list.origin_room = products.list_room
--     where products.name = new.product_name and products.list_room = new.product_list) * new.quantity);
-- begin
--   UPDATE tabs
--     SET pay_value = total_value + new_value 
--     WHERE number = new.tab_number and room = new.tab_room;
--     return new;
-- end;
-- $$ language plpgsql;

-- create or replace trigger add_to_tab_value_trigger
-- after insert on requests
-- for each row
-- execute procedure add_to_final_tab_value();

-- -- remove X do valor final baseado em N pedidos removidos da comanda
-- create or replace function remove_from_final_tab_value()
-- returns trigger as $$
-- declare
--   total_value decimal = (select (pay_value) from tabs where number = new.tab_number and room = new.tab_room); 
--   new_value decimal = ((
--     select (price) from products join product_list on product_list.name = products.list_name and product_list.origin_room = products.list_room
--     where products.name = new.product_name and products.list_room = new.product_list) * new.quantity);
-- begin
--   UPDATE tabs
--     SET pay_value = total_value - new_value 
--     WHERE number = new.tab_number and room = new.tab_room;
--     return new;
-- end;
-- $$ language plpgsql;

-- create or replace trigger remove_from_tab_value_trigger 
-- after delete on requests
-- for each row
-- execute procedure remove_from_final_tab_value();

-- atualiza o valor final da comando caso haja mudança em N pedidos da comanda
-- create or replace function update_final_tab_value()
-- returns trigger as $$
-- declare
--   total_value decimal = (select (pay_value) from tabs where number = new.tab_number and room = new.tab_room); 
--   old_value decimal = (old.value * old.quantity);
--   new_value decimal = ((
--     select (price) from products join product_list on product_list.name = products.list_name and product_list.origin_room = products.list_room
--     where products.name = new.product_name and products.list_room = new.product_list) * new.quantity);
-- begin
--   UPDATE tabs
--     SET pay_value = ((total_value - old_value) + new_value)
--     WHERE number = new.tab_number and room = new.tab_room;
--     return new;
-- end;
-- $$ language plgsql;

-- create or replace trigger update_final_tab_value_trigger 
-- after update on requests
-- for each row
-- execute procedure update_final_tab_value();

create table payed_tabs (
  id bigserial not null,
  number integer not null,
  room integer not null references rooms (id),
  value decimal(1000, 2) default 0,
  date date not null default current_date,
  table_number INTEGER default 0, 
  PRIMARY KEY(id, room)
);

create table payed_requests (
  room integer not null,
  tab_id serial not null,
  product_name text not null,
  quantity integer not null,
  PRIMARY KEY(product_name, tab_id, room),
  foreign key (room, tab_id) references payed_tabs (room, id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- tabelas experimentais ou para uso especifico dos desenvolvedores
create table how_many_register (
  users integer,
  rooms integer
);

CREATE OR REPLACE FUNCTION users_log_update()
RETURNS trigger AS $$
BEGIN
  update how_many_register set users=((select max(users) from how_many_register)+1); 
  RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER log_user_trigger 
AFTER INSERT ON users 
FOR EACH STATEMENT 
EXECUTE PROCEDURE users_log_update();

create table users_session (
	who int primary key,
	active_room int,
	securePS varchar unique not null,
	foreign key (who) references users (id),
	foreign key (active_room) references rooms (id)
);
