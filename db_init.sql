use defaultdb;

create table energy_history (
  reporting_date date,
  energy_kw float,
  PRIMARY KEY (reporting_date)
);

create table weather_history (
    reporting_date date,
    cloudiness int,
    day_minutes int,
    PRIMARY KEY (reporting_date)
);