BEGIN;
CREATE TABLE "public"."tbl_stocks"(
    "id" uuid NOT NULL,
    "created_at" TIMESTAMP not null,
    "paper_name" VARCHAR(255),
    "company_name" VARCHAR(255),
    "daily_rate" VARCHAR(25),
    "market_value" float8,
    PRIMARY KEY ("id")
);
COMMIT;