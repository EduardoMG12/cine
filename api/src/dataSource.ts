import { DataSource } from "typeorm";
import * as dotenv from "dotenv";

dotenv.config();

export const AppDataSource = new DataSource({
	type: "postgres",
	host: "localhost",
	port: Number(process.env.DATABASE_PORT) || 5432,
	username: process.env.DATABASE_USER || "username",
	password: process.env.DATABASE_PASSWORD || "password",
	database: process.env.DATABASE_NAME || "dbname",
	entities: [__dirname + "/entities/*.entity{.ts,.js}"],
	migrations: [__dirname + "/migrations/*{.ts,.js}"], // i don't no if need src
	synchronize: false,
	logging: true,
});
// see if __dirname is root or is api_organize_my_mind /**/*/migrations/*{.ts,.js}
// past json "migration:run": "bun run typeorm migration:run -d ./dataSource.ts"
