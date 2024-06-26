CREATE TABLE IF NOT EXISTS nasa_apod (
	id SERIAL PRIMARY KEY,
	copyright TEXT,
	apod_date DATE UNIQUE,
	explanation TEXT,
	hdurl TEXT,
	media_type TEXT,
	service_version TEXT,
	title TEXT,
	url TEXT
);
