import { check } from 'k6';
import { Client } from 'k6/x/sql';

export const options = {
  vus: 1,
  iterations: 1,
};

const db = new Client();

export function setup() {
  const postgresUser = __ENV.POSTGRES_USER;
  const postgresPassword = __ENV.POSTGRES_PASSWORD;
  const postgresDb = __ENV.POSTGRES_DB;
  const postgresHost = __ENV.POSTGRES_HOST;
  const postgresPort = __ENV.POSTGRES_PORT;

  if (!postgresUser || !postgresPassword || !postgresDb || !postgresHost || !postgresPort) {
    throw new Error(
      'Missing one or more required environment variables: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST, POSTGRES_PORT'
    );
  }

  const connectionString = `postgres://${postgresUser}:${postgresPassword}@${postgresHost}:${postgresPort}/${postgresDb}?sslmode=disable`;

  db.connect(connectionString);

  // Create a table if it doesn't exist
  const createTableQuery = `
    CREATE TABLE IF NOT EXISTS users (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      email VARCHAR(255) UNIQUE NOT NULL
    );
  `;
  db.exec(createTableQuery);

  return { db: db };
}

export default function (data) {
  const db = data.db;
  // Insert some data
  const insertQuery = `
    INSERT INTO users (name, email) VALUES
      ('John Doe', 'john.doe@example.com'),
      ('Jane Smith', 'jane.smith@example.com'),
      ('Peter Jones', 'peter.jones@example.com')
    ON CONFLICT (email) DO NOTHING;
  `;

  const result = db.exec(insertQuery);
  check(result, { 'Data inserted successfully': (r) => r.length === 0 });
}

export function teardown(data) {
  const db = data.db;
  db.close();
}
