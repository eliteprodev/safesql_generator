CREATE FUNCTION foo(bar text) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
DROP FUNCTION foo(text);