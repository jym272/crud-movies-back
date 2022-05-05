--
-- PostgreSQL database dump
--

-- Dumped from database version 13.6 (Ubuntu 13.6-0ubuntu0.21.10.1)
-- Dumped by pg_dump version 13.6 (Ubuntu 13.6-0ubuntu0.21.10.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: favorite_movies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.favorite_movies (
    id integer NOT NULL,
    user_id integer NOT NULL,
    movie_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: favorite_movies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.favorite_movies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: favorite_movies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.favorite_movies_id_seq OWNED BY public.favorite_movies.id;


--
-- Name: genres; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.genres (
    id integer NOT NULL,
    genre_name character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: genres_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.genres_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.genres_id_seq OWNED BY public.genres.id;


--
-- Name: movies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.movies (
    id integer NOT NULL,
    title character varying,
    description text,
    year integer,
    release_date date,
    runtime integer,
    rating integer,
    mpaa_rating character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    poster character varying,
    user_id integer
);


--
-- Name: movies_genres; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.movies_genres (
    id integer NOT NULL,
    movie_id integer,
    genre_id integer,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: movies_genres_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.movies_genres_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: movies_genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.movies_genres_id_seq OWNED BY public.movies_genres.id;


--
-- Name: movies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.movies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: movies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.movies_id_seq OWNED BY public.movies.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: favorite_movies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.favorite_movies ALTER COLUMN id SET DEFAULT nextval('public.favorite_movies_id_seq'::regclass);


--
-- Name: genres id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.genres ALTER COLUMN id SET DEFAULT nextval('public.genres_id_seq'::regclass);


--
-- Name: movies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies ALTER COLUMN id SET DEFAULT nextval('public.movies_id_seq'::regclass);


--
-- Name: movies_genres id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies_genres ALTER COLUMN id SET DEFAULT nextval('public.movies_genres_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: favorite_movies; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.favorite_movies (id, user_id, movie_id, created_at, updated_at) FROM stdin;
183	17	8	2022-05-02 00:20:28.481922	2022-05-02 00:20:28.482323
184	17	35	2022-05-02 00:20:30.414171	2022-05-02 00:20:30.414604
251	19	45	2022-05-03 02:52:19.989193	2022-05-03 02:52:19.989392
253	19	48	2022-05-03 03:09:50.321918	2022-05-03 03:09:50.322155
31	1	40	2022-05-01 07:06:18.454216	2022-05-01 07:06:18.45453
32	1	21	2022-05-01 07:06:36.456182	2022-05-01 07:06:36.456236
34	1	4	2022-05-01 17:46:37.516703	2022-05-01 17:46:37.51675
256	19	49	2022-05-03 03:09:53.583347	2022-05-03 03:09:53.583559
257	20	48	2022-05-03 14:43:46.592115	2022-05-03 14:43:46.592179
37	1	26	2022-05-01 18:00:25.075297	2022-05-01 18:00:25.075666
258	20	46	2022-05-03 20:27:38.881553	2022-05-03 20:27:38.881922
39	1	12	2022-05-01 21:44:47.444051	2022-05-01 21:44:47.444105
40	1	25	2022-05-01 22:47:01.473835	2022-05-01 22:47:01.474223
259	20	47	2022-05-03 20:27:41.304918	2022-05-03 20:27:41.30515
42	1	35	2022-05-01 22:47:55.298658	2022-05-01 22:47:55.299068
211	10	12	2022-05-02 16:02:34.572097	2022-05-02 16:02:34.572166
218	10	25	2022-05-02 17:35:38.386157	2022-05-02 17:35:38.38643
220	10	21	2022-05-02 18:06:32.695771	2022-05-02 18:06:32.69584
\.


--
-- Data for Name: genres; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.genres (id, genre_name, created_at, updated_at) FROM stdin;
1	Drama	2021-05-17 00:00:00	2021-05-17 00:00:00
2	Crime	2021-05-17 00:00:00	2021-05-17 00:00:00
3	Action	2021-05-17 00:00:00	2021-05-17 00:00:00
4	Comic Book	2021-05-17 00:00:00	2021-05-17 00:00:00
5	Sci-Fi	2021-05-17 00:00:00	2021-05-17 00:00:00
6	Mystery	2021-05-17 00:00:00	2021-05-17 00:00:00
7	Adventure	2021-05-17 00:00:00	2021-05-17 00:00:00
8	Comedy	2021-05-17 00:00:00	2021-05-17 00:00:00
9	Romance	2021-05-17 00:00:00	2021-05-17 00:00:00
10	Thriller	2022-05-02 18:51:56	2022-05-02 18:53:09
11	Fantasy	2022-05-03 02:35:35	2022-05-03 02:35:38
12	History	2022-05-03 02:35:59	2022-05-03 02:36:04
\.


--
-- Data for Name: movies; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.movies (id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at, poster, user_id) FROM stdin;
49	Nightcrawler	NIGHTCRAWLER is a thriller set in the nocturnal underbelly of contemporary Los Angeles. Jake Gyllenhaal stars as Lou Bloom, a driven young man desperate for work who discovers the high-speed world of L.A. crime journalism. Finding a group of freelance camera crews who film crashes, fires, murder and other mayhem, Lou muscles into the cut-throat, dangerous realm of nightcrawling - where each police siren wail equals a possible windfall and victims are converted into dollars and cents. Aided by Rene Russo as Nina, a veteran of the blood-sport that is local TV news, Lou blurs the line between observer and participant to become the star of his own story.	2014	2014-10-31	117	4	R	2022-05-03 03:09:18.647414	2022-05-03 03:09:18.647415	/j9HrX8f7GbZQm1BrBiR40uFQZSb.jpg	19
50	The Seventh Seal	A Knight and his squire are home from the crusades. Black Death is sweeping their country. As they approach home, Death appears to the knight and tells him it is his time. The knight challenges Death to a chess game for his life. The Knight and Death play as the cultural turmoil envelopes the people around them as they try, in different ways, to deal with the upheaval the plague has caused.	1958	1958-10-13	96	5	PG13	2022-05-03 03:16:04.918261	2022-05-03 03:16:04.918261	/wcZ21zrOsy0b52AfAF50XpTiv75.jpg	20
51	12 Angry Men	The defense and the prosecution have rested, and the jury is filing into the jury room to decide if a young man is guilty or innocent of murdering his father. What begins as an open-and-shut case of murder soon becomes a detective story that presents a succession of clues creating doubt, and a mini-drama of each of the jurors' prejudices and preconceptions about the trial, the accused, AND each other. Based on the play, all of the action takes place on the stage of the jury room.	1957	1957-10-04	96	5	PG13	2022-05-04 21:37:04.856379	2022-05-04 21:37:04.856379	/6PlhouMCYktJmdFwS9XtqRZaTqc.jpg	20
52	Face/Off	Sean Archer, a very tough, rugged FBI Agent, is still grieving for his dead son Michael. Archer believes that his son's killer is his sworn enemy and a very powerful criminal, Castor Troy. One day, Archer has finally cornered Castor, however, their fight has knocked out Troy cold. As Archer finally breathes easy over the capture of his enemy, he finds out that Troy has planted a bomb that will destroy the entire city of Los Angeles and all of its inhabitants. Unfortunately the only other person who knows its location is Castor's brother Pollux, and he refuses to talk. The solution, a special operation doctor that can cut off people's faces, and can place a person's face onto another person.	1997	1997-02-06	138	3	R	2022-05-04 21:38:39.781451	2022-05-04 21:38:39.781451	/aYZDYdiMym5xEkSs4BrJMXFdZ9g.jpg	20
44	The Machinist	Trevor Reznik is a machinist in a factory. An extreme case of insomnia has led to him not sleeping in a year, and his body withering away to almost nothing. He has an obsessive compulsion to write himself reminder notes and keep track of his dwindling weight, both scribbled on yellow stickies in his apartment. The only person he lets into his life in an emotional sense is Stevie, a prostitute, although he has an infatuation with Maria, a single mother waitress working in an airport diner. His co-workers don't associate with and mistrust him because of not knowing what is going on in his life that has led to his emaciated physical appearance. A workplace incident further alienates him with his coworkers, and in conjunction with some unfamiliar pieces of paper he finds in his apartment.	2004	2004-12-03	101	4	R	2022-05-03 02:11:11.156031	2022-05-03 02:25:00.88506	/ukMORT211Xi7I8mvtg9aXgzs374.jpg	19
45	The Northman	From visionary director Robert Eggers comes The Northman, an action-filled epic that follows a young Viking prince on his quest to avenge his father's murder.	2022	2022-02-22	137	5	R	2022-05-03 02:41:10.863517	2022-05-03 02:41:10.863517	/zhLKlUaF1SEpO58ppHIAyENkwgw.jpg	19
46	Dune	A mythic and emotionally charged hero's journey, "Dune" tells the story of Paul Atreides, a brilliant and gifted young man born into a great destiny beyond his understanding, who must travel to the most dangerous planet in the universe to ensure the future of his family and his people. As malevolent forces explode into conflict over the planet's exclusive supply of the most precious resource in existence-a commodity capable of unlocking humanity's greatest potential-only those who can conquer their fear will survive.	2021	2021-10-22	155	5	PG13	2022-05-03 03:00:02.209586	2022-05-03 03:00:02.209586	/d5NXSklXo0qyIYkgV94XAgMIckC.jpg	19
47	Logan	In 2029 the mutant population has shrunken significantly due to genetically modified plants designed to reduce mutant powers and the X-Men have disbanded. Logan, whose power to self-heal is dwindling, has surrendered himself to alcohol and now earns a living as a chauffeur. He takes care of the ailing old Professor X whom he keeps hidden away. One day, a female stranger asks Logan to drive a girl named Laura to the Canadian border. At first he refuses, but the Professor has been waiting for a long time for her to appear. Laura possesses an extraordinary fighting prowess and is in many ways like Wolverine. She is pursued by sinister figures working for a powerful corporation; this is because they made her, with Logan's DNA.	2017	2017-03-03	137	4	R	2022-05-03 03:04:36.564477	2022-05-03 03:04:36.564477	/fnbjcRDYn6YviCcePDnGdyAkYsB.jpg	19
48	American Psycho	It's the late 1980s. Twenty-seven year old Wall Streeter Patrick Bateman travels among a closed network of the proverbial beautiful people, that closed network in only they able to allow others like themselves in in a feeling of superiority. Patrick has a routinized morning regimen to maintain his appearance of attractiveness and fitness. He, like those in his network, are vain, narcissistic, egomaniacal and competitive, always having to one up everyone else in that presentation of oneself, but he, unlike the others, realizes that, for himself, all of these are masks to hide what is truly underneath, someone/something inhuman in nature.	2000	2000-04-01	102	3	R	2022-05-03 03:07:22.258754	2022-05-03 03:07:22.258754	/3ddHhfMlZHZCefHDeaP8FzSoH4Y.jpg	19
53	Solyaris	The Solaris mission has established a base on a planet that appears to host some kind of intelligence, but the details are hazy and very secret. After the mysterious demise of one of the three scientists on the base, the main character is sent out to replace him. He finds the station run-down and the two remaining scientists cold and secretive. When he also encounters his wife who has been dead for ten years, he begins to appreciate the baffling nature of the alien intelligence.	1972	1972-09-26	163	5	PG	2022-05-04 21:40:46.811975	2022-05-04 21:45:05.755416	/nVS59KVITySXPFV7ABRJrEYVJ3d.jpg	19
54	Сталкер	In a small, unnamed country there's an area called the Zone. It's an unusual area, and within its a place known as the Room, where it's believed wishes are granted. The government declared The Zone a no-go area and have sealed it off. This hasn't stopped people from entering the Zone. A writer, and a professor, want to reach the Zone. Their guide - a man known as a stalker, has a special relation with the Zone.	1980	1980-04-17	162	5	PG13	2022-05-04 21:46:20.416294	2022-05-04 21:47:20.525602	/lUE0Bp7wH0EterJ44qMRsqtKFnp.jpg	19
55	Terminator 2: Judgment Day	Over 10 years have passed since the first machine called The Terminator tried to kill Sarah Connor and her unborn son, John. The man who will become the future leader of the human resistance against the Machines is now a healthy young boy. However, another Terminator, called the T-1000, is sent back through time by the supercomputer Skynet. This new Terminator is more advanced and more powerful than its predecessor and its mission is to kill John Connor when he's still a child. However, Sarah and John do not have to face the threat of the T-1000 alone. Another Terminator (identical to the same model that tried and failed to kill Sarah Connor in 1984) is also sent back through time to protect them. Now, the battle for tomorrow has begun.	1991	1991-07-03	137	4	R	2022-05-04 21:49:05.540303	2022-05-04 21:49:05.540303	/weVXMD5QBGeQil4HEATZqAkXeEc.jpg	19
\.


--
-- Data for Name: movies_genres; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.movies_genres (id, movie_id, genre_id, created_at, updated_at) FROM stdin;
617	44	10	2022-05-03 02:25:00.920979	\N
618	44	1	2022-05-03 02:25:00.92352	\N
619	45	10	2022-05-03 02:41:10.865452	\N
620	45	11	2022-05-03 02:41:10.865935	\N
621	45	12	2022-05-03 02:41:10.867479	\N
622	45	1	2022-05-03 02:41:10.867904	\N
623	45	3	2022-05-03 02:41:10.869401	\N
624	45	7	2022-05-03 02:41:10.869846	\N
625	46	3	2022-05-03 03:00:02.235747	\N
626	46	5	2022-05-03 03:00:02.247533	\N
627	46	7	2022-05-03 03:00:02.24935	\N
628	46	1	2022-05-03 03:00:02.249959	\N
629	47	1	2022-05-03 03:04:36.580155	\N
630	47	3	2022-05-03 03:04:36.593292	\N
631	47	5	2022-05-03 03:04:36.594742	\N
632	47	10	2022-05-03 03:04:36.5968	\N
633	48	1	2022-05-03 03:07:22.282771	\N
634	48	2	2022-05-03 03:07:22.295787	\N
635	49	10	2022-05-03 03:09:18.65116	\N
636	49	1	2022-05-03 03:09:18.652335	\N
637	49	2	2022-05-03 03:09:18.654195	\N
638	50	11	2022-05-03 03:16:04.935579	\N
639	50	1	2022-05-03 03:16:04.950216	\N
640	51	1	2022-05-04 21:37:04.863197	\N
641	51	2	2022-05-04 21:37:04.866825	\N
642	52	2	2022-05-04 21:38:39.786253	\N
643	52	3	2022-05-04 21:38:39.788541	\N
644	52	5	2022-05-04 21:38:39.78944	\N
645	52	10	2022-05-04 21:38:39.79144	\N
670	53	1	2022-05-04 21:45:05.797171	\N
671	53	5	2022-05-04 21:45:05.798851	\N
672	53	6	2022-05-04 21:45:05.799353	\N
677	54	1	2022-05-04 21:47:20.550171	\N
678	54	5	2022-05-04 21:47:20.55221	\N
679	55	3	2022-05-04 21:49:05.551498	\N
680	55	5	2022-05-04 21:49:05.567982	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, username, password, created_at, updated_at) FROM stdin;
19	joclavijo@fi.uba.ar	google	2022-05-02 18:47:54.364468	2022-05-02 18:47:54.364468
20	jym272@gmail.com	google	2022-05-03 03:15:07.441797	2022-05-03 03:15:07.441798
\.


--
-- Name: favorite_movies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.favorite_movies_id_seq', 260, true);


--
-- Name: genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.genres_id_seq', 12, true);


--
-- Name: movies_genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.movies_genres_id_seq', 680, true);


--
-- Name: movies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.movies_id_seq', 55, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_id_seq', 20, true);


--
-- Name: favorite_movies favorite_movies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.favorite_movies
    ADD CONSTRAINT favorite_movies_pkey PRIMARY KEY (id);


--
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);


--
-- Name: movies_genres movies_genres_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT movies_genres_pkey PRIMARY KEY (id);


--
-- Name: movies movies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: movies_genres fk_movie_genries_genre_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT fk_movie_genries_genre_id FOREIGN KEY (genre_id) REFERENCES public.genres(id);


--
-- Name: movies_genres fk_movie_genries_movie_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT fk_movie_genries_movie_id FOREIGN KEY (movie_id) REFERENCES public.movies(id);


--
-- PostgreSQL database dump complete
--

