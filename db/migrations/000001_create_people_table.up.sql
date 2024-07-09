CREATE TABLE IF NOT EXISTS public.people (
                                             id SERIAL PRIMARY KEY,
                                             surname VARCHAR(255),
                                             name VARCHAR(255),
                                             patronymic VARCHAR(255),
                                             address VARCHAR(255),
                                             passport_serie VARCHAR(50),
                                             passport_number VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS public.task (
                                           id SERIAL PRIMARY KEY,
                                           peopleId INTEGER NOT NULL,
                                           startT TIMESTAMP NOT NULL,
                                           endT TIMESTAMP,
                                           FOREIGN KEY (peopleId) REFERENCES public.people(id)
);
