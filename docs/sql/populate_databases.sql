-- Active: 1675298682566@@127.0.0.1@3306@test_vibbra
USE test_vibbra;

INSERT INTO usuarios (name, email, login, password, lat, ing, address, city, state, zip_code)
VALUES 
("usuario de teste 1","user_test_1@gmail.com","user_test_1",
"$2a$10$0Rlmgd6nM99bfX/eqG/bt.ixig5uJwmusZgEJp5yTOaevrolK3ra2", 8000.55, 9000.22, 
"Rua Mockada, número 55", "Blumenau", "SC", 89041480);