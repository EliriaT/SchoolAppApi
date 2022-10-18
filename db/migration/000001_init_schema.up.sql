CREATE TABLE "User" (
                        "id" bigserial PRIMARY KEY,
                        "email" varchar UNIQUE NOT NULL,
                        "password" varchar NOT NULL,
                        "last_name" varchar NOT NULL,
                        "first_name" varchar NOT NULL,
                        "gender" varchar NOT NULL,
                        "phone_number" varchar,
                        "domicile" varchar,
                        "birth_date" date NOT NULL
);

CREATE TABLE "Role" (
                        "id" bigserial PRIMARY KEY,
                        "name" varchar NOT NULL
);

CREATE TABLE "UserRoles" (
                             "id" bigserial PRIMARY KEY,
                             "user_id" bigint,
                             "role_id" bigint,
                             "school_id" bigint
);

CREATE TABLE "UserRoleClass" (
                                 "id" bigserial PRIMARY KEY,
                                 "user_role_id" bigint,
                                 "class_id" bigint
);

CREATE TABLE "School" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar NOT NULL UNIQUE,
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Course" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar NOT NULL,
                          "semester_id" int,
                          "class_id" int,
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Class" (
                         "id" bigserial PRIMARY KEY,
                         "name" varchar NOT NULL,
                         "diriginte" bigint,
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
                          "week_day" date,
                          "classroom" varchar,
                          "created_by" bigint,
                          "updated_by" bigint,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE "Marks" (
                         "id" bigserial PRIMARY KEY,
                         "lesson_id" bigint,
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
                            "start_date" date,
                            "end_date" date,
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

CREATE INDEX ON "Marks" ("student_id");

CREATE INDEX ON "Marks" ("lesson_id", "student_id");

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

ALTER TABLE "Class" ADD FOREIGN KEY ("diriginte") REFERENCES "UserRoles" ("id");

ALTER TABLE "Class" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Class" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("teacher_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Lesson" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("lesson_id") REFERENCES "Lesson" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("student_id") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Marks" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Semester" ADD FOREIGN KEY ("created_by") REFERENCES "UserRoles" ("id");

ALTER TABLE "Semester" ADD FOREIGN KEY ("updated_by") REFERENCES "UserRoles" ("id");