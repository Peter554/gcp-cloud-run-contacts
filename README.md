# gcp-cloud-run-contacts

[![CI](https://github.com/Peter554/gcp-cloud-run-contacts/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/Peter554/gcp-cloud-run-contacts/actions/workflows/ci.yml)

An example app using GCP Cloud Run with Cloud SQL.

## Project ID

- Set the GitHub actions secret `GCP_PROJECT_ID`

## Database preparations

https://cloud.google.com/sql/docs/postgres/connect-run

- Create a Cloud SQL postgres instance.
- Connect to the instance (`gcloud sql connect`) and create the contacts table:

```sql
create table if not exists contacts (
    id serial primary key,
    name varchar(100),
    email varchar(100)
);
```

- Obtain the `INSTANCE_CONNECTION_NAME` (`gcloud sql instances describe`).
- Set the GitHub actions secret `GCP_SQL_CONNECTION_NAME`.
- Determine the postgres Data Source Name (DSN): `user=postgres password=<password> database=postgres host=/cloudsql/<INSTANCE_CONNECTION_NAME>`
- Set the GitHub actions secret `GCP_SQL_DSN`.

## Service account

- Create  a service account and obtain a JSON key (https://github.com/google-github-actions/deploy-cloudrun/)
- Set the GitHub actions secret `GCP_CREDENTIALS`.
