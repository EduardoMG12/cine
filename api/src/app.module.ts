import { Module } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { GraphQLModule } from "@nestjs/graphql";
import { ApolloDriver, type ApolloDriverConfig } from "@nestjs/apollo";
import { join } from "path";
import { AppResolver } from "./app.resolver";
import { UserModule } from "./user/user.module";
import { AuthModule } from "./auth/auth.module";

@Module({
	imports: [
		ConfigModule.forRoot({
			isGlobal: true,
			envFilePath: join(__dirname, "..", "..", ".env.api"),
		}),
		TypeOrmModule.forRootAsync({
			imports: [ConfigModule],
			useFactory: async (configService: ConfigService) => {
				const isTestEnv = process.env.NODE_ENV === "test";
				const databaseHost = isTestEnv
					? configService.get<string>("DATABASE_HOST")
					: configService.get<string>("DATABASE_HOST_ON_DOCKER");

				return {
					type: "postgres",
					// host: databaseHost || configService.get<string>("DATABASE_HOST"),
					host: "localhost",
					port: configService.get<number>("DATABASE_PORT"),
					username: configService.get<string>("DATABASE_USERNAME"),
					password: configService.get<string>("DATABASE_PASSWORD"),
					database: configService.get<string>("DATABASE_NAME"),
					entities: [join(__dirname, "**", "*.entity{.ts,.js}")],
					synchronize: true,
				};
			},
			inject: [ConfigService],
		}),
		GraphQLModule.forRoot<ApolloDriverConfig>({
			driver: ApolloDriver,
			autoSchemaFile: join(process.cwd(), "src/schema.gql"),
		}),
		UserModule,
		AuthModule,
	],
	providers: [AppResolver],
})
export class AppModule {}
