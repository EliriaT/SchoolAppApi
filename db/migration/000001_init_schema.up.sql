CREATE TABLE "User" (
                        "id" bigserial PRIMARY KEY,
                        "email" varchar UNIQUE NOT NULL,
                        "password" varchar NOT NULL,
                        "totp_secret" varchar NOT NULL,
                        "last_name" varchar NOT NULL,
                        "first_name" varchar NOT NULL,
                        "gender" varchar NOT NULL,
                        "phone_number" varchar,
                        "domicile" varchar,
                        "birth_date" date,
                        "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00',
                        "created_at" timestamp NOT NULL DEFAULT (now()),
                        "updated_at" timestamp

);

CREATE TABLE "Role" (
                        "id" bigserial PRIMARY KEY,
                        "name" varchar NOT NULL
);

CREATE TABLE "UserRoles" (
                             "id" bigserial PRIMARY KEY,
                             "user_id" bigint NOT NULL,
                             "role_id" bigint NOT NULL,
                             "school_id" bigint NOT NULL
);

CREATE TABLE "UserRoleClass" (
                                 "id" bigserial PRIMARY KEY,
                                 "user_role_id" bigint NOT NULL,
                                 "class_id" bigint NOT NULL
);

CREATE TABLE "School" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar UNIQUE NOT NULL,
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Course" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar NOT NULL,
                          teacher_id bigint,
                          "semester_id" int,
                          "class_id" int,
                          "dates" date[],
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Class" (
                         "id" bigserial PRIMARY KEY,
                         "name" varchar  NOT NULL,
                         "head_teacher" bigint UNIQUE NOT NULL,
                         "created_by" bigint,
                         "updated_by" bigint,
                         "created_at" timestamp DEFAULT (now()),
                         "updated_at" timestamp
);

CREATE TABLE "Lesson" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar NOT NULL,
                          "course_id" bigint,
                          "teacher_id" bigint,
                          "start_hour" time,
                          "end_hour" time,
                          "week_day" varchar,
                          "classroom" varchar,
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Marks" (
                         "id" bigserial PRIMARY KEY,
                         "course_id" bigint,
                         "mark_date" date,
                         "is_absent" boolean,
                         "mark" int,
                         "student_id" bigint,
                         "created_by" bigint,
                         "updated_by" bigint,
                         "created_at" timestamp DEFAULT (now()),
                         "updated_at" timestamp
);

CREATE TABLE "Semester" (
                            "id" bigserial PRIMARY KEY,
                            "name" varchar,
                            "start_date" time,
                            "end_date" time,
                            "created_by" bigint,
                            "updated_by" bigint,
                            "created_at" timestamp DEFAULT (now()),
                            "updated_at" timestamp
);

CREATE INDEX ON "User" ("email");

CREATE INDEX ON "User" ("last_name", "first_name");

CREATE INDEX ON "UserRoles" ("school_id");

CREATE INDEX ON "UserRoleClass" ("class_id", "user_role_id");

CREATE INDEX ON "Course" ("class_id");

CREATE INDEX ON "Course" ("class_id", "semester_id");

CREATE INDEX ON "Lesson" ("course_id");

CREATE INDEX ON "Lesson" ("course_id", "teacher_id");

COMMENT ON COLUMN "Marks"."mark" IS 'Bigger than 0, lower than 11';

ALTER TABLE "UserRoles" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id");

ALTER TABLE "UserRoles" ADD FOREIGN KEY ("role_id") REFERENCES "Role" ("id");

ALTER TABLE "UserRoles" ADD FOREIGN KEY ("school_id") REFERENCES "School" ("id");

ALTER TABLE "UserRoleClass" ADD FOREIGN KEY ("user_role_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "UserRoleClass" ADD FOREIGN KEY ("class_id") REFERENCES "Class" ("id");

ALTER TABLE "School" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "School" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Course" ADD FOREIGN KEY ("semester_id") REFERENCES "Semester" ("id");

ALTER TABLE "Course" ADD FOREIGN KEY ("class_id") REFERENCES "Class" ("id");

ALTER TABLE "Course" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Course" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");


ALTER TABLE "Course" ADD FOREIGN KEY ("teacher_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "Class" ADD FOREIGN KEY ("head_teacher") REFERENCES "UserRoles" ("id");

ALTER TABLE "Class" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Class" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("teacher_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("student_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Semester" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Semester" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

INSERT INTO "Role"(name) VALUES ('Admin'), ('Director'),  ('School_Manager'),
                                ('Head_Teacher'),  ('Teacher'), ('Student') ;

ALTER TABLE "Class" ADD CONSTRAINT "unique_classes" UNIQUE ("name", "head_teacher");