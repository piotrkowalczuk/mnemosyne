SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;


CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;


--
-- Mnemosyne Session
--
CREATE TABLE mnemosyne_session (
    id character varying(128) NOT NULL,
    data json NOT NULL,
    expire_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY mnemosyne_session
    ADD CONSTRAINT mnemosyne_session_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mnemosyne_session
    ADD CONSTRAINT mnemosyne_session_id_key UNIQUE (id);
