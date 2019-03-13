--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2
-- Dumped by pg_dump version 11.2

-- Started on 2019-03-13 20:28:52 +07

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE "Catte";
--
-- TOC entry 3138 (class 1262 OID 16394)
-- Name: Catte; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "Catte" WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C' LC_CTYPE = 'C';


ALTER DATABASE "Catte" OWNER TO postgres;

\connect "Catte"

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 196 (class 1259 OID 16395)
-- Name: rooms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rooms (
    roomid character varying NOT NULL,
    amount bigint,
    isactive boolean NOT NULL,
    numplayer integer,
    host character varying
);


ALTER TABLE public.rooms OWNER TO postgres;

--
-- TOC entry 197 (class 1259 OID 16401)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    userid character varying NOT NULL,
    username character varying NOT NULL,
    email character varying,
    source character varying(50) NOT NULL,
    password character varying,
    lastcheckin date,
    amount bigint,
    user3rdid character varying,
    dateofbirth date
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 3131 (class 0 OID 16395)
-- Dependencies: 196
-- Data for Name: rooms; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('7', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('1', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('2', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('3', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('4', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('5', 5000, false, 0, NULL);
INSERT INTO public.rooms (roomid, amount, isactive, numplayer, host) VALUES ('6', 10000, false, 0, NULL);


--
-- TOC entry 3132 (class 0 OID 16401)
-- Dependencies: 197
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (userid, username, email, source, password, lastcheckin, amount, user3rdid, dateofbirth) VALUES ('21c13177-1793-461d-9747-a4868c441bc4', 'Nguyen Cat Dinh', NULL, 'App', '$2a$10$LYYmD3u3mhC3PYQiR5JC9eocGmAyhrTEApRHwG7/EwZhjdaE8v3PC', NULL, 50000, NULL, NULL);


--
-- TOC entry 3007 (class 2606 OID 16408)
-- Name: rooms roomid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT roomid PRIMARY KEY (roomid);


--
-- TOC entry 3009 (class 2606 OID 16410)
-- Name: users userid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT userid PRIMARY KEY (userid);


-- Completed on 2019-03-13 20:28:52 +07

--
-- PostgreSQL database dump complete
--

