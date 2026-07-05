CREATE SEQUENCE IF NOT EXISTS public.goadmin_menu_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_operation_log_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_site_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_permissions_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_roles_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_session_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS public.goadmin_users_myid_seq START WITH 1 INCREMENT BY 1 NO MINVALUE MAXVALUE 99999999 CACHE 1;

CREATE TABLE IF NOT EXISTS public.goadmin_menu (
    id integer DEFAULT nextval('public.goadmin_menu_myid_seq'::regclass) PRIMARY KEY,
    parent_id integer DEFAULT 0 NOT NULL,
    type integer DEFAULT 0,
    "order" integer DEFAULT 0 NOT NULL,
    title character varying(50) NOT NULL,
    header character varying(100),
    plugin_name character varying(100) NOT NULL DEFAULT '',
    icon character varying(50) NOT NULL,
    uri character varying(3000) NOT NULL,
    uuid character varying(100),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_operation_log (
    id integer DEFAULT nextval('public.goadmin_operation_log_myid_seq'::regclass) PRIMARY KEY,
    user_id integer NOT NULL,
    path character varying(255) NOT NULL,
    method character varying(10) NOT NULL,
    ip character varying(45) NOT NULL,
    input text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_site (
    id integer DEFAULT nextval('public.goadmin_site_myid_seq'::regclass) PRIMARY KEY,
    key character varying(100) NOT NULL,
    value text NOT NULL,
    type integer DEFAULT 0,
    description character varying(3000),
    state integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_permissions (
    id integer DEFAULT nextval('public.goadmin_permissions_myid_seq'::regclass) PRIMARY KEY,
    name character varying(50) NOT NULL,
    slug character varying(50) NOT NULL,
    http_method character varying(255),
    http_path text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_role_menu (
    role_id integer NOT NULL,
    menu_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    PRIMARY KEY (role_id, menu_id)
);

CREATE TABLE IF NOT EXISTS public.goadmin_role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS public.goadmin_role_users (
    role_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    PRIMARY KEY (role_id, user_id)
);

CREATE TABLE IF NOT EXISTS public.goadmin_roles (
    id integer DEFAULT nextval('public.goadmin_roles_myid_seq'::regclass) PRIMARY KEY,
    name character varying NOT NULL,
    slug character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_session (
    id integer DEFAULT nextval('public.goadmin_session_myid_seq'::regclass) PRIMARY KEY,
    sid character varying(50) NOT NULL,
    "values" character varying(3000) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.goadmin_user_permissions (
    user_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE IF NOT EXISTS public.goadmin_users (
    id integer DEFAULT nextval('public.goadmin_users_myid_seq'::regclass) PRIMARY KEY,
    username character varying(100) NOT NULL,
    password character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    avatar character varying(255),
    remember_token character varying(100),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

INSERT INTO public.goadmin_menu (id, parent_id, type, "order", title, plugin_name, header, icon, uri, created_at, updated_at) VALUES
    (1, 0, 1, 2, 'Admin', '', NULL, 'fa-tasks', '', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 1, 1, 2, 'Users', '', NULL, 'fa-users', '/info/manager', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (3, 1, 1, 3, 'Roles', '', NULL, 'fa-user', '/info/roles', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (4, 1, 1, 4, 'Permission', '', NULL, 'fa-ban', '/info/permission', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (5, 1, 1, 5, 'Menu', '', NULL, 'fa-bars', '/menu', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (6, 1, 1, 6, 'Operation log', '', NULL, 'fa-history', '/info/op', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (7, 0, 1, 1, 'Dashboard', '', NULL, 'fa-bar-chart', '/', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (8, 0, 1, 3, 'FlowerDoro', '', NULL, 'fa-leaf', '', now(), now()),
    (9, 8, 1, 1, 'Flower Catalog', '', NULL, 'fa-pagelines', '/info/flowers', now(), now()),
    (10, 8, 1, 2, 'Garden Flowers', '', NULL, 'fa-envira', '/info/garden_flowers', now(), now()),
    (11, 8, 1, 3, 'Focus Sessions', '', NULL, 'fa-clock-o', '/info/focus_sessions', now(), now()),
    (12, 8, 1, 4, 'App Users', '', NULL, 'fa-mobile', '/info/users', now(), now()),
    (13, 8, 1, 5, 'Gardens', '', NULL, 'fa-tree', '/info/gardens', now(), now()),
    (14, 8, 1, 6, 'User Settings', '', NULL, 'fa-sliders', '/info/user_settings', now(), now())
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.goadmin_permissions (id, name, slug, http_method, http_path, created_at, updated_at) VALUES
    (1, 'All permission', '*', '', '*', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 'Dashboard', 'dashboard', 'GET,PUT,POST,DELETE', '/', '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.goadmin_roles (id, name, slug, created_at, updated_at) VALUES
    (1, 'Administrator', 'administrator', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 'Operator', 'operator', '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.goadmin_role_menu (role_id, menu_id, created_at, updated_at) VALUES
    (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (1, 7, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (1, 8, now(), now()),
    (1, 9, now(), now()),
    (1, 10, now(), now()),
    (1, 11, now(), now()),
    (1, 12, now(), now()),
    (1, 13, now(), now()),
    (1, 14, now(), now()),
    (2, 7, '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (role_id, menu_id) DO NOTHING;

INSERT INTO public.goadmin_role_permissions (role_id, permission_id, created_at, updated_at) VALUES
    (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (1, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO public.goadmin_role_users (role_id, user_id, created_at, updated_at) VALUES
    (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (role_id, user_id) DO NOTHING;

INSERT INTO public.goadmin_user_permissions (user_id, permission_id, created_at, updated_at) VALUES
    (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (user_id, permission_id) DO NOTHING;

INSERT INTO public.goadmin_users (id, username, password, name, avatar, remember_token, created_at, updated_at) VALUES
    (1, 'admin', '$2a$10$OxWYJJGTP2gi00l2x06QuOWqw5VR47MQCJ0vNKnbMYfrutij10Hwe', 'admin', '', 'tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
    (2, 'operator', '$2a$10$rVqkOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.', 'Operator', '', NULL, '2019-09-10 00:00:00', '2019-09-10 00:00:00')
ON CONFLICT (id) DO NOTHING;

SELECT setval('public.goadmin_menu_myid_seq', GREATEST((SELECT COALESCE(MAX(id), 14) FROM public.goadmin_menu), 14), true);
SELECT setval('public.goadmin_permissions_myid_seq', 2, true);
SELECT setval('public.goadmin_roles_myid_seq', 2, true);
SELECT setval('public.goadmin_users_myid_seq', 2, true);
