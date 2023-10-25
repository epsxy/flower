-- posts table
CREATE TABLE public.posts (
  id uuid NOT NULL,
  name VARCHAR(34) NOT NULL,
  description VARCHAR(514),
  created_at timestamp without time zone NOT NULL,
  user_id BIGINT NOT NULL,
);

ALTER TABLE ONLY public.posts
  ADD CONSTRAINT posts_pkey PRIMARY KEY (id);

-- users table
CREATE TABLE public.users (
  name VARCHAR(34) NOT NULL,
  id BIGINT NOT NULL AUTO_INCREMENT,
);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- comments table
CREATE TABLE public.comments (
  user_id BIGINT NOT NULL,
  content VARCHAR(514),
  post_id BIGINT NOT NULL AUTO_INCREMENT,
);

-- autonomous1 table
CREATE TABLE public.autonomous1 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- autonomous2 table
CREATE TABLE public.autonomous2 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- autonomous3 table
CREATE TABLE public.autonomous3 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked1 table
CREATE TABLE public.linked1 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked2 table
CREATE TABLE public.linked2 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked3 table
CREATE TABLE public.linked3 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked4 table
CREATE TABLE public.linked4 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked5 table
CREATE TABLE public.linked5 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- linked6 table
CREATE TABLE public.linked6 (
  id BIGINT NOT NULL,
  central_id BIGINT NOT NULL,
);

-- central link table
CREATE TABLE public.centrallink (
  id BIGINT NOT NULL,
);

ALTER TABLE ONLY public.linked1
  ADD CONSTRAINT linked1_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked2
  ADD CONSTRAINT linked2_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked3
  ADD CONSTRAINT linked3_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked4
  ADD CONSTRAINT linked4_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked5
  ADD CONSTRAINT linked5_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked6
  ADD CONSTRAINT linked6_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.linked1
  ADD CONSTRAINT linked1_pkey PRIMARY KEY (id);

ALTER TABLE public.linked1 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE public.linked2 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE public.linked3 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE public.linked4 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE public.linked5 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE public.linked6 ADD CONSTRAINT
  FOREIGN KEY (central_id)
  REFERENCES public.centrallink(id);

ALTER TABLE ONLY public.comments
  ADD CONSTRAINT comments_pkey PRIMARY KEY (user_id, post_id);

ALTER TABLE public.comments ADD CONSTRAINT
  FOREIGN KEY (post_id)
  REFERENCES public.posts(id);

ALTER TABLE public.comments ADD CONSTRAINT
  FOREIGN KEY (user_id)
  REFERENCES public.users(id);

ALTER TABLE public.posts ADD CONSTRAINT
  FOREIGN KEY (user_id)
  REFERENCES public.users(id);
